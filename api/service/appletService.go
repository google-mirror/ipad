package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/utils"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/lunny/log"
	"net/url"
	"regexp"
	"strconv"
)

// GetA8KeyService 授权链接
func GetA8KeyService(queryKey string, m model.GetA8KeyRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		resp, err := bizcgi.GetA8KeyRequest(wxAccount, m.OpCode, m.Scene, m.ReqUrl, baseinfo.ThrIdGetA8Key)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// JSLoginService 小程序授权
func JsLoginService(queryKey string, m model.AppletModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		resp, err := reqInvoker.JSLoginRequest(m.AppId)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// JSOperateWxDataService
func JSOperateWxDataService(queryKey string, m model.AppletModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		resp, err := reqInvoker.JSOperateWxDataRequest(m.AppId)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

func SdkOauthAuthorizeService(queryKey string, m model.AppletModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		resp, err := reqInvoker.SdkOauthAuthorizeRequest(m.AppId, m.SdkName, m.PackageName)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 二维码授权登录
func QRConnectAuthorizeService(queryKey string, m model.QRConnectAuthorizeModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		resp, err := reqInvoker.SendQRConnectAuthorize(m.QrUrl)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 二维码授权登录确认
func QRConnectAuthorizeConfirmService(queryKey string, m model.QRConnectAuthorizeModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		resp, err := reqInvoker.SendQRConnectAuthorizeConfirm(m.QrUrl)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 授权链接
func GetMpA8Service(queryKey string, m model.GetMpA8KeyModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		resp, err := reqInvoker.SendGetMpA8Request(m.Url, m.Opcode)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 授权公众号登录
func AuthMpLoginService(queryKey string, m model.GetMpA8KeyModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		m.Opcode = 2
		resp, err := reqInvoker.SendGetMpA8Request(m.Url, m.Opcode)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		log.Println(resp)
		if resp.BaseResponse.GetRet() == 0 {
			if m.Scene == 0 {
				v, cookie, _ := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
				//var passTicket = regexp.MustCompile(`var _pass_ticket = '([^']*)'`)
				unescape, _ := url.QueryUnescape(resp.GetFullURL())
				var passTicket = regexp.MustCompile(`&pass_ticket=([^']*)&`)
				log.Println(v)
				var appmsgToken = regexp.MustCompile(`window.appmsg_token = "([^']*)"`)
				//var qrticket = regexp.MustCompile(`qrticket: '([^']*)'`)
				data := url.Values{}
				/*data.Add("uin", "777")
				data.Add("key", "777")
				data.Add("pass_ticket", passTicket.FindStringSubmatch(unescape)[1])
				data.Add("appmsg_token", appmsgToken.FindStringSubmatch(v)[1])
				data.Add("f", "json")
				data.Add("param", "qrticket")
				data.Add("qrticket", qrticket.FindStringSubmatch(v)[1])*/
				urlresp := "https://mp.weixin.qq.com/mp/scanlogin?action=confirm&uin=777&key=777&pass_ticket=" + passTicket.FindStringSubmatch(unescape)[1] + "&wxtoken=&appmsg_token=" + appmsgToken.FindStringSubmatch(v)[1] + "&x5=0&f=json"
				rsp := utils.HttpPost(urlresp, data, cookie, wxAccount.GetUserInfo())
				//rsp := utils.HttpPost("https://mp.weixin.qq.com/wap/loginauthqrcode?action=confirm", data, cookie, wxAccount.GetUserInfo())
				return vo.NewSuccessObj(rsp, "授权登录成功")
			} else if m.Scene == 1 {
				v, cookie, _ := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
				var ticket = regexp.MustCompile(`_ticket = "(.*?)",`)
				var uuid = regexp.MustCompile(`_uuid = "(.*?)",`)
				var passTicket = regexp.MustCompile(`_pass_ticket = "(.*?)",`)
				var appmsgToken = regexp.MustCompile(`_appmsg_token = "(.*?)",`)
				var msgId = regexp.MustCompile(`_msgid = "(.*?)",`)
				var secondOpenId = regexp.MustCompile(`_second_openid = "(.*?)",`)
				data := url.Values{}
				data.Add("ticket", ticket.FindStringSubmatch(v)[1])
				data.Add("uuid", uuid.FindStringSubmatch(v)[1])
				data.Add("action", "check")
				data.Add("uin", "777")
				data.Add("key", "777")
				data.Add("pass_ticket", passTicket.FindStringSubmatch(v)[1])
				data.Add("appmsg_token", appmsgToken.FindStringSubmatch(v)[1])
				data.Add("code", "invalid")
				data.Add("type", "bind_second_admin")
				data.Add("msgid", msgId.FindStringSubmatch(v)[1])
				data.Add("second_openid", secondOpenId.FindStringSubmatch(v)[1])
				data.Add("expire_time_type", "0")
				data.Add("allow", "1")
				rsp := utils.HttpPost("https://mp.weixin.qq.com/safe/safeconfirm_reply", data, cookie, wxAccount.GetUserInfo())
				return vo.NewSuccessObj(rsp, "授权验证成功")
			} else if m.Scene == 2 {
				v, cookie, _ := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
				var ticket = regexp.MustCompile(`_ticket = "(.*?)",`)
				var uuid = regexp.MustCompile(`_uuid = "(.*?)",`)
				var passTicket = regexp.MustCompile(`_pass_ticket = "(.*?)",`)
				var appmsgToken = regexp.MustCompile(`_appmsg_token = "(.*?)",`)
				var msgId = regexp.MustCompile(`_msgid = "(.*?)",`)
				data := url.Values{}
				data.Add("ticket", ticket.FindStringSubmatch(v)[1])
				data.Add("uuid", uuid.FindStringSubmatch(v)[1])
				data.Add("action", "check")
				data.Add("uin", "777")
				data.Add("key", "777")
				data.Add("pass_ticket", passTicket.FindStringSubmatch(v)[1])
				data.Add("appmsg_token", appmsgToken.FindStringSubmatch(v)[1])
				data.Add("code", "invalid")
				data.Add("type", "appkey")
				data.Add("msgid", msgId.FindStringSubmatch(v)[1])
				data.Add("allow", "1")
				rsp := utils.HttpPost("https://mp.weixin.qq.com/safe/safeconfirm_reply", data, cookie, wxAccount.GetUserInfo())
				return vo.NewSuccessObj(rsp, "授权验证成功")
			} else if m.Scene == 3 {
				v, cookie, _ := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
				var passTicket = regexp.MustCompile(`pass_ticket: '(.*?)',`)
				var qrcheckTicket = regexp.MustCompile(`qrcheck_ticket: '(.*?)',`)
				var appmsgToken = regexp.MustCompile(`window.appmsg_token = "(.*?)";`)
				data := url.Values{}
				data.Add("f", "json")
				data.Add("action", "scan")
				data.Add("qrcheck_ticket", qrcheckTicket.FindStringSubmatch(v)[1])
				data.Add("uin", "777")
				data.Add("key", "777")
				data.Add("pass_ticket", passTicket.FindStringSubmatch(v)[1])
				data.Add("appmsg_token", appmsgToken.FindStringSubmatch(v)[1])
				rsp := utils.HttpPost("https://mp.weixin.qq.com/wap/qrcheckoper", data, cookie, wxAccount.GetUserInfo())
				return vo.NewSuccessObj(rsp, "授权验证成功")
			} else if m.Scene == 4 {
				v, cookie, _ := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
				//log.Println(v)
				var passTicket = regexp.MustCompile(`pass_ticket: '(.*?)',`)
				var appmsgToken = regexp.MustCompile(`appmsg_token: '(.*?)',`)
				var qrcheckTicket = regexp.MustCompile(`qrcheck_ticket: '(.*?)',`)
				data := url.Values{}
				data.Add("f", "json")
				data.Add("action", "confirm")
				data.Add("operate_type", "1")
				data.Add("qrcheck_ticket", qrcheckTicket.FindStringSubmatch(v)[1])
				data.Add("uin", "777")
				data.Add("key", "777")
				data.Add("pass_ticket", passTicket.FindStringSubmatch(v)[1])
				data.Add("appmsg_token", appmsgToken.FindStringSubmatch(v)[1])
				//log.Println(data)
				//log.Println(cookie)
				rsp := utils.HttpPost("https://mp.weixin.qq.com/mp/wapsafeqrcode", data, cookie, wxAccount.GetUserInfo())
				return vo.NewSuccessObj(rsp, "授权验证成功")
			}
		}
		return vo.NewSuccessObj(resp, "失败")
	})
}
