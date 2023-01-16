package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXTenPayRouter 支付
type WXTenPayRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (glqr *WXTenPayRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	var response wechat.TenPayResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
