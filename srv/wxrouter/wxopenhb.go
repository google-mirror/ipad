package wxrouter

import (
	"encoding/json"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXOpenHBRouter 接收红包响应处理器
type WXOpenHBRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXOpenHBRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 心跳包响应
	var hongbaoResp wechat.HongBaoRes
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &hongbaoResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 如果 返回错误码
	if hongbaoResp.GetErrorType() != 0 {
		baseutils.PrintLog("WXReceiveHBRouter err: " + hongbaoResp.GetErrorMsg())
		return nil, nil
	}

	// 解析
	retHongBaoOpenResp := &baseinfo.HongBaoOpenResp{}
	err = json.Unmarshal(hongbaoResp.GetRetText().GetBuffer(), retHongBaoOpenResp)
	if err != nil {
		return nil, err
	}
	// 抢红包成功 可以增加一条记录
	return hongbaoResp, nil
}
