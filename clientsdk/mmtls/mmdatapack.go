package mmtls

import (
	"feiyu.com/wx/clientsdk/baseutils"
)

// RecordHeadSerialize 序列化RecordHead
func RecordHeadSerialize(recordHead *RecordHead) []byte {
	retBytes := make([]byte, 0)
	// Type
	retBytes = append(retBytes, recordHead.Type)
	// Tag
	retBytes = append(retBytes, baseutils.Int16ToBytesBigEndian(recordHead.Tag)[0:]...)
	// size
	retBytes = append(retBytes, baseutils.Int16ToBytesBigEndian(recordHead.Size)[0:]...)

	return retBytes
}

// ClientHelloSerialize 序列化ClientHello
func ClientHelloSerialize(clientHello *ClientHello) []byte {
	bodyData := make([]byte, 0)
	// Type
	bodyData = append(bodyData, ClientHelloType)
	// Version
	bodyData = append(bodyData, baseutils.Int16ToBytesLittleEndian(clientHello.Version)[0:]...)
	// suiteCount
	suiteCount := byte(len(clientHello.CipherSuiteList))
	bodyData = append(bodyData, suiteCount)
	// suiteList
	suiteList := clientHello.CipherSuiteList
	for index := 0; index < int(suiteCount); index++ {
		// suiteCode
		bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(suiteList[index].SuiteCode)[0:]...)
	}
	// RandomBytes
	bodyData = append(bodyData, clientHello.RandomBytes[0:]...)
	// ClientGmtTime
	bodyData = append(bodyData, baseutils.Int32ToBytes(clientHello.ClientGmtTime)[0:]...)
	// Extensions
	extensionsData := ExtensionsSerialize(clientHello.ExtensionList)
	bodyData = append(bodyData, extensionsData[0:]...)

	// 返回数据
	retBytes := make([]byte, 0)
	totalLength := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(totalLength)[0:]...)
	retBytes = append(retBytes, bodyData[0:]...)
	return retBytes
}

// ExtensionsSerialize 序列化Extensions
func ExtensionsSerialize(extensionList []*Extension) []byte {
	retBytes := make([]byte, 0)

	// bodyData
	bodyData := make([]byte, 0)
	// MapSize
	extensionCount := byte(len(extensionList))
	bodyData = append(bodyData, extensionCount)
	for index := 0; index < int(extensionCount); index++ {
		// Extension TotalLength
		extensionLength := uint32(len(extensionList[index].ExtensionData))
		bodyData = append(bodyData, baseutils.Int32ToBytes(extensionLength)[0:]...)
		// extensionData
		bodyData = append(bodyData, extensionList[index].ExtensionData[0:]...)
	}

	// Extensions Size
	extensionsSize := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(extensionsSize)[0:]...)

	// ExtensionsData
	retBytes = append(retBytes, bodyData[0:]...)
	return retBytes
}

// PskSerialize 序列化Psk
func PskSerialize(psk *Psk) []byte {
	// BodyData
	bodyData := make([]byte, 0)
	// Type
	bodyData = append(bodyData, psk.Type)
	// TicketLifeTimeHint
	bodyData = append(bodyData, baseutils.Int32ToBytes(psk.TicketKLifeTimeHint)[0:]...)
	// MacValue
	macValueLen := uint16(len(psk.MacValue))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(macValueLen)[0:]...)
	bodyData = append(bodyData, psk.MacValue[0:]...)
	// KeyVersion
	bodyData = append(bodyData, baseutils.Int32ToBytes(psk.KeyVersion)[0:]...)
	// IV
	ivLen := uint16(len(psk.Iv))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(ivLen)[0:]...)
	bodyData = append(bodyData, psk.Iv[0:]...)
	// EncryptTicket
	encryptTicketLen := uint16(len(psk.EncryptedTicket))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(encryptTicketLen)[0:]...)
	bodyData = append(bodyData, psk.EncryptedTicket[0:]...)

	// 返回数据
	retBytes := make([]byte, 0)
	bodyLen := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(bodyLen)[0:]...)
	retBytes = append(retBytes, bodyData[0:]...)

	return retBytes
}

// ClientKeyOfferSerialize 序列化ClientKeyOffer
func ClientKeyOfferSerialize(clientKeyOffer *ClientKeyOffer) []byte {
	// BodyData
	bodyData := make([]byte, 0)
	// Version
	bodyData = append(bodyData, baseutils.Int32ToBytes(clientKeyOffer.Version)[0:]...)
	// PublicValue
	publicValueLen := uint16(len(clientKeyOffer.PublicValue))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(publicValueLen)[0:]...)
	if publicValueLen > 0 {
		bodyData = append(bodyData, clientKeyOffer.PublicValue[0:]...)
	}

	// 返回数据
	retBytes := make([]byte, 0)
	bodyDataLen := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(bodyDataLen)[0:]...)
	retBytes = append(retBytes, bodyData[0:]...)
	return retBytes
}

// EncryptedExtensionsSerialize 序列化EncryptedExtensions
func EncryptedExtensionsSerialize(encryptedExtensions *EncryptedExtensions) []byte {
	bodyData := make([]byte, 0)
	// Type
	bodyData = append(bodyData, EncryptedExtensionsType)
	// ExtensionList
	extensionsData := ExtensionsSerialize(encryptedExtensions.ExtensionList)
	bodyData = append(bodyData, extensionsData[0:]...)

	// 返回数据
	retBytes := make([]byte, 0)
	bodyDataLen := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(bodyDataLen)[0:]...)
	retBytes = append(retBytes, bodyData[0:]...)
	return retBytes
}

