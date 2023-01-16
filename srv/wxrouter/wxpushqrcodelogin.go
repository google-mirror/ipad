package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXPushQrCodeLoginRouter 二维码二次登录
type WXPushQrCodeLoginRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXPushQrCodeLoginRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析获取二维码响应包
	var getQrcodeResp wechat.PushLoginURLResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &getQrcodeResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	//第一次push 成功后没有点击确认 第二次push 会报 -2017 需要重新扫码登录
	if getQrcodeResp.GetBaseResponse().GetRet() == -2017 {
		// 更改登录状态为未登录状态
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateNoLogin)
		// Mysql 更新状态
		db.UpdateLoginStatus(currentUserInfo.UUID, int32(baseinfo.MMLoginQrcodeStateNone), "需要重新登录，pushLoginResp.RetCode -2017")
		return nil, nil
	}

	// 如果获取失败
	if getQrcodeResp.GetBaseResponse().GetRet() != 0 {
		return nil, nil
	}

	return getQrcodeResp, nil
}
