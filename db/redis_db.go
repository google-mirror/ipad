package db

import (
	"encoding/base64"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lunny/log"
	"strings"
)

type RedisDb struct {
}

// InitDB 初始化数据库
func (db *RedisDb) InitDB() {

}

// 提交登录日志
func (db *RedisDb) SetLoginLog(loginType string, userInfo *baseinfo.UserInfo, errMsg string, state int32) {
	var userName string
	if len(userInfo.LoginDataInfo.UserName) > 0 {
		userName = userInfo.LoginDataInfo.UserName
	} else {
		userName = userInfo.WxId
	}
	loginLog := &table.UserLoginLog{
		UUId:      userInfo.UUID,
		UserName:  userName,
		NickName:  userInfo.NickName,
		LoginType: loginType,
		RetCode:   state,
		ErrMsg:    errMsg,
	}
	loginLog.TargetIp = srvconfig.GlobalSetting.TargetIp
	LPUSH(userName+"-loginLog", loginLog)
}

// 获取登录日志
func (db *RedisDb) GetLoginJournal(userName string) []table.UserLoginLog {
	return nil
}

// 保存登录状态
func (db *RedisDb) UpdateLoginStatus(uuid string, state int32, errMsg string) {
	data := make(map[string]interface{})
	data["State"] = state //零值字段
	data["ErrMsg"] = errMsg
	v := db.GetUserInfoEntity(uuid) // *table.UserInfoEntity
	if v != nil && v.State != state {
		v.State = state
		v.ErrMsg = errMsg
		SETObj(v.UUID, v)
		_ = PublishSyncMsgLoginState(v.WxId, uint32(state))
	}
}

// 更新用户信息
func (db *RedisDb) UpdateUserInfo(userInfo *baseinfo.UserInfo) {
	if userInfo.WxId == "" {
		userInfo.WxId = userInfo.GetUserName()
	}
	var userInfoEntity table.UserInfoEntity
	SimpleCopyProperties(&userInfoEntity, userInfo)
	userInfoEntity.AutoAuthKey = base64.StdEncoding.EncodeToString(userInfo.AutoAuthKey)
	userInfoEntity.SyncKey = base64.StdEncoding.EncodeToString(userInfo.SyncKey)
	userInfoEntity.FavSyncKey = base64.StdEncoding.EncodeToString(userInfo.FavSyncKey)
	userInfoEntity.TargetIp = srvconfig.GlobalSetting.TargetIp
	userInfoEntity.UserName = userInfo.LoginDataInfo.UserName
	userInfoEntity.Password = userInfo.LoginDataInfo.PassWord
	userInfoEntity.State = int32(userInfo.GetLoginState())
	SETObj(userInfo.UUID, userInfoEntity)
}

// SaveUserInfo 保存用户信息
func (db *RedisDb) SaveUserInfo(userInfo *baseinfo.UserInfo) {

	if userInfo.WxId == "" {
		userInfo.WxId = userInfo.GetUserName()
	}
	var userInfoEntity table.UserInfoEntity
	SimpleCopyProperties(&userInfoEntity, userInfo)
	userInfoEntity.AutoAuthKey = base64.StdEncoding.EncodeToString(userInfo.AutoAuthKey)
	userInfoEntity.SyncKey = base64.StdEncoding.EncodeToString(userInfo.SyncKey)
	userInfoEntity.FavSyncKey = base64.StdEncoding.EncodeToString(userInfo.FavSyncKey)
	userInfoEntity.TargetIp = srvconfig.GlobalSetting.TargetIp
	userInfoEntity.UserName = userInfo.LoginDataInfo.UserName
	userInfoEntity.Password = userInfo.LoginDataInfo.PassWord

	SETObj(userInfo.UUID, userInfoEntity)
	//判断是62还是A16
	if strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") {
		//A16存redis
		key := fmt.Sprintf("%s%s", "wechat:a16DeviceInfo:", userInfo.WxId)

		error := SETExpirationObj(key, &userInfo.DeviceInfoA16, 60*60*24*6)
		if error != nil {
			log.Error("保存redis is error=" + error.Error())
		}
	} else {
		//62存DB
		var deviceInfoEntity table.DeviceInfoEntity
		deviceInfoEntity.WxId = userInfo.WxId
		SimpleCopyProperties(&deviceInfoEntity, userInfo.DeviceInfo)
		deviceInfoEntity.SoftTypeXML = base64.StdEncoding.EncodeToString([]byte(deviceInfoEntity.SoftTypeXML))
		deviceInfoEntity.ClientCheckDataXML = base64.StdEncoding.EncodeToString([]byte(deviceInfoEntity.ClientCheckDataXML))
		//62存redis
		key := fmt.Sprintf("%s%s", "wechat:62DeviceInfo:", userInfo.WxId)
		SETObj(key, deviceInfoEntity)
	}
}

// GetUserInfoEntity 获取登录信息
func (db *RedisDb) GetUserInfoEntity(uuid string) *table.UserInfoEntity {
	var userInfoEntity table.UserInfoEntity
	err := GETObj(uuid, &userInfoEntity)
	if err != nil {
		return nil
	}
	return &userInfoEntity
}

func (db *RedisDb) GetUSerInfoByUUID(uuid string) *baseinfo.UserInfo {
	userInfo := &baseinfo.UserInfo{}

	var userInfoEntity table.UserInfoEntity
	userInfoEntity.UUID = uuid

	err := GETObj(uuid, &userInfoEntity)
	if err != nil {
		return nil
	}
	//判断是62还是A16
	key := fmt.Sprintf("%s%s", "wechat:a16DeviceInfo:", userInfoEntity.WxId)
	exists, _ := Exists(key)
	if exists {
		//A16存redis
		deviceInfoA16 := &baseinfo.AndroidDeviceInfo{}
		error := GETObj(key, &deviceInfoA16)
		if error != nil {
			log.Error("保存redis is error=" + error.Error())
		}
		userInfo.DeviceInfoA16 = deviceInfoA16
	} else {
		var deviceInfoEntity table.DeviceInfoEntity
		deviceInfoEntity.WxId = userInfoEntity.WxId
		key := fmt.Sprintf("%s%s", "wechat:62DeviceInfo:", userInfoEntity.WxId)
		err := GETObj(key, &deviceInfoEntity)
		if err != nil {
			return nil
		}
		deviceInfo := &baseinfo.DeviceInfo{}
		SimpleCopyProperties(deviceInfo, deviceInfoEntity)
		decodeSoftTypeXML, _ := base64.StdEncoding.DecodeString(deviceInfo.SoftTypeXML)
		deviceInfo.SoftTypeXML = string(decodeSoftTypeXML)
		decodeClientCheckDataXML, _ := base64.StdEncoding.DecodeString(deviceInfo.ClientCheckDataXML)
		deviceInfo.ClientCheckDataXML = string(decodeClientCheckDataXML)
		userInfo.DeviceInfo = deviceInfo
	}

	SimpleCopyProperties(userInfo, userInfoEntity)
	userInfo.AutoAuthKey, _ = base64.StdEncoding.DecodeString(userInfoEntity.AutoAuthKey)
	userInfo.SyncKey, _ = base64.StdEncoding.DecodeString(userInfoEntity.SyncKey)
	userInfo.FavSyncKey, _ = base64.StdEncoding.DecodeString(userInfoEntity.FavSyncKey)
	userInfo.SetLoginState(uint32(userInfoEntity.State))
	return userInfo
}
