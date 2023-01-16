package wxface

import "feiyu.com/wx/clientsdk/baseinfo"

// IWXResponse 微信请求接口
type IWXResponse interface {
	GetWXUuidKey() string
	// 获取请求数据
	GetPackHeader() *baseinfo.PackHeader
	// 获取WXConncet
	//GetWXConncet() IWXConnect
}
