package wxrouter

import (
	"encoding/xml"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/lunny/log"

	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
)

// WXBatchGetFavItemRouter 获取单条收藏详情响应路由
type WXBatchGetFavItemRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXBatchGetFavItemRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentTaskMgr := currentWXConn.GetWXTaskMgr()
	//taskMgr, _ := currentTaskMgr.(*wxcore.WXTaskMgr)
	//_ = taskMgr.GetSnsTransTask()

	// 解析 获取单条收藏响应包
	var batchGetFavItemResp wechat.BatchGetFavItemResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &batchGetFavItemResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}
	// 解析单条收藏响应详情
	count := batchGetFavItemResp.GetCount()
	objectList := batchGetFavItemResp.GetObjectList()
	for index := uint32(0); index < count; index++ {
		tmpFavObject := objectList[index]
		// 反序列化
		objStr := tmpFavObject.GetObject()
		if len(objStr) <= 0 {
			continue
		}
		favItem := &baseinfo.FavItem{}
		err := xml.Unmarshal([]byte(objStr), favItem)
		if err != nil {
			log.Error("xml 解析error!")
			continue
		}

		// 自动转发
		favItem.FavItemID = tmpFavObject.GetFavId()
		//currentSnsTransTask.AddFavItem(favItem)
		//收藏发送mq
		//if favItem != nil {
		//	log.Info("----收藏id=", favItem.FavItemID)
		//	go db.PublishFavItem(currentWXAccount, favItem)
		//}
	}
	return batchGetFavItemResp, nil
}
