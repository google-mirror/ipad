package wxface

// IWXConnect 微信链接接口
type IWXConnect interface {
	// 开启
	Start() error
	// 关闭
	Stop()

	StopWithReConnect(isReConnect bool)

	// 设置微信链接ID
	SetWXConnID(wxConnID uint32)
	// 获取WX链接ID
	GetWXConnID() uint32
	GetWXUuidKey() string
	GetWXAccount() IWXAccount
	// 获取同步管理器
	//GetWXSyncMgr() IWXSyncMgr
	// 获取好友消息管理器
	GetWXFriendMsgMgr() IWXUserMsgMgr
	// 判断是否处于链接状态
	IsConnected() bool
	// 发送给消息队列去处理
	SendToWXMsgHandler(wxResp IWXResponse)
	// 添加到长链接请求队列
	SendToWXLongReqQueue(wxLongReq IWXLongRequest)
	// 等待 waitTimes后发送心跳包
	SendHeartBeatWaitingSeconds(seconds uint32)
}
