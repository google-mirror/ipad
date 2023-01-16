package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXGetLoginQrcodeRouter 获取二维码响应路由
type WXGetLoginQrcodeRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (glqr *WXGetLoginQrcodeRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentUserInfo := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey()).GetUserInfo()

	// 获取登录二维码响应
	var qrCodeResponse wechat.LoginQRCodeResponse
	packHeader := wxResp.GetPackHeader()
	err := clientsdk.ParseResponseData(currentUserInfo, packHeader, &qrCodeResponse)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return qrCodeResponse, nil
}
