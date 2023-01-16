package bizcgi

import (
	"encoding/json"
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"github.com/golang/protobuf/proto"
	"strconv"
)

// SendTextMsgReq 打包发送消息
func SendWxPayReq(wxAccount wxface.IWXAccount, reqItem *baseinfo.TenPayReqItem) (*wechat.TenPayResponse, error) {
	userInfo := wxAccount.GetUserInfo()
	var request wechat.TenPayRequest
	// baserequest
	baseReq := clientsdk.GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &reqItem.CgiCMD

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	var reqTextSKBuf wechat.SKBuiltinString_
	tmpLen := uint32(len(reqItem.ReqText))
	reqTextSKBuf.Len = &tmpLen
	reqTextSKBuf.Buffer = []byte(reqItem.ReqText)
	request.ReqText = &reqTextSKBuf

	// ReqTextWx
	var wxReqTextSKBuf wechat.SKBuiltinString_
	wxReqTextSKBuf.Buffer = []byte(reqItem.ReqText)
	wxReqTextSKBuf.Len = &tmpLen

	request.ReqTextWx = &wxReqTextSKBuf

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	reqData := clientsdk.Pack(userInfo, srcData, 385, 5)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 185,
		CgiUrl: "/cgi-bin/micromsg-bin/tenpay",
		Data:   reqData,
	}
	// 发送消息
	response, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		return nil, err
	}
	newSendMsgResponse := response.(*wechat.TenPayResponse)
	return newSendMsgResponse, nil
}

// SendCreatePreTransferReq 创建转账
func SendCreatePreTransferReq(wxAccount wxface.IWXAccount, toUserName string, fee uint, description string) (*baseinfo.PreTransferResp, error) {
	var req_text = "delay_confirm_flag=0&desc=" + description + "&fee=" + strconv.Itoa(int(fee)) + "&fee_type=CNY&pay_scene=31&receiver_name=" + toUserName + "&scene=31&transfer_scene=2"
	wcPaySign, err := clientsdk.TenPaySignDes3(req_text, "%^&*Tenpay!@#$")
	if err != nil {
		return nil, err
	}
	req_text += "&WCPaySign=" + wcPaySign
	tmpReqItem := &baseinfo.TenPayReqItem{}
	tmpReqItem.CgiCMD = 0x53
	tmpReqItem.ReqText = req_text
	tenPayResp, err := SendWxPayReq(wxAccount, tmpReqItem)
	if err != nil {
		return nil, err
	}
	retResp := &baseinfo.PreTransferResp{}
	retText := tenPayResp.GetRetText().GetBuffer()
	err = json.Unmarshal(retText, retResp)
	return retResp, err
}

// SendCreatePreTransferReq 确认转账
func SendConfirmPreTransferReq(wxAccount wxface.IWXAccount, bankType string, bankSerial string, reqKey string, payPassword string) (*wechat.TenPayResponse, error) {
	var req_text = "auto_deduct_flag=0&bank_type=" + bankType + "&bind_serial=" + bankSerial + "&busi_sms_flag=0&flag=3&passwd=" + payPassword + "&pay_scene=37&req_key=" + reqKey + "&use_touch=0"
	wcPaySign, err := clientsdk.TenPaySignDes3(req_text, "%^&*Tenpay!@#$")
	if err != nil {
		return nil, err
	}
	req_text += "&WCPaySign=" + wcPaySign
	tmpReqItem := &baseinfo.TenPayReqItem{}
	tmpReqItem.CgiCMD = 0
	tmpReqItem.ReqText = req_text
	return SendWxPayReq(wxAccount, tmpReqItem)
}
