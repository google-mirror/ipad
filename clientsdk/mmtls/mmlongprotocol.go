package mmtls

import (
	"errors"
	"github.com/lunny/log"
	"sync/atomic"

	"feiyu.com/wx/clientsdk/baseutils"
	"golang.org/x/net/proxy"
)

// InitMMTLSInfoLong 初始化MMTLS信息通过长链接
func InitMMTLSInfoLong(dialer proxy.Dialer, longHostName string, longPort string, shortHostName string, pskList []*Psk) (*MMInfo, error) {
	// 初始化MMInfo
	mmInfo := CreateNewMMInfo()
	mmInfo.ShortHost = shortHostName
	mmInfo.LongHost = longHostName
	mmInfo.LONGPort = longPort
	mmInfo.Dialer = dialer
	// 第一次设置为空，后面设置成前面握手返回的Psk列表
	mmInfo.ShortPskList = pskList
	// 随机生成ClientEcdhKeys,握手用
	mmInfo.ClientEcdhKeys = CreateClientEcdhKeys()

	// 开始握手
	err := MMHandShakeByLongLink(mmInfo)
	if err != nil {
		return nil, err
	}

	// 发送一次心跳包
	err = SendHeartBeat(mmInfo)
	if err != nil {
		return nil, err
	}

	// 握手成功，设置短链接的HOST 和 新的URL
	shortURL := []byte("/mmtls/")
	mmInfo.ShortURL = string(append(shortURL, []byte(baseutils.RandomSmallHexString(8))[0:]...))

	return mmInfo, nil
}

// MMHandShakeByLongLink 长链接握手(与短链接握手类似，只需选择其中一种握手方式，微信是采用长链接握手，每次登陆前都需要握手)
func MMHandShakeByLongLink(mmInfo *MMInfo) error {
	// 发送握手请求 - ClientHello
	clientHelloData := CreateHandShakeClientHelloData(mmInfo)
	sendData := CreateRecordData(ServerHandShakeType, clientHelloData)

	// 发送ClientHello
	err := MMTCPSendData(mmInfo, sendData)
	if err != nil {
		return err
	}

	// 接收响应
	retItems, err := MMTCPRecvItems(mmInfo)
	if err != nil {
		log.Info(err.Error())
		return err
	}

	// 处理握手信息
	clientFinishedData, err := DealHandShakePackItems(mmInfo, retItems, clientHelloData)
	if err != nil {
		return err
	}

	// 发送clientFinished, 长链接必须要发送，但服务器不会响应
	err = MMTCPSendData(mmInfo, clientFinishedData)
	if err != nil {
		return err
	}

	return nil
}

// SendHeartBeat 发送心跳包
func SendHeartBeat(mmInfo *MMInfo) error {
	// 发送心跳包
	heartData, err := GetHeartBeatData(mmInfo)
	err = MMTCPSendData(mmInfo, heartData)
	if err != nil {
		return err
	}

	// 接收心跳包响应
	retItem, err := MMTCPRecvOneItem(mmInfo)
	if err != nil {
		return err
	}

	// 解析心跳包响应数据，但不需处理
	_, err = MMLongUnPackData(mmInfo, retItem)
	if err != nil {
		return err
	}

	return nil
}

// GetHeartBeatData 获取心跳包数据
func GetHeartBeatData(mmInfo *MMInfo) ([]byte, error) {
	retBytes := make([]byte, 0)

	bodyData := make([]byte, 0)
	// LongPackTagInfo
	packTagInfo := &LongPackHeaderInfo{}
	packTagInfo.HeaderLen = 16
	packTagInfo.Version = MMLongVersion
	packTagInfo.Operation = MMLongOperationSmartHeartBeat
	packTagInfo.SequenceNumber = 0xffffffff
	tagInfoBytes := LongPackHeaderInfoSerialize(packTagInfo)

	// BodyData
	dataLength := uint32(len(tagInfoBytes) + 4)
	bodyData = append(bodyData, baseutils.Int32ToBytes(dataLength)[0:]...)
	bodyData = append(bodyData, tagInfoBytes[0:]...)

	// RecordHead
	tmpLen := uint16(dataLength + 16)
	recordHeader := GetRecordDataByLength(BodyType, tmpLen)

	clientSeq := atomic.AddUint32(&mmInfo.LONGClientSeq, 1) - 1
	// 加密
	tmpNonce := GetNonce(mmInfo.LongHdkfKey.EncodeNonce, uint32(clientSeq))
	tmpAad := []byte{0x00, 0x00, 0x00, 0x00}
	tmpAad = append(tmpAad, baseutils.Int32ToBytes(uint32(clientSeq))...)
	tmpAad = append(tmpAad, recordHeader...)
	encodeData, err := AesGcmEncrypt(mmInfo.LongHdkfKey.EncodeAesKey, tmpNonce, tmpAad, bodyData)
	if err != nil {
		return retBytes, err
	}

	// 组包
	retBytes = append(retBytes, recordHeader...)
	retBytes = append(retBytes, encodeData...)

	return retBytes, nil
}

