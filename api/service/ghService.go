package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/utils"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SearchGhService 搜索公众号
func SearchService(queryKey string, username string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendSearchContactRequest(0, 1, 2, username)
		if err != nil {
			return vo.NewFail("SearchContactRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// FollowerService 关注公众号
func FollowerService(queryKey string, ghId string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.VerifyUserRequest(1, "", 0, ghId, "", "")
		if err != nil {
			return vo.NewFail("VerifyUserRequestService err" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// ClickMenuService 操作菜单
func ClickMenuService(queryKey string, m model.ClickCommand) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		result, err := bizcgi.SendClickMenuReq(wxAccount, m.GhUsername, m.MenuId, m.MenuKey)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(result, "")
	})
}

// ReadArticleService 阅读文章
func ReadArticleService(queryKey string, articleUrl string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		resp, err := bizcgi.GetA8KeyRequest(wxAccount, 2, 4, articleUrl, baseinfo.ThrIdGetA8Key)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		if resp.BaseResponse.GetRet() == 0 {
			u, err := url.Parse(resp.GetFullURL())
			if err != nil {
				return vo.NewFail(err.Error())
			}
			m, _ := url.ParseQuery(u.RawQuery)
			v, cookies, err := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
			if err != nil {
				return vo.NewFail(err.Error())
			}
			var malluin, mallkey, wxtokenkey string
			for _, cookie := range cookies {
				if strings.HasPrefix(cookie, "malluin") {
					malluin = getRegxResult("malluin=(.*?);", cookie)
				} else if strings.HasPrefix(cookie, "mallkey") {
					malluin = getRegxResult(`mallkey=(.*?);`, cookie)
				} else if strings.HasPrefix(cookie, "wxtokenkey") {
					malluin = getRegxResult(`wxtokenkey=(.*?);`, cookie)
				}
			}
			passTicket := getRegxResult(`pass_ticket = '(.*?)';`, v)
			//_ = getRegxResult(`qrcheck_ticket = '(.*?)',`, v)
			devicetype := getRegxResult(`var devicetype = "(.+)";`, v)
			appmsg_token := getRegxResult(`window.appmsg_token = "(.+)";`, v)
			appmsg_type := getRegxResult(`var appmsg_type = "(.+)"`, v)
			msg_title := getRegxResult(`var msg_title = '(.+)'`, v)
			//_ = getRegxResult(`var nickname = "(.+)";`, v)
			ct := getRegxResult(`var ct = "(.+)"`, v)
			comment_id := getRegxResult(`var comment_id = "(.+)"`, v)
			msg_daily_idx := getRegxResult(`var msg_daily_idx = "(.+)"`, v)
			req_id := getRegxResult(`var req_id = '(.+)'`, v)
			data := url.Values{}
			data.Add("r", "0."+strconv.FormatInt(time.Now().Unix(), 10))
			data.Add("__biz", m["__biz"][0])
			data.Add("appmsg_type", appmsg_type)
			data.Add("mid", m["mid"][0])
			data.Add("sn", m["sn"][0])
			data.Add("idx", m["idx"][0])
			data.Add("title", msg_title)
			data.Add("ct", ct)
			data.Add("devicetype", devicetype)
			data.Add("version", m["version"][0])
			data.Add("is_need_ticket", "0")
			data.Add("is_need_ad", "0")
			data.Add("is_need_reward", "0")
			data.Add("both_ad", "0")
			data.Add("reward_uin_count", "0")
			data.Add("is_original", "0")
			data.Add("is_only_read", "1")
			data.Add("is_temp_url", "0")
			data.Add("item_show_type", "0")
			data.Add("tmp_version", "1")
			data.Add("more_read_type", "0")
			data.Add("appmsg_like_type", "2")
			data.Add("comment_id", comment_id)
			data.Add("msg_daily_idx", msg_daily_idx)
			data.Add("req_id", req_id)
			data.Add("pass_ticket", passTicket)
			//log.Println(data)
			postUrl := u.Scheme + "://" + u.Host + "/mp/getappmsgext?f=json&mock=&uin=" + malluin + "&key=" + mallkey + "&pass_ticket=" + passTicket + "&wxtoken=" + wxtokenkey + "&devicetype=" + devicetype + "&clientversion=" + m["version"][0] + "&__biz=" + m["__biz"][0] + "&appmsg_token=" + appmsg_token + "&x5=0&f=json"

			rsp := utils.HttpPost(postUrl, data, cookies, wxAccount.GetUserInfo())
			return vo.NewSuccessObj(rsp, "阅读成功")
		}
		return vo.NewSuccessObj(resp, "失败")
	})
}

func getRegxResult(rex string, str string) string {
	var regex = regexp.MustCompile(rex)
	regs := regex.FindStringSubmatch(str)
	if len(regs) > 1 {
		return regs[1]
	}
	return ""
}

// LikeArticleService 点赞文章
func LikeArticleService(queryKey string, articleUrl string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		resp, err := bizcgi.GetA8KeyRequest(wxAccount, 2, 4, articleUrl, baseinfo.ThrIdGetA8Key)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		if resp.BaseResponse.GetRet() == 0 {
			u, err := url.Parse(resp.GetFullURL())
			if err != nil {
				return vo.NewFail(err.Error())
			}
			m, _ := url.ParseQuery(u.RawQuery)
			v, cookies, err := utils.GetHTML(resp.GetFullURL(), resp.HttpHeader, wxAccount.GetUserInfo())
			if err != nil {
				return vo.NewFail(err.Error())
			}
			var malluin, mallkey, wxtokenkey string
			for _, cookie := range cookies {
				if strings.HasPrefix(cookie, "malluin") {
					malluin = getRegxResult("malluin=(.*?);", cookie)
				} else if strings.HasPrefix(cookie, "mallkey") {
					malluin = getRegxResult(`mallkey=(.*?);`, cookie)
				} else if strings.HasPrefix(cookie, "wxtokenkey") {
					malluin = getRegxResult(`wxtokenkey=(.*?);`, cookie)
				}
			}
			passTicket := getRegxResult(`pass_ticket = '(.*?)';`, v)
			//_ = getRegxResult(`qrcheck_ticket = '(.*?)',`, v)
			devicetype := getRegxResult(`var devicetype = "(.+)";`, v)
			appmsg_token := getRegxResult(`window.appmsg_token = "(.+)";`, v)
			//_ = getRegxResult(`var appmsg_type = "(.+)"`, v)
			//_ = getRegxResult(`var msg_title = '(.+)'`, v)
			//_ = getRegxResult(`var nickname = "(.+)";`, v)
			//_ = getRegxResult(`var ct = "(.+)"`, v)
			//_ = getRegxResult(`var comment_id = "(.+)"`, v)
			//_ = getRegxResult(`var msg_daily_idx = "(.+)"`, v)
			//_ = getRegxResult(`var req_id = '(.+)'`, v)
			data := url.Values{}
			data.Add("request_id", strconv.FormatInt(time.Now().Unix(), 10))
			data.Add("client_version", m["version"][0])
			data.Add("devicetype", devicetype)
			data.Add("client_version", m["version"][0])
			data.Add("is_temp_url", "0")
			data.Add("scene", "126")
			data.Add("subscene", "")
			data.Add("appmsg_like_type", "1")
			data.Add("item_show_type", "0")
			data.Add("prompted", "1")
			data.Add("style", "1")
			data.Add("action_type", "1")
			//log.Println(data)
			postUrl := u.Scheme + "://" + u.Host + "/mp/appmsg_like?__biz=" + m["__biz"][0] + "&mid=" + m["mid"][0] + "&idx=" + m["idx"][0] + "&like=1&f=json&appmsgid=" + m["mid"][0] + "&itemidx=" + m["idx"][0] + "&uin=" + malluin + "&key=" + mallkey + "&pass_ticket=" + passTicket + "&wxtoken=" + wxtokenkey + "&devicetype=" + devicetype + "&clientversion=" + m["version"][0] + "&appmsg_token=" + appmsg_token + "&x5=0&f=json"
			rsp := utils.HttpPost(postUrl, data, cookies, wxAccount.GetUserInfo())
			return vo.NewSuccessObj(rsp, "点赞成功")
		}
		return vo.NewSuccessObj(resp, "失败")
	})
}
