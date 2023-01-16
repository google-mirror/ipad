package wxmgr

import (
	"feiyu.com/wx/srv/wxface"
	"fmt"
	"github.com/gogf/gf/os/gcache"
	"github.com/lunny/log"
	"strconv"
	"time"

	"feiyu.com/wx/clientsdk/baseinfo"
)

var WxAccountMgr = NewWXAccountMgr()

// WXAccountMgr 微信链接管理器
type WXAccountMgr struct {
	wxAccountCache *gcache.Cache //管理的连接信息
}

// NewWXConnManager 创建一个WX链接管理
func NewWXAccountMgr() *WXAccountMgr {
	return &WXAccountMgr{
		wxAccountCache: gcache.New(1000),
	}
}

// Add 添加链接
func (wm *WXAccountMgr) Add(uuidKey string, accountAddr wxface.IWXAccount) {
	err := wm.wxAccountCache.Set(uuidKey, accountAddr, time.Minute*30)
	if err != nil {
		log.Error(err)
		return
	}
	// 打印链接数量
	wm.ShowConnectInfo()
}

// GetWXAccountByUserInfoUUID 根据UserInfoUUID获取微信链接
func (wm *WXAccountMgr) GetWXAccountByUserInfoUUID(userInfoUUID string) wxface.IWXAccount {
	wxAccount, _ := wm.wxAccountCache.Get(userInfoUUID)
	if wxAccount != nil {
		account, _ := (wxAccount).(wxface.IWXAccount)
		return account
	}
	log.Debug(fmt.Sprintf("GET wxAccount locfree Failed by %s  abandon the wxAccount get  !", userInfoUUID))
	//wxConn.GetWXAccount().SetLoginState(baseinfo.MMLoginStateLogout)
	return nil
}

// Remove 删除连接
func (wm *WXAccountMgr) Remove(uuidKey string) {
	//删除
	wm.wxAccountCache.Remove(uuidKey)
	// 打印链接数量
	wm.ShowConnectInfo()
}

// Len 获取当前连接
func (wm *WXAccountMgr) Len() int {
	len, _ := wm.wxAccountCache.Size()
	return len
}

// ClearWXConn 删除并停止所有链接
func (wm *WXAccountMgr) ClearWXConn() {
	wm.wxAccountCache.Clear()
	// 打印链接数量
	wm.ShowConnectInfo()
}

// ShowConnectInfo 打印链接情况
func (wm *WXAccountMgr) ShowConnectInfo() string {
	totalNum := uint32(0)
	noLoginNum := uint32(0)
	onlineNum := uint32(0)
	offlineNum := uint32(0)
	values, _ := wm.wxAccountCache.Values()
	for _, v := range values {
		totalNum++
		wxAccount := v.(wxface.IWXAccount)
		loginState := (wxAccount).GetLoginState()
		if loginState == baseinfo.MMLoginStateNoLogin {
			noLoginNum = noLoginNum + 1
		} else if loginState == baseinfo.MMLoginStateOnLine {
			onlineNum = onlineNum + 1
		} else if loginState == baseinfo.MMLoginStateOffLine {
			offlineNum = offlineNum + 1
		}
	}

	showText := time.Now().Format("2006-01-02 15:04:05")
	showText = showText + " 总链接数量: " + strconv.Itoa(int(totalNum)) + " 未登录数量:" + strconv.Itoa(int(noLoginNum))
	showText = showText + " 在线数量: " + strconv.Itoa(int(onlineNum)) + " 离线数量: " + strconv.Itoa(int(offlineNum))
	log.Println(showText)
	return showText
}
