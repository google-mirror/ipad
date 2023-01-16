package baseinfo

type GetA8KeyType uint32

const (
	GetA8Key      GetA8KeyType = 1
	ThrIdGetA8Key GetA8KeyType = 2
)

const (
	// ERR_SERVER_FILE_EXPIRED int32 = -5103059
	// MM_ERR_FORCE_QUIT int32 = -999999
	// MM_ERR_CLIENT int32 = -800000
	// MM_ERR_CHATROOM_PARTIAL_INVITE int32 = -2013
	// MM_ERR_CHATROOM_NEED_INVITE int32 = -2012
	// MM_ERR_CONNECT_INFO_URL_INVALID int32 = -2011
	// MM_ERR_CLIDB_ENCRYPT_KEYINFO_INVALID int32 = -2010
	// MM_ERR_LOGIN_URL_DEVICE_UNSAFE int32 = -2009
	// MM_ERR_COOKIE_KICK int32 = -2008
	// MM_ERR_LOGIN_QRCODE_UUID_EXPIRED int32 = -2007
	// MM_ERR_KEYBUF_INVALID int32 = -2006
	// MM_ERR_FORCE_REDIRECT int32 = -2005
	// MM_ERR_QRCODEVERIFY_BANBYEXPOSE int32 = -2004
	// MM_ERR_SHAKEBANBYEXPOSE int32 = -2003
	// MM_ERR_BOTTLEBANBYEXPOSE int32 = -2002
	// MM_ERR_LBSBANBYEXPOSE int32 = -2001
	// MM_ERR_LBSDATANOTFOUND int32 = -2000
	// MM_ERR_IMG_READ int32 = -1005
	// MM_ERR_FACING_CREATECHATROOM_RETRY int32 = -432
	// MM_ERR_RADAR_PASSWORD_SIMPLE int32 = -431
	// MM_ERR_REVOKEMSG_TIMEOUT int32 = -430
	// MM_ERR_FAV_ALREADY int32 = -400
	// MM_ERR_FILE_EXPIRED int32 = -352
	// MM_ERR_USER_NOT_VERIFYUSER int32 = -302

	//MMErrIdcRedirect MMErrIdcRedirect
	MMErrIdcRedirect int32 = -301

	//MMErrChangeKey
	MMErrChangeKey int32 = -305

	//MMErrDropped 用户主动退出或者服务器T下线
	MMErrDropped int32 = -2023

	// MM_ERR_REG_BUT_LOGIN int32 = -212
	// MM_ERR_UNBIND_MAIN_ACCT int32 = -206
	// MM_ERR_QQ_OK_NEED_MOBILE int32 = -205
	// MM_ERR_OTHER_MAIN_ACCT int32 = -204
	// MM_ERR_NODATA int32 = -203
	// MM_ERR_UNBIND_MOBILE_NEED_QQPWD int32 = -202
	// MM_ERR_QQ_BAN int32 = -201
	// MM_ERR_ACCOUNT_BAN int32 = -200
	// MM_ERR_QA_RELATION int32 = -153
	// MM_ERR_NO_QUESTION int32 = -152
	// MM_ERR_QUESTION_COUNT int32 = -151
	// MM_ERR_ANSWER_COUNT int32 = -150
	// MM_ERR_EMAIL_FORMAT int32 = -111
	// MM_ERR_BLOCK_BY_SPAM int32 = -106
	// MM_ERR_CERT_EXPIRED 类型：证书已过期，需要切换密钥
	MM_ERR_CERT_EXPIRED int32 = -102
	// MM_ERR_NO_RETRY int32 = -101
	// MM_ERR_AUTH_ANOTHERPLACE int32 = -100
	// MM_ERR_USER_NOT_SUPPORT int32 = -94
	// MM_ERR_SHAKE_TRAN_IMG_OTHER int32 = -93
	// MM_ERR_SHAKE_TRAN_IMG_CONTINUE int32 = -92
	// MM_ERR_SHAKE_TRAN_IMG_NOTFOUND int32 = -91
	// MM_ERR_SHAKE_TRAN_IMG_CANCEL int32 = -90
	// MM_ERR_BIZ_FANS_LIMITED int32 = -87
	// MM_ERR_BIND_EMAIL_SAME_AS_QMAIL int32 = -86
	// MM_ERR_BINDED_BY_OTHER int32 = -85
	// MM_ERR_HAS_BINDED int32 = -84
	// MM_ERR_HAS_UNBINDED int32 = -83
	// MM_ERR_ONE_BINDTYPE_LEFT int32 = -82
	// MM_ERR_NOTBINDQQ int32 = -81
	// MM_ERR_WEIBO_PUSH_TRANS int32 = -80
	// MM_ERR_NEW_USER int32 = -79
	// MM_ERR_SVR_MOBILE_FORMAT int32 = -78
	// MM_ERR_WRONG_SESSION_KEY int32 = -77
	// MM_ERR_UUID_BINDED int32 = -76
	// MM_ERR_ALPHA_FORBIDDEN int32 = -75
	// MM_ERR_MOBILE_NEEDADJUST int32 = -74
	// MM_ERR_TRYQQPWD int32 = -73
	// MM_ERR_NICEQQ_EXPIRED int32 = -72
	// MM_ERR_TOLIST_LIMITED int32 = -71
	// MM_ERR_GETMFRIEND_NOT_READY int32 = -70
	// MM_ERR_BIGBIZ_AUTH int32 = -69
	// MM_FACEBOOK_ACCESSTOKEN_UNVALID int32 = -68
	// MM_ERR_HAVE_BIND_FACEBOOK int32 = -67
	// MM_ERR_IS_NOT_OWNER int32 = -66
	// MM_ERR_UNBIND_REGBYMOBILE int32 = -65
	// MM_ERR_PARSE_MAIL int32 = -64
	// MM_ERR_GMAIL_IMAP int32 = -63
	// MM_ERR_GMAIL_WEBLOGIN int32 = -62
	// MM_ERR_GMAIL_ONLINELIMITE int32 = -61
	// MM_ERR_GMAIL_PWD int32 = -60
	// MM_ERR_UNSUPPORT_COUNTRY int32 = -59
	// MM_ERR_PICKBOTTLE_NOBOTTLE int32 = -58
	// MM_ERR_SEND_VERIFYCODE int32 = -57
	// MM_ERR_NO_BOTTLECOUNT int32 = -56
	// MM_ERR_NO_HDHEADIMG int32 = -55
	// MM_ERR_INVALID_HDHEADIMG_REQ_TOTAL_LEN int32 = -54
	// MM_ERR_HAS_NO_HEADIMG int32 = -53
	// MM_ERR_INVALID_GROUPCARD_CONTACT int32 = -52
	// MM_ERR_VERIFYCODE_NOTEXIST int32 = -51
	// MM_ERR_BINDUIN_BINDED int32 = -50
	// MM_ERR_NEED_QQPWD int32 = -49
	// MM_ERR_TICKET_NOTFOUND int32 = -48
	// MM_ERR_TICKET_UNMATCH int32 = -47
	// MM_ERR_NOTQQCONTACT int32 = -46
	// MM_ERR_BATCHGETCONTACTPROFILE_MODE int32 = -45
	// MM_ERR_NEED_VERIFY_USER int32 = -44
	// MM_ERR_USER_BIND_MOBILE int32 = -43
	// MM_ERR_USER_MOBILE_UNMATCH int32 = -42
	// MM_ERR_MOBILE_FORMAT int32 = -41
	// MM_ERR_UNMATCH_MOBILE int32 = -40
	// MM_ERR_MOBILE_NULL int32 = -39
	// MM_ERR_INVALID_UPLOADMCONTACT_OPMODE int32 = -38
	// MM_ERR_INVALID_BIND_OPMODE int32 = -37
	// MM_ERR_MOBILE_UNBINDED int32 = -36
	// MM_ERR_MOBILE_BINDED int32 = -35
	// MM_ERR_FREQ_LIMITED int32 = -34
	// MM_ERR_VERIFYCODE_TIMEOUT int32 = -33
	// MM_ERR_VERIFYCODE_UNMATCH int32 = -32
	// MM_ERR_NEEDSECONDPWD int32 = -31
	// MM_ERR_NEEDREG int32 = -30
	// MM_ERR_OIDBTIMEOUT int32 = -29
	// MM_ERR_BADEMAIL int32 = -28
	// MM_ERR_DOMAINDISABLE int32 = -27
	// MM_ERR_DOMAINMAXLIMITED int32 = -26
	// MM_ERR_DOMAINVERIFIED int32 = -25
	// MM_ERR_SPAM int32 = -24
	// MM_ERR_MEMBER_TOOMUCH int32 = -23
	// MM_ERR_BLACKLIST int32 = -22
	// MM_ERR_NOTCHATROOMCONTACT int32 = -21
	// MM_ERR_NOTMICROBLOGCONTACT int32 = -20
	// MM_ERR_NOTOPENPRIVATEMSG int32 = -19
	// MM_ERR_NOUPDATEINFO int32 = -18
	// MM_ERR_RECOMMENDEDUPDATE int32 = -17
	// MM_ERR_CRITICALUPDATE int32 = -16
	// MM_ERR_NICKNAMEINVALID int32 = -15
	// MM_ERR_USERNAMEINVALID int32 = -14

	//MMErrSessionTimeOut session超时，可能是由于手机端主动终止授权导致的
	MMErrSessionTimeOut int32 = -13

	// MM_ERR_UINEXIST int32 = -12
	// MM_ERR_NICKRESERVED int32 = -11
	// MM_ERR_USERRESERVED int32 = -10
	// MM_ERR_EMAILNOTVERIFY int32 = -9
	// MM_ERR_EMAILEXIST int32 = -8
	// MM_ERR_USEREXIST int32 = -7
	// MM_ERR_NEED_VERIFY int32 = -6
	// MM_ERR_ACCESS_DENIED int32 = -5
	// MM_ERR_NOUSER int32 = -4
	// MM_ERR_PASSWORD int32 = -3
	// MM_ERR_ARG int32 = -2
	// MM_ERR_SYS int32 = -1

	//CFatalError 致命错误
	CFatalError int32 = -100000000
	//CHttpError http错误
	CHttpError int32 = -100000001
	//CProtoError proto错误
	CProtoError int32 = -100000002
	//CParseResponseDataError 数据解析错误
	CParseResponseDataError int32 = -10000003

	//MMOk MMOk
	MMOk int32 = 0
	// MM_BOTTLE_ERR_UNKNOWNTYPE int32 = 15
	// MM_BOTTLE_COUNT_ERR int32 = 16
	// MM_BOTTLE_NOTEXIT int32 = 17
	// MM_BOTTLE_UINNOTMATCH int32 = 18
	// MM_BOTTLE_PICKCOUNTINVALID int32 = 19
	// MMSNS_RET_SPAM int32 = 201
	// MMSNS_RET_BAN int32 = 202
	// MMSNS_RET_PRIVACY int32 = 203
	// MMSNS_RET_COMMENT_HAVE_LIKE int32 = 204
	// MMSNS_RET_COMMENT_NOT_ALLOW int32 = 205
	// MMSNS_RET_CLIENTID_EXIST int32 = 206
	// MMSNS_RET_ISALL int32 = 207
	// MMSNS_RET_COMMENT_PRIVACY int32 = 208
	// MM_ERR_SHORTVIDEO_CANCEL int32 = 1000000

	//-----------------------------

	//CmdInvalid CmdInvalid
	CmdInvalid uint32 = 0
	//CmdIDModUserInfo 用户详情
	CmdIDModUserInfo uint32 = 1
	//CmdIDModContact CmdIdModContact
	CmdIDModContact uint32 = 2
	//CmdIDDelContact CmdIdDelContact
	CmdIDDelContact uint32 = 4
	//CmdIDAddMsg CmdIdAddMsg
	CmdIDAddMsg uint32 = 5
	//CmdIDModMsgStatus CmdIdModMsgStatus
	CmdIDModMsgStatus uint32 = 6
	//CmdIDDelChatContact CmdIdDelChatContact
	CmdIDDelChatContact uint32 = 7
	//CmdIDDelContactMsg CmdIdDelContactMsg
	CmdIDDelContactMsg uint32 = 8
	//CmdIDDelMsg CmdIdDelMsg
	CmdIDDelMsg uint32 = 9
	//CmdIDReport CmdIdReport
	CmdIDReport uint32 = 10
	//CmdIDOpenQQMicroBlog CmdIdOpenQQMicroBlog
	CmdIDOpenQQMicroBlog uint32 = 11
	//CmdIDCloseMicroBlog CmdIdCloseMicroBlog
	CmdIDCloseMicroBlog uint32 = 12
	//CmdIDModMicroBlog CmdIdModMicroBlog
	CmdIDModMicroBlog uint32 = 13
	//CmdIDModNotifyStatus CmdIdModNotifyStatus
	CmdIDModNotifyStatus uint32 = 14
	//CmdIDModChatRoomMember CmdIdModChatRoomMember
	CmdIDModChatRoomMember uint32 = 15
	//CmdIDQuitChatRoom CmdIdQuitChatRoom
	CmdIDQuitChatRoom uint32 = 16
	//CmdIDModContactDomainEmail CmdIdModContactDomainEmail
	CmdIDModContactDomainEmail uint32 = 17
	//CmdIDModUserDomainEmail CmdIdModUserDomainEmail
	CmdIDModUserDomainEmail uint32 = 18
	//CmdIDDelUserDomainEmail CmdIdDelUserDomainEmail
	CmdIDDelUserDomainEmail uint32 = 19
	//CmdIDModChatRoomNotify CmdIdModChatRoomNotify
	CmdIDModChatRoomNotify uint32 = 20
	//CmdIDPossibleFriend CmdIdPossibleFriend
	CmdIDPossibleFriend uint32 = 21
	//CmdIDInviteFriendOpen CmdIdInviteFriendOpen
	CmdIDInviteFriendOpen uint32 = 22
	//CmdIDFunctionSwitch CmdIdFunctionSwitch
	CmdIDFunctionSwitch uint32 = 23
	//CmdIDModQContact CmdIdModQContact
	CmdIDModQContact uint32 = 24
	//CmdIDModTContact CmdIdModTContact
	CmdIDModTContact uint32 = 25
	//CmdIDPsmStat CmdIdPsmStat
	CmdIDPsmStat uint32 = 26
	//CmdIDModChatRoomTopic CmdIdModChatRoomTopic
	CmdIDModChatRoomTopic uint32 = 27
	// MM_SYNCCMD_UPDATESTAT uint32 = 30
	// MM_SYNCCMD_MODDISTURBSETTING uint32 = 31
	// MM_SYNCCMD_DELETEBOTTLE uint32 = 32
	// MM_SYNCCMD_MODBOTTLECONTACT uint32 = 33
	// MM_SYNCCMD_DELBOTTLECONTACT uint32 = 34

	//CmdIDModUserImg 用户图像？
	CmdIDModUserImg uint32 = 35
	// MM_SYNCCMD_MODUSERIMG uint32 = 35
	// MM_SYNCCMD_KVSTAT uint32 = 36
	// NN_SYNCCMD_THEMESTAT uint32 = 37

	//CmdIDUserInfoExt 用户扩展数据
	CmdIDUserInfoExt uint32 = 44
	// MM_SYNCCMD_USERINFOEXT uint32 = 44

	// MMSnsSyncCmdObject 朋友圈同步到的对象
	MMSnsSyncCmdObject uint32 = 45

	// MM_SNS_SYNCCMD_ACTION uint32 = 46
	// MM_SYNCCMD_BRAND_SETTING uint32 = 47
	// MM_SYNCCMD_MODCHATROOMMEMBERDISPLAYNAME uint32 = 48
	// MM_SYNCCMD_MODCHATROOMMEMBERFLAG uint32 = 49
	// MM_SYNCCMD_WEBWXFUNCTIONSWITCH uint32 = 50
	// MM_SYNCCMD_MODSNSUSERINFO uint32 = 51
	// MM_SYNCCMD_MODSNSBLACKLIST uint32 = 52
	// MM_SYNCCMD_NEWDELMSG uint32 = 53
	// MM_SYNCCMD_MODDESCRIPTION uint32 = 54
	// MM_SYNCCMD_KVCMD uint32 = 55
	// MM_SYNCCMD_DELETE_SNS_OLDGROUP uint32 = 56

	// CmdIdMax uint32 = 201
	// MM_GAME_SYNCCMD_ADDMSG uint32 = 201

	//MMLoginUnknow MMLoginUnknow
	MMLoginUnknow int32 = -2
	//MMLoginError MMLoginError
	MMLoginError int32 = -1
	//MMLoginSuccess MMLoginSuccess
	MMLoginSuccess int32 = 0
	//MMLoginRedirect MMLoginRedirect
	MMLoginRedirect int32 = 1

	//MMAddFriendNoVerify 好友校验类型-不需要验证
	MMAddFriendNoVerify uint32 = 1
	//MMAddFriendWithVerify 好友校验类型-需要验证
	MMAddFriendWithVerify uint32 = 2
	//MMAddFriendAccept 通过好友验证
	MMAddFriendAccept uint32 = 3

	// MMAddFiendFromQQ 好友来源-QQ
	MMAddFiendFromQQ byte = 1
	// MMAddFiendFromMail 好友来源-邮箱
	MMAddFiendFromMail byte = 2
	// MMAddFiendFromWxName 好友来源-微信号
	MMAddFiendFromWxName byte = 3
	// MMAddFiendFromAddressBook 好友来源-通讯录
	MMAddFiendFromAddressBook byte = 13
	// MMAddFiendFromChatRoom 好友来源-群
	MMAddFiendFromChatRoom byte = 14
	// MMAddFiendFromPhone 好友来源-手机号
	MMAddFiendFromPhone byte = 15
	// MMAddFiendFromNear 好友来源-附近的人
	MMAddFiendFromNear byte = 18
	// MMAddFiendFromBottle 好友来源-漂流瓶
	MMAddFiendFromBottle byte = 25
	// MMAddFiendFromShake 好友来源-摇一摇
	MMAddFiendFromShake byte = 29
	// MMAddFiendFromQrcode 好友来源-二维码
	MMAddFiendFromQrcode byte = 30

	// MMVerifyUserErrPrivate 对方为私有设置 添加失败
	MMVerifyUserErrPrivate int32 = -2
	// MMVerifyUserErrNeedVerify 添加好友需要发送验证信息
	MMVerifyUserErrNeedVerify int32 = -44

	// ModUserSexMale 用户性别-男
	ModUserSexMale uint32 = 1
	// ModUserSexFemale 用户性别-女
	ModUserSexFemale uint32 = 2

	// MMStatusNotifyMarkChatRead 标记某个联系人或群聊的消息已读
	MMStatusNotifyMarkChatRead uint32 = 1
	// MMStatusNotifyEnterChat 进入聊天房间(联系人或群聊)
	MMStatusNotifyEnterChat uint32 = 2
	// MMStatusNotifyGetChatList 第一次登陆时获取 聊天项列表
	MMStatusNotifyGetChatList uint32 = 3
	// MMStatusNotifyGetAllChat 获取所有聊天项
	MMStatusNotifyGetAllChat uint32 = 4
	// MMStatusNotifyQuitChat 关闭微信
	MMStatusNotifyQuitChat uint32 = 5
	// MMStatusNotifyWechatResume 微信从后台 切换到 最前面
	MMStatusNotifyWechatResume uint32 = 7
	// MMStatusNotifyWechatToBackground 切换到后台
	MMStatusNotifyWechatToBackground uint32 = 8
	// MMStatusNotifyMark 标记 语音消息，朋友圈等状态
	MMStatusNotifyMark uint32 = 9

	// MMSnsOpCodeDelete 删除朋友圈
	MMSnsOpCodeDelete uint32 = 1
	// MMSnsOpCodeSetPrivate 设置朋友圈为私密文字(仅自己可见)
	MMSnsOpCodeSetPrivate uint32 = 2
	// MMSnsOpCodeSetPublic 设置朋友圈为公开信息(所有人可见)
	MMSnsOpCodeSetPublic uint32 = 3
	// MMSnsOpCodeDeleteComment 删除评论
	MMSnsOpCodeDeleteComment uint32 = 4
	// MMSnsOpCodeUnLike 取消点赞
	MMSnsOpCodeUnLike uint32 = 5

	// MMSnsCommentTypeLike 点赞
	MMSnsCommentTypeLike uint32 = 1
	// MMSnsCommentTypeComment 发表评论
	MMSnsCommentTypeComment uint32 = 2

	// MMSnsPrivacyPublic 朋友圈状态:公开
	MMSnsPrivacyPublic uint32 = 0
	// MMSnsPrivacyPrivate 朋友圈状态:不公开可指定好友可见
	MMSnsPrivacyPrivate uint32 = 1

	// MMSNSContentStyleImgAndText 图文朋友圈
	MMSNSContentStyleImgAndText uint32 = 1
	// MMSNSContentStyleText 文本朋友圈
	MMSNSContentStyleText uint32 = 2
	// MMSNSContentStyleRefer 引用
	MMSNSContentStyleRefer uint32 = 3
	// MMSNSContentStyleVideo 微信小视频
	MMSNSContentStyleVideo uint32 = 15

	// MMSNSMediaTypeImage 类型为图片
	MMSNSMediaTypeImage uint32 = 2
	// MMSNSMediaTypeVideo 类型为视频
	MMSNSMediaTypeVideo uint32 = 6

	// MMCdnDownMediaTypeImage Cdn下载类型：图片
	MMCdnDownMediaTypeImage uint32 = 2
	// MMCdnDownMediaTypeVedioImage Cdn下载类型：视频封面图片
	MMCdnDownMediaTypeVedioImage uint32 = 3
	// MMCdnDownMediaTypeVedio Cdn下载类型：视频
	MMCdnDownMediaTypeVedio uint32 = 4

	// MMHeadDeviceTypeIpadUniversal IpadUniversal类型
	MMHeadDeviceTypeIpadUniversal byte = 0x0d //13
	// MMHeadDeviceTypeIpadOthers 其它iPad类型
	MMHeadDeviceTypeIpadOthers byte = 0x01

	// MMAppRunStateNormal App在前面正常运行
	MMAppRunStateNormal byte = 0xff
	// MMAppRunStateBackgroundRun App在后台运行
	MMAppRunStateBackgroundRun byte = 0xfe
	// MMAppRunStateBackgroundWillSuspend App在后台将要暂停运行
	MMAppRunStateBackgroundWillSuspend byte = 0xf6

	// MMPackDataTypeCompressed 打包数据类型：压缩
	MMPackDataTypeCompressed byte = 1
	// MMPackDataTypeUnCompressed 打包类型：未压缩
	MMPackDataTypeUnCompressed byte = 2

	// MMSyncSceneTypeApnsNotify 同步场景：收到苹果推送/长链接推送同步通知(长链接)
	MMSyncSceneTypeApnsNotify uint32 = 1
	// MMSyncSceneTypeOnTokenLogin 同步场景：二次登陆
	MMSyncSceneTypeOnTokenLogin uint32 = 2
	// MMSyncSceneTypeBackGroundToForeGround 同步场景：后台切换到前台
	MMSyncSceneTypeBackGroundToForeGround uint32 = 3
	// MMSyncSceneTypeProcessStart 同步场景：进程开始
	MMSyncSceneTypeProcessStart uint32 = 4
	// MMSyncSceneTypeNeed 同步场景：需要同步
	MMSyncSceneTypeNeed uint32 = 7
	// MMSyncSceneTypeAfterManualAuthNotify 同步场景：扫码登陆后收到的同步通知
	MMSyncSceneTypeAfterManualAuthNotify uint32 = 10
	// MMSyncSceneTypeVOIPPushAwakeOld 同步场景：语音推送唤醒(老版本)
	MMSyncSceneTypeVOIPPushAwakeOld uint32 = 13
	// MMSyncSceneTypeVOIPPushAwake 同步场景：语音推送唤醒
	MMSyncSceneTypeVOIPPushAwake uint32 = 15
	// MMSyncSceneTypeSlientPushAwake 同步场景：静默推送唤醒
	MMSyncSceneTypeSlientPushAwake uint32 = 16

	// MMSyncMsgDigestTypeLongLink 长链接同步
	MMSyncMsgDigestTypeLongLink uint32 = 0
	// MMSyncMsgDigestTypeShortLink 短链接同步
	MMSyncMsgDigestTypeShortLink uint32 = 1

	// MMUUIDTypeUnArchive 设备UUID类型：未存档的，新的设备
	MMUUIDTypeUnArchive int = 1
	// MMUUIDTypeArchive 设备UUID类型：存档过的，使用过微信的设备
	MMUUIDTypeArchive int = 2
)

