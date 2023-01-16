package wxrouter

import (
	"feiyu.com/wx/srv/wxmgr"
	"strings"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
)

// WXGetQrcodeRouter 获取个人/群二维码响应
type WXGetQrcodeRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXGetQrcodeRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()

	// 解析获取二维码响应包
	var getQrcodeResp wechat.GetQRCodeResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &getQrcodeResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 如果获取失败
	if getQrcodeResp.GetBaseResponse().GetRet() != 0 {
		return nil, nil
	}

	// 判断是否是群聊
	tmpUserName := string(wxResp.GetPackHeader().ReqData)
	if !strings.HasSuffix(tmpUserName, "@chatroom") {
		return nil, nil
	}

	return getQrcodeResp, nil
}
