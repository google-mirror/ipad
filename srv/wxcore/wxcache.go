package wxcore

import (
	"feiyu.com/wx/srv/wxface"
	"github.com/gogf/gf/os/gcache"
	"time"
)

var WxInfoCache = NewWXCache()

// WXCache 缓存
type WXCache struct {
	// 二维码缓存信息
	infoCache *gcache.Cache
}

// NewWXCache 新建一个缓存对象
func NewWXCache() wxface.IWXCache {
	return &WXCache{
		infoCache: gcache.New(),
	}
}

// SetQrcodeInfo 设置Qrcode信息
func (wxc *WXCache) SetQrcodeInfo(uuid string, qrAesKey []byte) {
	wxc.infoCache.Set(uuid, qrAesKey, time.Second*250)
}

// GetQrcodeInfo 获取二维码信息
func (wxc *WXCache) GetQrcodeInfo(uuid string) []byte {
	result, _ := wxc.infoCache.Get(uuid)
	var qrAesKey []byte
	if result != nil {
		qrAesKey, _ = result.([]byte)
	}
	return qrAesKey
}
