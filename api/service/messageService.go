package service

import (
	"encoding/base64"
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/container/garray"
	"strconv"
	"strings"
	"time"
)

// AddMessageMgrService 添加要发送的消息进入消息管理器
func AddMessageMgrService(queryKey string, m model.SendMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		connect := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		// 获取好友消息助手
		friendMsgMgr := connect.GetWXFriendMsgMgr()
		// 发送文本消息
		if len(m.MsgItem) <= 0 {
			return vo.NewFail("没有要加入消息管理器的消息！")
		}
		for _, item := range m.MsgItem {
			//1 text 2 Image
			if item.MsgType == 1 {
				friendMsgMgr.AddNewTextMsg(item.MsgIds, item.TextContent, item.ToUserName)
			} else if item.MsgType == 2 {
				sImageBase := strings.Split(item.ImageContent, ",")
				if len(sImageBase) > 1 {
					item.ImageContent = sImageBase[1]
				}
				imageBytes, _ := base64.StdEncoding.DecodeString(item.ImageContent)
				friendMsgMgr.AddImageMsg(item.MsgIds, imageBytes, item.ToUserName)
			}
		}
		return vo.NewSuccess(gin.H{}, "已加入消息管理器队列")
	})
}
func SendTestService(queryKey string, m model.SendMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		return vo.NewSuccessObj(nil, "")
	})
}

// SendImageMessageService 发送图片消息
func SendImageMessageService(queryKey string, m model.SendMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		// 发送文本消息
		if len(m.MsgItem) <= 0 {
			return vo.NewFail("没有要加入消息管理器的消息！")
		}
		results := garray.New(true)
		for _, item := range m.MsgItem {
			sImageBase := strings.Split(item.ImageContent, ",")
			if len(sImageBase) > 1 {
				item.ImageContent = sImageBase[1]
			}
			imageBytes, _ := base64.StdEncoding.DecodeString(item.ImageContent)
			imageId := baseutils.Md5ValueByte(imageBytes, false)

			cdnUploadImageResp, err := reqInvoker.SendCdnUploadImageReuqest(imageBytes, item.ToUserName)
			if err != nil {
				results.Append(gin.H{
					"imageId":       imageId,
					"toUSerName":    item.ToUserName,
					"isSendSuccess": cdnUploadImageResp,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				results.Append(gin.H{
					"imageId":       imageId,
					"toUSerName":    item.ToUserName,
					"isSendSuccess": cdnUploadImageResp,
				})
			}
		}
		return vo.NewSuccessObj(results, "")
	})
}

// 发送图片New
func SendImageNewMessageService(queryKey string, m model.SendMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		// 发送文本消息
		if len(m.MsgItem) <= 0 {
			return vo.NewFail("没有要加入消息管理器的消息！")
		}
		results := garray.New(true)
		for _, item := range m.MsgItem {
			sImageBase := strings.Split(item.ImageContent, ",")
			if len(sImageBase) > 1 {
				item.ImageContent = sImageBase[1]
			}
			imageBytes, _ := base64.StdEncoding.DecodeString(item.ImageContent)
			imageId := baseutils.Md5ValueByte(imageBytes, false)
			resp, err := reqInvoker.SendUploadImageNewRequest(imageBytes, item.ToUserName)
			if err != nil {
				results.Append(gin.H{
					"imageId":    imageId,
					"toUSerName": item.ToUserName,
					"resp":       resp,
					"errMsg":     err.Error(),
				})
				continue
			} else {
				results.Append(gin.H{
					"imageId":    imageId,
					"toUSerName": item.ToUserName,
					"resp":       resp,
				})
			}
		}
		return vo.NewSuccessObj(results, "")
	})
}

