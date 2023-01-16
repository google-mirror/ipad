package wxcore

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXResponse 微信响应
type WXResponse struct {
	// 微信uuid key
	uuidKey    string
	packHeader *baseinfo.PackHeader
	//业务日志Id
	LogUUID string
}

// NewWXResponse 新建WXResponse
func NewWXResponse(uuidKey string, packHeader *baseinfo.PackHeader) wxface.IWXResponse {
	return &WXResponse{
		uuidKey:    uuidKey,
		packHeader: packHeader,
	}
}

// GetWXUuidKey 获取WX uuidKey
func (resp *WXResponse) GetWXUuidKey() string {
	return resp.uuidKey
}

// GetPackHeader 获取响应数据
func (resp *WXResponse) GetPackHeader() *baseinfo.PackHeader {
	return resp.packHeader
}

// GetWXConncet 获取WXConncet
func (resp *WXResponse) GetWXConncet() wxface.IWXConnect {
	return wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(resp.uuidKey)
}
