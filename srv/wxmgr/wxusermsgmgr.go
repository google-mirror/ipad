package wxmgr

import (
	"feiyu.com/wx/db"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/wxface"
	"github.com/gogo/protobuf/proto"
	"github.com/lunny/log"
	"time"
)

// WXUSerMsgMgr WX用户消息管理器
type WXUSerMsgMgr struct {
	wxConn    wxface.IWXConnect
	msgList   chan *USerMsgItem
	endChan   chan bool
	isStarted bool
}

// USerMsgItem MsgItem
type USerMsgItem struct {
	MsgId      string
	TextMsg    string
	ImageData  []byte
	MsgType    uint32
	TOUserName string
}

// NewWXFileHelperMgr 新建文件传输助手管理器
func NewWXUSerMsgMgr(wxConn wxface.IWXConnect) *WXUSerMsgMgr {
	return &WXUSerMsgMgr{
		wxConn:    wxConn,
		msgList:   make(chan *USerMsgItem, 100),
		endChan:   make(chan bool, 1),
		isStarted: false,
	}
}

// Start 开启
func (wxfhm *WXUSerMsgMgr) Start() {
	go wxfhm.startDealMsg()
	wxfhm.isStarted = true
}

// Stop 关闭
func (wxfhm *WXUSerMsgMgr) Stop() {
	wxfhm.endChan <- true
	wxfhm.isStarted = false
}

// AddNewMsg 新增提示
func (wxfhm *WXUSerMsgMgr) AddNewTextMsg(MsgId, newMsg, toUSerName string) {
	if wxfhm.isStarted {
		newMsgItem := &USerMsgItem{}
		newMsgItem.MsgId = MsgId
		newMsgItem.TOUserName = toUSerName
		newMsgItem.MsgType = 1
		newMsgItem.TextMsg = newMsg
		wxfhm.msgList <- newMsgItem
	}
}

// AddImageMsg 新增图片消息
func (wxfhm *WXUSerMsgMgr) AddImageMsg(MsgId string, imgData []byte, toUSerName string) {
	if wxfhm.isStarted {
		newMsgItem := &USerMsgItem{}
		newMsgItem.MsgId = MsgId
		newMsgItem.TOUserName = toUSerName
		newMsgItem.MsgType = 2
		newMsgItem.ImageData = imgData
		wxfhm.msgList <- newMsgItem
	}
}

// 处理消息
func (wxfhm *WXUSerMsgMgr) startDealMsg() {
	currentWXAccount := WxAccountMgr.GetWXAccountByUserInfoUUID(wxfhm.wxConn.GetWXUuidKey())
	userInfo := currentWXAccount.GetUserInfo()
	currentReqInvoker := currentWXAccount.GetWXReqInvoker()
	for {
		// 最少1秒发送一次
		time.Sleep(3 * time.Second)
		select {
		case newMsgItem := <-wxfhm.msgList:

			if newMsgItem.MsgType == 1 {
				// 文字 todo
				log.Info("用户[" + userInfo.WxId + "]-->发送至-[" + newMsgItem.TOUserName + "]-->消息内容为:[" + newMsgItem.TextMsg + "]")
				//resp, err := bizcgi.SendTextMsgReq(wxfhm.wxConn, newMsgItem.TOUserName, newMsgItem.TextMsg, []string{}, 1)
				errMsg := ""
				ret := 0
				//if err != nil {
				//	errMsg = err.Error()
				//	ret = -1
				//}
				//发布同步信息消息
				req := &wechat.NewSendMsgResponse{
					BaseResponse: &wechat.BaseResponse{
						Ret: proto.Int32(int32(ret)),
						ErrMsg: &wechat.SKBuiltinString{
							Str: proto.String(errMsg),
						},
					},
					//Count:           resp.Count,
					//ChatSendRetList: resp.ChatSendRetList,
				}
				_ = db.PublishTxtImagePush(currentWXAccount.GetUserInfo(), req, newMsgItem.MsgId)
			} else if newMsgItem.MsgType == 2 {
				// 图片
				log.Info("用户[" + userInfo.WxId + "]-->发送至-[" + newMsgItem.TOUserName + "]-->消息内容为:[图片]")
				resp, err := currentReqInvoker.SendCdnUploadImageReuqest(newMsgItem.ImageData, newMsgItem.TOUserName)
				errMsg := ""
				ret := 0
				if err != nil {
					errMsg = err.Error()
					ret = -1
				}
				if !resp {
					ret = -1
				}
				req := &wechat.NewSendMsgResponse{
					BaseResponse: &wechat.BaseResponse{
						Ret: proto.Int32(int32(ret)),
						ErrMsg: &wechat.SKBuiltinString{
							Str: proto.String(errMsg),
						},
					},
				}
				//发布同步信息消息
				_ = db.PublishTxtImagePush(currentWXAccount.GetUserInfo(), req, newMsgItem.MsgId)
			}
		case <-wxfhm.endChan:
			return
		}
	}
}
