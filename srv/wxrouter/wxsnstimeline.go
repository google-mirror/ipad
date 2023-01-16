package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXSnsTimeLineRouter 朋友圈路由
type WXSnsTimeLineRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (oplog *WXSnsTimeLineRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析退出登陆响应包
	var timeLineResp wechat.SnsTimeLineResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &timeLineResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return timeLineResp, nil
}
