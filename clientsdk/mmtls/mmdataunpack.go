package mmtls

import (
	"errors"

	"feiyu.com/wx/clientsdk/baseutils"
)

// DataPackerDeserialize 反序列化数据
func DataPackerDeserialize(data []byte) {
	totalLength := uint32(len(data))
	current := uint32(0)

	for current < totalLength {
		// recordHead
		recordHead := RecordHeadDeSerialize(data[current:])
		baseutils.ShowObjectValue(recordHead)
		current = current + 5
		offset := uint32(0)
		// tmpType
		pkgSize := baseutils.BytesToInt32(data[current+offset : current+offset+4])
		offset = offset + 4
		tmpType := data[current+offset]

		// ClientHelloType
		if tmpType == ClientHelloType {
			clientHello, _ := ClientHelloDeSerialize(data[current+offset : current+offset+pkgSize])
			baseutils.ShowObjectValue(clientHello)
			ShowMMTLSExtensions(clientHello.ExtensionList)
		}

		// ServerHelloType
		if tmpType == ServerHelloType {
			serverHello, _ := ServerHelloDeSerialize(data[current+offset : current+offset+pkgSize])
			baseutils.ShowObjectValue(serverHello)
			ShowMMTLSExtensions(serverHello.ExtensionList)
		}
		offset = offset + pkgSize
		current = current + offset
	}
}

// GetCipherSuiteInfoByCode 返回code对应的 CipherSuiteInfo
func GetCipherSuiteInfoByCode(code uint16) *CipherSuite {
	if code == 0xc02b {
		cipherSuite := &CipherSuite{}
		// SuiteCode
		cipherSuite.SuiteCode = code

		// cipherSuiteInfo
		cipherSuiteInfo := &CipherSuiteInfo{}
		cipherSuiteInfo.SuiteCode = code
		cipherSuiteInfo.Clipher1 = "ECDHE"
		cipherSuiteInfo.Clipher2 = "ECDSA"
		cipherSuiteInfo.Clipher3 = "SHA256"
		cipherSuiteInfo.Clipher4 = "AES_128_GCM"
		cipherSuiteInfo.Clipher5 = "AEAD"
		cipherSuiteInfo.Length1 = 16
		cipherSuiteInfo.Length2 = 0
		cipherSuiteInfo.Length3 = 12
		cipherSuite.SuiteInfo = cipherSuiteInfo
		return cipherSuite
	}

	if code == 0x00a8 {
		cipherSuite := &CipherSuite{}
		// SuiteCode
		cipherSuite.SuiteCode = code

		// cipherSuiteInfo
		cipherSuiteInfo := &CipherSuiteInfo{}
		cipherSuiteInfo.SuiteCode = code
		cipherSuiteInfo.Clipher1 = "PSK"
		cipherSuiteInfo.Clipher2 = "ECDSA"
		cipherSuiteInfo.Clipher3 = "SHA256"
		cipherSuiteInfo.Clipher4 = "AES_128_GCM"
		cipherSuiteInfo.Clipher5 = "AEAD"
		cipherSuiteInfo.Length1 = 16
		cipherSuiteInfo.Length2 = 0
		cipherSuiteInfo.Length3 = 12
		cipherSuite.SuiteInfo = cipherSuiteInfo
		return cipherSuite
	}

	return nil
}

// RecordHeadDeSerialize 反序列化RecordHead
func RecordHeadDeSerialize(data []byte) *RecordHead {
	retRecordHead := &RecordHead{}
	// 偏移
	current := uint32(0)

	// Type
	retRecordHead.Type = data[current]
	current = current + 1

	// Tag
	retRecordHead.Tag = baseutils.BytesToUint16BigEndian(data[current : current+2])
	current = current + 2

	// Size
	retRecordHead.Size = baseutils.BytesToUint16BigEndian(data[current : current+2])

	return retRecordHead
}

// ClientHelloDeSerialize 反序列化ClientHello
func ClientHelloDeSerialize(data []byte) (*ClientHello, error) {
	current := uint32(0)
	clientHello := &ClientHello{}

	tmptype := data[current]
	if tmptype != ClientHelloType {
		return nil, errors.New("ClientHelloDeSerialize err: not ClientHelloType")
	}
	current = current + 1

	// Version
	clientHello.Version = baseutils.BytesToUint16LittleEndian(data[current : current+2])
	current = current + 2

	// ciphersuiteSize
	ciphersuiteSize := data[current]
	current = current + 1

	// CipherSuiteList
	clientHello.CipherSuiteList = make([]*CipherSuite, ciphersuiteSize)
	for index := 0; index < int(ciphersuiteSize); index++ {
		suitecode := baseutils.BytesToUint16BigEndian(data[current : current+2])
		current = current + 2
		clientHello.CipherSuiteList[index] = GetCipherSuiteInfoByCode(suitecode)
	}

	// RandomBytes
	clientHello.RandomBytes = data[current : current+32]
	current = current + 32

	// ClientGmtTime
	clientHello.ClientGmtTime = baseutils.BytesToInt32(data[current : current+4])
	current = current + 4

	// ExtensionList
	clientHello.ExtensionList = ExtensionsDeSerialize(data[current:])

	return clientHello, nil
}

