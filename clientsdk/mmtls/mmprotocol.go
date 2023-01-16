package mmtls

import (
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"feiyu.com/wx/clientsdk/baseutils"
	"github.com/lunny/log"
	"github.com/wsddn/go-ecdh"
	"sync/atomic"
)

// CreateNewMMInfo 创建新的MMInfo
func CreateNewMMInfo() *MMInfo {
	mmInfo := &MMInfo{}
	mmInfo.ShortHost = "szshort.weixin.qq.com"
	mmInfo.LongHost = "szlong.weixin.qq.com"
	//mmInfo.ShortHost = "short.weixin.qq.com"
	//mmInfo.LongHost = "long.weixin.qq.com"
	mmInfo.LONGPort = "443"
	//if mmInfo.LONGPort == "443" {
	//	longPort := []string{"80","8080","443"}
	//	math_rand.Seed(time.Now().UnixNano())
	//	mmInfo.LONGPort = longPort[math_rand.Intn(3)]
	//}
	mmInfo.LONGClientSeq = 1
	mmInfo.LONGServerSeq = 1

	return mmInfo
}

// InitMMTLSInfoShort 如果使用MMTLS，每次登陆前都需要初始化MMTLSInfo信息
// pskList: 之前握手服务端返回的，要保存起来，后面握手时使用， 第一次握手传空数组即可
func InitMMTLSInfoShort(hostName string, pskList []*Psk) *MMInfo {
	// 初始化MMInfo
	mmInfo := CreateNewMMInfo()
	mmInfo.ShortPskList = pskList
	// 随机生成ClientEcdhKeys
	mmInfo.ClientEcdhKeys = CreateClientEcdhKeys()
	mmInfo.ShortHost = hostName
	// 握手
	mmInfo, err := MMHandShakeByShortLink(mmInfo, hostName)
	// 如果握手失败 就不使用MMTLS
	if err != nil {
		return nil
	}

	// 握手成功，设置好HOST 和 新的URL
	shortURL := []byte("/mmtls/")
	mmInfo.ShortHost = hostName
	mmInfo.ShortURL = string(append(shortURL, []byte(baseutils.RandomSmallHexString(8))[0:]...))

	return mmInfo
}

// CreateClientEcdhKeys 创建新的ClientEcdhKeys
func CreateClientEcdhKeys() *ClientEcdhKeys {
	// 随机
	clientEcdhKeys := &ClientEcdhKeys{}
	e := ecdh.NewEllipticECDH(elliptic.P256())
	priKey1, pubKey1, _ := e.GenerateKey(rand.Reader)
	priKey2, pubKey2, _ := e.GenerateKey(rand.Reader)
	clientEcdhKeys.PriKey1 = priKey1
	clientEcdhKeys.PriKey2 = priKey2
	clientEcdhKeys.PubKeyBuf1 = e.Marshal(pubKey1)
	clientEcdhKeys.PubKeyBuf2 = e.Marshal(pubKey2)

	return clientEcdhKeys
}

// MMHandShakeByShortLink 通过短链接握手
func MMHandShakeByShortLink(mmInfo *MMInfo, hostName string) (*MMInfo, error) {
	shortURL := []byte("/mmtls/")
	mmURL := append(shortURL, []byte(baseutils.RandomSmallHexString(8))[0:]...)
	mmInfo.ShortURL = string(mmURL)

	// 发送握手请求 - ClientHello
	clientHelloData := CreateHandShakeClientHelloData(mmInfo)
	sendData := CreateRecordData(ServerHandShakeType, clientHelloData)
	retBytes, err := MMHTTPPost(mmInfo, sendData)
	if err != nil {
		log.Error("短连接握手时报错:", err.Error())
		return nil, err
	}

	// 解析握手相应数据
	retItems, err := ParserMMtlsResponseData(retBytes)
	if err != nil {
		return nil, err
	}

	// 处理握手信息
	clientFinishData, err := DealHandShakePackItems(mmInfo, retItems, clientHelloData)
	_ = clientFinishData

	return mmInfo, nil
}

