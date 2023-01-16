package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXInitContactRouter 初始化联系人响应路由
type WXInitContactRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXInitContactRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentWXCache := currentWXConn.GetWXCache()

	// 解析初始化通讯录响应
	var initContactResp wechat.InitContactResp
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &initContactResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 设置ContactSeq
	//newContactSeq := initContactResp.GetCurrentWxcontactSeq()
	//currentWXCache.SetContactSeq(newContactSeq)
	return initContactResp, nil
}