// ExtensionsDeSerialize 反序列号Extensions
func ExtensionsDeSerialize(data []byte) []*Extension {
	// 初始化返回数组
	retExtensions := make([]*Extension, 0)

	// 初始化索引
	current := uint32(0)

	// totalLength
	totalLength := baseutils.BytesToInt32(data[current : current+4])
	current = current + 4

	// extensionCount := data[current]
	current = current + 1
	// ExtensionList
	for current-4 < totalLength {
		extension := &Extension{}

		// ExtensionData
		extensionLength := baseutils.BytesToInt32(data[current : current+4])
		current = current + 4

		// ExtensionType
		extension.ExtensionType = baseutils.BytesToUint16BigEndian(data[current : current+2])

		// ExtensionData
		extension.ExtensionData = data[current : current+extensionLength]

		// 放入列表
		retExtensions = append(retExtensions, extension)
		current = current + extensionLength
	}

	return retExtensions
}

// PskDeSerialize 反序列化Psk
func PskDeSerialize(data []byte) (*Psk, error) {
	tmpPsk := &Psk{}

	// current
	current := uint32(0)

	// Type
	tmpPsk.Type = data[current]
	current = current + 1

	// TicketKLifeTimeHint
	tmpPsk.TicketKLifeTimeHint = baseutils.BytesToInt32(data[current : current+4])
	current = current + 4

	// MacValue
	macValueLength := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	tmpPsk.MacValue = data[current : current+macValueLength]
	current = current + macValueLength

	// KeyVersion
	tmpPsk.KeyVersion = baseutils.BytesToInt32(data[current : current+4])
	current = current + 4

	// IV
	ivLength := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	tmpPsk.Iv = data[current : current+ivLength]
	current = current + ivLength

	// EncryptedTicket
	encryptedTicketLength := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	tmpPsk.EncryptedTicket = data[current : current+encryptedTicketLength]
	current = current + encryptedTicketLength
	if current != uint32(len(data)) {
		return nil, errors.New("PskDeSerialize err: current - startPos != pskTotalLength")
	}

	return tmpPsk, nil
}

// EncryptedExtensionsDeSerialize 反序列化EncryptedExtensions
func EncryptedExtensionsDeSerialize(data []byte) (*EncryptedExtensions, error) {
	current := uint32(0)
	retEncryptedExtensions := &EncryptedExtensions{}

	// Type
	tmptype := data[current]
	if tmptype != EncryptedExtensionsType {
		return nil, errors.New("EncryptedExtensionsDeSerialize err: not ServerHelloType")
	}
	current = current + 1

	// ExtensionList
	retEncryptedExtensions.ExtensionList = ExtensionsDeSerialize(data[current:])

	return retEncryptedExtensions, nil
}

// ServerHelloDeSerialize 反序列化ServerHello
func ServerHelloDeSerialize(data []byte) (*ServerHello, error) {
	current := uint32(0)
	serverHello := &ServerHello{}

	tmptype := data[current]
	if tmptype != ServerHelloType {
		return nil, errors.New("ServerHelloDeSerialize err: not ServerHelloType")
	}
	current = current + 1

	// Version
	serverHello.Version = baseutils.BytesToUint16LittleEndian(data[current : current+2])
	current = current + 2

	// CipherSuite
	suiteCode := baseutils.BytesToUint16BigEndian(data[current : current+2])
	current = current + 2
	serverHello.CipherSuite = GetCipherSuiteInfoByCode(suiteCode)

	// RandomBytes
	serverHello.RandomBytes = data[current : current+32]
	current = current + 32

	// ExtensionList
	serverHello.ExtensionList = ExtensionsDeSerialize(data[current:])

	return serverHello, nil
}

