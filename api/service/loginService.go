package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/utils"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/proxynet"
	"feiyu.com/wx/db"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"feiyu.com/wx/srv/wxmgr"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 提取62
func Get62DataService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		resp, err := utils.GenerateWxDat(wxAccount.GetUserInfo().DeviceInfo.Imei)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// GetLoginQrCodeNewService 自动生成UUID 可使用代理
func GetLoginQrCodeNewService(queryKey string, model model.GetLoginQrCodeModel) vo.DTO {
	return checkExIdPerform(queryKey, model.DeviceId, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//queryKey 为空字符串
		if wxAccount == nil {
			return vo.NewFail("逻辑错误！")
		}
		// 使用使用代理
		if len(model.Proxy) > 0 {
			// 设置代理
			proxyInfo := proxynet.ParseWXProxyInfo(model.Proxy)
			wxAccount.GetUserInfo().SetProxy(proxyInfo)
		}

		//
		//检查是否在线会进行一次二次登录如果用户退出或者掉线才能重新获取二维码
		if bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.DTO{
				Code: vo.FAIL_Bound,
				Data: nil,
				Text: "该链接已绑定微信号！",
			}
		}
		iwxConnect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		if iwxConnect == nil {
			iwxConnect = wxlink.NewWXConnect(queryKey)
		}
		if !iwxConnect.IsConnected() {
			err := iwxConnect.Start()
			if err != nil {
				return vo.NewFail(err.Error())
			}
		}
		//发送获取二维码请求
		loginQrCodeResp, err := bizcgi.SendGetLoginQrcodeRequest(wxAccount)
		if err != nil {
			return vo.NewFail("获取二维码失败！err:" + err.Error())
		}
		wxAccount.GetUserInfo().QrUuid = loginQrCodeResp.GetUuid()
		wxcore.WxInfoCache.SetQrcodeInfo(loginQrCodeResp.GetUuid(), loginQrCodeResp.GetAes().GetKey())
		resp := gin.H{
			"baseResp":  loginQrCodeResp.GetBaseResponse(),
			"Txt":       "建议返回data=之后内容自定义生成二维码",
			"QrUuid":    loginQrCodeResp.GetUuid(),
			"QrCodeUrl": "http://api.qrserver.com/v1/create-qr-code/?data=http://weixin.qq.com/x/" + loginQrCodeResp.GetUuid(),
		}
		//是否为新链接
		if newIWXConnect {
			resp["Key"] = wxAccount.GetUserInfo().UUID
		}

		return vo.NewSuccess(resp, "")

	})
}

// 发送验证码
func WxBindOpMobileForRegService(queryKey string, m model.WxBindOpMobileForModel) vo.DTO {
	return checkExIdPerform(queryKey, "", func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//queryKey 为空字符串
		if wxAccount == nil {
			return vo.NewFail("逻辑错误！")
		}
		// 使用使用代理
		if len(m.Proxy) > 0 {
			// 设置代理
			proxyInfo := proxynet.ParseWXProxyInfo(m.Proxy)
			wxAccount.GetUserInfo().SetProxy(proxyInfo)
		}
		//iwxConnect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		//if !iwxConnect.IsConnected() {
		//	err := iwxConnect.Start()
		//	if err != nil {
		//		return vo.NewFail(err.Error())
		//	}
		//
		//}
		//请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendWxBindOpMobileForRequest(m.OpCode, m.PhoneNumber, m.VerifyCode)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, wxAccount.GetUserInfo().UUID)
	})
}

