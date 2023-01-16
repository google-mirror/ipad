package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/utils"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	pb "feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"strconv"
	"strings"
)

// SetChatroomAnnouncement 设置群公告
func SetChatroomAnnouncementService(queryKey string, m model.UpdateChatroomAnnouncementModel) vo.DTO {
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

		resp, err := reqInvoker.SetChatRoomAnnouncementRequest(m.ChatRoomName, m.Content)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccess(gin.H{
			"baseResp": resp.GetBaseResponse(),
		}, "")
	})
}

// GetChatroomMemberDetailService 获取群成员详细
func GetChatroomMemberDetailService(queryKey string, m model.GetChatroomMemberDetailModel) vo.DTO {
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

		resp, err := reqInvoker.GetChatroomMemberDetailRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 获取群公告
func SetGetChatRoomInfoDetailService(queryKey string, m model.GetChatroomMemberDetailModel) vo.DTO {
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

		resp, err := reqInvoker.SetGetChatRoomInfoDetailRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 获取群详情
func GetChatRoomInfoService(queryKey string, m model.ChatRoomWxIdListModel) vo.DTO {
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
		resp, err := reqInvoker.SendGetContactRequest(m.ChatRoomWxIdList, nil, m.ChatRoomWxIdList, true)
		if err != nil {
			return vo.NewFail(err.Error())
		}

		return vo.NewSuccessObj(resp, "")
	})
}

// 设置群昵称
func SetChatroomNameService(queryKey string, m model.ChatroomNameModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		chatRoomNames := strings.Split(m.ChatRoomName, ",")
		// 获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		getContactResp, err := reqInvoker.SendGetContactRequest(chatRoomNames, nil, chatRoomNames, true)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		modContact := getContactResp.ContactList[0]
		modContact.NickName = &pb.SKBuiltinString{
			Str: proto.String(m.Nickname),
		}
		buffer, err := proto.Marshal(modContact)
		cmdItem := baseinfo.ModifyItem{
			CmdID: uint32(27),
			Len:   uint32(len(buffer)),
			Data:  buffer,
		}
		var cmdItems []*baseinfo.ModifyItem
		cmdItems = append(cmdItems, &cmdItem)
		error := reqInvoker.SendOplogRequest(cmdItems)
		if error != nil {
			return vo.NewFail(error.Error())
		}
		return vo.NewSuccessObj(nil, "成功")
	})
}

// 保存群聊
func MoveToContractService(queryKey string, m model.MoveContractModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		chatRoomNames := strings.Split(m.ChatRoomName, ",")
		// 获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		getContactResp, err := reqInvoker.SendGetContactRequest(chatRoomNames, nil, chatRoomNames, true)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		modContact := getContactResp.ContactList[0]
		bit := uint32(0)
		if m.Val == 1 {
			bit = *(modContact.BitVal) | uint32(1<<0)
		} else {
			bit = *(modContact.BitVal) &^ uint32(1<<0)
		}
		ModContactData := &pb.ModContact{
			UserName: &pb.SKBuiltinString{
				Str: &m.ChatRoomName,
			},
			NickName:  &pb.SKBuiltinString{},
			Pyinitial: &pb.SKBuiltinString{},
			QuanPin:   &pb.SKBuiltinString{},
			Sex:       proto.Int32(0),
			ImgBuf: &pb.SKBuiltinString_{
				Len: proto.Uint32(0),
			},
			BitMask: modContact.BitMask,
			BitVal:  proto.Uint32(bit),
			ImgFlag: proto.Uint32(0),
			Remark: &pb.SKBuiltinString{
				Str: modContact.Remark.Str,
			},
			RemarkPyinitial: &pb.SKBuiltinString{
				Str: modContact.RemarkPyinitial.Str,
			},
			RemarkQuanPin: &pb.SKBuiltinString{
				Str: modContact.RemarkQuanPin.Str,
			},
			ContactType:     proto.Uint32(0),
			ChatRoomNotify:  proto.Uint32(1),
			AddContactScene: proto.Uint32(0),
			ExtFlag:         proto.Uint32(0),
		}
		buffer, err := proto.Marshal(ModContactData)
		cmdItem := baseinfo.ModifyItem{
			CmdID: uint32(2),
			Len:   uint32(len(buffer)),
			Data:  buffer,
		}
		var cmdItems []*baseinfo.ModifyItem
		cmdItems = append(cmdItems, &cmdItem)
		error := reqInvoker.SendOplogRequest(cmdItems)
		if error != nil {
			return vo.NewFail(error.Error())
		}
		return vo.NewSuccessObj(nil, "成功")
	})
}

// QuitChatroomService 退出群聊
func QuitChatroomService(queryKey string, m model.GetChatroomMemberDetailModel) vo.DTO {
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

		err := reqInvoker.GetQuitChatroomRequest(m.ChatRoomName)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(nil, "发送退群请求成功")
	})
}

