package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"strconv"
)

// FavSyncService 同步收藏
func FavSyncService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 同步收藏
		resp, err := reqInvoker.SendFavSyncRequest()
		if err != nil {
			return vo.NewFail("FavSyncService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 获取收藏list
func GetFavListService(queryKey string, req model.FavInfoModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 同步收藏
		resp, err := reqInvoker.SendFavSyncListRequestResult(req.KeyBuf)
		if err != nil {
			return vo.NewFail("FavSyncService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// GetFavInfoService 获取收藏信息
func BatchDelFavItemService(queryKey string, m model.FavInfoModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 获取收藏消息
		resp, err := reqInvoker.SendBatchDelFavItemRequestResult(m.FavId)
		if err != nil {
			return vo.NewFail("BatchDelFavItemService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// GetFavInfoService 获取收藏信息
func GetFavInfoService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 获取收藏消息
		resp, err := reqInvoker.SendGetFavInfoRequestResult()
		if err != nil {
			return vo.NewFail("GetFavInfoService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// GetFavInfoService 获取收藏详细
func BatchGetFavItemService(queryKey string, m model.FavInfoModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 获取收藏消息
		resp, err := bizcgi.SendBatchGetFavItemReq(wxAccount, m.FavId)
		if err != nil {
			return vo.NewFail("GetFavInfoService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// ShareFavService 分享收藏
func ShareFavService(queryKey string, m model.ShareFavModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 分享收藏
		resp, err := bizcgi.SendShareFavRequest(wxAccount, m.FavId, m.ToUserName)
		if err != nil {
			return vo.NewFail("GetFavInfoService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// CheckFavCdnService 检测收藏cdn
func CheckFavCdnService(queryKey string, m model.CheckFavCdnModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		checkFavCdnItem := &baseinfo.CheckFavCdnItem{
			DataId:         m.DataId,
			DataSourceId:   m.DataSourceId,
			DataSourceType: m.DataSourceType,
			FullMd5:        m.FullMd5,
			FullSize:       m.FullSize,
			Head256Md5:     m.Head256Md5,
			IsThumb:        m.IsThumb,
		}
		// 分享收藏
		resp, err := bizcgi.SendCheckFavCdnRequest(wxAccount, checkFavCdnItem)
		if err != nil {
			return vo.NewFail("GetFavInfoService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
