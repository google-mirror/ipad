package websrv

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/srv/srvconfig"
)

func createParamData(paramMap map[string]string) []byte {
	retString := ""
	for key, value := range paramMap {
		retString = retString + key + "=" + value + "&"
	}
	retBytes := []byte(retString)
	count := len(retBytes)
	retBytes = retBytes[0 : count-1]
	return retBytes
}

func getBaseParamMap(taskInfo *TaskInfo) map[string]string {
	retMap := make(map[string]string)
	curTime := strconv.Itoa(int(time.Now().UnixNano() / 1000000000))
	head := curTime + taskInfo.SignKey
	hsWebTime := curTime + baseutils.RandomBigHexString(20)
	apiKey := baseutils.Md5Value(head)

	retMap["name"] = taskInfo.Name
	retMap["app_number"] = taskInfo.AppNumber
	retMap["apikey"] = apiKey
	retMap["hswebtime"] = hsWebTime
	retMap["time"] = curTime
	return retMap
}

// UploadTaskStatus 上传任务状态
func UploadTaskStatus(webTask *WebTask) error {
	baseParamMap := getBaseParamMap(webTask.TaskInfo)
	baseParamMap["name"] = webTask.UserInfo.NickName
	baseParamMap["task_id"] = webTask.TaskInfo.TaskID
	baseParamMap["account"] = webTask.UserInfo.WxId
	baseParamMap["headimg"] = webTask.UserInfo.HeadURL
	baseParamMap["status"] = webTask.Status
	if webTask.Status == TaskStateLoginSuccess {
		tmpToken := base64.StdEncoding.EncodeToString(webTask.UserInfo.AutoAuthKey)
		baseParamMap["token"] = tmpToken
		baseParamMap["data"] = ""
		baseParamMap["autoauth"] = "data=" + "token=" + tmpToken
	}
	paramBytes := createParamData(baseParamMap)
	respData, err := TaskPost("http://"+srvconfig.GlobalSetting.WebDomain+"/statistic/Wechat/upd.html", paramBytes)
	if err != nil {
		baseutils.PrintLog(err.Error())
		return err
	}
	resp := &UploadTaskStatusResp{}
	err = json.Unmarshal(respData, resp)
	if err != nil {
		return err
	}

	// 如果显示44
	if resp.Code == 44 {
		return nil
	}

	return nil
}

// ReportWechatStatus 上报状态
func ReportWechatStatus(webTask *WebTask, tmpType string) error {
	baseParamMap := getBaseParamMap(webTask.TaskInfo)
	baseParamMap["name"] = webTask.UserInfo.NickName
	baseParamMap["task_id"] = webTask.TaskInfo.TaskID
	baseParamMap["account"] = webTask.UserInfo.WxId
	baseParamMap["task_type"] = "2"
	baseParamMap["type"] = tmpType
	paramBytes := createParamData(baseParamMap)
	respData, err := TaskPost("http://"+srvconfig.GlobalSetting.WebDomain+"/statistic/Wechat/reportStatus", paramBytes)
	if err != nil {
		baseutils.PrintLog(err.Error())
		return err
	}
	resp := &UploadTaskStatusResp{}
	err = json.Unmarshal(respData, resp)
	if err != nil {
		return err
	}
	// 上报错误
	if resp.Code == 44 {
		return errors.New("上报错误")
	}
	return nil
}
