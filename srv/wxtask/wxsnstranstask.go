package wxtask

import (
	"encoding/json"
	"github.com/lunny/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"

	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/xmltool"

	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/defines"
	"feiyu.com/wx/srv/wxface"
)

// WXSnsTransTask 朋友圈转发管理器
type WXSnsTransTask struct {
	wxConn wxface.IWXConnect
	// 待转发的收藏
	favItemChan chan *baseinfo.FavItem
	// 同步转发
	syncItemChan chan *wechat.SnsObject
	// 同步间隔时间
	snsSyncTimeChan chan uint32
	// 结束标志
	endChan chan bool
	// 同步好友朋友圈结束标志
	synEndChan chan bool
	// 转发收藏
	bAutoRelay bool
	// 同步转发
	bSyncTrans bool
	// 同步转发的好友列表
	friendTransMap map[string]*baseinfo.FriendTransItem
	// 同步锁
	wxCurrentLock sync.RWMutex
}

// NewWXSnsTransTask 新建收藏转发任务器
func NewWXSnsTransTask(wxConn wxface.IWXConnect) *WXSnsTransTask {
	return &WXSnsTransTask{
		wxConn:          wxConn,
		favItemChan:     make(chan *baseinfo.FavItem, 100),
		syncItemChan:    make(chan *wechat.SnsObject, 100),
		friendTransMap:  make(map[string]*baseinfo.FriendTransItem),
		snsSyncTimeChan: make(chan uint32, 5),
		endChan:         make(chan bool, 1),
		synEndChan:      make(chan bool, 1),
		bAutoRelay:      false,
		bSyncTrans:      false,
	}
}

// Start 开启任务
func (transTask *WXSnsTransTask) Start() {
	// 开启任务
	go transTask.startTransTask()
	go transTask.startSnsSyncFriend()
}

// Stop 关闭任务
func (transTask *WXSnsTransTask) Stop() {
	transTask.endChan <- true
	//transTask.synEndChan <- true
}

// SetAutoRelay 设置自动转发
func (transTask *WXSnsTransTask) SetAutoRelay(bFlag bool) {
	transTask.bAutoRelay = bFlag
}

// IsAutoRelay 是否自动转发
func (transTask *WXSnsTransTask) IsAutoRelay() bool {
	return transTask.bAutoRelay
}

// SetSyncTrans 设置同步转发
func (transTask *WXSnsTransTask) SetSyncTrans(bFlag bool) {
	transTask.bSyncTrans = bFlag
	if !bFlag {
		transTask.ClearFriendTransMap()
	}
}

// IsSyncTrans 是否同步转发
func (transTask *WXSnsTransTask) IsSyncTrans() bool {
	return transTask.bSyncTrans
}

// AddSyncTransFriend 增加同步转发的好友
func (transTask *WXSnsTransTask) AddSyncTransFriend(friendWXID string) bool {
	transTask.wxCurrentLock.Lock()
	defer transTask.wxCurrentLock.Unlock()

	_, ok := transTask.friendTransMap[friendWXID]
	if !ok {
		// 不存在就添加
		friendTransItem := &baseinfo.FriendTransItem{}
		friendTransItem.FriendWXID = friendWXID
		friendTransItem.FirstPageMd5 = ""
		friendTransItem.CreateTime = 0
		friendTransItem.IsInited = false
		transTask.friendTransMap[friendWXID] = friendTransItem
		return true
	}
	return false
}

// GetTransFriendItem 根据微信ID 同步转发好友项
func (transTask *WXSnsTransTask) GetTransFriendItem(friendWXID string) *baseinfo.FriendTransItem {
	transTask.wxCurrentLock.Lock()
	defer transTask.wxCurrentLock.Unlock()

	friendTransItem, ok := transTask.friendTransMap[friendWXID]
	if !ok {
		return nil
	}
	return friendTransItem
}

// GetTransFriendCount 返回个数
func (transTask *WXSnsTransTask) GetTransFriendCount() int {
	transTask.wxCurrentLock.Lock()
	defer transTask.wxCurrentLock.Unlock()

	return len(transTask.friendTransMap)
}

// ClearFriendTransMap 清空同步转发的好友列表
func (transTask *WXSnsTransTask) ClearFriendTransMap() {
	transTask.wxCurrentLock.Lock()
	defer transTask.wxCurrentLock.Unlock()

	// 清空所有元素
	transTask.friendTransMap = make(map[string]*baseinfo.FriendTransItem)
}

// AddFavItem 新增收藏项
func (transTask *WXSnsTransTask) AddFavItem(favItem *baseinfo.FavItem) {
	transTask.favItemChan <- favItem
}

