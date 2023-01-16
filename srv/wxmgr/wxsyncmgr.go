package wxmgr

//
//import (
//	"feiyu.com/wx/clientsdk/baseinfo"
//	"feiyu.com/wx/srv/bizcgi"
//	"feiyu.com/wx/srv/utils"
//	"feiyu.com/wx/srv/wxface"
//)
//
//// WXSyncMgr 同步管理器(同步消息，同步收藏等等)
//type WXSyncMgr struct {
//	wxConn        wxface.IWXConnect
//	newSyncIDList chan uint32
//	favSyncIDList chan uint32
//	newInitIDList chan uint32
//	endNewChan    chan bool
//	endFavChan    chan bool
//	endInitChan   chan bool
//}
//
//// NewWXSyncMgr 新建同步管理器
//func NewWXSyncMgr(wxConn wxface.IWXConnect) wxface.IWXSyncMgr {
//	return &WXSyncMgr{
//		wxConn:        wxConn,
//		newSyncIDList: make(chan uint32, 100),
//		favSyncIDList: make(chan uint32, 100),
//		newInitIDList: make(chan uint32, 200),
//		endNewChan:    make(chan bool, 1),
//		endFavChan:    make(chan bool, 1),
//		endInitChan:   make(chan bool, 1),
//	}
//}
//
//// Start 开启管理器
//func (wxsm *WXSyncMgr) Start() {
//	//go wxsm.startNewSyncListener()
//	//go wxsm.startFavSyncListener()
//	//go wxsm.startInitSyncListener()
//}
//
//// Stop 关闭管理器
//func (wxsm *WXSyncMgr) Stop() {
//	wxsm.endNewChan <- true
//	wxsm.endFavChan <- true
//	wxsm.endInitChan <- true
//}
//
//// SendNewSyncRequest 发送同步消息请求
//func (wxsm *WXSyncMgr) SendNewSyncRequest() {
//	//处理异常
//	defer utils.TryE(wxsm.wxConn.GetWXUuidKey())
//	// 同步消息
//
//	_, _ = bizcgi.SendNewSyncRequest(wxsm.wxConn.GetWXAccount(), baseinfo.MMSyncSceneTypeNeed, false)
//}
//
//// SendFavSyncRequest 发送同步收藏请求
//func (wxsm *WXSyncMgr) SendFavSyncRequest() {
//	wxsm.SendNewSyncRequest()
//}
//
//func (wxsm *WXSyncMgr) SendSyncInitRequest() {
//	_, _ = wxsm.wxConn.GetWXReqInvoker().SendNewInitSyncRequest()
//}
//
