package wxlink

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/proxynet"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/websrv"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"strconv"
	"sync"
)

// WXAccount 代表微信帐号
type WXAccount struct {
	userInfo *baseinfo.UserInfo
	// 配置文件信息
	userProfile *wechat.GetProfileResponse
	// 登录状态
	loginState uint32
	iwxServer  *wxface.IWXServer
	// 请求调用器
	wxReqInvoker wxface.IWXReqInvoker
	// 首次登录初始化只执行一次
	onceInit sync.Once
}

// NewWXAccount 生成一个新的账户
func NewWXAccount(taskInfo *websrv.TaskInfo, proxyInfo *proxynet.WXProxyInfo, iwxServer *wxface.IWXServer, userInfo *baseinfo.UserInfo) *WXAccount {
	if userInfo == nil {
		userInfo = clientsdk.NewUserInfo(taskInfo.UUID, taskInfo.DeviceId, proxyInfo)
	}
	wxAccount := &WXAccount{
		userInfo:    userInfo,
		userProfile: nil,
		iwxServer:   iwxServer,
		loginState:  baseinfo.MMLoginStateNoLogin,
	}
	wxmgr.WxAccountMgr.Add(taskInfo.UUID, wxAccount)
	wxAccount.wxReqInvoker = NewWXLongReqInvoker(wxAccount)
	return wxAccount
}

func (wxAccount *WXAccount) GetWxServer() wxface.IWXServer {
	return *wxAccount.iwxServer
}

//func (wxAccount *WXAccount) GetWxConnect() wxface.IWXConnect {
//	connect := wxmgr.WxConnectMgr.GetWXAccountByUserInfoUUID(wxAccount.userInfo.UUID)
//	return connect
//}

// GetUserInfo 获取UserInfo
func (wxAccount *WXAccount) GetUserInfo() *baseinfo.UserInfo {
	return wxAccount.userInfo
}

// SetUserInfo 设置用户信息
func (wxAccount *WXAccount) SetUserInfo(info *baseinfo.UserInfo) {
	wxAccount.userInfo = info
}

// SetUserProfile 设置用户配置信息
func (wxAccount *WXAccount) SetUserProfile(userProfile *wechat.GetProfileResponse) {
	wxAccount.userProfile = userProfile
}

// GetUserProfile 获取帐号信息
func (wxAccount *WXAccount) GetUserProfile() *wechat.GetProfileResponse {
	return wxAccount.userProfile
}

// GetLoginState 获取登录状态
func (wxAccount *WXAccount) GetLoginState() uint32 {
	//return wxAccount.loginState
	return wxAccount.GetUserInfo().GetLoginState()
}

// SetLoginState 设置登录状态
func (wxAccount *WXAccount) SetLoginState(loginState uint32) {
	wxAccount.loginState = loginState
	wxAccount.userInfo.SetLoginState(loginState)
}

func (wxAccount *WXAccount) GetWXReqInvoker() wxface.IWXReqInvoker {
	return wxAccount.wxReqInvoker
}

// GetBindQueryNewReq GetBindQueryNewReq
func (wxAccount *WXAccount) GetBindQueryNewReq() (string, error) {
	tmpBindQueryNewReq := &baseinfo.BindQueryNewReq{}
	tmpBindQueryNewReq.BalanceVersion = wxAccount.userInfo.BalanceVersion
	tmpBindQueryNewReq.BindQueryScen = 1
	tmpBindQueryNewReq.BindTypeCond = "all_type"
	tmpBindQueryNewReq.ISRoot = 0
	tmpBindQueryNewReq.City = wxAccount.userProfile.GetUserInfo().GetCity()
	tmpBindQueryNewReq.Province = wxAccount.userProfile.GetUserInfo().GetProvince()
	tmpString := "balance_version=" + strconv.Itoa(int(tmpBindQueryNewReq.BalanceVersion))
	tmpString = tmpString + "&bind_query_scene=" + strconv.Itoa(int(tmpBindQueryNewReq.BindQueryScen))
	tmpString = tmpString + "&bind_type_cond=" + tmpBindQueryNewReq.BindTypeCond
	tmpString = tmpString + "&city=" + tmpBindQueryNewReq.City
	tmpString = tmpString + "&is_device_open_touch=" + strconv.Itoa(int(tmpBindQueryNewReq.ISDeviceOpenTouch))
	tmpString = tmpString + "&is_root=" + strconv.Itoa(int(tmpBindQueryNewReq.ISRoot))
	tmpString = tmpString + "&province=" + tmpBindQueryNewReq.Province
	wcPaySign, err := clientsdk.TenPaySignDes3(tmpString, "%^&*Tenpay!@#$")
	if err != nil {
		return "", err
	}
	tmpString = tmpString + "&WCPaySign=" + wcPaySign
	return tmpString, nil
}

// GetOnceInit 获取OnceInit
func (wxAccount *WXAccount) GetOnceInit() *sync.Once {
	return &wxAccount.onceInit
}
