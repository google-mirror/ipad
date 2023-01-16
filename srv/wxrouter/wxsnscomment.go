package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXSnsCommentRouter 发送朋友圈响应路由
type WXSnsCommentRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXSnsCommentRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析发送朋友圈响应包
	var snsCommentResp wechat.SnsCommentResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &snsCommentResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return snsCommentResp, nil
}