const (
	// MMRequestTypeForwardCdnImage 转发cdn图
	MMRequestTypeForwardCdnImage uint32 = 110
	// MMRequestTypeForwardCdnVideo
	MMRequestTypeForwardCdnVideo uint32 = 110
	// MMRequestTypeSearchContact 类型：搜索联系人
	MMRequestTypeSearchContact uint32 = 106
	// MMRequestTypeUploadMsgImg 类型：发送图片
	MMRequestTypeUploadMsgImg uint32 = 110
	// MMRequestTypeCreateChatRoom 类型：创建群聊
	MMRequestTypeCreateChatRoom uint32 = 119
	// MMRequestTypeAddChatRoomMember 类型：邀请好友进群
	MMRequestTypeAddChatRoomMember uint32 = 120
	// MMRequestTypeUploadVoice 类型：发送语音
	MMRequestTypeUploadVoice    uint32 = 127
	MMRequestTypeUploadVoiceNew uint32 = 329
	// MMRequestTypeDownloadVoice 类型：下载语音
	MMRequestTypeDownloadVoice uint32 = 128
	// MMRequestTypeUploadMContact 类型：上传通讯录
	MMRequestTypeUploadMContact uint32 = 133
	// MMRequestTypeVerifyUser 类型：添加/验证好友
	MMRequestTypeVerifyUser uint32 = 137
	// MMRequestTypeNewSync 类型：同步消息
	MMRequestTypeNewSync uint32 = 138
	// MMRequestTypeNewInit  类型：首次登录初始化
	MMRequestTypeNewInit uint32 = 139
	// MMRequestTypeGetMFriend 类型获取手机通讯录好友
	MMRequestTypeGetMFriend uint32 = 142
	// MMRequestTypeUploadHDHeadImg 类型：上传头像
	MMRequestTypeUploadHDHeadImg uint32 = 157
	// MMRequestTypeGetQrCode 类型：获取群/个人二维码
	MMRequestTypeGetQrCode uint32 = 168
	// MMRequestTypeDelChatRoomMember 类型：删除群成员
	MMRequestTypeDelChatRoomMember uint32 = 179
	// MMRequestTypeTransferChatRoomOwnerRequest 类型：转让群
	MMRequestTypeTransferChatRoomOwnerRequest uint32 = 990
	// MMRequestTypeSendEmoji 发生表情
	MMRequestTypeSendEmoji uint32 = 175
	// MMRequestTypeGetContact 类型：获取联系人信息
	MMRequestTypeGetContact uint32 = 182
	// MMRequestTypeMMSnsPost 类型：发朋友圈
	MMRequestTypeMMSnsPost uint32 = 209
	// MMRequestTypeMMSnsObjectDetail 类型：指定朋友圈详情
	MMRequestTypeMMSnsObjectDetail uint32 = 210
	// MMRequestTypeMMSnsTimeLine 取朋友圈首页
	MMRequestTypeMMSnsTimeLine uint32 = 211
	// MMRequestTypeMMSnsUserPage 类型：获取朋友圈信息
	MMRequestTypeMMSnsUserPage uint32 = 212
	// MMRequestTypeMMSnsComment 类型：点赞/评论朋友圈
	MMRequestTypeMMSnsComment uint32 = 213
	// MMRequestTypeMMSnsSync 类型：同步朋友圈
	MMRequestTypeMMSnsSync uint32 = 214
	// MMRequestTypeMMSnsObjectOp 类型：发朋友圈操作
	MMRequestTypeMMSnsObjectOp uint32 = 218
	// MMRequestTypeSendAppMsg 类型 :发送app消息
	MMRequestTypeSendAppMsg uint32 = 222
	// MMRequestTypeGetChatRoomInfoDetail 类型：获取聊天室详情
	MMRequestTypeGetChatRoomInfoDetail uint32 = 223
	// MMRequestTypeGetA8Key 类型；授权链接
	MMRequestTypeGetA8Key      uint32 = 233
	MMRequestTypeThrIdGetA8Key uint32 = 226

	// MMRequestTypeStatusNotify 类型：发送状态
	MMRequestTypeStatusNotify uint32 = 251
	// MMRequestTypeSecManualAuth 类型：安全登陆
	MMRequestTypeSecManualAuth uint32 = 252
	// MMRequestTypeLogout 类型：退出登陆
	MMRequestTypeLogout uint32 = 282
	// MMRequestTypeGetProfile 类型：获取帐号所有信息
	MMRequestTypeGetProfile uint32 = 302
	// MMRequestTypeGetCdnDNS 类型：获取CdnDNS信息
	MMRequestTypeGetCdnDNS uint32 = 379
	// MMRequestTypeGetCert 类型：获取密钥信息
	MMRequestTypeGetCert uint32 = 381
	// MMRequestTypeVerifyPassword 类型 ：验证密码
	MMRequestTypeVerifyPassword uint32 = 384
	// MMRequestTypeSetPassword 类型 ：修改密码
	MMRequestTypeSetPassword uint32 = 383
	// MMRequestTypeFavSync 类型：同步收藏
	MMRequestTypeFavSync uint32 = 400
	// MMRequestTypeBatchGetFavItem 类型：批量获取收藏项
	MMRequestTypeBatchGetFavItem uint32 = 402
	// MMRequestTypeBatchDelFavItem 类型：删除收藏
	MMRequestTypeBatchDelFavItem uint32 = 403
	// MMRequestTypeCheckFavCdn 类型：删除收藏
	MMRequestTypeCheckFavCdn uint32 = 404
	// MMRequestTypeGetFavInfo 类型：获取收藏信息
	MMRequestTypeGetFavInfo uint32 = 438
	MMRequestTypeShareFav   uint32 = 608
	// MMRequestTypeGetLoginQRCode 类型：获取二维码
	MMRequestTypeGetLoginQRCode uint32 = 502
	// MMRequestTypeCheckLoginQRCode 类型：检测二维码状态
	MMRequestTypeCheckLoginQRCode uint32 = 503
	// MMRequestTypePushQrLogin 类型：二维码二次登录
	MMRequestTypePushQrLogin uint32 = 654
	//MMRequestTypeHeartBeat 类型：心跳包
	MMRequestTypeHeartBeat uint32 = 518
	// MMRequestTypeNewSendMsg 类型：发送消息
	MMRequestTypeNewSendMsg uint32 = 522
	// MMRequestTypeGetOnlineInfo 类型：登录信息
	MMRequestTypeGetOnlineInfo uint32 = 526
	// MMRequestTypeGetChatRoomMemberDetail 类型：获取微信群成员信息列表
	MMRequestTypeGetChatRoomMemberDetail uint32 = 551
	// MMRequestTypeRevokeMsg 类型：撤回消息
	MMRequestTypeRevokeMsg uint32 = 594
	//MMRequestTypeInviteChatRoomMember 类型：邀请群成员
	MMRequestTypeInviteChatRoomMember uint32 = 610
	// MMRequestTypeAddContactLabel 类型：添加标签
	MMRequestTypeAddContactLabel uint32 = 635
	// MMRequestTypeDelContactLabel 类型：删除标签
	MMRequestTypeDelContactLabel uint32 = 636
	// MMRequestTypeUpdateContactLabel 类型：跟新标签名称
	MMRequestTypeUpdateContactLabel uint32 = 637
	// MMRequestTypeModifyContactLabelList 类型：修改好友标签列表
	MMRequestTypeModifyContactLabelList uint32 = 638
	// MMRequestTypeGetContactLabelList 类型：获取标签列表
	MMRequestTypeGetContactLabelList uint32 = 639
	// MMRequestTypeOplog 类型：Oplog
	MMRequestTypeOplog uint32 = 681
	// MMRequestTypeManualAuth 类型：扫码登陆
	MMRequestTypeManualAuth uint32 = 701
	// MMRequestTypeHybridManualAuth hybrid 登录
	MMRequestTypeHybridManualAuth uint32 = 252
	// MMRequestTypeAutoAuth 类型：token登陆
	MMRequestTypeAutoAuth uint32 = 702
	// MMRequestTypeInitContact 类型：初始化联系人
	MMRequestTypeInitContact uint32 = 851
	// MMRequestTypeBatchGetContactBriefInfo 类型：批量获取联系人信息
	MMRequestTypeBatchGetContactBriefInfo uint32 = 945
	// MMRequestTypeSetChatRoomAnnouncement 类型：修改群公告
	MMRequestTypeSetChatRoomAnnouncement uint32 = 993
	// MMRequestTypeJsLogin 类型：小程序授权
	MMRequestTypeJSLogin uint32 = 1029
	// MMRequestTypeJSOperateWxData 类型：小程序
	MMRequestTypeJSOperateWxData uint32 = 1133
	// MMRequestTypeSdkOauthAuthorize 类型：授权app应用
	MMRequestTypeSdkOauthAuthorize uint32 = 1388
	// MMRequestTypeBindQueryNew 查询红包支付信息
	MMRequestTypeBindQueryNew uint32 = 1501
	// MMRequestTypeQryListWxHB 类型：获取领取的红包列表信息
	MMRequestTypeQryListWxHB uint32 = 1514
	// MMRequestTypeReceiveWxHB 类型：接收微信红包
	MMRequestTypeReceiveWxHB uint32 = 1581
	// MMRequestTypeQryDetailWxHB 类型：查询红包领取详情
	MMRequestTypeQryDetailWxHB uint32 = 1585
	// MMRequestTypeOpenWxHB 类型：打开微信红包
	MMRequestTypeOpenWxHB uint32 = 1685
)