// ParserMMtlsResponseData 解析mmtls响应数据
func ParserMMtlsResponseData(data []byte) ([]*PackItem, error) {
	// RecodeHead *RecodeHead
	retItems := make([]*PackItem, 0)

	// 总数据大小
	totalLength := uint32(len(data))

	current := uint32(0)
	// 解析所有包
	for current < totalLength {
		packItem := &PackItem{}

		// recordHead
		if current+5 > totalLength {
			return retItems, errors.New("ParserMMtlsResponseData err: current+5 >= totalLength")
		}
		recordHead := RecordHeadDeSerialize(data[current:])
		packItem.RecordHead = data[current : current+5]
		current = current + 5
		// PackData
		// 判断数据是否有问题
		if current+uint32(recordHead.Size) > totalLength {
			return retItems, errors.New("ParserMMtlsResponseData err: current+uint32(recordHead.Size) >= totalLength")
		}
		packItem.PackData = data[current : current+uint32(recordHead.Size)]

		// current
		current = current + uint32(recordHead.Size)
		retItems = append(retItems, packItem)
	}

	return retItems, nil
}

// DealHandShakePackItems 解密packItems
func DealHandShakePackItems(mmInfo *MMInfo, packItems []*PackItem, clientHelloReq []byte) ([]byte, error) {
	retClientFinishData := make([]byte, 0)

	// 先解析 ServerHello
	secretKey, err := DealServerHello(mmInfo, packItems[0])
	if err != nil {
		return retClientFinishData, err
	}

	// 计算HashRet
	hashData := make([]byte, 0)
	hashData = append(hashData, clientHelloReq[0:]...)
	hashData = append(hashData, packItems[0].PackData...)
	hashRet := Sha256(hashData)

	// 密钥扩展
	message := []byte("handshake key expansion")
	message = append(message, hashRet...)
	aesKeyExpand := HkdfExpand(secretKey, message, 56)
	gcmAesKey := aesKeyExpand[0x10:0x20]
	oriNonce := aesKeyExpand[0x2c:]

	// 解密后面的包
	count := len(packItems)
	for index := 1; index < count; index++ {
		tmpPackItem := packItems[index]

		// 解密数据
		tmpNonce := GetNonce(oriNonce, uint32(index))
		tmpAad := []byte{0x00, 0x00, 0x00, 0x00}
		tmpAad = append(tmpAad, baseutils.Int32ToBytes(atomic.LoadUint32(&mmInfo.LONGServerSeq))...)
		tmpAad = append(tmpAad, tmpPackItem.RecordHead...)
		decodeData, err := AesGcmDecrypt(gcmAesKey, tmpNonce, tmpAad, tmpPackItem.PackData)
		// 设置解密后的数据
		tmpPackItem.PackData = decodeData
		if err != nil {
			return retClientFinishData, err
		}
		atomic.AddUint32(&mmInfo.LONGServerSeq, 1)
		//mmInfo.LONGServerSeq++

		// 处理CertificateVerifyType
		tmpType := decodeData[4]
		if tmpType == CertificateVerifyType {
			// 校验服务器
			flag, err := DealCertificateVerify(clientHelloReq, packItems[0].PackData, decodeData)
			if err != nil {
				return retClientFinishData, err
			}
			if !flag {
				return retClientFinishData, errors.New("DealHandShakePackItems err: CertificateVerify failed")
			}
		}

		// 处理NewSessionTicketType
		if tmpType == NewSessionTicketType {
			err := DealNewSessionTicket(mmInfo, decodeData)
			if err != nil {
				return retClientFinishData, err
			}
		}

		// 处理Server FinishType
		if tmpType == FinishedType {
			// 第一步验证ServerFinished数据
			tmpHashData := make([]byte, 0)
			tmpHashData = append(tmpHashData, clientHelloReq[0:]...)
			tmpHashData = append(tmpHashData, packItems[0].PackData[0:]...)
			tmpHashData = append(tmpHashData, packItems[1].PackData[0:]...)
			tmpHashData = append(tmpHashData, packItems[2].PackData[0:]...)
			tmpHashValue := Sha256(tmpHashData)
			serverFinished, err := FinishedDeSerialize(tmpPackItem.PackData[4:])
			if err != nil {
				return retClientFinishData, err
			}
			bSuccess := VerifyFinishedData(secretKey, tmpHashValue, serverFinished.VerifyData)
			if !bSuccess {
				return retClientFinishData, errors.New("DealHandShakePackItems err: Finished verify failed")
			}

			// 第二步生成ClientFinished数据，然后加密
			hkdfClientFinish := HkdfExpand(secretKey, []byte("client finished"), 32)
			hmacRet := HmacHash256(hkdfClientFinish, tmpHashValue)
			aesGcmParam := &AesGcmParam{}
			aesGcmParam.AesKey = aesKeyExpand[0x00:0x10]
			aesGcmParam.Nonce = aesKeyExpand[0x20:0x2c]
			// 创建Finished
			finished := CreateFinished(hmacRet)
			// 加密
			finishedData := FinishedSerialize(finished)
			clientSeq := atomic.AddUint32(&mmInfo.LONGClientSeq, 1) - 1
			encodeData, err := EncryptedReqData(aesGcmParam, finishedData, ServerHandShakeType, clientSeq)
			if err != nil {
				return retClientFinishData, err
			}
			retClientFinishData = CreateRecordData(ServerHandShakeType, encodeData)
			//mmInfo.LONGClientSeq++
			//atomic.AddUint32(&mmInfo.LONGClientSeq, 1)
			break
		}
	}

	// 计算扩展出来的用于后续加密的Key
	tmpExpandHashData := make([]byte, 0)
	tmpExpandHashData = append(tmpExpandHashData, clientHelloReq[0:]...)
	tmpExpandHashData = append(tmpExpandHashData, packItems[0].PackData[0:]...)
	tmpExpandHashData = append(tmpExpandHashData, packItems[1].PackData[0:]...)
	tmpExpandHashValue := Sha256(tmpExpandHashData)

	// PskAccessKey 短连接MMTLS密钥
	expandPskAccessData := []byte("PSK_ACCESS")
	expandPskAccessData = append(expandPskAccessData, tmpExpandHashValue[0:]...)
	mmInfo.PskAccessKey = HkdfExpand(secretKey, expandPskAccessData, 32)

	// AppDataKeyExtension 长链接MMTLS密钥
	tmpExpandHashData = append(tmpExpandHashData, packItems[2].PackData[0:]...)
	tmpLongHashValue := Sha256(tmpExpandHashData)
	expandedSecret := append([]byte("expanded secret"), tmpLongHashValue[0:]...)
	retExpandSecret := HkdfExpand(secretKey, expandedSecret, 32)
	appDataKeyData := append([]byte("application data key expansion"), tmpLongHashValue[0:]...)
	appDataKeyExtension := HkdfExpand(retExpandSecret, appDataKeyData, 56)
	mmInfo.LongHdkfKey = &HkdfKey56{}
	mmInfo.LongHdkfKey.EncodeAesKey = appDataKeyExtension[0x00:0x10]
	mmInfo.LongHdkfKey.DecodeAesKey = appDataKeyExtension[0x10:0x20]
	mmInfo.LongHdkfKey.EncodeNonce = appDataKeyExtension[0x20:0x2c]
	mmInfo.LongHdkfKey.DecodeNonce = appDataKeyExtension[0x2c:]

	// 返回ClientFinishData
	return retClientFinishData, nil
}

