package websrv

import "feiyu.com/wx/clientsdk/baseinfo"

const (
	// TaskStateLoginSuccess 登录成功状态
	TaskStateLoginSuccess = string("2")
	// TaskStateLogout 退出登录
	TaskStateLogout = string("6")
	// TaskStateCheckLoginQrcodeSuccess 检测状态
	TaskStateCheckLoginQrcodeSuccess = string("1")
	// TaskStateCheckLoginQrcodeStartLogin 用户点击了登录
	TaskStateCheckLoginQrcodeStartLogin = string("7")
	// TaskStateMismatching 账号不匹配
	TaskStateMismatching = string("9")

	// TaskTypeUploadStatus 上传状态
	TaskTypeUploadStatus uint32 = 1
	// TaskTypeReportHeart 发送心跳包
	TaskTypeReportHeart uint32 = 2
)

// TaskInfo 任务信息
type TaskInfo struct {
	UUID      string
	TaskID    string
	Name      string
	AppNumber string
	Account   string
	SignKey   string
	DeviceId  string
}

// UploadTaskStatusResp 上传状态响应
type UploadTaskStatusResp struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// WebTask web任务
type WebTask struct {
	UserInfo *baseinfo.UserInfo
	TaskInfo *TaskInfo
	Status   string
	Type     uint32
	Count    uint32
}
