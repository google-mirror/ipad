package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXNewSendMsgRouter 心跳包响应路由
type WXNewSendMsgRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (oplog *WXNewSendMsgRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析发送文本消息响应
	var newSendMsgResp wechat.NewSendMsgResponse
	packHeader := wxResp.GetPackHeader()
	err := clientsdk.ParseResponseData(currentUserInfo, packHeader, &newSendMsgResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return newSendMsgResp, nil
}
