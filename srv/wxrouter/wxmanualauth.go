package wxrouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXManualAuthRouter 扫码登陆响应路由
type WXManualAuthRouter struct {
	WXBaseRouter
}

// PreHandle 在处理conn业务之前的钩子方法
func (wxbr *WXManualAuthRouter) PreHandle(response wxface.IWXResponse) error {
	currentAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(response.GetWXUuidKey())
	tmpUserInfo := currentAccount.GetUserInfo()
	packHeader := response.GetPackHeader()
	if packHeader.RetCode == baseinfo.MM_ERR_CERT_EXPIRED {
		// 切换密钥登录
		tmpUserInfo.SwitchRSACert()
		return errors.New("请重试")
	}
	if packHeader.Data == nil {
		return errors.New("返回内容为空，无法反序列化")
	}
	return nil
}

// Handle 处理conn业务的方法
func (glqr *WXManualAuthRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	//currentCache := currentWXConn.GetWXCache()
	//currentInvoker := currentWXConn.GetWXReqInvoker()
	currentUserInfo := currentWXAccount.GetUserInfo()
	currentPackHeader := wxResp.GetPackHeader()
	//currentWXFileHelperMgr := currentWXConn.GetWXFileHelperMgr()

	// 解析扫码登陆响应
	var manualResponse wechat.ManualAuthResponse
	err := clientsdk.ParseResponseData(currentUserInfo, currentPackHeader, &manualResponse)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	retCode := manualResponse.GetBaseResponse().GetRet()
	// 重定向
	if retCode == baseinfo.MMErrIdcRedirect {
		resp := manualResponse.GetDnsInfo()
		for _, ip := range resp.GetBuiltinIplist().GetShortConnectIplist() {
			ipStr := strings.TrimRight(ip.GetIp(), "\u0000")
			if ipStr == "127.0.0.1" {
				continue
			}
			currentUserInfo.ShortHost = ipStr
			break
		}

		for _, ip := range resp.GetBuiltinIplist().GetLongConnectIplist() {
			ipStr := strings.TrimRight(ip.GetIp(), "\u0000")
			if ipStr == "127.0.0.1" {
				continue
			}
			currentUserInfo.LongHost = strings.TrimRight(ip.GetDomain(), "\u0000")
			currentUserInfo.LongPort = strconv.Itoa(int(ip.GetPort()))
			break
		}
		currentUserInfo.MMInfo = nil
		currentUserInfo.GetMMInfo().ShortHost = currentUserInfo.ShortHost
		currentUserInfo.GetMMInfo().LongHost = currentUserInfo.LongHost
		currentUserInfo.GetMMInfo().LONGPort = currentUserInfo.LongPort
		//提交登录日志
		db.SetLoginLog("ManualAuth", currentWXAccount.GetUserInfo(), fmt.Sprintf("重定向登录 ShortHost :%s,LongHost:%s", currentUserInfo.ShortHost, currentUserInfo.LongHost), retCode)

		// 关闭重新启动，再次发送登陆请求
		//currentWXConn.Stop()
		//_ = currentWXConn.Start()
		//todo 要补充deviceId登录时清空伪密码的逻辑
		if currentUserInfo.LoginDataInfo.NewPassWord == "" {
			return bizcgi.SendManualAuthByDeviceIdRequest(currentWXAccount)
		} else {
			//return clientsdk.SendHybridManualAutoRequest(currentWXAccount.GetUserInfo(), currentUserInfo.LoginDataInfo.NewPassWord, currentUserInfo.WxId, 0)
			return bizcgi.SendManualAuth(currentWXAccount, currentUserInfo.LoginDataInfo.NewPassWord, currentUserInfo.WxId)
		}
	} else if retCode == baseinfo.MMOk {

		//取基本信息
		accountInfo := manualResponse.GetAccountInfo()
		currentUserInfo.SetWxId(accountInfo.GetWxid())
		currentUserInfo.NickName = accountInfo.GetNickName()

		// 协商密钥
		currentUserInfo.ConsultSessionKey(manualResponse.GetAuthParam().GetEcdhKey().GetKey().GetBuffer(), manualResponse.AuthParam.SessionKey.Key)
		// AutoAuthKey
		currentUserInfo.SetAutoKey(manualResponse.AuthParam.AutoAuthKey.Buffer)
		// 随机一个字符红包AesKey
		currentUserInfo.GenHBKey()

		// 如果数据库有存储这个微信号的信息
		if len(currentUserInfo.FavSyncKey) <= 0 {
			// 刷新收藏同步Key
			oldUserInfo := db.GetUSerInfoByUUID(currentUserInfo.UUID)
			if oldUserInfo != nil {
				currentUserInfo.FavSyncKey = oldUserInfo.FavSyncKey
			}
		}
		x, _ := json.MarshalIndent(currentUserInfo, "", "\t")
		ioutil.WriteFile(currentUserInfo.NickName+".json", x, 0777)
		// 开始发送心跳包
		//currentWXConn.SendHeartBeatWaitingSeconds(10)
		// 设置登陆状态，发送提示
		currentWXAccount.SetLoginState(baseinfo.MMLoginStateOnLine)
		go func() {
			// 保存UserInfo
			db.SaveUserInfo(currentUserInfo)
			//Mysql  保存登录状态
			db.UpdateLoginStatus(currentUserInfo.UUID, int32(currentWXAccount.GetLoginState()), "登录成功！")
			// 获取账号的wxProfile
			bizcgi.SendGetProfileRequest(currentWXAccount)
			// 获取联系人标签列表
			//currentInvoker.SendGetContactLabelListRequest(false)
			// 获取CDNDns信息
			//todo
			//currentWXAccount.GetWxConnect().GetWXReqInvoker().SendGetCDNDnsRequest()
			// 开始发送二次登录包
			//currentWXConn.SendAutoAuthWaitingMinutes(60)
			// 初始化通讯录
			//contactSeq := currentCache.GetContactSeq()
			//currentInvoker.SendInitContactRequest(contactSeq)
			//currentCache.SetInitContactFinished(true)
			// 同步消息
			_, _ = bizcgi.SendNewSyncRequest(currentWXAccount, baseinfo.MMSyncSceneTypeNeed, false)
			// 打印当前链接状态
			wxmgr.WxConnectMgr.ShowConnectInfo()
			//currentWXFileHelperMgr.AddNewTipMsg("上线成功！")
			//currentWXFileHelperMgr.AddNewTipMsg("系统正在初始化...")
			//redis 发布消息 发布登录状态
			//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
			//Mysql 提交登录日志
			db.SetLoginLog("ManualAuth", currentWXAccount.GetUserInfo(), "登录成功！", retCode)

			/*time.Sleep(time.Second * 10)
			currentWXConn.Stop()*/
		}()
	} else {
		// 登陆失败
		errMsg := manualResponse.GetBaseResponse().GetErrMsg().GetStr()
		baseutils.PrintLog(errMsg)
		//Mysql 提交登录日志
		//db.SetLoginLog("ManualAuth", currentWXAccount, errMsg, retCode)
		//Mysql  保存用户信息
		//db.SaveUserInfo(currentUserInfo)
		//Mysql  保存登录状态
		//db.UpdateLoginStatus(currentUserInfo.UUID, retCode, errMsg)
		//currentWXConn.Stop()

		//redis 发布消息 发布登录状态
		//db.PublishLoginState(currentWXAccount.GetUserInfo().UUID, currentWXAccount.GetLoginState())
		return nil, errors.New("login failed：" + errMsg)
	}
	return manualResponse, nil
}