// NewSessionTicketDeSerialize 反序列化NewSessionTicket
func NewSessionTicketDeSerialize(data []byte) (*NewSessionTicket, error) {
	retNewSessionTicket := &NewSessionTicket{}

	// current
	current := uint32(0)

	// tmpType
	tmpType := data[current]
	if tmpType != NewSessionTicketType {
		return nil, errors.New("NewSessionTicketDeSerialize err: not NewSessionTicketType")
	}
	current = current + 1

	// pskListSize
	pskListSize := data[current]
	current = current + 1

	// PskList
	retNewSessionTicket.PskList = make([]*Psk, pskListSize)
	for index := 0; index < int(pskListSize); index++ {
		// pskTotalLength
		pskTotalLength := baseutils.BytesToInt32(data[current : current+4])
		current = current + 4

		// PskDeSerialize
		retPsk, err := PskDeSerialize(data[current : current+pskTotalLength])
		if err != nil {
			return nil, err
		}

		// Add to PskList
		retNewSessionTicket.PskList[index] = retPsk
		current = current + pskTotalLength
	}

	return retNewSessionTicket, nil
}

// CertificateVerifyDeSerialize CertificateVerifyDeSerialize
func CertificateVerifyDeSerialize(data []byte) (*CertificateVerify, error) {
	retCertificateVerify := &CertificateVerify{}

	// current
	current := uint32(0)

	// tmpType
	tmpType := data[current]
	if tmpType != CertificateVerifyType {
		return nil, errors.New("CertificateVerifyDeSerialize err: not CertificateVerifyType")
	}
	current = current + 1

	// SignatureSize
	size := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2

	// Signature
	retCertificateVerify.Signature = data[current : current+size]
	current = current + size

	// 判断数据是否完整解析
	if current != uint32(len(data)) {
		return nil, errors.New("CertificateVerifyDeSerialize err: current != uint32(len(data)")
	}

	return retCertificateVerify, nil
}