// AddSyncItem 新增同步项
func (transTask *WXSnsTransTask) AddSyncItem(snsObject *wechat.SnsObject) {
	transTask.syncItemChan <- snsObject
}

// startSnsSyncFriend 同步好友朋友圈
func (transTask *WXSnsTransTask) startSnsSyncFriend() {
	transTask.snsSyncTimeChan <- 60
	for {
		select {
		case waitTimes := <-transTask.snsSyncTimeChan:
			time.Sleep(time.Second * time.Duration(waitTimes))
			// 如果开启了自动转发
			if transTask.bSyncTrans {
				// 遍历看是否又跟新朋友圈
				transTask.checkSnsFriendList()
			}
			transTask.snsSyncTimeChan <- 60
			continue
		case <-transTask.synEndChan:
			return
		}
	}
}

func (transTask *WXSnsTransTask) checkSnsFriendList() {
	transTask.wxCurrentLock.Lock()
	defer transTask.wxCurrentLock.Unlock()
	// 遍历朋友圈
	currentReqInvoker := transTask.wxConn.GetWXReqInvoker()
	for _, friendItem := range transTask.friendTransMap {
		_, err := currentReqInvoker.SendSnsUserPageRequest(friendItem.FriendWXID, "", 0, false)
		if err != nil {
			break
		}
	}
}

// starttransTask 任务线程
func (transTask *WXSnsTransTask) startTransTask() {
	for {
		select {
		case currentFavItem := <-transTask.favItemChan:
			err := transTask.doFavTask(currentFavItem)
			if err != nil {
				baseutils.PrintLog("转发收藏失败")
			}
			continue
		case snsObject := <-transTask.syncItemChan:
			location := baseinfo.Location{}
			err := transTask.DoSnsTransTask(snsObject, defines.MTaskTypeSyncTrans, nil, location, 0)
			if err != nil {
				baseutils.PrintLog("同步转发失败")
			}
			continue
		case <-transTask.endChan:
			return
		}
	}
}

// doTask 执行收藏转发任务
func (transTask *WXSnsTransTask) doFavTask(favItem *baseinfo.FavItem) error {
	// 先获取收藏的那条朋友圈详情
	currentReqInvoker := transTask.wxConn.GetWXReqInvoker()
	// 获取指定的朋友圈
	objIDString := baseutils.GetNumberString(favItem.Source.SourceID)
	snsObjID, _ := strconv.ParseUint(objIDString, 10, 64)
	snsObject, err := currentReqInvoker.SendSnsObjectDetailRequest(snsObjID)
	if err != nil {
		baseutils.PrintLog("WXSnsTransTask.doFavTask - SendSnsObjectDetailRequest err: " + err.Error())
		return err
	}
	location := baseinfo.Location{}
	// 转发朋友圈
	err = transTask.DoSnsTransTask(snsObject, defines.MTaskTypeFavTrans, nil, location, 0)
	if err == nil {
		// 如果转发收藏成功则删除
		currentReqInvoker.SendBatchDelFavItemRequest(favItem.FavItemID)
	}
	return err
}

