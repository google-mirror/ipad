package bizcgi

import (
	"errors"
	"strings"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/cecdh"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
)

// SendAutoAuthRequest 发送Token登陆请求
func SendAutoAuthRequest(wxAccount wxface.IWXAccount) (*wechat.ManualAuthResponse, error) {
	log.Info("发起二次登录")
	userInfo := wxAccount.GetUserInfo()
	userInfo.EcPublicKey, userInfo.EcPrivateKey = cecdh.GenerateEccKey()
	autoAuthKey := &wechat.AutoAuthKey{}
	err := proto.Unmarshal(userInfo.AutoAuthKey, autoAuthKey)
	if err != nil {
		log.Warn(wxAccount.GetUserInfo().UUID, "二次登录失败！", err.Error())
		return nil, err
	}
	userInfo.SessionKey = autoAuthKey.EncryptKey.Buffer
	// 获取AutoAuthRsaReqData
	rsaReqData := clientsdk.GetAutoAuthRsaReqDataMarshal(userInfo)
	aesReqData := clientsdk.GetAutoAuthAesReqDataMarshal(userInfo)

	// 开始打包数据
	// 加密压缩 rsaReqData
	rsaEncodeData := baseutils.CompressAndRsaByVer(rsaReqData, userInfo.GetLoginRsaVer())
	rsaAesEncodeData := baseutils.CompressAes(userInfo.SessionKey, rsaReqData)

	// 加密压缩 aesReqData
	aesEncodeData := baseutils.CompressAes(userInfo.SessionKey, aesReqData)

	body := make([]byte, 0)
	// rsaReqBuflen
	tmpBuf := baseutils.Int32ToBytes(uint32(len(rsaReqData)))
	body = append(body, tmpBuf[0:]...)

	// aesReqBuf len
	tmpBuf = baseutils.Int32ToBytes(uint32(len(aesReqData)))
	body = append(body, tmpBuf[0:]...)

	// rsaEncodeData len
	tmpBuf = baseutils.Int32ToBytes(uint32(len(rsaEncodeData)))
	body = append(body, tmpBuf[0:]...)

	// rsaAesEncodeData len
	tmpBuf = baseutils.Int32ToBytes(uint32(len(rsaAesEncodeData)))
	body = append(body, tmpBuf[0:]...)

	// rsaEncodeData
	body = append(body, rsaEncodeData[0:]...)
	body = append(body, rsaAesEncodeData[0:]...)
	body = append(body, aesEncodeData[0:]...)
	// 发送请求
	sendData := clientsdk.Pack(userInfo, body, baseinfo.MMRequestTypeAutoAuth, 9)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 254,
		CgiUrl: "/cgi-bin/micromsg-bin/autoauth",
		Data:   sendData,
	}
	result, err := wxlink.WXShortSend(wxAccount, longReq)
	if err != nil {
		log.Error("二次登录失败", err.Error())
		return nil, err
	}
	response, _ := result.(wechat.ManualAuthResponse)
	return &response, nil
}

// SendManualAuth 发送登陆请求
func SendManualAuth(wxAccount wxface.IWXAccount, newpass string, wxid string) (*wechat.ManualAuthResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	// 序列化
	accountData, err := clientsdk.GetManualAuthAccountDataReq(userInfo, newpass, wxid)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") && userInfo.DeviceInfo != nil {
		return sendManualAuthByAccountData(wxAccount, accountData)
	}
	return SendManualAuthA16(wxAccount, accountData)
}

// SendManualAuthA16发送A16登录请求
func SendManualAuthA16(wxAccount wxface.IWXAccount, accountData []byte) (*wechat.ManualAuthResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	// 发送登陆请求
	sendData := clientsdk.GetManualAuthA16Req(userInfo, accountData)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 253,
		CgiUrl: "/cgi-bin/micromsg-bin/manualauth",
		Data:   sendData,
	}
	result, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		log.Error("登录失败", err.Error())
		return nil, err
	}
	response, _ := result.(wechat.ManualAuthResponse)
	return &response, nil
}

// sendManualAuthByAccountData 发送ManualAuth请求
func sendManualAuthByAccountData(wxAccount wxface.IWXAccount, accountData []byte) (*wechat.ManualAuthResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	// 发送登陆请求
	sendData := clientsdk.GetManualAuthByAccountDataReq(userInfo, accountData)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 253,
		CgiUrl: "/cgi-bin/micromsg-bin/manualauth",
		Data:   sendData,
	}
	result, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		log.Error("登录失败", err.Error())
		return nil, err
	}
	response, _ := result.(wechat.ManualAuthResponse)
	return &response, nil
}

