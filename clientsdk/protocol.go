package clientsdk

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	clientsdk "feiyu.com/wx/clientsdk/hybrid"
	"github.com/lunny/log"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/ccdata"
	"feiyu.com/wx/clientsdk/proxynet"
	"feiyu.com/wx/protobuf/wechat"
	"golang.org/x/net/proxy"

	"github.com/golang/protobuf/proto"
)

// NewUserInfo 新建一个UserInfo
func NewUserInfo(uuid string, deviceId string, proxyInfo *proxynet.WXProxyInfo) *baseinfo.UserInfo {
	//随机生成 deviceID
	return &baseinfo.UserInfo{
		UUID:       uuid,
		Uin:        0,
		WxId:       "",
		Session:    []byte{},
		SessionKey: baseutils.RandomStringByLength(16),
		ShortHost:  "szshort.weixin.qq.com",
		LongHost:   "szlong.weixin.qq.com",
		//ShortHost:      "short.weixin.qq.com",
		//LongHost:       "long.weixin.qq.com",
		SyncKey:        []byte{},
		BalanceVersion: 1589560770,
		DeviceInfo:     createDeviceInfo(deviceId),
		DeviceInfoA16:  createDeviceInfoA16(),
		WifiInfo:       createWifiInfo(),
		ProxyInfo:      proxyInfo,
		// 默认登录版本号 135
		LoginRsaVer: baseinfo.DefaultLoginRsaVer,
	}
}

// 新建Wifi信息
func createWifiInfo() *baseinfo.WifiInfo {
	retWifiInfo := &baseinfo.WifiInfo{}
	retWifiInfo.Name = string(baseutils.RandomString(8, 15))
	retWifiInfo.WifiBssID = baseutils.RandomBSSID()
	return retWifiInfo
}

// 生成A16设备信息
func createDeviceInfoA16() *baseinfo.AndroidDeviceInfo {
	deviceInfo := &baseinfo.AndroidDeviceInfo{}
	deviceInfo.BuildBoard = "bullhead"
	return deviceInfo
}

// CreateDeviceInfo 生成新的设备信息 ipad
func createDeviceInfo(deviceId string) *baseinfo.DeviceInfo {
	deviceInfo := &baseinfo.DeviceInfo{}
	if deviceId == "" && len(deviceId) < 2 {
		deviceInfo.Imei = baseutils.RandomSmallHexString(32)
		tmpDeviceID := baseutils.HexStringToBytes(deviceInfo.Imei)
		tmpDeviceID[0] = 0x49
		deviceInfo.DeviceID = tmpDeviceID
	} else {
		tmpImei := deviceId[2:]
		deviceInfo.Imei = baseutils.RandomSmallHexString(2) + tmpImei
		tmpDeviceID := baseutils.HexStringToBytes(deviceId)
		deviceInfo.DeviceID = tmpDeviceID
	}
	deviceInfo.DeviceName = "iPad" //iPhone
	deviceInfo.TimeZone = "8.00"
	deviceInfo.Language = "zh_CN" //
	deviceInfo.DeviceBrand = "Apple"
	deviceInfo.RealCountry = "CN"
	deviceInfo.IphoneVer = "iPad4,7" //iPhone4,7
	deviceInfo.BundleID = "com.tencent.xin"
	deviceInfo.OsTypeNumber = "13.5"                     //12.4.6
	deviceInfo.OsType = "iPad" + deviceInfo.OsTypeNumber //iPhone
	deviceInfo.CoreCount = 4                             // 4核
	deviceInfo.AdSource = baseutils.RandomUUID()
	deviceInfo.UUIDOne = baseutils.RandomUUID()
	deviceInfo.UUIDTwo = baseutils.RandomUUID()
	// 运营商名
	deviceInfo.CarrierName = "(null)"
	deviceInfo.SoftTypeXML = CreateSoftInfoXML(deviceInfo)
	// ClientCheckDataXML
	deviceInfo.ClientCheckDataXML = ccdata.CreateClientCheckDataXML(deviceInfo)
	return deviceInfo
}

