package service

import (
	"encoding/base64"
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/container/garray"
	_ "github.com/gogf/gf/container/garray"
	"strconv"
	"strings"
	"sync"
)

// SendSnsTimeLineRequestService 获取朋友圈主页
func SendSnsTimeLineRequestService(queryKey string, m model.GetSnsInfoModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		//获取朋友圈主页
		resp, err := reqInvoker.SendSnsTimeLineRequestResult(m.FirstPageMD5, m.MaxID)
		if err != nil {
			return vo.NewFail("SendSnsTimeLineRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 通过id获朋友圈详情
func SendSnsObjectDetailByIdService(queryKey string, m model.GetIdDetailModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		//获取好友朋友圈主页
		id, _ := strconv.ParseUint(m.Id, 0, 64)
		resp, err := reqInvoker.SendSnsObjectDetailRequest(id)
		if err != nil {
			return vo.NewFail("SendSnsTimeLineRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// SendSnsUserPageRequestService 获取指定人朋友圈
func SendSnsUserPageRequestService(queryKey string, m model.GetSnsInfoModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		//获取好友朋友圈主页
		resp, err := reqInvoker.SendSnsUserPageRequest(m.UserName, m.FirstPageMD5, m.MaxID, true)
		if err != nil {
			return vo.NewFail("SendSnsTimeLineRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 转发收藏朋友圈
func SendFavItemCircleService(queryKey string, m model.SendFavItemCircle) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 获取指定的朋友圈
		objIDString := baseutils.GetNumberString(m.SourceID)
		snsObjID, _ := strconv.ParseUint(objIDString, 10, 64)
		_, err := reqInvoker.SendSnsObjectDetailRequest(snsObjID)
		if err != nil {
			baseutils.PrintLog("WXSnsTransTask.doFavTask - SendSnsObjectDetailRequest err: " + err.Error())
			return vo.NewFail("操作失败！" + err.Error())
		}
		//放入消息队例
		//currentTaskMgr := iwxConnect.GetWXTaskMgr()
		//taskMgr, _ := currentTaskMgr.(*wxmgr.WXTaskMgr)
		//currentSnsTransTask := taskMgr.GetSnsTransTask()
		//// 转发朋友圈
		//err = currentSnsTransTask.DoSnsTransTask(snsObject, defines.MTaskTypeFavTrans, m.BlackList, m.Location, m.LocationVal)
		if err == nil {
			// 如果转发收藏成功则删除
			reqInvoker.SendBatchDelFavItemRequest(m.FavItemID)
			return vo.NewSuccessObj("ok", "操作成功！")
		}
		return vo.NewFail("操作失败！" + err.Error())
	})
}

// 一键转发朋友圈
func SendOneIdCircleService(queryKey string, m model.GetIdDetailModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 获取指定的朋友圈
		snsObjID, _ := strconv.ParseUint(m.Id, 10, 64)
		_, err := reqInvoker.SendSnsObjectDetailRequest(snsObjID)
		if err != nil {
			baseutils.PrintLog("WXSnsTransTask.doFavTask - SendSnsObjectDetailRequest err: " + err.Error())
			return vo.NewFail("操作失败！" + err.Error())
		}
		//放入消息队例
		//currentTaskMgr := iwxConnect.GetWXTaskMgr()
		//taskMgr, _ := currentTaskMgr.(*wxmgr.WXTaskMgr)
		//currentSnsTransTask := taskMgr.GetSnsTransTask()
		//// 转发朋友圈
		//err = currentSnsTransTask.DoSnsTransTask(snsObject, defines.MTaskTypeFavTrans, m.BlackList, m.Location, m.LocationVal)
		if err == nil {
			// 如果转发收藏成功则删除
			return vo.NewSuccessObj("ok", "转发成功！")
		}
		return vo.NewFail("操作失败！" + err.Error())
	})
}

// 获取收藏朋友圈详情
func GetCollectCircleService(queryKey string, m model.SendFavItemCircle) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		// 获取指定的朋友圈
		objIDString := baseutils.GetNumberString(m.SourceID)
		snsObjID, _ := strconv.ParseUint(objIDString, 10, 64)
		snsObject, err := reqInvoker.SendSnsObjectDetailRequest(snsObjID)
		if err != nil {
			baseutils.PrintLog("WXSnsTransTask.doFavTask - SendSnsObjectDetailRequest err: " + err.Error())
			return vo.NewFail("操作失败！" + err.Error())
		}
		return vo.NewSuccessObj(snsObject, "操作成功!")
	})
}

// SetBackgroundImageApi 设置朋友圈图片
func SetBackgroundImageService(queryKey string, m model.SetBackgroundImageModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//习近平
		content := fmt.Sprintf("<TimelineObject><id><![CDATA[0]]></id><username><![CDATA[%s]]></username><createTime><![CDATA[0]]></createTime><contentDescShowType>0</contentDescShowType><contentDescScene>0</contentDescScene><private><![CDATA[0]]></private><contentDesc></contentDesc><contentattr><![CDATA[0]]></contentattr><sourceUserName></sourceUserName><sourceNickName></sourceNickName><statisticsData></statisticsData><weappInfo><appUserName></appUserName><pagePath></pagePath></weappInfo><canvasInfoXml></canvasInfoXml><location poiClickableStatus=\"0\"  poiClassifyId=\"\"  poiScale=\"0\"  longitude=\"0.0\"  city=\"\"  poiName=\"\"  latitude=\"0.0\"  poiClassifyType=\"0\"  poiAddress=\"\" ></location><ContentObject><contentStyle><![CDATA[7]]></contentStyle><contentSubStyle><![CDATA[0]]></contentSubStyle><title>&#x0A;&#x0A;&#x0A;</title><description></description><contentUrl></contentUrl><mediaList><media><id><![CDATA[0]]></id><type><![CDATA[2]]></type><title></title><description></description><private><![CDATA[0]]></private><url type=\"1\" ><![CDATA[%s]]></url><thumb type=\"1\" ><![CDATA[%s]]></thumb>", wxAccount.GetUserInfo().WxId, m.Url, m.Url)
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		snsPostItem := baseinfo.SnsPostItem{
			Content:   content,
			Xml:       true,
			MediaList: make([]*baseinfo.SnsMediaItem, 0),
		}
		resp, err := reqInvoker.SendSnsPostRequestNew(&snsPostItem)
		if err != nil {
			return vo.NewFail("操作失败！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "操作成功！")
	})
}

// 下载视频
func DownloadMediaService(queryKey string, req model.DownloadMediaModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		tmpEncKey, _ := strconv.Atoi(req.Key)
		videoData, err := reqInvoker.SendCdnSnsVideoDownloadReuqest(uint64(tmpEncKey), req.URL)
		if err != nil {
			return vo.NewFail("下载失败！" + err.Error())
		}
		return vo.NewSuccessObj(videoData, "下载成功！")
	})
}