// SendGetLoginQrcodeRequest 获取登录二维码
func SendGetLoginQrcodeRequest(wxAccount wxface.IWXAccount) (*wechat.LoginQRCodeResponse, error) {
	// 发送请求
	tmpUserInfo := wxAccount.GetUserInfo()
	tmpUserInfo.DeviceInfo.OsType = baseinfo.DeviceTypeIpad
	tmpUserInfo.DeviceInfo.DeviceName = "iPad"
	tmpUserInfo.DeviceInfo.IphoneVer = "iPad4,7"
	reqData := clientsdk.GetLoginQRCodeReq(tmpUserInfo)
	//发送给长链接请求去处理  /cgi-bin/micromsg-bin/getloginqrcode
	longReq := &clientsdk.WXLongRequest{
		OpCode: mmtls.MMLongOperationGetQrcode,
		CgiUrl: "/cgi-bin/micromsg-bin/getloginqrcode",
		Data:   reqData,
	}
	result, err := wxlink.WXShortSend(wxAccount, longReq)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	qrCodeResponse := result.(wechat.LoginQRCodeResponse)
	return &qrCodeResponse, nil
}

// SendCheckLoginQrcodeRequest 检测二维码状态
func SendCheckLoginQrcodeRequest(wxAccount wxface.IWXAccount, qrcodeUUID string) (*wechat.LoginQRCodeNotify, error) {
	// 获取检测二维码状态请求
	tmpUserInfo := wxAccount.GetUserInfo()
	qrAesKey := wxcore.WxInfoCache.GetQrcodeInfo(tmpUserInfo.QrUuid)
	if qrAesKey == nil {
		return nil, errors.New("二维码失效")
	}
	reqData, err := clientsdk.GetCheckLoginQRCodeReq(tmpUserInfo, qrcodeUUID)
	if err != nil {
		return nil, err
	}

	// 发送给长链接请求去处理
	longReq := &clientsdk.WXLongRequest{
		OpCode: mmtls.MMLongOperationCheckQrcode,
		Data:   reqData,
		CgiUrl: "/cgi-bin/micromsg-bin/checkloginqrcode",
	}
	result, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}
	response, _ := result.(wechat.LoginQRCodeNotify)
	return &response, nil
}

// SendHeartBeatRequest 发送心跳包
func SendHeartBeatRequest(wxAccount wxface.IWXAccount, sync bool) (*wechat.HeartBeatResponse, error) {
	// 获取请求包
	tmpUserInfo := wxAccount.GetUserInfo()
	reqData, err := clientsdk.GetHeartBeatReq(tmpUserInfo)
	if err != nil {
		return nil, err
	}
	// 发送给长链接请求去处理
	longReq := &clientsdk.WXLongRequest{
		OpCode: mmtls.MMLongOperationHeartBeat,
		CgiUrl: "/cgi-bin/micromsg-bin/heartbeat",
		Data:   reqData,
	}
	result, err := wxlink.WXShortSend(wxAccount, longReq)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var response wechat.HeartBeatResponse
	if sync {
		response = result.(wechat.HeartBeatResponse)
	}
	return &response, nil
}

