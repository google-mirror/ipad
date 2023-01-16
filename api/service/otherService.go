package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"strconv"
)

// GetQrCodeService 获取二维码
func GetQrCodeService(queryKey string, m model.GetQrCodeModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGetQrCodeRequest(m.Id)
		if err != nil {
			return vo.NewFail("GetQrCodeService err:" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// 查看附近的人
func GetPeopleNearbyService(queryKey string, m model.PeopleNearbyModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		//取基本信息

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		//获取请求管理器
		reqInvoker := wxAccount.GetWXReqInvoker()
		resp, err := reqInvoker.SendGetPeopleNearbyResultRequest(m.Longitude, m.Latitude)
		if err != nil {
			return vo.NewFail("GetPeopleNearbyService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
