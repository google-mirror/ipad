package clientsdk

// WXLongRequest 微信长链接请求
type WXLongRequest struct {
	SeqId  uint32
	OpCode uint32
	CgiUrl string
	Data   []byte
}

// GetSeqId 获取SeqId
func (wxlq *WXLongRequest) GetSeqId() uint32 {
	return wxlq.SeqId
}

// SetSeqId 设置SeqId
func (wxlq *WXLongRequest) SetSeqId(seqId uint32) {
	wxlq.SeqId = seqId
}

// GetOpcode 获取Opcode
func (wxlq *WXLongRequest) GetOpcode() uint32 {
	return wxlq.OpCode
}

// GetOpcode 获取Opcode
func (wxlq *WXLongRequest) GetCgiUrl() string {
	return wxlq.CgiUrl
}

// GetData 获取数据
func (wxlq *WXLongRequest) GetData() []byte {
	return wxlq.Data
}
