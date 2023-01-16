package service

import (
	"encoding/hex"
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"strconv"
)

// logOutService 退出登录
func LogOutService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		err := bizcgi.SendLogoutRequest(wxAccount)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccess(gin.H{}, "退出成功！")
	})
}

// SendDelContactService 删除好友
func SendDelContactService(queryKey string, m model.DelContactModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		err := reqInvoker.SendDelContactRequest(m.DelUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccess(gin.H{}, "删除成功！")
	})
}

// OnlineInfoService 获取登录设备信息
func OnlineInfoService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendOnlineInfo()
		if err != nil {
			return vo.NewFail(err.Error())
		}

		deviceInfoList := make([]gin.H, 0)
		for _, deviceInfo := range resp.GetOnlineList() {
			deviceInfoList = append(deviceInfoList, gin.H{
				"deviceType":   deviceInfo.GetDeviceType(),
				"deviceId":     hex.EncodeToString(deviceInfo.GetDeviceID()),
				"clientKey":    deviceInfo.GetClientKey(),
				"onlineStatus": deviceInfo.GetOnlineStatus(),
			})
		}

		return vo.NewSuccessObj(gin.H{
			"onlineCount": resp.GetOnlineCount(),
			"onlineList":  deviceInfoList,
		}, "")
	})
}

// GetProfileService 获取个人资料详细信息
func GetProfileService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		/*err := reqInvoker.SendGetProfileRequest()
		if err != nil {
			return vo.NewFail("获取个人资料失败")
		}*/
		profile, err := reqInvoker.SendGetProfileNewRequest()
		if err != nil {
			return vo.NewFail("获取个人资料失败")
		}
		return vo.NewSuccessObj(profile, "")
	})
}

// SendModifyUserInfoRequest 修改资料
func SendModifyUserInfoRequestService(queryKey string, m model.ModifyUserInfo) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		err := reqInvoker.SendModifyUserInfoRequest(m.City, m.Country, m.NickName, m.Province, m.Signature, m.Sex, m.InitFlag)
		if err != nil {
			return vo.NewFail("修改资料失败！")
		}
		return vo.NewSuccessObj(nil, "修改资料成功！")
	})
}

// 修改名称
func UpdateNickNameService(queryKey string, m model.UpdateNickNameModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		err := reqInvoker.SendUpdateNickNameRequest(m.Scene, m.Val)
		if err != nil {
			return vo.NewFail("修改名称失败！")
		}
		return vo.NewSuccessObj(nil, "修改成功！")
	})
}

// 设置名称
func SetNickNameService(queryKey string, m model.UpdateNickNameModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		err := reqInvoker.SendUpdateNickNameRequest(m.Scene, m.Val)
		if err != nil {
			return vo.NewFail("修改失败！")
		}
		return vo.NewSuccessObj(nil, "修改成功！")
	})
}

// 修改姓别
func SetSexService(queryKey string, m model.UpdateSexModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		err := reqInvoker.SetSexService(m.Sex, m.Country, m.City, m.Province)
		if err != nil {
			return vo.NewFail("修改失败！")
		}
		return vo.NewSuccessObj(nil, "修改成功！")
	})
}

func GetRedisSyncMsgService(queryKey string, m model.GetSyncMsgModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendWxSyncMsg(m.Key)
		// 如果没有同步到数据则返回
		cmdList := resp.GetCmdList()
		syncCount := cmdList.GetCount()
		messageResp := new(table.SyncMessageResponse)
		// 遍历同步的信息和群
		itemList := cmdList.GetItemList()
		for index := uint32(0); index < syncCount; index++ {
			item := itemList[index]
			itemID := item.GetCmdId()
			// 同步到消息
			if itemID == baseinfo.CmdIDAddMsg {
				addMsg := &wechat.AddMsg{}
				err := proto.Unmarshal(item.CmdBuf.Data, addMsg)
				if err != nil {
					baseutils.PrintLog(err.Error())
					continue
				}
			}
			messageResp.SetMessage(item.GetCmdBuf().GetData(), int32(itemID))
		}
		//发布同步信息消息
		if err != nil {
			return vo.NewFail(err.Error())
		}
		messageResp.Key = resp.KeyBuf
		return vo.NewSuccessObj(*messageResp, "成功")
	})
}

