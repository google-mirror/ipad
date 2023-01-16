package wxlink

import (
	"errors"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/os/gtimer"
	"github.com/lunny/log"
	"strconv"
	"sync/atomic"
	"time"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/db"
	"feiyu.com/wx/srv/wxface"
)

// WXConnect 微信链接
type WXConnect struct {
	// 微信链接ID
	wxConnID uint32
	// 微信uuid key
	uuidKey string
	//// 同步请求管理器
	//wxSyncMgr wxface.IWXSyncMgr
	//用户消息管理器
	wxUserMsgMgr wxface.IWXUserMsgMgr
	// 心跳定时器
	heartBeatTimer *gtimer.Entry
	// 二次登录定时器
	//autoAuthTimer *gtimer.Entry
	// 是否连接着
	isConnected bool

	//重链数
	reset int64
}

// NewWXConnect 新的微信连接
func NewWXConnect(uuidKey string) wxface.IWXConnect {

	wxconn := &WXConnect{
		uuidKey:     uuidKey,
		isConnected: false,
		reset:       0,
	}
	//wxconn.wxSyncMgr = wxmgr.NewWXSyncMgr(wxconn)
	//用户消息管理器
	wxconn.wxUserMsgMgr = wxmgr.NewWXUSerMsgMgr(wxconn)
	return wxconn
}

func (wxconn *WXConnect) handleLongWriter(longReq wxface.IWXLongRequest) {
	wxAccount := wxconn.GetWXAccount()
	mmInfo := wxAccount.GetUserInfo().MMInfo
	seqId := longReq.GetSeqId()
	if seqId == 0 {
		seqId = atomic.LoadUint32(&mmInfo.LONGClientSeq)
	}
	err := mmtls.MMTCPSendReq(mmInfo, seqId, longReq.GetOpcode(), longReq.GetData())
	if err != nil {
		// 断开链接
		log.Printf("[%s],[%s] 断开链接 -  %s \n",
			wxconn.GetWXAccount().GetUserInfo().GetUserName(), wxconn.GetWXAccount().GetUserInfo().NickName, err.Error())
		wxconn.Stop()
	}
}

var lock = gmlock.New()

// 重新链接
func (wxconn *WXConnect) restartLong() {
	key := strconv.Itoa(int(wxconn.GetWXConnID()))
	if lock.TryLock(key) {
		defer lock.Unlock(key)
		wxAccount := wxconn.GetWXAccount()
		userInfo := wxAccount.GetUserInfo()
		username := userInfo.GetUserName()
		if username == "" {
			log.Println(key, "未登录成功无需重连")
			wxconn.Stop()
			return
		}
		// 断开链接
		flag := key + "->[" + username + "]->"
		//重试次数
		atomic.AddInt64(&wxconn.reset, 1)
		count := atomic.LoadInt64(&wxconn.reset)
		if count <= 3 {
			if wxAccount.GetUserInfo().DeviceInfo == nil {
				log.Println("a16重新开启长链接->", wxconn.isConnected)
			} else {
				log.Println("62重新开启长链接", wxconn.isConnected)
			}
			log.Println(flag, "重新开启长连接，重试次数=", count)
			if wxconn.isConnected {
				wxconn.isConnected = false
			}
			if userInfo.MMInfo.Conn != nil {
				userInfo.MMInfo.Conn.Close()
			}
			err := wxconn.startLongLink()
			if err != nil {
				log.Println(flag, "重新开启长连接失败", err.Error())
				if wxconn.reset < 3 {
					log.Println(flag, "等待", count, "s后重试")
					time.Sleep(time.Duration(count) * time.Second)
					lock.Unlock(key)
					wxconn.restartLong()
				} else {
					//发布长链失败消息
					messageResp := new(table.SyncMessageResponse)
					messageResp.ErrorMsg = "长连接3次重连失败"
					log.Println(flag, messageResp.ErrorMsg, err.Error())
					go db.PublishSyncMsgWxMessage(wxAccount.GetUserInfo(), *messageResp)
				}
				return
			} else {
				log.Println(flag, "重新开启长连接成功")
			}
		} else {
			log.Debug(flag, "自动重连次数已达到3次")
			wxconn.Stop()
		}
	}
}

