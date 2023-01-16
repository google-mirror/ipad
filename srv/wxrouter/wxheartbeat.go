package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/lunny/log"
)

// WXHeartBeatRouter 心跳包响应路由
type WXHeartBeatRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXHeartBeatRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析心跳包响应
	var hearBeatResp wechat.HeartBeatResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &hearBeatResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	if hearBeatResp.GetBaseResponse().GetRet() == 0 {
		log.Printf("[%s],[%s] HeartBeatSuccess next [%d]\n", currentUserInfo.GetUserName(), currentUserInfo.NickName, hearBeatResp.GetNextTime())
	} else {
		log.Printf("[%s],[%s] HeartBeatFail info [%s]\n", currentUserInfo.GetUserName(), currentUserInfo.NickName, hearBeatResp.GetBaseResponse().GetErrMsg())
	}
	//fmt.Println("心跳--->", hearBeatResp.GetBaseResponse().GetRet())
	// 等待 NextTime后再次发送心跳包
	//todo
	//currentWXConn.SendHeartBeatWaitingSeconds(hearBeatResp.GetNextTime())
	return hearBeatResp, nil
}
