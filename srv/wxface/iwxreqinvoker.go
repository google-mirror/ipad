package wxface

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
)

// IWXReqInvoker 微信请求调用器
type IWXReqInvoker interface {
	//发送登录短信
	SendWxBindOpMobileForRequest(OpCode int64, PhoneNumber string, VerifyCode string) (*wechat.BindOpMobileForRegResponse, error)
	// SendHybridManualAutoRequest
	SendHybridManualAutoRequest(newPass string, wxID string, ver byte) error
	//获取设备
	SendGetSafetyInfoRequest() (*wechat.GetSafetyInfoResponse, error)
	//删除设备
	SendDelSafeDeviceRequest(deviceUUID string) (*wechat.DelSafeDeviceResponse, error)
	//检测微信登录环境
	SendCheckCanSetAliasRequest() (*wechat.CheckCanSetAliasResp, error)
	//扫码登录新设备
	SendExtDeviceLoginConfirmGetRequest(url string) (*wechat.ExtDeviceLoginConfirmOKResponse, error)
	SendGetProfileNewRequest() (*wechat.GetProfileResponse, error)
	//同步消息
	SendWxSyncMsg(key string) (*wechat.NewSyncResponse, error)

	SendNewInitSyncRequest() (interface{}, error)
	// 会重新链接服务器 发送Token登陆请求
	//SendAutoAuthRequest() (interface{}, error)
	// 发送初始化联系人请求
	SendInitContactRequest(contactSeq uint32) error
	//分页获取联系人
	SendGetContactListPageRequest(CurrentWxcontactSeq uint32, CurrentChatRoomContactSeq uint32) (*wechat.InitContactResp, error)
	// 批量获取联系人详情
	SendBatchGetContactBriefInfoReq(userWxidList []string) error
	// 获取联系人信息列表
	SendGetContactRequest(userInfoList []string, antisPanTicketList []string, chatRoomWxidList []string, needResp bool) (*wechat.GetContactResponse, error)
	// 获取联系人信息列表
	SendGetContactRequestForHB(userInfoList string) (*wechat.GetContactResponse, error)
	// 获取联系人信息列表List
	SendGetContactRequestForList(userInfoList []string, roomWxIDList []string) (*wechat.GetContactResponse, error)
	//获取好友关系状态
	SendGetFriendRelationRequest(userName string) (*wechat.MMBizJsApiGetUserOpenIdResponse, error)
	// 接收红包
	SendReceiveWxHBRequest(hbItem *baseinfo.HongBaoItem) error
	// 打开红包
	SendOpenWxHBRequest(hbItem *baseinfo.HongBaoItem, timingIdentifier string) error
	//拆红包
	SendOpenRedEnvelopesRequest(hbItem *baseinfo.HongBaoItem) (*wechat.HongBaoRes, error)
	//创建红包
	SendWXCreateRedPacketRequest(hbItem *baseinfo.RedPacket) (*wechat.HongBaoRes, error)
	//查看红包详情
	SendRedEnvelopesDetailRequest(hbItem *baseinfo.HongBaoItem) (*wechat.HongBaoRes, error)
	//查看红包列表
	SendGetRedPacketListRequest(hbItem *baseinfo.GetRedPacketList) (*wechat.HongBaoRes, error)
	// 发送Oplog请求
	SendOplogRequest(modifyItems []*baseinfo.ModifyItem) error
	//发送企业Oplog请求
	SendQWOpLogRequest(cmdId int64, value []byte) error
	// 获取群/个人二维码
	SendGetQRCodeRequest(userName string) error
	// 同步收藏
	SendFavSyncRequest() (interface{}, error)
	SendFavSyncListRequestResult(keyBuf string) (*wechat.SyncResponse, error)
	// 获取收藏信息
	SendGetFavInfoRequest() error
	SendGetFavInfoRequestResult() (*wechat.GetFavInfoResponse, error)
	// 删除收藏
	SendBatchDelFavItemRequest(favID uint32) error
	SendBatchDelFavItemRequestResult(favID uint32) (*wechat.BatchDelFavItemResponse, error)
	// 获取CdnDns信息
	SendGetCDNDnsRequest() (interface{}, error)
	//上报设备
	SendReportstrategyRequest() (*wechat.GetReportStrategyResp, error)
	// 发送朋友圈
	SendSnsPostRequest(postItem *baseinfo.SnsPostItem) error
	// 发送朋友圈
	SendSnsPostRequestNew(postItem *baseinfo.SnsPostItem) (*wechat.SnsPostResponse, error)
	//设置朋友圈可见天数
	SetFriendCircleDays(postItem *model.SetFriendCircleDaysModel) error
	// 操作朋友圈
	SendSnsObjectOpRequest(opItems []*baseinfo.SnsObjectOpItem) (*wechat.SnsObjectOpResponse, error)
	// 获取指定好友朋友圈
	SendSnsUserPageRequest(userName string, firstPageMd5 string, maxID uint64, needResp bool) (*wechat.SnsUserPageResponse, error)
	// 同步转发朋友圈
	SendSnsPostRequestByXML(timeLineObj *baseinfo.TimelineObject, blackList []string) error
	// 获取指定的朋友圈详情
	SendSnsObjectDetailRequest(snsID uint64) (*wechat.SnsObject, error)
	// 获取朋友圈首页
	SendSnsTimeLineRequest(firstPageMD5 string, maxID uint64) error
	SendSnsTimeLineRequestResult(firstPageMD5 string, maxID uint64) (*wechat.SnsTimeLineResponse, error)
	// 发送评论/点赞请求
	SendSnsCommentRequest(commentItem *baseinfo.SnsCommentItem) error
	// 同步朋友圈
	SendSnsSyncRequest() error
	// 获取联系人标签列表
	SendGetContactLabelListRequest(needResp bool) (*wechat.GetContactLabelListResponse, error)
	// 添加标签
	SendAddContactLabelRequest(newLabelList []string, needResp bool) (*wechat.AddContactLabelResponse, error)
	// 删除标签
	SendDelContactLabelRequest(labelId string) (*wechat.DelContactLabelResponse, error)
	// 修改标签
	SendModifyLabelRequest(userLabelList []baseinfo.UserLabelInfoItem) (*wechat.ModifyContactLabelListResponse, error)
	// 查询钱包信息
	SendBindQueryNewRequest(reqItem *baseinfo.TenPayReqItem) error
	//获取余额以及银行卡信息
	SendBandCardRequest(reqItem *baseinfo.TenPayReqItem) (*wechat.TenPayResponse, error)
	//支付方法
	SendTenPayRequest(reqItem *baseinfo.TenPayReqItem) (*wechat.TenPayResponse, error)
	// 下载请求
	SendCdnDownloadReuqest(downItem *baseinfo.DownMediaItem) (*baseinfo.CdnDownloadResponse, error)
	//下载图片
	GetMsgBigImg(m model.GetMsgBigImgModel) (*wechat.GetMsgImgResponse, error)
	// Cdn上传高清图片
	SendCdnSnsUploadImageReuqest(imgData []byte) (*baseinfo.CdnSnsImageUploadResponse, error)
	//发送CDN朋友圈视频下载请求
	SendCdnSnsVideoDownloadReuqest(encKey uint64, tmpURL string) ([]byte, error)
	// 发送CDN朋友圈上传视频请求
	SendCdnSnsVideoUploadReuqest(videoData []byte, thumbData []byte) (*baseinfo.CdnSnsVideoUploadResponse, error)
	// Cdn发送图片给好友
	SendCdnUploadImageReuqest(imgData []byte, toUserName string) (bool, error)
	//发送图片
	SendUploadImageNewRequest(imgData []byte, toUserName string) (*wechat.UploadMsgImgResponse, error)
	// cdn发送视频
	SendCdnUploadVideoRequest(toUserName string, imgData string, videoData []byte) (*baseinfo.CdnMsgVideoUploadResponse, error)
	// Cdn发送图片给文件传输助手
	SendImageToFileHelper(imgData []byte) (bool, error)
	// 转发图片
	ForwardCdnImageRequest(item baseinfo.ForwardImageItem) (*wechat.UploadMsgImgResponse, error)
	// 转发视频
	ForwardCdnVideoRequest(item baseinfo.ForwardVideoItem) (*wechat.UploadVideoResponse, error)
	// 发送app消息
	SendAppMessage(msgXml, toUSerName string, contentType uint32) (*wechat.SendAppMsgResponse, error)
	// SendEmojiRequest 发送表情
	SendEmojiRequest(md5 string, toUSerName string, length int32) (*wechat.SendAppMsgResponse, error)
	//发送表情new  动图
	ForwardEmojiRequest(md5 string, toUSerName string, length int32) (*wechat.UploadEmojiResponse, error)
	//群发文字
	SendGroupMassMsgTextRequest(toUserName []string, content string) (*wechat.MassSendResponse, error)
	//群发图片
	SendGroupMassMsgImageRequest(toUserName []string, ImageBase64 []byte) (*wechat.MassSendResponse, error)
	//群拍一拍
	SendSendPatRequest(chatRoomName string, toUserName string, scene int64) (*wechat.SendPatResponse, error)
	//下载语音
	SendGetMsgVoiceRequest(toUserName, newMsgId, bufid string, length int) (*vo.DownloadVoiceData, error)
	//群发
	// 设置群公告
	SetChatRoomAnnouncementRequest(roomId, content string) (*wechat.SetChatRoomAnnouncementResponse, error)
	// 获取群成员详细
	GetChatroomMemberDetailRequest(roomId string) (*wechat.GetChatroomMemberDetailResponse, error)
	//获取群详细
	SetGetChatRoomInfoDetailRequest(roomId string) (*wechat.GetChatRoomInfoDetailRequest, error)
	// 退出群聊
	GetQuitChatroomRequest(chatRoomName string) error
	// 创建群
	SendCreateChatRoomRequest(topIc string, userList []string) (*wechat.CreateChatRoomResponse, error)
	// 邀请群成员
	SendInviteChatroomMembersRequest(chatRoomName string, userList []string) (*wechat.CreateChatRoomResponse, error)
	// 添加好友进群
	SendAddChatRoomMemberRequest(chatRoomName string, userList []string) (*wechat.AddChatRoomMemberResponse, error)
	//删除群成员
	SendDelDelChatRoomMemberRequest(chatRoomName string, delUserList []string) (*wechat.DelChatRoomMemberResponse, error)
	//转让群
	SendTransferGroupOwnerRequest(chatRoomName, newOwnerUserName string) (*wechat.TransferChatRoomOwnerResponse, error)
	//添加群管理
	SendAddChatroomAdminRequest(chatRoomName string, userList []string) (*wechat.AddChatRoomAdminResponse, error)
	//删除群管理
	SendDelChatroomAdminRequest(chatRoomName string, userList []string) (*wechat.DelChatRoomAdminResponse, error)
	//获取群例表
	SendWXSyncContactRequest() (*vo.GroupData, error)
	// 链接授权
	GetA8KeyRequest(opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*wechat.GetA8KeyResp, error)
	// 链接链接扫码进群
	GetA8KeyGroupRequest(opCode, scene uint32, reqUrl string, getType baseinfo.GetA8KeyType) (*wechat.GetA8KeyResp, error)
	// 小程序授权
	JSLoginRequest(appId string) (*wechat.JSLoginResponse, error)
	// 小程序授权
	JSOperateWxDataRequest(appId string) (*wechat.JSOperateWxDataResponse, error)
	// app 授权
	SdkOauthAuthorizeRequest(appId string, sdkName string, packageName string) (*wechat.SdkOauthAuthorizeConfirmNewResp, error)
	// 搜索好友
	SendSearchContactRequest(opCode, fromScene, searchScene uint32, userName string) (*wechat.SearchContactResponse, error)
	// 好友验证/加好友/关注公众号
	VerifyUserRequest(opCode uint32, verifyContent string, scene byte, V1, V2, ChatRoomUserName string) (*wechat.VerifyUserResponse, error)
	// 上传手机通讯录好友
	UploadMContact(mobile string, mobileList []string) (*wechat.UploadMContactResponse, error)
	// 获取手机通讯录好友
	GetMFriend() (*wechat.GetMFriendResponse, error)
	// 获取证书
	SendCertRequest() (*wechat.GetCertResponse, error)
	// 发送二维码授权请求
	SendQRConnectAuthorize(qrUrl string) (*wechat.QRConnectAuthorizeResp, error)
	// 发送二维码授权请求确认
	SendQRConnectAuthorizeConfirm(qrUrl string) (*wechat.SdkOauthAuthorizeConfirmNewResp, error)
	//授权链接
	SendGetMpA8Request(url string, opcode uint32) (*wechat.GetA8KeyResp, error)
	// 获取登录设备信息
	SendOnlineInfo() (*wechat.GetOnlineInfoResponse, error)
	// 获取二维码
	SendGetQrCodeRequest(id string) (*wechat.GetQRCodeResponse, error)
	//查看附近的人
	SendGetPeopleNearbyResultRequest(longitude float32, latitude float32) (*wechat.LbsResponse, error)
	// 撤销消息
	SendRevokeMsgRequest(newMsgId string, clientMsgId uint64, toUserName string) (*wechat.RevokeMsgResponse, error)
	// 删除好友
	SendDelContactRequest(userName string) error
	// 修改资料
	SendModifyUserInfoRequest(city, country, nickName, province, signature string, sex uint32, initFlag uint32) error
	//修改昵称
	SendUpdateNickNameRequest(cmd uint32, val string) error
	//设置姓名
	SetNickNameService(cmd uint32, val string) error
	//设置性别
	SetSexService(val uint32, country string, city string, province string) error
	//修改加好友需要验证属性
	UpdateAutopassRequest(SwitchType uint32) error
	//修改头像
	UploadHeadImage(base64 string) (*wechat.UploadHDHeadImgResponse, error)
	// 修改密码
	SendChangePwdRequest(oldPwd, NewPwd string, OpCode uint32) (*wechat.BaseResponse, error)
	// 修改备注
	SendModifyRemarkRequest(userName string, remarkName string) error
	// 发送语音
	SendUploadVoiceRequest(toUserName string, voiceData string, voiceSecond, voiceFormat int32) (*wechat.UploadVoiceResponse, error)
	//设置微信号
	SetWechatRequest(alisa string) (*wechat.GeneralSetResponse, error)
	//设置微信步数
	UpdateStepNumberRequest(number uint64) (*wechat.UploadDeviceStepResponse, error)
	//换绑手机
	SendBindingMobileRequest(mobile, verifyCode string) (*wechat.BindOpMobileResponse, error)
	//发送验证码
	SendVerifyMobileRequest(mobile string, opcode uint32) (*wechat.BindOpMobileResponse, error)
	//获取步数列表
	SendGetUserRankLikeCountRequest(rankId string) (*wechat.GetUserRankLikeCountResponse, error)
	//提取企业 wx 详情
	SendQWContactRequest(openIm, chatRoom, t string) (*wechat.GetQYContactResponse, error)
	//提取全部的企业通寻录
	SendQWSyncContactRequest() (*wechat.GetQYContactResponse, error)
	//备注企业
	SendQWRemarkRequest(toUserName string, name string) error
	//创建企业群
	SendQWCreateChatRoomRequest(userList []string) (*wechat.CreateQYChatRoomResponese, error)
	//搜手机或企业对外名片链接提取验证
	SendQWSearchContactRequest(tg string, fromScene uint64, userName string) (*wechat.SearchQYContactResponse, error)
	//向企业微信打招呼
	SendQWApplyAddContactRequest(toUserName, v1, Content string) error
	//单向加企业微信
	SendQWAddContactRequest(toUserName, v1, Content string) error
	//拉取企业微信群
	SendQWSyncChatRoomRequest(key string) (*vo.QYChatroomContactVo, error)
	//转让企业群
	SendQWChatRoomTransferOwnerRequest(chatRoomName string, toUserName string) (*wechat.BaseResponse, error)
	//直接拉好友进群
	SendQWAddChatRoomMemberRequest(chatRoomName string, toUserName []string) (*wechat.QYAddChatRoomMemberResponse, error)
	//发送群邀请链接
	SendQWInviteChatRoomMemberRequest(chatRoomName string, toUserName []string) (*wechat.BaseResponse, error)
	//删除企业群群成员
	SendQWDelChatRoomMemberRequest(chatRoomName string, toUserName []string) (*wechat.QYDelChatRoomMemberResponse, error)
	//提取企业群全部成员
	SendQWGetChatRoomMemberRequest(chatRoomName string) (*wechat.GetQYChatroomMemberDetailResponse, error)
	//提取企业群名称公告设定等信息
	SendQWGetChatroomInfoRequest(chatRoomName string) (*wechat.QYChatroomContactResponse, error)
	//提取企业群二维码
	SendQWGetChatRoomQRRequest(chatRoomName string) (*wechat.QYGetQRCodeResponse, error)
	//增加企业管理员
	SendQWAppointChatRoomAdminRequest(chatRoomName string, toUserName []string) (*wechat.TransferChatRoomOwnerResponse, error)
	//移除企业群管理员
	SendQWDelChatRoomAdminRequest(chatRoomName string, toUserName []string) (*wechat.TransferChatRoomOwnerResponse, error)
	//同意进企业群
	SendQWAcceptChatRoomRequest(link string, opcode uint32) (*wechat.GetA8KeyResp, error)
	//设定企业群管理审核进群
	SendQWAdminAcceptJoinChatRoomSetRequest(chatRoomName string, p int64) (*wechat.TransferChatRoomOwnerResponse, error)
	//群管理批准进企业群 1
	SendQWAdminAcceptJoinChatRoomRequest(chatRoomName, key, toUserName string, toUserNames []string) (*wechat.TransferChatRoomOwnerResponse, error)
	//修改企业群名称
	SendQWModChatRoomNameRequest(chatRoomName, name string) (*wechat.TransferChatRoomOwnerResponse, error)
	//修改成员在群中呢称
	SendQWModChatRoomMemberNickRequest(chatRoomName, name string) (*wechat.TransferChatRoomOwnerResponse, error)
	//发布企业群公告
	SendQWChatRoomAnnounceRequest(chatRoomName, Announcement string) (*wechat.TransferChatRoomOwnerResponse, error)
	//删除企业群
	SendQWDelChatRoomRequest(chatRoomName string) (*wechat.TransferChatRoomOwnerResponse, error)
	//视频号搜索
	SendGetFinderSearchRequest(Index uint32, Userver int32, UserKey string, Uuid string) (*wechat.FinderSearchResponse, error)
	//视频号个人中心
	SendFinderUserPrepareRequest(uServer int32) (*wechat.FinderUserPrepareResponse, error)
	//视频号关注
	SendFinderFollowRequest(FinderUserName string, OpType int32, RefObjectId string, Cook string, Userver int32, PosterUsername string) (*wechat.FinderFollowResponse, error)
	//视频号首页
	TargetUserPageRequest(target string, lastBuffer string) (*wechat.FinderUserPageResponse, error)
}
