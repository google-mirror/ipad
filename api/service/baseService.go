package service

import (
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/proxynet"
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/server"
	"feiyu.com/wx/srv/websrv"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"feiyu.com/wx/srv/wxmgr"
	"fmt"
	"github.com/gogf/guuid"
	"github.com/lunny/log"
	"reflect"
	"strings"
)

type BusinessFunc = func(account wxface.IWXAccount, newAccount bool) vo.DTO
type USerInfoCallFunc = func(account wxface.IWXAccount, invoker wxface.IWXReqInvoker, state uint32) vo.DTO

func CreateWXAccountByQueryKey(queryKey, proxy string, userInfo *baseinfo.UserInfo) *wxlink.WXAccount {
	var proxyInfo *proxynet.WXProxyInfo
	if proxy != "" {
		proxyInfo = proxynet.ParseWXProxyInfo(proxy)
	}
	wxAccount := wxlink.NewWXAccount(&websrv.TaskInfo{
		UUID: queryKey,
	}, proxyInfo, server.WxServer, userInfo)
	return wxAccount
}

// 检查实例Id是否存在 不存在创建新的账号
func checkExIdPerformIp(queryKey string, proxy string, businessFunc BusinessFunc) vo.DTO {
	//查询的queryKey为空创建一个链接实例
	if queryKey == "" {
		queryKey = guuid.New().String()
	}
	//查询该链接是否存在
	iwxAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(queryKey)
	if iwxAccount != nil {
		//执行回调方法
		return businessFunc(iwxAccount, false)
	} else {
		//如果链接管理器不存在该链接查询数据库是否存在
		dbUserInfo := db.GetUSerInfoByUUID(queryKey)
		//数据库存在该链接数据 重新实例化一个链接对象
		if dbUserInfo != nil {
			//创建一个用户信息
			iwxAccount = CreateWXAccountByQueryKey(queryKey, proxy, dbUserInfo)
			//设置用户信息
			//iwxAccount.SetUserInfo(dbUserInfo)
			return businessFunc(iwxAccount, false)
		} else {
			//创建新一个用户信息
			wxAccount := CreateWXAccountByQueryKey(queryKey, proxy, nil)
			return businessFunc(wxAccount, true)
		}
	}
}

// 检查实例Id是否存在 不存在创建新的链接
func checkExIdPerform(queryKey, deviceId string, businessFunc BusinessFunc) vo.DTO {
	//查询的queryKey为空创建一个链接实例
	if queryKey == "" {
		return businessFunc(nil, true)
	}
	//查询该链接是否存在
	iwxAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(queryKey)
	if iwxAccount != nil {
		//执行回调方法
		return businessFunc(iwxAccount, false)
	} else {
		//如果链接管理器不存在该链接查询数据库是否存在
		dbUserInfo := db.GetUSerInfoByUUID(queryKey)
		//数据库存在该链接数据 重新实例化一个链接对象
		if dbUserInfo != nil {
			//创建一个用户信息
			wxAccount := CreateWXAccountByQueryKey(queryKey, "", dbUserInfo)
			////创建一个新链接
			//wxConnect := wxlink.NewWXConnect(queryKey)
			//设置用户信息
			//wxAccount.SetUserInfo(dbUserInfo)
			return businessFunc(wxAccount, false)
		} else {
			//创建新一个用户信息
			wxAccount := wxlink.NewWXAccount(&websrv.TaskInfo{
				UUID:     queryKey,
				DeviceId: deviceId,
			}, nil, server.WxServer, nil)
			wxmgr.WxAccountMgr.Add(queryKey, wxAccount)
			////创建一个新链接
			//wxConnect := wxlink.NewWXConnect(server.WxServer)
			return businessFunc(wxAccount, true)
			/*return vo.DTO{
				Code: vo.FAIL_UUID,
				Data: nil,
				Text: fmt.Sprintf("%s 该链接不存在！",queryKey),
			}*/
		}
	}
}

