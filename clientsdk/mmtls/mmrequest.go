package mmtls

import (
	"errors"
	"time"

	"feiyu.com/wx/clientsdk/baseutils"
)

// CreateRecordData 根据请求创建完整的mmtls数据包
func CreateRecordData(recordType byte, data []byte) []byte {
	recordHead := &RecordHead{}
	recordHead.Type = recordType
	recordHead.Tag = 0xF103
	recordHead.Size = uint16(len(data))

	// 组包返回
	retBytes := make([]byte, 0)
	retBytes = append(retBytes, RecordHeadSerialize(recordHead)[0:]...)
	retBytes = append(retBytes, data[0:]...)
	return retBytes
}

// GetRecordDataByLength 根据长度获取RecordData
func GetRecordDataByLength(recordType byte, len uint16) []byte {
	recordHead := &RecordHead{}
	recordHead.Type = recordType
	recordHead.Tag = 0xF103
	recordHead.Size = len
	return RecordHeadSerialize(recordHead)
}

// CreateHandShakeClientHelloData 创建ClientHello数据包
func CreateHandShakeClientHelloData(mmInfo *MMInfo) []byte {
	clientHello := &ClientHello{}
	// Version
	clientHello.Version = 0xF103
	// CipherSuiteList
	clientHello.CipherSuiteList = make([]*CipherSuite, 2)
	cipherSuite1 := &CipherSuite{}
	cipherSuite1.SuiteCode = 0xC02B
	clientHello.CipherSuiteList[0] = cipherSuite1
	cipherSuite2 := &CipherSuite{}
	cipherSuite2.SuiteCode = 0xA8
	clientHello.CipherSuiteList[1] = cipherSuite2
	// RandomBytes
	clientHello.RandomBytes = baseutils.RandomBytes(32)
	// ClientGmtTime
	clientHello.ClientGmtTime = (uint32)(time.Now().UnixNano() / 1000000000)
	// ExtensionList
	extensionList := make([]*Extension, 0)
	pskCount := len(mmInfo.ShortPskList)
	if pskCount > 1 {
		// 握手是用最后一个
		extensionList = append(extensionList, CreatePreSharedKeyExtensionData(mmInfo.ShortPskList[pskCount-1]))
	}

	// 随机生成两队ECDHKey
	extensionList = append(extensionList, CreateClientKeyShareExtensionData(mmInfo.ClientEcdhKeys))
	clientHello.ExtensionList = extensionList

	return ClientHelloSerialize(clientHello)
}

// CreatePreSharedKeyExtensionData 创建CreatePreSharedKeyExtension数据
func CreatePreSharedKeyExtensionData(psk *Psk) *Extension {
	preSharedKeyExtension := &PreSharedKeyExtension{}
	preSharedKeyExtension.PskList = make([]*Psk, 1)
	// 选取前面协商的最后一个Psk
	preSharedKeyExtension.PskList[0] = psk
	// 序列化
	retExtension := PreSharedKeyExtensionSerialize(preSharedKeyExtension)
	return retExtension
}

// CreateClientKeyShareExtensionData 创建ClientKeyShareExtension数据
func CreateClientKeyShareExtensionData(clientEcdhKeys *ClientEcdhKeys) *Extension {
	// retExtension
	clientKeyShareExtension := &ClientKeyShareExtension{}

	// ClientKeyOfferList
	clientKeyShareExtension.ClientKeyOfferList = make([]*ClientKeyOffer, 2)
	// 随机第一个EcdhKey
	clientKeyShareExtension.ClientKeyOfferList[0] = CreateClientKeyOfferData(1, clientEcdhKeys.PubKeyBuf1)
	// 随机第一个EcdhKey
	clientKeyShareExtension.ClientKeyOfferList[1] = CreateClientKeyOfferData(2, clientEcdhKeys.PubKeyBuf2)
	// CertificateVersion
	clientKeyShareExtension.CertificateVersion = 1

	// 返回序列化的ClientKeyShareExtension
	return ClientKeyShareExtensionSerialize(clientKeyShareExtension)
}

// CreateClientKeyOfferData 创建CreateClientKeyOffer数据
func CreateClientKeyOfferData(version uint32, publicKey []byte) *ClientKeyOffer {
	clientKeyOffser := &ClientKeyOffer{}
	clientKeyOffser.PublicValue = publicKey
	clientKeyOffser.Version = version

	return clientKeyOffser
}

// CreateClientHelloData 创建ClientHello数据包
func CreateClientHelloData(mmInfo *MMInfo) (*ClientHello, error) {
	clientHello := &ClientHello{}
	// Version
	clientHello.Version = 0xF103
	// CipherSuiteList
	clientHello.CipherSuiteList = make([]*CipherSuite, 1)
	cipherSuite := &CipherSuite{}
	cipherSuite.SuiteCode = 0xA8
	clientHello.CipherSuiteList[0] = cipherSuite
	// RandomBytes
	clientHello.RandomBytes = baseutils.RandomBytes(32)
	// ClientGmtTime
	clientHello.ClientGmtTime = (uint32)(time.Now().UnixNano() / 1000000000)
	// ExtensionList
	extensionList := make([]*Extension, 0)
	pskCount := len(mmInfo.ShortPskList)
	if pskCount <= 0 {
		return nil, errors.New("CreateClientHelloData error: mmInfo.PskList empty")
	}
	extensionList = append(extensionList, CreatePreSharedKeyExtensionData(mmInfo.ShortPskList[0]))
	clientHello.ExtensionList = extensionList

	return clientHello, nil
}

// CreateEarlyEncryptDataExtension 创建EarlyEncryptDataExtension
func CreateEarlyEncryptDataExtension() *Extension {
	retEarlyEncryptDataExtension := &EarlyEncryptDataExtension{}
	retEarlyEncryptDataExtension.ClientGmtTime = (uint32)(time.Now().UnixNano() / 1000000000)
	return EarlyEncryptDataExtensionSerialize(retEarlyEncryptDataExtension)
}

// CreateEncryptedExtensions 创建EncryptedExtensions
func CreateEncryptedExtensions() *EncryptedExtensions {
	retEncryptedExtensions := &EncryptedExtensions{}

	// ExtensionList
	retEncryptedExtensions.ExtensionList = make([]*Extension, 1)
	retEncryptedExtensions.ExtensionList[0] = CreateEarlyEncryptDataExtension()

	return retEncryptedExtensions
}

// CreateFinished 创建Finish包
func CreateFinished(verifyData []byte) *Finished {
	retFinished := &Finished{}
	retFinished.VerifyData = verifyData
	return retFinished
}

// GetAlertData 获取Alert数据
func GetAlertData() []byte {
	return []byte{0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x01}
}
