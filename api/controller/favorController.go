package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// FavSyncApi 同步收藏
func FavSyncApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	result := service.FavSyncService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 获取收藏list
func GetFavListApi(ctx *gin.Context) {
	reqModel := new(model.FavInfoModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetFavListService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// GetFavInfoApi 获取收藏信息
func GetFavInfoApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.GetFavInfoService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// BatchDelFavItemApi 删除收藏
func BatchDelFavItemApi(ctx *gin.Context) {
	reqModel := new(model.FavInfoModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.BatchDelFavItemService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// BatchGetFavItemApi 获取收藏详细
func BatchGetFavItemApi(ctx *gin.Context) {
	reqModel := new(model.FavInfoModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.BatchGetFavItemService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ShareFavServiceApi 分享收藏
func ShareFavServiceApi(ctx *gin.Context) {
	reqModel := new(model.ShareFavModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ShareFavService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// CheckFavCdnServiceApi 检测收藏cdn
func CheckFavCdnServiceApi(ctx *gin.Context) {
	reqModel := new(model.CheckFavCdnModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.CheckFavCdnService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