// 设置朋友圈可见天数
func SetFriendCircleDaysService(queryKey string, postItem model.SetFriendCircleDaysModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		err := reqInvoker.SetFriendCircleDays(&postItem)
		if err != nil {
			return vo.NewFail("设置朋友圈可见天数失败！" + err.Error())
		}
		return vo.NewSuccessObj("ok", "设置朋友圈可见天数成功！")
	})
}

// SendFriendCircle发送朋友圈
func SendFriendCircleService(queryKey string, postItem model.SnsPostItemModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		snsPostItem := baseinfo.SnsPostItem{
			MediaList: make([]*baseinfo.SnsMediaItem, 0),
		}

		//对结构体进行复制
		StructCopy(&snsPostItem, &postItem)
		if postItem.LocationInfo != nil {
			snsPostItem.LocationInfo = (*baseinfo.SnsLocationInfo)(postItem.LocationInfo)
		}
		if len(postItem.MediaList) > 0 {
			for _, item := range postItem.MediaList {
				var snsMediaItem baseinfo.SnsMediaItem
				StructCopy(&snsMediaItem, item)
				snsPostItem.MediaList = append(snsPostItem.MediaList, &snsMediaItem)
			}

		}
		resp, err := reqInvoker.SendSnsPostRequestNew(&snsPostItem)
		if err != nil {
			return vo.NewFail("发送朋友圈失败！" + err.Error())
		}

		return vo.NewSuccessObj(resp, "发送朋友圈成功！")

	})
}

// SendFriendCircleByXMlService 根据XML发送朋友圈
func SendFriendCircleByXMlService(queryKey string, postItem baseinfo.TimelineObject) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		err := reqInvoker.SendSnsPostRequestByXML(&postItem, []string{})
		if err != nil {
			return vo.NewFail("发送朋友圈失败！" + err.Error())
		}

		return vo.NewSuccessObj(nil, "发送朋友圈成功！")

	})
}