// DealServerHello 处理ServerHello
func DealServerHello(mmInfo *MMInfo, packItem *PackItem) ([]byte, error) {
	// 解析ServerHello
	serverHello, err := ServerHelloDeSerialize(packItem.PackData[4:])
	if err != nil {
		return []byte{}, err
	}

	// 解析ServerKeyShare
	serverKeyShareExtension, err := ServerKeyShareExtensionDeSerialize(serverHello.ExtensionList[0].ExtensionData)
	if err != nil {
		return []byte{}, err
	}

	// 解析ServerPublicKey
	ecdhTool := ecdh.NewEllipticECDH(elliptic.P256())
	serverPubKey, isOk := ecdhTool.Unmarshal(serverKeyShareExtension.PublicValue)
	if !isOk {
		return []byte{}, errors.New("DecodePackItems ecdhTool.Unmarshal(serverKeyShareExtension.PublicValue) failed")
	}

	// 根据NameGroup 决定使用哪个Privakey
	ecdhPriKey := mmInfo.ClientEcdhKeys.PriKey1
	if serverKeyShareExtension.KeyOfferNameGroup == 2 {
		ecdhPriKey = mmInfo.ClientEcdhKeys.PriKey2
	}

	// 协商密钥
	secretKey, err := ecdhTool.GenerateSharedSecret(ecdhPriKey, serverPubKey) //服务器公钥和本地第一个私钥协商出安全密钥
	if err != nil {
		return []byte{}, err
	}
	return Sha256(secretKey), nil
}

