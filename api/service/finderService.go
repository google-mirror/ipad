package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/google/uuid"
	"strconv"
)

// 搜索
func GetFinderSearchService(queryKey string, req model.FinderSearchModel) vo.DTO {
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
		resp, err := reqInvoker.SendGetFinderSearchRequest(req.Index, req.Userver, req.UserKey, uuid.New().String())
		if err != nil {
			return vo.NewFail("GetFinderSearchService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 视频号中心
func FinderUserPrepareService(queryKey string, req model.FinderUserPrepareModel) vo.DTO {
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
		resp, err := reqInvoker.SendFinderUserPrepareRequest(req.Userver)
		if err != nil {
			return vo.NewFail("FinderUserPrepareService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 视频号关注取消
func FinderFollowService(queryKey string, req model.FinderFollowModel) vo.DTO {
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
		resp, err := reqInvoker.SendFinderFollowRequest(req.FinderUserName, req.OpType, req.RefObjectId, req.Cook, req.Userver, req.PosterUsername)
		if err != nil {
			return vo.NewFail("FinderUserPrepareService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 查看指定人首页
func TargetUserPage(queryKey string, req model.TargetUserPageParam) vo.DTO {
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
		resp, err := reqInvoker.TargetUserPageRequest(req.Target, req.LastBuffer)
		if err != nil {
			return vo.NewFail("FinderUserPrepareService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