// SendTextMessageService 发送文本消息
func SendTextMessageService(queryKey string, m model.SendMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		// 发送文本消息
		if len(m.MsgItem) <= 0 {
			return vo.NewFail("没有要加入消息管理器的消息！")
		}
		results := garray.New(true)
		for _, item := range m.MsgItem {
			resp, err := bizcgi.SendTextMsgReq(wxAccount, item.ToUserName, item.TextContent, item.AtWxIDList, item.MsgType)
			if err != nil {
				results.Append(gin.H{
					"toUSerName":    item.ToUserName,
					"textContent":   item.TextContent,
					"isSendSuccess": false,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				results.Append(gin.H{
					"textContent":   item.TextContent,
					"toUSerName":    item.ToUserName,
					"resp":          resp,
					"isSendSuccess": true,
				})
			}
		}
		return vo.NewSuccessObj(results, "")
	})
}

// SendTextMessageService 分享名片
func SendShareCardService(queryKey string, m model.SendShareCardModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		resp, err := bizcgi.SendShareCardReq(wxAccount, m.ToUserName, m.Id, m.Nickname, m.Alias)
		var result gin.H
		if err != nil {
			result = gin.H{
				"toUserName":    m.ToUserName,
				"isSendSuccess": false,
				"errMsg":        err.Error(),
			}
		} else {
			result = gin.H{
				"toUserName":    m.ToUserName,
				"resp":          resp,
				"isSendSuccess": true,
			}
		}

		return vo.NewSuccessObj(result, "")
	})
}

// ForwardImageMessageService 转发图片
func ForwardImageMessageService(queryKey string, m model.ForwardMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		results := garray.New(true)

		if len(m.ForwardImageList) <= 0 {
			return vo.NewFail("没有要进行转发的数据。")
		}
		for i, item := range m.ForwardImageList {
			var cdnItem baseinfo.ForwardImageItem
			StructCopy(&cdnItem, &item)
			resp, err := reqInvoker.ForwardCdnImageRequest(cdnItem)
			if err != nil {
				results.Append(gin.H{
					"toUSerName":    item.ToUserName,
					"cdnMidImgUrl":  item.CdnMidImgUrl,
					"isSendSuccess": false,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				isSendSuccess := false
				if resp.GetBaseResponse().GetRet() == 0 {
					isSendSuccess = true
				}
				results.Append(gin.H{
					"cdnMidImgUrl":  item.CdnMidImgUrl,
					"toUSerName":    item.ToUserName,
					"isSendSuccess": isSendSuccess,
					"resp":          resp,
					"retCode":       resp.GetBaseResponse().GetRet(),
					"errMsg":        resp.GetBaseResponse().GetErrMsg().GetStr(),
				})
			}
			//每两天延迟1秒
			if i != 0 && i%2 == 0 {
				time.Sleep(time.Second * 1)
			}
		}
		return vo.NewSuccessObj(results.Interfaces(), "")
	})
}

// ForwardVideoMessageService 转发视频
func ForwardVideoMessageService(queryKey string, m model.ForwardMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		results := garray.New(true)

		if len(m.ForwardVideoList) <= 0 {
			return vo.NewFail("没有要进行转发的数据。")
		}
		for i, item := range m.ForwardVideoList {
			var cdnItem baseinfo.ForwardVideoItem
			StructCopy(&cdnItem, &item)
			resp, err := reqInvoker.ForwardCdnVideoRequest(cdnItem)
			if err != nil {
				results.Append(gin.H{
					"toUSerName":    item.ToUserName,
					"CdnVideoUrl":   item.CdnVideoUrl,
					"isSendSuccess": false,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				isSendSuccess := false
				if resp.GetBaseResponse().GetRet() == 0 {
					isSendSuccess = true
				}
				results.Append(gin.H{
					"CdnVideoUrl":   item.CdnVideoUrl,
					"toUSerName":    item.ToUserName,
					"isSendSuccess": isSendSuccess,
					"resp":          resp,
					"retCode":       resp.GetBaseResponse().GetRet(),
					"errMsg":        resp.GetBaseResponse().GetErrMsg().GetStr(),
				})
			}
			//每两天延迟1秒
			if i != 0 && i%2 == 0 {
				time.Sleep(time.Second * 1)
			}
		}
		return vo.NewSuccessObj(results.Interfaces(), "")
	})
}

