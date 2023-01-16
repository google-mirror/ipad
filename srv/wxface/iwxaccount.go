package wxface

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"sync"
)

// IWXAccount 微信账户
type IWXAccount interface {
	GetWxServer() IWXServer
	// GetUserInfo 获取UserInfo
	GetUserInfo() *baseinfo.UserInfo
	SetUserInfo(info *baseinfo.UserInfo)
	// GetUserProfile 获取帐号信息
	GetUserProfile() *wechat.GetProfileResponse
	// SetUserProfile 设置用户配置信息
	SetUserProfile(userProfile *wechat.GetProfileResponse)
	// GetLoginState 获取登录状态
	GetLoginState() uint32
	// SetLoginState 设置登录状态
	SetLoginState(loginState uint32)
	GetWXReqInvoker() IWXReqInvoker
	GetOnceInit() *sync.Once
}