// startLongWriter 开启长连接接受数据
func (wxconn *WXConnect) startLongReader() {
	wxAccount := wxconn.GetWXAccount()
	mmInfo := wxAccount.GetUserInfo().MMInfo
	for wxconn.IsConnected() {
		// 接收数据
		recvInfo, err := mmtls.MMTCPRecvData(mmInfo)
		if err != nil {
			log.Printf("[%s],[%s] 长连接出错 -  %s \n",
				wxconn.GetWXAccount().GetUserInfo().GetUserName(), wxAccount.GetUserInfo().NickName, err.Error())
			wxconn.restartLong()
			break
		}
		var packHeader *baseinfo.PackHeader
		//log.Println("---->",recvInfo.HeaderInfo.Operation)
		// 系统推送-需要同步消息 - 1000000000 1000000238
		if recvInfo.HeaderInfo.Operation < 1000000000 {
			//loginState := wxAccount.GetLoginState()
			//if loginState != baseinfo.MMLoginStateOnLine {
			//	continue
			//}

			// 同步收藏 todo 暂不需要
			if recvInfo.HeaderInfo.Operation == baseinfo.MMLongOperatorTypeFavSync {
				//packHeader := &baseinfo.PackHeader{}
				//packHeader.ReqData = recvInfo.RespData
				//packHeader.RetCode = 0
				//packHeader.URLID = baseinfo.MMLongOperatorTypeFavSync
				continue
			} else
			// TODO 发送同步请求
			if recvInfo.HeaderInfo.Operation == 24 {
				packHeader = &baseinfo.PackHeader{}
				packHeader.ReqData = recvInfo.RespData
				packHeader.RetCode = 0
				packHeader.URLID = 24
			} else {
				continue
			}
		} else {
			// 解包响应数据
			packHeader, err = clientsdk.DecodePackHeader(recvInfo.RespData, nil)
		}
		if packHeader == nil {
			log.Error("[", wxAccount.GetUserInfo().GetUserName(), "]",
				"startLongReader - DecodePackHeader packHeader is nil")
			wxconn.Stop()
			break
		}
		if err != nil {
			// TODO 接受消息出错，断开链接-token登陆
			log.Warn("[", wxAccount.GetUserInfo().GetUserName(), "]",
				"startLongReader - DecodePackHeader", packHeader.RetCode, err.Error())
			// 网络异常超时 不需要关闭
			if packHeader.RetCode != baseinfo.MMErrSessionTimeOut {
				wxconn.Stop()
				handleRecv(wxconn, recvInfo, packHeader)
				break
			}
		} else {
			//如果同步到消息，改为0
			atomic.AddInt64(&wxconn.reset, 0)
		}
		handleRecv(wxconn, recvInfo, packHeader)
	}
}

func handleRecv(wxconn *WXConnect, recvInfo *mmtls.LongRecvInfo, packHeader *baseinfo.PackHeader) {
	// 发送给微信消息处理器
	packHeader.SeqId = recvInfo.HeaderInfo.SequenceNumber
	wxResp := wxcore.NewWXResponse(wxconn.uuidKey, packHeader)
	wxconn.SendToWXMsgHandler(wxResp)
}