// CreateChatRoomService 创建群
func CreateChatRoomService(queryKey string, m model.CreateChatRoomModel) vo.DTO {
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
		resp, err := reqInvoker.SendCreateChatRoomRequest(m.TopIc, m.UserList)
		if err != nil {
			return vo.NewFail("创建群失败！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// InviteChatroomMembersService 邀请群成员
func InviteChatroomMembersService(queryKey string, m model.InviteChatroomMembersModel) vo.DTO {
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
		resp, err := reqInvoker.SendInviteChatroomMembersRequest(m.ChatRoomName, m.UserList)
		if err != nil {
			return vo.NewFail("InviteChatroomMembersService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// AddChatRoomMemberService 添加群成员
func AddChatRoomMemberService(queryKey string, m model.InviteChatroomMembersModel) vo.DTO {
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
		resp, err := reqInvoker.SendAddChatRoomMemberRequest(m.ChatRoomName, m.UserList)
		if err != nil {
			return vo.NewFail("InviteChatroomMembersService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 删除群成员
func SendDelDelChatRoomMemberService(queryKey string, m model.InviteChatroomMembersModel) vo.DTO {
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
		resp, err := reqInvoker.SendDelDelChatRoomMemberRequest(m.ChatRoomName, m.UserList)
		if err != nil {
			return vo.NewFail("InviteChatroomMembersService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 转让群
func SendTransferGroupOwnerService(queryKey string, m model.TransferGroupOwnerModel) vo.DTO {
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
		resp, err := reqInvoker.SendTransferGroupOwnerRequest(m.ChatRoomName, m.NewOwnerUserName)
		if err != nil {
			return vo.NewFail("InviteChatroomMembersService！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 扫码入群
func ScanIntoUrlGroupService(queryKey string, m model.ScanIntoUrlGroupModel) vo.DTO {
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
		if len(m.Url) == 0 {
			return vo.NewFail("ScanIntoUrlGroupService url == empty")
		}
		resp, err := reqInvoker.GetA8KeyGroupRequest(2, 4, m.Url, baseinfo.GetA8Key)
		if err != nil {
			return vo.NewFail("ScanIntoUrlGroupService err:" + err.Error())
		}
		if resp.GetBaseResponse().GetRet() != 0 {
			return vo.NewSuccess(gin.H{
				"isJoinSuccess": false,
				"resp":          resp,
			}, "进群失败")
		}
		deviceType := 0
		if wxAccount.GetUserInfo().DeviceInfo != nil {
			deviceType = 1
		}
		body, err := utils.ScanIntoGrouppost(resp.GetFullURL(), deviceType, wxAccount.GetUserInfo())
		if err != nil && strings.Index(err.Error(), "@chatroom") != -1 {
			return vo.NewSuccess(gin.H{
				"isJoinSuccess": true,
				"body":          body,
				"chatroomUrl":   err.Error(),
				"fullUrl":       resp.GetFullURL(),
				"resp":          resp,
			}, "进群成功")
		}
		msg := ""
		if strings.Index(body, "频繁") != -1 {
			msg = ",操作太频繁，请稍后再试!"
		}
		if strings.Index(body, "二维码已过期") != -1 {
			msg = ",二维码已过期！"
		}
		if strings.Index(body, "该群聊邀请已过期") != -1 {
			msg = ",该群聊邀请已过期！"
		}
		return vo.NewSuccess(gin.H{
			"isJoinSuccess": true,
			"body":          body,
			"fullUrl":       resp.GetFullURL(),
			"resp":          resp,
		}, "进群失败"+msg)

	})
}

// 设置群聊邀请开关
func SetChatroomAccessVerifyService(queryKey string, m model.SetChatroomAccessVerifyModel) vo.DTO {
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
		v := uint32(0)
		if m.Enable {
			v = uint32(2)
		}
		req := &wechat.ModChatRoomAccessVerifyRequest{
			ChatRoomName: proto.String(m.ChatRoomName),
			Status:       proto.Uint32(v),
		}
		buffer, err := proto.Marshal(req)
		cmdItem := baseinfo.ModifyItem{
			CmdID: 66,
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

// 添加群管理员
func AddChatroomAdminService(queryKey string, m model.ChatroomMemberModel) vo.DTO {
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
		resp, err := reqInvoker.SendAddChatroomAdminRequest(m.ChatRoomName, m.UserList)
		if err != nil {
			return vo.NewFail("SendAddChatroomAdminRequest！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 删除群管理
func DelChatroomAdminService(queryKey string, m model.ChatroomMemberModel) vo.DTO {
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
		resp, err := reqInvoker.SendDelChatroomAdminRequest(m.ChatRoomName, m.UserList)
		if err != nil {
			return vo.NewFail("SendDelChatroomAdminRequest！err :" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
