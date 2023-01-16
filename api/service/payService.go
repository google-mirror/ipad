package service

import (
	"encoding/json"
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxface"
	"strconv"
)

// 获取银行卡信息
func GetBandCardListService(queryKey string) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		tmpReqItem := &baseinfo.TenPayReqItem{}
		tmpReqItem.CgiCMD = 72
		tmpReqItem.ReqText = ""
		tenPayResp, err := reqInvoker.SendBandCardRequest(tmpReqItem)
		if err != nil {
			return vo.NewFail("GetBandCardListService err:" + err.Error())
		}
		// 解析响应
		retResp := &baseinfo.TenPayResp{}
		retText := tenPayResp.GetRetText().GetBuffer()
		err = json.Unmarshal(retText, retResp)
		if err != nil {
			return vo.NewFail("查询QB信息失败")
		}
		return vo.NewSuccessObj(retResp, "")
	})
}

// 生成自定义二维码
func GeneratePayQCodeService(queryKey string, req model.GeneratePayQCodeModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		var tenpayUrl = "delay_confirm_flag=0&desc=" + req.Name + "&fee=" + req.Money + "&fee_type=CNY&pay_scene=31&receiver_name=" + wxAccount.GetUserInfo().WxId + "&scene=31&transfer_scene=2"
		wcPaySign, err := clientsdk.TenPaySignDes3(tenpayUrl, "%^&*Tenpay!@#$")
		if err != nil {
			return vo.NewFail("no")
		}
		tenpayUrl += "&WCPaySign=" + wcPaySign
		tmpReqItem := &baseinfo.TenPayReqItem{}
		tmpReqItem.CgiCMD = 94
		tmpReqItem.ReqText = tenpayUrl
		tenPayResp, err := reqInvoker.SendTenPayRequest(tmpReqItem)
		if err != nil {
			return vo.NewFail("GetBandCardListService err:" + err.Error())
		}
		// 解析响应
		retResp := &baseinfo.GeneratePayQCodeResp{}
		retText := tenPayResp.GetRetText().GetBuffer()
		err = json.Unmarshal(retText, retResp)
		if err != nil {
			return vo.NewFail("查询QB信息失败")
		}
		return vo.NewSuccessObj(retResp, "")
	})
}

// 确认收款
func CollectmoneyServie(queryKey string, req model.CollectmoneyModel) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		tenpayUrl := "invalid_time=" + req.InvalidTime + "&op=confirm&total_fee=0&trans_id=" + req.TransFerId + "&transaction_id=" + req.TransactionId + "&username=" + req.ToUserName
		wcPaySign, err := clientsdk.TenPaySignDes3(tenpayUrl, "%^&*Tenpay!@#$")
		if err != nil {
			return vo.NewFail("no")
		}
		tenpayUrl += "&WCPaySign=" + wcPaySign
		tmpReqItem := &baseinfo.TenPayReqItem{}
		tmpReqItem.CgiCMD = 85
		tmpReqItem.ReqText = tenpayUrl
		tenPayResp, err := reqInvoker.SendTenPayRequest(tmpReqItem)
		if err != nil {
			return vo.NewFail("CollectmoneyServie err:" + err.Error())
		}
		return vo.NewSuccessObj(tenPayResp, "")
	})
}

// 拆红包
func OpenRedEnvelopesService(queryKey string, req baseinfo.HongBaoItem) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.SendOpenRedEnvelopesRequest(&req)
		if err != nil {
			return vo.NewFail("CollectmoneyServie err:" + err.Error())
		}
		return vo.NewSuccessObj(rsp, "")
	})
}

// 创建红包
func WXCreateRedPacketService(queryKey string, req baseinfo.RedPacket) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.SendWXCreateRedPacketRequest(&req)
		if err != nil {
			return vo.NewFail("CollectmoneyServie err:" + err.Error())
		}
		return vo.NewSuccessObj(rsp, "")
	})
}

// 查看红包详情
func QueryRedEnvelopesDetailService(queryKey string, req baseinfo.HongBaoItem) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.SendRedEnvelopesDetailRequest(&req)
		if err != nil {
			return vo.NewFail("QueryRedEnvelopesDetailService err:" + err.Error())
		}
		return vo.NewSuccessObj(string(rsp.GetRetText().GetBuffer()), "")
	})
}

// 查看红包领取列表
func GetRedPacketListService(queryKey string, req baseinfo.GetRedPacketList) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {

		loginState := wxAccount.GetLoginState()
		//判断在线情况
		if loginState == baseinfo.MMLoginStateNoLogin {
			return vo.NewFail("该账号需要重新登录！loginState == MMLoginStateNoLogin ")
		} else if !bizcgi.CheckOnLineStatus(wxAccount) {
			return vo.NewFail("账号离线,自动上线失败！loginState == " + strconv.Itoa(int(wxAccount.GetLoginState())))
		}
		reqInvoker := wxAccount.GetWXReqInvoker()
		rsp, err := reqInvoker.SendGetRedPacketListRequest(&req)
		if err != nil {
			return vo.NewFail("QueryRedEnvelopesDetailService err:" + err.Error())
		}
		return vo.NewSuccessObj(string(rsp.GetRetText().GetBuffer()), "")
	})
}

func CreatePreTransferService(queryKey string, req *baseinfo.CreatePreTransfer) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		tenPayResp, err := bizcgi.SendCreatePreTransferReq(wxAccount, req.ToUserName, req.Fee, req.Description)
		if err != nil {
			return vo.NewFail("CreatePreTransferService err:" + err.Error())
		}
		return vo.NewSuccessObj(tenPayResp, "")
	})
}

func ConfirmPreTransferService(queryKey string, req *baseinfo.ConfirmPreTransfer) vo.DTO {
	return checkExIdPerformNoCreateConnect(queryKey, func(wxAccount wxface.IWXAccount, newIWXConnect bool) vo.DTO {
		tenPayResp, err := bizcgi.SendConfirmPreTransferReq(wxAccount, req.BankType, req.BankSerial, req.ReqKey, req.PayPassword)
		if err != nil {
			return vo.NewFail("ConfirmPreTransferService err:" + err.Error())
		}
		return vo.NewSuccessObj(tenPayResp, "")
	})
}
