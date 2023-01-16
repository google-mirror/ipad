package wxcore

import (
	"errors"
	"feiyu.com/wx/common"
	"feiyu.com/wx/srv/utils"
	"feiyu.com/wx/srv/wxface"
	"github.com/gogf/gf/os/grpool"
	"github.com/lunny/log"
	"reflect"
	"strconv"
)

// WXMsgHandler 微信响应管理器
type WXMsgHandler struct {
	wxRouterMap map[uint32]wxface.IWXRouter //存放每个MsgId 所对应的处理方法的map属性
}

// NewWXMsgHandler 新建微信消息处理器
func NewWXMsgHandler() *WXMsgHandler {
	return &WXMsgHandler{
		wxRouterMap: make(map[uint32]wxface.IWXRouter),
	}
}

// AddRouter 增加微信消息路由
func (wxmh *WXMsgHandler) AddRouter(respID uint32, wxRouter wxface.IWXRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := wxmh.wxRouterMap[respID]; ok {
		return
	}
	//2 添加msg与api的绑定关系
	wxmh.wxRouterMap[respID] = wxRouter
}

// GetRouterByRespID 根据响应ID获取对应的路由
func (wxmh *WXMsgHandler) GetRouterByRespID(urlID uint32) wxface.IWXRouter {
	handler, ok := wxmh.wxRouterMap[urlID]
	if !ok {
		return nil
	}
	return handler
}

// doMsgHandler 马上以阻塞方式处理消息
func (wxmh *WXMsgHandler) doMsgHandler(response wxface.IWXResponse) {
	defer utils.TryE(response.GetWXUuidKey())
	var result interface{}
	var err error
	handler, ok := wxmh.wxRouterMap[response.GetPackHeader().URLID]
	defer func() {
		key := response.GetWXUuidKey() + "-" + strconv.Itoa(int(response.GetPackHeader().SeqId))
		common.Send(key, &common.Pack{
			Content: result,
			Err:     err,
		})
	}()
	if !ok {
		result, err = nil, errors.New("URLID is 0, No handler")
		return
	}
	//执行对应处理方法
	err = handler.PreHandle(response)
	if err == nil {
		result, err = handler.Handle(response)
		if err == nil {
			err = handler.PostHandle(response)
			if err == nil {
				return
			}
		}
	}
	log.Error(reflect.TypeOf(handler), "->处理结果时出错", err.Error())
}

// SendWXRespToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (wxmh *WXMsgHandler) SendWXRespToTaskQueue(response wxface.IWXResponse) {
	grpool.Add(func() {
		defer utils.TryE("")
		wxmh.doMsgHandler(response)
	})

}