// SendAppMessageService 发送app消息
func SendAppMessageService(queryKey string, m model.AppMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		results := garray.New(true)
		if len(m.AppList) == 0 {
			return vo.NewFail("没有要进行发送的数据。")
		}
		for _, item := range m.AppList {
			resp, err := reqInvoker.SendAppMessage(item.ContentXML, item.ToUserName, item.ContentType)
			if err != nil {
				results.Append(gin.H{
					"contentXML":    item.ContentXML,
					"toUserName":    item.ToUserName,
					"contentType":   item.ContentType,
					"isSendSuccess": false,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				isSendSuccess, retCode := false, resp.GetBaseResponse().GetRet()
				if resp.GetBaseResponse().GetRet() == 0 {
					isSendSuccess = true
				}
				results.Append(gin.H{
					"contentXML":    item.ContentXML,
					"toUserName":    item.ToUserName,
					"contentType":   item.ContentType,
					"isSendSuccess": isSendSuccess,
					"resp":          resp,
					"retCode":       retCode,
				})

			}
		}
		return vo.NewSuccessObj(results.Interfaces(), "")
	})
}

// SendEmojiMessageService 发送表情
func SendEmojiMessageService(queryKey string, m model.SendEmojiMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		results := garray.New(true)

		if len(m.EmojiList) <= 0 {
			return vo.NewFail("没有要进行转发的数据。")
		}
		for i, item := range m.EmojiList {
			resp, err := reqInvoker.SendEmojiRequest(item.EmojiMd5, item.ToUserName, item.EmojiSize)
			if err != nil {
				results.Append(gin.H{
					"toUSerName":    item.ToUserName,
					"EmojiMd5":      item.EmojiMd5,
					"isSendSuccess": false,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				isSendSuccess := false
				if resp.GetBaseResponse().GetRet() == 0 {
					isSendSuccess = true
				}
				results.Append(gin.H{
					"EmojiMd5":      item.EmojiMd5,
					"toUSerName":    item.ToUserName,
					"isSendSuccess": isSendSuccess,
					"retCode":       resp.GetBaseResponse().GetRet(),
					"errMsg":        resp.GetBaseResponse().GetErrMsg().GetStr(),
				})
			}
			//每两天延迟1秒
			if i != 0 && i%2 == 0 {
				time.Sleep(time.Second * 1)
			}
		}
		return vo.NewSuccessObj(results.Interfaces(), "")

	})
}

// ForwardEmojiService 发送表情&包含动图
func ForwardEmojiService(queryKey string, m model.SendEmojiMessageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()

		results := garray.New(true)

		if len(m.EmojiList) <= 0 {
			return vo.NewFail("没有要进行转发的数据。")
		}
		for i, item := range m.EmojiList {
			resp, err := reqInvoker.ForwardEmojiRequest(item.EmojiMd5, item.ToUserName, item.EmojiSize)
			if err != nil {
				results.Append(gin.H{
					"toUSerName":    item.ToUserName,
					"EmojiMd5":      item.EmojiMd5,
					"isSendSuccess": false,
					"errMsg":        err.Error(),
				})
				continue
			} else {
				isSendSuccess := false
				if resp.GetBaseResponse().GetRet() == 0 {
					isSendSuccess = true
				}
				results.Append(gin.H{
					"EmojiMd5":      item.EmojiMd5,
					"toUSerName":    item.ToUserName,
					"isSendSuccess": isSendSuccess,
					"retCode":       resp.GetBaseResponse().GetRet(),
					"errMsg":        resp.GetBaseResponse().GetErrMsg().GetStr(),
				})
			}
			//每两天延迟1秒
			if i != 0 && i%2 == 0 {
				time.Sleep(time.Second * 1)
			}
		}
		return vo.NewSuccessObj(results.Interfaces(), "")
	})
}