// CreateSoftInfoXML CreateSoftInfoXML
func CreateSoftInfoXML(deviceInfo *baseinfo.DeviceInfo) string {
	// 生成DeviceInfoXML
	var retString string
	retString = retString + "<softtype>"
	retString = retString + "<k3>" + deviceInfo.OsTypeNumber + "</k3>"
	retString = retString + "<k9>" + deviceInfo.DeviceName + "</k9>"
	retString = retString + "<k10>" + strconv.Itoa(int(deviceInfo.CoreCount)) + "</k10>"
	retString = retString + "<k19>" + deviceInfo.UUIDOne + "</k19>"
	retString = retString + "<k20>" + deviceInfo.UUIDTwo + "</k20>"
	retString = retString + "<k22>" + deviceInfo.CarrierName + "</k22>"
	retString = retString + "<k24>" + baseutils.BuildRandomMac() + "/k24"
	retString = retString + "<k33>微信</k33>"
	// <k47>: 网络类型 1-wifi
	retString = retString + "<k47>1</k47>"
	// <k50>: 是否越狱 0-非越狱 1-越狱
	retString = retString + "<k50>0</k50>"
	retString = retString + "<k51>" + deviceInfo.BundleID + "</k51>"
	retString = retString + "<k54>" + deviceInfo.IphoneVer + "</k54>"
	// <k61>: 设备UUID是新的设备，还是老的设备
	retString = retString + "<k61>" + strconv.Itoa(1) + "</k61>"
	retString = retString + "</softtype>"

	return retString
}

// GetEncryptUserInfo 获取加密信息
func GetEncryptUserInfo(userInfo *baseinfo.UserInfo) string {
	if userInfo.WifiInfo == nil {
		userInfo.WifiInfo = createWifiInfo()
	}
	tmpString := "wifissid=" + userInfo.WifiInfo.Name
	tmpString = tmpString + "&wifibssid=" + userInfo.WifiInfo.WifiBssID
	timeStamp := strconv.Itoa(int(time.Now().UnixNano() / 1000))
	tmpString = tmpString + "&ssid_timestamp=" + timeStamp
	srcBytes := []byte(tmpString)
	srcBytes = append(srcBytes, 0)
	userInfo.GenHBKey()
	encData := baseutils.AesEncryptECB(srcBytes, userInfo.HBAesKey)
	return base64.StdEncoding.EncodeToString(encData)
}

// GetDialer 获取代理
func GetDialer(userInfo *baseinfo.UserInfo) proxy.Dialer {
	if userInfo.Dialer != nil {
		return userInfo.Dialer
	}
	if userInfo.ProxyInfo == nil {
		return nil
	}
	return userInfo.ProxyInfo.GetDialer()
}

// Paser62Data 解析62数据
func Parse62Data(data62 string) (string, error) {
	retList := strings.Split(strings.ToUpper(data62), "6E756C6C5F1020")
	if len(retList) < 2 {
		return "", errors.New("InitLoginDataInfo err: loginDataInfo.Data Split 6E756C6C5F1020 error")
	}

	// 截取64位数据
	data64Str := retList[1][0:64]
	// 设置Imei
	return string(baseutils.HexStringToBytes(data64Str)), nil
}

// CreatePackHead 创建包头
func CreatePackHead(userInfo *baseinfo.UserInfo, compressType byte, urlID uint32, srcData []byte, encodeData []byte, zipLen uint32, encodeType byte, encodeVersion uint32) *baseinfo.PackHeader {
	retHeader := &baseinfo.PackHeader{}

	// Signature
	retHeader.Signature = 0xbf
	retHeader.CompressType = compressType
	retHeader.EncodeType = encodeType << 4
	retHeader.ServerVersion = baseinfo.ClientVersion
	retHeader.Uin = userInfo.Uin
	retHeader.Session = userInfo.Session
	retHeader.URLID = urlID
	retHeader.SrcLen = uint32(len(srcData))
	retHeader.ZipLen = zipLen
	retHeader.EncodeVersion = encodeVersion
	retHeader.HeadDeviceType = baseinfo.MMHeadDeviceTypeIpadUniversal
	retHeader.CheckSum = 0x00
	// 如果有压缩，则计算Sum值
	if retHeader.CompressType == baseinfo.MMPackDataTypeCompressed {
		retHeader.CheckSum = CalcHeadCheckSum(userInfo.Uin, userInfo.CheckSumKey, srcData)
	}

	// test
	/*if retHeader.URLID == 213 {
		retHeader.HeadDeviceType = 0x00
	}*/

	retHeader.RunState = baseinfo.MMAppRunStateNormal
	retHeader.RqtCode = baseutils.CalcMsgCrcForString_807(baseutils.Md5ValueByte(encodeData, false))
	retHeader.EndFlag = 0x00
	retHeader.Data = encodeData

	return retHeader
}