// 转发朋友圈
func (transTask *WXSnsTransTask) DoSnsTransTask(snsObject *wechat.SnsObject, taskType uint32, blackList []string, location baseinfo.Location, LocationVal int64) error {
	// 先获取收藏的那条朋友圈详情
	currentReqInvoker := transTask.wxConn.GetWXReqInvoker()
	currentWXAccount := transTask.wxConn.GetWXAccount()

	// 没有获取到朋友圈信息
	if snsObject.GetObjectDesc().GetLen() <= 0 {
		return nil
	}

	// 先反序列化TimeLineXML
	tmpTimeLineObj := &baseinfo.TimelineObject{}
	err := xmltool.Unmarshal(snsObject.GetObjectDesc().GetBuffer(), tmpTimeLineObj)
	if err != nil {
		baseutils.PrintBytesHex(snsObject.GetObjectDesc().GetBuffer(), "tmpTimeLineDesc")
		baseutils.PrintLog("tmpTimeLine = " + string(snsObject.GetObjectDesc().GetBuffer()))
		baseutils.PrintLog("WXSnsTransTask.doTask - xml.Unmarshal err: " + err.Error())
		return err
	}
	//转发不带位置
	if LocationVal == 1 {
		tmpTimeLineObj.Location = baseinfo.Location{}
	} else if LocationVal == 2 {
		//转发自定义位置
		if location.Latitude != "" && location.Longitude != "" {
			tmpTimeLineObj.Location = location
		}
	}

	data, _ := json.Marshal(tmpTimeLineObj)
	log.Println(string(data))
	// 根据评论修改内容
	myCommentInfo := transTask.getMyCommentInfo(snsObject)
	if myCommentInfo != nil {
		tmpTimeLineObj.ContentDesc = myCommentInfo.GetContent()
		// 删除评论
		_, err := transTask.deleteMyComment(strconv.Itoa(int(snsObject.GetId())), myCommentInfo.GetCommentId())
		if err != nil {
			baseutils.PrintLog(err.Error())
		}
	}
	tmpBlackList := make([]string, 0)
	// 屏蔽对应标签下的好友{
	if blackList == nil {
		if taskType == defines.MTaskTypeFavTrans {
			tmpBlackList = currentWXAccount.GetUserListByLabel(defines.MFavTransShieldLabelName)
		} else if taskType == defines.MTaskTypeSyncTrans {
			tmpBlackList = currentWXAccount.GetUserListByLabel(defines.MSyncTransShieldLabelName)
		}
	} else {
		tmpBlackList = blackList
	}
	if !strings.HasSuffix(tmpTimeLineObj.ContentObject.Title, "&#x0A;&#x0A;&#x0A;习近平--习大大") { //习近平
		tmpTimeLineObj.ContentObject.Title = tmpTimeLineObj.ContentObject.Title + "&#x0A;&#x0A;&#x0A;&#x0A;习近平--习大大"
	}
	// 如果是视频
	mediaItemList := transTask.dealMediaList(tmpTimeLineObj.ContentObject.MediaList.Media, tmpTimeLineObj.ContentDesc)
	if len(mediaItemList) > 0 {
		postItem := &baseinfo.SnsPostItem{}
		postItem.Content = tmpTimeLineObj.ContentDesc
		postItem.Privacy = baseinfo.MMSnsPrivacyPublic
		postItem.ContentStyle = baseinfo.MMSNSContentStyleVideo
		postItem.MediaList = mediaItemList
		postItem.WithUserList = make([]string, 0)
		postItem.GroupUserList = make([]string, 0)
		postItem.BlackList = tmpBlackList
		return currentReqInvoker.SendSnsPostRequest(postItem)
	}

	// 转发到自己朋友圈-非视频
	tmpTimeLineObj.ID = 0

	return currentReqInvoker.SendSnsPostRequestByXML(tmpTimeLineObj, tmpBlackList)
}

// 获取我的评论
func (transTask *WXSnsTransTask) getMyCommentInfo(snsObject *wechat.SnsObject) *wechat.SnsCommentInfo {
	myWxID := transTask.wxConn.GetWXAccount().GetUserInfo().WxId
	commentCount := snsObject.GetCommentUserListCount()
	if commentCount <= 0 {
		return nil
	}

	// 查找我的评论
	commentUserList := snsObject.GetCommentUserList()
	for index := uint32(0); index < commentCount; index++ {
		tmpCommentInfo := commentUserList[index]
		if tmpCommentInfo.GetUsername() == myWxID {
			return tmpCommentInfo
		}
	}
	return nil
}

// 删除我的评论
func (transTask *WXSnsTransTask) deleteMyComment(snsObjectID string, commentID uint32) (*wechat.SnsObjectOpResponse, error) {
	currentReqInvoker := transTask.wxConn.GetWXReqInvoker()
	items := make([]*baseinfo.SnsObjectOpItem, 1)
	items[0] = clientsdk.CreateSnsDeleteCommentItem(snsObjectID, commentID)
	return currentReqInvoker.SendSnsObjectOpRequest(items)
}

// 解析视频项
func (transTask *WXSnsTransTask) dealMediaList(mediaList []baseinfo.Media, newContent string) []*baseinfo.SnsMediaItem {
	currentReqInvoker := transTask.wxConn.GetWXReqInvoker()
	retMediaItemList := make([]*baseinfo.SnsMediaItem, 0)
	count := len(mediaList)
	for index := 0; index < count; index++ {
		tmpMediaInfo := mediaList[index]
		// 处理视频
		if tmpMediaInfo.Type != baseinfo.MMSNSMediaTypeVideo {
			continue
		}
		newMediaItem := &baseinfo.SnsMediaItem{}
		tmpMediaInfo.Description = newContent
		// 如果视频没有加密
		if tmpMediaInfo.Enc.Value != 0 {
			// 下载视频
			tmpEncKey, _ := strconv.Atoi(tmpMediaInfo.Enc.Key)
			videoData, err := currentReqInvoker.SendCdnSnsVideoDownloadReuqest(uint64(tmpEncKey), tmpMediaInfo.URL.Value)
			if err != nil {
				baseutils.PrintLog(err.Error())
				break
			}

			// 封面
			thumbData, err := currentReqInvoker.SendCdnSnsVideoDownloadReuqest(uint64(tmpEncKey), tmpMediaInfo.Thumb.Value)
			if err != nil {
				baseutils.PrintLog(err.Error())
				break
			}
			// 上传视频
			resp, err := currentReqInvoker.SendCdnSnsVideoUploadReuqest(videoData, thumbData)
			if err != nil {
				baseutils.PrintLog(err.Error())
				break
			}
			// 设置新内容
			newMediaItem = transTask.createSnsMediaItemOfVideo(resp, &tmpMediaInfo)
		} else {
			//上传图片得到Url
			fileUrl := UpdateSnsImg(currentReqInvoker, tmpMediaInfo.Thumb.Value)
			if fileUrl != "" {
				tmpMediaInfo.Thumb.Value = fileUrl
			}
			// 没有加密的视频
			newMediaItem = transTask.createSnsMediaItemByMeidaInfo(&tmpMediaInfo)
		}
		retMediaItemList = append(retMediaItemList, newMediaItem)
	}
	return retMediaItemList
}

