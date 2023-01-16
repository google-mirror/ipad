package baseinfo

// SysMsg 系统消息
type SysMsg struct {
	Type      string    `xml:"type,attr"`
	RevokeMsg RevokeMsg `xml:"revokemsg"`
}

// RevokeMsg 撤回消息
type RevokeMsg struct {
	Session    string `xml:"session"`
	MsgID      uint32 `xml:"msgid"`
	NewMsgID   int64  `xml:"newmsgid"`
	ReplaceMsg string `xml:"replacemsg"`
}
