package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXFavSyncRouter 同步收藏响应路由
type WXFavSyncRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXFavSyncRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentReqInvoker := currentWXConn.GetWXReqInvoker()

	// 解析同步收藏响应包
	var favSyncResp wechat.FavSyncResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &favSyncResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 保存同步key下次使用
	// 解析同步收藏列表
	//cmdList := favSyncResp.GetCmdList()
	//tmpCount := cmdList.GetCount()
	tmpKeyBuf := favSyncResp.KeyBuf.GetBuffer()
	if len(tmpKeyBuf) <= 0 {
		//todo
		//currentReqInvoker.SendGetFavInfoRequest()
		return favSyncResp, nil
	}
	currentUserInfo.FavSyncKey = tmpKeyBuf
	//todo
	//currentReqInvoker.SendGetFavInfoRequest()
	//// 如果还没开启转发收藏
	//if !currentSnsTransTask.IsAutoRelay() {
	//	return nil,nil
	//}

	//itemList := cmdList.GetItemList()
	//for index := uint32(0); index < tmpCount; index++ {
	//	tmpItem := itemList[index]
	//	// 判断是否是新增的收藏
	//	tmpCmdID := tmpItem.GetCmdId()
	//	if tmpCmdID != baseinfo.MMFavSyncCmdAddItem {
	//		return nil,nil
	//	}
	//	// 处理新增的收藏
	//	addItem := &wechat.AddFavItem{}
	//	err := proto.Unmarshal(tmpItem.GetCmdBuf().GetData(), addItem)
	//	if err != nil {
	//		return nil,err
	//	}
	//	// 获取收藏详情
	//	return bizcgi.SendBatchGetFavItemReq(currentWXAccount, favInfoCache.LastFavID())
	//}
	return favSyncResp, nil
}
