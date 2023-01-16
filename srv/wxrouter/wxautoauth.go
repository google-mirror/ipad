package wxrouter

import (
	"errors"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"fmt"
	"github.com/lunny/log"
)

// WXAutoAuthRouter token登陆响应路由
type WXAutoAuthRouter struct {
	WXBaseRouter
}

// PreHandle 在处理conn业务之前的钩子方法
func (wxbr *WXAutoAuthRouter) PreHandle(response wxface.IWXResponse) error {
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(response.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	packHeader := response.GetPackHeader()
	//二次登录失败需要重新登录
	if packHeader.GetRetCode() != 0 {
		switch packHeader.RetCode {
		case baseinfo.MM_ERR_CERT_EXPIRED:
			// 切换密钥登录
			currentUserInfo.SwitchRSACert()
			_, _ = bizcgi.SendAutoAuthRequest(currentWXAccount)
			break
		case baseinfo.MMErrSessionTimeOut: // Session 会话过期
			currentWXAccount.SetLoginState(baseinfo.MMLoginStateNoLogin)
			log.Debug("WXAutoAuthRouter MMErrSessionTimeOut,err: " + currentWXAccount.GetUserInfo().GetUserName())
			db.UpdateLoginStatus(currentWXAccount.GetUserInfo().UUID, int32(currentWXAccount.GetLoginState()), "二次登录失败需要重新登录")
		default:
			//defer currentWXConn.Stop()
		}
		return errors.New("二次登录失败需要重新登录")
	}
	return nil
}

// Handle 处理conn业务的方法
func (glqr *WXAutoAuthRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	fmt.Printf("%+v\n", currentUserInfo, wxResp.GetWXUuidKey())
	//currentInvoker := currentWXConn.GetWXReqInvoker()
	packHeader := wxResp.GetPackHeader()

	// 解析token登陆响应
	var manualResponse wechat.ManualAuthResponse
	err := clientsdk.ParseResponseData(currentUserInfo, packHeader, &manualResponse)
	if err != nil {
		// 请求出问题了，应该关闭链接
		log.Debug("AutoAuth", currentWXAccount.GetUserInfo().GetUserName(), err.Error())
		//currentWXConn.Stop()
		return nil, err
	}

	retCode := manualResponse.GetBaseResponse().GetRet()
	// Mysql 提交登录日志
	db.SetLoginLog("AutoAuth", currentWXAccount.GetUserInfo(), manualResponse.GetBaseResponse().GetErrMsg().GetStr(), retCode)

	switch retCode {
	case baseinfo.MMOk: //success
		WXAutoAuthSuccess(&manualResponse, currentWXAccount)
		//redis 发布消息 发布登录状态
		//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateOnLine)
		//初始化
		log.Debug("AutoAuth", currentWXAccount.GetUserInfo().GetUserName(), "成功")
		return manualResponse, nil
	case baseinfo.MMErrDropped: //出现用户主动退出获取被T下线在线状态不存在需要调用Push
		// -2023 登录出现错误可重新登录
		/*if strings.Contains(manualResponse.GetBaseResponse().GetErrMsg().GetStr(), "登录出现错误") {
			currentWXConn.SendAutoAuthWaitingMinutes(5)
			return nil
		}*/
		//log.Println(hex.EncodeToString(currentUserInfo.AutoAuthKey))
		// 登录状态改为登录后退出
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateLogout)
		// 保存登录状态到数据库
		db.UpdateLoginStatus(currentUserInfo.UUID, int32(currentWXAccount.GetLoginState()), "你已退出微信")
		// 关闭重新启动，再次发送登陆请求
		//currentWXConn.Stop()
		//redis 发布消息 发布登录状态
		//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
		log.Debug("WXAutoAuthRouter retCode = - 2023,err: " + manualResponse.GetBaseResponse().GetErrMsg().GetStr())
		return nil, errors.New("WXAutoAuthRouter retCode = - 2023,err: 用户已主动退出")
	case -100:
		// 关闭重新启动，再次发送登陆请求
		//currentWXConn.Stop()
		// 登录状态改为未登录
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateNoLogin)
		errMsg := manualResponse.GetBaseResponse().GetErrMsg().GetStr()
		// 保存登录状态到数据库
		db.UpdateLoginStatus(currentUserInfo.UUID, int32(currentUserInfo.GetLoginState()), errMsg)
		//redis 发布消息 发布登录状态
		//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
		log.Debug("WXAutoAuthRouter retCode = - 100,err: " + manualResponse.GetBaseResponse().GetErrMsg().GetStr())
		return nil, errors.New("WXAutoAuthRouter retCode = - 100,err: " + manualResponse.GetBaseResponse().GetErrMsg().GetStr())
	case -6:
		// 登录状态改为未登录
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateNoLogin)
		// 关闭重新启动，再次发送登陆请求
		//currentWXConn.Stop()
		// 保存登录状态到数据库
		errMsg := manualResponse.GetBaseResponse().GetErrMsg().GetStr()
		db.UpdateLoginStatus(currentUserInfo.UUID, int32(currentUserInfo.GetLoginState()), errMsg)
		//redis 发布消息 发布登录状态
		//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
		log.Debug("WXAutoAuthRouter err: " + errMsg)
		return nil, errors.New("WXAutoAuthRouter err: " + errMsg)
	default:
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateNoLogin)
		// 关闭重新启动，再次发送登陆请求
		//currentWXConn.Stop()
		// 保存登录状态到数据库
		errMsg := manualResponse.GetBaseResponse().GetErrMsg().GetStr()
		db.UpdateLoginStatus(currentUserInfo.UUID, int32(currentUserInfo.GetLoginState()), errMsg)
		//redis 发布消息 发布登录状态
		//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
		log.Debug("WXAutoAuthRouter err: " + errMsg)
		return nil, errors.New("WXAutoAuthRouter err: " + errMsg)
	}

}

func WXAutoAuthSuccess(manualResponse *wechat.ManualAuthResponse, account wxface.IWXAccount) {
	userInfo := account.GetUserInfo()
	//reqInvoker := wxAccount.GetWXReqInvoker()
	// 获取aesKey
	userInfo.ConsultSessionKey(manualResponse.AuthParam.EcdhKey.Key.GetBuffer(), manualResponse.AuthParam.SessionKey.Key)
	// SetAutoKey
	userInfo.SetAutoKey(manualResponse.AuthParam.AutoAuthKey.Buffer)
	// SetNetworkSect
	userInfo.SetNetworkSect(manualResponse.DnsInfo)

	// 登录成功可与服务器重新建立长链接
	//if !connect.IsConnected() {
	//	err := connect.Start()
	//	if err != nil {
	//		return
	//	}
	//}

	// 是否需要初始化cdn信息
	if userInfo.CheckCdn() {
		// 获取CDNDns信息 todo
		//_,_ = reqInvoker.SendGetCDNDnsRequest()
	}
	// 获取账号的wxProfile
	_ = bizcgi.SendGetProfileRequest(account)
	// 登录状态改为在线
	account.SetLoginState(baseinfo.MMLoginStateOnLine)
	// 发送心跳
	//connect.SendHeartBeatWaitingSeconds(0)
	// 等待发送二次登录十分钟
	//connect.SendAutoAuthWaitingMinutes(60)

	// 保存UserInfo
	db.UpdateUserInfo(userInfo)
	// 更新状态到数据库
	db.UpdateLoginStatus(userInfo.UUID, int32(account.GetLoginState()), "登录成功！")
}
