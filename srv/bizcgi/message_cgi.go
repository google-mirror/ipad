package bizcgi

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
)

// SendTextMsgReq 打包发送消息
func SendTextMsgReq(wxAccount wxface.IWXAccount, toUserName string, content string, atWxIDList []string, ContentType int) (*wechat.NewSendMsgResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	reqData := clientsdk.GetTextMsReq(userInfo, toUserName, content, atWxIDList, ContentType)
	longReq := &clientsdk.WXLongRequest{
		OpCode: mmtls.MMLongOperationNewSendMessage,
		CgiUrl: "/cgi-bin/micromsg-bin/newsendmsg",
		Data:   reqData,
	}
	// 发送消息
	response, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		return nil, err
	}
	newSendMsgResponse := response.(wechat.NewSendMsgResponse)
	return &newSendMsgResponse, nil
}

// SendShareCardReq 打包分享名片
func SendShareCardReq(wxAccount wxface.IWXAccount, toUserName string, id string, nickname string, alias string) (*wechat.NewSendMsgResponse, error) {
	contentType := 42
	if nickname == "" {
		nickname = id
	}
	content := "<?xml version=\"1.0\"?>\n<msg bigheadimgurl=\"\" smallheadimgurl=\"\" username=\"" + id + "\" nickname=\"" + nickname + "\" fullpy=\"\" shortpy=\"\" alias=\"" + alias + "\" imagestatus=\"0\" scene=\"17\" province=\"\" city=\"\" sign=\"\" sex=\"2\" certflag=\"0\" certinfo=\"\" brandIconUrl=\"\" brandHomeUrl=\"\" brandSubscriptConfigUrl=\"\" brandFlags=\"0\" regionCode=\"CN\" />\n"
	return SendTextMsgReq(wxAccount, toUserName, content, nil, contentType)
}

// SendTextMsgToFileHelperRequest 发送消息给文件传输助手
func SendTextMsgToFileHelperRequest(wxAccount wxface.IWXAccount, content string) error {
	_, error := SendTextMsgReq(wxAccount, baseinfo.FileHelperWXID, content, nil, 1)
	return error
}

// 同步消息
func SendNewInitSyncRequest(wxAccount wxface.IWXAccount, sync bool) (*table.SyncMessageResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	mgr := userInfo.SyncKeyMgr()
	if mgr.CurKey() == nil {
		mgr.SetCurKey(&wechat.BufferT{
			ILen:   proto.Uint32(0),
			Buffer: make([]byte, 0),
		})
	}
	if mgr.MaxKey() == nil {
		mgr.SetMaxKey(&wechat.BufferT{
			ILen:   proto.Uint32(0),
			Buffer: make([]byte, 0),
		})
	}
	req := &wechat.NewInitRequest{}
	req.BaseRequest = clientsdk.GetBaseRequest(userInfo)
	//req.BaseRequest.Scene = proto.Uint32(0)
	req.UserName = proto.String(userInfo.GetUserName())
	req.CurrentSynckey = mgr.CurKey()
	req.MaxSynckey = mgr.MaxKey()
	Language := "id" //zh_CN
	if userInfo.DeviceInfo != nil {
		Language = userInfo.DeviceInfo.Language
	}
	req.Language = proto.String(Language)

	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendData := clientsdk.Pack(userInfo, srcData, baseinfo.MMRequestTypeNewInit, 5)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 27,
		CgiUrl: "/cgi-bin/micromsg-bin/newinit",
		Data:   sendData,
	}
	// 发送消息
	response, err := wxlink.WXSend(wxAccount, longReq, sync)
	if err != nil {
		return nil, err
	}
	newSendMsgResponse, _ := response.(*table.SyncMessageResponse)
	return newSendMsgResponse, nil
}

// 同步消息
func NewSyncHistoryMessageRequest(queryKey string, wxAccount wxface.IWXAccount, scene uint32, syncKey string) (*table.SyncMessageResponse, error) {
	userInfo := wxAccount.GetUserInfo()

	if len(userInfo.SyncKey) == 0 || len(userInfo.SyncHistoryKey) == 0 {
		RecoverOnLineStatus(wxAccount)
	}

	sendData := clientsdk.GetNewSyncHistoryMessageReq(userInfo, scene, syncKey)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 121,
		CgiUrl: "/cgi-bin/micromsg-bin/newsync",
		Data:   sendData,
	}

	// 发送消息
	response, err := wxlink.WXSyncSend(wxAccount, longReq)

	if err != nil {
		return nil, err
	}
	newSendMsgResponse := response.(*table.SyncMessageResponse)
	return newSendMsgResponse, nil
}

// SendNewSyncRequest 发送同步请求
func SendNewSyncRequest(wxAccount wxface.IWXAccount, scene uint32, sync bool) (interface{}, error) {
	log.Debug("准备拉取消息")
	// 发送请求
	tmpUserInfo := wxAccount.GetUserInfo()
	reqData := clientsdk.GetNewSyncReq(tmpUserInfo, scene, false)
	//发送给长链接请求去处理
	longReq := &clientsdk.WXLongRequest{
		OpCode: 26,
		Data:   reqData,
		CgiUrl: "/cgi-bin/micromsg-bin/newsync",
	}
	result, err := wxlink.WXSend(wxAccount, longReq, sync)
	if err != nil {
		return nil, err
	}
	return result, nil
}