// SendGetProfileRequest 获取微信账号配置信息
func SendGetProfileRequest(wxAccount wxface.IWXAccount) error {
	// 发送请求
	tmpUserInfo := wxAccount.GetUserInfo()
	reqData := clientsdk.GetProfileReq(tmpUserInfo)

	// 发送给长链接请求去处理
	longReq := &clientsdk.WXLongRequest{
		OpCode: mmtls.MMLongOperationGetProfile,
		Data:   reqData,
		CgiUrl: "/cgi-bin/micromsg-bin/getprofile",
	}
	_, err := wxlink.WXSend(wxAccount, longReq, false)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// 发送登录请求根据设备Id进行登录
func SendManualAuthByDeviceIdRequest(wxAccount wxface.IWXAccount) (interface{}, error) {
	//检查是否在线会进行一次二次登录如果用户退出或者掉线才能重新登录
	//if CheckOnLineStatus(wxAccount) {
	//	return nil,errors.New("该链接已绑定微信号！")
	//}
	// 根据数据库缓存跟新UserInfo
	tmpUserInfo := wxAccount.GetUserInfo()
	tmpDeviceId := tmpUserInfo.LoginDataInfo.LoginData
	userName := tmpUserInfo.LoginDataInfo.UserName
	password := tmpUserInfo.LoginDataInfo.PassWord
	if strings.HasPrefix(tmpDeviceId, "62") {
		tmpDeviceId, err := clientsdk.Parse62Data(tmpDeviceId)
		if err != nil {
			return nil, err
		}
		tmpUserInfo.DeviceInfo.SetDeviceId(tmpDeviceId)
	} else if !strings.HasPrefix(tmpDeviceId, "49") && len(tmpDeviceId) >= 32 {
		tmpUserInfo.DeviceInfo.SetDeviceId("49" + tmpDeviceId[2:])
	} else if strings.HasPrefix(tmpDeviceId, "49") && len(tmpDeviceId) >= 32 {
		tmpUserInfo.DeviceInfo.SetDeviceId(tmpDeviceId)
	} else if strings.HasPrefix(tmpDeviceId, "A") {
		tmpUserInfo.DeviceInfoA16.DeviceId = []byte(tmpDeviceId[:15])
		tmpUserInfo.DeviceInfoA16.DeviceIdStr = tmpDeviceId
		tmpUserInfo.DeviceInfo = nil
		return dataA16Login(wxAccount, tmpUserInfo, password, userName)
	} else {
		//短信登录
		tmpUserInfo.Ticket = tmpUserInfo.LoginDataInfo.Ticket
		return SmsLogin(wxAccount, tmpUserInfo, password, userName)
		//return errors.New("SendManualAuthByDeviceIdRequest err: deviceId something the matter")
	}
	// 发送62登录请求
	tmpUserInfo.DeviceInfoA16 = nil
	return data62Login(wxAccount, tmpUserInfo, password, userName)
}

// A16登录请求
func dataA16Login(wxAccount wxface.IWXAccount, tmpUserInfo *baseinfo.UserInfo, password string, userName string) (interface{}, error) {
	//进行登录操作
	//packHeader, err := clientsdk.SendHybridManualAutoRequest(tmpUserInfo, password, userName, 146)
	packHeader, err := clientsdk.SendManualAuth(tmpUserInfo, password, userName)
	if err != nil {
		// 断开链接
		return nil, err
	}
	//var manualResponse wechat.UnifyAuthResponse
	var manualResponse wechat.ManualAuthResponse
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &manualResponse)
	if err != nil {
		log.Info("error SendManualAuth", err.Error())
		// 断开链接
		return nil, err
	}
	retCode := manualResponse.GetBaseResponse().GetRet()
	if retCode == baseinfo.MMErrIdcRedirect {
		wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
		dealRouter := wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
		return dealRouter.Handle(wxResp)
	} else if retCode == baseinfo.MMLoginSuccess {
		// 发送给微信消息处理器
		wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
		wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
		return nil, nil
	} else {
		// 断开链接
		//defer wxmgr.WxConnectMgr.Stop(wxAccount)

	}
	// 发送给微信消息处理器
	return nil, errors.New("login failed：" + manualResponse.GetBaseResponse().GetErrMsg().GetStr())
}

func SmsLogin(wxAccount wxface.IWXAccount, tmpUserInfo *baseinfo.UserInfo, password string, userName string) (interface{}, error) {
	var manualResponse wechat.ManualAuthResponse
	tmpUserInfo.DeviceInfoA16 = nil
	//packHeader, err := clientsdk.SendManualAuth(tmpUserInfo, password, userName)
	packHeader, err := clientsdk.SendHybridManualAutoRequest(tmpUserInfo, password, userName, 146)
	if err != nil {
		// 断开链接
		log.Info("短信登录断开", err.Error())
		return nil, err
	}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &manualResponse)
	if err != nil {
		// 断开链接
		log.Info("短信登录断开", err.Error())
		return nil, err
	}
	retCode := manualResponse.GetBaseResponse().GetRet()
	if retCode == baseinfo.MMErrIdcRedirect {
		wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
		//wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey()).SendToWXMsgHandler(wxResp)
		dealRouter := wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
		return dealRouter.Handle(wxResp)
	} else if retCode == baseinfo.MMLoginSuccess {
		// 发送给微信消息处理器
		wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
		wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
		return nil, nil
	} else {
		// 断开链接
		//defer wxmgr.WxConnectMgr.Stop(wxAccount)
	}
	// 发送给微信消息处理器
	return nil, errors.New("login failed：" + manualResponse.GetBaseResponse().GetErrMsg().GetStr())
}