// SendChangePwdRequestService 更改密码
func SendChangePwdRequestService(queryKey string, m model.SendChangePwdRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendChangePwdRequest(m.OldPass, m.NewPass, m.OpCode)
		if err != nil {
			return vo.NewFail("SendChangePwdRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 上传头像
func UploadHeadImageService(queryKey string, m model.UploadHeadImageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.UploadHeadImage(m.Base64)
		if err != nil {
			return vo.NewFail("SendChangePwdRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "成功")
	})
}

// endModifyRemarkRequestService 修改备注
func SendModifyRemarkRequestService(queryKey string, m model.SendModifyRemarkRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		err := reqInvoker.SendModifyRemarkRequest(m.UserName, m.RemarkName)
		if err != nil {
			return vo.NewFail("SendChangePwdRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(nil, "成功")
	})
}

// 修改加好友需要验证属性
func UpdateAutopassService(queryKey string, m model.UpdateAutopassModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		err := reqInvoker.UpdateAutopassRequest(m.SwitchType)
		if err != nil {
			return vo.NewFail("UpdateAutopassService err:" + err.Error())
		}
		return vo.NewSuccessObj(nil, "成功")
	})
}

// 设置微信号
func SetWechatService(queryKey string, m model.AlisaModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")

		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.SetWechatRequest(m.Alisa)
		if err != nil {
			return vo.NewFail("UpdateAutopassService err:" + err.Error())
		}
		return vo.NewSuccessObj(rsp, "成功")
	})
}

// 修改步数
func UpdateStepNumberService(queryKey string, m model.UpdateStepNumberModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")
		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.UpdateStepNumberRequest(m.Number)
		if err != nil {
			return vo.NewFail("UpdateStepNumberService err:" + err.Error())
		}
		return vo.NewSuccessObj(rsp, "成功")
	})
}

// 获取步数
func GetUserRankLikeCountService(queryKey string, m model.UserRankLikeModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")
		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.SendGetUserRankLikeCountRequest(m.RankId)
		if err != nil {
			return vo.NewFail("UpdateStepNumberService err:" + err.Error())
		}
		return vo.NewSuccessObj(rsp, "成功")
	})
}

// 设置添加我的方式
func SetFunctionSwitchService(queryKey string, m model.WxFunctionSwitchModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")
		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		req := &wechat.FunctionSwitch{
			FunctionId:  proto.Uint32(m.Function),
			SwitchValue: proto.Uint32(m.Value),
		}
		buffer, err := proto.Marshal(req)
		cmdItem := baseinfo.ModifyItem{
			CmdID: 0x17,
			Len:   uint32(len(buffer)),
			Data:  buffer,
		}
		var cmdItems []*baseinfo.ModifyItem
		cmdItems = append(cmdItems, &cmdItem)
		err = reqInvoker.SendOplogRequest(cmdItems)
		if err != nil {
			return vo.NewFail("SetFunctionSwitchService err:" + err.Error())
		}
		return vo.NewSuccessObj("ok", "成功")
	})
}

// 设置拍一拍名称
func SetSendPatService(queryKey string, m model.SetSendPatModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")
		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		req := &wechat.PatMod{
			Value: proto.Int64(8),
			Name:  proto.String(m.Value),
		}
		buffer, err := proto.Marshal(req)
		cmdItem := baseinfo.ModifyItem{
			CmdID: 222,
			Len:   uint32(len(buffer)),
			Data:  buffer,
		}
		var cmdItems []*baseinfo.ModifyItem
		cmdItems = append(cmdItems, &cmdItem)
		err = reqInvoker.SendOplogRequest(cmdItems)
		if err != nil {
			return vo.NewFail("SetSendPatService err:" + err.Error())
		}
		return vo.NewSuccessObj("ok", "成功")
	})
}

// 换绑手机
func BindingMobileService(queryKey string, m model.BindMobileModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")
		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendBindingMobileRequest(m.Mobile, m.VerifyCode)
		if err != nil {
			return vo.NewFail("SetSendPatService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "成功")
	})
}

// 发送验证码
func SendVerifyMobileService(queryKey string, m model.SendVerifyMobileModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("你已退出登录")
		} else if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendVerifyMobileRequest(m.Mobile, m.Opcode)
		if err != nil {
			return vo.NewFail("SetSendPatService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "成功")
	})
}
