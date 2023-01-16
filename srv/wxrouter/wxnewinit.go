package wxrouter

import (
	"bytes"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/db"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXNewSyncRouter 获取二维码响应路由
type WXNewInitRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXNewInitRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentWXCache := currentWXConn.GetWXCache()
	// 同步响应
	var initResp wechat.NewInitResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &initResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 跟新同步Key
	syncKey := initResp.GetCurrentSyncKey().GetBuffer()
	syncKeyMgr := currentUserInfo.SyncKeyMgr()
	syncKeyMgr.SetMaxKey(initResp.GetMaxSyncKey())
	syncKeyMgr.SetMaxKey(initResp.GetCurrentSyncKey())

	//保存SyncKey
	if len(currentUserInfo.SyncKey) <= 0 || !bytes.Equal(currentUserInfo.SyncKey, syncKey) {
		currentUserInfo.SyncKey = syncKey
		db.UpdateUserInfo(currentUserInfo)
	}

	if initResp.GetContinueFlag() <= 0 {
		//if !currentWXCache.IsInitNewSyncFinished() {
		//	currentWXCache.SetInitNewSyncFinished(true)
		//}
	}

	// 如果没有同步到数据则返回
	cmdList := initResp.GetCmdList()
	syncCount := initResp.GetCmdCount()
	//log.Info(initResp.GetContinueFlag())

	//log.Info(initResp.GetContinueFlag(), syncCount)
	//redis 发布结构体
	messageResp := new(table.SyncMessageResponse)
	// 遍历同步的信息和群
	itemList := cmdList
	for index := uint32(0); index < syncCount; index++ {
		item := itemList[index]
		itemID := item.GetCmdId()
		messageResp.SetMessage(item.GetCmdBuf().GetData(), int32(itemID))
	}
	//发布同步信息消息
	_ = db.PublishSyncMsgWxMessage(currentWXAccount.GetUserInfo(), *messageResp)

	//todo 注释后可能造成长时间未同步后 同步到的消息不全
	// 如果数量超过10条则继续同步
	//if !currentWXCache.IsInitNewSyncFinished() {
	//	_, _ = bizcgi.SendNewInitSyncRequest(currentWXAccount, false)
	//}
	return initResp, nil
}

