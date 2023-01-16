package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetFinderSearchApi(ctx *gin.Context) {
	reqModel := new(model.FinderSearchModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetFinderSearchService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 视频号中心
func FinderUserPrepareApi(ctx *gin.Context) {
	reqModel := new(model.FinderUserPrepareModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.FinderUserPrepareService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 关注取消
func FinderFollowApi(ctx *gin.Context) {
	reqModel := new(model.FinderFollowModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.FinderFollowService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 查看指定人首页
func TargetUserPageApi(ctx *gin.Context) {
	reqModel := new(model.TargetUserPageParam)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.TargetUserPage(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
