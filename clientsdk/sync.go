package clientsdk

import "feiyu.com/wx/protobuf/wechat"

type SyncResponseV struct {
	ModUserInfos    []wechat.ModUserInfo    //CmdId = 1
	ModContacts     []wechat.ModContact     //CmdId = 2
	DelContacts     []wechat.DelContact     //CmdId = 4
	ModUserImgs     []wechat.ModUserImg     //CmdId = 35
	FunctionSwitchs []wechat.FunctionSwitch //CmdId = 23
	UserInfoExts    []wechat.UserInfoExt    //CmdId = 44
	AddMsgs         []wechat.AddMsg         //CmdId = 5
	ContinueFlag    uint32
	KeyBuf          wechat.SKBuiltinBufferT
	Status          uint32
	Continue        uint32
	Time            uint32
	UnknownCmdId    string
	Remarks         string
}
