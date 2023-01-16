package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"github.com/gogf/guuid"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetLoginQrCodeApi 获取登录二维码接口
func GetLoginQrCodeApi(ctx *gin.Context) {
	queryKey, _ := ctx.GetQuery("key")
	/*	if !isExist {
		ctx.JSON(http.StatusOK,vo.NewFailUUId(""))
		return
	}*/
	result := service.GetLoginQrCodeService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

func GetLoginQrCodeNewApi(ctx *gin.Context) {
	reqModel := new(model.GetLoginQrCodeModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist {
		queryKey = guuid.New().String()
	}
	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.GetLoginQrCodeNewService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取验证码
func WxBindOpMobileForRegApi(ctx *gin.Context) {
	reqModel := new(model.WxBindOpMobileForModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist {
		queryKey = guuid.New().String()
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.WxBindOpMobileForRegService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 扫码登录新设备
func ExtDeviceLoginConfirmGetApi(ctx *gin.Context) {
	reqModel := new(model.ExtDeviceLoginConfirmModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist {
		queryKey = guuid.New().String()
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ExtDeviceLoginConfirmGetService(queryKey, reqModel.Url)
	ctx.JSON(http.StatusOK, result)
}

// 提取62
func Get62DataApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist {
		queryKey = guuid.New().String()
	}
	result := service.Get62DataService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 辅助新手机登录
func PhoneDeviceLoginApi(ctx *gin.Context) {
	reqModel := new(model.PhoneLoginModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.PhoneDeviceLoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// DeviceIdLoginApi  62账号密码登录
func DeviceIdLoginApi(ctx *gin.Context) {
	reqModel := new(model.DeviceIdLoginModel)
	queryKey, _ := ctx.GetQuery("key")
	/*	if !isExist {
		ctx.JSON(http.StatusOK,vo.NewFailUUId(""))
		return
	}*/
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.DeviceIdLoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 短信登录
func SmsLoginApi(ctx *gin.Context) {
	reqModel := new(model.DeviceIdLoginModel)
	queryKey, _ := ctx.GetQuery("key")
	if !validateData(ctx, &reqModel) {
		return
	}
	reqModel.Password = "strdm@," + reqModel.Password
	result := service.SmsLoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// a16数据号登录
func A16LoginApi(ctx *gin.Context) {
	reqModel := new(model.DeviceIdLoginModel)
	queryKey, _ := ctx.GetQuery("key")
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.A16LoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 62LoginNew新疆号登录
func LoginNewApi(ctx *gin.Context) {
	reqModel := new(model.DeviceIdLoginModel)
	queryKey, _ := ctx.GetQuery("key")
	if !validateData(ctx, &reqModel) {
		return
	}
	reqModel.Type = 1
	result := service.A16LoginService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// CheckLoginStatusApi 检测扫码状态
func CheckLoginStatusApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	uuid, _ := ctx.GetQuery("uuid")
	result := service.CheckLoginQrCodeStatusService(queryKey, uuid)
	ctx.JSON(http.StatusOK, result)
}

// 初始化状态
func GetInItStatusApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	result := service.GetInItStatusService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// WakeUpLoginApi 唤醒登录
func WakeUpLoginApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	result := service.WakeUpLoginService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// GetLoginStatusApi 获取在线状态
func GetLoginStatusApi(ctx *gin.Context) {
	autoLogin := false
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	isAutoLogin, isExist := ctx.GetQuery("autoLogin")
	if isExist {
		autoLogin = strings.Contains(isAutoLogin, "true")
	}

	result := service.GetLoginStatusService(queryKey, false, autoLogin)
	ctx.JSON(http.StatusOK, result)
}

// 获取设备list
func GetSafetyInfoApi(ctx *gin.Context) {
	queryKey, _ := ctx.GetQuery("key")
	result := service.GetSafetyInfoService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 删除设备
func DelSafeDeviceApi(ctx *gin.Context) {
	reqModel := new(model.DelSafeDeviceModel)
	queryKey, _ := ctx.GetQuery("key")
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.DelSafeDeviceService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 检测微信登录环境
func CheckCanSetAliasApi(ctx *gin.Context) {
	queryKey, _ := ctx.GetQuery("key")
	result := service.CheckCanSetAliasService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 打印链接数量
func IWXConnectMgrApi(ctx *gin.Context) {
	result := service.IWXConnectMgrService()
	ctx.JSON(http.StatusOK, result)
}