// HTTPHandlerDeSerialize 反序列化HttpHandler
func HTTPHandlerDeSerialize(data []byte) (*HTTPHandler, error) {
	retHTTPHandler := &HTTPHandler{}

	current := 0
	// URL
	urlLength := int(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	if urlLength > 0 {
		retHTTPHandler.URL = string(data[current : current+urlLength])
		current = current + int(urlLength)
	}

	// Host
	hostLength := int(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	if urlLength > 0 {
		retHTTPHandler.Host = string(data[current : current+hostLength])
		current = current + int(hostLength)
	}

	// MMPkg
	mmpkgLength := int(baseutils.BytesToInt32(data[current : current+4]))
	current = current + 4
	if mmpkgLength > 0 {
		retHTTPHandler.MMPkg = data[current : current+int(mmpkgLength)]
		current = current + mmpkgLength
	}

	// 判断数据是否正常完整解析
	if current != len(data) {
		return nil, errors.New("HTTPHandlerDeSerialize err: current != len(data)")
	}

	return retHTTPHandler, nil
}

// FinishedDeSerialize 反序列化 Finished
func FinishedDeSerialize(data []byte) (*Finished, error) {
	retFinished := &Finished{}

	// current
	current := uint32(0)

	// tmpType
	tmpType := data[current]
	if tmpType != FinishedType {
		return nil, errors.New("FinishedDeSerialize err: not FinishedType")
	}
	current = current + 1

	// VerifyData
	verifyDataLen := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	retFinished.VerifyData = data[current : current+verifyDataLen]
	current = current + verifyDataLen

	// 判断数据是否完整解析
	if current != uint32(len(data)) {
		return nil, errors.New("FinishedDeSerialize err: current != uint32(len(data)")
	}

	return retFinished, nil
}

// ------------------ 各类Extension反序列化 ------------------

// PreSharedKeyExtensionDeSerialize 反序列化PreSharedKeyExtension
func PreSharedKeyExtensionDeSerialize(data []byte) (*PreSharedKeyExtension, error) {
	retPreSharedKeyExtensions := &PreSharedKeyExtension{}
	current := uint32(0)
	tmpType := baseutils.BytesToUint16BigEndian(data[current : current+2])
	if tmpType != PreSharedKeyExtensionType {
		return nil, errors.New("PreSharedKeyExtensionDeSerialize err: not PreSharedKeyExtensionType")
	}
	current = current + 2

	// pskCount
	pskCount := data[current]
	current = current + 1

	// PskList
	retPreSharedKeyExtensions.PskList = make([]*Psk, pskCount)
	for index := 0; index < int(pskCount); index++ {
		// pskTotalLength
		pskTotalLength := baseutils.BytesToInt32(data[current : current+4])
		current = current + 4

		// PskDeSerialize
		retPsk, err := PskDeSerialize(data[current : current+pskTotalLength])
		if err != nil {
			return nil, err
		}

		// Add to PskList
		retPreSharedKeyExtensions.PskList[index] = retPsk
		current = current + pskTotalLength
	}

	return retPreSharedKeyExtensions, nil
}

// ClientKeyShareExtensionDeSerialize 解析ClientKeyShareExtension
func ClientKeyShareExtensionDeSerialize(data []byte) (*ClientKeyShareExtension, error) {
	retClientKeyShareExtension := &ClientKeyShareExtension{}
	current := uint32(0)
	tmpType := baseutils.BytesToUint16BigEndian(data[current : current+2])
	if tmpType != ClientKeyShareType {
		return nil, errors.New("ClientKeyShareExtensionDeSerialize err: tmpType != ClientKeyShareType")
	}
	current = current + 2

	// clientKeyOfferCount
	clientKeyOfferCount := data[current]
	current = current + 1

	retClientKeyShareExtension.ClientKeyOfferList = make([]*ClientKeyOffer, clientKeyOfferCount)
	// ClientKeyOfferList
	for index := 0; index < int(clientKeyOfferCount); index++ {
		clientKeyOffer := &ClientKeyOffer{}
		clientKeyOfferTotalLength := baseutils.BytesToInt32(data[current : current+4])
		current = current + 4
		startPos := current

		// Version
		clientKeyOffer.Version = baseutils.BytesToInt32(data[current : current+4])
		current = current + 4

		// PublicValue
		publicValueSize := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
		current = current + 2
		clientKeyOffer.PublicValue = data[current : current+publicValueSize]
		current = current + publicValueSize

		if current-startPos != clientKeyOfferTotalLength {
			return nil, errors.New("ClientKeyShareExtensionDeSerialize err: current - startPos != clientKeyOfferTotalLength")
		}
		retClientKeyShareExtension.ClientKeyOfferList[index] = clientKeyOffer
	}

	// CertificateVersion
	retClientKeyShareExtension.CertificateVersion = baseutils.BytesToInt32(data[current : current+4])

	return retClientKeyShareExtension, nil
}

// ServerKeyShareExtensionDeSerialize 反序列化ServerKeyShareExtension
func ServerKeyShareExtensionDeSerialize(data []byte) (*ServerKeyShareExtension, error) {
	retServerKeyShareExtension := &ServerKeyShareExtension{}
	// 索引
	current := uint32(0)

	// tmpType
	tmpType := baseutils.BytesToUint16BigEndian(data[current : current+2])
	if tmpType != ServerKeyShareType {
		return nil, errors.New("ServerKeyShareExtensionDeSerialize err: tmpType != ServerKeyShareType")
	}
	current = current + 2

	// KeyOfferNameGroup
	retServerKeyShareExtension.KeyOfferNameGroup = baseutils.BytesToInt32(data[current : current+4])
	current = current + 4

	// PublicValue
	publicValueSize := uint32(baseutils.BytesToUint16BigEndian(data[current : current+2]))
	current = current + 2
	retServerKeyShareExtension.PublicValue = data[current : current+publicValueSize]

	return retServerKeyShareExtension, nil
}

// EarlyEncryptedDataExtensionDeSerialize 反序列化EarlyEncryptedDataExtension
func EarlyEncryptedDataExtensionDeSerialize(data []byte) (*EarlyEncryptDataExtension, error) {
	retEarlyEncryptDataExtension := &EarlyEncryptDataExtension{}
	// 索引
	current := uint32(0)
	// Type
	tmpType := baseutils.BytesToUint16BigEndian(data[current : current+2])
	if tmpType != EarlyEncryptDataType {
		return nil, errors.New("EarlyEncryptedDataExtensionDeSerialize err: tmpType != EarlyEncryptDataType")
	}
	current = current + 2
	// ClientGmtTime
	retEarlyEncryptDataExtension.ClientGmtTime = baseutils.BytesToInt32(data[current : current+4])

	return retEarlyEncryptDataExtension, nil
}

// LongPackHeaderInfoDeSerialize 序列化LongPackHeaderInfo
func LongPackHeaderInfoDeSerialize(data []byte) (*LongPackHeaderInfo, error) {
	retLongPackHeaderInfo := &LongPackHeaderInfo{}

	current := 0
	// HeaderLen
	retLongPackHeaderInfo.HeaderLen = baseutils.BytesToUint16BigEndian(data[current : current+2])
	current = current + 2
	// Version
	retLongPackHeaderInfo.Version = baseutils.BytesToUint16BigEndian(data[current : current+2])
	current = current + 2
	// Operation
	retLongPackHeaderInfo.Operation = baseutils.BytesToInt32(data[current : current+4])
	current = current + 4
	// SequenceNumber
	retLongPackHeaderInfo.SequenceNumber = baseutils.BytesToInt32(data[current : current+4])
	current = current + 4

	if current != len(data) {
		return retLongPackHeaderInfo, errors.New("LongPackHeaderInfoDeSerialize err: current != len(data)")
	}
	return retLongPackHeaderInfo, nil
}
