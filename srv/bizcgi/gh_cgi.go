package bizcgi

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
)

// SendClickMenuReq 打包操作菜单
func SendClickMenuReq(wxAccount wxface.IWXAccount, ghUsername string, menuId string, menuKey string) (*wechat.ClickCommandResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	content := "#bizmenu#<info><id>" + menuId + "</id><key>" + menuKey + "</key><status>menu_click</status><content></content></info>"
	request := &wechat.ClickCommandRequest{
		BaseRequest: clientsdk.GetBaseRequest(userInfo),
		BizUserName: &ghUsername,
		ClickType:   proto.Uint32(1),
		ClickInfo:   &content,
	}
	// 打包数据
	src, _ := proto.Marshal(request)
	sendData := clientsdk.Pack(userInfo, src, 359, 5)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 176,
		CgiUrl: "/cgi-bin/micromsg-bin/clickcommand",
		Data:   sendData,
	}
	result, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	clickCommandResponse := result.(*wechat.ClickCommandResponse)
	return clickCommandResponse, nil
}
