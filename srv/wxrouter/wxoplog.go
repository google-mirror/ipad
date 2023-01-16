package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXOplogRouter 心跳包响应路由
type WXOplogRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (oplog *WXOplogRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析退出登陆响应包
	var oplogResp wechat.OplogResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &oplogResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return oplogResp, nil
}