// DealCertificateVerify 处理CertificateVerify数据: 校验服务器-判断是不是微信服务器，请求返回数据有没有被串改
func DealCertificateVerify(clientHelloData []byte, serverHelloData []byte, data []byte) (bool, error) {
	// 解析数据
	totalSize := baseutils.BytesToInt32(data[0:4])
	certificateVerify, err := CertificateVerifyDeSerialize(data[4 : 4+totalSize])
	if err != nil {
		return false, err
	}

	// 合并请求数据
	message := make([]byte, 0)
	message = append(message, clientHelloData[0:]...)
	message = append(message, serverHelloData[0:]...)
	message = Sha256(message)

	// 校验数据
	flag, err := ECDSAVerifyData(message, certificateVerify.Signature)
	if err != nil {
		return false, err
	}

	return flag, nil
}

// DealNewSessionTicket 处理NewSessionTicket数据
func DealNewSessionTicket(mmInfo *MMInfo, data []byte) error {
	// 解析数据
	totalSize := baseutils.BytesToInt32(data[0:4])
	newSessionTicket, err := NewSessionTicketDeSerialize(data[4 : 4+totalSize])
	if err != nil {
		return err
	}

	mmInfo.ShortPskList = newSessionTicket.PskList
	return nil
}

// ----------- 上面是握手阶段 -----------
// ----------- 接下来是发送请求 -----------

// EncryptedReqData EncryptedReqData
func EncryptedReqData(aesGcmParam *AesGcmParam, data []byte, recordHeadType byte, clientSeq uint32) ([]byte, error) {
	tmpNonce := GetNonce(aesGcmParam.Nonce, clientSeq)
	tmpHead := GetRecordDataByLength(recordHeadType, uint16(len(data)+0x10))
	tmpAad := []byte{0x00, 0x00, 0x00, 0x00}
	tmpAad = append(tmpAad, baseutils.Int32ToBytes(clientSeq)...)
	tmpAad = append(tmpAad, tmpHead[0:]...)
	encodeData, err := AesGcmEncrypt(aesGcmParam.AesKey, tmpNonce, tmpAad, data)
	if err != nil {
		return []byte{}, err
	}
	return encodeData, nil
}