// UploadFriendCircleImagesService 上传朋友圈图片
func UploadFriendCircleImageService(queryKey string, m model.UploadFriendCircleModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		uploadRespArray := garray.New(true)
		wg := new(sync.WaitGroup)

		if len(m.ImageDataList) <= 0 {
			return vo.NewFail("没有要上传的图片！")
		}

		for _, imageData := range m.ImageDataList {
			sImageBase := strings.Split(imageData, ",")
			if len(sImageBase) > 1 {
				imageData = sImageBase[1]
			}
			imageBuffer, _ := base64.StdEncoding.DecodeString(imageData)
			//生成一个Md5
			imageId := baseutils.Md5ValueByte(imageBuffer, false)
			//log.Println(imageId)
			//查询数据库是否存在该Id
			//上传图片
			wg.Add(1)
			go func(image []byte, id string) {
				defer wg.Done()
				upImageResp, err := reqInvoker.SendCdnSnsUploadImageReuqest(image)
				if err != nil {
					uploadRespArray.Append(gin.H{
						"imageId": imageId,
						"errMgs":  err.Error(),
					})
				} else {
					uploadRespArray.Append(gin.H{
						"imageId": imageId,
						"resp":    upImageResp,
					})
				}
			}(imageBuffer, imageId)
		}
		wg.Wait()
		return vo.NewSuccessObj(uploadRespArray.Interfaces(), "")
	})
}

// SendSnsObjectOpRequestService 朋友圈操作
func SendSnsObjectOpRequestService(queryKey string, m model.SendSnsObjectOpRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}

		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		if len(m.SnsObjectOpList) <= 0 {
			return vo.NewFail("没有要操作的Id！")
		}

		opItems := make([]*baseinfo.SnsObjectOpItem, 0)
		for _, item := range m.SnsObjectOpList {
			var opItem baseinfo.SnsObjectOpItem
			StructCopy(&opItem, &item)
			opItems = append(opItems, &opItem)

		}
		resp, err := reqInvoker.SendSnsObjectOpRequest(opItems)
		if err != nil {
			return vo.NewFail(err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// SendSnsCommentRequestService 点赞/评论
func SendSnsCommentRequestService(queryKey string, m model.SendSnsCommentRequestModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()

		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()

		respArray := garray.New(true)
		wg := new(sync.WaitGroup)

		if len(m.SnsCommentList) <= 0 {
			return vo.NewFail("没有数据！")
		}
		for _, item := range m.SnsCommentList {
			commItem := &baseinfo.SnsCommentItem{}
			if item.OpType == baseinfo.MMSnsCommentTypeLike { //点赞
				id, _ := strconv.ParseUint(item.ItemID, 0, 64)
				commItem = clientsdk.CreateSnsCommentLikeItem(id, item.ToUserName)
			} else if item.OpType == baseinfo.MMSnsCommentTypeComment { // 评论
				//oldDevieType:=wxAccount.GetUserInfo().DeviceInfo.OsType
				//wxAccount.GetUserInfo().DeviceInfo.OsType="wechat"
				//_=reqInvoker.SendAutoAuthRequest()
				/*defer func() {
					wxAccount.GetUserInfo().DeviceInfo.OsType=oldDevieType
					_=reqInvoker.SendAutoAuthRequest()
				}()*/
				//评论能用
				if m.Tx {
					item.Content = item.Content + "\n\n~~~~~~~~~~\n每日一练习,\n大大有进步"
				}
				id, _ := strconv.ParseUint(item.ItemID, 0, 64)
				commItem = clientsdk.CreateSnsCommentItem(id, item.ToUserName, item.Content, nil)
				commItem.ReplyItem = &baseinfo.ReplyCommentItem{}
				StructCopy(commItem.ReplyItem, &item.ReplyItem)
			}

			wg.Add(1)
			go func(commentItem *baseinfo.SnsCommentItem) {
				defer wg.Done()
				err := reqInvoker.SendSnsCommentRequest(commentItem)
				if err != nil {
					respArray.Append(gin.H{
						"data":             commentItem,
						"isCommentSuccess": false,
						"errMsg":           err.Error(),
					})
				} else {
					respArray.Append(gin.H{
						"data":             commentItem,
						"isCommentSuccess": true,
						"errMsg":           "",
					})
				}

			}(commItem)
		}
		wg.Wait()
		return vo.NewSuccessObj(respArray.Interfaces(), "")

	})
}
