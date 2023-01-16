package wxlink

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/lunny/log"
	"strings"
	"time"
)

// WXReqInvoker 微信请求调用器
type WXReqInvoker struct {
	wxAccount *WXAccount
}

// NewWXLongReqInvoker 新建一个请求调用器
func NewWXLongReqInvoker(wxAccount *WXAccount) wxface.IWXReqInvoker {
	return &WXReqInvoker{
		wxAccount: wxAccount,
	}
}

// 发送短信
func (wxqi *WXReqInvoker) SendWxBindOpMobileForRequest(OpCode int64, PhoneNumber string, VerifyCode string) (*wechat.BindOpMobileForRegResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendWxBindOpMobileForRegRequest(tmpUserInfo, OpCode, PhoneNumber, VerifyCode)
	if err != nil {
		return nil, err
	}
	WXBindOpMobileForReg := &wechat.BindOpMobileForRegResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, WXBindOpMobileForReg)
	if err != nil {
		return nil, err
	}
	if WXBindOpMobileForReg.BaseResponse.GetRet() == -301 {
		for _, tg := range WXBindOpMobileForReg.GetNewHostList().GetList() {
			log.Println(tg.GetOrigin(), tg.GetSubstitute())
			if tg.GetOrigin() == "extshort.weixin.qq.com" {
				/*G.MmtlsHost = tg.GetRedirect()*/
				tmpUserInfo.MMInfo.ShortHost = tg.GetSubstitute()
			}
			//ShortHost
			if tg.GetOrigin() == "long.weixin.qq.com" {
				/*G.Mmtlsip = tg.GetRedirect()*/
				tmpUserInfo.MMInfo.LongHost = tg.GetSubstitute()

			}
			/*log.Println("MmtlsHost", WXBindOpMobileForReg.MmtlsHost)
			log.Println("Mmtlsip", G.Mmtlsip)*/
		}
	}
	return WXBindOpMobileForReg, nil
}

