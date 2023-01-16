package wxrouter

import (
	"errors"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/lunny/log"
)

// WXBaseRouter 实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type WXBaseRouter struct{}

// PreHandle 在处理conn业务之前的钩子方法
func (wxbr *WXBaseRouter) PreHandle(response wxface.IWXResponse) error {
	//currentWXConn := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(response.GetWXUuidKey())
	packHeader := response.GetPackHeader()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(response.GetWXUuidKey())
	if packHeader.RetCode == baseinfo.MMErrSessionTimeOut {
		flag := string(response.GetWXUuidKey()) + "->[" + currentWXAccount.GetUserInfo().GetUserName() + "]->"
		log.Println(flag, "sessionTimeOut触发二次登录")
		_, err := bizcgi.SendAutoAuthRequest(currentWXAccount)
		if err != nil {
			log.Println("token auth登录请求error", err.Error())
			//currentWXConn.Stop()
			return err
		}
		return errors.New("重新连接")
	}
	if packHeader.Data == nil {
		return errors.New("返回内容为空，无法反序列化")
	}
	return nil
}

// Handle 处理conn业务的方法
func (wxbr *WXBaseRouter) Handle(response wxface.IWXResponse) (interface{}, error) {
	return nil, nil
}

// PostHandle 处理conn业务之后的钩子方法
func (wxbr *WXBaseRouter) PostHandle(response wxface.IWXResponse) error {
	return nil
}
