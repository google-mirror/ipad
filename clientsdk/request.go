package clientsdk

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"feiyu.com/wx/api/utils"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	clientsdk "feiyu.com/wx/clientsdk/hybrid"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/protobuf/wechat"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/lunny/log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SendLoginQRCodeRequest 获取登陆二维码
func SendLoginQRCodeRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	sendData := GetLoginQRCodeReq(userInfo)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getloginqrcode", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendManualAuth 发送ManualAuth请求
func sendManualAuthByAccountData(userInfo *baseinfo.UserInfo, accountData []byte) (*baseinfo.PackHeader, error) {

	// 发送登陆请求
	sendData := GetManualAuthByAccountDataReq(userInfo, accountData)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/manualauth", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendManualAuthA16发送A16登录请求
func SendManualAuthA16(userInfo *baseinfo.UserInfo, accountData []byte) (*baseinfo.PackHeader, error) {
	// 发送登陆请求
	sendData := GetManualAuthA16Req(userInfo, accountData)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/manualauth", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendManualAuth 发送登陆请求
func SendManualAuth(userInfo *baseinfo.UserInfo, newpass string, wxid string) (*baseinfo.PackHeader, error) {
	// 序列化
	accountData, err := GetManualAuthAccountDataReq(userInfo, newpass, wxid)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") && userInfo.DeviceInfo != nil {
		return sendManualAuthByAccountData(userInfo, accountData)
	}
	return SendManualAuthA16(userInfo, accountData)
}

// 获取DeviceToken IOS
func SendIosDeviceTokenRequest(userInfo *baseinfo.UserInfo) (*wechat.TrustResp, error) {
	hecData, hec := GetIosDeviceTokenReq(userInfo)
	recvData, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/fpinitnl", hecData)
	if err != nil {
		return &wechat.TrustResp{}, err
	}
	if len(recvData) <= 31 {
		return &wechat.TrustResp{}, errors.New(hex.EncodeToString(recvData))
	}
	ph := hec.HybridEcdhPackIosUn(recvData)
	DTResp := &wechat.TrustResp{}
	_ = proto.Unmarshal(ph.Data, DTResp)
	return DTResp, nil
}

// 获取DeviceToken
func SendAndroIdDeviceTokenRequest(userInfo *baseinfo.UserInfo) (*wechat.TrustResp, error) {
	sendData, hec := GetAndroIdDeviceTokenReq(userInfo)
	recvData, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/fpinitnl", sendData)
	if err != nil {
		return &wechat.TrustResp{}, err
	}
	if len(recvData) <= 31 {
		return &wechat.TrustResp{}, errors.New(hex.EncodeToString(recvData))
	}
	ph := hec.HybridEcdhPackAndroidUn(recvData)
	DTResp := &wechat.TrustResp{}
	_ = proto.Unmarshal(ph.Data, DTResp)
	return DTResp, nil
}

// 二次登录-new
func SendSecautouthRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	//开始组头
	retData, secKeyMgr, err := GetSecautouthReq(userInfo)
	//log.Println(hex.EncodeToString(retData))
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secautoauth", retData)
	if err != nil {
		return nil, err
	}
	packHeader, err := DecodePackHeader(resp, nil)
	if err != nil {
		return nil, err
	}
	packHeader.Data, err = clientsdk.HybridEcdhDecrypt(packHeader.Data, secKeyMgr.PriKey, secKeyMgr.PubKey, secKeyMgr.FinalSha256)
	if err != nil {
		return nil, err
	}
	return packHeader, err
}

// Secautoauth二次登录
func SecautoauthRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	hecData, hec, err := GetSecautoauthReq(userInfo)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secautoauth", hecData)
	if err != nil {
		return nil, err
	}
	if len(resp) <= 31 {
		log.Info("您已退出微信/session过期")
		return nil, err
	}
	ph := hec.HybridEcdhPackIosUn(resp)
	//解包
	loginRes := wechat.UnifyAuthResponse{}
	err = proto.Unmarshal(ph.Data, &loginRes)
	return nil, nil
}

// SendAutoAuthRequest 发送token登陆请求
func SendAutoAuthRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	// 发送请求
	sendData, err := GetAutoAuthReq(userInfo)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/autoauth", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendPushQrLoginNotice 二维码二次登录