// 检查实例Id是否存在 链接不存在返回错误不创建新链接
func checkExIdPerformNoCreateConnect(queryKey string, businessFunc BusinessFunc) vo.DTO {
	//查询的queryKey为空创建一个链接实例
	if queryKey == "" {
		return businessFunc(nil, true)
	}

	//查询该链接是否存在
	iwxAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(queryKey)
	if iwxAccount != nil {
		//执行回调方法
		return businessFunc(iwxAccount, false)
	} else {
		//如果链接管理器不存在该链接查询数据库是否存在
		dbUserInfo := db.GetUSerInfoByUUID(queryKey)
		//数据库存在该链接数据 重新实例化一个链接对象
		if dbUserInfo != nil {
			//dbUserInfo.CheckSumKey = []byte{}
			//dbUserInfo.SyncKey = []byte{}
			//创建一个用户信息
			wxAccount := CreateWXAccountByQueryKey(queryKey, "", dbUserInfo)
			////创建一个新链接
			//wxConnect := wxlink.NewWXConnect(queryKey)
			//设置用户信息
			//wxAccount.SetUserInfo(dbUserInfo)
			return businessFunc(wxAccount, false)
		} else {
			/*//创建新一个用户信息
			wxAccount := srv.NewWXAccount(&websrv.TaskInfo{
				UUID: queryKey,
			}, nil)
			//创建一个新链接
			wxConnect := wxcore.NewWXConnect(wxmgr.WxServer)
			return businessFunc(wxConnect,true)*/
			return vo.DTO{
				Code: vo.FAIL_UUID,
				Data: nil,
				Text: fmt.Sprintf("%s 该链接不存在！", queryKey),
			}
		}
	}
}

func DeepFields(ifaceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < ifaceType.NumField(); i++ {
		v := ifaceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

func StructCopy(DstStructPtr interface{}, SrcStructPtr interface{}) {
	srcv := reflect.ValueOf(SrcStructPtr)
	dstv := reflect.ValueOf(DstStructPtr)
	srct := reflect.TypeOf(SrcStructPtr)
	dstt := reflect.TypeOf(DstStructPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := DeepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}

/*
*
获取DeviceToken
*/
func checkDeviceToken(tmpUserInfo *baseinfo.UserInfo) {
	//如果是android
	if strings.HasPrefix(tmpUserInfo.LoginDataInfo.LoginData, "A") {
		key := fmt.Sprintf("%s%s", "wechat:deviceTokenA16:", tmpUserInfo.LoginDataInfo.UserName)
		exists, _ := db.Exists(key)
		if exists {
			//A16存redis
			trustRes := &wechat.TrustResp{}
			error := db.GETObj(key, &trustRes)
			if error != nil {
				log.Error("redis deviceToken is error=" + error.Error())
			}
			tmpUserInfo.DeviceInfoA16.DeviceToken = trustRes
		} else {
			tmpUserInfo.DeviceInfoA16.DeviceId = []byte(tmpUserInfo.LoginDataInfo.LoginData[:15])
			tmpUserInfo.DeviceInfoA16.DeviceIdStr = tmpUserInfo.LoginDataInfo.LoginData
			deviceTokenRsp, err := clientsdk.SendAndroIdDeviceTokenRequest(tmpUserInfo)
			if err != nil {
				log.Error("android 请求 deviceTokenRequest error!")
			}
			//保存5天
			db.SETExpirationObj(key, &deviceTokenRsp, 60*60*24*5)
			tmpUserInfo.DeviceInfoA16.DeviceToken = deviceTokenRsp
		}
	} else if strings.HasPrefix(tmpUserInfo.LoginDataInfo.LoginData, "62") {
		key := fmt.Sprintf("%s%s", "wechat:deviceTokenIos:", tmpUserInfo.LoginDataInfo.UserName)
		exists, _ := db.Exists(key)
		if exists {
			//ios存redis
			trustRes := &wechat.TrustResp{}
			error := db.GETObj(key, &trustRes)
			if error != nil {
				log.Error("ios redis deviceTokenIos is error=" + error.Error())
			}
			tmpUserInfo.DeviceInfo.DeviceToken = trustRes
		} else {
			deviceTokenRsp, err := clientsdk.SendIosDeviceTokenRequest(tmpUserInfo)
			if err != nil {
				log.Error("ios 请求 deviceTokenRequest error!")
				return
			}
			//保存5天
			db.SETExpirationObj(key, &deviceTokenRsp, 60*60*24*5)
			tmpUserInfo.DeviceInfo.DeviceToken = deviceTokenRsp
		}
	}
}
