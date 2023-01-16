package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetQrCodeApi 获取二维码
func GetQrCodeApi(ctx *gin.Context) {
	reqModel := new(model.GetQrCodeModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetQrCodeService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 查看附近的人
func GetPeopleNearbyApi(ctx *gin.Context) {
	reqModel := new(model.PeopleNearbyModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetPeopleNearbyService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