// SendHybridManualAutoRequest
// ver hybrid 密钥版本 145 和146
func (wxqi *WXReqInvoker) SendHybridManualAutoRequest(newPass string, wxID string, ver byte) error {
	// 根据数据库缓存跟新UserInfo
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	// 发送请求
	packHeader, err := clientsdk.SendHybridManualAutoRequest(tmpUserInfo, newPass, wxID, ver)
	if err != nil {
		// 断开链接
		//wxmgr.WxConnectMgr.Stop(wxqi.wxAccount)
		return err
	}

	tmpUserInfo.LoginDataInfo.NewPassWord = newPass
	//tmpUserInfo.LoginDataInfo.UserName = wxID
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

//// SendAutoAuthRequest 发送Token登陆请求
//func (wxqi *WXReqInvoker) SendAutoAuthRequest() (interface{}, error) {
//	log.Info("发起二次登录")
//	// 重新链接，然后发送token登陆请求
//	/*wxmgr.WxConnectMgr.Stop(wxqi.wxAccount)
//	err := wxqi.wxconn.Start()
//	if err != nil {
//		return err
//	}
//
//	*/
//	// 发送请求
//send:
//	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
//	packHeader, err := clientsdk.SendAutoAuthRequest(tmpUserInfo)
//	if err != nil {
//		// 断开链接
//		log.Info("二次登录失败", err.Error())
//		//二次登录失败需要重新登录
//		if packHeader != nil {
//			switch packHeader.RetCode {
//			case baseinfo.MM_ERR_CERT_EXPIRED:
//				// 切换密钥登录
//				tmpUserInfo.SwitchRSACert()
//				goto send
//			case baseinfo.MMErrSessionTimeOut: // Session 会话过期
//				wxqi.wxAccount.SetLoginState(baseinfo.MMLoginStateNoLogin)
//				db.UpdateLoginStatus(wxqi.wxAccount.GetUserInfo().UUID, int32(wxqi.wxAccount.GetLoginState()), "二次登录失败需要重新登录")
//			default:
//				defer wxmgr.WxConnectMgr.Stop(wxqi.wxAccount)
//			}
//
//		}
//		return nil,err
//	}
//	// 发送给微信消息处理器
//	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
//	dealRouter := wxqi.wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
//	return dealRouter.Handle(wxResp)
//}

// 获取设备
func (wxqi *WXReqInvoker) SendGetSafetyInfoRequest() (*wechat.GetSafetyInfoResponse, error) {
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetSafetyInfoRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetSafetyInfoResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 删除设备
func (wxqi *WXReqInvoker) SendDelSafeDeviceRequest(deviceUUID string) (*wechat.DelSafeDeviceResponse, error) {
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendDelSafeDeviceRequest(tmpUserInfo, deviceUUID)
	if err != nil {

		return nil, err
	}
	response := &wechat.DelSafeDeviceResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 检测微信登录环境
func (wxqi *WXReqInvoker) SendCheckCanSetAliasRequest() (*wechat.CheckCanSetAliasResp, error) {
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendCheckCanSetAliasRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	response := &wechat.CheckCanSetAliasResp{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 扫码登录新设备
func (wxqi *WXReqInvoker) SendExtDeviceLoginConfirmGetRequest(url string) (*wechat.ExtDeviceLoginConfirmOKResponse, error) {
	// 判断是否与微信服务器握手成功
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendExtDeviceLoginConfirmGetRequest(tmpUserInfo, url)
	if err != nil {
		//if packHeader != nil &&
		//	packHeader.CheckSessionOut() {
		//	// 断开链接, 发送token登陆
		//	wxqi.wxconn.SendAutoAuthWaitingMinutes(4)
		//}
		return nil, err
	}
	responseH := &wechat.ExtDeviceLoginConfirmGetResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, responseH)
	if err != nil {
		return nil, err
	}
	if responseH.BaseResponse.GetRet() == 0 {
		time.Sleep(time.Second * 2)
		packHeader, err = clientsdk.ExtDeviceLoginConfirmOk(tmpUserInfo, url)
		if err != nil {
			return nil, err
		}
		response := &wechat.ExtDeviceLoginConfirmOKResponse{}
		err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, err
}

// 同步消息
func (wxqi *WXReqInvoker) SendWxSyncMsg(key string) (*wechat.NewSyncResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendWxSyncMsg(tmpUserInfo, key)
	if err != nil {

		return nil, err
	}
	response := &wechat.NewSyncResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SendNewInitSyncRequest 首次登录初始化
func (wxqi *WXReqInvoker) SendNewInitSyncRequest() (interface{}, error) {
	log.Debug("准备初始化消息")
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendNewInitSyncRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	dealRouter := wxqi.wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
	return dealRouter.Handle(wxResp)
}

// SendGetProfileRequest 获取微信账号配置信息
func (wxqi *WXReqInvoker) SendGetProfileNewRequest() (*wechat.GetProfileResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetProfileRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	response := &wechat.GetProfileResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	return response, err
}

// SendInitContactRequest 发送初始化联系人请求
func (wxqi *WXReqInvoker) SendInitContactRequest(contactSeq uint32) error {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendInitContactReq(tmpUserInfo, contactSeq)
	if err != nil {

		return err
	}

	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// 分页获取联系人
func (wxqi *WXReqInvoker) SendGetContactListPageRequest(CurrentWxcontactSeq uint32, CurrentChatRoomContactSeq uint32) (*wechat.InitContactResp, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendContactListPageRequest(tmpUserInfo, CurrentWxcontactSeq, CurrentChatRoomContactSeq)
	if err != nil {

		return nil, err
	}
	response := &wechat.InitContactResp{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	return response, err
}

// SendBatchGetContactBriefInfoReq 批量获取联系人
func (wxqi *WXReqInvoker) SendBatchGetContactBriefInfoReq(userWxidList []string) error {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendBatchGetContactBriefInfoReq(tmpUserInfo, userWxidList)
	if err != nil {

		return err
	}
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendGetContactRequest 获取联系人信息列表
func (wxqi *WXReqInvoker) SendGetContactRequest(userInfoList []string, antisPanTicketList []string, chatRoomWxidList []string, needResp bool) (*wechat.GetContactResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetContactRequest(tmpUserInfo, userInfoList, antisPanTicketList, chatRoomWxidList)
	if err != nil {

		return nil, err
	}

	if needResp {
		response := &wechat.GetContactResponse{}
		// 解析token登陆响应
		err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
		return response, err
	}

	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil, nil
}

// SendGetContactRequestForHB 获取联系人信息列表
func (wxqi *WXReqInvoker) SendGetContactRequestForHB(userWxid string) (*wechat.GetContactResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	userInfoList := make([]string, 1)
	userInfoList[0] = userWxid
	packHeader, err := clientsdk.SendGetContactRequest(tmpUserInfo, userInfoList, []string{}, []string{})
	if err != nil {

		return nil, err
	}

	// 解析获取联系人响应
	getContactResp := new(wechat.GetContactResponse)
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, getContactResp)
	if err != nil {

		return nil, err
	}
	return getContactResp, nil
}

// SendGetContactRequestForList 获取联系人信息列表List
func (wxqi *WXReqInvoker) SendGetContactRequestForList(userInfoList []string, roomWxIDList []string) (*wechat.GetContactResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetContactRequest(tmpUserInfo, userInfoList, []string{}, roomWxIDList)
	if err != nil {

		return nil, err
	}

	// 解析获取联系人响应
	getContactResp := new(wechat.GetContactResponse)
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, getContactResp)
	if err != nil {

		return nil, err
	}
	return getContactResp, nil
}

// 获取好友关系状态
func (wxqi *WXReqInvoker) SendGetFriendRelationRequest(userName string) (*wechat.MMBizJsApiGetUserOpenIdResponse, error) {
	// 发送请求
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetFriendRelationReq(tmpUserInfo, userName)
	if err != nil {

		return nil, err
	}
	// 解析获取联系人响应
	resp := new(wechat.MMBizJsApiGetUserOpenIdResponse)
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, resp)
	if err != nil {

		return nil, err
	}
	return resp, nil
}

// 创建红包
func (wxqi *WXReqInvoker) SendWXCreateRedPacketRequest(hbItem *baseinfo.RedPacket) (*wechat.HongBaoRes, error) {
	// 开始
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendWXCreateRedPacket(wxqi.wxAccount.GetUserInfo(), hbItem)
	if err != nil {

		return nil, err
	}
	var hongbaoResp wechat.HongBaoRes
	errs := clientsdk.ParseResponseData(tmpUserInfo, packHeader, &hongbaoResp)
	if errs != nil {
		return nil, errs
	}
	return &hongbaoResp, nil
}

// SendReceiveWxHBRequest 拆红包
func (wxqi *WXReqInvoker) SendOpenRedEnvelopesRequest(hbItem *baseinfo.HongBaoItem) (*wechat.HongBaoRes, error) {
	// 开始接收红包
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	currentProfile := wxqi.wxAccount.GetUserProfile()
	hongBaoReceiverItem := new(baseinfo.HongBaoReceiverItem)
	hongBaoReceiverItem.CgiCmd = baseinfo.MMHongBaoReqCgiCmdReceiveWxhb       // 接收红包
	hongBaoReceiverItem.Province = currentProfile.GetUserInfo().GetProvince() // 当前帐号设置的省
	hongBaoReceiverItem.City = currentProfile.GetUserInfo().GetCity()         // 当前帐号设置的市
	hongBaoReceiverItem.InWay = baseinfo.MMHongBaoReqInAwayGroup              // 群红包
	hongBaoReceiverItem.NativeURL = hbItem.NativeURL                          // nativeurl
	hongBaoReceiverItem.HongBaoURLItem = hbItem.URLItem
	packHeader, err := clientsdk.SendReceiveWxHB(wxqi.wxAccount.GetUserInfo(), hongBaoReceiverItem)
	if err != nil {

		return nil, err
	}
	var hongbaoResp wechat.HongBaoRes
	errs := clientsdk.ParseResponseData(tmpUserInfo, packHeader, &hongbaoResp)
	if errs != nil {
		return nil, errs
	}
	// 解析
	retHongBaoReceiveResp := &baseinfo.HongBaoReceiverResp{}
	err = json.Unmarshal(hongbaoResp.GetRetText().GetBuffer(), retHongBaoReceiveResp)
	if err != nil {
		return nil, err
	}
	// 发送给微信消息处理器
	rsp, er := wxqi.SendOpenWxHBNewRequest(hbItem, retHongBaoReceiveResp.TimingIdentifier)
	if er != nil {
		return nil, er
	}
	return rsp, nil
}

// SendOpenWxHBRequest 打开红包
func (wxqi *WXReqInvoker) SendOpenWxHBNewRequest(hbItem *baseinfo.HongBaoItem, timingIdentifier string) (*wechat.HongBaoRes, error) {
	// 构造参数
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	currentProfile := wxqi.wxAccount.GetUserProfile()
	hongBaoOpenItem := new(baseinfo.HongBaoOpenItem)
	hongBaoOpenItem.CgiCmd = baseinfo.MMHongBaoReqCgiCmdOpenWxhb                   // 打开红包
	hongBaoOpenItem.Province = currentProfile.GetUserInfo().GetProvince()          // 当前帐号设置的省
	hongBaoOpenItem.City = currentProfile.GetUserInfo().GetCity()                  // 当前帐号设置的市
	hongBaoOpenItem.NativeURL = hbItem.NativeURL                                   // nativeurl
	hongBaoOpenItem.HongBaoURLItem = hbItem.URLItem                                // 解析URL的各个字段信息
	hongBaoOpenItem.HeadImg = currentProfile.GetUserInfoExt().GetSmallHeadImgUrl() // 当前帐号的小头像URL
	hongBaoOpenItem.TimingIdentifier = timingIdentifier                            // 接收红包返回的
	hongBaoOpenItem.NickName = currentProfile.GetUserInfo().GetNickName().GetStr() // 当前帐号的昵称
	packHeader, err := clientsdk.SendOpenWxHB(wxqi.wxAccount.GetUserInfo(), hongBaoOpenItem)
	if err != nil {

		return nil, err
	}
	//wechat.HongBaoRes{}
	// 解析获取联系人响应
	resp := new(wechat.HongBaoRes)
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, resp)
	if err != nil {

		return nil, err
	}
	return resp, nil
}

// SendReceiveWxHBRequest 接收红包
func (wxqi *WXReqInvoker) SendReceiveWxHBRequest(hbItem *baseinfo.HongBaoItem) error {
	// 开始接收红包
	currentProfile := wxqi.wxAccount.GetUserProfile()
	hongBaoReceiverItem := new(baseinfo.HongBaoReceiverItem)
	hongBaoReceiverItem.CgiCmd = baseinfo.MMHongBaoReqCgiCmdReceiveWxhb       // 接收红包
	hongBaoReceiverItem.Province = currentProfile.GetUserInfo().GetProvince() // 当前帐号设置的省
	hongBaoReceiverItem.City = currentProfile.GetUserInfo().GetCity()         // 当前帐号设置的市
	hongBaoReceiverItem.InWay = baseinfo.MMHongBaoReqInAwayGroup              // 群红包
	hongBaoReceiverItem.NativeURL = hbItem.NativeURL                          // nativeurl
	hongBaoReceiverItem.HongBaoURLItem = hbItem.URLItem
	_, err := clientsdk.SendReceiveWxHB(wxqi.wxAccount.GetUserInfo(), hongBaoReceiverItem)
	if err != nil {

		return err
	}
	// 发送给微信消息处理器
	return nil
}

// SendOpenWxHBRequest 打开红包
func (wxqi *WXReqInvoker) SendOpenWxHBRequest(hbItem *baseinfo.HongBaoItem, timingIdentifier string) error {
	// 构造参数
	currentProfile := wxqi.wxAccount.GetUserProfile()
	hongBaoOpenItem := new(baseinfo.HongBaoOpenItem)
	hongBaoOpenItem.CgiCmd = baseinfo.MMHongBaoReqCgiCmdOpenWxhb                   // 打开红包
	hongBaoOpenItem.Province = currentProfile.GetUserInfo().GetProvince()          // 当前帐号设置的省
	hongBaoOpenItem.City = currentProfile.GetUserInfo().GetCity()                  // 当前帐号设置的市
	hongBaoOpenItem.NativeURL = hbItem.NativeURL                                   // nativeurl
	hongBaoOpenItem.HongBaoURLItem = hbItem.URLItem                                // 解析URL的各个字段信息
	hongBaoOpenItem.HeadImg = currentProfile.GetUserInfoExt().GetSmallHeadImgUrl() // 当前帐号的小头像URL
	hongBaoOpenItem.TimingIdentifier = timingIdentifier                            // 接收红包返回的
	hongBaoOpenItem.NickName = currentProfile.GetUserInfo().GetNickName().GetStr() // 当前帐号的昵称
	packHeader, err := clientsdk.SendOpenWxHB(wxqi.wxAccount.GetUserInfo(), hongBaoOpenItem)
	if err != nil {

		return err
	}
	//wechat.HongBaoRes{}
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(wxqi.wxAccount.GetUserInfo().UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendRedEnvelopesDetailRequest 查看红包详情
func (wxqi *WXReqInvoker) SendRedEnvelopesDetailRequest(hbItem *baseinfo.HongBaoItem) (*wechat.HongBaoRes, error) {
	// 构造参数
	currentProfile := wxqi.wxAccount.GetUserProfile()
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	hongBaoOpenItem := new(baseinfo.HongBaoOpenItem)
	hongBaoOpenItem.CgiCmd = 5181                                                  // 查看红包
	hongBaoOpenItem.Province = currentProfile.GetUserInfo().GetProvince()          // 当前帐号设置的省
	hongBaoOpenItem.City = currentProfile.GetUserInfo().GetCity()                  // 当前帐号设置的市
	hongBaoOpenItem.NativeURL = hbItem.NativeURL                                   // nativeurl
	hongBaoOpenItem.HongBaoURLItem = hbItem.URLItem                                // 解析URL的各个字段信息
	hongBaoOpenItem.HeadImg = currentProfile.GetUserInfoExt().GetSmallHeadImgUrl() // 当前帐号的小头像URL
	hongBaoOpenItem.NickName = currentProfile.GetUserInfo().GetNickName().GetStr() // 当前帐号的昵称
	packHeader, err := clientsdk.SendRedEnvelopeWxHB(wxqi.wxAccount.GetUserInfo(), hongBaoOpenItem)
	if err != nil {

		return nil, err
	}
	resp := new(wechat.HongBaoRes)
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, resp)
	if err != nil {

		return nil, err
	}
	return resp, nil
}

// SendGetRedPacketListRequest 查看红包领取列表
func (wxqi *WXReqInvoker) SendGetRedPacketListRequest(hbItem *baseinfo.GetRedPacketList) (*wechat.HongBaoRes, error) {
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetRedPacketListRequest(wxqi.wxAccount.GetUserInfo(), hbItem)
	if err != nil {

		return nil, err
	}
	resp := new(wechat.HongBaoRes)
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, resp)
	if err != nil {

		return nil, err
	}
	return resp, nil
}

// 发送图片
func (wxqi *WXReqInvoker) SendUploadImageNewRequest(imgData []byte, toUserName string) (*wechat.UploadMsgImgResponse, error) {
	// 发送消息
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendUploadImageNewRequest(tmpUserInfo, imgData, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.UploadMsgImgResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SendOplogRequest 发送Oplog请求
func (wxqi *WXReqInvoker) SendOplogRequest(modifyItems []*baseinfo.ModifyItem) error {
	// 发送消息
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendOplogRequest(tmpUserInfo, modifyItems)
	if err != nil {

		return err
	}

	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// 发送企业oplog
func (wxqi *WXReqInvoker) SendQWOpLogRequest(cmdId int64, value []byte) error {
	// 发送消息
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWOpLogRequest(tmpUserInfo, cmdId, value)
	if err != nil {

		return err
	}
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendGetQRCodeRequest 获取群/个人二维码
func (wxqi *WXReqInvoker) SendGetQRCodeRequest(userName string) error {
	// 发送消息
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetQRCodeRequest(tmpUserInfo, userName)
	if err != nil {

		return err
	}
	// 发送给微信消息处理器
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendFavSyncRequest 同步收藏
func (wxqi *WXReqInvoker) SendFavSyncRequest() (interface{}, error) {
	// 同步收藏
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendFavSyncRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	dealRouter := wxqi.wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
	return dealRouter.Handle(wxResp)
}

// 获取收藏list
func (wxqi *WXReqInvoker) SendFavSyncListRequestResult(keyBuf string) (*wechat.SyncResponse, error) {
	// 同步收藏
	wxqi.SendFavSyncRequest()
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendFavSyncListRequest(tmpUserInfo, keyBuf)
	if err != nil {

		return nil, err
	}
	wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	var response wechat.FavSyncResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &response)
	if err != nil {
		return nil, err
	}
	var List []wechat.AddFavItem
	for _, v := range response.CmdList.ItemList {
		if *v.CmdId == uint32(200) {
			var data wechat.AddFavItem
			_ = proto.Unmarshal(v.CmdBuf.Data, &data)
			List = append(List, data)
		}
	}
	rep := &wechat.SyncResponse{
		Ret:    *response.Ret,
		List:   List,
		KeyBuf: *response.KeyBuf,
	}
	return rep, nil
}

// SendGetFavInfoRequest 获取收藏信息
func (wxqi *WXReqInvoker) SendGetFavInfoRequest() error {

	// 同步收藏
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetFavInfoRequest(tmpUserInfo)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}
func (wxqi *WXReqInvoker) SendGetFavInfoRequestResult() (*wechat.GetFavInfoResponse, error) {

	// 同步收藏
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetFavInfoRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	var response wechat.GetFavInfoResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
	if err != nil {
		return nil, err
	}
	// 发送给路由处理
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return &response, nil
}

// SendBatchDelFavItemRequest 删除收藏
func (wxqi *WXReqInvoker) SendBatchDelFavItemRequest(favID uint32) error {

	// 删除单条收藏详情
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendBatchDelFavItemRequest(tmpUserInfo, favID)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}
func (wxqi *WXReqInvoker) SendBatchDelFavItemRequestResult(favID uint32) (*wechat.BatchDelFavItemResponse, error) {

	// 删除单条收藏详情
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendBatchDelFavItemRequest(tmpUserInfo, favID)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	var response wechat.BatchDelFavItemResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SendGetCDNDnsRequest 获取CdnDns信息
func (wxqi *WXReqInvoker) SendGetCDNDnsRequest() (interface{}, error) {

	// 获取CdnDns信息
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetCDNDnsRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	dealRouter := wxqi.wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
	return dealRouter.Handle(wxResp)
}

// 上报设备
func (wxqi *WXReqInvoker) SendReportstrategyRequest() (*wechat.GetReportStrategyResp, error) {
	// 获取CdnDns信息
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendReportstrategyRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	var resp wechat.GetReportStrategyResp
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 设置朋友圈可见天数
func (wxqi *WXReqInvoker) SetFriendCircleDays(postItem *model.SetFriendCircleDaysModel) error {
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, _ := clientsdk.SendGetProfileRequest(tmpUserInfo)
	userInfo := &wechat.GetProfileResponse{}
	clientsdk.ParseResponseData(tmpUserInfo, packHeader, userInfo)
	userInfo.UserInfoExt.SnsUserInfo.SnsFlagex = proto.Uint32(postItem.Function)
	userInfo.UserInfoExt.SnsUserInfo.SnsPrivacyRecent = proto.Uint32(postItem.Value)
	reqBuf, _ := proto.Marshal(userInfo.UserInfoExt.SnsUserInfo)
	modifyItem := baseinfo.ModifyItem{
		CmdID: 0x33,
		Len:   uint32(len(reqBuf)),
		Data:  reqBuf,
	}
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{&modifyItem})
}

// SendSnsPostRequest 发送朋友圈
func (wxqi *WXReqInvoker) SendSnsPostRequestNew(postItem *baseinfo.SnsPostItem) (*wechat.SnsPostResponse, error) {

	// 发送朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsPostRequest(tmpUserInfo, postItem)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	/*wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)*/
	var snsPostResp wechat.SnsPostResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &snsPostResp)
	if err != nil {
		return nil, err
	}

	return &snsPostResp, nil
}

// SendSnsPostRequest 发送朋友圈
func (wxqi *WXReqInvoker) SendSnsPostRequest(postItem *baseinfo.SnsPostItem) error {

	// 发送朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsPostRequest(tmpUserInfo, postItem)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendSnsObjectOpRequest 操作朋友圈
func (wxqi *WXReqInvoker) SendSnsObjectOpRequest(opItems []*baseinfo.SnsObjectOpItem) (*wechat.SnsObjectOpResponse, error) {

	// 操作朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsObjectOpRequest(tmpUserInfo, opItems)
	if err != nil {

		return nil, err
	}

	response := &wechat.SnsObjectOpResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	return response, err
}

// SendSnsPostRequestByXML 同步转发朋友圈
func (wxqi *WXReqInvoker) SendSnsPostRequestByXML(timeLineObj *baseinfo.TimelineObject, blackList []string) error {

	// 转发朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsPostRequestByXML(tmpUserInfo, timeLineObj, blackList)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendSnsUserPageRequest 获取指定好友朋友圈
func (wxqi *WXReqInvoker) SendSnsUserPageRequest(userName string, firstPageMd5 string, maxID uint64, needResp bool) (*wechat.SnsUserPageResponse, error) {

	// 获取指定好友朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsUserPageRequest(tmpUserInfo, userName, firstPageMd5, maxID)
	if err != nil {
		//if packHeader != nil && packHeader.RetCode == baseinfo.MMRequestRetSessionTimeOut {
		//	// 断开链接, 发送token登陆
		//	_,err = bizcgi.SendAutoAuthRequest(wxqi.wxconn)
		//	if err != nil {
		//		return nil, err
		//	}
		//	packHeader, err = clientsdk.SendSnsUserPageRequest(tmpUserInfo, userName, firstPageMd5, maxID)
		//} else {
		//	return nil, err
		//}
		return nil, err
	}

	// 是否需要自己解析
	if needResp {
		response := &wechat.SnsUserPageResponse{}
		// 解析token登陆响应
		err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
		return response, err
	}
	// 如果不需要则发给路由去解析
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil, nil
}

// SendSnsObjectDetailRequest 获取指定的朋友圈详情
func (wxqi *WXReqInvoker) SendSnsObjectDetailRequest(snsID uint64) (*wechat.SnsObject, error) {

	// 获取指定好友朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsObjectDetailRequest(tmpUserInfo, snsID)
	if err != nil {

		return nil, err
	}
	response := &wechat.SnsObjectDetailResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response.GetObject(), nil
}

// SendSnsSyncRequest 同步朋友圈
func (wxqi *WXReqInvoker) SendSnsSyncRequest() error {

	// 同步朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsSyncRequest(tmpUserInfo)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendSnsTimeLineRequest 获取朋友圈首页
func (wxqi *WXReqInvoker) SendSnsTimeLineRequest(firstPageMD5 string, maxID uint64) error {

	// 取朋友圈首页
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsTimeLineRequest(tmpUserInfo, firstPageMD5, maxID)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}
func (wxqi *WXReqInvoker) SendSnsTimeLineRequestResult(firstPageMD5 string, maxID uint64) (*wechat.SnsTimeLineResponse, error) {

	// 取朋友圈首页
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsTimeLineRequest(tmpUserInfo, firstPageMD5, maxID)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	var response wechat.SnsTimeLineResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SendSnsCommentRequest 发送评论/点赞请求
func (wxqi *WXReqInvoker) SendSnsCommentRequest(commentItem *baseinfo.SnsCommentItem) error {

	// 点赞/评论
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSnsCommentRequest(tmpUserInfo, commentItem)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendGetContactLabelListRequest 获取联系人标签列表
func (wxqi *WXReqInvoker) SendGetContactLabelListRequest(needResp bool) (*wechat.GetContactLabelListResponse, error) {

	// 获取标签列表
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetContactLabelListRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	// 自己解析
	if needResp {
		wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
		var response wechat.GetContactLabelListResponse
		err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil, nil
}

// SendAddContactLabelRequest 添加标签
func (wxqi *WXReqInvoker) SendAddContactLabelRequest(newLabelList []string, needResp bool) (*wechat.AddContactLabelResponse, error) {

	// 添加标签列表
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendAddContactLabelRequest(tmpUserInfo, newLabelList)
	if err != nil {

		return nil, err
	}
	// 自己解析
	if needResp {
		wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
		var response wechat.AddContactLabelResponse
		err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil, nil
}

// SendDelContactLabelRequest 删除标签
func (wxqi *WXReqInvoker) SendDelContactLabelRequest(labelId string) (*wechat.DelContactLabelResponse, error) {

	// 添加标签列表
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendDelContactLabelRequest(tmpUserInfo, labelId)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	var response wechat.DelContactLabelResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

// SendModifyLabelRequest 修改标签
func (wxqi *WXReqInvoker) SendModifyLabelRequest(userLabelList []baseinfo.UserLabelInfoItem) (*wechat.ModifyContactLabelListResponse, error) {

	// 添加标签列表
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendModifyLabelRequest(tmpUserInfo, userLabelList)
	if err != nil {

		return nil, err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	var response wechat.ModifyContactLabelListResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, wxResp.GetPackHeader(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SendBindQueryNewRequest 查询钱包信息
func (wxqi *WXReqInvoker) SendBindQueryNewRequest(reqItem *baseinfo.TenPayReqItem) error {

	// 添加标签列表
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendBindQueryNewRequest(tmpUserInfo, reqItem)
	if err != nil {

		return err
	}
	wxResp := wxcore.NewWXResponse(tmpUserInfo.UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxqi.wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// 获取银行卡信息
func (wxqi *WXReqInvoker) SendBandCardRequest(reqItem *baseinfo.TenPayReqItem) (*wechat.TenPayResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendBindQueryNewRequest(tmpUserInfo, reqItem)
	if err != nil {

		return nil, err
	}
	var response wechat.TenPayResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// 支付方法
func (wxqi *WXReqInvoker) SendTenPayRequest(reqItem *baseinfo.TenPayReqItem) (*wechat.TenPayResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendTenPayRequest(tmpUserInfo, reqItem)
	if err != nil {

		return nil, err
	}
	var response wechat.TenPayResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SendCdnDownloadReuqest Cdn下载请求
func (wxqi *WXReqInvoker) SendCdnDownloadReuqest(downItem *baseinfo.DownMediaItem) (*baseinfo.CdnDownloadResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	return clientsdk.SendCdnDownloadReuqest(tmpUserInfo, downItem)
}

// 获取图片
func (wxqi *WXReqInvoker) GetMsgBigImg(m model.GetMsgBigImgModel) (*wechat.GetMsgImgResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, error := clientsdk.GetMsgBigImg(tmpUserInfo, m)
	if error != nil {
		return nil, error
	}
	var response wechat.GetMsgImgResponse
	err := clientsdk.ParseResponseData(tmpUserInfo, packHeader, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SendCdnSnsUploadImageReuqest Cdn上传高清图片
func (wxqi *WXReqInvoker) SendCdnSnsUploadImageReuqest(imgData []byte) (*baseinfo.CdnSnsImageUploadResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	if tmpUserInfo.SNSDnsInfo == nil {
		wxqi.SendGetCDNDnsRequest()
	}

	return clientsdk.SendCdnSnsUploadImageReuqest(tmpUserInfo, imgData)
}

// SendCdnSnsVideoDownloadReuqest 发送CDN朋友圈视频下载请求
func (wxqi *WXReqInvoker) SendCdnSnsVideoDownloadReuqest(encKey uint64, tmpURL string) ([]byte, error) {

	// 发送朋友圈
	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	return clientsdk.SendCdnSnsVideoDownloadReuqest(tmpUserInfo, encKey, tmpURL)
}

// SendCdnSnsVideoUploadReuqest 发送CDN朋友圈上传视频请求
func (wxqi *WXReqInvoker) SendCdnSnsVideoUploadReuqest(videoData []byte, thumbData []byte) (*baseinfo.CdnSnsVideoUploadResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	if tmpUserInfo.DNSInfo == nil {
		_, _ = wxqi.SendGetCDNDnsRequest()
	}
	return clientsdk.SendCdnSnsVideoUploadReuqest(tmpUserInfo, videoData, thumbData)
}

// SendCdnUploadImageReuqest 发送图片给文件助手
func (wxqi *WXReqInvoker) SendCdnUploadImageReuqest(imgData []byte, toUserName string) (bool, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	if tmpUserInfo.DNSInfo == nil {
		_, _ = wxqi.SendGetCDNDnsRequest()
	}

	return clientsdk.SendCdnUploadImageReuqest(tmpUserInfo, toUserName, imgData)
}

// 发送图片
func (wxqi *WXReqInvoker) SendCdnUploadImageRequest(imgData []byte, toUserName string) (bool, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	if tmpUserInfo.DNSInfo == nil {
		_, _ = wxqi.SendGetCDNDnsRequest()
	}

	return clientsdk.SendCdnUploadImageReuqest(tmpUserInfo, toUserName, imgData)
}

func (wxqi *WXReqInvoker) SendCdnUploadVideoRequest(toUserName string, imgData string, videoData []byte) (*baseinfo.CdnMsgVideoUploadResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	if tmpUserInfo.DNSInfo == nil {
		_, _ = wxqi.SendGetCDNDnsRequest()
	}

	return clientsdk.SendCdnUploadVideoRequest(tmpUserInfo, toUserName, imgData, videoData)
}

// SendImageToFileHelper 发送图片给文件助手
func (wxqi *WXReqInvoker) SendImageToFileHelper(imgData []byte) (bool, error) {

	return wxqi.SendCdnUploadImageReuqest(imgData, baseinfo.FileHelperWXID)
}

// ForwardCdnImageRequest 转发Cdn图片
func (wxqi *WXReqInvoker) ForwardCdnImageRequest(item baseinfo.ForwardImageItem) (*wechat.UploadMsgImgResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.ForwardCdnImageRequest(tmpUserInfo, item)
	if err != nil {

		return nil, err
	}

	response := &wechat.UploadMsgImgResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ForwardCdnVideoRequest 转发Cdn视频
func (wxqi *WXReqInvoker) ForwardCdnVideoRequest(item baseinfo.ForwardVideoItem) (*wechat.UploadVideoResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.ForwardCdnVideoRequest(tmpUserInfo, item)
	if err != nil {

		return nil, err
	}

	response := &wechat.UploadVideoResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendAppMessage 发送app信息
func (wxqi *WXReqInvoker) SendAppMessage(msgXml, toUSerName string, contentType uint32) (*wechat.SendAppMsgResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	// 发送app
	packHeader, err := clientsdk.SendAppMsgRequest(tmpUserInfo, contentType, toUSerName, msgXml)
	if err != nil {

		return nil, err
	}

	response := &wechat.SendAppMsgResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendEmojiRequest 发送表情
func (wxqi *WXReqInvoker) SendEmojiRequest(md5 string, toUSerName string, length int32) (*wechat.SendAppMsgResponse, error) {

	// 生成xml
	//msgXMl := `<appmsg appid=""  sdkver="0"><title></title><des></des><action></action><type>8</type><showtype>0</showtype><soundtype>0</soundtype><mediatagname></mediatagname><messageext></messageext><messageaction></messageaction><content></content><contentattr>0</contentattr><url></url><lowurl></lowurl><dataurl></dataurl><lowdataurl></lowdataurl><songalbumurl></songalbumurl><songlyric></songlyric><appattach><totallen>1060941</totallen><attachid>0:0:ce58baf1002411bdafd299a689cadfe4</attachid><emoticonmd5>ce58baf1002411bdafd299a689cadfe4</emoticonmd5><fileext>pic</fileext><cdnthumbaeskey></cdnthumbaeskey><aeskey></aeskey></appattach><extinfo></extinfo><sourceusername></sourceusername><sourcedisplayname></sourcedisplayname><thumburl></thumburl><md5></md5><statextstr></statextstr><directshare>0</directshare></appmsg><fromusername></fromusername>` //clientsdk.CreateSendEmojiMsgXMl(md5, length)
	msgXMl := clientsdk.CreateSendEmojiMsgXMl(md5, length)
	return wxqi.SendAppMessage(msgXMl, toUSerName, 8)
}

// 发送表情new
func (wxqi *WXReqInvoker) ForwardEmojiRequest(md5 string, toUSerName string, length int32) (*wechat.UploadEmojiResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.ForwardEmojiRequest(tmpUserInfo, toUSerName, md5, length)
	if err != nil {

		return nil, err
	}
	response := &wechat.UploadEmojiResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 下载语音
func (wxqi *WXReqInvoker) SendGetMsgVoiceRequest(toUserName, newMsgId, bufid string, length int) (*vo.DownloadVoiceData, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	return clientsdk.SendGetMsgVoiceRequest(tmpUserInfo, toUserName, newMsgId, bufid, length)
}

// 群发文字
func (wxqi *WXReqInvoker) SendGroupMassMsgTextRequest(toUSerName []string, content string) (*wechat.MassSendResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGroupMassMsgText(tmpUserInfo, toUSerName, content)
	if err != nil {

		return nil, err
	}
	response := &wechat.MassSendResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 群发图片
func (wxqi *WXReqInvoker) SendGroupMassMsgImageRequest(toUSerName []string, ImageBase64 []byte) (*wechat.MassSendResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGroupMassMsgImage(tmpUserInfo, toUSerName, ImageBase64)
	if err != nil {

		return nil, err
	}
	response := &wechat.MassSendResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 群拍一拍
func (wxqi *WXReqInvoker) SendSendPatRequest(chatRoomName string, toUserName string, scene int64) (*wechat.SendPatResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendSendPatRequest(tmpUserInfo, chatRoomName, toUserName, scene)
	if err != nil {

		return nil, err
	}
	response := &wechat.SendPatResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SetChatRoomAnnouncementRequest 设置群公告
func (wxqi *WXReqInvoker) SetChatRoomAnnouncementRequest(roomId, content string) (*wechat.SetChatRoomAnnouncementResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.SetChatRoomAnnouncementRequest(tmpUserInfo, roomId, content)
	if err != nil {

		return nil, err
	}

	response := &wechat.SetChatRoomAnnouncementResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetChatroomMemberDetailRequest 获取群成员
func (wxqi *WXReqInvoker) GetChatroomMemberDetailRequest(roomId string) (*wechat.GetChatroomMemberDetailResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetChatroomMemberDetailRequest(tmpUserInfo, roomId)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetChatroomMemberDetailResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// 获取群详细
func (wxqi *WXReqInvoker) SetGetChatRoomInfoDetailRequest(roomId string) (*wechat.GetChatRoomInfoDetailRequest, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SetGetChatRoomInfoDetailRequest(tmpUserInfo, roomId)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetChatRoomInfoDetailRequest{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// 添加群管理
func (wxqi *WXReqInvoker) SendAddChatroomAdminRequest(chatRoomName string, userList []string) (*wechat.AddChatRoomAdminResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendAddChatroomAdmin(tmpUserInfo, chatRoomName, userList)
	if err != nil {

		return nil, err
	}

	response := &wechat.AddChatRoomAdminResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// 删除群管理
func (wxqi *WXReqInvoker) SendDelChatroomAdminRequest(chatRoomName string, userList []string) (*wechat.DelChatRoomAdminResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.SendDelChatroomAdminRequest(tmpUserInfo, chatRoomName, userList)
	if err != nil {

		return nil, err
	}
	response := &wechat.DelChatRoomAdminResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 获取群列表
func (wxqi *WXReqInvoker) SendWXSyncContactRequest() (*vo.GroupData, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	list := make([]wechat.ModContact, 0)
	v := int64(0)
	_ = initGroup(wxqi, tmpUserInfo, []byte(""), &list, v)
	return &vo.GroupData{
		Count: int64(len(list)),
		List:  list,
	}, nil
}

func initGroup(wxqi *WXReqInvoker, tmpUserInfo *baseinfo.UserInfo, key []byte, list *[]wechat.ModContact, v int64) error {
	packHeader, err := clientsdk.SendWXSyncContactRequest(tmpUserInfo, key)
	if err != nil {

		return err
	}
	response := &wechat.NewSyncResponse{}
	// 解析token登陆响应 -cmdList.list
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return err
	}
	i := len(*list)
	wxId := tmpUserInfo.WxId
	for _, item := range response.CmdList.ItemList {
		if item.GetCmdId() == baseinfo.CmdIDModContact || item.GetCmdId() == baseinfo.CmdIDAddMsg {
			contact := new(wechat.ModContact)
			err := proto.Unmarshal(item.GetCmdBuf().GetData(), contact)
			if err != nil {
				continue
			}
			// 判断contact是否是群 == 0 不是群
			if contact.GetChatroomVersion() == 0 {
				continue
			}
			//被移除群聊
			if contact.GetChatRoomNotify() == 0 {
				log.Println("消息免打扰群, 群wxid = ", contact.GetUserName().GetStr(), " 群昵称：", contact.GetNickName().GetStr())
			} else {
				log.Println("微信群, 群wxid = ", contact.GetUserName().GetStr(), " 群昵称：", contact.GetNickName().GetStr())
			}
			userName := contact.GetUserName().GetStr()
			if strings.HasSuffix(userName, "@chatroom") {
				add := false
				if contact.NewChatroomData.GetMemberCount() > 0 {
					for _, v := range contact.NewChatroomData.ChatroomMemberList {
						if wxId == v.GetUserName() {
							add = true
						}
					}
				}
				if add {
					*list = append(*list, *contact)
				}
			}
		}
	}
	v++
	if len(*list) > i || (len(*list) == 0 && v <= 1) {
		key = response.KeyBuf.Buffer
		tmpUserInfo.SyncKey = key
		_ = initGroup(wxqi, tmpUserInfo, key, list, v)
	}
	return nil
}

// SendCreateChatRoomRequest 创建群请求
func (wxqi *WXReqInvoker) SendCreateChatRoomRequest(topIc string, userList []string) (*wechat.CreateChatRoomResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetCreateChatRoomEntity(tmpUserInfo, topIc, userList)
	if err != nil {

		return nil, err
	}

	response := &wechat.CreateChatRoomResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendDelDelChatRoomMemberRequest 删除群成员
func (wxqi *WXReqInvoker) SendDelDelChatRoomMemberRequest(chatRoomName string, delUserList []string) (*wechat.DelChatRoomMemberResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.DelDelChatRoomMemberRequest(tmpUserInfo, chatRoomName, delUserList)
	if err != nil {

		return nil, err
	}

	response := &wechat.DelChatRoomMemberResponse{}
	// 解析响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendTransferGroupOwnerRequest 转让群
func (wxqi *WXReqInvoker) SendTransferGroupOwnerRequest(chatRoomName, newOwnerUserName string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetTransferGroupOwnerRequest(tmpUserInfo, chatRoomName, newOwnerUserName)
	if err != nil {

		return nil, err
	}

	response := &wechat.TransferChatRoomOwnerResponse{}
	// 解析响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetQuitChatroomRequest 退出群聊请求
func (wxqi *WXReqInvoker) GetQuitChatroomRequest(chatRoomName string) error {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	qutiChatRoomItem := clientsdk.CreateQutiChatRoomItem(chatRoomName, tmpUserInfo.GetUserName())
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{qutiChatRoomItem})
}

// GetInviteChatroomMembersRequest 邀请群成员
func (wxqi *WXReqInvoker) SendInviteChatroomMembersRequest(chatRoomName string, userList []string) (*wechat.CreateChatRoomResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetInviteChatroomMembersRequest(tmpUserInfo, chatRoomName, userList)
	if err != nil {

		return nil, err
	}

	response := &wechat.CreateChatRoomResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendAddChatRoomMemberRequest 发送拉人请求
func (wxqi *WXReqInvoker) SendAddChatRoomMemberRequest(chatRoomName string, userList []string) (*wechat.AddChatRoomMemberResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetAddChatRoomMemberRequest(tmpUserInfo, chatRoomName, userList)
	if err != nil {

		return nil, err
	}

	response := &wechat.AddChatRoomMemberResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetA8KeyRequest 授权链接
func (wxqi *WXReqInvoker) GetA8KeyRequest(opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*wechat.GetA8KeyResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetA8KeyRequest(tmpUserInfo, opCode, scene, reqUrl, getType)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetA8KeyResp{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetA8KeyRequest 授权进群链接
func (wxqi *WXReqInvoker) GetA8KeyGroupRequest(opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*wechat.GetA8KeyResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetA8KeyGroupRequest(tmpUserInfo, opCode, scene, reqUrl, getType)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetA8KeyResp{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// JSLoginRequest 小程序授权
func (wxqi *WXReqInvoker) JSLoginRequest(appId string) (*wechat.JSLoginResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.JSLoginRequest(tmpUserInfo, appId)
	if err != nil {

		return nil, err
	}

	response := &wechat.JSLoginResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// JSOperateWxDataRequest
func (wxqi *WXReqInvoker) JSOperateWxDataRequest(appId string) (*wechat.JSOperateWxDataResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.JSOperateWxDataRequest(tmpUserInfo, appId)
	if err != nil {

		return nil, err
	}

	response := &wechat.JSOperateWxDataResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SdkOauthAuthorizeRequest 授权 App应用
func (wxqi *WXReqInvoker) SdkOauthAuthorizeRequest(appId string, sdkName string, packageName string) (*wechat.SdkOauthAuthorizeConfirmNewResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetSdkOauthAuthorizeConfirmRequest(tmpUserInfo, appId, sdkName, packageName)
	if err != nil {

		return nil, err
	}

	response := &wechat.SdkOauthAuthorizeConfirmNewResp{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendSearchContactRequest 搜索联系人
func (wxqi *WXReqInvoker) SendSearchContactRequest(opCode, fromScene, searchScene uint32, userName string) (*wechat.SearchContactResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.SendSearchContactRequest(tmpUserInfo, opCode, fromScene, searchScene, userName)
	if err != nil {

		return nil, err
	}

	response := &wechat.SearchContactResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// VerifyUserRequest 添加好友 关注公众号 同意好友添加
func (wxqi *WXReqInvoker) VerifyUserRequest(opCode uint32, verifyContent string, scene byte, V1, V2, ChatRoomUserName string) (*wechat.VerifyUserResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.VerifyUserRequest(tmpUserInfo, opCode, verifyContent, scene, V1, V2, ChatRoomUserName)
	if err != nil {

		return nil, err
	}

	response := &wechat.VerifyUserResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UploadMContact 上传手机通讯录
func (wxqi *WXReqInvoker) UploadMContact(mobile string, mobileList []string) (*wechat.UploadMContactResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.UploadMContactRequest(tmpUserInfo, mobile, mobileList)
	if err != nil {

		return nil, err
	}

	response := &wechat.UploadMContactResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetMFriend 获取手机通讯录好友
func (wxqi *WXReqInvoker) GetMFriend() (*wechat.GetMFriendResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetMFriendRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetMFriendResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (wxqi *WXReqInvoker) SendCertRequest() (*wechat.GetCertResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetCertRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetCertResponse{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendQRConnectAuthorize 发送二维码授权请求
func (wxqi *WXReqInvoker) SendQRConnectAuthorize(qrUrl string) (*wechat.QRConnectAuthorizeResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetQRConnectAuthorizeRequest(tmpUserInfo, qrUrl)
	if err != nil {

		return nil, err
	}

	response := &wechat.QRConnectAuthorizeResp{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendQRConnectAuthorize 发送二维码授权请求确认
func (wxqi *WXReqInvoker) SendQRConnectAuthorizeConfirm(qrUrl string) (*wechat.SdkOauthAuthorizeConfirmNewResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetQRConnectAuthorizeConfirmRequest(tmpUserInfo, qrUrl)
	if err != nil {

		return nil, err
	}

	response := &wechat.SdkOauthAuthorizeConfirmNewResp{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// 授权链接
func (wxqi *WXReqInvoker) SendGetMpA8Request(url string, opcode uint32) (*wechat.GetA8KeyResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetMpA8Request(tmpUserInfo, url, opcode)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetA8KeyResp{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SendOnlineInfoRequest 获取登录设备信息
func (wxqi *WXReqInvoker) SendOnlineInfo() (*wechat.GetOnlineInfoResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetOnlineInfoRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetOnlineInfoResponse{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendGetQrCode 发送获取二维码请求
func (wxqi *WXReqInvoker) SendGetQrCodeRequest(id string) (*wechat.GetQRCodeResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	if id == "" {
		id = tmpUserInfo.WxId
	}
	packHeader, err := clientsdk.GetQrCodeRequest(tmpUserInfo, id)
	if err != nil {

		return nil, err
	}

	response := &wechat.GetQRCodeResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendGetPeopleNearbyResult 查看附近的人
func (wxqi *WXReqInvoker) SendGetPeopleNearbyResultRequest(longitude float32, latitude float32) (*wechat.LbsResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetPeopleNearbyResultRequest(tmpUserInfo, longitude, latitude)
	if err != nil {

		return nil, err
	}
	response := &wechat.LbsResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SendRevokeMsgRequest 撤销消息
func (wxqi *WXReqInvoker) SendRevokeMsgRequest(newMsgId string, clientMsgId uint64, toUserName string) (*wechat.RevokeMsgResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetRevokeMsgRequest(tmpUserInfo, newMsgId, clientMsgId, toUserName)
	if err != nil {

		return nil, err
	}

	response := &wechat.RevokeMsgResponse{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendDelContactRequest 删除好友
func (wxqi *WXReqInvoker) SendDelContactRequest(userName string) error {
	resp, err := wxqi.SendGetContactRequest([]string{userName}, []string{}, []string{}, true)
	if err != nil {
		return err
	}
	if resp.GetBaseResponse().GetRet() != 0 || len(resp.GetContactList()) <= 0 {
		return errors.New("获取用户详细失败！")
	}

	modifyItem := clientsdk.CreateDeleteFriendField(resp.GetContactList()[0])
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{modifyItem})
}

// SendModifyUserInfoRequest 修改用户资料
func (wxqi *WXReqInvoker) SendModifyUserInfoRequest(city, country, nickName, province, signature string, sex uint32, initFlag uint32) error {
	modUserInfo := wxqi.wxAccount.GetUserProfile().GetUserInfo()
	if modUserInfo == nil {
		modUserInfo = &wechat.ModUserInfo{}
	} else {
		initFlag = modUserInfo.GetBitFlag()
	}
	if nickName != "" {
		modUserInfo.NickName = &wechat.SKBuiltinString{
			Str: proto.String(nickName),
		}
	}
	if city != "" {
		modUserInfo.City = proto.String(city)
	}
	if province != "" {
		modUserInfo.Province = proto.String(province)
	}
	if signature != "" {
		modUserInfo.Signature = proto.String(signature)
	}
	modUserInfo.Sex = proto.Uint32(sex)
	modifyItem := clientsdk.CreateModifyUserInfoField(modUserInfo, initFlag, nickName)
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{modifyItem})
}

// 修改昵称
func (wxqi *WXReqInvoker) SendUpdateNickNameRequest(cmd uint32, val string) error {
	newData := &wechat.ModInfo{
		Cmd:   proto.Uint32(cmd),
		Value: proto.String(val),
	}
	data, marshalErr := proto.Marshal(newData)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return marshalErr
	}
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{&baseinfo.ModifyItem{
		CmdID: uint32(64),
		Len:   uint32(len(data)),
		Data:  data,
	}})
}

// 修改姓名
func (wxqi *WXReqInvoker) SetNickNameService(cmd uint32, val string) error {
	newData := &wechat.ModInfo{
		Cmd:   proto.Uint32(cmd),
		Value: proto.String(val),
	}
	data, marshalErr := proto.Marshal(newData)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return marshalErr
	}
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{&baseinfo.ModifyItem{
		CmdID: uint32(64),
		Len:   uint32(len(data)),
		Data:  data,
	}})
}

// 修改性别
func (wxqi *WXReqInvoker) SetSexService(val uint32, country string, city string, province string) error {
	rsp, _ := wxqi.SendGetProfileNewRequest()
	modUserInfo := rsp.UserInfo
	req := wechat.ModUserInfo{
		BitFlag: proto.Uint32(2178),
		Sex:     modUserInfo.Sex,
		UserName: &wechat.SKBuiltinString{
			Str: proto.String(modUserInfo.GetUserName().GetStr()),
		},
		NickName: &wechat.SKBuiltinString{
			Str: proto.String(modUserInfo.NickName.GetStr()),
		},
		BindUin:    proto.Uint32(modUserInfo.GetBindUin()),
		BindEmail:  modUserInfo.BindEmail,
		BindMobile: modUserInfo.BindMobile,
		Status:     modUserInfo.Status,
		ImgLen:     modUserInfo.ImgLen,
		Province:   modUserInfo.Province,
		City:       modUserInfo.City,
		Signature:  modUserInfo.Signature,
		PluginFlag: modUserInfo.PluginFlag,
		Country:    modUserInfo.Country,
	}
	if modUserInfo.GetSex() != val {
		req.Sex = proto.Uint32(val)
	}
	if country != "" || city != "" || province != "" {
		req.Country = proto.String(country)
		req.City = proto.String(city)
		req.Province = proto.String(province)
	}
	data, marshalErr := proto.Marshal(&req)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.SetSexService failed: ", marshalErr)
		return marshalErr
	}
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{&baseinfo.ModifyItem{
		CmdID: uint32(1),
		Len:   uint32(len(data)),
		Data:  data,
	}})
}

// 修改头像
func (wxqi *WXReqInvoker) UploadHeadImage(base64 string) (*wechat.UploadHDHeadImgResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.UploadHeadImage(tmpUserInfo, base64)
	if err != nil {

		return nil, err
	}
	response := &wechat.UploadHDHeadImgResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 修改加好友需要验证属性
func (wxqi *WXReqInvoker) UpdateAutopassRequest(SwitchType uint32) error {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	modifyItem := clientsdk.CreateBlackSnsItem(tmpUserInfo.GetUserName(), SwitchType)
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{modifyItem})
}

func (wxqi *WXReqInvoker) verifyPwdRequest(oldPwd string) (*wechat.VerifyPwdResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.GetVerifyPwdRequest(tmpUserInfo, oldPwd)
	if err != nil {

		return nil, err
	}

	response := &wechat.VerifyPwdResponse{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (wxqi *WXReqInvoker) setPwdRequest(ticket, newPwd string, OpCode uint32) (*wechat.SetPwdResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.SendSetPwdRequest(tmpUserInfo, ticket, newPwd, OpCode)
	if err != nil {

		return nil, err
	}

	response := &wechat.SetPwdResponse{}

	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendChangePwdRequest 更改密码
func (wxqi *WXReqInvoker) SendChangePwdRequest(oldPwd, NewPwd string, OpCode uint32) (*wechat.BaseResponse, error) {
	ticket := ""
	if OpCode == 0 {
		verifyPwdResp, err := wxqi.verifyPwdRequest(oldPwd)
		if err != nil {
			return nil, err
		}
		if verifyPwdResp.GetBaseResponse().GetRet() != 0 {
			return verifyPwdResp.GetBaseResponse(), nil
		}
		ticket = verifyPwdResp.GetTicket()
	}
	resp, err := wxqi.setPwdRequest(ticket, NewPwd, OpCode)
	if err != nil {
		return nil, err
	}
	return resp.GetBaseResponse(), nil
}

// SendModifyRemarkRequest 修改备注
func (wxqi *WXReqInvoker) SendModifyRemarkRequest(userName string, remarkName string) error {
	resp, err := wxqi.SendGetContactRequest([]string{userName}, []string{}, []string{}, true)
	if err != nil {
		return err
	}
	if resp.GetBaseResponse().GetRet() != 0 || len(resp.GetContactList()) <= 0 {
		return errors.New("获取用户详细失败！")
	}
	modContact := resp.GetContactList()[0]
	modContact.GetRemark().Str = proto.String(remarkName)
	modifyItem := clientsdk.CreateModifyFriendField(modContact)
	return wxqi.SendOplogRequest([]*baseinfo.ModifyItem{modifyItem})
}

// SendUploadVoiceRequest 发送语音消息
/*func (wxqi *WXReqInvoker) SendUploadVoiceRequest(toUserName string, voiceData []byte, voiceSecond, voiceFormat uint32) (*wechat.UploadVoiceResponse, error) {


	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	// 计算语音长度
	if voiceSecond <= 0 {
		voiceSecond = (uint32)(float32(len(voiceData)) / 1.8 / 1024)
		if voiceSecond >= 60 {
			voiceSecond = 59
		}
	}
	var startPos = uint32(0)
	//mmtls结构限制，块长度ushort最大长度65535
	var dataLenth = uint32(65000)
	var totalLen = uint32(len(voiceData))
	var clientMsgId = fmt.Sprintf("%v", time.Now().Unix())
	response := &wechat.UploadVoiceResponse{}

	for startPos != totalLen {
		count := uint32(0)
		if totalLen-startPos > dataLenth {
			count = dataLenth
		} else {
			count = totalLen - startPos
		}
		updateData := voiceData[startPos : startPos+count]
		packHeader, err := clientsdk.SendUploadVoiceRequest(
			tmpUserInfo, toUserName, updateData, totalLen, startPos, clientMsgId, voiceSecond*1000, voiceFormat)
		if err != nil {

			return nil, err
		}

		err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
		if err != nil {
			return nil, err
		}
		startPos += count
	}
	return response, nil
}*/
func (wxqi *WXReqInvoker) SendUploadVoiceRequest(toUserName string, voiceData string, voiceSecond, voiceFormat int32) (*wechat.UploadVoiceResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	VoiceData := strings.Split(voiceData, ",")
	var VoiceBase64 []byte
	if len(VoiceData) > 1 {
		VoiceBase64, _ = base64.StdEncoding.DecodeString(VoiceData[1])
	} else {
		VoiceBase64, _ = base64.StdEncoding.DecodeString(voiceData)
	}
	VoiceStream := bytes.NewBuffer(VoiceBase64)
	Startpos := 0
	datalen := 65000
	datatotalength := VoiceStream.Len()
	ClientImgId := fmt.Sprintf("%s—%v", tmpUserInfo.WxId, time.Now().UnixNano())
	// 计算语音长度
	if voiceSecond <= 0 {
		voiceSecond = (int32)(float32(len(voiceData)) / 1.8 / 1024)
		if voiceSecond >= 60 {
			voiceSecond = 59
		}
	}
	for {
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}

		endFlag := 0
		if Startpos+count >= datatotalength {
			endFlag = 1
		}
		Databuff := make([]byte, count)
		_, _ = VoiceStream.Read(Databuff)
		packHeader, err := clientsdk.SendUploadVoiceNewRequest(
			tmpUserInfo, toUserName, Startpos, Databuff, ClientImgId, voiceSecond, voiceFormat, endFlag)
		if err != nil {

			return nil, err
		}
		response := &wechat.UploadVoiceResponse{}
		err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
		if err != nil {
			return nil, err
		}
		if response.GetEndFlag() == 1 {
			return response, nil
		}
		Startpos += count
	}
	return nil, nil

	/*response := &wechat.UploadVoiceResponse{}

	for startPos != totalLen {
		count := uint32(0)
		if totalLen-startPos > dataLenth {
			count = dataLenth
		} else {
			count = totalLen - startPos
		}
		updateData := voiceData[startPos : startPos+count]
		packHeader, err := clientsdk.SendUploadVoiceRequest(
			tmpUserInfo, toUserName, updateData, totalLen, startPos, clientMsgId, voiceSecond*1000, voiceFormat)
		if err != nil {

			return nil, err
		}

		err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
		if err != nil {
			return nil, err
		}
		startPos += count
	}
	return response, nil*/
}

// 设置微信号
func (wxqi *WXReqInvoker) SetWechatRequest(alisa string) (*wechat.GeneralSetResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SetWechatRequest(tmpUserInfo, alisa)
	if err != nil {

		return nil, err
	}
	response := &wechat.GeneralSetResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 设置微信步数
func (wxqi *WXReqInvoker) UpdateStepNumberRequest(number uint64) (*wechat.UploadDeviceStepResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.GetBoundHardDeviceRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetBoundHardDevicesResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	device := response.GetDeviceList()[0]
	//上传步数
	packHeaderUpd, errUpd := clientsdk.UploadStepSetRequestRequest(tmpUserInfo, device.HardDevice.GetDeviceId(), device.HardDevice.GetDeviceType(), number)
	if errUpd != nil {
		return nil, errUpd
	}
	responseUpdate := &wechat.UploadDeviceStepResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeaderUpd, responseUpdate)
	if err != nil {
		return nil, err
	}
	return responseUpdate, nil
}

// 获取步数列表
func (wxqi *WXReqInvoker) SendGetUserRankLikeCountRequest(rankId string) (*wechat.GetUserRankLikeCountResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetUserRankLikeCountRequest(tmpUserInfo, rankId)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetUserRankLikeCountResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 搜手机或企业对外名片链接提取验证
func (wxqi *WXReqInvoker) SendQWSearchContactRequest(tg string, fromScene uint64, userName string) (*wechat.SearchQYContactResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.SendQWSearchContactRequest(tmpUserInfo, tg, fromScene, userName)
	if err != nil {

		return nil, err
	}

	response := &wechat.SearchQYContactResponse{}
	// 解析token登陆响应
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 提取企业 wx 详情
func (wxqi *WXReqInvoker) SendQWContactRequest(openIm, chatRoom, t string) (*wechat.GetQYContactResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWContactRequest(tmpUserInfo, openIm, chatRoom, t)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetQYContactResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 提取全部的企业通寻录
func (wxqi *WXReqInvoker) SendQWSyncContactRequest() (*wechat.GetQYContactResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWSyncContactRequest(tmpUserInfo)
	if err != nil {

		return nil, err
	}
	response := &wechat.QYSyncRespone{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	list := make([]*wechat.ContactInfo, 0)
	for _, val := range response.List.List {
		if val.GetCmdid().String() == wechat.SyncCmdID_OpenimContact.String() {
			rps := &wechat.ContactInfo{}
			_ = proto.Unmarshal(val.Cmdg.Data, rps)
			list = append(list, rps)
		}
	}
	rsp := &wechat.GetQYContactResponse{
		Continue:    response.Continue,
		ContactList: list,
	}

	return rsp, nil
}

// 换绑手机
func (wxqi *WXReqInvoker) SendBindingMobileRequest(mobile, verifyCode string) (*wechat.BindOpMobileResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendBindingMobileRequest(mobile, verifyCode, tmpUserInfo)
	if err != nil {

		return nil, err
	}
	response := &wechat.BindOpMobileResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 发送手机验证码
func (wxqi *WXReqInvoker) SendVerifyMobileRequest(mobile string, opcode uint32) (*wechat.BindOpMobileResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendVerifyMobileRequest(mobile, opcode, tmpUserInfo)
	if err != nil {

		return nil, err
	}
	response := &wechat.BindOpMobileResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (wxqi *WXReqInvoker) SendQWRemarkRequest(toUserName string, name string) error {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()

	packHeader, err := clientsdk.SendQWSyncContactRequest(tmpUserInfo)
	if err != nil {

		return err
	}
	response := &wechat.GetQYContactResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return err
	}
	return nil
}

// 创建企业群
func (wxqi *WXReqInvoker) SendQWCreateChatRoomRequest(userList []string) (*wechat.CreateQYChatRoomResponese, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWCreateChatRoomRequest(tmpUserInfo, userList)
	if err != nil {

		return nil, err
	}
	response := &wechat.CreateQYChatRoomResponese{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 向企业微信打招呼
func (wxqi *WXReqInvoker) SendQWApplyAddContactRequest(toUserName, v1, Content string) error {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	_, err := clientsdk.SendQWApplyAddContactRequest(tmpUserInfo, toUserName, v1, Content)
	if err != nil {

		return err
	}
	return nil
}

// 单向加企业微信
func (wxqi *WXReqInvoker) SendQWAddContactRequest(toUserName, v1, Content string) error {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	_, err := clientsdk.SendQWAddContactRequest(tmpUserInfo, toUserName, v1, Content)
	if err != nil {

		return err
	}
	return nil
}

// 提取所有微信企业群
func (wxqi *WXReqInvoker) SendQWSyncChatRoomRequest(key string) (*vo.QYChatroomContactVo, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWSyncChatRoomRequest(tmpUserInfo, key)
	if err != nil {

		return nil, err
	}
	response := &wechat.QYSyncRespone{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	list := make([]*wechat.QYChatroomContactInfo, 0)

	for _, val := range response.List.List {
		if val.GetCmdid() == 0x193 {
			rps := &wechat.QYChatroomContactInfo{}
			_ = proto.Unmarshal(val.Cmdg.Data, rps)
			list = append(list, rps)
		}
	}
	rsp := &vo.QYChatroomContactVo{
		Key:  string(response.Key),
		List: list,
	}
	return rsp, nil
}

// 转让企业微信群
func (wxqi *WXReqInvoker) SendQWChatRoomTransferOwnerRequest(chatRoomName string, toUserName string) (*wechat.BaseResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWChatRoomTransferOwnerRequest(tmpUserInfo, chatRoomName, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.BaseResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 直接拉好友进群
func (wxqi *WXReqInvoker) SendQWAddChatRoomMemberRequest(chatRoomName string, toUserName []string) (*wechat.QYAddChatRoomMemberResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWAddChatRoomMemberRequest(tmpUserInfo, chatRoomName, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.QYAddChatRoomMemberResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (wxqi *WXReqInvoker) SendQWInviteChatRoomMemberRequest(chatRoomName string, toUserName []string) (*wechat.BaseResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWInviteChatRoomMemberRequest(tmpUserInfo, chatRoomName, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.BaseResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 删除企业群群成员
func (wxqi *WXReqInvoker) SendQWDelChatRoomMemberRequest(chatRoomName string, toUserName []string) (*wechat.QYDelChatRoomMemberResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWDelChatRoomMemberRequest(tmpUserInfo, chatRoomName, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.QYDelChatRoomMemberResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 提取企业群全部成员
func (wxqi *WXReqInvoker) SendQWGetChatRoomMemberRequest(chatRoomName string) (*wechat.GetQYChatroomMemberDetailResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWGetChatRoomMemberRequest(tmpUserInfo, chatRoomName)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetQYChatroomMemberDetailResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 提取企业群名称公告设定等信息
func (wxqi *WXReqInvoker) SendQWGetChatroomInfoRequest(chatRoomName string) (*wechat.QYChatroomContactResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWGetChatroomInfoRequest(tmpUserInfo, chatRoomName)
	if err != nil {

		return nil, err
	}
	response := &wechat.QYChatroomContactResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 提取企业群二维码
func (wxqi *WXReqInvoker) SendQWGetChatRoomQRRequest(chatRoomName string) (*wechat.QYGetQRCodeResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWGetChatRoomQRRequest(tmpUserInfo, chatRoomName)
	if err != nil {

		return nil, err
	}
	response := &wechat.QYGetQRCodeResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 增加企业管理员
func (wxqi *WXReqInvoker) SendQWAppointChatRoomAdminRequest(chatRoomName string, toUserName []string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWAppointChatRoomAdminRequest(tmpUserInfo, chatRoomName, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 移除企业群管理员
func (wxqi *WXReqInvoker) SendQWDelChatRoomAdminRequest(chatRoomName string, toUserName []string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWDelChatRoomAdminRequest(tmpUserInfo, chatRoomName, toUserName)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 同意进企业群
func (wxqi *WXReqInvoker) SendQWAcceptChatRoomRequest(link string, opcode uint32) (*wechat.GetA8KeyResp, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWAcceptChatRoomRequest(tmpUserInfo, link, opcode)
	if err != nil {

		return nil, err
	}
	response := &wechat.GetA8KeyResp{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 设定企业群管理审核进群
func (wxqi *WXReqInvoker) SendQWAdminAcceptJoinChatRoomSetRequest(chatRoomName string, p int64) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	req := &wechat.QYAdminAcceptJoinChatRoomSet{
		G: proto.String(chatRoomName),
		P: proto.Int64(p),
	}
	buffer, err := proto.Marshal(req)
	packHeader, err := clientsdk.SendQWOpLogRequest(tmpUserInfo, 0x10, buffer)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 群管理批准进企业群
func (wxqi *WXReqInvoker) SendQWAdminAcceptJoinChatRoomRequest(chatRoomName, key, toUserName string, toUserNames []string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendQWAdminAcceptJoinChatRoomRequest(tmpUserInfo, chatRoomName, key, toUserName, toUserNames)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 修改企业群名称
func (wxqi *WXReqInvoker) SendQWModChatRoomNameRequest(chatRoomName, name string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	req := &wechat.QYModChatRoomTopicRequest{
		G: proto.String(chatRoomName),
		P: proto.String(name),
	}
	buffer, err := proto.Marshal(req)
	packHeader, err := clientsdk.SendQWOpLogRequest(tmpUserInfo, 8, buffer)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 修改成员在群中呢称
func (wxqi *WXReqInvoker) SendQWModChatRoomMemberNickRequest(chatRoomName, name string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	req := &wechat.QYModChatRoomTopicRequest{
		G: proto.String(chatRoomName),
		P: proto.String(name),
	}
	buffer, err := proto.Marshal(req)
	packHeader, err := clientsdk.SendQWOpLogRequest(tmpUserInfo, 10, buffer)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 发布企业群公告
func (wxqi *WXReqInvoker) SendQWChatRoomAnnounceRequest(chatRoomName, Announcement string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	req := &wechat.QYModChatRoomTopicRequest{
		G: proto.String(chatRoomName),
		P: proto.String(Announcement),
	}
	buffer, err := proto.Marshal(req)
	packHeader, err := clientsdk.SendQWOpLogRequest(tmpUserInfo, 9, buffer)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 删除企业群
func (wxqi *WXReqInvoker) SendQWDelChatRoomRequest(chatRoomName string) (*wechat.TransferChatRoomOwnerResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	jsonbyte, _ := json.Marshal(chatRoomName)
	packHeader, err := clientsdk.SendQWOpLogRequest(tmpUserInfo, 14, jsonbyte)
	if err != nil {

		return nil, err
	}
	response := &wechat.TransferChatRoomOwnerResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 视频号搜索
func (wxqi *WXReqInvoker) SendGetFinderSearchRequest(Index uint32, Userver int32, UserKey string, Uuid string) (*wechat.FinderSearchResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendGetFinderSearchRequest(tmpUserInfo, Index, Userver, UserKey, Uuid)
	if err != nil {

		return nil, err
	}
	response := &wechat.FinderSearchResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 视频号个人中心
func (wxqi *WXReqInvoker) SendFinderUserPrepareRequest(uServer int32) (*wechat.FinderUserPrepareResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendFinderUserPrepareRequest(tmpUserInfo, uServer)
	if err != nil {

		return nil, err
	}
	response := &wechat.FinderUserPrepareResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 视频号关注，取消关注
func (wxqi *WXReqInvoker) SendFinderFollowRequest(FinderUserName string, OpType int32, RefObjectId string, Cook string, Userver int32, PosterUsername string) (*wechat.FinderFollowResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendFinderFollowRequest(tmpUserInfo, FinderUserName, OpType, RefObjectId, Cook, Userver, PosterUsername)
	if err != nil {

		return nil, err
	}
	response := &wechat.FinderFollowResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 查看视频号首页
func (wxqi *WXReqInvoker) TargetUserPageRequest(target string, lastBuffer string) (*wechat.FinderUserPageResponse, error) {

	tmpUserInfo := wxqi.wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendTargetUserPageRequest(tmpUserInfo, target, lastBuffer)
	if err != nil {

		return nil, err
	}
	response := &wechat.FinderUserPageResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
