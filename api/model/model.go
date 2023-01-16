package model

import (
	"feiyu.com/wx/clientsdk/baseinfo"
)

type MessageItem struct {
	MsgIds       string
	ToUserName   string
	TextContent  string
	ImageContent string
	MsgType      int //1 Text 2 Image
	AtWxIDList   []string
}

type SendMessageModel struct {
	MsgItem []MessageItem
}

type SendShareCardModel struct {
	ToUserName string
	Id         string
	Nickname   string
	Alias      string
}

// SnsLocationInfo 朋友圈地址项
type SnsLocationInfoModel struct {
	City               string
	Longitude          string
	Latitude           string
	PoiName            string
	PoiAddress         string
	PoiScale           int32
	PoiInfoURL         string
	PoiClassifyID      string
	PoiClassifyType    uint32
	PoiClickableStatus uint32
}

// SnsMediaItem 朋友圈媒体项
type SnsMediaItemModel struct {
	ID            uint32
	Type          uint32
	Title         string
	Description   string
	Private       uint32
	UserData      string
	SubType       uint32
	URL           string
	URLType       string
	Thumb         string
	ThumType      string
	SizeWidth     string
	SizeHeight    string
	TotalSize     string
	VideoWidth    string
	VideoHeight   string
	MD5           string
	VideoMD5      string
	VideoDuration float64
}

type DownloadMediaModel struct {
	Key string
	URL string
}
type TransmitFriendCircleModel struct {
	SourceID string
}

// 设置朋友圈可见天数
type SetFriendCircleDaysModel struct {
	Function uint32
	Value    uint32
}

// SnsPostItem 发送朋友圈需要的信息
type SnsPostItemModel struct {
	ContentStyle  uint32 // 纯文字/图文/引用/视频
	ContentUrl    string
	Description   string
	Privacy       uint32                // 是否仅自己可见
	Content       string                // 文本内容
	MediaList     []*SnsMediaItemModel  // 图片/视频列表
	WithUserList  []string              // 提醒好友看列表
	GroupUserList []string              // 可见好友列表
	BlackList     []string              // 不可见好友列表
	LocationInfo  *SnsLocationInfoModel // 发送朋友圈的位置信息
}

type GetSnsInfoModel struct {
	UserName     string
	FirstPageMD5 string
	MaxID        uint64
}
type GetIdDetailModel struct {
	Id          string
	BlackList   []string
	LocationVal int64
	Location    baseinfo.Location
}
type SetBackgroundImageModel struct {
	Url string
}
type SendFavItemCircle struct {
	SourceID    string
	FavItemID   uint32
	BlackList   []string
	LocationVal int64
	Location    baseinfo.Location
}

// UploadFriendCircleModel 上传朋友圈图片视频信息
type UploadFriendCircleModel struct {
	ImageDataList []string
	VideoDataList []string
}

type ForwardImageItem struct {
	ToUserName      string
	AesKey          string
	CdnMidImgUrl    string
	CdnMidImgSize   int32
	CdnThumbImgSize int32
}

type ForwardVideoItem struct {
	ToUserName     string
	AesKey         string
	CdnVideoUrl    string
	Length         int
	PlayLength     int
	CdnThumbLength int
}

type ForwardMessageModel struct {
	ForwardImageList []ForwardImageItem
	ForwardVideoList []ForwardVideoItem
}

type SendEmojiItem struct {
	ToUserName string
	EmojiMd5   string
	EmojiSize  int32
}
type SendEmojiMessageModel struct {
	EmojiList []SendEmojiItem
}

type RevokeMsgModel struct {
	NewMsgId    string
	ClientMsgId uint64
	ToUserName  string
}

type AppMessageItem struct {
	ToUserName  string
	ContentXML  string
	ContentType uint32
}

type AppMessageModel struct {
	AppList []AppMessageItem
}

type UpdateChatroomAnnouncementModel struct {
	ChatRoomName string
	Content      string
}

type GetChatroomMemberDetailModel struct {
	ChatRoomName string
}

