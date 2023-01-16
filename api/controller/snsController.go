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

// 下载朋友圈视频
func DownloadMediaApi(ctx *gin.Context) {
	reqModel := new(model.DownloadMediaModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.DownloadMediaService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 设置朋友圈可见天数
func SetFriendCircleDaysApi(ctx *gin.Context) {
	reqModel := new(model.SetFriendCircleDaysModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetFriendCircleDaysService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendFriendCircleApi 发送朋友圈
func SendFriendCircleApi(ctx *gin.Context) {
	reqModel := new(model.SnsPostItemModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendFriendCircleService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendFriendCircleByXMlApi 发送朋友圈XML结构
func SendFriendCircleByXMlApi(ctx *gin.Context) {
	reqModel := new(baseinfo.TimelineObject)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendFriendCircleByXMlService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// UploadFriendCircleImageApi 上传图片信息
func UploadFriendCircleImageApi(ctx *gin.Context) {
	reqModel := new(model.UploadFriendCircleModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, " ") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.UploadFriendCircleImageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendSnsCommentRequestApi 点赞评论
func SendSnsCommentRequestApi(ctx *gin.Context) {
	reqModel := new(model.SendSnsCommentRequestModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendSnsCommentRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendSnsObjectOpRequestApi 朋友圈操作
func SendSnsObjectOpRequestApi(ctx *gin.Context) {
	reqModel := new(model.SendSnsObjectOpRequestModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendSnsObjectOpRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendSnsTimeLineRequestApi 获取朋友圈主页
func SendSnsTimeLineRequestApi(ctx *gin.Context) {
	reqModel := new(model.GetSnsInfoModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendSnsTimeLineRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendSnsUserPageRequestApi 获取指定人朋友圈
func SendSnsUserPageRequestApi(ctx *gin.Context) {
	reqModel := new(model.GetSnsInfoModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SendSnsUserPageRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendSnsObjectDetailByIdApi 获取指定id朋友圈
func SendSnsObjectDetailByIdApi(ctx *gin.Context) {
	reqModel := new(model.GetIdDetailModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendSnsObjectDetailByIdService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SetBackgroundImageApi设置朋友圈背景图片
func SetBackgroundImageApi(ctx *gin.Context) {
	reqModel := new(model.SetBackgroundImageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SetBackgroundImageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 转发收藏朋友圈id
func SendFavItemCircleApi(ctx *gin.Context) {
	reqModel := new(model.SendFavItemCircle)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendFavItemCircleService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 一键转发朋友圈
func SendOneIdCircleApi(ctx *gin.Context) {
	reqModel := new(model.GetIdDetailModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendOneIdCircleService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取收藏朋友圈详情
func GetCollectCircleApi(ctx *gin.Context) {
	reqModel := new(model.SendFavItemCircle)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetCollectCircleService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
