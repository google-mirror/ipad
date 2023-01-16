package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// LogOutRequestApi 退出登录
func LogOutRequestApi(ctx *gin.Context) {

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.LogOutService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 获取缓存在redis中的消息
func GetRedisSyncMsgApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	reqModel := new(model.GetSyncMsgModel)
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetRedisSyncMsgService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// GetOnlineInfoApi 获取在线设备信息
func GetOnlineInfoApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.OnlineInfoService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// GetProfileApi 获取个人资料信息
func GetProfileApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.GetProfileService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// DelContactApi 删除好友
func DelContactApi(ctx *gin.Context) {
	reqModel := new(model.DelContactModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendDelContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ModifyUserInfoRequestApi 修改资料
func ModifyUserInfoRequestApi(ctx *gin.Context) {
	reqModel := new(model.ModifyUserInfo)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendModifyUserInfoRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改名称
func UpdateNickNameApi(ctx *gin.Context) {
	reqModel := new(model.UpdateNickNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.UpdateNickNameService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改名称
func SetNickNameApi(ctx *gin.Context) {
	reqModel := new(model.UpdateNickNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	reqModel.Scene = 1
	result := service.SetNickNameService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改性别
func SetSexApi(ctx *gin.Context) {
	reqModel := new(model.UpdateSexModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetSexService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改签名
func SetSignatureApi(ctx *gin.Context) {
	reqModel := new(model.UpdateNickNameModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	reqModel.Scene = 2
	result := service.SetNickNameService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ChangePwdRequestRequestApi 更改密码
func ChangePwdRequestRequestApi(ctx *gin.Context) {
	reqModel := new(model.SendChangePwdRequestModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendChangePwdRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 上传头像
func UploadHeadImageApi(ctx *gin.Context) {
	reqModel := new(model.UploadHeadImageModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.UploadHeadImageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

func UpdateAutopassApi(ctx *gin.Context) {
	reqModel := new(model.UpdateAutopassModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.UpdateAutopassService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendModifyRemarkRequestApi 修改备注
func SendModifyRemarkRequestApi(ctx *gin.Context) {
	reqModel := new(model.SendModifyRemarkRequestModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendModifyRemarkRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 设置微信号
func SetWechatApi(ctx *gin.Context) {
	reqModel := new(model.AlisaModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetWechatService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 修改步数
func UpdateStepNumberApi(ctx *gin.Context) {
	reqModel := new(model.UpdateStepNumberModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.UpdateStepNumberService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取步数据列表
func GetUserRankLikeCountApi(ctx *gin.Context) {
	reqModel := new(model.UserRankLikeModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetUserRankLikeCountService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

/*//修改加好友需要验证属性
func UpdateAutopassApi(ctx *gin.Context) {
	reqModel := new(req.UpdateAutopassModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.UpdateAutopassService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
*/
//设置添加我的方式
func SetFunctionSwitchApi(ctx *gin.Context) {
	reqModel := new(model.WxFunctionSwitchModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetFunctionSwitchService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 设置拍一拍名称
func SetSendPatApi(ctx *gin.Context) {
	reqModel := new(model.SetSendPatModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetSendPatService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 换绑手机号
func BindingMobileApi(ctx *gin.Context) {
	reqModel := new(model.BindMobileModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.BindingMobileService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 发送手机验证码,
func SendVerifyMobileApi(ctx *gin.Context) {
	reqModel := new(model.SendVerifyMobileModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendVerifyMobileService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
