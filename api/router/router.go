package router

import (
	"feiyu.com/wx/api/controller"
	"github.com/gin-gonic/gin"
)

type SetMiddleWare = func(engine *gin.Engine)

func SetUpRouter(middleware SetMiddleWare, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	//获取Gin实例
	r := gin.New()
	//设置中间
	if middleware != nil {
		middleware(r)
	}
	//加载模版下的所有文件
	r.LoadHTMLGlob("api/templates/*")

	//设置静态文件目录
	r.Static("static", "api/static")

	setApi_V1(r)

	setTemplate(r)
	return r
}

func setTemplate(engine *gin.Engine) {
}

// 设置url Version:v1
func setApi_V1(engine *gin.Engine) {

	ver := "/v1"

	//登录
	login := engine.Group(ver + "/login")
	{
		login.GET("/GetIWXConnect", controller.IWXConnectMgrApi)
		login.GET("/CheckCanSetAlias", controller.CheckCanSetAliasApi)
		login.POST("/LoginNew", controller.LoginNewApi)
		login.POST("/SmsLogin", controller.SmsLoginApi)
		login.POST("/A16Login", controller.A16LoginApi)
		login.POST("/DeviceLogin", controller.DeviceIdLoginApi)
		login.GET("/GetLoginQrCode", controller.GetLoginQrCodeApi)
		login.GET("/CheckLoginStatus", controller.CheckLoginStatusApi)
		login.GET("/GetInItStatus", controller.GetInItStatusApi)
		login.GET("/WakeUpLogin", controller.WakeUpLoginApi)
		login.GET("/GetLoginStatus", controller.GetLoginStatusApi)
		login.POST("/GetLoginQrCodeNew", controller.GetLoginQrCodeNewApi)
		login.POST("/PhoneDeviceLogin", controller.PhoneDeviceLoginApi)
		login.GET("/Get62Data", controller.Get62DataApi)
		login.POST("/WxBindOpMobileForReg", controller.WxBindOpMobileForRegApi)
		login.POST("/ExtDeviceLoginConfirmGet", controller.ExtDeviceLoginConfirmGetApi)
	}
	//equipment
	equipment := engine.Group(ver + "/equipment")
	{
		equipment.POST("/GetSafetyInfo", controller.GetSafetyInfoApi)
		equipment.POST("/DelSafeDevice", controller.DelSafeDeviceApi)
	}
	//message
	message := engine.Group(ver + "/message")
	{
		message.POST("/test", controller.TestApi)
		message.POST("/AddMessageMgr", controller.AddMessageMgrApi)
		message.POST("/SendImageMessage", controller.SendImageMessageApi)
		message.POST("/SendImageNewMessage", controller.SendImageNewMessageApi)
		message.POST("/SendTextMessage", controller.SendTextMessageApi)
		message.POST("/SendShareCard", controller.SendShareCardApi)
		message.POST("/ForwardImageMessage", controller.ForwardImageMessageApi)
		message.POST("/ForwardVideoMessage", controller.ForwardVideoMessageApi)
		message.POST("/SendEmojiMessage", controller.SendEmojiMessageApi)
		message.POST("/ForwardEmoji", controller.ForwardEmojiApi)
		message.POST("/SendAppMessage", controller.SendAppMessageApi)
		message.POST("/RevokeMsg", controller.RevokeMsgApi)
		message.POST("/SendVoice", controller.UploadVoiceRequestApi)
		message.POST("/CdnUploadVideo", controller.CdnUploadVideoRequestApi)
		message.POST("/SendCdnDownload", controller.SendCdnDownloadApi)
		message.POST("/GetMsgBigImg", controller.GetMsgBigImgApi)
		message.POST("/NewSyncHistoryMessage", controller.NewSyncHistoryMessageApi)
		message.POST("/GetMsgVoice", controller.GetMsgVoiceApi)
		message.POST("/GroupMassMsgText", controller.GroupMassMsgTextApi)
		message.POST("/GroupMassMsgImage", controller.GroupMassMsgImageApi)
	}
	// sns
	sns := engine.Group(ver + "/sns")
	{
		sns.POST("/DownloadMedia", controller.DownloadMediaApi)
		sns.POST("/SetFriendCircleDays", controller.SetFriendCircleDaysApi)
		sns.POST("/SendFriendCircle", controller.SendFriendCircleApi)
		sns.POST("/SendFriendCircleByXMl", controller.SendFriendCircleByXMlApi)
		sns.POST("/UploadFriendCircleImage", controller.UploadFriendCircleImageApi)
		sns.POST("/SendSnsComment", controller.SendSnsCommentRequestApi)
		sns.POST("/SendSnsObjectOp", controller.SendSnsObjectOpRequestApi)
		sns.POST("/SendSnsTimeLine", controller.SendSnsTimeLineRequestApi)
		sns.POST("/SendSnsUserPage", controller.SendSnsUserPageRequestApi)
		sns.POST("/SendSnsObjectDetailById", controller.SendSnsObjectDetailByIdApi)
		sns.POST("/SetBackgroundImage", controller.SetBackgroundImageApi)
		sns.POST("/SendFavItemCircle", controller.SendFavItemCircleApi)
		sns.POST("/SendOneIdCircle", controller.SendOneIdCircleApi)
		sns.POST("/GetCollectCircle", controller.GetCollectCircleApi)
	}

	//group
	group := engine.Group(ver + "/group")
	{
		group.POST("/SetChatroomAnnouncement", controller.SetChatroomAnnouncementApi)
		group.POST("/GetChatroomMemberDetail", controller.GetChatroomMemberDetailApi)
		group.POST("/QuitChatroom", controller.GetQuitChatroomApi)
		group.POST("/CreateChatRoom", controller.CreateChatRoomApi)
		group.POST("/InviteChatroomMembers", controller.InviteChatroomMembersApi)
		group.POST("/AddChatRoomMembers", controller.AddChatRoomMembersApi)
		group.POST("/SendDelDelChatRoomMember", controller.SendDelDelChatRoomMemberApi)
		group.POST("/ScanIntoUrlGroup", controller.ScanIntoUrlGroupApi)
		group.POST("/SendTransferGroupOwner", controller.SendTransferGroupOwnerApi)
		group.POST("/SetGetChatRoomInfoDetail", controller.SetGetChatRoomInfoDetailApi)
		group.POST("/GetChatRoomInfo", controller.GetChatRoomInfoApi)
		group.POST("/MoveToContract", controller.MoveToContractApi)
		group.POST("/SetChatroomAccessVerify", controller.SetChatroomAccessVerifyApi)
		group.POST("/AddChatroomAdmin", controller.AddChatroomAdminApi)
		group.POST("/DelChatroomAdmin", controller.DelChatroomAdminApi)
		group.POST("/SetChatroomName", controller.SetChatroomNameApi)
		group.POST("/SendPat", controller.SendPatApi)
		group.GET("/GroupList", controller.GroupListApi)

	}

	user := engine.Group(ver + "/user")
	{ //
		user.GET("/LogOut", controller.LogOutRequestApi)
		user.POST("/GetContactList", controller.GetContactListApi)
		//user.GET("/GetFriendList", controller.GetFriendListApi)
		//user.GET("/GroupList", controller.GetGroupListApi)
		//user.GET("/GetGHList", controller.GetGHListApi)
		user.POST("/GetRedisSyncMsg", controller.GetRedisSyncMsgApi)
		user.GET("/GetMFriend", controller.GetMFriendApi)
		user.POST("/GetContactDetailsList", controller.GetContactContactApi)
		user.POST("/GetFriendRelation", controller.GetFriendRelationApi)
		user.POST("/GetFriendRelations", controller.GetFriendRelationsApi)
		user.POST("/UploadMContact", controller.UploadMContactApi)
		user.GET("/GetOnlineInfo", controller.GetOnlineInfoApi)
		user.GET("/GetProfile", controller.GetProfileApi)
		user.POST("/DelContact", controller.DelContactApi)
		user.POST("/ModifyUserInfo", controller.ModifyUserInfoRequestApi)
		user.POST("/UpdateNickName", controller.UpdateNickNameApi)
		user.POST("/SetNickName", controller.SetNickNameApi)
		user.POST("/SetSignature", controller.SetSignatureApi)
		user.POST("/SetSexDq", controller.SetSexApi)
		user.POST("/ChangePwd", controller.ChangePwdRequestRequestApi)
		user.POST("/UploadHeadImage", controller.UploadHeadImageApi)
		user.POST("/UpdateAutoPass", controller.UpdateAutopassApi)
		user.POST("/ModifyRemark", controller.SendModifyRemarkRequestApi)
		user.POST("/SetWechat", controller.SetWechatApi)
		user.POST("/UpdateStepNumber", controller.UpdateStepNumberApi)
		user.POST("/GetUserRankLikeCount", controller.GetUserRankLikeCountApi)
		user.POST("/SetFunctionSwitch", controller.SetFunctionSwitchApi)
		user.POST("/SetSendPat", controller.SetSendPatApi)
		user.POST("/BindingMobile", controller.BindingMobileApi)
		user.POST("/SendVerifyMobile", controller.SendVerifyMobileApi)
	}

	applet := engine.Group(ver + "/applet")
	{
		applet.POST("/FollowGH", controller.FollowGHApi)
		applet.POST("/GetA8Key", controller.GetA8KeyApi)
		applet.POST("/JsLogin", controller.JSLoginApi)
		applet.POST("/JSOperateWxData", controller.JSOperateWxDataApi)
		applet.POST("/SdkOauthAuthorize", controller.SdkOauthAuthorizeApi)
		applet.POST("/QRConnectAuthorize", controller.QRConnectAuthorizeApi)
		applet.POST("/QRConnectAuthorizeConfirm", controller.QRConnectAuthorizeConfirmApi)
		applet.POST("/GetMpA8Key", controller.GetMpA8KeyApi)
		applet.POST("/AuthMpLogin", controller.AuthMpLoginApi)
	}

	other := engine.Group(ver + "/other")
	{
		other.POST("/GetQrCode", controller.GetQrCodeApi)
		other.POST("/GetPeopleNearby", controller.GetPeopleNearbyApi)
	}
	favor := engine.Group(ver + "/favor")
	{
		favor.GET("/FavSync", controller.FavSyncApi)
		favor.POST("/GetFavList", controller.GetFavListApi)
		favor.POST("/GetFavItemId", controller.BatchGetFavItemApi)
		favor.POST("/ShareFav", controller.ShareFavServiceApi)
		favor.POST("/CheckFavCdn", controller.CheckFavCdnServiceApi)
		favor.POST("/BatchDelFavItem", controller.BatchDelFavItemApi)
	}
	label := engine.Group(ver + "/label")
	{
		label.GET("/GetContactLabelList", controller.GetContactLabelListApi)
		label.POST("/AddContactLabel", controller.AddContactLabelRequestApi)
		label.POST("/DelContactLabel", controller.DelContactLabelRequestApi)
		label.POST("/ModifyLabel", controller.ModifyLabelRequestApi)
		//label.POST("/GetWXFriendListByLabel", controller.GetWXFriendListByLabelIDApi)
	}
	friend := engine.Group(ver + "/friend")
	{
		friend.POST("/SearchContact", controller.SearchContactRequestApi)
		friend.POST("/VerifyUser", controller.VerifyUserRequestApi)
		friend.POST("/AgreeAdd", controller.AgreeAddApi)
	}

	pay := engine.Group(ver + "/pay")
	{
		pay.POST("/GetBandCardList", controller.GetBandCardListApi)
		pay.POST("/GeneratePayQCode", controller.GeneratePayQCodeApi)
		pay.POST("/Collectmoney", controller.CollectmoneyApi)
		pay.POST("/WXCreateRedPacket", controller.WXCreateRedPacketApi)
		pay.POST("/OpenRedEnvelopes", controller.OpenRedEnvelopesApi)
		pay.POST("/GetRedEnvelopesDetail", controller.QueryRedEnvelopesDetailApi)
		pay.POST("/GetRedPacketList", controller.GetRedPacketListApi)
		pay.POST("/CreatePreTransfer", controller.CreatePreTransferApi)
		pay.POST("/ConfirmPreTransfer", controller.ConfirmPreTransferApi)
	}
	//视频号
	finder := engine.Group(ver + "/Finder")
	{
		finder.POST("/FinderSearch", controller.GetFinderSearchApi)
		finder.POST("/FinderUserPrepare", controller.FinderUserPrepareApi)
		finder.POST("/FinderFollow", controller.FinderFollowApi)
		finder.POST("/TargetUserPage", controller.TargetUserPageApi)
	}
	//公众号
	gh := engine.Group(ver + "/Gh")
	{
		gh.GET("/Search", controller.SendSearchApi)
		gh.POST("/Follower", controller.SendFollowApi)
		gh.POST("/ClickMenu", controller.SendClickMenuApi)
		gh.POST("/ReadArticle", controller.SendReadArticleApi)
		gh.POST("/LikeArticle", controller.SendLikeArticleApi)
	}
	qy := engine.Group(ver + "/qy")
	{
		qy.POST("/QWSearchContact", controller.QWSearchContactApi)
		qy.POST("/QWApplyAddContact", controller.QWApplyAddContactApi)
		qy.POST("QWAddContact", controller.QWAddContactApi)
		qy.POST("/QWContact", controller.QWContactApi)
		qy.POST("/QWSyncContact", controller.QWSyncContactApi)
		qy.POST("/QWRemark", controller.QWRemarkApi)
		qy.POST("/QWCreateChatRoom", controller.QWCreateChatRoomApi)
		qy.POST("/QWSyncChatRoom", controller.QWSyncChatRoomApi)
		qy.POST("/QWChatRoomTransferOwner", controller.QWChatRoomTransferOwnerApi)
		qy.POST("/QWAddChatRoomMember", controller.QWAddChatRoomMemberApi)
		qy.POST("/QWInviteChatRoomMember", controller.QWInviteChatRoomMemberApi)
		qy.POST("/QWDelChatRoomMember", controller.QWDelChatRoomMemberApi)
		qy.POST("/QWGetChatRoomMember", controller.QWGetChatRoomMemberApi)
		qy.POST("/QWGetChatroomInfo", controller.QWGetChatroomInfoApi)
		qy.POST("/QWGetChatRoomQR", controller.QWGetChatRoomQRApi)
		qy.POST("/QWAppointChatRoomAdmin", controller.QWAppointChatRoomAdminApi)
		qy.POST("/QWDelChatRoomAdmin", controller.QWDelChatRoomAdminApi)
		qy.POST("/QWAcceptChatRoom", controller.QWAcceptChatRoomRequestApi)
		qy.POST("/QWAdminAcceptJoinChatRoomSet", controller.QWAdminAcceptJoinChatRoomSetApi)
		qy.POST("/QWModChatRoomName", controller.QWModChatRoomNameApi)
		qy.POST("/QWModChatRoomMemberNick", controller.QWModChatRoomMemberNickApi)
		qy.POST("/QWChatRoomAnnounce", controller.QWChatRoomAnnounceApi)
		qy.POST("/QWDelChatRoom", controller.QWDelChatRoomApi)
	}

}
