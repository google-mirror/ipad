package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/container/garray"
	"strconv"
	"sync"
)

// GetContactListService 获取联系人
func GetContactListService(queryKey string, req model.GetContactListModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		//分页获取好友列表
		resp, err := reqInvoker.SendGetContactListPageRequest(req.CurrentWxcontactSeq, req.CurrentChatRoomContactSeq)
		if err != nil {
			return vo.NewSuccess(gin.H{
				"errMsg":  resp.GetBaseResponse().GetErrMsg().GetStr(),
				"retCode": resp.GetBaseResponse().GetRet(),
			}, "")
		}
		return vo.NewSuccess(gin.H{
			"ContactList": resp,
			"errMsg":      resp.GetBaseResponse().GetErrMsg().GetStr(),
			"retCode":     resp.GetBaseResponse().GetRet(),
		}, "")

	})
}

// FollowGHService 关注公众号
func FollowGHService(queryKey string, m model.FollowGHModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		if len(m.GHList) <= 0 {
			return vo.NewFail("没有要关注公众号。")
		}

		respArray := garray.New(true)
		wg := sync.WaitGroup{}

		for _, item := range m.GHList {
			wg.Add(1)
			go func(m model.VerifyUserItem) {
				defer wg.Done()
				resp, err := reqInvoker.VerifyUserRequest(1, "", m.Scene, m.Gh, m.Gh, "")
				if err != nil {
					respArray.Append(gin.H{
						"GH":              m.Gh,
						"isFollowSuccess": false,
						"errMsg":          err.Error(),
					})
					return
				}
				isFollowSuccess := false
				if resp.GetBaseResponse().GetRet() == 0 {
					isFollowSuccess = true
				}

				respArray.Append(gin.H{
					"GH":              m.Gh,
					"isFollowSuccess": isFollowSuccess,
					"errMsg":          resp.GetBaseResponse().GetErrMsg().GetStr(),
					"retCode":         resp.GetBaseResponse().GetRet(),
				})

			}(item)

		}
		wg.Wait()
		return vo.NewSuccessObj(respArray.Interfaces(), "")
	})
}

// UploadMContactService 上传手机通讯录好友
func UploadMContactService(queryKey string, m model.UploadMContactModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		if len(m.MobileList) == 0 {
			return vo.NewFail("没有要上传的手机号！")
		}
		/*lolDeviceInfoOsType := wxAccount.GetUserInfo().DeviceInfo.OsType
		lolDeviceInfoOsTypeNumber := wxAccount.GetUserInfo().DeviceInfo.OsTypeNumber
		userInfo := wxAccount.GetUserInfo().DeviceInfo
		if !strings.HasSuffix(userInfo.OsType, "HUAWEI android") {
			userInfo.OsType = "HUAWEI android"
			userInfo.OsTypeNumber = "android 9.0"
			reqInvoker.SendAutoAuthRequest()
		}*/
		resp, err := reqInvoker.UploadMContact(m.Mobile, m.MobileList)
		if err != nil {
			return vo.NewFail("UploadMContactService - " + err.Error())
		}
		/*if strings.HasSuffix(userInfo.OsType, "HUAWEI android") {
			userInfo := wxAccount.GetUserInfo().DeviceInfo
			userInfo.OsType = lolDeviceInfoOsType
			userInfo.OsTypeNumber = lolDeviceInfoOsTypeNumber
			reqInvoker.SendAutoAuthRequest()
		}*/
		return vo.NewSuccessObj(resp, "上传成功！")
	})
}

// GetMFriendService 获取手机通讯录好友
func GetMFriendService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.GetMFriend()
		if err != nil {
			return vo.NewFail("GetMFriendService - " + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// VerifyUserRequestService // 添加好友 // 关注公众号 // 同意好友请求
func VerifyUserRequestService(queryKey string, m model.VerifyUserRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.VerifyUserRequest(m.OpCode, m.VerifyContent, m.Scene, m.V1, m.V2, m.ChatRoomUserName)
		if err != nil {
			return vo.NewFail("VerifyUserRequestService err" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
