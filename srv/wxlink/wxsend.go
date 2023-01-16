package wxlink

import (
	"errors"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/srv/utils"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"fmt"
	"github.com/lunny/log"
	"sync/atomic"
)

// WXSyncSend 发送同步数据
func WXSyncSend(wxAccount wxface.IWXAccount, req wxface.IWXLongRequest) (interface{}, error) {
	return WXSend(wxAccount, req, true)
}

// WXSend 发送数据
func WXSend(wxAccount wxface.IWXAccount, req wxface.IWXLongRequest, sync bool) (interface{}, error) {
	defer utils.TryE(wxAccount.GetUserInfo().GetUserName())
	var result interface{}
	var err error
	userinfo := wxAccount.GetUserInfo()
	wxconn := wxmgr.WxConnectMgr.GetWXConnectByUserInfoUUID(wxAccount.GetUserInfo().UUID)
	if wxconn != nil && wxconn.IsConnected() && req.GetOpcode() > 0 {
		req.SetSeqId(atomic.LoadUint32(&userinfo.GetMMInfo().LONGClientSeq))
		result, err = WXShortSend(wxAccount, req)
		////发送给长链接请求去处理
		//key := wxconn.GetWXUuidKey() + "-" + strconv.Itoa(int(req.GetSeqId()))
		//event := common.NewEvent(key)
		//common.AddEvent(key, event)
		//defer common.RemoveEvent(key)
		//waiter := event.AddWaiter()
		//log.Debugf("进行长连接请求：%s,%d", key, req.GetOpcode())
		//wxconn.SendToWXLongReqQueue(req)
		//if sync {
		//	var pack *common.Pack
		//	pack, err = event.Wait(waiter, time.Second*3)
		//	if err != nil {
		//		return nil, err
		//	}
		//	result, err = pack.Content, pack.Err
		//}
		//if err != nil {
		//	return nil, err
		//}
	} else {
		result, err = WXShortSend(wxAccount, req)
	}
	return result, err
}

// WXShortSend 发送数据
func WXShortSend(wxAccount wxface.IWXAccount, req wxface.IWXLongRequest) (interface{}, error) {
	var result interface{}
	var err error
	userinfo := wxAccount.GetUserInfo()
	log.Debug("进行短连接请求：", req.GetCgiUrl())
	fmt.Printf("%+v\n", userinfo)
	resp, err := mmtls.MMHTTPPostData(userinfo.GetMMInfo(), req.GetCgiUrl(), req.GetData())
	if err != nil {
		return nil, err
	}
	packHeader, err := clientsdk.DecodePackHeader(resp, nil)
	if err != nil {
		//if packHeader != nil && (packHeader.RetCode == baseinfo.MMRequestRetSessionTimeOut) {
		//	// token登陆
		//	wxconn := wxmgr.WxConnectMgr.GetWXConnectByWXID(wxAccount.GetUserInfo().UUID)
		//	wxconn.SendAutoAuthWaitingMinutes(0)
		//}
		return nil, err
	}
	// 发送给微信消息处理器
	handler := wxAccount.GetWxServer().GetWXMsgHandler().GetRouterByRespID(packHeader.URLID)
	if handler == nil {
		return nil, errors.New("未找到结果处理handler")
	}
	wxResp := wxcore.NewWXResponse(wxAccount.GetUserInfo().UUID, packHeader)
	//执行对应处理方法
	err = handler.PreHandle(wxResp)
	if err == nil {
		result, err = handler.Handle(wxResp)
		if err == nil {
			err = handler.PostHandle(wxResp)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, err
}