// 62登录
func data62Login(wxAccount wxface.IWXAccount, tmpUserInfo *baseinfo.UserInfo, password string, userName string) (interface{}, error) {
	// 发送请求
	var manualResponse wechat.ManualAuthResponse
	//packHeader, err := clientsdk.SendHybridManualAutoRequest(tmpUserInfo, password, userName, 146)
	packHeader, err := clientsdk.SendManualAuth(tmpUserInfo, password, userName)
	//err = proto.Unmarshal(packHeader.Data, &manualResponse)
	if err != nil {
		// 断开链接
		log.Info("62登录断开", err.Error())
		//wxmgr.WxConnectMgr.Stop(wxAccount)
		return nil, err
	}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, &manualResponse)
	if err != nil {
		// 断开链接
		log.Info("62登录断开", err.Error())
		//wxmgr.WxConnectMgr.Stop(wxAccount)
		return nil, err
	}
	retCode := manualResponse.GetBaseResponse().GetRet()
	if retCode == baseinfo.MMErrIdcRedirect {
		wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
		dealRouter := wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
		return dealRouter.Handle(wxResp)
	} else if retCode == baseinfo.MMLoginSuccess {
		// 发送给微信消息处理器
		//go loginSuccess(wxqi,tmpUserInfo)
		wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
		wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
		return nil, nil
	} else {
		// 断开链接
		//defer wxmgr.WxConnectMgr.Stop(wxAccount)
	}
	// 发送给微信消息处理器
	return nil, errors.New("login failed：" + manualResponse.GetBaseResponse().GetErrMsg().GetStr())
}

//
////登录成功，做一些操作
//func loginSuccess(wxqi *WXReqInvoker, tmpUserInfo *baseinfo.UserInfo) {
//	if tmpUserInfo.DeviceInfo != nil {
//		resp, err := wxqi.SendReportstrategyRequest()
//		if err != nil {
//			log.Error("上报设备信息error", err.Error())
//		}
//		log.Println(resp)
//	}
//}

// SendLogoutRequest 退出登陆
func SendLogoutRequest(wxAccount wxface.IWXAccount) error {
	// 发送消息
	tmpUserInfo := wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendLogOutRequest(tmpUserInfo)
	if err != nil {
		//if packHeader != nil && packHeader.CheckSessionOut() {
		//	_,err = bizcgi.SendAutoAuthRequest(wxqi.wxconn)
		//	if err != nil {
		//		return err
		//	}
		//}
		// 断开链接, 发送token登陆
		return err
	}

	wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
	wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID).SendToWXMsgHandler(wxResp)
	return nil
}

// SendPushQrLoginNotice 发送二维码二次登录请求
func SendPushQrLoginNotice(wxAccount wxface.IWXAccount) (*wechat.PushLoginURLResponse, error) {
	// 重新链接，然后发送二维码二次登陆请求
	//wxmgr.WxConnectMgr.Stop(wxqi.wxAccount)
	//err := wxqi.wxconn.Start()
	//if err != nil {
	//	return nil, err
	//}

	// 发送请求
	tmpUserInfo := wxAccount.GetUserInfo()
	packHeader, err := clientsdk.SendPushQrLoginNotice(tmpUserInfo)
	if err != nil {
		if packHeader != nil && packHeader.GetRetCode() == baseinfo.MMRequestRetSessionTimeOut {
			return nil, errors.New("SendPushQrLoginNotice err: packHeader.RetCode == baseinfo.MMRequestRetSessionTimeOut ")
		}

		// 断开链接
		//wxmgr.WxConnectMgr.Stop(wxqi.wxAccount)
		return nil, err
	}

	// 获取登录二维码响应
	qrCodeResponse := &wechat.PushLoginURLResponse{}
	err = clientsdk.ParseResponseData(tmpUserInfo, packHeader, qrCodeResponse)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//wxmgr.WxConnectMgr.Stop(wxqi.wxAccount)
		return nil, err
	}

	//让路由处理时间
	wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
	dealRouter := wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
	dealRouter.Handle(wxResp)

	return qrCodeResponse, nil
}
