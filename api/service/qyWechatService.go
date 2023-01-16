package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/gogo/protobuf/proto"
	"strconv"
)

func QWContactService(queryKey string, m *model.QWContactModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendQWContactRequest(m.ToUserName, m.ChatRoom, m.T)
		if err != nil {
			return vo.NewFail("查询失败")
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 提取全部的企业通寻录
func QWSyncContactService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendQWSyncContactRequest()
		if err != nil {
			return vo.NewFail("查询失败")
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 备注企业 wxid
func QWRemarkService(queryKey string, m *model.QWRemarkModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		req := &wechat.QYModChatRoomTopicRequest{
			G: proto.String(m.ToUserName),
			P: proto.String(m.Name),
		}
		buffer, err := proto.Marshal(req)
		err = reqInvoker.SendQWOpLogRequest(3, buffer)
		if err != nil {
			return vo.NewFail("失败")
		}
		return vo.NewSuccessObj("ok", "操作成功")
	})
}

// 创建企业微信
func QWCreateChatRoomService(queryKey string, m *model.QWCreateModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendQWCreateChatRoomRequest(m.ToUserName)
		if err != nil {
			return vo.NewFail("失败")
		}
		return vo.NewSuccessObj(resp, "操作成功")
	})
}

// 搜手机或企业对外名片链接提取验证
func QWSearchContactService(queryKey string, m model.SearchContactModel) vo.DTO {
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
		resp, err := reqInvoker.SendQWSearchContactRequest(m.Tg, m.FromScene, m.UserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 向企业微信打招呼
func QWApplyAddContactService(queryKey string, m model.QWApplyAddContactModel) vo.DTO {
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
		err := reqInvoker.SendQWApplyAddContactRequest(m.UserName, m.V1, m.Content)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj("ok", "操作成功")
	})
}

// 单向加企业微信
func QWAddContactService(queryKey string, m model.QWApplyAddContactModel) vo.DTO {
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
		err := reqInvoker.SendQWAddContactRequest(m.UserName, m.V1, m.Content)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj("ok", "操作成功")
	})
}

// 提取全部企业微信群
func QWSyncChatRoomService(queryKey string, m model.QWSyncChatRoomModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWSyncChatRoomRequest(m.Key)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 转让企业群
func QWChatRoomTransferOwnerService(queryKey string, m model.QWChatRoomTransferOwnerModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWChatRoomTransferOwnerRequest(m.ChatRoomName, m.ToUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 直接拉好友进群
func QWAddChatRoomMemberService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWAddChatRoomMemberRequest(m.ChatRoomName, m.ToUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 发送群邀请链接
func QWInviteChatRoomMemberService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWInviteChatRoomMemberRequest(m.ChatRoomName, m.ToUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 删除企业群成员
func QWDelChatRoomMemberService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWDelChatRoomMemberRequest(m.ChatRoomName, m.ToUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 提取企业群全部成员
func QWGetChatRoomMemberService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWGetChatRoomMemberRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// -提取企业群名称公告设定等信息
func QWGetChatroomInfoService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWGetChatroomInfoRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 提取企业群二维码
func QWGetChatRoomQRService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWGetChatRoomQRRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 增加企业管理员
func QWAppointChatRoomAdminService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWAppointChatRoomAdminRequest(m.ChatRoomName, m.ToUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 移除群管理
func QWDelChatRoomAdminService(queryKey string, m model.QWAddChatRoomMemberModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWDelChatRoomAdminRequest(m.ChatRoomName, m.ToUserName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 同意进企业群
func QWAcceptChatRoomRequestService(queryKey string, m model.QWAcceptChatRoomModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWAcceptChatRoomRequest(m.Link, m.Opcode)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 设定企业群
func QWAdminAcceptJoinChatRoomSetService(queryKey string, m model.QWAdminAcceptJoinChatRoomSetModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWAdminAcceptJoinChatRoomSetRequest(m.ChatRoomName, m.P)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 修改企业群名称
func QWModChatRoomNameService(queryKey string, m model.QWModChatRoomNameModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWModChatRoomNameRequest(m.ChatRoomName, m.Name)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 修改成员在群中呢称
func QWModChatRoomMemberNickService(queryKey string, m model.QWModChatRoomNameModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWModChatRoomMemberNickRequest(m.ChatRoomName, m.Name)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 发布企业群公告
func QWChatRoomAnnounceService(queryKey string, m model.QWModChatRoomNameModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWChatRoomAnnounceRequest(m.ChatRoomName, m.Name)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}

// 删除企业群
func SendQWDelChatRoomService(queryKey string, m model.QWModChatRoomNameModel) vo.DTO {
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
		rsp, err := reqInvoker.SendQWDelChatRoomRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(rsp, "操作成功")
	})
}
