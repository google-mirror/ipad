package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXCheckFavCdnRouter token登陆响应路由
type WXCheckFavCdnRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (glqr *WXCheckFavCdnRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentInvoker := currentWXConn.GetWXReqInvoker()

	// 解析token登陆响应
	var checkFavCdnResponse wechat.CheckCDNResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &checkFavCdnResponse)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return checkFavCdnResponse, nil
}
