package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXSnsUserPageRouter 获取好友朋友圈路由
type WXSnsUserPageRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXSnsUserPageRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析 获取好友朋友圈响应
	var userPageResponse wechat.SnsUserPageResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &userPageResponse)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return userPageResponse, nil
}
