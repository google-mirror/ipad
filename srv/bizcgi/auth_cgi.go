package bizcgi

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
	"time"
)

// GetA8KeyRequest 授权链接
func GetA8KeyRequest(wxAccount wxface.IWXAccount, opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*wechat.GetA8KeyResp, error) {
	userInfo := wxAccount.GetUserInfo()
	req := wechat.GetA8KeyRequest{
		BaseRequest: clientsdk.GetBaseRequest(userInfo),
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
	cgiCode := uint32(0)
	switch getType {
	case baseinfo.ThrIdGetA8Key:
		cgiCode = 388
		cgi = baseinfo.MMRequestTypeThrIdGetA8Key
		cgiUrl = "/cgi-bin/micromsg-bin/3rd-geta8key"
	default:
		cgiCode = 155
		cgi = baseinfo.MMRequestTypeGetA8Key
		cgiUrl = "/cgi-bin/micromsg-bin/geta8key"
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendData := clientsdk.Pack(userInfo, srcData, uint32(cgi), 5)
	longReq := &clientsdk.WXLongRequest{
		OpCode: cgiCode,
		CgiUrl: cgiUrl,
		Data:   sendData,
	}
	result, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	response, _ := result.(wechat.GetA8KeyResp)
	return &response, nil
}

// CheckOnLineStatus 检查在线状态
func CheckOnLineStatus(wxAccount wxface.IWXAccount) bool {
	state := wxAccount.GetLoginState()
	//先这样 后面考虑是否区分未登录和退出
	if state == baseinfo.MMLoginStateNoLogin ||
		state == baseinfo.MMLoginStateLogout ||
		//离线状态
		state == baseinfo.MMLoginStateOffLine {
		return false
	} else if state == baseinfo.MMLoginStateOnLine {
		//heartBeat,err := SendHeartBeatRequest(wxAccount, true)
		//if err == nil && heartBeat.GetBaseResponse().GetRet() == 0 {
		//	return true
		//}
		return true
	} else {
		//真正掉线
		return false
	}
}

// RecoverOnLineStatus 恢复在线状态
func RecoverOnLineStatus(wxAccount wxface.IWXAccount) bool {
	state := wxAccount.GetLoginState()
	// 先发送心跳 如果心跳发送成功则等待10分钟后再调用二次 ，心跳发送失败直接调用二次确认在线状态
	// 这里的心跳要求最好是同步的
	var err error
	if wxAccount.GetUserInfo().MMInfo != nil {
		_, err = SendHeartBeatRequest(wxAccount, true)
	}
	if err != nil || wxAccount.GetUserInfo().MMInfo == nil {
		_, err = SendAutoAuthRequest(wxAccount)
		if err != nil ||
			wxAccount.GetLoginState() == baseinfo.MMLoginStateNoLogin ||
			wxAccount.GetLoginState() == baseinfo.MMLoginStateLogout {
			return false
		}
	}

	if state == baseinfo.MMLoginStateOffLine {
		wxAccount.SetLoginState(baseinfo.MMLoginStateOnLine)
		db.UpdateLoginStatus(wxAccount.GetUserInfo().UUID, int32(wxAccount.GetLoginState()), "重新上线成功！")
	}
	wxconn := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
	if !wxconn.IsConnected() {
		err = wxconn.Start()
		if err != nil {
			log.Println(wxAccount.GetUserInfo().UUID, "启动长链接失败！", err.Error())
			return false
		}
		//wxconn.SendHeartBeatWaitingSeconds(1)
		//wxconn.SendAutoAuthWaitingMinutes(1)
		//wxconn.GetWXReqInvoker().SendAutoAuthRequest()
		_, _ = SendNewInitSyncRequest(wxAccount, false)
		//wxconn.GetWXSyncMgr().SendFavSyncRequest()
	}
	go func() {
		// 获取账号的wxProfile
		_ = SendGetProfileRequest(wxAccount)
	}()
	return true
}
