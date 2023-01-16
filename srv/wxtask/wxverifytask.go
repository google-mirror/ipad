package wxtask

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/wxface"
)

// WXVerifyTask 加我/拉我进群聊是否需要验证
type WXVerifyTask struct {
	wxConn wxface.IWXConnect
	// 结束标志
	bNeedVerify bool
}

// NewWXVerifyTask 新建朋友圈任务管理器
func NewWXVerifyTask(wxConn wxface.IWXConnect) *WXVerifyTask {
	return &WXVerifyTask{
		wxConn:      wxConn,
		bNeedVerify: true,
	}
}

// IsNeedVerify 判断是否需要验证
func (wxvt *WXVerifyTask) IsNeedVerify() bool {
	return wxvt.bNeedVerify
}

// SetNeedVerify 设置被添加，被拉入群聊时是否需要验证
func (wxvt *WXVerifyTask) SetNeedVerify(needVerify bool) {
	wxvt.bNeedVerify = needVerify
	currentReqInvoker := wxvt.wxConn.GetWXReqInvoker()
	// MMConcealAddNoNeedVerify 不需要验证， MMConcealAddNeedVerify 需要验证
	tmpNeedVerify := baseinfo.MMConcealAddNoNeedVerify
	if needVerify {
		tmpNeedVerify = baseinfo.MMConcealAddNeedVerify
	}
	needVerifyItem := clientsdk.CreateFunctionSwitchItem(baseinfo.MMAddMeNeedVerifyType, tmpNeedVerify)
	items := make([]*baseinfo.ModifyItem, 1)
	items[0] = needVerifyItem
	currentReqInvoker.SendOplogRequest(items)
}
