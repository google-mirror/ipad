package db

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db/table"
)

// IDBInvoker 持久化管理
type IDBInvoker interface {
	InitDB()
	SetLoginLog(loginType string, userInfo *baseinfo.UserInfo, errMsg string, state int32)
	GetLoginJournal(userName string) []table.UserLoginLog
	//保存登录状态
	UpdateLoginStatus(uuid string, state int32, errMsg string)
	//更新用户信息
	UpdateUserInfo(userInfo *baseinfo.UserInfo)
	// SaveUserInfo 保存用户信息
	SaveUserInfo(userInfo *baseinfo.UserInfo)
	// GetUserInfoEntity 获取登录信息
	GetUserInfoEntity(uuid string) *table.UserInfoEntity
	// GetUSerInfoByUUID 从数据库获取UserInfo
	GetUSerInfoByUUID(uuid string) *baseinfo.UserInfo
}
