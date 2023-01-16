package service

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"strconv"
)

// GetContactLabelListRequestService 获取标签列表
func GetContactLabelListRequestService(queryKey string) vo.DTO {
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
		resp, err := reqInvoker.SendGetContactLabelListRequest(true)
		if err != nil {
			return vo.NewFail("GetContactLabelListRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// AddContactLabelRequestService 添加标签
func AddContactLabelRequestService(queryKey string, m model.LabelModel) vo.DTO {
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
		if len(m.LabelNameList) == 0 {
			return vo.NewFail("没有要添加的标签")
		}
		resp, err := reqInvoker.SendAddContactLabelRequest(m.LabelNameList, true)
		if err != nil {
			return vo.NewFail("AddContactLabelRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// DelContactLabelRequestService 删除标签
func DelContactLabelRequestService(queryKey string, m model.LabelModel) vo.DTO {
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
		resp, err := reqInvoker.SendDelContactLabelRequest(m.LabelId)
		if err != nil {
			return vo.NewFail("SendDelContactLabelRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}

// ModifyLabelRequestService 修改标签
func ModifyLabelRequestService(queryKey string, m model.LabelModel) vo.DTO {
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
		resp, err := reqInvoker.SendModifyLabelRequest(m.UserLabelList)
		if err != nil {
			return vo.NewFail("SendDelContactLabelRequestService！" + err.Error())
		}
		return vo.NewSuccessObj(resp, "")
	})
}
