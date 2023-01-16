package server

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/srv/wxcore"
	"feiyu.com/wx/srv/wxface"
	"feiyu.com/wx/srv/wxrouter"
)

// 微信服务器
var WxServer = InitWXServerRouter()

// WXServer 微信服务器，处理与微信的交互
type WXServer struct {
	wxMsgHandler wxface.IWXMsgHandler
}

// NewWXServer 新建微信服务对象
func NewWXServer() *WXServer {
	return &WXServer{
		wxMsgHandler: wxcore.NewWXMsgHandler(),
	}
}

// GetWXMsgHandler 获取微信消息管理器
func (wxs *WXServer) GetWXMsgHandler() wxface.IWXMsgHandler {
	return wxs.wxMsgHandler
}

// AddWXRouter 添加微信消息路由
func (wxs *WXServer) AddWXRouter(funcID uint32, wxRouter wxface.IWXRouter) {
	wxs.wxMsgHandler.AddRouter(funcID, wxRouter)
}

// 初始化微信响应路由
func InitWXServerRouter() *wxface.IWXServer {
	var wxServer wxface.IWXServer
	// 开启服务器
	wxServer = NewWXServer()
	// 注册微信响应路由
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetLoginQRCode, new(wxrouter.WXGetLoginQrcodeRouter)) // 获取登陆二维码响应
	wxServer.AddWXRouter(baseinfo.MMRequestTypeCheckLoginQRCode, new(wxrouter.WXCheckQrcodeRouter))  // 检测二维码状态请求
	wxServer.AddWXRouter(baseinfo.MMRequestTypePushQrLogin, new(wxrouter.WXPushQrCodeLoginRouter))   // 唤醒登录
	wxServer.AddWXRouter(baseinfo.MMRequestTypeManualAuth, new(wxrouter.WXManualAuthRouter))
	wxServer.AddWXRouter(baseinfo.MMRequestTypeHybridManualAuth, new(wxrouter.WXManualAuthRouter))                          // 登录
	wxServer.AddWXRouter(baseinfo.MMRequestTypeLogout, new(wxrouter.WXLogoutRouter))                                        // 退出登录
	wxServer.AddWXRouter(baseinfo.MMRequestTypeAutoAuth, new(wxrouter.WXAutoAuthRouter))                                    // token登陆
	wxServer.AddWXRouter(baseinfo.MMRequestTypeNewSync, new(wxrouter.WXNewSyncRouter))                                      // 同步消息，联系人
	wxServer.AddWXRouter(139, new(wxrouter.WXNewInitRouter))                                                                // 首次登录初始化
	wxServer.AddWXRouter(baseinfo.MMRequestTypeHeartBeat, new(wxrouter.WXHeartBeatRouter))                                  // 心跳包
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetProfile, new(wxrouter.WXGetProfileRouter))                                // 获取帐号信息
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetCdnDNS, new(wxrouter.WXGetCDNDnsRouter))                                  // 获取CDNDns信息
	wxServer.AddWXRouter(baseinfo.MMRequestTypeInitContact, new(wxrouter.WXInitContactRouter))                              // 初始化联系人
	wxServer.AddWXRouter(baseinfo.MMRequestTypeBatchGetContactBriefInfo, new(wxrouter.WXBatchGetContactBriefInfoReqRouter)) // 初始化联系人
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetContact, new(wxrouter.WXGetContactRouter))                                // 批量获取联系人信息
	wxServer.AddWXRouter(baseinfo.MMRequestTypeReceiveWxHB, new(wxrouter.WXReceiveHBRouter))                                // 接收红包
	wxServer.AddWXRouter(baseinfo.MMRequestTypeOpenWxHB, new(wxrouter.WXOpenHBRouter))                                      // 打开红包
	wxServer.AddWXRouter(baseinfo.MMRequestTypeOplog, new(wxrouter.WXOplogRouter))                                          // Oplog请求
	wxServer.AddWXRouter(baseinfo.MMRequestTypeNewSendMsg, new(wxrouter.WXNewSendMsgRouter))                                // 发送文本消息
	wxServer.AddWXRouter(baseinfo.MMRequestTypeFavSync, new(wxrouter.WXFavSyncRouter))                                      // 同步收藏
	wxServer.AddWXRouter(baseinfo.MMRequestTypeShareFav, new(wxrouter.WXShareFavRouter))                                    // 分享收藏
	wxServer.AddWXRouter(baseinfo.MMRequestTypeCheckFavCdn, new(wxrouter.WXCheckFavCdnRouter))                              // 分享收藏
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetFavInfo, new(wxrouter.WXGetFavInfoRouter))                                // 获取收藏信息
	wxServer.AddWXRouter(baseinfo.MMRequestTypeBatchGetFavItem, new(wxrouter.WXBatchGetFavItemRouter))                      // 获取单条收藏
	wxServer.AddWXRouter(baseinfo.MMRequestTypeMMSnsPost, new(wxrouter.WXSnsPostRouter))                                    // 发送朋友圈
	wxServer.AddWXRouter(baseinfo.MMRequestTypeMMSnsSync, new(wxrouter.WXSnsSyncRouter))                                    // 同步朋友圈
	wxServer.AddWXRouter(baseinfo.MMRequestTypeMMSnsUserPage, new(wxrouter.WXSnsUserPageRouter))                            // 朋友圈
	wxServer.AddWXRouter(baseinfo.MMRequestTypeMMSnsTimeLine, new(wxrouter.WXSnsTimeLineRouter))                            // 同步朋友圈
	wxServer.AddWXRouter(baseinfo.MMRequestTypeMMSnsComment, new(wxrouter.WXSnsCommentRouter))                              // 评论/点赞朋友圈
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetContactLabelList, new(wxrouter.WXGetContactLabelListRouter))              // 获取标签列表
	wxServer.AddWXRouter(baseinfo.MMRequestTypeAddContactLabel, new(wxrouter.WXAddContactLabelRouter))                      // 新增标签列表
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetQrCode, new(wxrouter.WXGetQrcodeRouter))                                  // 获取二维码
	wxServer.AddWXRouter(baseinfo.MMRequestTypeBindQueryNew, new(wxrouter.WXBindQueryNewRouter))                            // 获取钱包信息
	wxServer.AddWXRouter(baseinfo.MMRequestTypeThrIdGetA8Key, new(wxrouter.WXGetA8KeyRouter))                               // 获取a8key
	wxServer.AddWXRouter(baseinfo.MMRequestTypeGetA8Key, new(wxrouter.WXGetA8KeyRouter))                                    // 获取a8key
	wxServer.AddWXRouter(385, new(wxrouter.WXTenPayRouter))                                                                 // 支付
	wxServer.AddWXRouter(24, new(wxrouter.WXPullMsgRouter))                                                                 // 拉取消息
	return &wxServer
}