// GetLoginQrCodeService获取登录二维码
func GetLoginQrCodeService(queryKey string) vo.DTO {
	return checkExIdPerform(queryKey, "", func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//queryKey 为空字符串
		if wxAccount == nil {
			return vo.NewFail("逻辑错误！")
		}

		//
		//检查是否在线会进行一次二次登录如果用户退出或者掉线才能重新获取二维码
		if bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.DTO{
				Code: vo.FAIL_Bound,
				Data: nil,
				Text: "该链接已绑定微信号！",
			}
		}

		//iwxConnect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		////iwxConnect.Stop()
		//if !iwxConnect.IsConnected() {
		//	err := iwxConnect.Start()
		//	if err != nil {
		//		return vo.NewFail(err.Error())
		//	}
		//
		//}
		wxAccount.GetUserInfo().Session = []byte{}
		wxAccount.GetUserInfo().Uin = 0
		wxAccount.GetUserInfo().CheckSumKey = []byte{}

		//发送获取二维码请求
		loginQrCodeResp, err := bizcgi.SendGetLoginQrcodeRequest(wxAccount)
		if err != nil {
			return vo.NewFail("获取二维码失败！err:" + err.Error())
		}
		wxAccount.GetUserInfo().QrUuid = loginQrCodeResp.GetUuid()
		wxcore.WxInfoCache.SetQrcodeInfo(loginQrCodeResp.GetUuid(), loginQrCodeResp.GetAes().GetKey())

		resp := gin.H{
			"baseResp":  loginQrCodeResp.GetBaseResponse(),
			"Txt":       "建议返回data=之后内容自定义生成二维码",
			"QrUuid":    loginQrCodeResp.GetUuid(),
			"QrCodeUrl": "http://api.qrserver.com/v1/create-qr-code/?data=http://weixin.qq.com/x/" + loginQrCodeResp.GetUuid(),
		}
		//是否为新链接
		if newIWXConnect {
			resp["Key"] = wxAccount.GetUserInfo().UUID
		}

		return vo.NewSuccess(resp, "")

	})
}

// CheckLoginQrCodeStatusService 检测扫码状态
func CheckLoginQrCodeStatusService(queryKey string, uuid string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		result, err := bizcgi.SendCheckLoginQrcodeRequest(wxAccount, uuid)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(result, "")
	})
}

// 初始化状态
func GetInItStatusService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		// 获取二级缓存器
		//iwxCache := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID).GetWXCache()
		msg := "初始化未完成!"
		//if iwxCache.IsInitFinished() {
		//	msg = "初始化完成!"
		//}
		return vo.NewSuccessObj(nil, msg)
	})
}

// WakeUpLoginService 发送唤醒登录
func WakeUpLoginService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//判断在线情况
		if bizcgi.CheckOnLineStatus(wxAccount) {
			_ = bizcgi.SendLogoutRequest(wxAccount)
			//time.Sleep(time.Second * 2)
		}
		//connect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		//connect.GetWXAccount().GetUserInfo().Session = []byte{}
		wxAccount.GetUserInfo().Uin = 0
		wxAccount.GetUserInfo().CheckSumKey = []byte{}

		//发送唤醒登录
		pushResp, err := bizcgi.SendPushQrLoginNotice(wxAccount)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		if pushResp.GetBaseResponse().GetRet() == 0 {
			wxAccount.GetUserInfo().QrUuid = pushResp.GetUUID()
			wxcore.WxInfoCache.SetQrcodeInfo(pushResp.GetUUID(), pushResp.GetNotifyKey().GetBuffer())
		} else {
			log.Warn("唤醒失败关闭长连接")
			//connect.Stop()
			return vo.NewSuccessObj(pushResp, "发送唤醒登录失败！")
		}
		return vo.NewSuccessObj(pushResp, "发送唤醒登录成功！")
	})
}

// GetLoginStatusService 获取登录状态
func GetLoginStatusService(queryKey string, loginJournal bool, autoLogin bool) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		var errMsg string
		log.Info("检查登录状态")
		//取登录状态
		loginState := wxAccount.GetLoginState()
		//判断在线情况 掉线或者离线会通过Token 登录确认掉线状态
		if wxAccount != nil {
			if bizcgi.CheckOnLineStatus(wxAccount) {
				errMsg = "账号在线状态良好！"
			} else {
				if autoLogin {
					if bizcgi.RecoverOnLineStatus(wxAccount) {
						errMsg = "账号在线状态良好！"
					} else {
						errMsg = "账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState()))
					}
				} else if loginState == baseinfo.MMLoginStateNoLogin {
					errMsg = "该账号需要重新登录！loginState == MMLoginStateNoLogin "
				} else {
					errMsg = "账号已离线！"
					//errMsg = db.GetUSerLoginErrMsg(wxAccount.GetUserInfo().GetUserName())
				}
			}
		} else {
			errMsg = "未获取到长链接，状态未知！"
		}

		logs := make([]table.UserLoginLog, 0)
		//登录日志
		if loginJournal {
			//Mysql 从数据库查询记录
			logs = db.GetLoginJournal(wxAccount.GetUserInfo().GetUserName())
		}
		targetIp := ""
		userInfoEntity := db.GetUserInfoEntity(queryKey)
		if userInfoEntity != nil {
			targetIp = userInfoEntity.TargetIp
		}
		return vo.NewSuccess(gin.H{
			"loginState":  int32(wxAccount.GetLoginState()),
			"targetIp":    targetIp,
			"loginErrMsg": errMsg,
			"loginJournal": gin.H{
				"count": len(logs),
				"logs":  logs,
			},
		}, "")
	})
}