// startLong 开启长链接
func (wxconn *WXConnect) startLongLink() error {
	// 先进行长链接握手
	wxAccount := wxconn.GetWXAccount()
	userInfo := wxAccount.GetUserInfo()
	dialer := clientsdk.GetDialer(userInfo)
	tmpMMInfo, err := mmtls.InitMMTLSInfoLong(dialer, userInfo.LongHost, userInfo.LongPort, userInfo.ShortHost, nil)
	if err != nil {
		return err
	}

	userInfo.MMInfo = tmpMMInfo
	wxconn.isConnected = true
	// 启动长链接接收协程
	go wxconn.startLongReader()
	//需要一次立即的心跳 保证消息接收
	wxconn.SendHeartBeatRequest()
	//防止上面的出错 不继续心跳
	wxconn.SendHeartBeatWaitingSeconds(100)
	//wxconn.autoAuthTimer = time.NewTimer(time.Minute * 10)
	//wxconn.wxSyncMgr.Start()
	return nil
}

// Start 开启微信链接任务
func (wxconn *WXConnect) Start() error {
	// 如果是链接状态
	if wxconn.isConnected {
		return nil
	}
	wxAccount := wxconn.GetWXAccount()
	userInfo := wxAccount.GetUserInfo()
	// 判断微信信息是否为空
	if userInfo == nil {
		return errors.New("wxconn.Start() err: userInfo == nil")
	}
	atomic.AddInt64(&wxconn.reset, 0)
	// 开启长链接
	err := wxconn.startLongLink()
	if err != nil {
		return err
	}
	wxconn.wxUserMsgMgr.Start()
	wxmgr.WxConnectMgr.Add(wxconn)
	log.Printf("[%s] 开始长连接状态！\n", wxconn.GetWXUuidKey())
	// 设置下二次登录时间12
	//userInfo.UpdateLastAuthTime()
	return nil
}

// Stop 关闭链接
func (wxconn *WXConnect) Stop() {
	wxconn.StopWithReConnect(true)
}

func (wxconn *WXConnect) StopWithReConnect(isReConnect bool) {
	key := strconv.Itoa(int(wxconn.GetWXConnID()))
	wxAccount := wxconn.GetWXAccount()
	flag := key + "->[" + wxAccount.GetUserInfo().GetUserName() + "]->"
	if lock.TryLock(key) {
		defer lock.Unlock(key)
		if wxconn.isConnected == false {
			return
		}
		log.Debug(flag, "长链将被关闭")
		if log.OutputLevel() == 0 {
			glog.PrintStack(1)
		}
		if !isReConnect {
			wxconn.clearReset()
		}
		// 断开链接
		wxconn.isConnected = false
		log.Println("关闭链接Stop->", wxAccount.GetLoginState())
		// 关闭长链接
		userInfo := wxAccount.GetUserInfo()
		if userInfo.MMInfo.Conn != nil {
			userInfo.MMInfo.Conn.Close()
		}
		// 设置成离线状态
		if wxAccount.GetLoginState() == baseinfo.MMLoginStateOnLine {
			wxAccount.SetLoginState(baseinfo.MMLoginStateOffLine)
			db.UpdateUserInfo(userInfo)
			db.UpdateLoginStatus(userInfo.UUID, int32(userInfo.GetLoginState()), "已离线")
		}
		//wxconn.wxSyncMgr.Stop()
		wxconn.wxUserMsgMgr.Stop()
		wxmgr.WxConnectMgr.Remove(wxconn)
		// 立即过期
		if wxconn.heartBeatTimer != nil {
			wxconn.heartBeatTimer.Close()
		}
		//if wxconn.autoAuthTimer != nil {
		//	wxconn.autoAuthTimer.Close()
		//}
		log.Printf("[%s],[%s],[%s] 退出！\n", userInfo.GetUserName(), userInfo.NickName, userInfo.UUID)
	} else {
		log.Debug(flag, "长链关闭时未获取到锁")
	}
}

func (wxconn *WXConnect) clearReset() {
	atomic.AddInt64(&wxconn.reset, 3)
}

// SetWXConnID 设置微信链接ID
func (wxconn *WXConnect) SetWXConnID(wxConnID uint32) {
	wxconn.wxConnID = wxConnID
}

