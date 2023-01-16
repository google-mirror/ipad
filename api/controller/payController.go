package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 获取银行卡信息
func GetBandCardListApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.GetBandCardListService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 生成自定义二维码
func GeneratePayQCodeApi(ctx *gin.Context) {
	reqModel := new(model.GeneratePayQCodeModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GeneratePayQCodeService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 确定收款
func CollectmoneyApi(ctx *gin.Context) {
	reqModel := new(model.CollectmoneyModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.CollectmoneyServie(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 创建红包
func WXCreateRedPacketApi(ctx *gin.Context) {
	reqModel := new(baseinfo.RedPacket)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.WXCreateRedPacketService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 拆红包
func OpenRedEnvelopesApi(ctx *gin.Context) {
	reqModel := new(baseinfo.HongBaoItem)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.OpenRedEnvelopesService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 查看红包详情
func QueryRedEnvelopesDetailApi(ctx *gin.Context) {
	reqModel := new(baseinfo.HongBaoItem)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.QueryRedEnvelopesDetailService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 查看红包领取列表
func GetRedPacketListApi(ctx *gin.Context) {
	reqModel := new(baseinfo.GetRedPacketList)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetRedPacketListService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 创建转账
func CreatePreTransferApi(ctx *gin.Context) {
	reqModel := new(baseinfo.CreatePreTransfer)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.CreatePreTransferService(queryKey, reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 确认转账
func ConfirmPreTransferApi(ctx *gin.Context) {
	reqModel := new(baseinfo.ConfirmPreTransfer)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ConfirmPreTransferService(queryKey, reqModel)
	ctx.JSON(http.StatusOK, result)
}
