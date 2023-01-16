package wxrouter

import (
	"bytes"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/db"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/lunny/log"
)

// WXNewSyncRouter 获取二维码响应路由
type WXNewSyncRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXNewSyncRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	log.Debug("处理拉取到的消息")
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentWXCache := currentWXConn.GetWXCache()
	// 同步响应
	var syncResp wechat.NewSyncResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &syncResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 跟新同步Key
	syncKey := syncResp.GetKeyBuf().GetBuffer()
	//保存SyncKey
	if len(currentUserInfo.SyncKey) <= 0 || !bytes.Equal(currentUserInfo.SyncKey, syncResp.GetKeyBuf().GetBuffer()) {
		currentUserInfo.SyncKey = syncKey
		go db.UpdateUserInfo(currentUserInfo)
	}

	// 如果没有同步到数据则返回bInitNewSyncFinished
	cmdList := syncResp.GetCmdList()
	syncCount := cmdList.GetCount()
	//if syncCount <= 0 {
	//	if !currentWXCache.IsInitNewSyncFinished() {
	//		currentWXCache.SetInitNewSyncFinished(true)
	//	}
	//}
	//log.Info(syncResp.GetContinueFlag(), syncCount)
	//redis 发布结构体
	messageResp := new(table.SyncMessageResponse)
	// 遍历同步的信息和群
	itemList := cmdList.GetItemList()
	for index := uint32(0); index < syncCount; index++ {
		item := itemList[index]
		itemID := item.GetCmdId()
		messageResp.SetMessage(item.GetCmdBuf().GetData(), int32(itemID))
	}
	messageResp.Key = syncResp.GetKeyBuf()
	//发布同步信息消息
	go db.PublishSyncMsgWxMessage(currentWXAccount.GetUserInfo(), *messageResp)
	// 如果数量超过10条则继续同步
	//todo 注释后可能造成长时间未同步后 同步到的消息不全
	//if !currentWXCache.IsInitNewSyncFinished() {
	//	_, _ = bizcgi.SendNewSyncRequest(currentWXAccount, baseinfo.MMSyncSceneTypeNeed, false)
	//}
	return messageResp, nil
}
