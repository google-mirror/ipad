package wxrouter

import (
	"encoding/json"
	"errors"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXReceiveHBRouter 接收红包响应处理器
type WXReceiveHBRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXReceiveHBRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
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
		return nil, errors.New(hongbaoResp.GetErrorMsg())
	}

	// 解析
	retHongBaoReceiveResp := &baseinfo.HongBaoReceiverResp{}
	err = json.Unmarshal(hongbaoResp.GetRetText().GetBuffer(), retHongBaoReceiveResp)
	if err != nil {
		return nil, err
	}
	return hongbaoResp, nil
}
