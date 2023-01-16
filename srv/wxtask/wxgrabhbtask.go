package wxtask

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/wxface"
)

// WXGrabHBTask 抢红包管理器
type WXGrabHBTask struct {
	wxConn wxface.IWXConnect
	// 待领红包项列表
	hongBaoItemChan chan *baseinfo.HongBaoItem
	// 可不可以抢红包
	canDoNext chan bool
	// 结束标志
	endChan chan bool
	// 是否正在运行
	isRunning bool
	// 当前在抢的红包
	currentHBItem *baseinfo.HongBaoItem
	// autoGrap
	autoGropFlag bool
}

// NewWXGrabHBTask 新建抢红包管理器
func NewWXGrabHBTask(wxConn wxface.IWXConnect) *WXGrabHBTask {
	return &WXGrabHBTask{
		wxConn:          wxConn,
		hongBaoItemChan: make(chan *baseinfo.HongBaoItem, 100),
		canDoNext:       make(chan bool, 1),
		endChan:         make(chan bool, 1),
		isRunning:       false,
		autoGropFlag:    false,
	}
}

func (ghbm *WXGrabHBTask) grapHB() {
	for {
		select {
		case ghbm.currentHBItem = <-ghbm.hongBaoItemChan:
			// 抢红包操作
			err := ghbm.wxConn.GetWXReqInvoker().SendReceiveWxHBRequest(ghbm.currentHBItem)
			if err != nil {
				// 抢下一个红包
				ghbm.GrapNext()
			}
			return
		case <-ghbm.endChan:
			//ghbm.endChan <- true
			//close(ghbm.endChan)
			return
		}
	}
}

// startGrap 开始抢红包
func (ghbm *WXGrabHBTask) startGrap() {
	for {
		select {
		case <-ghbm.canDoNext:
			ghbm.grapHB()
		case <-ghbm.endChan:
			return
		}
	}
}

// Start 开始抢红包协程
func (ghbm *WXGrabHBTask) Start() {
	go ghbm.startGrap()
	ghbm.canDoNext <- true
}

// Stop 结束抢红包协程
func (ghbm *WXGrabHBTask) Stop() {
	ghbm.endChan <- true
	//close(ghbm.endChan)
}

// GrapNext 抢下一个红包
func (ghbm *WXGrabHBTask) GrapNext() {
	ghbm.canDoNext <- true
}

// AddHBItem 添加红包项
func (ghbm *WXGrabHBTask) AddHBItem(hbItem *baseinfo.HongBaoItem) {
	ghbm.hongBaoItemChan <- hbItem
}

// GetCurrentHBItem 获取当前正在抢的红包项
func (ghbm *WXGrabHBTask) GetCurrentHBItem() *baseinfo.HongBaoItem {
	return ghbm.currentHBItem
}

// IsAutoGrap 是否开启了自动抢红包功能
func (ghbm *WXGrabHBTask) IsAutoGrap() bool {
	return ghbm.autoGropFlag
}

// SetAutoGrap 设置自动抢开关
func (ghbm *WXGrabHBTask) SetAutoGrap(bFlag bool) {
	ghbm.autoGropFlag = bFlag
}