const (
	// MMUserInfoStateNew 微信号状态：新建状态
	MMUserInfoStateNew uint32 = 0
	// MMUserInfoStateOnline 微信号状态：在线
	MMUserInfoStateOnline uint32 = 1
	// MMUserInfoStateOffline 微信号状态：离线
	MMUserInfoStateOffline uint32 = 2
)

const (
	// MMLoginQrcodeStateNone 登陆二维码状态：未空状态
	MMLoginQrcodeStateNone uint32 = 0
	// MMLoginQrcodeStateScaned 登陆二维码状态：扫描
	MMLoginQrcodeStateScaned uint32 = 1
	// MMLoginQrcodeStateSure 登陆二维码状态：点击了确定登陆
	MMLoginQrcodeStateSure uint32 = 2
)

// AppMsg 消息类型
const (
	// MMAppMsgTypePayInfo App消息类型：支付
	MMAppMsgTypePayInfo uint32 = 2001

	// MMPayInfoSceneIDHongBao App支付类型ID：微信红包
	MMPayInfoSceneIDHongBao uint32 = 1002
)

const (
	// MMZombieFanCheckStateNone 检测僵尸粉状态：未开始检测
	MMZombieFanCheckStateNone uint32 = 0
	// MMZombieFanCheckStateIning 检测僵尸粉状态：正在检测中
	MMZombieFanCheckStateIning uint32 = 1
	// MMZombieFanCheckStateFinish 检测僵尸粉状态：检测完毕
	MMZombieFanCheckStateFinish uint32 = 2
)