/*// dealModContact 处理同步到的联系人信息
func dealModContact(wxConn wxface.IWXConnect, modContact *wechat.ModContact) {
	currentWXAccount := wxConn.GetWXAccount()
	userName := modContact.GetUserName().GetStr()

	// 处理群
	if strings.HasSuffix(userName, "@chatroom") {
		memberCount := modContact.GetNewChatroomData().GetMemberCount()
		if memberCount <= 0 {
			// 被移除群聊
			currentWXAccount.RemoveWXGroup(userName)
		} else {
			// 新增群聊
			currentWXAccount.AddWXGroup(modContact)
		}
		return
	}
	// 公众号
	hasExternalInfo := false
	if len(modContact.GetCustomizedInfo().GetExternalInfo()) > 0 {
		hasExternalInfo = true
	}
	if strings.HasPrefix(userName, "gh_") || hasExternalInfo {
		return
	}
	// 好友
	currentWXAccount.AddWXFriendContact(modContact)
}

// dealAddMsg 处理同步到的消息
func dealAddMsg(wxConn wxface.IWXConnect, addMsg *wechat.AddMsg) {
	msgType := addMsg.GetMsgType()
	msgContent := addMsg.GetContent().GetStr()
	fromUserName := addMsg.GetFromUserName().GetStr()
	toUserName := addMsg.GetToUserName().GetStr()

	// 判断是否是命令
	if toUserName == baseinfo.FileHelperWXID {
		// 处理命令
		if msgType == baseinfo.MMAddMsgTypeText {
			//dealCommand(wxConn, msgContent)
		}

		// 名片类型
		if msgType == baseinfo.MMAddMsgTypeCard {
			//dealPersionalCard(wxConn, addMsg)
		}
		return
	}

	// 判断 是否是群红包
	if msgType == baseinfo.MMAddMsgTypeRefer {
		if strings.HasSuffix(toUserName, "@chatroom") {
			dealGroupHB(wxConn, msgContent, toUserName)
		} else if strings.HasSuffix(fromUserName, "@chatroom") {
			dealGroupHB(wxConn, msgContent, fromUserName)
		}
		return
	}

	// 防消息撤回管理器
	currentTaskMgr := wxConn.GetWXTaskMgr()
	taskMgr, _ := currentTaskMgr.(*wxcore.WXTaskMgr)
	currentRevokeTask := taskMgr.GetRevokeTask()
	// 文本消息
	if msgType == baseinfo.MMAddMsgTypeText {
		// 如果开启了消息防撤回
		if currentRevokeTask.IsAvoidRevoke() {
			currentRevokeTask.AddNewMsg(addMsg)
		}
		return
	}

	// 系统消息
	if msgType == baseinfo.MMAddMsgTypeSystemMsg {
		tmpSysMsg := &baseinfo.SysMsg{}
		err := xml.Unmarshal([]byte(msgContent), tmpSysMsg)
		if err != nil {
			return
		}

		// 如果是撤回消息
		if tmpSysMsg.Type == "revokemsg" {
			if currentRevokeTask.IsAvoidRevoke() {
				currentRevokeTask.OnRevokeMsg(tmpSysMsg.RevokeMsg)
			}
			return
		}
	}
}

// dealCommand 处理收到的命令
func dealCommand(wxConn wxface.IWXConnect, cmdText string) {
	currentReqInvoker := wxConn.GetWXReqInvoker()
	currentWXCache := wxConn.GetWXCache()
	currentTaskMgr := wxConn.GetWXTaskMgr()
	taskMgr, _ := currentTaskMgr.(*wxcore.WXTaskMgr)
	currentGrapHBTask := taskMgr.GetGrabHBTask()
	currentGroupTask := taskMgr.GetGroupTask()
	currentSnsTransTask := taskMgr.GetSnsTransTask()
	currentRevokeTask := taskMgr.GetRevokeTask()
	currentSnsTask := taskMgr.GetSnsTask()
	currentVerifyTask := taskMgr.GetVerifyTask()
	currentWXFileHelperMgr := wxConn.GetWXFileHelperMgr()

	// 退出登陆
	if cmdText == "000" {
		currentWXFileHelperMgr.AddNewTipMsg("已退出系统")
		currentReqInvoker.SendLogoutRequest()
		return
	}

	// 开启自动抢红包
	if cmdText == "101" {
		bAutoGrap := currentGrapHBTask.IsAutoGrap()
		if bAutoGrap {
			currentGrapHBTask.SetAutoGrap(false)
			currentWXFileHelperMgr.AddNewTipMsg("自动抢红包已关闭")
		} else {
			currentGrapHBTask.SetAutoGrap(true)
			currentWXFileHelperMgr.AddNewTipMsg("自动抢红包已开启")
		}
		return
	}

	// 自动同步转发、收藏转发
	if cmdText == "201" {
		isAutoRelay := currentSnsTransTask.IsAutoRelay()
		if isAutoRelay {
			currentSnsTransTask.SetAutoRelay(false)
			currentWXFileHelperMgr.AddNewTipMsg("已关闭收藏转发")
		} else {
			currentSnsTransTask.SetAutoRelay(true)
			currentWXFileHelperMgr.AddNewTipMsg("已开启收藏转发")
		}
		return
	}

	// 自动同步转发、收藏转发
	if cmdText == "301" {
		isSyncTrans := currentSnsTransTask.IsSyncTrans()
		if isSyncTrans {
			currentSnsTransTask.SetSyncTrans(false)
			currentWXFileHelperMgr.AddNewTipMsg("已关闭同步转发,并清空了同步转发好友列表")
		} else {
			currentSnsTransTask.SetSyncTrans(true)
			currentWXFileHelperMgr.AddNewTipMsg("已开启同步转发")
		}
		return
	}

	// 自动点赞朋友圈
	if cmdText == "401" {
		isAutoThumbUP := currentSnsTask.IsAutoThumbUP()
		if isAutoThumbUP {
			currentSnsTask.SetAutoThumbUP(false)
			currentWXFileHelperMgr.AddNewTipMsg("已关闭自动点赞朋友圈")
		} else {
			currentSnsTask.SetAutoThumbUP(true)
			currentWXFileHelperMgr.AddNewTipMsg("已开启自动点赞朋友圈")
			// 初始化
			currentReqInvoker.SendSnsTimeLineRequest("", 0)
		}
		return
	}

	// 自动保存所有群聊到通讯录
	if cmdText == "501" {
		// 开始 dump到
		currentGroupTask.StartSaveToAddressBook(true)
		return
	}

	// 自动取消保存所有通讯录群聊
	if cmdText == "502" {
		// 开始 dump到
		currentGroupTask.StartSaveToAddressBook(false)
		return
	}

	// 开启/关闭消息防撤回功能
	if cmdText == "601" {
		// 开始 dump到
		isAvoidRevoke := currentRevokeTask.IsAvoidRevoke()
		if isAvoidRevoke {
			currentRevokeTask.SetAvoidRevoke(false)
			currentWXFileHelperMgr.AddNewTipMsg("消息防撤回已关闭")
		} else {
			currentRevokeTask.SetAvoidRevoke(true)
			currentWXFileHelperMgr.AddNewTipMsg("消息防撤回已开启")
		}
		return
	}

	// 删除30条朋友圈
	if cmdText == "701" {
		currentSnsTask.DeleteTenSnsObject()
		return
	}

	// 开启/关闭自动通过好友验证
	if cmdText == "801" {
		needVerify := currentVerifyTask.IsNeedVerify()
		if needVerify {
			currentVerifyTask.SetNeedVerify(false)
			currentWXFileHelperMgr.AddNewTipMsg("自动通过验证已开启")
		} else {
			currentVerifyTask.SetNeedVerify(true)
			currentWXFileHelperMgr.AddNewTipMsg("自动通过验证已关闭")
		}
		return
	}

	// 开启/关闭自动通过好友验证
	if cmdText == "901" {
		currentWXFileHelperMgr.AddNewTipMsg("开始获取QB信息")
		reqText, err := wxConn.GetWXAccount().GetBindQueryNewReq()
		if err != nil {
			currentWXFileHelperMgr.AddNewTipMsg("获取QB信息失败")
			return
		}
		tmpReqItem := &baseinfo.TenPayReqItem{}
		tmpReqItem.CgiCMD = 72
		tmpReqItem.ReqText = reqText
		err = currentReqInvoker.SendBindQueryNewRequest(tmpReqItem)
		if err != nil {
			currentWXFileHelperMgr.AddNewTipMsg("获取QB信息失败")
			return
		}
		return
	}

	// 发送使用说明
	currentWXCache.SendUsage()
}

// 处理名片信息
func dealPersionalCard(wxConn wxface.IWXConnect, addMsg *wechat.AddMsg) {
	currentReqInvoker := wxConn.GetWXReqInvoker()
	currentTaskMgr := wxConn.GetWXTaskMgr()
	taskMgr, _ := currentTaskMgr.(*wxcore.WXTaskMgr)
	currentSnsTransTask := taskMgr.GetSnsTransTask()
	currentWXFileHelperMgr := wxConn.GetWXFileHelperMgr()

	// 先反序列化CardInfo
	msgContent := addMsg.GetContent().GetStr()
	cardInfo := &baseinfo.CardInfo{}
	err := xml.Unmarshal([]byte(msgContent), cardInfo)
	if err != nil {
		baseutils.PrintLog(err.Error())
		return
	}
	// 如果是公众号-不予理会
	if strings.HasPrefix(cardInfo.UserName, "gh_") {
		return
	}
	// 如果开启了自动转发功能
	if currentSnsTransTask.IsSyncTrans() {
		// 如果不存在则增加，然后初始化
		if currentSnsTransTask.AddSyncTransFriend(cardInfo.UserName) {
			currentReqInvoker.SendSnsUserPageRequest(cardInfo.UserName, "", 0, false)
		}
		showText := "成功绑定名片【" + cardInfo.NickName + "】\n"
		showText = showText + "当前绑定" + strconv.Itoa(currentSnsTransTask.GetTransFriendCount()) + "个"
		currentWXFileHelperMgr.AddNewTipMsg(showText)
	}
}

// 处理自动抢红包操作
func dealGroupHB(wxConn wxface.IWXConnect, content string, groupWXID string) {
	// 判断是否开启了自动抢红包功能
	currentTaskMgr := wxConn.GetWXTaskMgr()
	taskMgr, _ := currentTaskMgr.(*wxcore.WXTaskMgr)
	currentGrapHBMgr := taskMgr.GetGrabHBTask()
	if !currentGrapHBMgr.IsAutoGrap() {
		return
	}

	// 解析引用的消息
	tmpMsg := new(baseinfo.Msg)
	err := xml.Unmarshal([]byte(content), tmpMsg)
	if err != nil {
		return
	}

	// 判断是否是红包
	if tmpMsg.APPMsg.MsgType != baseinfo.MMAppMsgTypePayInfo {
		return
	}

	// 判断是否红包类型
	if tmpMsg.APPMsg.WCPayInfo.SceneID != baseinfo.MMPayInfoSceneIDHongBao {
		return
	}

	// 这里要判断 该群是否有消息置顶
	currentReqInvoker := wxConn.GetWXReqInvoker()
	resp, err := currentReqInvoker.SendGetContactRequestForHB(groupWXID)
	if err != nil || resp.GetBaseResponse().GetRet() != 0 {
		return
	}
	tmpBitVal := resp.GetContactList()[0].GetBitVal()
	if tmpBitVal&baseinfo.MMBitValChatOnTop == 0 {
		return
	}

	// 开始抢红包
	hbItem := new(baseinfo.HongBaoItem)
	hbItem.NativeURL = tmpMsg.APPMsg.WCPayInfo.NativeURL
	hongBaoURLItem, err := clientsdk.ParseHongBaoURL(hbItem.NativeURL, groupWXID)
	if err != nil {
		return
	}
	hbItem.URLItem = hongBaoURLItem
	currentGrapHBMgr.AddHBItem(hbItem)
}

// 处理删除联系人
func dealDelContact(wxConn wxface.IWXConnect, userName string) {
	currentWXAccount := wxConn.GetWXAccount()
	bRet := currentWXAccount.RemoveWXFriendID(userName)
	if !bRet {
		bRet = currentWXAccount.RemoveWXGhID(userName)
		if !bRet {
			currentWXAccount.RemoveWXGroup(userName)
		}
	}
}
*/
