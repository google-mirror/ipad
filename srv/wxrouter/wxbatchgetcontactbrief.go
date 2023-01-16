package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXBatchGetContactBriefInfoReqRouter 批量获取联系人信息
type WXBatchGetContactBriefInfoReqRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (wxbgcbirr *WXBatchGetContactBriefInfoReqRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析批量获取联系人信息响应包
	var briefInfoResp wechat.BatchGetContactBriefInfoResp
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &briefInfoResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	return briefInfoResp, nil
}

func isSystemWXID(userName string) bool {
	if userName == "qqmail" ||
		userName == "qmessage" ||
		userName == "mphelper" ||
		userName == "filehelper" ||
		userName == "weixin" ||
		userName == "floatbottle" ||
		userName == "fmessage" ||
		userName == "medianote" {
		return true
	}
	return false
}
