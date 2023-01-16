package db

import (
	"encoding/base64"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lunny/log"
	"strings"
	"time"
)

type MysqlDB struct {
	// MysqlDB MysqlDB
	Mysql *gorm.DB
}

// InitDB 初始化数据库
func (db *MysqlDB) InitDB() {
	//所有数据到我这来
	var mysqlUrl = srvconfig.GlobalSetting.MysqlConnectStr
	//获取系统
	//sysType := runtime.GOOS
	//if sysType == "linux" {
	//	if srvconfig.GlobalSetting.Dt != true {
	//		mysqlUrl = "wechat_mmtls:ye24Zxzf6mRMmzaG@tcp(119.45.28.143:3306)/wechat_mmtls?charset=utf8mb4&parseTime=true&loc=Local"
	//	}
	//}
	mysql, err := gorm.Open("mysql", mysqlUrl)
	db.Mysql = mysql
	mysql.DB().SetConnMaxLifetime(time.Minute * 1)
	//设置连接池中的最大闲置连接数
	mysql.DB().SetMaxIdleConns(10)
	//设置数据库的最大连接数量
	mysql.DB().SetMaxOpenConns(200)
	if err != nil {
		fmt.Println("failed to connect database:", err)
		_ = mysql.Close()
		return
	}
	fmt.Println("connect database success")
	mysql.SingularTable(true)
	//自动建表
	if !mysql.HasTable(&table.UserInfoEntity{}) {
		mysql.AutoMigrate(&table.UserInfoEntity{})
	}
	if !mysql.HasTable(&table.DeviceInfoEntity{}) {
		mysql.AutoMigrate(&table.DeviceInfoEntity{})
	}
	if !mysql.HasTable(&table.UserLoginLog{}) {
		mysql.AutoMigrate(&table.UserLoginLog{})
	}
	if !mysql.HasTable(&table.UserBusinessLog{}) {
		mysql.AutoMigrate(&table.UserBusinessLog{})
	}
	fmt.Println("create table success")
	mysql.LogMode(false)
}

// 提交登录日志
func (db *MysqlDB) SetLoginLog(loginType string, userInfo *baseinfo.UserInfo, errMsg string, state int32) {
	var userName string
	if len(userInfo.LoginDataInfo.UserName) > 0 {
		userName = userInfo.LoginDataInfo.UserName
	} else {
		userName = userInfo.WxId
	}
	log := &table.UserLoginLog{
		UUId:      userInfo.UUID,
		UserName:  userName,
		NickName:  userInfo.NickName,
		LoginType: loginType,
		RetCode:   state,
		ErrMsg:    errMsg,
	}
	log.TargetIp = srvconfig.GlobalSetting.TargetIp
	db.Mysql.Save(log)
}

// 获取登录日志
func (db *MysqlDB) GetLoginJournal(userName string) []table.UserLoginLog {
	logs := make([]table.UserLoginLog, 0)

	db.Mysql.Where("user_name=?", userName).Find(&logs)

	return logs
}

// 获取登录错误信息
func (db *MysqlDB) GetUSerLoginErrMsg(userName string) string {
	userInfoEntity := &table.UserInfoEntity{
		WxId: userName,
	}
	db.Mysql.Model(userInfoEntity).First(userInfoEntity)
	return userInfoEntity.ErrMsg
}

// 保存登录状态
func (db *MysqlDB) UpdateLoginStatus(uuid string, state int32, errMsg string) {
	data := make(map[string]interface{})
	data["State"] = state //零值字段
	data["ErrMsg"] = errMsg
	v := db.GetUserInfoEntity(uuid) // *table.UserInfoEntity
	if v != nil && v.State != state {
		db.Mysql.Model(&table.UserInfoEntity{
			WxId: v.WxId,
		}).Update(data)
		_ = PublishSyncMsgLoginState(v.WxId, uint32(state))
	}
}

// 保存初始化联系人结果
func (db *MysqlDB) UpdateInitContactStatus(userName string, state int32) {
	db.Mysql.Update(&table.UserInfoEntity{
		WxId:  userName,
		State: state,
	})
}

