package wxrouter

import (
	"errors"
	"feiyu.com/wx/srv/bizcgi"
	"fmt"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/gogo/protobuf/proto"
	"github.com/lunny/log"
)

// WXCheckQrcodeRouter 检测二维码状态响应路由
type WXCheckQrcodeRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (cqr *WXCheckQrcodeRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {

	//currentWXConn := wxResp.GetWXConncet()
	currentAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentAccount.GetUserInfo()

	// 解析检测二维码响应
	var checkResp wechat.CheckLoginQRCodeResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &checkResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.StopWithReConnect(false)
		return nil, err
	}
	// 解密出现问题，说明协议出现了问题
	qrAesKey := wxcore.WxInfoCache.GetQrcodeInfo(currentUserInfo.QrUuid)
	if checkResp.GetBaseResponse().GetRet() == 0 && qrAesKey != nil {
		retBytes, err := baseutils.AesDecryptByteKey(checkResp.LoginQrcodeNotifyPkg.NotifyData.Buffer, qrAesKey)
		if err != nil {
			// 请求出问题了，应该关闭链接
			//currentWXConn.Stop()
			return nil, err
		}

		lgQrNotify := &wechat.LoginQRCodeNotify{}
		err = proto.Unmarshal(retBytes, lgQrNotify)
		if err != nil {
			// 请求出问题了，应该关闭链接
			//currentWXConn.StopWithReConnect(false)
			return nil, err
		}
		log.Info("qrCheckStatus", lgQrNotify.GetUuid(), lgQrNotify.GetState(), lgQrNotify.GetEffectiveTime())
		currentUserInfo.HeadURL = lgQrNotify.GetHeadImgUrl()
		if lgQrNotify.GetState() == 2 {
			// 确认是否可以登录
			/*if checkIsLogin(currentUserInfo, lgQrNotify) {
				currentWXConn.Stop()
				log.Printf("[%s] 已经登录！无需再次登录[%s]。\n", currentUserInfo.WxId, currentUserInfo.UUID)
				return errors.New(currentUserInfo.WxId + "已经登录！无需再次登录。")
			}*/
			// 扫码成功发送登录包
			currentUserInfo.LoginDataInfo.NewPassWord = lgQrNotify.GetWxnewpass()
			currentUserInfo.LoginDataInfo.UserName = lgQrNotify.GetWxid()

			t, err := clientsdk.SendIosDeviceTokenRequest(currentAccount.GetUserInfo())
			if err != nil {
				fmt.Println("get token err due to", err.Error())
				return nil, err
			}
			currentAccount.GetUserInfo().DeviceInfo.DeviceToken = t

			fmt.Println("get token :%v", t)

			//_, err = clientsdk.SendHybridManualAutoRequest(currentAccount.GetUserInfo(), currentUserInfo.LoginDataInfo.NewPassWord, currentUserInfo.WxId, 0)
			_, err = bizcgi.SendManualAuth(currentAccount, lgQrNotify.GetWxnewpass(), lgQrNotify.GetWxid())
			if err != nil {
				log.Debug(lgQrNotify.GetWxid(), "扫码成功尝试登录报错", err.Error())
				return nil, err
			}
		} else if lgQrNotify.GetState() == 4 {
			return nil, errors.New("WXCheckQrcodeRouter err: 二维码失效")
		}
		return *lgQrNotify, nil
	}
	return nil, errors.New("二维码失效")
}

// checkIsLogin 检测是否登录
func checkIsLogin(userInfo *baseinfo.UserInfo, loginNotify *wechat.LoginQRCodeNotify) bool {
	entity := db.GetUserInfoEntity(userInfo.UUID)
	if entity != nil &&
		(entity.State == int32(baseinfo.MMLoginStateOnLine)) {
		// 设置下返回已登录信息
		userInfo.WxId = entity.WxId
		//db.AddCheckStatusCache(userInfo.UUID, &baseinfo.CheckLoginQrCodeResult{
		//	LoginQRCodeNotify:   loginNotify,
		//	OthersInServerLogin: true,
		//	TargetServer:        entity.TargetIp,
		//	UUId:                entity.UUID,
		//})
		return true
	}
	// 可以登录
	userInfo.WxId = loginNotify.GetWxid()
	userInfo.NickName = loginNotify.GetNickName()
	return false
}
