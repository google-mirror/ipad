package wxtask

import (
	"sync"
	"time"

	"feiyu.com/wx/clientsdk/baseinfo"

	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
)

// WXRevokeTask 消息防撤回管理器
type WXRevokeTask struct {
	wxConn wxface.IWXConnect
	// 缓存三分钟内的消息，只缓存文本消息
	addMsgList []*wechat.AddMsg
	// 防撤回开关
	bAvoidRevoke bool
	// 结束标志
	isRunning bool

	// 读写数据的读写锁
	revokeLock sync.RWMutex
}

// NewWXRevokeTask 新建消息防撤回任务管理器
func NewWXRevokeTask(wxConn wxface.IWXConnect) *WXRevokeTask {
	return &WXRevokeTask{
		wxConn:       wxConn,
		addMsgList:   make([]*wechat.AddMsg, 0),
		isRunning:    false,
		bAvoidRevoke: false,
	}
}

// Start 开启防消息撤回任务
func (wxrt *WXRevokeTask) Start() {
	wxrt.isRunning = true
	go wxrt.startRemoveOldMsg()
}

// Stop 停止防消息撤回任务
func (wxrt *WXRevokeTask) Stop() {
	wxrt.isRunning = false
	wxrt.addMsgList = make([]*wechat.AddMsg, 0)
}

// AddNewMsg 新增缓存消息
func (wxrt *WXRevokeTask) AddNewMsg(addMsg *wechat.AddMsg) {
	wxrt.addMsgList = append(wxrt.addMsgList, addMsg)
}

// 移除没用的消息
func (wxrt *WXRevokeTask) startRemoveOldMsg() {
	for {
		// 如果已经停止运行
		if !wxrt.isRunning {
			break
		}
		// 一分钟清除一次缓存消息
		time.Sleep(time.Second * 60)
		// 如果已经停止运行
		if !wxrt.isRunning {
			break
		}
		wxrt.checkAndRemoveOldMsg()
	}
}

// 移除过时的消息
func (wxrt *WXRevokeTask) checkAndRemoveOldMsg() {
	wxrt.revokeLock.Lock()
	defer wxrt.revokeLock.Unlock()

	// 如果时间相差超过3分钟则清除
	currentTime := uint32(time.Now().UnixNano() / 1000000000)
	msgCount := len(wxrt.addMsgList)
	for index := 0; index < msgCount; index++ {
		tmpAddMsg := wxrt.addMsgList[index]
		createTime := tmpAddMsg.GetCreateTime()

		// 判断是否超过180秒(三分钟)
		tmpTime := currentTime - createTime
		if tmpTime > 180 {
			wxrt.addMsgList = append(wxrt.addMsgList[:index], wxrt.addMsgList[index+1:]...)
			index = index - 1
			msgCount = msgCount - 1
		}
	}
}

// OnRevokeMsg 某人移除了消息
func (wxrt *WXRevokeTask) OnRevokeMsg(revokeMsg baseinfo.RevokeMsg) {
	wxrt.revokeLock.Lock()
	defer wxrt.revokeLock.Unlock()

	currentWXFileHelperMgr := wxrt.wxConn.GetWXFileHelperMgr()
	msgCount := len(wxrt.addMsgList)
	for index := 0; index < msgCount; index++ {
		tmpAddMsg := wxrt.addMsgList[index]
		if tmpAddMsg.GetNewMsgId() == revokeMsg.NewMsgID &&
			tmpAddMsg.GetMsgId() == revokeMsg.MsgID {

			strMsg := tmpAddMsg.GetContent().GetStr()
			showMsg := revokeMsg.ReplaceMsg + ": " + strMsg
			currentWXFileHelperMgr.AddNewTipMsg(showMsg)
			return
		}
	}
}

// SetAvoidRevoke 设置是否防消息撤回
func (wxrt *WXRevokeTask) SetAvoidRevoke(bFlag bool) {
	wxrt.bAvoidRevoke = bFlag
}

// IsAvoidRevoke 是否防消息撤回
func (wxrt *WXRevokeTask) IsAvoidRevoke() bool {
	return wxrt.bAvoidRevoke
}