// DeviceIdLoginService 账号密码登录
func DeviceIdLoginService(queryKey string, m model.DeviceIdLoginModel) vo.DTO {
	return checkExIdPerformIp(queryKey, m.Proxy, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		if m.DeviceId == "" {
			key := fmt.Sprintf("%s%s", "wechat:device62:", m.UserName)
			exists, _ := db.Exists(key)
			if exists {
				db.GETObj(key, &m)
			} else {
				//生成62数据
				m.DeviceId, _ = utils.GenerateWxDat(baseutils.RandomSmallHexString(32))
				db.SETExpirationObj(key, m, 60*60*24*1)
			}
		}
		if m.DeviceId == "" || m.UserName == "" || m.Password == "" {
			return vo.NewFail("登录数据错误！")
		}
		//queryKey 为空字符串
		//if wxAccount == nil {
		//	wxAccount = CreateWXConnectByQueryKey(guuid.New().String(), m.Proxy, nil)
		//}
		//取用户信息
		//iwxConnect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		userInfo := wxAccount.GetUserInfo()
		//检查是否在线会进行一次二次登录如果用户退出或者掉线才能重新获取二维码
		if bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.DTO{
				Code: vo.FAIL_Bound,
				Data: nil,
				Text: "该链接已绑定微信号！",
			}
		}

		//if !iwxConnect.IsConnected() {
		//	err := iwxConnect.Start()
		//	if err != nil {
		//		return vo.NewFail(err.Error())
		//	}
		//
		//}
		userInfo.LoginDataInfo = baseinfo.LoginDataInfo{
			UserName:  m.UserName,
			PassWord:  m.Password,
			LoginData: m.DeviceId,
			Language:  m.DeviceInfo.Language,
		}
		//判断如果是A16数据
		if strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") {
			userInfo.DeviceInfoA16.Imei = m.DeviceInfo.ImeI
			userInfo.DeviceInfoA16.AndriodId = m.DeviceInfo.AndroidId
			userInfo.DeviceInfoA16.PhoneModel = m.DeviceInfo.Model
			userInfo.DeviceInfoA16.Manufacturer = m.DeviceInfo.Manufacturer
		}
		//判断语言 62加进去
		if userInfo.LoginDataInfo.Language != "" {
			userInfo.DeviceInfo.Language = userInfo.LoginDataInfo.Language
			userInfo.DeviceInfo.RealCountry = userInfo.LoginDataInfo.Language
		}
		checkDeviceToken(userInfo)
		//err := reqInvoker.SendManualAuthByDeviceIdRequest()
		_, err := bizcgi.SendManualAuthByDeviceIdRequest(wxAccount)
		if err != nil {
			return vo.NewFail("发送登录请求失败！err - " + err.Error())
		}
		return vo.NewSuccess(gin.H{"uuid": wxAccount.GetUserInfo().UUID, "data": m.DeviceId}, "发送登录请求成功!")
	})
}

// 短信登录
func SmsLoginService(queryKey string, m model.DeviceIdLoginModel) vo.DTO {
	return checkExIdPerform(queryKey, m.DeviceId, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		if m.UserName == "" || m.Password == "" {
			return vo.NewFail("登录数据错误！")
		}
		//queryKey 为空字符串
		//if iwxConnect == nil {
		//	iwxConnect = CreateWXConnectByQueryKey(guuid.New().String(), m.Proxy, nil)
		//}
		//检查是否在线会进行一次二次登录如果用户退出或者掉线才能重新获取二维码
		if bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.DTO{
				Code: vo.FAIL_Bound,
				Data: nil,
				Text: "该链接已绑定微信号！",
			}
		}
		//iwxConnect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		//if !iwxConnect.IsConnected() {
		//	err := iwxConnect.Start()
		//	if err != nil {
		//		return vo.NewFail(err.Error())
		//	}
		//
		//}
		userInfo := wxAccount.GetUserInfo()
		userInfo.LoginDataInfo = baseinfo.LoginDataInfo{
			UserName: m.UserName,
			PassWord: m.Password,
			Ticket:   m.Ticket,
		}
		_, err := bizcgi.SendManualAuthByDeviceIdRequest(wxAccount)
		if err != nil {
			return vo.NewFail("发送登录请求失败！err - " + err.Error())
		}

		return vo.NewSuccess(gin.H{"uuid": wxAccount.GetUserInfo().UUID}, "发送登录请求成功!")
	})
}

