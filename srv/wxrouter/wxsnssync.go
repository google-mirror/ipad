package wxrouter

import (
	"errors"
	"feiyu.com/wx/srv/wxmgr"
	"strconv"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
)

// WXSnsSyncRouter 同步朋友圈
type WXSnsSyncRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (ssr *WXSnsSyncRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	//currentTaskMgr := currentWXConn.GetWXTaskMgr()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//taskMgr, _ := currentTaskMgr.(*wxcore.WXTaskMgr)
	//currentSnsTransTask := taskMgr.GetSnsTransTask()
	//currentReqInvoker := currentWXConn.GetWXReqInvoker()
	// 解析 同步朋友圈响应包
	var snsSyncResp wechat.SnsSyncResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &snsSyncResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 如果请求失败
	retCode := snsSyncResp.GetBaseResponse().GetRet()
	errMsg := snsSyncResp.GetBaseResponse().GetErrMsg().GetStr()
	if retCode != baseinfo.MMOk {
		return nil, errors.New("WXSnsSyncRouter err:" + strconv.Itoa(int(retCode)) + " msg = " + errMsg)
	}

	// 跟新key
	currentUserInfo.SnsSyncKey = snsSyncResp.GetKeyBuf().GetBuffer()
	// 遍历同步到的朋友圈
	tmpCmdList := snsSyncResp.GetCmdList()
	tmpCount := tmpCmdList.GetCount()
	if tmpCount > 0 {
		//todo
		//currentReqInvoker.SendSnsSyncRequest()
	}
	return snsSyncResp, nil
}
