package wxface

import "feiyu.com/wx/clientsdk/baseinfo"

// IWXGrabHBMgr 自动抢红包管理器
type IWXGrabHBMgr interface {
	// 开始抢红包协程
	Start()
	// 结束抢红包协程
	Stop()
	// 抢下一个红包
	GrapNext()
	// 添加红包项
	AddHBItem(hbItem *baseinfo.HongBaoItem)
	// 获取当前正在抢的红包项
	GetCurrentHBItem() *baseinfo.HongBaoItem
	// 是否开启自动抢红包功能
	IsAutoGrap() bool
	// 设置自动抢红包开关
	SetAutoGrap(bFlag bool)
}
