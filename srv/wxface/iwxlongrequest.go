package wxface

// IWXLongRequest 微信请求
type IWXLongRequest interface {
	GetSeqId() uint32
	SetSeqId(seqId uint32)
	GetOpcode() uint32
	GetCgiUrl() string
	GetData() []byte
}