// PackHeaderSerialize 序列化PackHeader
func PackHeaderSerialize(packHeader *baseinfo.PackHeader, needCookie bool) []byte {
	retBytes := make([]byte, 0)
	retBytes = append(retBytes, packHeader.Signature)
	retBytes = append(retBytes, 0)
	encodeType := packHeader.EncodeType
	if needCookie {
		packHeader.EncodeType = packHeader.EncodeType + 0xf
	}
	retBytes = append(retBytes, packHeader.EncodeType)
	retBytes = append(retBytes, baseutils.Int32ToBytes(packHeader.ServerVersion)[0:]...)
	retBytes = append(retBytes, baseutils.Int32ToBytes(packHeader.Uin)[0:]...)
	if needCookie {
		retBytes = append(retBytes, packHeader.Session[0:]...)
	}
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.URLID)[0:]...)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.SrcLen)[0:]...)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.ZipLen)[0:]...)
	// hybrid
	if encodeType>>4 == 12 {
		retBytes = append(retBytes, []byte{byte(packHeader.HybridKeyVer)}...)
	}
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.EncodeVersion)[0:]...)
	retBytes = append(retBytes, packHeader.HeadDeviceType)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.CheckSum)[0:]...)
	retBytes = append(retBytes, packHeader.RunState)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.RqtCode)[0:]...)
	retBytes = append(retBytes, packHeader.EndFlag)
	headLen := byte(len(retBytes))
	retBytes[1] = packHeader.CompressType + headLen<<2
	//log.Println(hex.EncodeToString(retBytes))
	retBytes = append(retBytes, packHeader.Data[0:]...)
	return retBytes
}

// Pack 打包加密数据
func Pack(userInfo *baseinfo.UserInfo, src []byte, urlID uint32, encodeType byte) []byte {
	retData := make([]byte, 0)
	if encodeType == 7 || encodeType == 17 {
		//加密类型7
		encodeData := src
		if encodeType == 7 {
			encodeData = baseutils.NoCompressRsaByVer(src, userInfo.GetLoginRsaVer())
		}
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, src, encodeData, uint32(len(src)), 7, userInfo.GetLoginRsaVer())
		retData = PackHeaderSerialize(packHeader, false)
	} else if encodeType == 5 {
		// 加密类型5
		zipBytes := baseutils.CompressByteArray(src)
		encodeData := baseutils.AesEncrypt(zipBytes, userInfo.SessionKey)
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeCompressed, urlID, src, encodeData, uint32(len(zipBytes)), encodeType, 0)
		retData = PackHeaderSerialize(packHeader, true)
	} else if encodeType == 9 {
		// 加密类型9
		encodeData := src
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, src, encodeData, uint32(len(src)), encodeType, userInfo.GetLoginRsaVer())
		retData = PackHeaderSerialize(packHeader, true)
	} else if encodeType == 1 {
		// 加密类型1
		encodeData := baseutils.NoCompressRsaByVer(src, userInfo.GetLoginRsaVer())
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, src, encodeData, uint32(len(src)), encodeType, userInfo.GetLoginRsaVer())
		retData = PackHeaderSerialize(packHeader, true)
	} else if encodeType == 12 {
		secKeyMgr := NewSecLoginKeyMgrByVer(146)
		reqData := src
		//加密
		encrypt, epKey, token, ecdhpairkey, err := clientsdk.HybridEncrypt(reqData, secKeyMgr.WeChatPubKey)
		if err != nil {
			log.Error("加密 error", err.Error())
		}
		ecdhPacket := &wechat.EcdhPacket{
			Type: proto.Uint32(1),
			Key: &wechat.BufferT{
				ILen:   proto.Uint32(415),
				Buffer: ecdhpairkey.PubKey,
			},
			Token:        token,
			Url:          proto.String(""),
			ProtobufData: encrypt,
		}
		secKeyMgr.PubKey = ecdhpairkey.PubKey
		secKeyMgr.PriKey = ecdhpairkey.PriKey
		secKeyMgr.SourceData = reqData
		secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, epKey[24:]...)
		secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, reqData...)
		ecdhDataPacket, err := proto.Marshal(ecdhPacket)
		if err != nil {
			log.Error("ecdhDataPacket error", err.Error())
		}
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, ecdhDataPacket, ecdhDataPacket, uint32(len(ecdhDataPacket)), encodeType, uint32(0x4e))
		//设置Hybrid 加密密钥版本
		packHeader.HybridKeyVer = secKeyMgr.WeChatPubKeyVersion
		//开始组头
		retData := PackHeaderSerialize(packHeader, false)
		//log.Println(hex.EncodeToString(retData))
		/*resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secautoauth", retData)
		if err != nil {
			log.Error("mmtls error", err.Error())
		}*/
		/*packHeader, err = DecodePackHeader(resp, nil)
		if err != nil {
			log.Error("ecdhDataPacket error", err.Error())
		}
		packHeader.Data, err = clientsdk.HybridEcdhDecrypt(packHeader.Data, secKeyMgr.PriKey, secKeyMgr.PubKey, secKeyMgr.FinalSha256)
		if err != nil {
			log.Error("HybridEcdhDecrypt error", err.Error())
		}
		return packHeader, err*/
		return retData
	}
	return retData
}

