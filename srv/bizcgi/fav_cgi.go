package bizcgi

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxlink"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
)

// SendBatchGetFavItemReq 获取单条收藏
func SendBatchGetFavItemReq(wxAccount wxface.IWXAccount, favID uint32) (*wechat.BatchGetFavItemResponse, error) {
	userInfo := wxAccount.GetUserInfo()

	var request wechat.BatchGetFavItemRequest
	// baseRequest
	baseReq := clientsdk.GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// Count
	favIDCount := uint32(1)
	request.Count = &favIDCount

	// FavIdList
	request.FavIdList = make([]byte, 0)
	tmpBytes := baseutils.EncodeVByte32(favID)
	request.FavIdList = append(request.FavIdList, tmpBytes[0:]...)

	srcData, _ := proto.Marshal(&request)
	sendData := clientsdk.Pack(userInfo, srcData, baseinfo.MMRequestTypeBatchGetFavItem, 5)
	longReq := &clientsdk.WXLongRequest{
		OpCode: 0,
		CgiUrl: "/cgi-bin/micromsg-bin/batchgetfavitem",
		Data:   sendData,
	}
	result, err := wxlink.WXSyncSend(wxAccount, longReq)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	response, _ := result.(wechat.BatchGetFavItemResponse)
	return &response, nil
}

// SendCheckFavCdnRequest 检测收藏cdn
func SendCheckFavCdnRequest(wxAccount wxface.IWXAccount, favCdn *baseinfo.CheckFavCdnItem) (*interface{}, error) {
	// 获取单条收藏详情
	tmpUserInfo := wxAccount.GetUserInfo()
	var request wechat.CheckCDNRequest

	// baseRequest
	baseReq := clientsdk.GetBaseRequest(tmpUserInfo)
	request.BaseRequest = baseReq
	var count = uint32(1)
	request.Count = &count
	list := make([]*wechat.CheckCDN, 0)
	checkCdn := &wechat.CheckCDN{
		DataId:         &favCdn.DataId,
		DataSourceId:   &favCdn.DataSourceId,
		DataSourceType: &favCdn.DataSourceType,
		FullMd5:        &favCdn.FullMd5,
		FullSize:       &favCdn.FullSize,
		Head256Md5:     &favCdn.Head256Md5,
		IsThumb:        &favCdn.IsThumb,
	}
	list = append(list, checkCdn)
	request.List = list
	// Count
	favIDCount := uint32(1)
	request.Count = &favIDCount

	srcData, _ := proto.Marshal(&request)
	retData := clientsdk.Pack(tmpUserInfo, srcData, 404, 5)
	longReq := &clientsdk.WXLongRequest{
		CgiUrl: "/cgi-bin/micromsg-bin/checkcdn",
		Data:   retData,
	}
	response, err := wxlink.WXShortSend(wxAccount, longReq)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func SendShareFavRequest(wxAccount wxface.IWXAccount, favID uint32, toUsername string) (*interface{}, error) {
	// 获取单条收藏详情
	tmpUserInfo := wxAccount.GetUserInfo()
	var request wechat.ShareFavRequest

	// baseRequest
	baseReq := clientsdk.GetBaseRequest(tmpUserInfo)
	request.BaseRequest = baseReq
	var tmpScene = uint32(0)
	request.Scene = &tmpScene
	request.ToUser = &toUsername
	// Count
	favIDCount := uint32(1)
	request.Count = &favIDCount

	// FavIdList
	request.FavIdList = make([]uint32, 0)
	request.FavIdList = append(request.FavIdList, favID)
	srcData, _ := proto.Marshal(&request)
	retData := clientsdk.Pack(tmpUserInfo, srcData, 608, 5)
	longReq := &clientsdk.WXLongRequest{
		CgiUrl: "/cgi-bin/micromsg-bin/sharefav",
		Data:   retData,
	}
	response, err := wxlink.WXShortSend(wxAccount, longReq)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