// 红包相关
const (
	// MMHongBaoReqCgiCmdReceiveWxhb 红包请求类型：接收红包
	MMHongBaoReqCgiCmdReceiveWxhb uint32 = 3
	// MMHongBaoReqCgiCmdOpenWxhb 红包请求类型：打开红包
	MMHongBaoReqCgiCmdOpenWxhb uint32 = 4
	// MMHongBaoReqCgiCmdQryDetailWxhb 红包请求类型：查看红包领取详情
	MMHongBaoReqCgiCmdQryDetailWxhb uint32 = 5
	// MMHongBaoReqCgiCmdQryListWxhb 红包请求类型：查看领取的红包列表
	MMHongBaoReqCgiCmdQryListWxhb uint32 = 6

	// MMTenPayReqOutputTypeJSON 红包请求响应的数据格式类型：JSON
	MMTenPayReqOutputTypeJSON uint32 = 1

	// MMHongBaoReqInAwayGroup 接收群红包
	MMHongBaoReqInAwayGroup uint32 = 0
	// MMHongBaoReqInAwayPersonal 接收私人转的红包
	MMHongBaoReqInAwayPersonal uint32 = 1
)

const (
	// MMAddMsgTypeText 消息类型：文本消息
	MMAddMsgTypeText uint32 = 1
	// MMAddMsgTypeImage 消息类型：图片消息
	MMAddMsgTypeImage uint32 = 3
	// MMAddMsgTypeCard 消息类型：名片
	MMAddMsgTypeCard uint32 = 42
	// MMAddMsgTypeRefer 消息类型：引用
	MMAddMsgTypeRefer uint32 = 49
	// MMAddMsgTypeStatusNotify 消息类型：状态通知
	MMAddMsgTypeStatusNotify uint32 = 51
	// MMAddMsgTypeSystemMsg 消息类型：系统消息
	MMAddMsgTypeSystemMsg uint32 = 10002
)