type ChatRoomWxIdListModel struct {
	ChatRoomWxIdList []string
}
type SetChatroomAccessVerifyModel struct {
	ChatRoomName string
	Enable       bool
}
type ChatroomMemberModel struct {
	UserList     []string
	ChatRoomName string
}
type ChatroomNameModel struct {
	Nickname     string
	ChatRoomName string
}
type SendPatModel struct {
	ChatRoomName string
	ToUserName   string
	Scene        int64
}
type GroupListModel struct {
	Key string
}
type SearchContactModel struct {
	Tg        string
	FromScene uint64
	UserName  string
}
type QWApplyAddContactModel struct {
	UserName string
	V1       string
	Content  string
}
type QWSyncChatRoomModel struct {
	Key string
}
type QWChatRoomTransferOwnerModel struct {
	ChatRoomName string
	ToUserName   string
}
type QWAddChatRoomMemberModel struct {
	ChatRoomName string
	ToUserName   []string
}
type QWAcceptChatRoomModel struct {
	Link   string
	Opcode uint32
}
type QWAdminAcceptJoinChatRoomSetModel struct {
	ChatRoomName string
	P            int64
}
type QWModChatRoomNameModel struct {
	ChatRoomName string
	Name         string
}
type MoveContractModel struct {
	ChatRoomName string
	Val          uint32
}
type CreateChatRoomModel struct {
	TopIc    string
	UserList []string
}

type TransferGroupOwnerModel struct {
	ChatRoomName     string
	NewOwnerUserName string
}

type InviteChatroomMembersModel struct {
	ChatRoomName string
	UserList     []string
}

type ScanIntoUrlGroupModel struct {
	Url string
}

type SnsObjectOpItem struct {
	SnsObjID string // 朋友圈ID
	OpType   uint32 // 操作码
	DataLen  uint32 // 其它数据长度
	Data     []byte // 其它数据
	Ext      uint32
}

type SendSnsObjectOpRequestModel struct {
	SnsObjectOpList []SnsObjectOpItem
}

// ReplyCommentItem 回覆的评论项
type ReplyCommentItem struct {
	UserName string // 评论的微信ID
	NickName string // 发表评论的昵称
	OpType   uint32 // 操作类型：评论/点赞
	Source   uint32 // source
}

// SnsCommentItem 朋友圈项：发表评论/点赞
type SnsCommentItem struct {
	OpType         uint32           // 操作类型：评论/点赞
	ItemID         string           // 朋友圈项ID
	ToUserName     string           // 好友微信ID
	Content        string           // 评论内容
	CreateTime     uint32           // 创建时间
	ReplyCommentID uint32           // 回复的评论ID
	ReplyItem      ReplyCommentItem // 回覆的评论项
}

type SendSnsCommentRequestModel struct {
	SnsCommentList []SnsCommentItem
	Tx             bool
}

type SnsObjectOpRequestModel struct {
	SnsObjectOpList []SnsObjectOpItem
}

type DeviceIdLoginModel struct {
	Proxy      string
	DeviceId   string
	UserName   string
	Password   string
	Ticket     string
	Type       int
	DeviceInfo DeviceInfo
}
type DeviceInfo struct {
	Language     string
	Model        string
	AndroidId    string
	Manufacturer string
	ImeI         string
}
type DelSafeDeviceModel struct {
	DeviceUUID string
}

type ExtDeviceLoginModel struct {
	QrConnect string
}

type GetA8KeyRequestModel struct {
	OpCode uint32
	Scene  uint32
	ReqUrl string
}

type AppletModel struct {
	AppId       string
	SdkName     string
	PackageName string
}

type VerifyUserItem struct {
	Gh    string
	Scene byte
}

// 获取通讯录好友
type GetContactListModel struct {
	CurrentWxcontactSeq       uint32
	CurrentChatRoomContactSeq uint32
}

type FollowGHModel struct {
	GHList []VerifyUserItem
}

type GetLoginQrCodeModel struct {
	Proxy    string
	DeviceId string
}

type PhoneLoginModel struct {
	Url string
}

type UploadMContactModel struct {
	MobileList []string
	Mobile     string
}

type QRConnectAuthorizeModel struct {
	QrUrl string
}
type GetMpA8KeyModel struct {
	Url    string
	Opcode uint32
	Scene  int64
}

type GetQrCodeModel struct {
	Id      string
	Recover bool
}
type PeopleNearbyModel struct {
	Longitude float32
	Latitude  float32
}
type CollectmoneyModel struct {
	InvalidTime   string
	TransFerId    string
	TransactionId string
	ToUserName    string
}
type OpenRedEnvelopesModel struct {
	NativeUrl string
}
type FavInfoModel struct {
	FavId  uint32
	KeyBuf string
}

