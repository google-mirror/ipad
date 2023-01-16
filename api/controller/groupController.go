package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// SetChatroomAnnouncementApi 设置群公告
func SetChatroomAnnouncementApi(ctx *gin.Context) {
	reqModel := new(model.UpdateChatroomAnnouncementModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetChatroomAnnouncementService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取群成员详细
func GetChatroomMemberDetailApi(ctx *gin.Context) {
	reqModel := new(model.GetChatroomMemberDetailModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetChatroomMemberDetailService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

func GetQuitChatroomApi(ctx *gin.Context) {
	reqModel := new(model.GetChatroomMemberDetailModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QuitChatroomService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// CreateChatRoomApi 创建群请求
func CreateChatRoomApi(ctx *gin.Context) {
	reqModel := new(model.CreateChatRoomModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.CreateChatRoomService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// InviteChatroomMembersApi 邀请群成员
func InviteChatroomMembersApi(ctx *gin.Context) {
	reqModel := new(model.InviteChatroomMembersModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.InviteChatroomMembersService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// AddChatRoomMemberApi 添加群成员
func AddChatRoomMembersApi(ctx *gin.Context) {
	reqModel := new(model.InviteChatroomMembersModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.AddChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ScanIntoUrlGroupApi 扫码入群
func ScanIntoUrlGroupApi(ctx *gin.Context) {
	reqModel := new(model.ScanIntoUrlGroupModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ScanIntoUrlGroupService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 添加好友进群
func SendAddChatRoomMemberApi(ctx *gin.Context) {
	reqModel := new(model.InviteChatroomMembersModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.AddChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 删除群成员
func SendDelDelChatRoomMemberApi(ctx *gin.Context) {
	reqModel := new(model.InviteChatroomMembersModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendDelDelChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 转让群
func SendTransferGroupOwnerApi(ctx *gin.Context) {
	reqModel := new(model.TransferGroupOwnerModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendTransferGroupOwnerService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取群公告
func SetGetChatRoomInfoDetailApi(ctx *gin.Context) {
	reqModel := new(model.GetChatroomMemberDetailModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetGetChatRoomInfoDetailService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取群详情
func GetChatRoomInfoApi(ctx *gin.Context) {
	reqModel := new(model.ChatRoomWxIdListModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetChatRoomInfoService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取群聊
func MoveToContractApi(ctx *gin.Context) {
	reqModel := new(model.MoveContractModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.MoveToContractService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 设置群聊邀请开关
func SetChatroomAccessVerifyApi(ctx *gin.Context) {
	reqModel := new(model.SetChatroomAccessVerifyModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetChatroomAccessVerifyService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 添加群管理员
func AddChatroomAdminApi(ctx *gin.Context) {
	reqModel := new(model.ChatroomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.AddChatroomAdminService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 删除群管理员
func DelChatroomAdminApi(ctx *gin.Context) {
	reqModel := new(model.ChatroomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.DelChatroomAdminService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 设置群昵称
func SetChatroomNameApi(ctx *gin.Context) {
	reqModel := new(model.ChatroomNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetChatroomNameService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 群拍一拍功能
func SendPatApi(ctx *gin.Context) {
	reqModel := new(model.SendPatModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendPatService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取群例表
func GroupListApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	result := service.NewSyncGroupMessageService(queryKey)
	ctx.JSON(http.StatusOK, result)
}
