package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 提取企业 wx 详情
func QWContactApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	reqModel := new(model.QWContactModel)
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWContactService(queryKey, reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 提取全部的企业通寻录
func QWSyncContactApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	result := service.QWSyncContactService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 备注企业 wxid
func QWRemarkApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	reqModel := new(model.QWRemarkModel)
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWRemarkService(queryKey, reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 创建企业群
func QWCreateChatRoomApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	reqModel := new(model.QWCreateModel)
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWCreateChatRoomService(queryKey, reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 搜手机或企业对外名片链接提取验证
func QWSearchContactApi(ctx *gin.Context) {
	reqModel := new(model.SearchContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWSearchContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 向企业微信打招呼
func QWApplyAddContactApi(ctx *gin.Context) {
	reqModel := new(model.QWApplyAddContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWApplyAddContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 单向加企业微信
func QWAddContactApi(ctx *gin.Context) {
	reqModel := new(model.QWApplyAddContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWAddContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 提取全部企业微信群-
func QWSyncChatRoomApi(ctx *gin.Context) {
	reqModel := new(model.QWSyncChatRoomModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWSyncChatRoomService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 转让企业群
func QWChatRoomTransferOwnerApi(ctx *gin.Context) {
	reqModel := new(model.QWChatRoomTransferOwnerModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWChatRoomTransferOwnerService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 直接拉朋友进企业群
func QWAddChatRoomMemberApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWAddChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 发送群邀请链接
func QWInviteChatRoomMemberApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWInviteChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 删除企业群成员
func QWDelChatRoomMemberApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWDelChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 提取企业群全部成员
func QWGetChatRoomMemberApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWGetChatRoomMemberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 提取企业群名称公告设定等信息
func QWGetChatroomInfoApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWGetChatroomInfoService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 提取企业群二维码
func QWGetChatRoomQRApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWGetChatRoomQRService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 增加企业管理员
func QWAppointChatRoomAdminApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWAppointChatRoomAdminService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 移除群管理员
func QWDelChatRoomAdminApi(ctx *gin.Context) {
	reqModel := new(model.QWAddChatRoomMemberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWDelChatRoomAdminService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 同意进企业群
func QWAcceptChatRoomRequestApi(ctx *gin.Context) {
	reqModel := new(model.QWAcceptChatRoomModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWAcceptChatRoomRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 设定企业群管理审核进群
func QWAdminAcceptJoinChatRoomSetApi(ctx *gin.Context) {
	reqModel := new(model.QWAdminAcceptJoinChatRoomSetModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWAdminAcceptJoinChatRoomSetService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改企业群名称
func QWModChatRoomNameApi(ctx *gin.Context) {
	reqModel := new(model.QWModChatRoomNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWModChatRoomNameService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改成员在群中呢称
func QWModChatRoomMemberNickApi(ctx *gin.Context) {
	reqModel := new(model.QWModChatRoomNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWModChatRoomMemberNickService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 发布企业群公告
func QWChatRoomAnnounceApi(ctx *gin.Context) {
	reqModel := new(model.QWModChatRoomNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QWChatRoomAnnounceService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 删除企业群
func QWDelChatRoomApi(ctx *gin.Context) {
	reqModel := new(model.QWModChatRoomNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendQWDelChatRoomService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