// HTTPHandlerSerialize 序列化HTTPHandler
func HTTPHandlerSerialize(httpHandler *HTTPHandler) []byte {
	bodyData := make([]byte, 0)
	// URL
	urlLength := uint16(len(httpHandler.URL))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(urlLength)[0:]...)
	bodyData = append(bodyData, []byte(httpHandler.URL)[0:]...)

	// Host
	hostLength := uint16(len(httpHandler.Host))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(hostLength)[0:]...)
	bodyData = append(bodyData, []byte(httpHandler.Host)[0:]...)

	// MMPkg
	mmpkgLength := uint32(len(httpHandler.MMPkg))
	bodyData = append(bodyData, baseutils.Int32ToBytes(mmpkgLength)[0:]...)
	bodyData = append(bodyData, httpHandler.MMPkg[0:]...)

	// 返回数据
	retBytes := make([]byte, 0)
	bodyDataLen := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(bodyDataLen)[0:]...)
	retBytes = append(retBytes, bodyData[0:]...)
	return retBytes
}

// FinishedSerialize 序列化Finished包
func FinishedSerialize(finished *Finished) []byte {
	bodyData := make([]byte, 0)
	// Type
	bodyData = append(bodyData, FinishedType)
	// VerifyData
	verifyDataLen := uint16(uint32(len(finished.VerifyData)))
	bodyData = append(bodyData, baseutils.Int16ToBytesBigEndian(verifyDataLen)[0:]...)
	bodyData = append(bodyData, finished.VerifyData[0:]...)

	// 返回数据
	retBytes := make([]byte, 0)
	bodyDataLen := uint32(len(bodyData))
	retBytes = append(retBytes, baseutils.Int32ToBytes(bodyDataLen)[0:]...)
	retBytes = append(retBytes, bodyData[0:]...)
	return retBytes
}

// --------------- 各类Extension的序列化 ---------------

// PreSharedKeyExtensionSerialize 序列化PreSharedKeyExtension
func PreSharedKeyExtensionSerialize(preSharedKeyExtension *PreSharedKeyExtension) *Extension {
	// ExtensionBytes
	extensionBytes := make([]byte, 0)
	// PreSharedKeyExtensionType
	extensionBytes = append(extensionBytes, baseutils.Int16ToBytesBigEndian(PreSharedKeyExtensionType)[0:]...)
	// pskCount
	pskCount := byte(len(preSharedKeyExtension.PskList))
	extensionBytes = append(extensionBytes, pskCount)
	// PskList
	for index := 0; index < int(pskCount); index++ {
		psk := preSharedKeyExtension.PskList[index]
		pskData := PskSerialize(psk)
		extensionBytes = append(extensionBytes, pskData[0:]...)
	}

	// 返回数据
	retExtension := &Extension{}
	retExtension.ExtensionType = PreSharedKeyExtensionType
	retExtension.ExtensionData = extensionBytes
	return retExtension
}

// ClientKeyShareExtensionSerialize 序列化ClientKeyShareExtension
func ClientKeyShareExtensionSerialize(clientKeyShareExtension *ClientKeyShareExtension) *Extension {
	// ExtensionBytes
	extensionBytes := make([]byte, 0)
	// ClientKeyShareType
	extensionBytes = append(extensionBytes, baseutils.Int16ToBytesBigEndian(ClientKeyShareType)[0:]...)
	// KeyOfferCount
	keyOfferCount := byte(len(clientKeyShareExtension.ClientKeyOfferList))
	extensionBytes = append(extensionBytes, keyOfferCount)
	// KeyOfferList
	for index := 0; index < int(keyOfferCount); index++ {
		keyOffer := clientKeyShareExtension.ClientKeyOfferList[index]
		keyOfferData := ClientKeyOfferSerialize(keyOffer)
		extensionBytes = append(extensionBytes, keyOfferData[0:]...)
	}
	// CertificateVersion
	extensionBytes = append(extensionBytes, baseutils.Int32ToBytes(clientKeyShareExtension.CertificateVersion)[0:]...)

	// 返回数据
	retExtension := &Extension{}
	retExtension.ExtensionType = ClientKeyShareType
	retExtension.ExtensionData = extensionBytes
	return retExtension
}

// EarlyEncryptDataExtensionSerialize 序列化EarlyEncryptDataExtension
func EarlyEncryptDataExtensionSerialize(earlyEncryptDataExtension *EarlyEncryptDataExtension) *Extension {
	// ExtensionBytes
	extensionBytes := make([]byte, 0)
	// ClientKeyShareType
	extensionBytes = append(extensionBytes, baseutils.Int16ToBytesBigEndian(EarlyEncryptDataType)[0:]...)
	// ClientGmtTime
	extensionBytes = append(extensionBytes, baseutils.Int32ToBytes(earlyEncryptDataExtension.ClientGmtTime)[0:]...)

	// 返回数据
	retExtension := &Extension{}
	retExtension.ExtensionType = EarlyEncryptDataType
	retExtension.ExtensionData = extensionBytes
	return retExtension
}

// LongPackHeaderInfoSerialize 序列化LongPackHeaderInfo
func LongPackHeaderInfoSerialize(packTagInfo *LongPackHeaderInfo) []byte {
	retBytes := make([]byte, 0)
	// Type
	retBytes = append(retBytes, baseutils.Int16ToBytesBigEndian(packTagInfo.HeaderLen)[0:]...)
	// Version
	retBytes = append(retBytes, baseutils.Int16ToBytesBigEndian(packTagInfo.Version)[0:]...)
	// Operation
	retBytes = append(retBytes, baseutils.Int32ToBytes(packTagInfo.Operation)[0:]...)
	// SequenceNumber
	retBytes = append(retBytes, baseutils.Int32ToBytes(packTagInfo.SequenceNumber)[0:]...)

	return retBytes
}
