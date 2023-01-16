package wxrouter

import (
	"feiyu.com/wx/clientsdk"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/defines"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxmgr"
)

// WXGetContactLabelListRouter 获取联系人标签路由
type WXGetContactLabelListRouter struct {
	WXBaseRouter
}

// Handle 处理conn业务的方法
func (hbr *WXGetContactLabelListRouter) Handle(wxResp wxface.IWXResponse) (interface{}, error) {
	//currentWXConn := wxResp.GetWXConncet()
	currentWXAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(wxResp.GetWXUuidKey())
	currentUserInfo := currentWXAccount.GetUserInfo()
	//currentReqInvoker := currentWXConn.GetWXReqInvoker()

	// 解析退出登陆响应包
	var getContableListResp wechat.GetContactLabelListResponse
	err := clientsdk.ParseResponseData(currentUserInfo, wxResp.GetPackHeader(), &getContableListResp)
	if err != nil {
		// 请求出问题了，应该关闭链接
		//currentWXConn.Stop()
		return nil, err
	}

	// 先判断是否有 收藏转发屏蔽，同步转发屏蔽
	hasFavLabel := false
	hasSyncLabel := false
	labelCount := getContableListResp.GetLabelCount()
	labelPairList := getContableListResp.GetLabelPairList()
	// 判断是否有收藏转发屏蔽组、同步转发屏蔽，如果没有就添加
	for index := uint32(0); index < labelCount; index++ {
		tmpLabelPair := labelPairList[index]
		if tmpLabelPair.GetLabelName() == defines.MFavTransShieldLabelName {
			hasFavLabel = true
			continue
		}

		if tmpLabelPair.GetLabelName() == defines.MSyncTransShieldLabelName {
			hasSyncLabel = true
			continue
		}
	}

	// 创建 屏蔽组
	newLabelNameList := []string{}
	if !hasFavLabel {
		newLabelNameList = append(newLabelNameList, defines.MFavTransShieldLabelName)
	}
	if !hasSyncLabel {
		newLabelNameList = append(newLabelNameList, defines.MSyncTransShieldLabelName)
	}
	if len(newLabelNameList) > 0 {
		//todo
		//currentReqInvoker.SendAddContactLabelRequest(newLabelNameList, false)
	}
	return getContableListResp, nil
}
