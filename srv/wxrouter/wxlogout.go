package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXLogoutRouter 心跳包响应路由
type WXLogoutRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXLogoutRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析退出登陆响应包
	var logoutResp wechat.LogOutResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &logoutResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	currentUserInfo.SetLoginState(baseinfo.MMLoginQrcodeStateNone)
	db.UpdateLoginStatus(currentUserInfo.UUID, int32(currentUserInfo.GetLoginState()), "LogOut")
	// 退出登陆成功
	//currentWXConn.Stop()
	return logoutResp, nil
}