// DecryptedRecvData 解析响应数据包
func DecryptedRecvData(aesGcmParam *AesGcmParam, recvItem *PackItem, serverSeq uint32) ([]byte, error) {
	tmpNonce := GetNonce(aesGcmParam.Nonce, serverSeq)
	tmpAad := []byte{0x00, 0x00, 0x00, 0x00}
	tmpAad = append(tmpAad, baseutils.Int32ToBytes(serverSeq)...)
	tmpAad = append(tmpAad, recvItem.RecordHead[0:]...)
	encodeData, err := AesGcmDecrypt(aesGcmParam.AesKey, tmpNonce, tmpAad, recvItem.PackData)
	if err != nil {
		return []byte{}, err
	}
	return encodeData, nil
}

// CreateSendPackItems 创建发送的请求项列表
func CreateSendPackItems(mmInfo *MMInfo, httpHandler *HTTPHandler) ([]*PackItem, error) {
	retItems := make([]*PackItem, 0)

	// ClientHelloItem
	clientHelloItem := &PackItem{}
	clientHello, err := CreateClientHelloData(mmInfo)
	if err != nil {
		return nil, err
	}
	clientHelloData := ClientHelloSerialize(clientHello)
	clientHelloItem.RecordHead = GetRecordDataByLength(ClientHandShakeType, uint16(len(clientHelloData)))
	clientHelloItem.PackData = clientHelloData

	// EncryptedExtensionsItem
	encryptedExtensionsItem := &PackItem{}
	encryptedExtensions := CreateEncryptedExtensions()
	encryptedExtensionsData := EncryptedExtensionsSerialize(encryptedExtensions)
	encryptedExtensionsItem.RecordHead = GetRecordDataByLength(ClientHandShakeType, uint16(len(encryptedExtensionsData)))
	encryptedExtensionsItem.PackData = encryptedExtensionsData

	// HTTPHandlerItem
	httpHandlerItem := &PackItem{}
	httpHandlerData := HTTPHandlerSerialize(httpHandler)
	httpHandlerItem.RecordHead = GetRecordDataByLength(BodyType, uint16(len(httpHandlerData)))
	httpHandlerItem.PackData = httpHandlerData

	// AlertItem
	alertItem := &PackItem{}
	alertData := GetAlertData()
	alertItem.RecordHead = GetRecordDataByLength(AlertType, uint16(len(alertData)))
	alertItem.PackData = alertData

	// 返回数据
	retItems = append(retItems, clientHelloItem)
	retItems = append(retItems, encryptedExtensionsItem)
	retItems = append(retItems, httpHandlerItem)
	retItems = append(retItems, alertItem)
	return retItems, nil
}

// MMHTTPPackData MMPackData
func MMHTTPPackData(mmInfo *MMInfo, items []*PackItem) ([]byte, error) {
	// 密钥扩展
	sha256Value := Sha256(items[0].PackData)
	expandSecretData := []byte("early data key expansion")
	expandSecretData = append(expandSecretData, sha256Value[0:]...)
	tmpHkdfValue := HkdfExpand(mmInfo.PskAccessKey, expandSecretData, 28)
	aesGcmParam := &AesGcmParam{}
	aesGcmParam.AesKey = tmpHkdfValue[0x00:0x10]
	aesGcmParam.Nonce = tmpHkdfValue[0x10:0x1c]

	// 加密EncryptedExtensions
	encryptData, err := EncryptedReqData(aesGcmParam, items[1].PackData, ClientHandShakeType, 1)
	if err != nil {
		return []byte{}, err
	}
	partData2 := CreateRecordData(ClientHandShakeType, encryptData)

	// 加密HTTPHandler
	httpHandlerEncryptData, err := EncryptedReqData(aesGcmParam, items[2].PackData, BodyType, 2)
	if err != nil {
		return []byte{}, err
	}
	partData3 := CreateRecordData(BodyType, httpHandlerEncryptData)

	// 加密Alert
	alertDataEncryptData, err := EncryptedReqData(aesGcmParam, items[3].PackData, AlertType, 3)
	if err != nil {
		return []byte{}, err
	}
	partData4 := CreateRecordData(AlertType, alertDataEncryptData)

	// 返回数据
	retData := make([]byte, 0)
	retData = append(retData, items[0].RecordHead[0:]...)
	retData = append(retData, items[0].PackData[0:]...)
	retData = append(retData, partData2[0:]...)
	retData = append(retData, partData3[0:]...)
	retData = append(retData, partData4[0:]...)
	return retData, err
}