// MMLongPackData 长链接方式 打包请求数据
func MMLongPackData(mmInfo *MMInfo, seqId uint32, opCode uint32, data []byte) ([]byte, error) {
	retBytes := make([]byte, 0)

	bodyData := make([]byte, 0)
	// LongPackTagInfo
	packTagInfo := &LongPackHeaderInfo{}
	packTagInfo.HeaderLen = 16
	packTagInfo.Version = MMLongVersion
	packTagInfo.Operation = opCode
	packTagInfo.SequenceNumber = seqId
	tagInfoBytes := LongPackHeaderInfoSerialize(packTagInfo)

	// BodyData
	dataLength := uint32(len(data) + len(tagInfoBytes) + 4)
	bodyData = append(bodyData, baseutils.Int32ToBytes(dataLength)[0:]...)
	bodyData = append(bodyData, tagInfoBytes[0:]...)
	bodyData = append(bodyData, data[0:]...)

	// RecordHead
	tmpLen := uint16(dataLength + 16)
	recordHeader := GetRecordDataByLength(BodyType, tmpLen)

	// 加密
	tmpNonce := GetNonce(mmInfo.LongHdkfKey.EncodeNonce, atomic.LoadUint32(&mmInfo.LONGClientSeq))
	tmpAad := []byte{0x00, 0x00, 0x00, 0x00}
	tmpAad = append(tmpAad, baseutils.Int32ToBytes(atomic.LoadUint32(&mmInfo.LONGClientSeq))...)
	tmpAad = append(tmpAad, recordHeader...)
	encodeData, err := AesGcmEncrypt(mmInfo.LongHdkfKey.EncodeAesKey, tmpNonce, tmpAad, bodyData)
	if err != nil {
		return retBytes, err
	}

	// 组包，ClientSeq 索引值+1
	atomic.AddUint32(&mmInfo.LONGClientSeq, 1)
	//mmInfo.LONGClientSeq++
	retBytes = append(retBytes, recordHeader...)
	retBytes = append(retBytes, encodeData...)
	return retBytes, nil
}

// MMLongUnPackData 长链接方式 解包数据
func MMLongUnPackData(mmInfo *MMInfo, packItem *PackItem) (*LongRecvInfo, error) {
	// 解密
	tmpNonce := GetNonce(mmInfo.LongHdkfKey.DecodeNonce, atomic.LoadUint32(&mmInfo.LONGServerSeq))
	tmpAad := []byte{0x00, 0x00, 0x00, 0x00}
	tmpAad = append(tmpAad, baseutils.Int32ToBytes(atomic.LoadUint32(&mmInfo.LONGServerSeq))...)
	tmpAad = append(tmpAad, packItem.RecordHead...)
	DecodeData, err := AesGcmDecrypt(mmInfo.LongHdkfKey.DecodeAesKey, tmpNonce, tmpAad, packItem.PackData)
	if err != nil {
		return nil, err
	}

	// TotalLength
	current := 0
	totalLength := baseutils.BytesToInt32(DecodeData[current : current+4])
	current = current + 4
	if totalLength != uint32(len(DecodeData)) {
		return nil, errors.New("MMLongUnPackData err: totalLength != uint32(len(DecodeData))")
	}

	headerInfo, err := LongPackHeaderInfoDeSerialize(DecodeData[current : current+12])
	current = current + 12
	if err != nil {
		return nil, err
	}

	retLongRecvInfo := &LongRecvInfo{}
	retLongRecvInfo.HeaderInfo = headerInfo
	// 如果大于 说明是业务包，如果等于则是长链接心跳包
	if totalLength > uint32(headerInfo.HeaderLen) {
		retLongRecvInfo.RespData = DecodeData[current:]
	}

	// ServerSeq 索引值+1
	//mmInfo.LONGServerSeq = mmInfo.LONGServerSeq + 1
	atomic.AddUint32(&mmInfo.LONGServerSeq, 1)
	return retLongRecvInfo, nil
}