// Pack 打包加密数据
func Pack12(userInfo *baseinfo.UserInfo, src []byte, urlID uint32, encodeType byte) ([]byte, *SecLoginKeyMgr) {
	retData := make([]byte, 0)
	secKeyMgr := NewSecLoginKeyMgrByVer(146)
	reqData := src
	//加密
	encrypt, epKey, token, ecdhpairkey, err := clientsdk.HybridEncrypt(reqData, secKeyMgr.WeChatPubKey)
	if err != nil {
		log.Error("加密 error", err.Error())
	}
	ecdhPacket := &wechat.EcdhPacket{
		Type: proto.Uint32(1),
		Key: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: ecdhpairkey.PubKey,
		},
		Token:        token,
		Url:          proto.String(""),
		ProtobufData: encrypt,
	}
	secKeyMgr.PubKey = ecdhpairkey.PubKey
	secKeyMgr.PriKey = ecdhpairkey.PriKey
	secKeyMgr.SourceData = reqData
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, epKey[24:]...)
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, reqData...)
	ecdhDataPacket, err := proto.Marshal(ecdhPacket)
	if err != nil {
		log.Error("ecdhDataPacket error", err.Error())
	}
	packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, ecdhDataPacket, ecdhDataPacket, uint32(len(ecdhDataPacket)), encodeType, uint32(0x4e))
	//设置Hybrid 加密密钥版本
	packHeader.HybridKeyVer = secKeyMgr.WeChatPubKeyVersion
	//开始组头
	retData = PackHeaderSerialize(packHeader, false)
	//log.Println(hex.EncodeToString(retData))
	/*resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secautoauth", retData)
	if err != nil {
		log.Error("mmtls error", err.Error())
	}*/
	/*packHeader, err = DecodePackHeader(resp, nil)
	if err != nil {
		log.Error("ecdhDataPacket error", err.Error())
	}
	packHeader.Data, err = clientsdk.HybridEcdhDecrypt(packHeader.Data, secKeyMgr.PriKey, secKeyMgr.PubKey, secKeyMgr.FinalSha256)
	if err != nil {
		log.Error("HybridEcdhDecrypt error", err.Error())
	}
	return packHeader, err*/
	return retData, secKeyMgr
}

// DecodePackHeader DecodePackHeader
func DecodePackHeader(respData []byte, reqData []byte) (*baseinfo.PackHeader, error) {
	packHeader := &baseinfo.PackHeader{}
	packHeader.ReqData = reqData
	packHeader.RetCode = 0
	// 如果数据长度小于等于32, 则表明请求出错
	if len(respData) <= 32 {
		packHeader.RetCode = GetRespErrorCode(respData)
		return packHeader, errors.New("DecodePackHeader err: len(respData) <= 32")
	}

	current := 0
	packHeader.Signature = respData[current]
	current++
	packHeader.HeadLength = (respData[current]) >> 2
	packHeader.CompressType = (respData[current]) & 3
	current++
	packHeader.EncodeType = respData[current] >> 4
	sessionLen := int(respData[current] & 0x0f)
	current++
	packHeader.ServerVersion = baseutils.BytesToInt32(respData[current : current+4])
	current = current + 4
	packHeader.Uin = baseutils.BytesToInt32(respData[current : current+4])
	current = current + 4
	if sessionLen > 0 {
		packHeader.Session = respData[current : current+sessionLen]
		current = current + sessionLen
	}
	retLen := uint32(0)
	packHeader.URLID, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.SrcLen, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.ZipLen, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.EncodeVersion, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.HeadDeviceType = respData[current]
	current = current + 1
	packHeader.CheckSum, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.RunState = respData[current]
	current = current + 1
	packHeader.RqtCode, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.EndFlag = respData[current]
	current = current + 1
	// // 后面还有一个字节-- 可能是7.10新版本增加的一个字节，待后面分析
	// current = current + 1
	// if current != int(packHeader.HeadLength) {
	// 	return nil, errors.New("DecodePackHeader failed current != int(packHeader.HeadLength")
	// }
	packHeader.Data = respData[packHeader.HeadLength:]
	return packHeader, nil
}