const (
	// MMBitValGroupSaveInAddressBook 保存群聊到通讯录
	MMBitValGroupSaveInAddressBook uint32 = 0x1
	// MMBitValChatOnTop 聊天置顶
	MMBitValChatOnTop uint32 = 0x800
)

const (
	// MMLoginStateNoLogin 未登录状态
	MMLoginStateNoLogin uint32 = 0
	// MMLoginStateOnLine 登录后在线
	MMLoginStateOnLine uint32 = 1
	// MMLoginStateOffLine 登录后离线
	MMLoginStateOffLine uint32 = 2
	// MMLoginStateLogout 登录后退出
	MMLoginStateLogout uint32 = 3
	// MMLoginStateLoginErr 登录失败！
	MMLoginStateLoginErr uint32 = 4
)

const (
	// MMLbsLifeOpcodeNormal 获取地址列表操作类型：非自动
	MMLbsLifeOpcodeNormal uint32 = 0
	// MMLbsLifeOpcodeAuto 获取地址列表操作类型：自动获取
	MMLbsLifeOpcodeAuto uint32 = 1
)

const (
	// MMFavSyncCmdAddItem 新增收藏项
	MMFavSyncCmdAddItem uint32 = 200

	// MMFavItemTypeText 收藏类型：文字
	MMFavItemTypeText uint32 = 1
	// MMFavItemTypeImage 收藏类型：图片
	MMFavItemTypeImage uint32 = 2
	// MMFavItemTypeShare 收藏类型：分享
	MMFavItemTypeShare uint32 = 5
	// MMFavItemTypeVedio 收藏类型：视频
	MMFavItemTypeVedio uint32 = 16
)

