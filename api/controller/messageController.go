package controller

import (
	req "feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AddTextMessageApi 添加要发送的文本消息进入管理器
func AddMessageMgrApi(ctx *gin.Context) {
	reqModel := new(req.SendMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.AddMessageMgrService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendImageMessageApi 发送图片消息
func SendImageMessageApi(ctx *gin.Context) {
	reqModel := new(req.SendMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendImageMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 发送图片消息New
func SendImageNewMessageApi(ctx *gin.Context) {
	reqModel := new(req.SendMessageModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Keys
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendImageNewMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

func TestApi(ctx *gin.Context) {
	reqModel := new(req.SendMessageModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendTestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendTextMessageApi 发送文本消息
func SendTextMessageApi(ctx *gin.Context) {
	reqModel := new(req.SendMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendTextMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendShareCardApi 分享名片
func SendShareCardApi(ctx *gin.Context) {
	reqModel := new(req.SendShareCardModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendShareCardService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ForwardImageMessageApi 转发图片
func ForwardImageMessageApi(ctx *gin.Context) {
	reqModel := new(req.ForwardMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ForwardImageMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// ForwardVideoMessageApi 转发视频
func ForwardVideoMessageApi(ctx *gin.Context) {
	reqModel := new(req.ForwardMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ForwardVideoMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendEmojiMessageApi 发送表情
func SendEmojiMessageApi(ctx *gin.Context) {
	reqModel := new(req.SendEmojiMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendEmojiMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 转发表情，包含动图
func ForwardEmojiApi(ctx *gin.Context) {
	reqModel := new(req.SendEmojiMessageModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.ForwardEmojiService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SendAppMessageApi 发送App消息
func SendAppMessageApi(ctx *gin.Context) {
	reqModel := new(req.AppMessageModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendAppMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// RevokeMsgApi 撤销消息
func RevokeMsgApi(ctx *gin.Context) {
	reqModel := new(req.RevokeMsgModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.RevokeMsgService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// UploadVoiceRequestApi 发送语音
func UploadVoiceRequestApi(ctx *gin.Context) {
	reqModel := new(req.SendUploadVoiceRequestModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.UploadVoiceRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// CdnUploadVideoRequestApi 上传视频
func CdnUploadVideoRequestApi(ctx *gin.Context) {
	reqModel := new(req.CdnUploadVideoRequest)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendCdnUploadVideoRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 下载 请求
func SendCdnDownloadApi(ctx *gin.Context) {
	reqModel := new(req.DownMediaModel)

	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.SendCdnDownloadService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取图片
func GetMsgBigImgApi(ctx *gin.Context) {
	reqModel := new(req.GetMsgBigImgModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetMsgBigImgService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 同步历史消息
func NewSyncHistoryMessageApi(ctx *gin.Context) {
	reqModel := new(req.SyncModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.NewSyncHistoryMessageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 群发接口
func GroupMassMsgTextApi(ctx *gin.Context) {
	reqModel := new(req.GroupMassMsgTextModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GroupMassMsgTextService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 群发图片
func GroupMassMsgImageApi(ctx *gin.Context) {
	reqModel := new(req.GroupMassMsgImageModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GroupMassMsgImageService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 下载语音消息
func GetMsgVoiceApi(ctx *gin.Context) {
	reqModel := new(req.DownloadVoiceModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetMsgVoiceService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