// GetWXConnID 获取WX链接ID
func (wxconn *WXConnect) GetWXConnID() uint32 {
	return wxconn.wxConnID
}

// GetWXUuidKey 获取WX uuidKey
func (wxconn *WXConnect) GetWXUuidKey() string {
	return wxconn.uuidKey
}

// GetWXAccount 获取微信帐号信息
func (wxconn *WXConnect) GetWXAccount() wxface.IWXAccount {
	return wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxconn.uuidKey)
}

// GetWXSyncMgr 获取同步管理器
//func (wxconn *WXConnect) GetWXSyncMgr() wxface.IWXSyncMgr {
//	return wxconn.wxSyncMgr
//}

// GetWXFriendMsgMgr 好友消息管理器
func (wxconn *WXConnect) GetWXFriendMsgMgr() wxface.IWXUserMsgMgr {
	return wxconn.wxUserMsgMgr
}

// IsConnected 判断是否链接状态
func (wxconn *WXConnect) IsConnected() bool {
	return wxconn.isConnected
}

// SendToWXMsgHandler 发送给消息队列去处理
func (wxconn *WXConnect) SendToWXMsgHandler(wxResp wxface.IWXResponse) {
	wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey()).GetWxServer().GetWXMsgHandler().SendWXRespToTaskQueue(wxResp)
}

// SendToWXLongReqQueue 添加到长链接请求队列
func (wxconn *WXConnect) SendToWXLongReqQueue(wxLongReq wxface.IWXLongRequest) {
	//wxconn.longReqQueue <- wxLongReq
	grpool.Add(func() {
		wxconn.handleLongWriter(wxLongReq)
	})
}
func (wxconn *WXConnect) SendHeartBeatRequest() {
	wxAccount := wxconn.GetWXAccount()
	if wxAccount.GetLoginState() == baseinfo.MMLoginStateOnLine {
		// 获取请求包
		tmpUserInfo := wxAccount.GetUserInfo()
		reqData, err := clientsdk.GetHeartBeatReq(tmpUserInfo)
		if err != nil {
			return
		}
		// 发送给长链接请求去处理
		longReq := &clientsdk.WXLongRequest{
			OpCode: mmtls.MMLongOperationHeartBeat,
			CgiUrl: "/cgi-bin/micromsg-bin/heartbeat",
			Data:   reqData,
		}
		_, err = WXSend(wxAccount, longReq, true)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}

// SendHeartBeatWaitingSeconds 添加到微信心跳包队列
func (wxconn *WXConnect) SendHeartBeatWaitingSeconds(seconds uint32) {
	//wxconn.heartBeatTimer.Reset(time.Second * time.Duration(seconds))
	if wxconn.heartBeatTimer != nil {
		wxconn.heartBeatTimer.Close()
	}
	if seconds == 0 {
		wxconn.heartBeatTimer = gtimer.AddOnce(time.Second*time.Duration(seconds), wxconn.SendHeartBeatRequest)
	} else {
		wxconn.heartBeatTimer = gtimer.AddSingleton(time.Second*time.Duration(seconds), wxconn.SendHeartBeatRequest)
	}
}

//func (wxconn *WXConnect) SendAutoAuth() {
//	wxAccount := wxconn.GetWXAccount()
//	if wxAccount.GetLoginState() == baseinfo.MMLoginStateOnLine {
//		_,_ = bizcgi.SendAutoAuthRequest(wxconn)
//	}
//}
//
//// SendAutoAuthWaitingMinutes 添加到微信二次包包队列
//func (wxconn *WXConnect) SendAutoAuthWaitingMinutes(minutes uint32) {
//	//wxconn.autoAuthTimer.Reset(time.Minute * time.Duration(minutes))
//	if wxconn.autoAuthTimer != nil {
//		wxconn.autoAuthTimer.Close()
//	}
//	wxconn.autoAuthTimer = gtimer.AddOnce(time.Minute * time.Duration(minutes), wxconn.SendAutoAuth)
//}
