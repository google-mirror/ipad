package wxface

type IWXUserMsgMgr interface {
	Start()
	Stop()
	AddNewTextMsg(MsgId, newMsg, toUSerName string)
	AddImageMsg(MsgId string, imgData []byte, toUSerName string)
}
