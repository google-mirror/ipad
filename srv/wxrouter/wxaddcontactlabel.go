package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXAddContactLabelRouter 心跳包响应路由
type WXAddContactLabelRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXAddContactLabelRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentReqInvoker := currentWXConn.GetWXReqInvoker()

	// 解析退出登陆响应包
	var addLabelResp wechat.AddContactLabelResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &addLabelResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 跟新标签列表
	//_, _ = currentReqInvoker.SendGetContactLabelListRequest(false)
	return addLabelResp, nil
}