// MMDecodeResponseData 解码响应数据
func MMDecodeResponseData(mmInfo *MMInfo, sendItems []*PackItem, respData []byte) ([]byte, error) {
	retData := make([]byte, 0)

	// 解析 对响应数据进行分包
	recvItems, err := ParserMMtlsResponseData(respData)
	if err != nil {
		return retData, err
	}

	if len(recvItems) < 4 {
		return retData, errors.New("MMDecodeResponseData err: recvItems Length < 4")
	}

	// 密钥扩展 用于后面的解密
	shaData := make([]byte, 0)
	shaData = append(shaData, sendItems[0].PackData[0:]...)
	shaData = append(shaData, sendItems[1].PackData[0:]...)
	shaData = append(shaData, recvItems[0].PackData[0:]...)
	sha256Value := Sha256(shaData)
	expandSecretData := []byte("handshake key expansion")
	expandSecretData = append(expandSecretData, sha256Value[0:]...)
	tmpHkdfValue := HkdfExpand(mmInfo.PskAccessKey, expandSecretData, 28)
	aesGcmParam := &AesGcmParam{}
	aesGcmParam.AesKey = tmpHkdfValue[0x00:0x10]
	aesGcmParam.Nonce = tmpHkdfValue[0x10:0x1c]

	// 解密剩下的包
	count := len(recvItems)
	for index := 1; index < count; index++ {
		// 解密Finished数据包
		decodeData, err := DecryptedRecvData(aesGcmParam, recvItems[index], uint32(index))
		if err != nil {
			return retData, err
		}

		// RecordHeadType
		recordHeadType := recvItems[index].RecordHead[0]
		// ServerHandShakeType 校验收到的数据是否完整，是否又被串改
		if recordHeadType == ServerHandShakeType {
			// 判断数据长度是否正常
			current := 0
			totalLength := int(baseutils.BytesToInt32(decodeData[current : current+4]))
			current = current + 4
			if totalLength < 0 {
				return retData, errors.New("MMDecodeResponseData err: totalLength < 0")
			}

			// ReceiveSubType
			subType := decodeData[current]
			// FinishedType 校验数据是否正常
			if subType == FinishedType {
				// 反序列化
				finished, err := FinishedDeSerialize(decodeData[current : current+totalLength])
				if err != nil {
					return retData, err
				}
				bSuccess := VerifyFinishedData(mmInfo.PskAccessKey, sha256Value, finished.VerifyData)
				if !bSuccess {
					return retData, errors.New("MMDecodeResponseData err: VerifyFinishedData failed")
				}
			}
		}

		// 解析响应数据
		if recordHeadType == BodyType {
			retData = append(retData, decodeData[0:]...)
		}

		// 解析AlertType
		if recordHeadType == AlertType {
			// 关闭连接的数据包(对长链接有用) 固定为0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x01
		}
	}

	return retData, nil
}

// VerifyFinishedData 校验服务端返回数据是否正确
func VerifyFinishedData(aesKey []byte, shaValue []byte, finishedData []byte) bool {
	count := len(finishedData)
	message := []byte("server finished")
	tmpHkdfValue := HkdfExpand(aesKey, message, count)
	verifyData := HmacHash256(tmpHkdfValue, shaValue)

	// 比对结果是否一致
	for index := 0; index < count; index++ {
		if verifyData[index] != finishedData[index] {
			return false
		}
	}

	return true
}
