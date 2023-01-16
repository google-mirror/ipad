package wxface

// IWXMsgHandler 微信响应处理器
type IWXMsgHandler interface {
	AddRouter(respID uint32, wxRouter IWXRouter) // 为消息添加具体的处理逻辑
	GetRouterByRespID(urlID uint32) IWXRouter    // 获取对应的路由处理器
	SendWXRespToTaskQueue(response IWXResponse)  // 将消息交给TaskQueue,由worker进行处理
}