// ParseResponseData 解析相应数据
func ParseResponseData(userInfo *baseinfo.UserInfo, packHeader *baseinfo.PackHeader, response proto.Message) error {
	//  判断包体长度是否大于0
	if len(packHeader.Data) <= 0 {
		log.Error("ParseResponseData err: len(packHeader.Data) <= 0")
		return errors.New("ParseResponseData err: len(packHeader.Data) <= 0")
	}
	var decptBody []byte
	var err error
	if packHeader.EncodeType == 12 {
		decptBody = packHeader.Data
	} else if packHeader.EncodeType == 5 {
		// 解密
		decptBody, err = baseutils.AesDecrypt(packHeader.Data, userInfo.SessionKey)
		if err != nil {
			return err
		}
		// 判断是否有压缩
		if packHeader.CompressType == baseinfo.MMPackDataTypeCompressed {
			if decptBody != nil {
				//log.Println(hex.EncodeToString(decptBody))
				decptBody, err = baseutils.UnzipByteArray(decptBody)
				if err != nil {
					log.Error("ParseResponseData err:", err.Error(), packHeader.URLID)
					return err
				}
			} else {
				return errors.New("decptBody err: len(decptBody) == nil")
			}
		}
	} else {
		// 解密
		decptBody, err = baseutils.AesDecrypt(packHeader.Data, userInfo.SessionKey)
		//log.Println(hex.EncodeToString(decptBody))
		if err != nil {
			//log.Error("ParseResponseData err:", err.Error())
			return err
		}
		// 判断是否有压缩
		if packHeader.CompressType == baseinfo.MMPackDataTypeCompressed {
			if decptBody != nil {
				decptBody, err = baseutils.UnzipByteArray(decptBody)
				if err != nil {
					//log.Println(hex.EncodeToString(packHeader.Data))
					log.Error("ParseResponseData err:", err.Error(), packHeader.URLID)
					return err
				}
			} else {
				return errors.New("decptBody err: len(decptBody) == nil")
			}
		}
	}

	// 更新UserInfo
	if len(packHeader.Session) > 6 {
		userInfo.Session = packHeader.Session
	}

	if packHeader.Uin != 0 {
		userInfo.Uin = packHeader.Uin
	}
	//log.Println(hex.EncodeToString(decptBody))
	// 解包ProtoBuf
	err = proto.Unmarshal(decptBody, response)
	if err != nil {
		log.Error("ParseResponseData err:", err.Error())
		return err
	}

	return nil
}

// GetBaseRequest 获取baserequest
func GetBaseRequest(userInfo *baseinfo.UserInfo) *wechat.BaseRequest {
	ret := &wechat.BaseRequest{}
	ret.SessionKey = []byte(userInfo.SessionKey)
	ret.Uin = &userInfo.Uin
	if !strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") && userInfo.DeviceInfo != nil {
		ret.DeviceId = userInfo.DeviceInfo.DeviceID
		ret.ClientVersion = &baseinfo.ClientVersion
		ret.OsType = &userInfo.DeviceInfo.OsType
		ret.Scene = proto.Uint32(0)
		//log.Info("ios is base request")
	} else {
		ret.ClientVersion = &baseinfo.AndroidClientVersion
		ret.OsType = &baseinfo.AndroidDeviceType
		ret.DeviceId = userInfo.DeviceInfoA16.DeviceId
		ret.Scene = proto.Uint32(1)
		//log.Info("android is base request")
	}
	return ret
}

