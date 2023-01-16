package wxface

// IWXServer 微信服务
type IWXServer interface {
	GetWXMsgHandler() IWXMsgHandler
	AddWXRouter(funcID uint32, wxRouter IWXRouter)
}
