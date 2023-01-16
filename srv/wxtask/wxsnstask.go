package wxtask

import (
	"strconv"
	"time"

	"feiyu.com/wx/protobuf/wechat"

	"feiyu.com/wx/clientsdk/baseinfo"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/srv/wxface"
)

// WXSnsTask 朋友圈任务
type WXSnsTask struct {
	wxConn wxface.IWXConnect
	// 需要自动点赞到朋友圈列表
	snsObjChan chan *wechat.SnsObject
	// 朋友圈刷新时间
	snsFlushTimeChan chan uint32
	// 标记当前最新的朋友圈创建时间
	currentSnsCreateTime uint32
	// 自动点赞开启标志
	bAutoThumbUP bool
	// 结束标志
	endChan chan bool
	// 结束标志
	isRunning bool
}

// NewWXSnsTask 新建朋友圈任务管理器
func NewWXSnsTask(wxConn wxface.IWXConnect) *WXSnsTask {
	return &WXSnsTask{
		wxConn:               wxConn,
		snsObjChan:           make(chan *wechat.SnsObject, 100),
		snsFlushTimeChan:     make(chan uint32, 5),
		currentSnsCreateTime: 0,
		bAutoThumbUP:         false,
		isRunning:            false,
	}
}

// Start 开启任务
func (wxst *WXSnsTask) Start() {
	wxst.endChan = make(chan bool, 1)
	// 定时刷新朋友圈
	go wxst.startFlushSns()
	// 开启自动点赞线程
	go wxst.startThumbUP()
}

// Stop 关闭任务
func (wxst *WXSnsTask) Stop() {
	//wxst.endChan <- true
	close(wxst.endChan)
}

// SetAutoThumbUP 设置是否自动点赞朋友圈
func (wxst *WXSnsTask) SetAutoThumbUP(bFlag bool) {
	wxst.bAutoThumbUP = bFlag
	if !bFlag {
		wxst.currentSnsCreateTime = 0
	}
}

// IsAutoThumbUP 是否开启了自动点赞朋友圈
func (wxst *WXSnsTask) IsAutoThumbUP() bool {
	return wxst.bAutoThumbUP
}

// SetCurrentCreateTime 设置最新的朋友圈创建时间
func (wxst *WXSnsTask) SetCurrentCreateTime(tmpTime uint32) {
	wxst.currentSnsCreateTime = tmpTime
}

// GetCurrentCreateTime 获取最新的朋友圈创建时间
func (wxst *WXSnsTask) GetCurrentCreateTime() uint32 {
	return wxst.currentSnsCreateTime
}

// AddCommentItemObj 新增自动点赞项
func (wxst *WXSnsTask) AddCommentItemObj(snsObj *wechat.SnsObject) {
	wxst.snsObjChan <- snsObj
}

// 定时刷新朋友圈
func (wxst *WXSnsTask) startFlushSns() {
	currentReqInvoker := wxst.wxConn.GetWXReqInvoker()
	// 1分钟刷新一次
	wxst.snsFlushTimeChan <- 60
	for {
		select {
		case waitTimes := <-wxst.snsFlushTimeChan:
			time.Sleep(time.Second * time.Duration(waitTimes))
			// 如果开启了自动转发
			if wxst.bAutoThumbUP {
				// 刷新朋友圈首页
				currentReqInvoker.SendSnsTimeLineRequest("", 0)
			}
			wxst.snsFlushTimeChan <- 60
			continue
		case <-wxst.endChan:
			return
		}
	}
}

// 自动点赞朋友圈
func (wxst *WXSnsTask) startThumbUP() {
	currentReqInvoker := wxst.wxConn.GetWXReqInvoker()
	for {
		time.Sleep(time.Second * 1)
		select {
		case snsObj := <-wxst.snsObjChan:
			commentLikeItem := clientsdk.CreateSnsCommentLikeItem(snsObj.GetId(), snsObj.GetUsername())
			currentReqInvoker.SendSnsCommentRequest(commentLikeItem)
		case <-wxst.endChan:
			return
		}
	}
}

// SetRunningFlag 设置运行状态
func (wxst *WXSnsTask) setRunningFlag(isRunning bool) {
	wxst.isRunning = isRunning
}

// DeleteTenSnsObject 删除30条朋友圈
func (wxst *WXSnsTask) DeleteTenSnsObject() {
	currentWXFileHelperMgr := wxst.wxConn.GetWXFileHelperMgr()
	// 如果正在执行
	if wxst.isRunning {
		currentWXFileHelperMgr.AddNewTipMsg("正在清理朋友圈请稍等")
		return
	}
	go func() {
		wxst.setRunningFlag(true)
		defer wxst.setRunningFlag(false)
		currentWXFileHelperMgr.AddNewTipMsg("开始清理朋友圈...")
		// 获取30条朋友圈
		objIDList, err := wxst.getTenSnsObjectByCount(30)
		if err != nil {
			currentWXFileHelperMgr.AddNewTipMsg("清空朋友圈太频繁，稍后再试")
			return
		}

		// 清理朋友圈
		count := len(objIDList)
		for index := 0; index < count; index++ {
			resp, err := wxst.DeleteSnsByID(strconv.Itoa(int(objIDList[index])))
			if err != nil {
				currentWXFileHelperMgr.AddNewTipMsg("清理太频繁，成功清理" + strconv.Itoa(index) + "条朋友圈")
				return
			}

			// 如果返回出错
			if resp.GetBaseResponse().GetRet() != 0 {
				currentWXFileHelperMgr.AddNewTipMsg("清理太频繁，成功清理" + strconv.Itoa(index) + "条朋友圈")
				return
			}
			// 500毫秒清理一次
			time.Sleep(time.Millisecond * 500)
		}
		currentWXFileHelperMgr.AddNewTipMsg("清理完成, 成功清理" + strconv.Itoa(count) + "条朋友圈")
		return
	}()
}

// DeleteSnsByID 根据朋友圈ID删除朋友圈
func (wxst *WXSnsTask) DeleteSnsByID(snsIDList string) (*wechat.SnsObjectOpResponse, error) {
	currentReqInvoker := wxst.wxConn.GetWXReqInvoker()
	items := make([]*baseinfo.SnsObjectOpItem, 0)
	deleteItem := clientsdk.CreateSnsDeleteItem(snsIDList)
	items = append(items, deleteItem)
	return currentReqInvoker.SendSnsObjectOpRequest(items)
}

// getAllSnsObject 获取所有的朋友圈
func (wxst *WXSnsTask) getTenSnsObjectByCount(needCount uint32) ([]uint64, error) {
	currentReqInvoker := wxst.wxConn.GetWXReqInvoker()
	currentWXAccount := wxst.wxConn.GetWXAccount()
	tmpUserInfo := currentWXAccount.GetUserInfo()

	retList := make([]uint64, 0)
	currentCount := uint32(0)
	tmpMaxID := uint64(0)
	for currentCount < needCount {
		// 获取10条朋友圈
		resp, err := currentReqInvoker.SendSnsUserPageRequest(tmpUserInfo.WxId, "", tmpMaxID, true)
		if err != nil {
			return nil, err
		}

		objectList := resp.GetObjectList()
		count := len(objectList)
		for index := 0; index < count; index++ {
			tmpObject := objectList[index]
			retList = append(retList, tmpObject.GetId())
		}
		if count < 10 {
			break
		}
		currentCount = currentCount + uint32(count)
		tmpMaxID = objectList[count-1].GetId()
	}
	return retList, nil
}