// CalcCheckCode 计算校验码
// wxId : 微信ID
// currentTime : 时间
// return uint32 计算出来的校验码
func CalcCheckCode(wxID string, currentTime time.Time) uint32 {
	var tmpTime = currentTime.Format("051504010206")
	if len(wxID) <= 0 {
		tmpTime = tmpTime + "fffffff"
	}
	if len(wxID) > 0 {
		wxIDMd5 := baseutils.Md5Value(wxID)
		tmpTime = tmpTime + baseutils.StringCut(wxIDMd5, 0, 7)
	}

	misSecond := currentTime.UnixNano() / 1000000
	var modValue = uint32(misSecond % 65535)
	var modValue2 = int(misSecond%7 + 100)
	modHexString := string(baseutils.UInt32To16Bytes(modValue))
	mod2String := strconv.Itoa(modValue2)
	tmpTime = tmpTime + modHexString + mod2String

	return baseutils.HashCode(tmpTime)
}

// WithSeqidCalcCheckCode 根据seqid计算
func WithSeqidCalcCheckCode(wxID string, seqid int64) uint32 {
	currentTime := time.Unix(0, seqid)
	var tmpTime = currentTime.Format("051504010206")
	if len(wxID) <= 0 {
		tmpTime = tmpTime + "fffffff"
	}
	if len(wxID) > 0 {
		wxIDMd5 := baseutils.Md5Value(wxID)
		tmpTime = tmpTime + baseutils.StringCut(wxIDMd5, 0, 7)
	}

	misSecond := currentTime.UnixNano() / 1000000
	var modValue = uint32(misSecond % 65535)
	var modValue2 = int(misSecond%7 + 100)
	modHexString := string(baseutils.UInt32To16Bytes(modValue))
	mod2String := strconv.Itoa(modValue2)
	tmpTime = tmpTime + modHexString + mod2String

	return baseutils.HashCode(tmpTime)
}

// GetErrStrByErrCode GetErrStrByErrCode
func GetErrStrByErrCode(errCode int32) string {
	if errCode == -2 {
		return "error: args not whole or args type error"
	}
	if errCode == -13 {
		return "error: session timeout"
	}
	if errCode == -102 {
		return "error: cert expired"
	}
	if errCode == -306 {
		return "error: ecdh failed rollback"
	}
	if errCode == -3001 || errCode == -3003 {
		return "error: need get dns"
	}
	if errCode == -3002 {
		return "error: MM_ERR_IDCDISASTER"
	}

	return "unknow"
}

// GetRespErrorCode 当response 小于32个字节时，调用这个接口获取响应的错误码
func GetRespErrorCode(data []byte) int32 {
	tmpData := make([]byte, 0)
	tmpData = append(tmpData, data[2:6]...)

	tmpRet := binary.BigEndian.Uint32(tmpData)
	return int32(tmpRet)
}

// CalcHeadCheckSum 计算HeadCheckSum值
func CalcHeadCheckSum(uin uint32, checkSumKey []byte, srcData []byte) uint32 {
	uinBytes := baseutils.Int32ToBytes(uin)
	tmpBytes := append(uinBytes, checkSumKey[0:]...)
	md5Value := baseutils.Md5Value16(tmpBytes)

	dataLen := uint32(len(srcData))
	dataLenBytes := baseutils.Int32ToBytes(dataLen)
	tmpBytes = append(dataLenBytes, checkSumKey[0:]...)
	tmpBytes = append(tmpBytes, md5Value[0:]...)
	md5Value = baseutils.Md5Value16(tmpBytes)

	tmpBytes = append([]byte{}, md5Value[0:]...)
	// 计算返回
	/*tmpSum := baseutils.Adler32(1, tmpBytes)
	return baseutils.Adler32(tmpSum, srcData)*/
	return crc32.ChecksumIEEE(tmpBytes)
}

// TenPaySignDes3 支付相关的加密算法
func TenPaySignDes3(srcData string, encKey string) (string, error) {
	srcMD5Bytes := []byte(baseutils.Md5ValueByte([]byte(srcData), true))
	keyMD5Bytes := baseutils.Md5ValueByte([]byte(encKey), true)
	desKey := baseutils.HexStringToBytes(keyMD5Bytes)

	encBytes := make([]byte, 0)
	for index := 0; index < 4; index++ {
		currentOffset := index * 8
		tmpSrcData := srcMD5Bytes[currentOffset : currentOffset+8]
		encData, err := baseutils.Encrypt3DES(tmpSrcData, desKey)
		if err != nil {
			return "", err
		}
		encBytes = append(encBytes, encData...)
	}
	return baseutils.BytesToHexString(encBytes, true), nil
}
