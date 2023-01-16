package controller

import (
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func validateData(ctx *gin.Context, model interface{}) bool {
	err := ctx.ShouldBindJSON(&model)
	if err != nil {
		ctx.JSON(http.StatusOK, vo.DTO{
			Code: vo.FAIL_DATA,
			Data: nil,
			Text: "\"提交数据错误!\"",
		})
		ctx.Abort()
		return false
	}
	return true
}