// RevokeMsgService 撤销消息
func RevokeMsgService(queryKey string, m model.RevokeMsgModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendRevokeMsgRequest(m.NewMsgId, m.ClientMsgId, m.ToUserName)
		if err != nil {
			return vo.NewFail("RevokeMsgService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// UploadVoiceRequestService 发送语音
func UploadVoiceRequestService(queryKey string, m model.SendUploadVoiceRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		/*voiceData, err := base64.StdEncoding.DecodeString(m.VoiceData)
		if err != nil {
			return vo.NewFail("UploadVoiceRequestService base64.StdEncoding.DecodeString err" + err.Error())
		}*/
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendUploadVoiceRequest(m.ToUserName, m.VoiceData, m.VoiceSecond, m.VoiceFormat)
		if err != nil {
			return vo.NewFail("UploadVoiceRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// UploadVoiceRequestService 发送视频
func SendCdnUploadVideoRequestService(queryKey string, m model.CdnUploadVideoRequest) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendCdnUploadVideoRequest(m.ToUserName, m.ThumbData, m.VideoData)
		if err != nil {
			return vo.NewFail("UploadVoiceRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 下载请求
func SendCdnDownloadService(queryKey string, m model.DownMediaModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		item := baseinfo.DownMediaItem{
			AesKey:   m.AesKey,
			FileURL:  m.FileURL,
			FileType: m.FileType,
		}
		_, errorCdn := reqInvoker.SendGetCDNDnsRequest()
		if errorCdn != nil {
			return vo.NewFail("SendGetCDNDnsRequest err:" + errorCdn.Error())
		}
		resp, err := reqInvoker.SendCdnDownloadReuqest(&item)
		if err != nil {
			return vo.NewFail("UploadVoiceRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 获取图片请求
func GetMsgBigImgService(queryKey string, m model.GetMsgBigImgModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.GetMsgBigImg(m)
		if err != nil {
			return vo.NewFail("UploadVoiceRequestService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 群发文字
func GroupMassMsgTextService(queryKey string, m model.GroupMassMsgTextModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGroupMassMsgTextRequest(m.ToUserName, m.Content)
		if err != nil {
			return vo.NewFail("GroupMassMsgTextService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 群发图片
func GroupMassMsgImageService(queryKey string, m model.GroupMassMsgImageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		imageBytes, _ := base64.StdEncoding.DecodeString(m.ImageBase64)
		resp, err := reqInvoker.SendGroupMassMsgImageRequest(m.ToUserName, imageBytes)
		if err != nil {
			return vo.NewFail("GroupMassMsgImageService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 群拍一拍
func SendPatService(queryKey string, m model.SendPatModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendSendPatRequest(m.ChatRoomName, m.ToUserName, m.Scene)
		if err != nil {
			return vo.NewFail("GroupMassMsgTextService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 下载语音
func GetMsgVoiceService(queryKey string, m model.DownloadVoiceModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGetMsgVoiceRequest(m.ToUserName, m.NewMsgId, m.Bufid, m.Length)
		if err != nil {
			return vo.NewFail("GetMsgVoiceService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 同步历史消息
func NewSyncHistoryMessageService(queryKey string, m model.SyncModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		NewSyncResponse, err := bizcgi.NewSyncHistoryMessageRequest(queryKey, wxAccount, m.Scene, m.SyncKey)
		if err != nil {
			return vo.NewFail("NewSyncHistoryMessageService err:" + err.Error())
		}
		if NewSyncResponse.Key != nil {
			wxAccount.GetUserInfo().SyncHistoryKey = NewSyncResponse.Key.Buffer
		}
		return vo.NewSuccessObj(NewSyncResponse, "成功")
	})
}

// 同步历史群消息
func NewSyncGroupMessageService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendWXSyncContactRequest()
		if err != nil {
			return vo.NewFail("NewSyncHistoryMessageService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "成功")
	})
}
