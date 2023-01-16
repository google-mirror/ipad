package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/lunny/log"
	"strconv"
	"strings"
)

// SearchContactRequestService 搜索好友
func SearchContactRequestService(queryKey string, m model.SearchContactRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendSearchContactRequest(m.OpCode, m.FromScene, m.SearchScene, m.UserName)
		if err != nil {
			return vo.NewFail("SearchContactRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 获取好友关系
func GetFriendRelationService(queryKey string, m model.GetFriendRelationModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGetFriendRelationRequest(m.UserName)
		if err != nil {
			return vo.NewFail("GetFriendRelationService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 获取好友关系
func GetFriendRelationsService(queryKey string, m model.GetFriendRelationModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		usernames := strings.Split(m.UserName, ",")
		if len(usernames) > 0 {
			respMap := make(map[string]*wechat.MMBizJsApiGetUserOpenIdResponse)
			for _, value := range usernames {
				resp, err := reqInvoker.SendGetFriendRelationRequest(value)
				if err != nil {
					log.Error("GetFriendRelationService err:" + err.Error())
				}
				respMap[value] = resp
			}
			return vo.NewSuccessObj(respMap, "")
		}
		return vo.NewFail("usernames can not null")
	})
}

// 查联系人详情
func GetContactContactService(queryKey string, m model.BatchGetContactModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGetContactRequestForList(m.UserNames, m.RoomWxIDList)
		if err != nil {
			return vo.NewFail("SearchContactRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
