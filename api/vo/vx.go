package vo

import "feiyu.com/wx/protobuf/wechat"

type QYChatroomContactVo struct {
	List []*wechat.QYChatroomContactInfo
	Key  string
}
type DownloadVoiceData struct {
	Base64      []byte
	VoiceLength uint32
}
type GroupData struct {
	Count int64
	List  interface{}
}