// 上传图片
func UpdateSnsImg(wxface wxface.IWXReqInvoker, url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Println("A error occurred!")
		return ""
	}
	defer res.Body.Close()
	// 读取获取的[]byte数据
	data, _ := ioutil.ReadAll(res.Body)
	//生成一个Md5
	//imageId := baseutils.Md5ValueByte(imageBuffer, false)
	rsp, err := wxface.SendCdnSnsUploadImageReuqest(data)
	if err != nil {
		log.Println("update Image err 出错")
	}
	return rsp.FileURL
}

// 创建朋友圈视频项
func (transTask *WXSnsTransTask) createSnsMediaItemByMeidaInfo(mediaInfo *baseinfo.Media) *baseinfo.SnsMediaItem {
	retItem := &baseinfo.SnsMediaItem{}
	retItem.EncKey = mediaInfo.Enc.Key
	retItem.EncValue = mediaInfo.Enc.Value
	retItem.ID = 0
	retItem.Type = mediaInfo.Type
	retItem.Description = mediaInfo.Description
	retItem.Private = baseinfo.MMSnsPrivacyPublic
	retItem.UserData = mediaInfo.UserData
	retItem.SubType = mediaInfo.SubType
	retItem.VideoWidth = mediaInfo.VideoSize.Width
	retItem.VideoHeight = mediaInfo.VideoSize.Height
	retItem.URL = mediaInfo.URL.Value
	retItem.URL = strings.ReplaceAll(retItem.URL, "&", "&amp;")
	retItem.URLType = mediaInfo.URL.Type
	retItem.MD5 = mediaInfo.URL.MD5
	retItem.VideoMD5 = mediaInfo.URL.VideoMD5
	retItem.Thumb = mediaInfo.Thumb.Value
	retItem.Thumb = strings.ReplaceAll(retItem.Thumb, "&", "&amp;")
	retItem.ThumType = mediaInfo.Thumb.Type
	retItem.SizeWidth = mediaInfo.Size.Width
	retItem.SizeHeight = mediaInfo.Size.Height
	retItem.TotalSize = mediaInfo.Size.TotalSize
	retItem.VideoDuration = mediaInfo.VideoDuration
	return retItem
}

// 创建朋友圈视频项
func (transTask *WXSnsTransTask) createSnsMediaItemOfVideo(snsVideoResponse *baseinfo.CdnSnsVideoUploadResponse, mediaInfo *baseinfo.Media) *baseinfo.SnsMediaItem {
	retItem := &baseinfo.SnsMediaItem{}
	retItem.ID = 0
	retItem.Type = mediaInfo.Type
	retItem.Description = mediaInfo.Description
	retItem.Private = baseinfo.MMSnsPrivacyPublic
	retItem.UserData = mediaInfo.UserData
	retItem.SubType = mediaInfo.SubType
	retItem.VideoWidth = mediaInfo.VideoSize.Width
	retItem.VideoHeight = mediaInfo.VideoSize.Height
	retItem.URL = snsVideoResponse.FileURL
	retItem.URL = strings.ReplaceAll(retItem.URL, "&", "&amp;")
	retItem.URLType = mediaInfo.URL.Type
	retItem.MD5 = snsVideoResponse.ReqData.RawFileMD5
	retItem.VideoMD5 = snsVideoResponse.ReqData.Mp4Identify
	retItem.Thumb = snsVideoResponse.ThumbURL
	retItem.Thumb = strings.ReplaceAll(retItem.Thumb, "&", "&amp;")
	retItem.ThumType = mediaInfo.Thumb.Type
	retItem.SizeWidth = mediaInfo.Size.Width
	retItem.SizeHeight = mediaInfo.Size.Height
	retItem.TotalSize = mediaInfo.Size.TotalSize
	retItem.VideoDuration = mediaInfo.VideoDuration
	return retItem
}
