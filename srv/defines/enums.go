package defines

const (
	// EGroupTaskStateIdel 群聊任务状态：空闲
	EGroupTaskStateIdel uint32 = 0
	// EGroupTaskStateSave 群聊任务状态：保存群聊
	EGroupTaskStateSave uint32 = 1
	// EGroupTaskStateUnSave 群聊任务状态：取消保存群聊
	EGroupTaskStateUnSave uint32 = 2

	// EGroupTaskStateDownQrcode 群聊任务状态：下载群二维码
	EGroupTaskStateDownQrcode uint32 = 3
)

const (
	// MTaskTypeFavTrans 收藏转发
	MTaskTypeFavTrans uint32 = 0
	// MTaskTypeSyncTrans 同步转发
	MTaskTypeSyncTrans uint32 = 0
)

const (
	// MFavTransShieldLabelName 收藏转发屏蔽标签名
	MFavTransShieldLabelName string = "收藏转发屏蔽"
	// MSyncTransShieldLabelName 同步转发屏蔽标签名
	MSyncTransShieldLabelName string = "同步转发屏蔽"
)
