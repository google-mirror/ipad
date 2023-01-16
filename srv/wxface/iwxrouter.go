package wxface

// IWXRouter 微信消息处理路由
type IWXRouter interface {
	PreHandle(response IWXResponse) error             //在处理conn业务之前的钩子方法
	Handle(response IWXResponse) (interface{}, error) //处理conn业务的方法
	PostHandle(response IWXResponse) error            //处理conn业务之后的钩子方法
}