func SendPushQrLoginNotice(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	sendData := GetPushQrLoginNoticeReq(userInfo)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/pushloginurl", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetHeartBeatShortRequest 短链接心跳
func GetHeartBeatShortRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	sendData, err := GetHeartBeatReq(userInfo)
	if err != nil {
		return nil, err
	}
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/heartbeat", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 同步消息
func NewSyncHistoryMessageRequest(userInfo *baseinfo.UserInfo, scene uint32, syncKey string) (*baseinfo.PackHeader, error) {
	sendData := GetNewSyncHistoryMessageReq(userInfo, scene, syncKey)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newsync", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendNewSyncRequest 发送同步信息请求
func SendNewSyncRequest(userInfo *baseinfo.UserInfo, scene uint32) (*baseinfo.PackHeader, error) {
	sendData := GetNewSyncReq(userInfo, scene, true)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newsync", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 同步消息
func SendWxSyncMsg(userInfo *baseinfo.UserInfo, key string) (*baseinfo.PackHeader, error) {
	sendData := GetWxSyncMsgReq(userInfo, key)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newsync", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendGetProfileRequest 发送获取帐号所有信息请求
func SendGetProfileRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	sendData := GetProfileReq(userInfo)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getprofile", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取设备
func SendGetSafetyInfoRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	sendEncodeData := GetSafetyInfoReq(userInfo)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getsafetyinfo", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 删除设备
func SendDelSafeDeviceRequest(userInfo *baseinfo.UserInfo, uuid string) (*baseinfo.PackHeader, error) {
	sendEncodeData := GetDelSafeDeviceReq(userInfo, uuid)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delsafedevice", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 检测微信登录环境
func SendCheckCanSetAliasRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	baseRequest := GetBaseRequest(userInfo)
	var req = wechat.CheckCanSetAliasReq{
		BaseRequest: baseRequest,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 926, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/checkcansetalias", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 扫码登录新设备
func SendExtDeviceLoginConfirmGetRequest(userInfo *baseinfo.UserInfo, url string) (*baseinfo.PackHeader, error) {
	Url := strings.Replace(url, "https", "http", -1)
	req := &wechat.ExtDeviceLoginConfirmGetRequest{
		LoginUrl:   proto.String(Url),
		DeviceName: proto.String(baseinfo.DeviceTypeIos),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, 971, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/extdeviceloginconfirmget", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 新设备扫码确认登录
func ExtDeviceLoginConfirmOk(userInfo *baseinfo.UserInfo, url string) (*baseinfo.PackHeader, error) {
	Url := strings.Replace(url, "https", "http", -1)
	req := &wechat.ExtDeviceLoginConfirmOKRequest{
		LoginUrl:    proto.String(Url),
		SessionList: proto.String(""),
		SyncMsg:     proto.Uint64(1),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, 972, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/extdeviceloginconfirmok", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendInitContactReq 初始化联系人列表
func SendInitContactReq(userInfo *baseinfo.UserInfo, contactSeq uint32) (*baseinfo.PackHeader, error) {
	var request wechat.InitContactReq

	// Username
	request.Username = &userInfo.WxId
	// CurrentWxcontactSeq
	request.CurrentWxcontactSeq = &contactSeq
	// CurrentChatRoomContactSeq
	roomContactSeq := uint32(0)
	request.CurrentChatRoomContactSeq = &roomContactSeq

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeInitContact, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/initcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

func SendContactListPageRequest(userInfo *baseinfo.UserInfo, CurrentWxcontactSeq uint32, CurrentChatRoomContactSeq uint32) (*baseinfo.PackHeader, error) {
	var request wechat.InitContactReq

	// Username
	request.Username = &userInfo.WxId
	// CurrentWxcontactSeq
	request.CurrentWxcontactSeq = &CurrentWxcontactSeq
	request.CurrentChatRoomContactSeq = &CurrentChatRoomContactSeq

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeInitContact, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/initcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendBatchGetContactBriefInfoReq 批量获取联系人信息
func SendBatchGetContactBriefInfoReq(userInfo *baseinfo.UserInfo, userNameList []string) (*baseinfo.PackHeader, error) {
	var request wechat.BatchGetContactBriefInfoReq
	request.ContactUsernameList = userNameList
	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeBatchGetContactBriefInfo, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/batchgetcontactbriefinfo", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

func SendGetFriendRelationReq(userInfo *baseinfo.UserInfo, userName string) (*baseinfo.PackHeader, error) {
	var request wechat.MMBizJsApiGetUserOpenIdRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(1)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	request.AppId = proto.String("wx7c8d593b2c3a7703")
	request.UserName = proto.String(userName)
	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	//获取好友关系状态
	sendEncodeData := Pack(userInfo, srcData, 1177, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/usrmsg/mmbizjsapi_getuseropenid", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendGetContactRequest 获取指定微信号信息请求, userWxID:联系人ID  roomWxID：群ID
func SendGetContactRequest(userInfo *baseinfo.UserInfo, userWxIDList []string, antisPanTicketList []string, roomWxIDList []string) (*baseinfo.PackHeader, error) {
	var request wechat.GetContactRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// userCount
	var userCount = uint32(len(userWxIDList))
	request.UserCount = &userCount
	// UserNameList
	userNameList := make([]*wechat.SKBuiltinString, userCount)
	// 遍历
	for index := uint32(0); index < userCount; index++ {
		userNameItem := new(wechat.SKBuiltinString)
		userNameItem.Str = &userWxIDList[index]
		userNameList[index] = userNameItem
	}
	request.UserNameList = userNameList

	// AntispamTicketCount
	antispamTicketCount := uint32(len(antisPanTicketList))
	request.AntispamTicketCount = &antispamTicketCount
	// AntispamTicket
	tmpAntispamTicketList := make([]*wechat.SKBuiltinString, antispamTicketCount)
	for index := uint32(0); index < antispamTicketCount; index++ {
		antispamTicket := new(wechat.SKBuiltinString)
		antispamTicket.Str = &antisPanTicketList[index]
		tmpAntispamTicketList[index] = antispamTicket
	}
	request.AntispamTicket = tmpAntispamTicketList

	// FromChatRoomCount
	fromChatRoomCount := uint32(len(roomWxIDList))
	request.FromChatRoomCount = &fromChatRoomCount
	// FromChatRoom
	fromChatRoomList := make([]*wechat.SKBuiltinString, fromChatRoomCount)
	for index := uint32(0); index < fromChatRoomCount; index++ {
		fromChatRoom := new(wechat.SKBuiltinString)
		fromChatRoom.Str = &roomWxIDList[index]
		fromChatRoomList[index] = fromChatRoom
	}
	request.FromChatRoom = fromChatRoomList

	// GetContactScene
	var getContactScene = uint32(0)
	request.GetContactScene = &getContactScene

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetContact, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 创建红包
func SendWXCreateRedPacket(userInfo *baseinfo.UserInfo, hbItem *baseinfo.RedPacket) (*baseinfo.PackHeader, error) {
	var request wechat.HongBaoReq
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// CgiCmd
	request.CgiCmd = proto.Uint32(0)
	// OutPutType
	request.OutPutType = proto.Uint32(0)
	// ReqText
	strReqText := string("")
	strReqText = strReqText + "city=Guangzhou&"
	strReqText = strReqText + "hbType=" + strconv.Itoa(int(hbItem.RedType)) + "&"
	strReqText = strReqText + "headImg=" + "&"
	strReqText = strReqText + "inWay=" + strconv.Itoa(int(hbItem.From)) + "&"
	strReqText = strReqText + "needSendToMySelf=0" + "&"
	strReqText = strReqText + "nickName=" + url.QueryEscape(userInfo.NickName) + "&"
	strReqText = strReqText + "perValue=" + strconv.Itoa(int(hbItem.Amount)) + "&"
	strReqText = strReqText + "province=Guangdong" + "&"
	strReqText = strReqText + "sendUserName=" + userInfo.WxId + "&"
	strReqText = strReqText + "totalAmount=" + strconv.Itoa(int(hbItem.Amount*hbItem.Count)) + "&"
	strReqText = strReqText + "totalNum=" + strconv.Itoa(int(hbItem.Count)) + "&"
	strReqText = strReqText + "username=" + hbItem.Username + "&"
	strReqText = strReqText + "wishing=" + url.QueryEscape(hbItem.Content)
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText
	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, 1575, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmpay-bin/requestwxhb", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendReceiveWxHB 发送接收红包请求
func SendReceiveWxHB(userInfo *baseinfo.UserInfo, hongBaoReceiverItem *baseinfo.HongBaoReceiverItem) (*baseinfo.PackHeader, error) {
	var request wechat.HongBaoReq

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &hongBaoReceiverItem.CgiCmd

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	strReqText := string("")
	strReqText = strReqText + "agreeDuty=0&"
	strReqText = strReqText + "channelId=" + hongBaoReceiverItem.HongBaoURLItem.ChannelID + "&"
	strReqText = strReqText + "city=" + hongBaoReceiverItem.City + "&"
	strReqText = strReqText + "encrypt_key=" + baseutils.EscapeURL(userInfo.HBAesKeyEncrypted) + "&"
	strReqText = strReqText + "encrypt_userinfo=" + baseutils.EscapeURL(GetEncryptUserInfo(userInfo)) + "&"
	strReqText = strReqText + "inWay=" + strconv.Itoa(int(hongBaoReceiverItem.InWay)) + "&"
	strReqText = strReqText + "msgType=" + hongBaoReceiverItem.HongBaoURLItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoReceiverItem.NativeURL) + "&"
	strReqText = strReqText + "province=" + hongBaoReceiverItem.Province + "&"
	strReqText = strReqText + "sendId=" + hongBaoReceiverItem.HongBaoURLItem.SendID
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeReceiveWxHB, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmpay-bin/receivewxhb", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendOpenWxHB 发送领取红包请求
func SendOpenWxHB(userInfo *baseinfo.UserInfo, hongBaoOpenItem *baseinfo.HongBaoOpenItem) (*baseinfo.PackHeader, error) {
	var request wechat.HongBaoReq

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &hongBaoOpenItem.CgiCmd

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	strReqText := string("")
	strReqText = strReqText + "channelId=" + hongBaoOpenItem.HongBaoURLItem.ChannelID + "&"
	strReqText = strReqText + "city=" + hongBaoOpenItem.City + "&"
	strReqText = strReqText + "encrypt_key=" + baseutils.EscapeURL(userInfo.HBAesKeyEncrypted) + "&"
	strReqText = strReqText + "encrypt_userinfo=" + baseutils.EscapeURL(GetEncryptUserInfo(userInfo)) + "&"
	strReqText = strReqText + "headImg=" + baseutils.EscapeURL(hongBaoOpenItem.HeadImg) + "&"
	strReqText = strReqText + "msgType=" + hongBaoOpenItem.HongBaoURLItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoOpenItem.NativeURL) + "&"
	strReqText = strReqText + "nickName=" + baseutils.HongBaoStringToBytes(hongBaoOpenItem.NickName) + "&"
	strReqText = strReqText + "province=" + hongBaoOpenItem.Province + "&"
	strReqText = strReqText + "sendId=" + hongBaoOpenItem.HongBaoURLItem.SendID + "&"
	strReqText = strReqText + "sessionUserName=" + hongBaoOpenItem.HongBaoURLItem.SendUserName + "&"
	strReqText = strReqText + "timingIdentifier=" + hongBaoOpenItem.TimingIdentifier
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeOpenWxHB, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmpay-bin/openwxhb", sendEncodeData)
	if err != nil {
		return nil, err
	}

	return DecodePackHeader(resp, nil)
}

// SendOpenWxHB 发送查看红包请求
func SendRedEnvelopeWxHB(userInfo *baseinfo.UserInfo, hongBaoOpenItem *baseinfo.HongBaoOpenItem) (*baseinfo.PackHeader, error) {
	var request wechat.HongBaoReq

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &hongBaoOpenItem.CgiCmd

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	strReqText := string("")
	strReqText = strReqText + "agreeDuty=1" + "&"
	strReqText = strReqText + "inWay=1" + "&"
	strReqText = strReqText + "channelId=" + hongBaoOpenItem.HongBaoURLItem.ChannelID + "&"
	strReqText = strReqText + "msgType=" + hongBaoOpenItem.HongBaoURLItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoOpenItem.NativeURL) + "&"
	strReqText = strReqText + "sendId=" + hongBaoOpenItem.HongBaoURLItem.SendID + "&"
	strReqText = strReqText + "sessionUserName=" + hongBaoOpenItem.HongBaoURLItem.SendUserName + "&"
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText
	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeOpenWxHB, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmpay-bin/receivewxhb", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 查看红包领取列表
func SendGetRedPacketListRequest(userInfo *baseinfo.UserInfo, hongBaoOpenItem *baseinfo.GetRedPacketList) (*baseinfo.PackHeader, error) {
	var request wechat.HongBaoReq
	if hongBaoOpenItem.Limit == 0 {
		hongBaoOpenItem.Limit = 11
	}
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// CgiCmd
	request.CgiCmd = proto.Uint32(5)
	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType
	// ReqText
	strReqText := string("")
	strReqText = strReqText + "channelId=" + hongBaoOpenItem.HongBaoItem.ChannelID + "&"
	strReqText = strReqText + "msgType=" + hongBaoOpenItem.HongBaoItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoOpenItem.NativeURL) + "&province=&"
	strReqText = strReqText + "sendId=" + hongBaoOpenItem.HongBaoItem.SendID + "&"
	strReqText = strReqText + "limit=" + strconv.FormatInt(hongBaoOpenItem.Limit, 10) + "&"
	strReqText = strReqText + "offset=" + strconv.FormatInt(hongBaoOpenItem.Offset, 10)
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText
	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeQryDetailWxHB, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmpay-bin/qrydetailwxhb", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendTextMsg 发送文本消息请求 toWxid：接受人微信id，content：消息内容，atWxIDList：@用户微信id列表（toWxid只能是群的wxid，content应为：@用户昵称 @用户昵称 消息内容）
func SendTextMsg(userInfo *baseinfo.UserInfo, toWxid string, content string, atWxIDList []string, ContentType int) (*baseinfo.PackHeader, error) {
	// 构造请求
	var request wechat.NewSendMsgRequest
	var count uint32 = 1
	request.MsgCount = &count
	var msgRequestNewList []*wechat.MicroMsgRequestNew = make([]*wechat.MicroMsgRequestNew, count)
	var msgRequestNew wechat.MicroMsgRequestNew
	var recvierString wechat.SKBuiltinString
	recvierString.Str = &toWxid
	msgRequestNew.ToUserName = &recvierString // 设置接收人wxid
	msgRequestNew.Content = &content          // 设置发送内容
	if ContentType == 0 {
		ContentType = 1
	}
	var tmpType uint32 = uint32(ContentType)
	msgRequestNew.Type = &tmpType // 发送的类型
	currentTime := time.Now()
	misSecond := currentTime.UnixNano() / 1000000
	var seconds = uint32(misSecond / 1000)
	msgRequestNew.CreateTime = &seconds // 设置时间 秒为单位
	seqID := time.Now().UnixNano() / int64(time.Millisecond)
	var tmpCheckCode = WithSeqidCalcCheckCode(toWxid, seqID)
	msgRequestNew.ClientMsgId = &tmpCheckCode // 设置校验码

	// atUserList
	var atUserStr = string("")
	size := len(atWxIDList)
	if size > 0 {
		atUserStr = atUserStr + "<msgsource><atuserlist>"
		for index := int(0); index < size; index++ {
			atUserStr = atUserStr + atWxIDList[index]
			if index < size-1 {
				atUserStr = atUserStr + ","
			}
		}
		atUserStr = atUserStr + "</atuserlist></msgsource>"
		log.Println(atUserStr)
		msgRequestNew.MsgSource = &atUserStr
	}
	msgRequestNewList[0] = &msgRequestNew
	request.ChatSendList = msgRequestNewList

	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeNewSendMsg, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newsendmsg", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 发送图片v1.1
func SendUploadImageNewRequest(userInfo *baseinfo.UserInfo, imgData []byte, toUserName string) (*baseinfo.PackHeader, error) {
	// 构造请求
	var protobufdata []byte
	imgStream := bytes.NewBuffer(imgData)
	Startpos := 0
	datalen := 50000
	datatotalength := imgStream.Len()
	ClientImgId := fmt.Sprintf("%v_%v", userInfo.WxId, time.Now().Unix())
	I := 0
	for {
		Startpos = I * datalen
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}
		Databuff := make([]byte, count)
		_, _ = imgStream.Read(Databuff)
		request := &wechat.UploadMsgImgRequest{
			BaseRequest: GetBaseRequest(userInfo),
			ClientImgId: &wechat.SKBuiltinString{
				Str: proto.String(ClientImgId),
			},
			SenderWxid: &wechat.SKBuiltinString{
				Str: proto.String(userInfo.WxId),
			},
			RecvWxid: &wechat.SKBuiltinString{
				Str: proto.String(toUserName),
			},
			TotalLen: proto.Uint32(uint32(datatotalength)),
			StartPos: proto.Uint32(uint32(Startpos)),
			DataLen:  proto.Uint32(uint32(len(Databuff))),
			Data: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			MsgType:    proto.Uint32(3),
			EncryVer:   proto.Uint32(0),
			ReqTime:    proto.Uint32(uint32(time.Now().Unix())),
			MessageExt: proto.String("png"),
		}
		//序列化
		srcData, _ := proto.Marshal(request)
		sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeForwardCdnImage, 5)
		// 发送请求
		rsp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadmsgimg", sendData)
		if err != nil {
			break
		}
		protobufdata = rsp
		I++
	}
	return DecodePackHeader(protobufdata, nil)
}

// 发送企业oplog
func SendQWOpLogRequest(userInfo *baseinfo.UserInfo, cmdId int64, value []byte) (*baseinfo.PackHeader, error) {
	var request wechat.QYOpLogRequest
	request.Type = proto.Int64(cmdId)
	request.V = value
	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, 806, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/openimoplog", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendOplogRequest 发送修改帐号信息请求
func SendOplogRequest(userInfo *baseinfo.UserInfo, modifyItems []*baseinfo.ModifyItem) (*baseinfo.PackHeader, error) {
	var request wechat.OplogRequest
	// CmdList
	var oplog wechat.CmdList
	count := uint32(len(modifyItems))
	oplog.Count = &count

	// ItemList
	cmdItemList := make([]*wechat.CmdItem, count)
	var index = uint32(0)
	for ; index < count; index++ {
		//Item
		cmdItem := &wechat.CmdItem{}
		cmdItem.CmdId = &modifyItems[index].CmdID

		cmdBuf := &wechat.DATA{}
		cmdBuf.Len = &modifyItems[index].Len
		cmdBuf.Data = modifyItems[index].Data
		cmdItem.CmdBuf = cmdBuf
		cmdItemList[index] = cmdItem
	}
	oplog.ItemList = cmdItemList
	request.Oplog = &oplog

	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeOplog, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/oplog", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendGetQRCodeRequest 获取二维码
func SendGetQRCodeRequest(userInfo *baseinfo.UserInfo, userName string) (*baseinfo.PackHeader, error) {
	var request wechat.GetQRCodeRequest

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// opCode
	opcode := uint32(0)
	request.Opcode = &opcode

	// style
	style := uint32(0)
	request.Style = &style

	// UserName
	var userNameSKBuffer wechat.SKBuiltinString
	userNameSKBuffer.Str = &userName
	request.UserName = &userNameSKBuffer

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetQrCode, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getqrcode", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, []byte(userName))
}

// SendLogOutRequest 发送登出请求
func SendLogOutRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var request wechat.LogOutRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// 打包数据
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, 282, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/logout", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSnsPostRequest 发送朋友圈
func SendSnsPostRequest(userInfo *baseinfo.UserInfo, postItem *baseinfo.SnsPostItem) (*baseinfo.PackHeader, error) {
	var request wechat.SnsPostRequest
	zeroValue32 := uint32(0)
	zeroValue64 := uint64(0)

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// ObjectDesc
	objectDescData := CreateSnsPostItemXML(userInfo.WxId, postItem)
	if postItem.Xml {
		objectDescData = []byte(postItem.Content)
	}

	length := uint32(len(objectDescData))
	var objectDesc wechat.SKBuiltinString_
	objectDesc.Len = &length
	objectDesc.Buffer = objectDescData
	request.ObjectDesc = &objectDesc

	// WithUserListCount
	withUserListCount := uint32(len(postItem.WithUserList))
	request.WithUserListCount = &withUserListCount
	// WithUserList
	request.WithUserList = make([]*wechat.SKBuiltinString, withUserListCount)
	index := uint32(0)
	for ; index < withUserListCount; index++ {
		withUser := &wechat.SKBuiltinString{}
		withUser.Str = &postItem.WithUserList[index]
		request.WithUserList[index] = withUser
	}

	// BlackListCount
	blackListCount := uint32(len(postItem.BlackList))
	request.BlackListCount = &blackListCount
	// BlackList
	request.BlackList = make([]*wechat.SKBuiltinString, blackListCount)
	index = uint32(0)
	for ; index < blackListCount; index++ {
		blackUser := &wechat.SKBuiltinString{}
		blackUser.Str = &postItem.BlackList[index]
		request.BlackList[index] = blackUser
	}

	// GroupUserCount
	groupUserCount := uint32(len(postItem.GroupUserList))
	request.GroupUserCount = &groupUserCount
	// GroupUser
	request.GroupUser = make([]*wechat.SKBuiltinString, groupUserCount)
	index = uint32(0)
	for ; index < groupUserCount; index++ {
		groupUser := &wechat.SKBuiltinString{}
		groupUser.Str = &postItem.GroupUserList[index]
		request.GroupUser[index] = groupUser
	}

	// otherFields
	bgImageType := uint32(1)
	request.PostBgimgType = &bgImageType
	request.ObjectSource = &zeroValue32
	request.ReferId = &zeroValue64
	request.Privacy = &postItem.Privacy
	request.SyncFlag = &zeroValue32

	// ClientId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	clientID := string("sns_post_")
	clientID = clientID + userInfo.WxId + "_" + tmpTimeStr + "_0"
	request.ClientId = &clientID

	// groupCount
	request.GroupCount = &zeroValue32
	request.GroupIds = make([]*wechat.SnsGroup, zeroValue32)

	// mediaInfoCount MediaInfo
	mediaInfoCount := uint32(len(postItem.MediaList))
	request.MediaInfoCount = &mediaInfoCount
	request.MediaInfo = make([]*wechat.MediaInfo, mediaInfoCount)
	for index := uint32(0); index < mediaInfoCount; index++ {
		mediaInfo := &wechat.MediaInfo{}
		source := uint32(2)
		mediaInfo.Source = &source

		// MediaType
		mediaType := uint32(1)
		if postItem.MediaList[index].Type == baseinfo.MMSNSMediaTypeImage {
			mediaType = 1
		}
		mediaInfo.MediaType = &mediaType

		// VideoPlayLength
		mediaInfo.VideoPlayLength = &zeroValue32

		// SessionId
		currentTime := int(time.Now().UnixNano() / 1000000)
		sessionID := "memonts-" + strconv.Itoa(currentTime)
		mediaInfo.SessionId = &sessionID

		// startTime
		startTime := uint32(time.Now().UnixNano() / 1000000000)
		mediaInfo.StartTime = &startTime

		request.MediaInfo[index] = mediaInfo
	}

	// SnsPostOperationFields
	var postOperationFields wechat.SnsPostOperationFields
	postOperationFields.ContactTagCount = &zeroValue32
	postOperationFields.TempUserCount = &zeroValue32
	request.SnsPostOperationFields = &postOperationFields

	// clientcheckdata
	var extSpamInfo wechat.SKBuiltinString_
	if userInfo.DeviceInfo != nil {
		extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	} else {
		extSpamInfo.Buffer = GetExtPBSpamInfoDataA16(userInfo)
	}
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	request.ExtSpamInfo = &extSpamInfo

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsPost, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnspost", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSnsPostRequestByXML 通过XML的信息来发送朋友圈
func SendSnsPostRequestByXML(userInfo *baseinfo.UserInfo, timeLineObj *baseinfo.TimelineObject, blackList []string) (*baseinfo.PackHeader, error) {
	var request wechat.SnsPostRequest
	zeroValue32 := uint32(0)
	zeroValue64 := uint64(0)

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// WithUserListCount
	withUserListCount := uint32(0)
	request.WithUserListCount = &withUserListCount
	// WithUserList
	request.WithUserList = make([]*wechat.SKBuiltinString, withUserListCount)

	// BlackListCount
	tmpCount := uint32(len(blackList))
	request.BlackListCount = &tmpCount
	request.BlackList = make([]*wechat.SKBuiltinString, tmpCount)
	// BlackList
	for index := uint32(0); index < tmpCount; index++ {
		tmpSKBuiltinString := &wechat.SKBuiltinString{}
		tmpSKBuiltinString.Str = &blackList[index]
		request.BlackList[index] = tmpSKBuiltinString
	}

	// GroupUserCount
	groupUserCount := uint32(0)
	request.GroupUserCount = &groupUserCount
	// GroupUser
	request.GroupUser = make([]*wechat.SKBuiltinString, groupUserCount)

	// otherFields
	bgImageType := uint32(1)
	request.PostBgimgType = &bgImageType
	request.ObjectSource = &zeroValue32
	request.ReferId = &zeroValue64
	request.Privacy = &timeLineObj.Private
	request.SyncFlag = &zeroValue32

	// ClientId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	clientID := string("sns_post_")
	clientID = clientID + userInfo.WxId + "_" + tmpTimeStr + "_0"
	request.ClientId = &clientID

	// groupCount
	request.GroupCount = &zeroValue32
	request.GroupIds = make([]*wechat.SnsGroup, zeroValue32)

	// mediaInfoCount MediaInfo
	mediaInfoCount := uint32(len(timeLineObj.ContentObject.MediaList.Media))
	request.MediaInfoCount = &mediaInfoCount
	request.MediaInfo = make([]*wechat.MediaInfo, mediaInfoCount)
	for index := uint32(0); index < mediaInfoCount; index++ {
		tmpMediaItem := timeLineObj.ContentObject.MediaList.Media[index]

		// 解析Source
		mediaInfo := &wechat.MediaInfo{}
		tmpSource := baseutils.ParseInt(tmpMediaItem.URL.Type)
		mediaInfo.Source = &tmpSource
		// MediaType
		mediaType := tmpMediaItem.Type - 1
		mediaInfo.MediaType = &mediaType
		// VideoPlayLength
		playLength := uint32(tmpMediaItem.VideoDuration)
		mediaInfo.VideoPlayLength = &playLength
		// SessionId
		currentTime := int(time.Now().UnixNano() / 1000000)
		sessionID := "memonts-" + strconv.Itoa(currentTime)
		mediaInfo.SessionId = &sessionID

		// startTime
		startTime := uint32(time.Now().UnixNano() / 1000000000)
		mediaInfo.StartTime = &startTime
		request.MediaInfo[index] = mediaInfo
	}

	// ID和UserName置为0
	timeLineObj.UserName = userInfo.WxId
	timeLineObj.CreateTime = uint32(int(time.Now().UnixNano() / 1000000000))
	// ObjectDesc
	objectDescData, err := xml.Marshal(timeLineObj)
	if err != nil {
		return nil, err
	}
	str := string(objectDescData)
	str = strings.ReplaceAll(str, "token=\"\"", "")
	str = strings.ReplaceAll(str, "key=\"\"", "")
	str = strings.ReplaceAll(str, "enc_idx=\"\"", "")
	str = strings.ReplaceAll(str, "md5=\"\"", "")
	str = strings.ReplaceAll(str, "videomd5=\"\"", "")
	str = strings.ReplaceAll(str, "video", "")
	objectDescData = []byte(str)
	length := uint32(len(objectDescData))
	var objectDesc wechat.SKBuiltinString_
	objectDesc.Len = &length
	objectDesc.Buffer = objectDescData
	request.ObjectDesc = &objectDesc

	// SnsPostOperationFields
	var postOperationFields wechat.SnsPostOperationFields
	postOperationFields.ContactTagCount = &zeroValue32
	postOperationFields.TempUserCount = &zeroValue32
	request.SnsPostOperationFields = &postOperationFields

	// clientcheckdata
	var extSpamInfo wechat.SKBuiltinString_
	if userInfo.DeviceInfo != nil {
		extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	} else {
		extSpamInfo.Buffer = GetExtPBSpamInfoDataA16(userInfo)
	}
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	request.ExtSpamInfo = &extSpamInfo

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsPost, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnspost", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSnsObjectOpRequest 发送朋友圈操作
func SendSnsObjectOpRequest(userInfo *baseinfo.UserInfo, opItems []*baseinfo.SnsObjectOpItem) (*baseinfo.PackHeader, error) {
	var request wechat.SnsObjectOpRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// OpCount
	opCount := uint32(len(opItems))
	request.OpCount = &opCount

	// OpList
	request.OpList = make([]*wechat.SnsObjectOp, opCount)
	index := uint32(0)
	for ; index < opCount; index++ {
		snsObject := &wechat.SnsObjectOp{}
		id, _ := strconv.ParseUint(opItems[index].SnsObjID, 0, 64)
		snsObject.Id = &id
		snsObject.OpType = &opItems[index].OpType
		if opItems[index].DataLen > 0 {
			skBuffer := &wechat.SKBuiltinString_{}
			skBuffer.Len = &opItems[index].DataLen
			skBuffer.Buffer = opItems[index].Data
		}
		if opItems[index].Ext != 0 {
			extInfo := &wechat.SnsObjectOpExt{
				Id: &opItems[index].Ext,
			}
			CommnetId, _ := proto.Marshal(extInfo)
			snsObject.Ext = &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(CommnetId))),
				Buffer: CommnetId,
			}
		}
		request.OpList[index] = snsObject
	}

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsObjectOp, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnsobjectop", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

func Uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}

// SendSnsUserPageRequest 发送 获取朋友圈信息 请求
func SendSnsUserPageRequest(userInfo *baseinfo.UserInfo, userName string, firstPageMd5 string, maxID uint64) (*baseinfo.PackHeader, error) {
	var request wechat.SnsUserPageRequest
	var zeroValue64 = uint64(0)
	var zeroValue32 = uint32(0)

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// 其它参数
	request.Username = &userName
	request.FirstPageMd5 = &firstPageMd5
	request.MaxId = &maxID
	request.MinFilterId = &zeroValue64
	request.LastRequestTime = &zeroValue32
	request.FilterType = &zeroValue32

	// 打包数据
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsUserPage, 5)

	// 发送请求
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnsuserpage", sendData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, srcData)
}

// SendSnsCommentRequest 发送评论/点赞请求
func SendSnsCommentRequest(userInfo *baseinfo.UserInfo, commentItem *baseinfo.SnsCommentItem) (*baseinfo.PackHeader, error) {
	var request wechat.SnsCommentRequest
	zeroValue32 := uint32(0)
	zeroValue64 := int64(0)
	emptyString := string("")

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(1)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	request.BaseRequest.OsType = proto.String("wechat")

	// ClientId
	clientID := string("wcc:")
	clientID = clientID + userInfo.WxId + "-" + strconv.Itoa(int(commentItem.CreateTime)) + "-0"
	request.ClientId = &clientID

	// Action
	request.Action = &wechat.SnsActionGroup{}
	request.Action.Id = &commentItem.ItemID
	request.Action.ParentId = proto.Uint64(0)
	request.Action.ClientId = &emptyString
	request.Action.ObjectCreateTime = &zeroValue32

	// Action.CurrentAction
	request.Action.CurrentAction = &wechat.SnsAction{}
	request.Action.CurrentAction.FromUsername = &userInfo.WxId
	request.Action.CurrentAction.FromNickname = &userInfo.NickName
	request.Action.CurrentAction.ToUsername = &commentItem.ToUserName
	request.Action.CurrentAction.ToNickname = &commentItem.ToUserName
	request.Action.CurrentAction.Type = &commentItem.OpType
	request.Action.CurrentAction.Source = &zeroValue32
	request.Action.CurrentAction.ReplyCommentId = &commentItem.ReplyCommentID
	request.Action.CurrentAction.CreateTime = &commentItem.CreateTime
	if commentItem.OpType == baseinfo.MMSnsCommentTypeComment {
		request.Action.CurrentAction.Content = &commentItem.Content
	}
	request.Action.CurrentAction.CommentId = &zeroValue32
	request.Action.CurrentAction.ReplyCommentId2 = &zeroValue64
	request.Action.CurrentAction.CommentId2 = &zeroValue64
	request.Action.CurrentAction.IsNotRichText = &zeroValue32
	request.Action.CurrentAction.CommentFlag = &zeroValue32

	// Action.ReferAction
	if commentItem.ReplyItem != nil {
		request.Action.ReferAction = &wechat.SnsAction{}
		request.Action.ReferAction.FromUsername = &commentItem.ReplyItem.UserName
		request.Action.ReferAction.FromNickname = &commentItem.ReplyItem.NickName
		request.Action.ReferAction.ToUsername = &emptyString
		request.Action.ReferAction.ToNickname = &emptyString
		request.Action.ReferAction.Type = &commentItem.ReplyItem.OpType
		request.Action.ReferAction.Content = &emptyString
		request.Action.ReferAction.Source = &commentItem.ReplyItem.Source
		request.Action.ReferAction.IsNotRichText = &zeroValue32
		request.Action.ReferAction.CreateTime = &zeroValue32
		request.Action.ReferAction.ReplyCommentId = &zeroValue32
		request.Action.ReferAction.CommentId = &zeroValue32
		request.Action.ReferAction.ReplyCommentId2 = &zeroValue64
		request.Action.ReferAction.CommentId2 = &zeroValue64
		request.Action.ReferAction.CommentFlag = &zeroValue32
	}

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsComment, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnscomment", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取收藏lit
func SendFavSyncListRequest(userInfo *baseinfo.UserInfo, keyBuf string) (*baseinfo.PackHeader, error) {
	var request wechat.FavSyncRequest
	selector := uint32(1)
	request.Selector = &selector
	var skBufferT wechat.SKBuiltinString_
	skBufferT.Len = proto.Uint32(0)
	/*favSyncKeyLen := uint32(len(userInfo.FavSyncKey))
	skBufferT.Len = &favSyncKeyLen
	skBufferT.Buffer = userInfo.FavSyncKey*/
	if keyBuf != "" {
		key, _ := base64.StdEncoding.DecodeString(keyBuf)
		skBufferT.Buffer = key
		skBufferT.Len = proto.Uint32(uint32(len(key)))
	}
	request.KeyBuf = &skBufferT
	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeFavSync, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/favsync", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendFavSyncRequest 同步收藏
func SendFavSyncRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var request wechat.FavSyncRequest

	// selector
	selector := uint32(1)
	request.Selector = &selector

	var skBufferT wechat.SKBuiltinString_
	favSyncKeyLen := uint32(len(userInfo.FavSyncKey))
	skBufferT.Len = &favSyncKeyLen
	skBufferT.Buffer = userInfo.FavSyncKey
	request.KeyBuf = &skBufferT

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeFavSync, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/favsync", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendGetFavInfoRequest 获取 收藏信息
func SendGetFavInfoRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var request wechat.GetFavInfoRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetFavInfo, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getfavinfo", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendBatchGetFavItemRequest 获取单条收藏
func SendBatchGetFavItemRequest(userInfo *baseinfo.UserInfo, favID uint32) (*baseinfo.PackHeader, error) {
	var request wechat.BatchGetFavItemRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// Count
	favIDCount := uint32(1)
	request.Count = &favIDCount

	// FavIdList
	request.FavIdList = make([]byte, 0)
	tmpBytes := baseutils.EncodeVByte32(favID)
	request.FavIdList = append(request.FavIdList, tmpBytes[0:]...)

	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeBatchGetFavItem, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/batchgetfavitem", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendBatchDelFavItemRequest 删除收藏项
func SendBatchDelFavItemRequest(userInfo *baseinfo.UserInfo, favID uint32) (*baseinfo.PackHeader, error) {
	var request wechat.BatchDelFavItemRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// Count
	tmpCount := uint32(1)
	request.Count = &tmpCount

	// FavIdList
	request.FavIdList = make([]byte, 0)
	tmpBytes := baseutils.EncodeVByte32(favID)
	request.FavIdList = append(request.FavIdList, tmpBytes[0:]...)

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeBatchDelFavItem, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/batchdelfavitem", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 上报
func SendReportstrategyRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	request := &wechat.GetReportStrategyReq{
		BaseRequest: GetBaseRequest(userInfo),
		DeviceBrand: proto.String("iPad Mini 2G (WiFi)<iPad4,4>"),
		OsName:      proto.String("Apple"),
		DeviceModel: proto.String("IOS"),
		OsVersion:   proto.String(userInfo.DeviceInfo.OsTypeNumber),
		LanguageVer: proto.String(userInfo.DeviceInfo.Language),
	}
	// 打包数据 发送
	srcData, _ := proto.Marshal(request)
	sendEncodeData := Pack(userInfo, srcData, 308, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/reportstrategy", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendGetCDNDnsRequest 获取该帐号的CdnDns信息
func SendGetCDNDnsRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var request wechat.GetCDNDnsRequest
	emptyString := string("")

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// ClientIp
	request.ClientIp = &emptyString

	// Scene
	scene := uint32(1)
	request.Scene = &scene

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetCdnDNS, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getcdndns", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSnsObjectDetailRequest SendSnsObjectDetailRequest
func SendSnsObjectDetailRequest(userInfo *baseinfo.UserInfo, snsID uint64) (*baseinfo.PackHeader, error) {
	var request wechat.SnsObjectDetailRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	baseReq.ClientVersion = proto.Uint32(0x16070228)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// ID
	request.Id = &snsID

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsObjectDetail, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnsobjectdetail", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSnsSyncRequest 同步朋友圈
func SendSnsSyncRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var request wechat.SnsSyncRequest
	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// Selector
	tmpSelector := uint32(509)
	request.Selector = &tmpSelector

	// KeyBuf 第一次使用同步消息Key
	tmpKeyBuffer := userInfo.SnsSyncKey
	if len(tmpKeyBuffer) <= 0 {
		tmpKeyBuffer = userInfo.SyncKey
	}
	tmpLen := uint32(len(tmpKeyBuffer))
	var tmpKeyBuf wechat.SKBuiltinString_
	tmpKeyBuf.Buffer = tmpKeyBuffer
	tmpKeyBuf.Len = &tmpLen
	request.KeyBuf = &tmpKeyBuf

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsSync, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnssync", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendGetContactLabelListRequest 获取设置好的联系人标签列表
func SendGetContactLabelListRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var request wechat.GetContactLabelListRequest

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetContactLabelList, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getcontactlabellist", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendAddContactLabelRequest 发送添加标签请求
func SendAddContactLabelRequest(userInfo *baseinfo.UserInfo, newLabelList []string) (*baseinfo.PackHeader, error) {
	var request wechat.AddContactLabelRequest
	labelID := uint32(0)

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// LabelCount
	labelCount := uint32(len(newLabelList))
	request.LabelCount = &labelCount

	// LabelPairList
	request.LabelPairList = make([]*wechat.LabelPair, labelCount)
	for index := uint32(0); index < labelCount; index++ {
		labelPair := &wechat.LabelPair{}
		labelPair.LabelName = &newLabelList[index]
		labelPair.LabelId = &labelID
		request.LabelPairList[index] = labelPair
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeAddContactLabel, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/addcontactlabel", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendDelContactLabelRequest 删除标签
func SendDelContactLabelRequest(userInfo *baseinfo.UserInfo, labelId string) (*baseinfo.PackHeader, error) {
	req := wechat.DelContactLabelRequest{
		BaseRequest: GetBaseRequest(userInfo),
		LabelIdlist: proto.String(labelId),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeDelContactLabel, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delcontactlabel", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendModifyLabelRequest 修改标签请求
func SendModifyLabelRequest(userInfo *baseinfo.UserInfo, userLabelList []baseinfo.UserLabelInfoItem) (*baseinfo.PackHeader, error) {
	_userLabelList := make([]*wechat.UserLabelInfo, 0)
	for _, item := range userLabelList {
		_userLabelList = append(_userLabelList, &wechat.UserLabelInfo{
			UserName:    proto.String(item.UserName),
			LabelIdlist: proto.String(item.LabelIDList),
		})
	}
	req := wechat.ModifyContactLabelListRequest{
		BaseRequest:       GetBaseRequest(userInfo),
		UserCount:         proto.Uint32(uint32(len(_userLabelList))),
		UserLabelInfoList: _userLabelList,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeModifyContactLabelList, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/modifycontactlabellist", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendBindQueryNewRequest SendBindQueryNewRequest
func SendBindQueryNewRequest(userInfo *baseinfo.UserInfo, reqItem *baseinfo.TenPayReqItem) (*baseinfo.PackHeader, error) {
	var request wechat.TenPayRequest
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &reqItem.CgiCMD

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	var reqTextSKBuf wechat.SKBuiltinString_
	tmpLen := uint32(len(reqItem.ReqText))
	reqTextSKBuf.Len = &tmpLen
	reqTextSKBuf.Buffer = []byte(reqItem.ReqText)
	request.ReqText = &reqTextSKBuf

	// ReqTextWx
	var wxReqTextSKBuf wechat.SKBuiltinString_
	tmpText := "encrypt_key=" + userInfo.HBAesKeyEncrypted
	tmpText = tmpText + "&encrypt_userinfo=" + GetEncryptUserInfo(userInfo)
	tmpWXLen := uint32(len(tmpText))
	wxReqTextSKBuf.Len = &tmpWXLen
	wxReqTextSKBuf.Buffer = []byte(tmpText)
	request.ReqTextWx = &wxReqTextSKBuf

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeBindQueryNew, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmpay-bin/tenpay/bindquerynew", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 确定收款
func SendTenPayRequest(userInfo *baseinfo.UserInfo, reqItem *baseinfo.TenPayReqItem) (*baseinfo.PackHeader, error) {
	var request wechat.TenPayRequest
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &reqItem.CgiCMD

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	var reqTextSKBuf wechat.SKBuiltinString_
	tmpLen := uint32(len(reqItem.ReqText))
	reqTextSKBuf.Len = &tmpLen
	reqTextSKBuf.Buffer = []byte(reqItem.ReqText)
	request.ReqText = &reqTextSKBuf

	// ReqTextWx
	var wxReqTextSKBuf wechat.SKBuiltinString_
	wxReqTextSKBuf.Buffer = []byte(reqItem.ReqText)
	wxReqTextSKBuf.Len = &tmpLen

	request.ReqTextWx = &wxReqTextSKBuf

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, 385, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/tenpay", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSnsTimeLineRequest 发送获取朋友圈请求
func SendSnsTimeLineRequest(userInfo *baseinfo.UserInfo, firstPageMD5 string, maxID uint64) (*baseinfo.PackHeader, error) {
	/*var request wechat.SnsTimeLineRequest
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// ClientLatestId
	tmpLatestID := uint64(0)
	request.ClientLatestId = &tmpLatestID
	// FirstPageMd5
	request.FirstPageMd5 = &firstPageMD5
	// LastRequestTime
	lastRequestTime := uint32(0)
	request.LastRequestTime = &lastRequestTime
	// MAXID
	request.MaxId = &maxID
	// MinFilterId
	minFilterID := uint64(0)
	request.MinFilterId = &minFilterID
	// NetworkType
	netWorkType := uint32(1)
	request.NetworkType = &netWorkType*/

	req := &wechat.SnsTimeLineRequest{
		BaseRequest:     GetBaseRequest(userInfo),
		ClientLatestId:  proto.Uint64(0),
		FirstPageMd5:    proto.String(firstPageMD5),
		LastRequestTime: proto.Uint32(0),
		MaxId:           proto.Uint64(maxID),
		MinFilterId:     proto.Uint64(0),
		NetworkType:     proto.Uint32(1),
	}
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	req.BaseRequest = baseReq
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsTimeLine, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mmsnstimeline", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendAppMsgRequest 发送App消息
func SendAppMsgRequest(userInfo *baseinfo.UserInfo, contentType uint32, toUserName, xml string) (*baseinfo.PackHeader, error) {
	req := wechat.SendAppMsgRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Msg: &wechat.AppMsg{
			FromUserName: proto.String(userInfo.GetUserName()),
			AppId:        proto.String(""),
			SdkVersion:   proto.Uint32(0),
			ToUserName:   proto.String(toUserName),
			Type:         proto.Uint32(contentType),
			Content:      proto.String(xml),
			CreateTime:   proto.Uint32(uint32(time.Now().Unix())),
			ClientMsgId:  proto.String(fmt.Sprintf("%s_%v", toUserName, time.Now().Unix())),
			Source:       proto.Int32(0),
			RemindId:     proto.Int32(0),
			MsgSource:    proto.String(""),
			Thumb: &wechat.BufferT{
				ILen:   proto.Uint32(0),
				Buffer: []byte{},
			},
		},
		FromSence:     proto.String(""),
		DirectShare:   proto.Int32(0),
		SendMsgTicket: proto.String(""),
	}

	// 打包发送数据
	src, _ := proto.Marshal(&req)

	sendEncodeData := Pack(userInfo, src, baseinfo.MMRequestTypeSendAppMsg, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/sendappmsg", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendEmojiRequest 发生表情
func SendEmojiRequest(userInfo *baseinfo.UserInfo, toUserName, Md5 string, length int32) (*baseinfo.PackHeader, error) {
	baseRequest := GetBaseRequest(userInfo)
	//baseRequest.Scene = proto.Uint32(0)
	req := wechat.UploadEmojiRequest{
		BaseRequest:    baseRequest,
		EmojiItemCount: proto.Int32(1),
		EmojiItem: []*wechat.EmojiUploadInfoReq{
			{
				MD5:      proto.String(Md5),
				StartPos: proto.Int32(length),
				TotalLen: proto.Int32(length),
				EmojiBuffer: &wechat.BufferT{
					ILen:   proto.Uint32(0),
					Buffer: []byte{},
				},
				Type:        proto.Int32(2),
				ToUserName:  proto.String(toUserName),
				ClientMsgID: proto.String(fmt.Sprintf("%d", time.Now().UnixNano()/1000/1000)),
			},
		},
	}

	// 打包发送数据
	src, _ := proto.Marshal(&req)

	sendEncodeData := Pack(userInfo, src, baseinfo.MMRequestTypeSendEmoji, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/sendemoji", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 发送表情 new 包含动图
func ForwardEmojiRequest(userInfo *baseinfo.UserInfo, toUserName, Md5 string, length int32) (*baseinfo.PackHeader, error) {
	baseRequest := GetBaseRequest(userInfo)
	//baseRequest.Scene = proto.Uint32(0)
	req := wechat.UploadEmojiRequest{
		BaseRequest:    baseRequest,
		EmojiItemCount: proto.Int32(1),
		EmojiItem: []*wechat.EmojiUploadInfoReq{
			{
				MD5:      proto.String(Md5),
				StartPos: proto.Int32(0),
				TotalLen: proto.Int32(length),
				EmojiBuffer: &wechat.BufferT{
					ILen:   proto.Uint32(0),
					Buffer: []byte{},
				},
				Type:        proto.Int32(1),
				ToUserName:  proto.String(toUserName),
				ClientMsgID: proto.String(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}

	// 打包发送数据
	src, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, src, baseinfo.MMRequestTypeSendEmoji, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/sendemoji", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 下载语音
func SendGetMsgVoiceRequest(userInfo *baseinfo.UserInfo, toUserName, NewMsgIds, Bufid string, Length int) (*vo.DownloadVoiceData, error) {
	I := 0
	Startpos := 0
	datalen := 50000
	Databuff := make([]byte, Length+1000)
	var VoiceLength uint32
	NewMsgId, _ := strconv.ParseUint(NewMsgIds, 10, 64)
	MasterBufId, _ := strconv.ParseUint(Bufid, 10, 64)
	resp := new(wechat.DownloadVoiceResponse)
	for {
		Startpos = I * datalen
		count := 0
		if Length-Startpos > datalen {
			count = Length
		} else {
			count = Length - Startpos
		}
		if count < 0 {
			break
		}
		req := wechat.DownloadVoiceRequest{
			BaseRequest:  GetBaseRequest(userInfo),
			MsgId:        proto.Uint32(0),
			Offset:       proto.Uint32(uint32(Startpos)),
			Length:       proto.Uint32(uint32(count)),
			NewMsgId:     proto.Uint64(NewMsgId),
			ChatRoomName: proto.String(toUserName),
			MasterBufId:  proto.Uint64(MasterBufId),
		}
		// 打包发送数据
		srcData, _ := proto.Marshal(&req)
		sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeDownloadVoice, 5)
		res, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/downloadvoice", sendEncodeData)
		if err != nil {
			break
		}
		header, err := DecodePackHeader(res, nil)
		if err != nil {
			break
		}
		err = ParseResponseData(userInfo, header, resp)
		if err != nil || resp.GetBaseResponse().GetRet() != 0 {
			break
		}
		DataStream := bytes.NewBuffer(resp.GetData().GetBuffer())
		_, _ = DataStream.Read(Databuff)
		VoiceLength = resp.GetVoiceLength()
		I++
	}
	return &vo.DownloadVoiceData{
		Base64:      Databuff,
		VoiceLength: VoiceLength,
	}, nil
}

// 群发图片
func SendGroupMassMsgImage(userInfo *baseinfo.UserInfo, toUSerName []string, ImageBase64 []byte) (*baseinfo.PackHeader, error) {
	baseRequest := GetBaseRequest(userInfo)
	baseRequest.Scene = proto.Uint32(0)
	toList := strings.Join(toUSerName, ";")
	tolistmd5 := baseutils.MD5ToLower(toList)
	ClientImgId := fmt.Sprintf("%v_%v", time.Now().Unix(), tolistmd5)
	var protobufdata []byte
	imgStream := bytes.NewBuffer(ImageBase64)
	Startpos := 0
	datalen := 50000
	datatotalength := imgStream.Len()
	I := 0
	for {
		Startpos = I * datalen
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}
		Databuff := make([]byte, count)
		_, _ = imgStream.Read(Databuff)
		req := wechat.MassSendRequest{
			BaseRequest: baseRequest,
			ToList:      proto.String(toList),
			ToListMd5:   proto.String(tolistmd5),
			ClientId:    proto.String(ClientImgId),
			MsgType:     proto.Uint64(3),
			MediaTime:   proto.Uint64(0),
			DataBuffer: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			DataStartPos:  proto.Uint64(uint64(Startpos)),
			DataTotalLen:  proto.Uint64(uint64(len(Databuff))),
			ThumbTotalLen: proto.Uint64(0),
			ThumbStartPos: proto.Uint64(0),
			ThumbData: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(0),
				Buffer: []byte{},
			},
			CameraType:   proto.Uint64(2),
			VideoSource:  proto.Uint64(0),
			ToListCount:  proto.Uint64(uint64(len(toUSerName))),
			IsSendAgain:  proto.Uint64(0),
			CompressType: proto.Uint64(1),
			VoiceFormat:  proto.Uint64(0),
		}
		// 打包发送数据
		src, _ := proto.Marshal(&req)
		sendEncodeData := Pack(userInfo, src, 193, 5)
		rsp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/masssend", sendEncodeData)
		if err != nil {
			return nil, err
		}
		protobufdata = rsp
		I++
	}
	return DecodePackHeader(protobufdata, nil)
}

func SendGroupMassMsgText(userInfo *baseinfo.UserInfo, toUSerName []string, content string) (*baseinfo.PackHeader, error) {
	baseRequest := GetBaseRequest(userInfo)
	//baseRequest.Scene = proto.Uint32(0)
	toList := strings.Join(toUSerName, ";")
	tolistmd5 := baseutils.MD5ToLower(toList)
	Databuff := []byte(content)
	ClientImgId := fmt.Sprintf("%v_%v", time.Now().Unix(), tolistmd5)
	req := wechat.MassSendRequest{
		BaseRequest: baseRequest,
		ToList:      proto.String(toList),
		ToListMd5:   proto.String(tolistmd5),
		ClientId:    proto.String(ClientImgId),
		MsgType:     proto.Uint64(1),
		MediaTime:   proto.Uint64(0),
		DataBuffer: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(uint32(len(Databuff))),
			Buffer: Databuff,
		},
		DataStartPos:  proto.Uint64(0),
		DataTotalLen:  proto.Uint64(uint64(len(Databuff))),
		ThumbTotalLen: proto.Uint64(0),
		ThumbStartPos: proto.Uint64(0),
		ThumbData: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(0),
			Buffer: []byte{},
		},
		CameraType:   proto.Uint64(2),
		VideoSource:  proto.Uint64(0),
		ToListCount:  proto.Uint64(uint64(len(toUSerName))),
		IsSendAgain:  proto.Uint64(1),
		CompressType: proto.Uint64(0),
		VoiceFormat:  proto.Uint64(0),
	}
	// 打包发送数据
	src, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, src, 193, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/masssend", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 群拍一拍
func SendSendPatRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName string, scene int64) (*baseinfo.PackHeader, error) {
	ClientImgId := fmt.Sprintf("%v_%v_%v", userInfo.WxId, toUserName, time.Now().Unix())
	req := wechat.SendPatRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		FromUsername:   proto.String(userInfo.WxId),
		ChatUsername:   proto.String(chatRoomName),
		PattedUsername: proto.String(toUserName),
		ClientMsgId:    proto.String(ClientImgId),
		Scene:          proto.Int64(scene),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 849, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/sendpat", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SetChatRoomAnnouncementRequest 设置群公告
func SetChatRoomAnnouncementRequest(userInfo *baseinfo.UserInfo, roomId, content string) (*baseinfo.PackHeader, error) {
	req := wechat.SetChatRoomAnnouncementRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		ChatRoomName: proto.String(roomId),
		Announcement: proto.String(content),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeSetChatRoomAnnouncement, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/setchatroomannouncement", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取群详细
func SetGetChatRoomInfoDetailRequest(userInfo *baseinfo.UserInfo, roomId string) (*baseinfo.PackHeader, error) {
	req := wechat.GetChatRoomInfoDetailRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		ChatRoomName: proto.String(roomId),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetChatRoomInfoDetail, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getchatroominfodetail", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

//保存群聊操作
/*func SetMoveToContractRequest(userInfo *baseinfo.UserInfo, ChatRoomName string, Val uint32) (*baseinfo.PackHeader, error) {
	UserNameListSplit := strings.Split(ChatRoomName, ",")


	req := wechat.GetChatRoomInfoDetailRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		ChatRoomName: proto.String(roomId),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetChatRoomInfoDetail, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getchatroominfodetail", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}*/

// GetChatroomMemberDetailRequest 获取群成员详细
func GetChatroomMemberDetailRequest(userInfo *baseinfo.UserInfo, roomId string) (*baseinfo.PackHeader, error) {
	req := wechat.GetChatroomMemberDetailRequest{
		BaseRequest:   GetBaseRequest(userInfo),
		ChatroomWxid:  proto.String(roomId),
		ClientVersion: proto.Uint32(0),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetChatRoomMemberDetail, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getchatroommemberdetail", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetCreateChatRoomEntity 创建群
func GetCreateChatRoomEntity(userInfo *baseinfo.UserInfo, topIc string, userList []string) (*baseinfo.PackHeader, error) {
	userList = append([]string{userInfo.GetUserName()}, userList...)
	memberList := make([]*wechat.MemberReq, 0)
	for _, user := range userList {
		memberList = append(memberList, &wechat.MemberReq{
			MemberName: &wechat.SKBuiltinString{
				Str: proto.String(user),
			},
		})
	}
	req := wechat.CreateChatRoomRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Topic: &wechat.SKBuiltinString{
			Str: proto.String(topIc),
		},
		MemberCount: proto.Uint32(uint32(len(memberList))),
		MemberList:  memberList,
		Scene:       proto.Uint32(0),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeCreateChatRoom, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/createchatroom", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 添加群管理
func SendAddChatroomAdmin(userInfo *baseinfo.UserInfo, chatRoomName string, userList []string) (*baseinfo.PackHeader, error) {
	req := wechat.AddChatRoomAdminRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		ChatRoomName: proto.String(chatRoomName),
		UserNameList: userList,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 889, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/addchatroomadmin", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 删除群管理
func SendDelChatroomAdminRequest(userInfo *baseinfo.UserInfo, chatRoomName string, userList []string) (*baseinfo.PackHeader, error) {
	req := wechat.DelChatRoomAdminRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		ChatRoomName: proto.String(chatRoomName),
		UserNameList: userList,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 259, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delchatroomadmin", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取群列表
func SendWXSyncContactRequest(userInfo *baseinfo.UserInfo, key []byte) (*baseinfo.PackHeader, error) {
	keyBuf := userInfo.SyncKey
	if key != nil {
		keyBuf = key
	}
	osType := ""
	if userInfo.DeviceInfoA16 != nil {
		osType = baseinfo.AndroidDeviceType
	} else {
		osType = userInfo.DeviceInfo.OsType
	}
	req := wechat.NewSyncRequest{
		Oplog: &wechat.CmdList{
			Count: proto.Uint32(0),
		},
		DeviceType:    proto.String(osType),
		Scene:         proto.Uint32(3), //有时候03，有时候01
		Selector:      proto.Uint32(7), //  7
		SyncMsgDigest: proto.Uint32(baseinfo.MMSyncMsgDigestTypeShortLink),
		KeyBuf: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(uint32(len(keyBuf))),
			Buffer: keyBuf,
		},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 138, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newsync", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 搜手机或企业对外名片链接提取验证
func SendQWSearchContactRequest(userInfo *baseinfo.UserInfo, tg string, fromScene uint64, userName string) (*baseinfo.PackHeader, error) {
	req := wechat.SearchQYContactRequest{}
	if utils.IsMobile(tg) {
		req.Tg = proto.String(tg)
		req.FromScene = proto.Uint64(1)
	} else {
		req.UserName = proto.String(tg)
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 372, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/searchopenimcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetAddChatRoomMemberRequest 拉人
func GetAddChatRoomMemberRequest(userInfo *baseinfo.UserInfo, chatRoomName string, userList []string) (*baseinfo.PackHeader, error) {
	memberList := make([]*wechat.MemberReq, 0)
	for _, user := range userList {
		memberList = append(memberList, &wechat.MemberReq{
			MemberName: &wechat.SKBuiltinString{
				Str: proto.String(user),
			},
		})
	}

	req := wechat.AddChatRoomMemberRequest{
		BaseRequest: GetBaseRequest(userInfo),
		MemberCount: proto.Uint32(uint32(len(memberList))),
		MemberList:  memberList,
		ChatRoomName: &wechat.SKBuiltinString{
			Str: proto.String(chatRoomName),
		},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeAddChatRoomMember, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/addchatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// DelDelChatRoomMember 删除群成员
func DelDelChatRoomMemberRequest(userInfo *baseinfo.UserInfo, chatRoomName string, delUserList []string) (*baseinfo.PackHeader, error) {
	memberList := make([]*wechat.DelMemberReq, 0)
	for _, user := range delUserList {
		memberList = append(memberList, &wechat.DelMemberReq{
			MemberName: &wechat.SKBuiltinString{
				Str: proto.String(user),
			},
		})
	}

	req := wechat.DelChatRoomMemberRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		MemberCount:  proto.Uint32(uint32(len(memberList))),
		MemberList:   memberList,
		ChatRoomName: proto.String(chatRoomName),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeDelChatRoomMember, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delchatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetTransferGroupOwnerRequest 转让群
func GetTransferGroupOwnerRequest(userInfo *baseinfo.UserInfo, chatRoomName, newOwnerUserName string) (*baseinfo.PackHeader, error) {
	req := wechat.TransferChatRoomOwnerRequest{
		BaseRequest:      GetBaseRequest(userInfo),
		ChatRoomName:     proto.String(chatRoomName),
		NewOwnerUserName: proto.String(newOwnerUserName),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeTransferChatRoomOwnerRequest, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/transferchatroomowner", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetInviteChatroomMembersEntity 邀请群成员
func GetInviteChatroomMembersRequest(userInfo *baseinfo.UserInfo, chatRoomName string, userList []string) (*baseinfo.PackHeader, error) {
	memberList := make([]*wechat.MemberReq, 0)
	for _, user := range userList {
		memberList = append(memberList, &wechat.MemberReq{
			MemberName: &wechat.SKBuiltinString{
				Str: proto.String(user),
			},
		})
	}
	req := wechat.InviteChatRoomMemberRequest{
		BaseRequest: GetBaseRequest(userInfo),
		MemberCount: proto.Uint32(uint32(len(memberList))),
		MemberList:  memberList,
		ChatRoomName: &wechat.SKBuiltinString{
			Str: proto.String(chatRoomName),
		},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeInviteChatRoomMember, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/invitechatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetQrCodeRequest 获取二维码
func GetQrCodeRequest(userInfo *baseinfo.UserInfo, id string) (*baseinfo.PackHeader, error) {
	opcode := 1
	//如果是群opcode=0
	if strings.HasSuffix(id, "@chatroom") || strings.HasSuffix(id, "@im.chatroom") {
		opcode = 0
	}
	req := wechat.GetQRCodeRequest{
		BaseRequest: GetBaseRequest(userInfo),
		UserName: &wechat.SKBuiltinString{
			Str: proto.String(id),
		},
		Style:  proto.Uint32(0),
		Opcode: proto.Uint32(uint32(opcode)),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetQrCode, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getqrcode", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 查看附近的人
func SendGetPeopleNearbyResultRequest(userInfo *baseinfo.UserInfo, longitude float32, latitude float32) (*baseinfo.PackHeader, error) {
	req := &wechat.LbsRequest{
		BaseRequest: GetBaseRequest(userInfo),
		GPSSource:   proto.Int64(0),
		Latitude:    proto.Float32(latitude),
		Longitude:   proto.Float32(longitude),
		OpCode:      proto.Uint64(1),
		Precision:   proto.Int64(65),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, 148, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/lbsfind", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetA8KeyRequest 授权链接
// @opCode ：
//
//	enum GetA8KeyOpCode
//	{
//	    MMGETA8KEY_OPENAPI = 1,
//	    MMGETA8KEY_QZONE = 3,
//	    MMGETA8KEY_REDIRECT = 2
//	}
//
// @scene ：
//
//		  enum GetA8KeyScene
//	   {
//	       MMGETA8KEY_SCENE_UNKNOW,
//	       MMGETA8KEY_SCENE_MSG,
//	       MMGETA8KEY_SCENE_TIMELINE,
//	       MMGETA8KEY_SCENE_PROFILE,
//	       MMGETA8KEY_SCENE_QRCODE,
//	       MMGETA8KEY_SCENE_QZONE,
//	       MMGETA8KEY_SCENE_OAUTH,
//	       MMGETA8KEY_SCENE_OPEN,
//	       MMGETA8KEY_SCENE_PLUGIN,
//	       MMGETA8KEY_SCENE_JUMPURL,
//	       MMGETA8KEY_SCENE_SHAKETV,
//	       MMGETA8KEY_SCENE_SCANBARCODE,
//	       MMGETA8KEY_SCENE_SCANIMAGE,
//	       MMGETA8KEY_SCENE_SCANSTREETVIEW,
//	       MMGETA8KEY_SCENE_FAV,
//	       MMGETA8KEY_SCENE_MMBIZ,
//	       MMGETA8KEY_SCENE_QQMAIL,
//	       MMGETA8KEY_SCENE_LINKEDIN,
//	       MMGETA8KEY_SCENE_SHAKETV_DETAIL,
//	       MMGETA8KEY_SCENE_BIZHOMEPAGE,
//	       MMGETA8KEY_SCENE_USBCONNECT,
//	       MMGETA8KEY_SCENE_SHORT_URL,
//	       MMGETA8KEY_SCENE_WIFI,
//	       MMGETA8KEY_SCENE_OUTSIDE_DEEPLINK,
//	       MMGETA8KEY_SCENE_PUSH_LOGIN_URL
//	   }
func GetA8KeyRequest(userInfo *baseinfo.UserInfo, opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*baseinfo.PackHeader, error) {
	req := wechat.GetA8KeyRequest{
		BaseRequest: GetBaseRequest(userInfo),
		OpCode:      proto.Uint32(opCode), //2
		ReqUrl: &wechat.SKBuiltinString{ //7
			Str: proto.String(reqUrl),
		},
		Scene:     proto.Uint32(scene), //4
		UserName:  proto.String(userInfo.GetUserName()),
		BundleID:  proto.String(""),
		NetType:   proto.String("WiFi"),
		FontScale: proto.Uint32(118), //15  118
		RequestId: proto.Uint64(uint64(time.Now().Unix())),
		CodeType:  proto.Uint32(19),
		//CodeType:    proto.Uint32(15),
		//CodeVersion: proto.Uint32(5),
		OuterUrl: proto.String(""),
		SubScene: proto.Uint32(1),
	}
	//req.FontScale=proto.Uint32(118)
	//req.CodeVersion=proto.Uint32(5)
	cgi := uint32(0)
	cgiUrl := ""

	switch getType {
	case baseinfo.ThrIdGetA8Key:
		cgi = baseinfo.MMRequestTypeThrIdGetA8Key
		cgiUrl = "/cgi-bin/micromsg-bin/3rd-geta8key"
	default:
		cgi = baseinfo.MMRequestTypeGetA8Key
		cgiUrl = "/cgi-bin/micromsg-bin/geta8key"
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, uint32(cgi), 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), cgiUrl, sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 群授权链接
func GetA8KeyGroupRequest(userInfo *baseinfo.UserInfo, opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*baseinfo.PackHeader, error) {
	req := wechat.GetA8KeyRequest{
		BaseRequest: GetBaseRequest(userInfo),
		OpCode:      proto.Uint32(opCode), //2
		A2Key: &wechat.SKBuiltinBufferT{ //3
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		AppID: &wechat.SKBuiltinString{ //4
			Str: proto.String(""),
		},
		Scope: &wechat.SKBuiltinString{ //5
			Str: proto.String(""),
		},
		State: &wechat.SKBuiltinString{ //6
			Str: proto.String(""),
		},
		ReqUrl: &wechat.SKBuiltinString{ //7
			Str: proto.String(reqUrl),
		},
		Scene:       proto.Uint32(scene),                     //10
		BundleID:    proto.String(""),                        //12
		A2KeyNew:    []byte{},                                //13
		FontScale:   proto.Uint32(118),                       //15
		NetType:     proto.String("WiFi"),                    //17
		CodeType:    proto.Uint32(19),                        //18
		CodeVersion: proto.Uint32(8),                         //19
		RequestId:   proto.Uint64(uint64(time.Now().Unix())), //20
		OuterUrl:    proto.String(""),                        //24
		SubScene:    proto.Uint32(1),                         //25
	}
	//判断Url是否是企业群
	getType = baseinfo.GetA8Key
	if strings.HasPrefix(reqUrl, "https://c.weixin.com") {
		req.CodeVersion = proto.Uint32(5)
		getType = baseinfo.ThrIdGetA8Key
	}
	cgi := uint32(0)
	cgiUrl := ""
	switch getType {
	case baseinfo.ThrIdGetA8Key:
		cgi = baseinfo.MMRequestTypeThrIdGetA8Key
		cgiUrl = "/cgi-bin/micromsg-bin/3rd-geta8key"
	default:
		cgi = baseinfo.MMRequestTypeGetA8Key
		cgiUrl = "/cgi-bin/micromsg-bin/geta8key"
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, uint32(cgi), 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), cgiUrl, sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// JSLoginRequest 授权小程序 返回Code
// @appId :要授权wxappid
func JSLoginRequest(userInfo *baseinfo.UserInfo, appId string) (*baseinfo.PackHeader, error) {
	req := wechat.JSLoginRequest{
		BaseRequest: GetBaseRequest(userInfo),
		AppId:       proto.String(appId),
		Scope:       proto.String("snsapi_login"),
		LoginType:   proto.Int32(4),
		Url:         proto.String("https://open.weixin.qq.com/connect/confirm?uuid=021extwqXdPRRlbJ"),
		VersionType: proto.Int32(0),
		WxaExternalInfo: &wechat.WxaExternalInfo{
			Scene:     proto.Int32(1001),
			SourceEnv: proto.Int32(1),
		},
	}
	/*req := wechat.JSLoginRequest{
		BaseRequest: GetBaseRequest(userInfo),
		AppId:       proto.String(appId),
		LoginType:   proto.Int32(1),
		VersionType: proto.Int32(0),
		WxaExternalInfo: &wechat.WxaExternalInfo{
			Scene:     proto.Int32(1066),
			SourceEnv: proto.Int32(1),
		},
	}*/
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeJSLogin, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/js-login", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// JSOperateWxDataRequest 授权小程序后返回 encryptedData,iv等信息
func JSOperateWxDataRequest(userInfo *baseinfo.UserInfo, appId string) (*baseinfo.PackHeader, error) {
	req := wechat.JSOperateWxDataRequest{
		BaseRequest: GetBaseRequest(userInfo),
		AppId:       proto.String(appId),
		Data:        []byte("{\"with_credentials\":true,\"from_component\":true,\"data\":{\"lang\":\"zh_CN\"},\"api_name\":\"webapi_getuserinfo\"}"),
		GrantScope:  proto.String("scope.userInfo"),
		Opt:         proto.Int32(1),
		VersionType: proto.Int32(0),
		WxaExternalInfo: &wechat.WxaExternalInfo{
			Scene:     proto.Int32(1001),
			SourceEnv: proto.Int32(2),
		},
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeJSOperateWxData, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/js-operatewxdata", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SdkOauthAuthorizeRequest 授权app应用
func SdkOauthAuthorizeRequest(userInfo *baseinfo.UserInfo, appId string, sdkName string, packageName string) (*baseinfo.PackHeader, error) {
	req := wechat.SdkOauthAuthorizeReq{
		BaseRequest: GetBaseRequest(userInfo),
		AppId:       proto.String(appId),
		Tag3:        proto.String("snsapi_userinfo"),
		Tag4:        proto.String(sdkName),     //wechat_sdk_demo_test
		Tag5:        proto.String(packageName), //"com.yimu.renwuxiongObject"
		Tag8:        proto.String(""),
		Tag9:        proto.String(""),
		Tag10:       proto.String(""),
		Tag11:       proto.String(""),
		Tag12:       proto.Uint32(0),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeSdkOauthAuthorize, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/sdk_oauth_authorize", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSearchContactRequest 搜索联系人
func SendSearchContactRequest(userInfo *baseinfo.UserInfo, opCode, fromScene, searchScene uint32, userName string) (*baseinfo.PackHeader, error) {
	req := wechat.SearchContactRequest{
		BaseRequest: GetBaseRequest(userInfo),
		UserName: &wechat.SKBuiltinString{
			Str: proto.String(userName),
		},
		OpCode:      proto.Uint32(opCode),
		FromScene:   proto.Uint32(fromScene),
		SearchScene: proto.Uint32(searchScene),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeSearchContact, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/searchcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// VerifyUserRequest 好友验证
func VerifyUserRequest(userInfo *baseinfo.UserInfo, opCode uint32, verifyContent string, scene byte, V1, V2, ChatRoomUserName string) (*baseinfo.PackHeader, error) {
	// clientcheckdata
	var extSpamInfo wechat.SKBuiltinString_
	if userInfo.DeviceInfo != nil {
		extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	} else {
		extSpamInfo.Buffer = GetExtPBSpamInfoDataA16(userInfo)
	}

	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen

	//ChatRoomUserName
	userTicket := V2
	if ChatRoomUserName != "" {
		userTicket = ""
	}
	req := wechat.VerifyUserRequest{
		BaseRequest:        GetBaseRequest(userInfo),
		OpCode:             proto.Uint32(opCode),
		VerifyUserListSize: proto.Uint32(1),
		VerifyUserList: []*wechat.VerifyUser{{
			Value:               proto.String(V1),
			VerifyUserTicket:    proto.String(userTicket),
			AntispamTicket:      proto.String(V2),
			FriendFlag:          proto.Uint32(0),
			ChatRoomUserName:    proto.String(ChatRoomUserName),
			SourceUserName:      proto.String(""),
			SourceNickName:      proto.String(""),
			ScanQrcodeFromScene: proto.Uint32(0),
			ReportInfo:          proto.String(""),
			OuterUrl:            proto.String(""),
			SubScene:            proto.Uint32(0),
		}},
		VerifyContent:  proto.String(verifyContent),
		SceneListCount: proto.Uint32(1),
		SceneList:      []byte{scene},
		ExtSpamInfo:    &extSpamInfo,
		//NeedConfirm:    proto.Uint32(1),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeVerifyUser, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/verifyuser", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// UploadMContact 上传通讯录
func UploadMContactRequest(userInfo *baseinfo.UserInfo, mobile string, mobiles []string) (*baseinfo.PackHeader, error) {
	mobileList := make([]*wechat.Mobile, 0)
	for _, mobile := range mobiles {
		mobileList = append(mobileList, &wechat.Mobile{
			V: proto.String(mobile),
		})
	}
	req := wechat.UploadMContactRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		UserName:       proto.String(userInfo.WxId),
		Opcode:         proto.Int32(1),
		Mobile:         proto.String(mobile),
		MobileListSize: proto.Int32(int32(len(mobileList))),
		MobileList:     mobileList,
		EmailListSize:  proto.Int32(0),
		EmailList:      []*wechat.MEmail{},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeUploadMContact, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadmcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetMFriendRequest 获取通讯录
func GetMFriendRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	req := wechat.GetMFriendRequest{
		BaseRequest: GetBaseRequest(userInfo),
		OpType:      proto.Uint32(0),
		MD5:         proto.String(uuid.New().String()),
		Scene:       proto.Uint32(0),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetMFriend, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getmfriend", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetCertRequest 获取证书
func GetCertRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	//userInfo.SessionKey = baseutils.RandomBytes(16)
	req := wechat.GetCertRequest{
		BaseRequest: GetBaseRequest(userInfo),
		AesEncryptKey: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(16),
			Buffer: userInfo.SessionKey[:16],
		},
		Version: proto.Uint32(135),
	}
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 381, 7)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getcert", sendEncodeData)
	log.Println(hex.EncodeToString(resp))
	log.Println(hex.EncodeToString(userInfo.SessionKey))
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetSdkOauthAuthorizeConfirmRequest
func GetSdkOauthAuthorizeConfirmRequest(userInfo *baseinfo.UserInfo, AppId, appName, appNamePack string) (*baseinfo.PackHeader, error) {
	req := wechat.SdkOauthAuthorizeConfirmNewReq{
		BaseRequest:     GetBaseRequest(userInfo),
		Opt:             proto.Uint32(1),
		Scope:           []string{"snsapi_userinfo"},
		AppId:           proto.String(AppId),
		State:           proto.String(appName),
		BundleId:        proto.String(appNamePack),
		AvatarId:        proto.Uint32(0),
		UniversalLink:   proto.String(""),
		OpenSdkVersion:  proto.String(""),
		SdkToken:        proto.String(""),
		OpenSdkBundleId: proto.String(""),
		SdkTokenChk:     proto.Uint32(0),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 1346, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/sdk_oauth_authorize_confirm", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetQRConnectAuthorizeRequest 获取授权二维码链接组包
func GetQRConnectAuthorizeRequest(userInfo *baseinfo.UserInfo, url string) (*baseinfo.PackHeader, error) {
	req := wechat.QRConnectAuthorizeReq{
		BaseRequest: GetBaseRequest(userInfo),
		OAuthUrl:    proto.String(url),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 2543, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/qrconnect_authorize", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetQRConnectAuthorizeConfirmRequest 授权二维码链接确认组包
func GetQRConnectAuthorizeConfirmRequest(userInfo *baseinfo.UserInfo, url string) (*baseinfo.PackHeader, error) {
	req := wechat.QRConnectAuthorizeConfirmReq{
		BaseRequest: GetBaseRequest(userInfo),
		OAuthUrl:    proto.String(url),
		Opt:         proto.Uint32(1),
		Scope:       []string{"snsapi_login"},
		AvatarId:    proto.Uint32(0),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 1137, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/qrconnect_authorize_confirm", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 授权链接
func SendGetMpA8Request(userInfo *baseinfo.UserInfo, url string, opcode uint32) (*baseinfo.PackHeader, error) {
	req := wechat.GetA8KeyRequest{
		BaseRequest: GetBaseRequest(userInfo),
		CodeType:    proto.Uint32(19),
		CodeVersion: proto.Uint32(9),
		Flag:        proto.Uint32(0),
		FontScale:   proto.Uint32(100),
		NetType:     proto.String("WIFI"),
		OpCode:      proto.Uint32(opcode),
		UserName:    proto.String(userInfo.WxId),
		ReqUrl: &wechat.SKBuiltinString{
			Str: proto.String(url),
		},
		FriendQq: proto.Uint32(0),
		Scene:    proto.Uint32(4),
		SubScene: proto.Uint32(1),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 233, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/mp-geta8key", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetOnlineInfoRequest 获取登录信息组包
func GetOnlineInfoRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	req := wechat.GetOnlineInfoRequest{
		BaseRequest: GetBaseRequest(userInfo),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetOnlineInfo, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getonlineinfo", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// GetRevokeMsgRequest 撤销消息
func GetRevokeMsgRequest(userInfo *baseinfo.UserInfo, newMsgId string, clientMsgId uint64, toUserName string) (*baseinfo.PackHeader, error) {
	msgId, _ := strconv.ParseUint(newMsgId, 10, 64)
	req := wechat.RevokeMsgRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		ClientMsgId:    proto.String(fmt.Sprintf("%v_%v", time.Now().Unix(), newMsgId)),
		NewClientMsgId: proto.Uint32(uint32(time.Now().Unix())),
		CreateTime:     proto.Uint32(uint32(time.Now().Unix())),
		SvrMsgId:       proto.Uint64(clientMsgId),
		FromUserName:   proto.String(userInfo.GetUserName()),
		ToUserName:     proto.String(toUserName),
		IndexOfRequest: proto.Uint32(26),
		SvrNewMsgId:    proto.Uint64(msgId),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeRevokeMsg, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/revokemsg", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// UploadHeadImage 修改头像
func UploadHeadImage(userInfo *baseinfo.UserInfo, base64Image string) (*baseinfo.PackHeader, error) {

	ImgData := strings.Split(base64Image, ",")
	var ImgBase64 []byte
	if len(ImgData) > 1 {
		ImgBase64, _ = base64.StdEncoding.DecodeString(ImgData[1])
	} else {
		ImgBase64, _ = base64.StdEncoding.DecodeString(base64Image)
	}
	ImgStream := bytes.NewBuffer(ImgBase64)
	Startpos := 0
	datalen := 30000
	datatotalength := ImgStream.Len()
	ImgHash := GetFileMD5Hash(ImgBase64)
	I := 0
	resps := make([]byte, 0)
	for {
		Startpos = I * datalen
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}
		Databuff := make([]byte, count)
		_, _ = ImgStream.Read(Databuff)
		req := wechat.UploadHDHeadImgRequest{
			BaseRequest: GetBaseRequest(userInfo),
			TotalLen:    proto.Uint32(uint32(datatotalength)),
			StartPos:    proto.Uint32(uint32(Startpos)),
			HeadImgType: proto.Uint32(1),
			Data: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			ImgHash: proto.String(ImgHash),
		}
		// 打包发送数据
		srcData, _ := proto.Marshal(&req)
		sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeUploadHDHeadImg, 5)
		resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadhdheadimg", sendEncodeData)
		if err != nil {
			return nil, err
		}
		resps = resp
		I++
	}
	// 打包发送数据
	return DecodePackHeader(resps, nil)
}

func GetFileMD5Hash(Data []byte) string {
	hash := md5.New()
	hash.Write(Data)
	retVal := hash.Sum(nil)
	return hex.EncodeToString(retVal)
}

// GetVerifyPwdRequest 验证密码
func GetVerifyPwdRequest(userInfo *baseinfo.UserInfo, pwd string) (*baseinfo.PackHeader, error) {
	req := wechat.VerifyPwdRequest{
		BaseRequest: GetBaseRequest(userInfo),
		OpCode:      proto.Uint32(1),
		Pwd1:        proto.String(baseutils.Md5Value(pwd)),
		Pwd2:        proto.String(baseutils.Md5Value(pwd)),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeVerifyPassword, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newverifypasswd", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendSetPwdRequest 修改密码
func SendSetPwdRequest(userInfo *baseinfo.UserInfo, ticket, newPwd string, OpCode uint32) (*baseinfo.PackHeader, error) {
	req := wechat.SetPwdRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Password:    proto.String(baseutils.Md5Value(newPwd)), //2
		AutoAuthKey: &wechat.BufferT{ //4
			ILen:   proto.Uint32(uint32(len(userInfo.AutoAuthKey))),
			Buffer: userInfo.AutoAuthKey,
		},
	}
	if OpCode == 0 {
		req.Ticket = proto.String(ticket) //3
	}
	if OpCode != 0 {
		req.OpCode = proto.Uint32(OpCode)
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeSetPassword, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newsetpasswd", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendUploadVoiceRequest 发送语音
func SendUploadVoiceNewRequest(userInfo *baseinfo.UserInfo, ToWxId string, Startpos int, Databuff []byte, ClientImgId string, VoiceTime int32, Type int32, endFlag int) (*baseinfo.PackHeader, error) {
	var req = wechat.UploadVoiceRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		FromUserName: proto.String(userInfo.WxId),
		ToUserName:   proto.String(ToWxId),
		Offset:       proto.Uint32(uint32(Startpos)),
		Length:       proto.Uint32(uint32(len(Databuff))),
		ClientMsgId:  proto.String(ClientImgId),
		MsgId:        proto.Uint32(0),
		VoiceLength:  proto.Int32(VoiceTime * 1000),
		VoiceFormat:  proto.Int32(Type),
		Data: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(uint32(len(Databuff))),
			Buffer: Databuff,
		},
		EndFlag: proto.Uint32(uint32(endFlag)),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeUploadVoiceNew, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadvoice", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// SendUploadVoiceRequest 发送语音
/*func SendUploadVoiceRequest(userInfo *baseinfo.UserInfo, toUserName string, data []byte, totalLen, startPos uint32, clientImgId string, voiceLen uint32, voiceFormat uint32) (*baseinfo.PackHeader, error) {
	endFlag := uint32(0)
	if startPos+uint32(len(data)) >= totalLen {
		endFlag = 1
	}

	var req = wechat.UploadVoiceRequest{
		BaseRequest:  GetBaseRequest(userInfo),
		FromUserName: proto.String(userInfo.GetUserName()),
		ToUserName:   proto.String(toUserName),
		ClientMsgId:  proto.String(clientImgId),
		VoiceFormat:  proto.Uint32(voiceFormat),
		VoiceLength:  proto.Uint32(voiceLen),
		Length:       proto.Uint32(uint32(len(data))),
		Data: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(uint32(len(data))),
			Buffer: data,
		},
		Offset:  proto.Uint32(startPos),
		EndFlag: proto.Uint32(endFlag),
		MsgId:   proto.Uint32(0),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeUploadVoice, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadvoice", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}*/

// SendNewInitSyncRequest 首次登录初始化
func SendNewInitSyncRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	mgr := userInfo.SyncKeyMgr()
	if mgr.CurKey() == nil {
		mgr.SetCurKey(&wechat.BufferT{
			ILen:   proto.Uint32(0),
			Buffer: make([]byte, 0),
		})
	}
	if mgr.MaxKey() == nil {
		mgr.SetMaxKey(&wechat.BufferT{
			ILen:   proto.Uint32(0),
			Buffer: make([]byte, 0),
		})
	}
	req := &wechat.NewInitRequest{}
	req.BaseRequest = GetBaseRequest(userInfo)
	//req.BaseRequest.Scene = proto.Uint32(0)
	req.UserName = proto.String(userInfo.GetUserName())
	req.CurrentSynckey = mgr.CurKey()
	req.MaxSynckey = mgr.MaxKey()
	Language := "id" //zh_CN
	if userInfo.DeviceInfo != nil {
		Language = userInfo.DeviceInfo.Language
	}
	req.Language = proto.String(Language)

	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeNewInit, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/newinit", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 设置微信号
func SetWechatRequest(userInfo *baseinfo.UserInfo, alisa string) (*baseinfo.PackHeader, error) {
	var req = wechat.GeneralSetRequest{
		BaseRequest: GetBaseRequest(userInfo),
		SetType:     proto.Int32(1),
		SetValue:    proto.String(alisa),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 177, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/generalset", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取设备
func GetBoundHardDeviceRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	var req = wechat.GetBoundHardDevicesRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Version:     proto.Uint32(0),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 539, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getboundharddevices", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 换绑定手机
func SendBindingMobileRequest(mobile, verifyCode string, userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	SafeDeviceName := "iPhone"
	SafeDeviceType := "iPhone"
	if userInfo.DeviceInfo != nil {
		SafeDeviceName = userInfo.DeviceInfo.DeviceName
		SafeDeviceType = userInfo.DeviceInfo.OsType
	} else {
		SafeDeviceName = baseinfo.AndroidDeviceType
		SafeDeviceType = baseinfo.AndroidDeviceType
	}
	var req = wechat.BindOpMobileRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		UserName:       proto.String(userInfo.WxId),
		Mobile:         proto.String(mobile),
		Opcode:         proto.Uint32(19),
		Verifycode:     proto.String(verifyCode),
		SafeDeviceName: proto.String(SafeDeviceName),
		SafeDeviceType: proto.String(SafeDeviceType),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 132, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/bindopmobile", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 发送验证码
func SendVerifyMobileRequest(mobile string, opcode uint32, userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	SafeDeviceName := "iPhone"
	SafeDeviceType := "iPhone"
	var req = wechat.BindOpMobileRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		UserName:       proto.String(userInfo.WxId),
		Mobile:         proto.String(mobile),
		DialFlag:       proto.Uint32(0),
		Opcode:         proto.Uint32(opcode),
		SafeDeviceName: proto.String(SafeDeviceName),
		SafeDeviceType: proto.String(SafeDeviceType),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 132, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/bindopmobile", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取步数列表
func SendGetUserRankLikeCountRequest(userInfo *baseinfo.UserInfo, rankId string) (*baseinfo.PackHeader, error) {
	var req = wechat.GetUserRankLikeCountRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Username:    proto.String(userInfo.WxId),
		LatestRank:  proto.Bool(true),
		RankId:      proto.String(rankId),
		AppUsername: proto.String("wx7fa037cc7dfabad5"),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 1042, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmbiz-bin/rank/getuserranklike", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 上传步数
func UploadStepSetRequestRequest(userInfo *baseinfo.UserInfo, deviceID string, deviceType string, number uint64) (*baseinfo.PackHeader, error) {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).Unix()
	var req = wechat.UploadDeviceStepRequest{
		BaseRequest: GetBaseRequest(userInfo),
		DeviceID:    proto.String(deviceID),
		DeviceType:  proto.String(deviceType),
		FromTime:    proto.Uint32(uint32(startTime)),
		ToTime:      proto.Uint32(uint32(endTime)),
		StepCount:   proto.Uint32(uint32(number)),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 1261, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/mmoc-bin/hardware/uploaddevicestep", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 提取企业 wx 详情
func SendQWContactRequest(userInfo *baseinfo.UserInfo, openIm, chatRoom, t string) (*baseinfo.PackHeader, error) {
	req := wechat.GetQYContactRequest{}
	req.Wxid = proto.String(openIm)
	if chatRoom != "" {
		req.Room = proto.String(chatRoom)
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 881, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getopenimcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 提取全部的企业通寻录
func SendQWSyncContactRequest(userInfo *baseinfo.UserInfo) (*baseinfo.PackHeader, error) {
	ck := make([]*wechat.SyncKey_, 0)
	ck = append(ck, &wechat.SyncKey_{
		SyncKey: proto.Int64(0),
		Type:    proto.Uint32(400),
	})
	key := wechat.SyncMsgKey{
		Len: proto.Uint32(1),
		MsgKey: &wechat.SyncKey{
			Size: proto.Uint32(1),
			Type: ck,
		},
	}
	keyMar, _ := proto.Marshal(&key)
	var req = wechat.QYSyncRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Selector:    proto.Int64(0x200000),
		Key:         keyMar,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 810, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/openimsync", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 创建企业群
func SendQWCreateChatRoomRequest(userInfo *baseinfo.UserInfo, userList []string) (*baseinfo.PackHeader, error) {
	memberlist := make([]*wechat.Openimcontact, 0)
	for _, val := range userList {
		memberlist = append(memberlist, &wechat.Openimcontact{
			UserName: proto.String(val),
		})
	}
	var req = wechat.CreateQYChatRoomRequest{
		MemberList: memberlist,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 371, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/createopenimchatroom", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 向企业微信打招呼
func SendQWApplyAddContactRequest(userInfo *baseinfo.UserInfo, toUserName, v1, Content string) (*baseinfo.PackHeader, error) {
	var req = wechat.QYVerifyUserRequest{
		Wxid:    proto.String(toUserName),
		V1:      proto.String(v1),
		Content: proto.String(Content),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 0xF3, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/sendopenimverifyrequest", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 单向加企业微信
func SendQWAddContactRequest(userInfo *baseinfo.UserInfo, toUserName, v1, Content string) (*baseinfo.PackHeader, error) {
	var req = wechat.QYVAddUserRequest{
		Wxid: proto.String(toUserName),
		V1:   proto.String(v1),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 667, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/addopenimcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 提取所有企业群
func SendQWSyncChatRoomRequest(userInfo *baseinfo.UserInfo, key string) (*baseinfo.PackHeader, error) {
	ck := make([]*wechat.SyncKey_, 0)
	ck = append(ck, &wechat.SyncKey_{
		SyncKey: proto.Int64(0),
		Type:    proto.Uint32(0),
	})
	keys := wechat.SyncMsgKey{
		Len: proto.Uint32(1),
		MsgKey: &wechat.SyncKey{
			Size: proto.Uint32(1),
			Type: ck,
		},
	}
	keyMar, _ := proto.Marshal(&keys)
	var req = wechat.QYSyncRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Selector:    proto.Int64(2097152), //0x200000  2097152
		Key:         keyMar,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 810, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/openimsync", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 转让企业微信群
func SendQWChatRoomTransferOwnerRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName string) (*baseinfo.PackHeader, error) {
	var req = wechat.QWTransferChatRoomOwnerRequest{
		Username: proto.String(chatRoomName),
		Owner: &wechat.Openimcontact{
			UserName: proto.String(toUserName),
		},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 811, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/modopenimchatroomowner", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 直接拉好友进群
func SendQWAddChatRoomMemberRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName []string) (*baseinfo.PackHeader, error) {
	list := make([]*wechat.Openimcontact, 0)
	for _, val := range toUserName {
		list = append(list, &wechat.Openimcontact{
			UserName: proto.String(val),
		})
	}
	var req = wechat.QYAddChatRoomRequest{
		ChatRoomName: proto.String(chatRoomName),
		MemberList:   list,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 0x32E, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/addopenimchatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 发送群邀请链接
func SendQWInviteChatRoomMemberRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName []string) (*baseinfo.PackHeader, error) {
	list := make([]*wechat.Openimcontact, 0)
	for _, val := range toUserName {
		list = append(list, &wechat.Openimcontact{
			UserName: proto.String(val),
		})
	}
	var req = wechat.InviteQYChatRoomRequest{
		ChatRoomName: proto.String(chatRoomName),
		MemberList:   list,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 887, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/inviteopenimchatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 删除企业群群成员
func SendQWDelChatRoomMemberRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName []string) (*baseinfo.PackHeader, error) {
	list := make([]*wechat.Openimcontact, 0)
	for _, val := range toUserName {
		list = append(list, &wechat.Openimcontact{
			UserName: proto.String(val),
		})
	}
	var req = wechat.QYDelChatRoomMemberRequest{
		Username:   proto.String(chatRoomName),
		MemberList: list,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 0x3AF, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delopenimchatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 提取企业群全部成员
func SendQWGetChatRoomMemberRequest(userInfo *baseinfo.UserInfo, chatRoomName string) (*baseinfo.PackHeader, error) {
	var req = wechat.GetQYChatroomMemberDetailRequest{
		ChatroomUserName: proto.String(chatRoomName),
		ClientVersion:    proto.Uint64(0),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 942, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getopenimchatroommemberdetail", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 提取企业群名称公告设定等信息
func SendQWGetChatroomInfoRequest(userInfo *baseinfo.UserInfo, chatRoomName string) (*baseinfo.PackHeader, error) {
	var req = wechat.Openimcontact{
		UserName: proto.String(chatRoomName),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 407, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getopenimchatroomcontact", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 提取企业群二维码
func SendQWGetChatRoomQRRequest(userInfo *baseinfo.UserInfo, chatRoomName string) (*baseinfo.PackHeader, error) {
	var req = wechat.Openimcontact{
		UserName: proto.String(chatRoomName),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 890, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getopenimchatroomqrcode", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 增加企业管理员
func SendQWAppointChatRoomAdminRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName []string) (*baseinfo.PackHeader, error) {
	var req = wechat.QYChatRoomAdminRequest{
		ChatRoomName: proto.String(chatRoomName),
		MemberList:   toUserName,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 776, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/addopenimchatroomadmin", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 移除企业群管理员
func SendQWDelChatRoomAdminRequest(userInfo *baseinfo.UserInfo, chatRoomName string, toUserName []string) (*baseinfo.PackHeader, error) {
	var req = wechat.QYChatRoomAdminRequest{
		ChatRoomName: proto.String(chatRoomName),
		MemberList:   toUserName,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 3677, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delopenimchatroomadmin", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 同意进企业群
func SendQWAcceptChatRoomRequest(userInfo *baseinfo.UserInfo, link string, opcode uint32) (*baseinfo.PackHeader, error) {
	var req = wechat.GetA8KeyRequest{}
	req.BaseRequest = GetBaseRequest(userInfo)
	req.CodeType = proto.Uint32(0)
	req.CodeVersion = proto.Uint32(8)
	req.Flag = proto.Uint32(0)
	req.FontScale = proto.Uint32(100)
	req.NetType = proto.String("WIFI")
	req.OpCode = proto.Uint32(opcode)
	req.UserName = proto.String(userInfo.WxId)
	req.ReqUrl = &wechat.SKBuiltinString{
		Str: proto.String(link),
	}
	req.FriendQq = proto.Uint32(0)
	req.Scene = proto.Uint32(37)
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 233, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/geta8key", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 设定企业群管理审核进群
func SendQWAdminAcceptJoinChatRoomSetRequest(userInfo *baseinfo.UserInfo, link string, url string) (*baseinfo.PackHeader, error) {
	var req = wechat.QYChatRoomAdminRequest{}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 3677, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/delopenimchatroomadmin", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 群管理批准进企业群
func SendQWAdminAcceptJoinChatRoomRequest(userInfo *baseinfo.UserInfo, chatRoomName, key, toUserName string, toUserNames []string) (*baseinfo.PackHeader, error) {
	list := make([]*wechat.Openimcontact, 0)
	for _, val := range toUserNames {
		list = append(list, &wechat.Openimcontact{
			UserName: proto.String(val),
		})
	}
	var req = wechat.QYAdminAddRequest{
		Room: proto.String(chatRoomName),
		Key:  proto.String(key),
		Username: &wechat.Openimcontact{
			UserName: proto.String(toUserName),
		},
		Usernamelist: list,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 0x3AD, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/approveaddopenimchatroommember", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 视频号搜索
func SendGetFinderSearchRequest(userInfo *baseinfo.UserInfo, Index uint32, Userver int32, UserKey string, Uuid string) (*baseinfo.PackHeader, error) {
	var req = wechat.FinderSearchRequest{
		BaseRequest: GetBaseRequest(userInfo),
		UserKey:     proto.String(UserKey),
		Offset:      proto.Uint32(Index),
		Scene:       proto.Uint32(0),
		Uuid:        proto.String(Uuid),
		FinderTxRequest: &wechat.FinderTxRequest{
			Userver: proto.Int32(Userver),
		},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 3820, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/findersearch", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 视频号个人中心
func SendFinderUserPrepareRequest(userInfo *baseinfo.UserInfo, Userver int32) (*baseinfo.PackHeader, error) {
	var req = wechat.FinderUserPrepareRequest{
		BaseRequest: GetBaseRequest(userInfo),
		Scene:       proto.Int32(0),
		FinderTxRequest: &wechat.FinderTxRequest{
			Userver: proto.Int32(Userver),
		},
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 3761, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/finderuserprepare", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 视频号关注or取消关注
func SendFinderFollowRequest(userInfo *baseinfo.UserInfo, FinderUserName string, OpType int32, RefObjectId string, Cook string, Userver int32, PosterUsername string) (*baseinfo.PackHeader, error) {
	//refId, err := strconv.ParseUint(RefObjectId, 10, 64)
	T := time.Now().Unix()
	var req = wechat.FinderFollowRequest{
		FinderUsername: proto.String(FinderUserName),
		OpType:         proto.Int32(OpType),
		RefObjectId:    proto.Uint64(0),
		PosterUsername: proto.String(""),
		FinderReq: &wechat.FinderTxRequest{
			Userver: proto.Int32(Userver),
			Scene:   proto.Int32(6),
			T:       proto.Int32(1),
			G: &wechat.FinderZd{
				G1: proto.String(fmt.Sprintf("Finder_Enter%v", T)),
				G2: proto.String(fmt.Sprintf("4-%v", T)),
				G3: proto.String(fmt.Sprintf(`{"sessionId":"109_%v#$0_%v#"}`, T, T)),
			},
			Tg: proto.Int64(time.Now().Unix()),
		},
		Cook:      proto.String(Cook),
		EnterType: proto.Int32(7),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 3867, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/finderfollow", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 查看视频号首页
func SendTargetUserPageRequest(userInfo *baseinfo.UserInfo, target string, lastBuffer string) (*baseinfo.PackHeader, error) {
	T := time.Now().Unix()
	LastBuffer := []byte{}
	if lastBuffer != "" && lastBuffer != "string" {
		key, _ := base64.StdEncoding.DecodeString(lastBuffer)
		LastBuffer = key
	}
	var req = wechat.FinderUserPageRequest{
		Username:      proto.String(target),
		MaxId:         proto.Uint64(0),
		FirstPageMd5:  proto.String(""),
		NeedFansCount: proto.Int(0),
		FinderBasereq: &wechat.FinderBaseRequest{
			Userver: proto.Int(12),
			Scene:   proto.Int(20),
			ExptFla: proto.Uint32(1),
			CtxInfo: &wechat.ClientContextInfo{
				ContextId:         proto.String(fmt.Sprintf("Finder_Enter%v", T)),
				ClickTabContextId: proto.String(fmt.Sprintf("4-%v", T)),
				ClientReportBuff:  proto.String(fmt.Sprintf(`{"sessionId":"143_%v#1900580_%v|$2_%v#"}`, T, T, T)),
			},
		},
		LastBuffer: LastBuffer,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 3736, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/finderuserpage", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取验证码
func SendWxBindOpMobileForRegRequest(userInfo *baseinfo.UserInfo, opcode int64, mobile, verifycode string) (*baseinfo.PackHeader, error) {
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	var strClientSeqID = string(userInfo.DeviceInfo.Imei + "-" + tmpTimeStr)
	var req = wechat.BindOpMobileForRegRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		Mobile:         proto.String(mobile),
		Opcode:         proto.Int64(opcode),
		Verifycode:     proto.String(verifycode),
		SafeDeviceName: proto.String("iPhone"),
		SafeDeviceType: proto.String("iPhone"),
		RandomEncryKey: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(uint32(len(userInfo.SessionKey))),
			Buffer: userInfo.SessionKey,
		},
		Language:          proto.String("zh_CN"),
		InputMobileReTrYs: proto.Uint32(0),
		AdjustRet:         proto.Uint32(0),
		ClientSeqID:       proto.String(strClientSeqID),
		DialLang:          proto.String(""),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 145, 7)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/bindopmobileforreg", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}
