package wxrouter

import (
	"encoding/json"
	"errors"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXBindQueryNewRouter 查询钱包信息响应
type WXBindQueryNewRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXBindQueryNewRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析钱包信息响应
	var tenPayResp wechat.TenPayResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &tenPayResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 判断是否失败
	if tenPayResp.GetRetText().GetLen() <= 0 {
		return tenPayResp, errors.New("查询QB信息失败")
	}

	// 解析响应
	retResp := &baseinfo.TenPayResp{}
	retText := tenPayResp.GetRetText().GetBuffer()
	err = json.Unmarshal(retText, retResp)
	if err != nil {
		return nil, err
	}
	return retResp, nil
}
