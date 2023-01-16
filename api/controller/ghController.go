package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// SendSearchApi 搜索公众号
func SendSearchApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	username, isExist := ctx.GetQuery("username")
	result := service.SearchService(queryKey, username)
	ctx.JSON(http.StatusOK, result)
}

// SendFollowApi 关注公众号
func SendFollowApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	ghId, isExist := ctx.GetQuery("ghId")

	result := service.FollowerService(queryKey, ghId)
	ctx.JSON(http.StatusOK, result)
}

// SendClickMenuApi 操作公众号菜单
func SendClickMenuApi(ctx *gin.Context) {
	reqModel := new(model.ClickCommand)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ClickMenuService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendReadArticleApi 阅读公众号文章
func SendReadArticleApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	url, isExist := ctx.GetQuery("url")
	result := service.ReadArticleService(queryKey, url)
	ctx.JSON(http.StatusOK, result)
}

// SendLikeArticleApi 点赞公众号文章
func SendLikeArticleApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	url, isExist := ctx.GetQuery("url")
	result := service.LikeArticleService(queryKey, url)
	ctx.JSON(http.StatusOK, result)
}
