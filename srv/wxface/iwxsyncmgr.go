package wxface

// IWXSyncMgr 消息同步管理器
type IWXSyncMgr interface {
	Start()
	Stop()

	// 发送同步请求
	SendNewSyncRequest()
	SendFavSyncRequest()
	SendSyncInitRequest()
}
