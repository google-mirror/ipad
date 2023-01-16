package wxmgr

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/lunny/log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"feiyu.com/wx/clientsdk/baseinfo"

	"feiyu.com/wx/srv/wxface"
)

var WxConnectMgr = NewWXConnManager()

// WXConnectMgr 微信链接管理器
type WXConnectMgr struct {
	canUseConnIDList []uint32 // 删掉/回收后的connID
	currentWxConnID  uint32
	wxConnectMap     *gmap.Map //管理的连接信息
}

// NewWXConnManager 创建一个WX链接管理
func NewWXConnManager() wxface.IWXConnectMgr {
	return &WXConnectMgr{
		canUseConnIDList: make([]uint32, 0),
		currentWxConnID:  0,
		wxConnectMap:     gmap.New(true),
	}
}

// Add 添加链接
func (wm *WXConnectMgr) Add(wxConnect wxface.IWXConnect) {
	newConnID := atomic.AddUint32(&wm.currentWxConnID, 1)
	wxConnect.SetWXConnID(newConnID)
	wxAccount := WxAccountMgr.GetWXAccountByUserInfoUUID(wxConnect.GetWXUuidKey())
	wm.wxConnectMap.Set(wxAccount.GetUserInfo().UUID, wxConnect)
	// 打印链接数量
	wm.ShowConnectInfo()
}

// GetWXConnectByUserInfoUUID 根据UserInfoUUID获取微信链接
func (wm *WXConnectMgr) GetWXConnectByUserInfoUUID(userInfoUUID string) wxface.IWXConnect {
	wxConn := wm.wxConnectMap.Get(userInfoUUID)
	if wxConn != nil {
		return wxConn.(wxface.IWXConnect)
	}
	log.Debug(fmt.Sprintf("GET Connection locfree Failed by %s  abandon the conntection get  !", userInfoUUID))
	//wxConn.GetWXAccount().SetLoginState(baseinfo.MMLoginStateLogout)
	return nil
}

// GetWXConnectByWXID 根据WXID获取微信链接
func (wm *WXConnectMgr) GetWXConnectByWXID(wxid string) wxface.IWXConnect {
	var tryCoon wxface.IWXConnect
	f := func(k, v interface{}) bool {
		tmp := v.(wxface.IWXConnect)
		if tmp != nil {
			wxAccount := WxAccountMgr.GetWXAccountByUserInfoUUID(tmp.GetWXUuidKey())
			tmpUserInfo := wxAccount.GetUserInfo()
			if tmpUserInfo != nil && strings.Compare(tmpUserInfo.WxId, wxid) == 0 {
				tryCoon = v.(wxface.IWXConnect)
				return false
			}
		}
		return true
	}
	wm.wxConnectMap.Iterator(f)
	return tryCoon
}

// Start 开始连接
func (wm *WXConnectMgr) Start(wxAccount wxface.IWXAccount) {
	//删除
	currentUserInfo := wxAccount.GetUserInfo()
	wxConn := wm.GetWXConnectByUserInfoUUID(currentUserInfo.UUID)
	if wxConn != nil {
		wxConn.Start()
	}
}

// Stop 停止连接
func (wm *WXConnectMgr) Stop(wxAccount wxface.IWXAccount) {
	//删除
	currentUserInfo := wxAccount.GetUserInfo()
	wxConn := wm.GetWXConnectByUserInfoUUID(currentUserInfo.UUID)
	if wxConn != nil {
		wxConn.Stop()
	}
}

// Remove 删除连接
func (wm *WXConnectMgr) Remove(wxconn wxface.IWXConnect) {
	//删除
	wxAccount := WxAccountMgr.GetWXAccountByUserInfoUUID(wxconn.GetWXUuidKey())
	currentUserInfo := wxAccount.GetUserInfo()
	wm.wxConnectMap.Remove(currentUserInfo.UUID)
	wm.Stop(wxAccount)
	//wm.canUseConnIDList = append(wm.canUseConnIDList, wxconn.GetWXConnID())
	currentUserInfo = nil
	// 打印链接数量
	wm.ShowConnectInfo()
}

// Len 获取当前连接
func (wm *WXConnectMgr) Len() int {
	return wm.wxConnectMap.Size()
}

// ClearWXConn 删除并停止所有链接
func (wm *WXConnectMgr) ClearWXConn() {
	wm.wxConnectMap.Clear()
	// 打印链接数量
	wm.ShowConnectInfo()
}

// ShowConnectInfo 打印链接情况
func (wm *WXConnectMgr) ShowConnectInfo() string {
	totalNum := uint32(0)
	noLoginNum := uint32(0)
	onlineNum := uint32(0)
	offlineNum := uint32(0)
	f := func(k, v interface{}) bool {
		totalNum++
		wxConn := v.(wxface.IWXConnect)
		wxAccount := WxAccountMgr.GetWXAccountByUserInfoUUID(wxConn.GetWXUuidKey())
		loginState := wxAccount.GetLoginState()
		if loginState == baseinfo.MMLoginStateNoLogin {
			noLoginNum = noLoginNum + 1
		} else if loginState == baseinfo.MMLoginStateOnLine {
			onlineNum = onlineNum + 1
		} else if loginState == baseinfo.MMLoginStateOffLine {
			offlineNum = offlineNum + 1
		}
		return true
	}
	wm.wxConnectMap.Iterator(f)
	showText := time.Now().Format("2006-01-02 15:04:05")
	showText = showText + " 总链接数量: " + strconv.Itoa(int(totalNum)) + " 未登录数量:" + strconv.Itoa(int(noLoginNum))
	showText = showText + " 在线数量: " + strconv.Itoa(int(onlineNum)) + " 离线数量: " + strconv.Itoa(int(offlineNum))
	log.Println(showText)
	return showText
}
