package wxrouter

import (
	"feiyu.com/wx/srv/wxmgr"
	"strings"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
)

// WXGetContactRouter 批量获取联系人响应路由
type WXGetContactRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXGetContactRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentWXUserInfo := currentWXAccount.GetUserInfo()
	//currentWXCache := currentWXConn.GetWXCache()
	//currentReqInvoker := currentWXConn.GetWXReqInvoker()

	// 解析获取联系人响应
	getContactResp := new(wechat.GetContactResponse)
	err := clientsdk.ParseResponseData(currentWXUserInfo, wxResp.GetPackHeader(), getContactResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	verifyUserList := getContactResp.GetVerifyUserValidTicketList()
	contactList := getContactResp.GetContactList()
	count := len(contactList)
	for index := 0; index < count; index++ {
		tmpContact := contactList[index]
		// 联系人
		if tmpContact.GetPersonalCard() == 1 {
			// 如果是僵死粉
			if strings.HasPrefix(verifyUserList[index].GetAntispamticket(), "v2_") {
				baseutils.PrintLog("僵尸粉：" + tmpContact.GetNickName().GetStr())
			}
			continue
		}
		// 群
		if tmpContact.GetChatroomVersion() > 0 {
		}
	}

	// 再一次获取剩余列表
	//userNameList := currentWXCache.GetNextInitContactWxidList(20)
	//if len(userNameList) > 0 {
	//	currentReqInvoker.SendGetContactRequest(userNameList, []string{}, []string{}, false)
	//} else {
	//	baseutils.PrintLog("检测僵尸粉完毕")
	//}
	return getContactResp, nil
}
