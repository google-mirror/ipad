package mmtls

const (
	// ClientHandShakeType ClientHandShakeType
	ClientHandShakeType byte = 25
	// ServerHandShakeType ServerHandShakeType
	ServerHandShakeType byte = 22
	// BodyType BodyType
	BodyType byte = 23
	// AlertType AlertType
	AlertType byte = 21

	// ClientHelloType ClientHelloType
	ClientHelloType byte = 1
	// ServerHelloType ServerHelloType
	ServerHelloType byte = 2
	// NewSessionTicketType NewSessionTicketType
	NewSessionTicketType byte = 4
	// EncryptedExtensionsType EncryptedExtensionsType
	EncryptedExtensionsType byte = 8
	// CertificateVerifyType CertificateVerifyType
	CertificateVerifyType byte = 15
	// FinishedType FinishedType
	FinishedType byte = 20

	// PreSharedKeyExtensionType PreSharedKeyExtensionType
	PreSharedKeyExtensionType uint16 = 15
	// ClientKeyShareType ClientKeyShareType
	ClientKeyShareType uint16 = 16
	// ServerKeyShareType ServerKeyShareType
	ServerKeyShareType uint16 = 17
	// EarlyEncryptDataType EarlyEncryptDataType
	EarlyEncryptDataType uint16 = 18

	// MaxCipherSuiteSize MaxCipherSuiteSize
	MaxCipherSuiteSize uint32 = 2
	// FixedRandomSize FixedRandomSize
	FixedRandomSize uint32 = 32
	// MaxNewSessionTicketPskSize MaxNewSessionTicketPskSize
	MaxNewSessionTicketPskSize uint32 = 2
	// MaxSignatureSize MaxSignatureSize
	MaxSignatureSize uint32 = 2048
	// MaxFinishedVerifyDataSize MaxFinishedVerifyDataSize
	MaxFinishedVerifyDataSize uint32 = 2048
	// MaxDataPackSize MaxDataPackSize
	MaxDataPackSize uint32 = 0x8000000

	// MaxExtensionSize MaxExtensionSize
	MaxExtensionSize uint32 = 256
	// MaxKeyOfferSize MaxKeyOfferSize
	MaxKeyOfferSize uint32 = 256
	// MaxPublicValueSize MaxPublicValueSize
	MaxPublicValueSize uint32 = 256

	// FixedRecordHeadSize FixedRecordHeadSize
	FixedRecordHeadSize uint32 = 5

	// MMLongVersion 长链接版本号
	MMLongVersion uint16 = 1

	// MMLongOperationSmartHeartBeat 心跳包
	MMLongOperationSmartHeartBeat uint32 = 0x06
	// MMLongOperationSmartHeartBeatBackUp 心跳包请求-备用
	MMLongOperationSmartHeartBeatBackUp uint32 = 0x0c
	// MMLongOperationServerPush 服务端推送消息
	MMLongOperationServerPush uint32 = 0x7a
	// MMLongOperationGetOnlineInfo 发起GetOnlineInfo请求
	MMLongOperationGetOnlineInfo uint32 = 0xcd
	// MMLongOperationCheckQrcode 检测二维码状态请求
	MMLongOperationGetQrcode       uint32 = 232  //502
	MMLongOperationCheckQrcode     uint32 = 0xe9 // 503
	MMLongOperationGetProfile      uint32 = 118
	MMLongOperationSendMessage     uint32 = 2
	MMLongOperationNewSendMessage  uint32 = 237
	MMLongOperationBatchGetFavItem uint32 = 32769
	MMLongOperationSync            uint32 = 26
	MMLongOperationSnsPost         uint32 = 97
	MMLongOperationSnsTimeLine     uint32 = 98
	// MMLongOperationRequest 发起请求
	MMLongOperationRequest uint32 = 0xed
	// MMLongOperationHeartBeat 发送心跳包请求
	MMLongOperationHeartBeat uint32 = 518 //   0xee
	//MMLongOperationSystemPush 系统推送
	MMLongOperationSystemPush uint32 = 0x18

	// MMLongSystemPushTypeSync 系统推送：需要同步
	MMLongSystemPushTypeSync int32 = 2
	// MMLongSystemPushTypeLogout 系统推送：被退出登陆
	MMLongSystemPushTypeLogout int32 = -1
)
