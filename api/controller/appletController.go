package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetA8KeyApi 授权链接
func GetA8KeyApi(ctx *gin.Context) {
	reqModel := new(model.GetA8KeyRequestModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.GetA8KeyService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// JSLoginApi 小程序授权
func JSLoginApi(ctx *gin.Context) {
	reqModel := new(model.AppletModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.JsLoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// JSOperateWxDataApi 小程序授权
func JSOperateWxDataApi(ctx *gin.Context) {
	reqModel := new(model.AppletModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.JSOperateWxDataService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SdkOauthAuthorizeApi app 应用授权
func SdkOauthAuthorizeApi(ctx *gin.Context) {
	reqModel := new(model.AppletModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SdkOauthAuthorizeService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// QRConnectAuthorize 二维码授权请求
func QRConnectAuthorizeApi(ctx *gin.Context) {
	reqModel := new(model.QRConnectAuthorizeModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.QRConnectAuthorizeService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// QRConnectAuthorizeConfirmApi 二维码授权确认
func QRConnectAuthorizeConfirmApi(ctx *gin.Context) {
	reqModel := new(model.QRConnectAuthorizeModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.QRConnectAuthorizeConfirmService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 授权链接
func GetMpA8KeyApi(ctx *gin.Context) {
	reqModel := new(model.GetMpA8KeyModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.GetMpA8Service(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 授权公众号登录
func AuthMpLoginApi(ctx *gin.Context) {
	reqModel := new(model.GetMpA8KeyModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.AuthMpLoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