// 更新同步消息Key
func (db *MysqlDB) UpdateSyncMsgKey() {

}

// 更新用户信息
func (db *MysqlDB) UpdateUserInfo(userInfo *baseinfo.UserInfo) {
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
	db.Mysql.Model(new(table.UserInfoEntity)).Update(&userInfoEntity)

}

// SaveUserInfo 保存用户信息
func (db *MysqlDB) SaveUserInfo(userInfo *baseinfo.UserInfo) {

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

	db.Mysql.Save(&userInfoEntity)
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
		db.Mysql.Save(&deviceInfoEntity)
	}
}

// 获取所有登录用户
func (db *MysqlDB) QueryListUserInfo() []table.UserInfoEntity {
	var userInfoEntity = make([]table.UserInfoEntity, 0)
	db.Mysql.Where("state=?", 1).Find(&userInfoEntity)
	return userInfoEntity
}

// GetUserInfoEntity 获取登录信息
func (db *MysqlDB) GetUserInfoEntity(uuid string) *table.UserInfoEntity {
	var userInfoEntity table.UserInfoEntity
	if err := db.Mysql.Where("uuid=?", uuid).First(&userInfoEntity).Error; err != nil {
		return nil
	}
	return &userInfoEntity
}

// GetUserInfoEntity 获取登录信息
func (db *MysqlDB) GetUserInfoEntityByWxId(wxId string) *table.UserInfoEntity {
	var userInfoEntity table.UserInfoEntity
	if err := db.Mysql.Where("wxId=?", wxId).First(&userInfoEntity).Error; err != nil {
		return nil
	}
	return &userInfoEntity
}

func (db *MysqlDB) GetUSerInfoByUUID(uuid string) *baseinfo.UserInfo {
	userInfo := &baseinfo.UserInfo{}

	var userInfoEntity table.UserInfoEntity
	userInfoEntity.UUID = uuid

	if err := db.Mysql.Model(&table.UserInfoEntity{}).Where(&userInfoEntity).First(&userInfoEntity).Error; err != nil {
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
		if err := db.Mysql.First(&deviceInfoEntity).Error; err != nil {
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

// GetUserInfoByWXID 从数据库获取UserInfo
func (db *MysqlDB) GetUserInfoByWXID(wxID string) *baseinfo.UserInfo {
	userInfo := &baseinfo.UserInfo{}
	var userInfoEntity table.UserInfoEntity
	userInfoEntity.WxId = wxID
	if err := db.Mysql.First(&userInfoEntity).Error; err != nil {
		return nil
	}
	key := fmt.Sprintf("%s%s", "wechat:a16DeviceInfo:", wxID)
	exists, _ := Exists(key)
	if exists {
		//A16存redis
		key := fmt.Sprintf("%s%s", "wechat:a16DeviceInfo:", wxID)
		deviceInfoA16 := &baseinfo.AndroidDeviceInfo{}
		error := GETObj(key, &deviceInfoA16)
		if error != nil {
			log.Error("获取redis is error=" + error.Error())
		}
		userInfo.DeviceInfoA16 = deviceInfoA16
	} else {
		var deviceInfoEntity table.DeviceInfoEntity
		deviceInfoEntity.WxId = userInfoEntity.WxId
		if err := db.Mysql.First(&deviceInfoEntity).Error; err != nil {
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
	return userInfo
}

// GetDeviceInfo 根据URLkey查询DeviceInfo 查不到返回nil
func (db *MysqlDB) GetDeviceInfo(queryKey string) *baseinfo.DeviceInfo {
	var userInfoEntity table.UserInfoEntity
	db.Mysql.Where("querykey = ?", "queryKey").First(&userInfoEntity)
	var count int
	db.Mysql.Model(&table.UserInfoEntity{}).Where("querykey = ?", "queryKey").Count(&count)
	if count == 0 {
		return nil
	}

	deviceInfo := &baseinfo.DeviceInfo{}
	_ = SimpleCopyProperties(deviceInfo, userInfoEntity)
	return deviceInfo
}
