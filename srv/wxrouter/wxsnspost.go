package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXSnsPostRouter 发送朋友圈响应路由
type WXSnsPostRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXSnsPostRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析发送朋友圈响应包
	var snsPostResp wechat.SnsPostResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &snsPostResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	return snsPostResp, nil
}
