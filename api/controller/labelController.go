package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// SendSnsUserPageRequestApi 获取标签列表
func GetContactLabelListApi(ctx *gin.Context) {

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.GetContactLabelListRequestService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// AddContactLabelRequestApi 添加列表
func AddContactLabelRequestApi(ctx *gin.Context) {
	reqModel := new(model.LabelModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.AddContactLabelRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// DelContactLabelRequestApi 删除标签
func DelContactLabelRequestApi(ctx *gin.Context) {
	reqModel := new(model.LabelModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.DelContactLabelRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ModifyLabelRequestApi 修改标签
func ModifyLabelRequestApi(ctx *gin.Context) {
	reqModel := new(model.LabelModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.ModifyLabelRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

////获取标签下所有好友
//func GetWXFriendListByLabelIDApi(ctx *gin.Context) {
//	reqModel := new(model.LabelModel)
//	queryKey, isExist := ctx.GetQuery("key")
//	if !isExist || strings.Trim(queryKey, "") == "" {
//		//确保每次都有Key
//		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
//		return
//	}
//	if !validateData(ctx, &reqModel) {
//		return
//	}
//
//	result := service.GetWXFriendListByLabelIDService(queryKey, *reqModel)
//	ctx.JSON(http.StatusOK, result)
//}
