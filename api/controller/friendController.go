package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetContactListApi 获取全部联系人
func GetContactListApi(ctx *gin.Context) {
	reqModel := new(model.GetContactListModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetContactListService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

//// GetFriendListApi 获取好友列表
//func GetFriendListApi(ctx *gin.Context) {
//	queryKey, isExist := ctx.GetQuery("key")
//	if !isExist || strings.Trim(queryKey, "") == "" {
//		//确保每次都有Key
//		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
//		return
//	}
//
//	result := service.GetFriendListService(queryKey)
//	ctx.JSON(http.StatusOK, result)
//}
//
//
//// GetGHListApi 获取好友列表
//func GetGHListApi(ctx *gin.Context) {
//	queryKey, isExist := ctx.GetQuery("key")
//	if !isExist || strings.Trim(queryKey, "") == "" {
//		//确保每次都有Key
//		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
//		return
//	}
//
//	result := service.GetGHListService(queryKey)
//	ctx.JSON(http.StatusOK, result)
//}

// FollowGHApi 关注公众号
func FollowGHApi(ctx *gin.Context) {

	reqModel := new(model.FollowGHModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.FollowGHService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// UploadMContactApi
func UploadMContactApi(ctx *gin.Context) {
	reqModel := new(model.UploadMContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.UploadMContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// GetMFriendApi
func GetMFriendApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.GetMFriendService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 获取联系人详情
func GetContactContactApi(ctx *gin.Context) {
	reqModel := new(model.BatchGetContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetContactContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取好友关系
func GetFriendRelationApi(ctx *gin.Context) {
	reqModel := new(model.GetFriendRelationModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetFriendRelationService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取好友关系
func GetFriendRelationsApi(ctx *gin.Context) {
	reqModel := new(model.GetFriendRelationModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetFriendRelationsService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SearchContactRequestApi 搜索联系人
func SearchContactRequestApi(ctx *gin.Context) {
	reqModel := new(model.SearchContactRequestModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SearchContactRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// VerifyUserRequestApi 验证用户
func VerifyUserRequestApi(ctx *gin.Context) {
	reqModel := new(model.VerifyUserRequestModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.VerifyUserRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 同意好友请求
func AgreeAddApi(ctx *gin.Context) {
	reqModel := new(model.VerifyUserRequestModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	if reqModel.Scene == 0 {
		reqModel.Scene = 0x06
	}
	result := service.VerifyUserRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