const (
	// MMLongOperatorTypeFavSync 操作类型: 同步收藏
	MMLongOperatorTypeFavSync uint32 = 192
)

const (
	// MMRequestRetSessionTimeOut 链接失效
	MMRequestRetSessionTimeOut int32 = -2
	// MMRequestRetMMTLSError mmtls错误
	MMRequestRetMMTLSError int32 = -1
)

const (
	// MMConcealAddNoNeedVerify 别人添加我时不需要验证
	MMConcealAddNoNeedVerify uint32 = 0
	// MMConcealAddNeedVerify 别人添加我时需要验证
	MMConcealAddNeedVerify uint32 = 1

	// MMSwitchFunctionOFF 开关-关闭
	MMSwitchFunctionOFF uint32 = 1
	// MMSwitchFunctionON 开关-打开
	MMSwitchFunctionON uint32 = 2

	// MMAddMeNeedVerifyType 通过手机号查找我
	MMAddMeNeedVerifyType uint32 = 4
	// MMFindMeByPhoneType 通过手机号查找我
	MMFindMeByPhoneType uint32 = 8
	// MMFindMeByWxIDType 通过微信号查找我
	MMFindMeByWxIDType uint32 = 25
	// MMFindMeByGroupType 通过群聊添加我
	MMFindMeByGroupType uint32 = 38
	// MMFindMeByMyQRCodeType 通过我的二维码添加我
	MMFindMeByMyQRCodeType uint32 = 39
	// MMFindMeByCardType 通过名片添加我
	MMFindMeByCardType uint32 = 40
)

const (
	// MMChatroomAccessTypeNeedVerify 进群需要验证
	MMChatroomAccessTypeNeedVerify uint32 = 1
)
