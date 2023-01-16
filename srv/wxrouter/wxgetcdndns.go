package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXGetCDNDnsRouter 获取CDN dns信息
type WXGetCDNDnsRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXGetCDNDnsRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析获取CdnDns信息响应包
	var getCdnDNSResp wechat.GetCDNDnsResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &getCdnDNSResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 跟新DNS信息
	currentUserInfo.APPDnsInfo = getCdnDNSResp.GetAppDnsInfo()
	currentUserInfo.FAKEDnsInfo = getCdnDNSResp.GetFakeDnsInfo()
	currentUserInfo.SNSDnsInfo = getCdnDNSResp.GetSnsDnsInfo()
	currentUserInfo.DNSInfo = getCdnDNSResp.GetDnsInfo()
	return getCdnDNSResp, nil
}
