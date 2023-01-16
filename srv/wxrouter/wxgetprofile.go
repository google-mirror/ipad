package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXGetProfileRouter 获取帐号信息路由
type WXGetProfileRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXGetProfileRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析账号Profile响应
	userProfileResp := new(wechat.GetProfileResponse)
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), userProfileResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 设置昵称和头像
	currentUserInfo.NickName = userProfileResp.GetUserInfo().GetNickName().GetStr()
	currentUserInfo.HeadURL = userProfileResp.GetUserInfoExt().GetSmallHeadImgUrl()
	// 设置UserProfile信息
	currentWXAccount.SetUserProfile(userProfileResp)
	return userProfileResp, nil
}
