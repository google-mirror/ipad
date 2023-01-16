package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXGetA8KeyRouter getA8key路由
type WXGetA8KeyRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (glqr *WXGetA8KeyRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentWXUserInfo := currentWXAccount.GetUserInfo()

	// 解析获取联系人响应
	getA8KeyResp := new(wechat.GetA8KeyResp)
	err := clientsdk.ParseResponseData(currentWXUserInfo, wxResp.GetPackHeader(), getA8KeyResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return *getA8KeyResp, nil
}
