package wxrouter

import (
	"errors"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/lunny/log"
)

// WXPullMsgRouter 主动拉取消息
type WXPullMsgRouter struct {
	WXBaseRouter
}

// PreHandle 在处理conn业务之前的钩子方法
func (wxbr *WXPullMsgRouter) PreHandle(response wxface.IWXResponse) error {
	//currentWXConn := wxResp.GetWXConncet()
	currentAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(response.GetWXUuidKey())
	log.Debug("长连接收到需要拉取消息通知")
	syncKey := currentAccount.GetUserInfo().SyncKey
	log.Debug("拉取消息时syncKey为：", syncKey)
	var err error
	if syncKey == nil || len(syncKey) == 0 {
		//todo
		//currentWXConn.GetOnceInit().Do(func() {
		log.Debug("初始化syncKey")
		_, err = bizcgi.SendNewInitSyncRequest(currentAccount, false)
		//})
	} else {
		_, err = bizcgi.SendNewSyncRequest(currentAccount, baseinfo.MMSyncSceneTypeNeed, false)
	}
	if err == nil {
		return nil
	}
	return errors.New("拉取消息失败")
}

// Handle 处理conn业务的方法
func (cqr *WXPullMsgRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {

	//if response.GetPackHeader().URLID == 24 {
	//	wxAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(response.GetWXUuidKey())
	//	syncKey := wxAccount.GetUserInfo().SyncKey
	//	if syncKey == nil || len(syncKey) == 0 {
	//		wxAccount.GetOnceInit().Do(func() {
	//			wxAccount.GetWXReqInvoker().SendNewInitSyncRequest()
	//		})
	//	} else {
	//		bizcgi.SendNewSyncRequest(wxAccount, baseinfo.MMSyncSceneTypeNeed, false)
	//	}
	//	return
	//}
	return nil, nil
}
