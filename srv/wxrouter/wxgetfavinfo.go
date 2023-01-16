package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXGetFavInfoRouter 心跳包响应路由
type WXGetFavInfoRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXGetFavInfoRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析退出登陆响应包
	var favInfo wechat.GetFavInfoResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &favInfo)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return favInfo, nil
}
