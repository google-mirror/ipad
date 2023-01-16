package db

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"github.com/gogf/gf/os/gcache"
	"time"
)

var cache *gcache.Cache

func init() {
	cache = gcache.New()
}

// AddCheckStatusCache 添加扫码状态
func AddCheckStatusCache(queryKey string, val *baseinfo.CheckLoginQrCodeResult) {
	//_ = PublishSyncMsgCheckLogin(queryKey, val)

	isContains, _ := cache.Contains(queryKey)
	if !isContains {
		cache.Set(queryKey, val, time.Second*200)
		return
	} else {
		cache.Remove(queryKey)
	}
	if val.Ret == -3003 {
		cache.Set(queryKey, val, time.Second*1000)
		return
	}
	//扫码时间到
	if val.GetEffectiveTime() == 0 {
		//checkStatusCache.Remove(queryKey)
		return
	}
	cache.Set(queryKey, val, time.Second*time.Duration(val.GetEffectiveTime()))
}

func GetCheckStatusCache(queryKey string) *baseinfo.CheckLoginQrCodeResult {
	isContains, _ := cache.Contains(queryKey)
	if !isContains {
		return nil
	}
	obj, _ := cache.Get(queryKey)
	return obj.(*baseinfo.CheckLoginQrCodeResult)
}

// 获取
func GetCaChe() *gcache.Cache {
	return cache
}