type ShareFavModel struct {
	FavId      uint32
	ToUserName string
}

type CheckFavCdnModel struct {
	DataId         string
	DataSourceId   string
	DataSourceType uint32
	FullMd5        string
	FullSize       uint32
	Head256Md5     string
	IsThumb        uint32
}

type LabelModel struct {
	UserLabelList []baseinfo.UserLabelInfoItem
	// del labelId
	LabelId       string
	LabelNameList []string
}

type GetSyncMsgModel struct {
	Key string
}
type DelContactModel struct {
	DelUserName string
}

type ModifyUserInfo struct {
	City      string
	Country   string
	NickName  string
	Province  string
	Signature string
	Sex       uint32
	InitFlag  uint32
}
type UpdateNickNameModel struct {
	Scene uint32
	Val   string
}
type UpdateSexModel struct {
	Sex      uint32
	City     string
	Province string
	Country  string
}
type UploadHeadImageModel struct {
	Base64 string
}

type SendChangePwdRequestModel struct {
	OldPass, NewPass string
	OpCode           uint32
}
type SendModifyRemarkRequestModel struct {
	UserName   string
	RemarkName string
}

// 设置微信号
type AlisaModel struct {
	Alisa string
}

type UpdateAutopassModel struct {
	SwitchType uint32
}

// 设置添加方式
type WxFunctionSwitchModel struct {
	Function uint32
	Value    uint32
}
type SetSendPatModel struct {
	Value string
}
type BindMobileModel struct {
	Mobile     string
	VerifyCode string
}
type SendVerifyMobileModel struct {
	Mobile string
	Opcode uint32
}

// 修改步数
type UpdateStepNumberModel struct {
	Number uint64
}

type UserRankLikeModel struct {
	RankId string
}

// BatchGetContact
type BatchGetContactModel struct {
	UserNames    []string
	RoomWxIDList []string
}
type SearchContactRequestModel struct {
	OpCode, FromScene, SearchScene uint32
	UserName                       string
}

type GetFriendRelationModel struct {
	UserName string
}
type SendUploadVoiceRequestModel struct {
	ToUserName               string
	VoiceData                string `json:"VoiceData"`
	VoiceSecond, VoiceFormat int32
}

type CdnUploadVideoRequest struct {
	ToUserName string
	VideoData  []byte // 视频数据
	ThumbData  string // ThumbData
}
type GetMsgBigImgModel struct {
	Datatotalength int
	ToWxid         string
	MsgId          uint32
	NewMsgId       uint64
}
type GroupMassMsgTextModel struct {
	ToUserName []string
	Content    string
}
type GroupMassMsgImageModel struct {
	ToUserName  []string
	ImageBase64 string
}
type SyncModel struct {
	Scene   uint32
	SyncKey string
}

type DownloadVoiceModel struct {
	ToUserName string
	NewMsgId   string
	Bufid      string
	Length     int
}
type DownMediaModel struct {
	AesKey   string
	FileURL  string
	FileType uint32
}
type VerifyUserRequestModel struct {
	OpCode                   uint32
	VerifyContent            string
	Scene                    byte
	V1, V2, ChatRoomUserName string
}
type GeneratePayQCodeModel struct {
	Name  string
	Money string
}

type QWContactModel struct {
	ToUserName string
	ChatRoom   string
	T          string
}

type QWRemarkModel struct {
	ToUserName string
	Name       string
}
type QWCreateModel struct {
	ToUserName []string
}
type FinderSearchModel struct {
	Index   uint32
	Userver int32
	UserKey string
	Uuid    string
}

type FinderUserPrepareModel struct {
	Userver int32
}

type FinderFollowModel struct {
	FinderUserName string
	PosterUsername string
	OpType         int32
	RefObjectId    string
	Cook           string
	Userver        int32
}
type TargetUserPageParam struct {
	Wxid       string
	Target     string
	LastBuffer string
}

type WxBindOpMobileForModel struct {
	OpCode      int64
	PhoneNumber string
	VerifyCode  string
	Reg         uint64
	Proxy       string
}

type ExtDeviceLoginConfirmModel struct {
	Url string
}