// a16登录
func A16LoginService(queryKey string, m model.DeviceIdLoginModel) vo.DTO {
	return checkExIdPerform(queryKey, m.DeviceId, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		if m.DeviceId == "" || m.UserName == "" || m.Password == "" {
			return vo.NewFail("登录数据错误！")
		}
		//queryKey 为空字符串
		//if iwxConnect == nil {
		//	iwxConnect = CreateWXConnectByQueryKey(guuid.New().String(), m.Proxy, nil)
		//}
		//检查是否在线会进行一次二次登录如果用户退出或者掉线才能重新获取二维码
		if bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.DTO{
				Code: vo.FAIL_Bound,
				Data: nil,
				Text: "该链接已绑定微信号！",
			}
		}
		//iwxConnect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		//if !iwxConnect.IsConnected() {
		//	err := iwxConnect.Start()
		//	if err != nil {
		//		return vo.NewFail(err.Error())
		//	}
		//
		//}
		userInfo := wxAccount.GetUserInfo()
		userInfo.LoginDataInfo = baseinfo.LoginDataInfo{
			UserName:  m.UserName,
			PassWord:  m.Password,
			LoginData: m.DeviceId,
			NewType:   m.Type,
		}
		//判断如果是A16数据
		if strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") {
			userInfo.DeviceInfoA16.Imei = m.DeviceInfo.ImeI
			userInfo.DeviceInfoA16.AndriodId = m.DeviceInfo.AndroidId
			userInfo.DeviceInfoA16.PhoneModel = m.DeviceInfo.Model
			userInfo.DeviceInfoA16.Manufacturer = m.DeviceInfo.Manufacturer
		}
		checkDeviceToken(userInfo)
		_, err := bizcgi.SendManualAuthByDeviceIdRequest(wxAccount)
		if err != nil {
			return vo.NewFail("发送登录请求失败！err - " + err.Error())
		}

		return vo.NewSuccess(gin.H{"uuid": wxAccount.GetUserInfo().UUID}, "发送登录请求成功!")
	})
}

// 铺助新手机登录
func PhoneDeviceLoginService(queryKey string, m model.PhoneLoginModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		resp, err := bizcgi.GetA8KeyRequest(wxAccount, 2, 4, m.Url, baseinfo.ThrIdGetA8Key)
		//flag:=strings.Contains(resp.GetFullURL(), "https://login.weixin.qq.com")
		if err != nil {
			return vo.NewFail(err.Error())
		}
		if resp.BaseResponse.GetRet() == 0 {
			urlv := resp.GetFullURL()
			v, _, _ := utils.GetHTML(urlv, resp.HttpHeader, wxAccount.GetUserInfo())
			//log.Println(v)
			var action = regexp.MustCompile(`action\="([^"]*)"`)
			if len(action.FindStringSubmatch(v)) <= 0 {
				return vo.NewSuccessObj(nil, "二维码已失效！")
			}
			postUrl := "https://login.weixin.qq.com" + action.FindStringSubmatch(v)[1]
			rsp, ck, _ := utils.GetHTML(postUrl, resp.HttpHeader, wxAccount.GetUserInfo())
			//log.Println(rsp)
			apiUrl := "https://login.weixin.qq.com" + action.FindStringSubmatch(rsp)[1]
			log.Println(apiUrl)
			data := url.Values{}
			time.Sleep(1 * time.Second)
			rspv := utils.HttpPost(postUrl, data, ck, wxAccount.GetUserInfo())
			return vo.NewSuccessObj(rspv, "")
		}
		return vo.NewSuccessObj(nil, "失败!")
	})
}

// 获取设备
func GetSafetyInfoService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGetSafetyInfoRequest()
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 删除设备
func DelSafeDeviceService(queryKey string, m model.DelSafeDeviceModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendDelSafeDeviceRequest(m.DeviceUUID)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 检测微信登录环境
func CheckCanSetAliasService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendCheckCanSetAliasRequest()
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 扫码登录新设备
func ExtDeviceLoginConfirmGetService(queryKey string, url string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendExtDeviceLoginConfirmGetRequest(url)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

func IWXConnectMgrService() vo.DTO {
	//获取实例管理器
	return vo.NewSuccessObj(wxmgr.WxConnectMgr.ShowConnectInfo(), "")
}
