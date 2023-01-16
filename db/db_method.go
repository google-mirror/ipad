package db

import (
	"errors"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
)

// DB DB
var DB IDBInvoker

// InitDB 初始化数据库
func InitDB() {
	//所有数据到我这来
	var dbType = srvconfig.GlobalSetting.DbType
	if dbType == "mysql" {
		DB = &MysqlDB{}
	} else {
		DB = &RedisDb{}
	}
	DB.InitDB()
}

// 提交登录日志
func SetLoginLog(loginType string, userInfo *baseinfo.UserInfo, errMsg string, state int32) {
	DB.SetLoginLog(loginType, userInfo, errMsg, state)
}

// 获取登录日志
func GetLoginJournal(userName string) []table.UserLoginLog {
	return DB.GetLoginJournal(userName)
}

// 保存登录状态
func UpdateLoginStatus(uuid string, state int32, errMsg string) {
	DB.UpdateLoginStatus(uuid, state, errMsg)
}

// 更新同步消息Key
func UpdateSyncMsgKey() {

}

// 更新用户信息
func UpdateUserInfo(userInfo *baseinfo.UserInfo) {
	DB.UpdateUserInfo(userInfo)
}

// SaveUserInfo 保存用户信息
func SaveUserInfo(userInfo *baseinfo.UserInfo) {
	DB.SaveUserInfo(userInfo)
}

// GetUserInfoEntity 获取登录信息
func GetUserInfoEntity(uuid string) *table.UserInfoEntity {
	return DB.GetUserInfoEntity(uuid)
}

func GetUSerInfoByUUID(uuid string) *baseinfo.UserInfo {
	return DB.GetUSerInfoByUUID(uuid)
}

//// GetDeviceInfo 根据URLkey查询DeviceInfo 查不到返回nil
//func GetDeviceInfo(queryKey string) *baseinfo.DeviceInfo {
//	return DB.GetDeviceInfo(queryKey)
//}

// SimpleCopyProperties 拷贝属性
func SimpleCopyProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}
		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}
	return nil
}
