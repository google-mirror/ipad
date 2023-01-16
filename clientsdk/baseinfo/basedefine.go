package baseinfo

import (
	"encoding/base64"
	"time"

	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/cecdh"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/clientsdk/proxynet"
	"feiyu.com/wx/protobuf/wechat"
	"golang.org/x/net/proxy"
)

// SDKVersion 本协议SDK版本号
var SDKVersion = string("1.0.0")

// ClientVersion 微信版本号
var ClientVersion = uint32(0x18001621) //0x17000C2B   0x17000f26 0x18000720
var PlistVersion = uint32(0x18001621)  //plist-version
// MicroMessenger/7.0.12(0x17000c21)
// MicroMessenger/7.0.14(0x17000e2e)
// MicroMessenger/7.0.15(0x17000f26)
// MicroMessenger/7.0.17(0x17001124)   iphone  0X1700112a
// 7.0.18  0x17001231
// 7.0.21  0x17001520
// ServerVersion 微信服务端版本号
var DeviceVersionNumber = "8.0.22(0x18001621)" //版本请求头
var ServerVersion = uint32(0x18001621)         //0x17000C2B  0x17000f26

var DeviceTypeIos = "iPhone iOS12.5.1" // "iPad iOS13.5"  如果是ipad换绑手机发送验证码会提示版本低问题

var DeviceTypeIpad = "iPad iOS13.5" //如果是扫码，则只能用Ipad设备

// 安卓版本号 -  7.019  654315572  654316592   720  0x27001439   7.022   0x27001636
var AndroidClientVersion = uint32(0x28000736)

// 安卓设备类型
var AndroidDeviceType = "android-28"

// LoginRsaVer 登陆用到的RSA版本号
var LoginRsaVer = uint32(135)

var XJLoginRSAVer = uint32(133)

// DefaultLoginRsaVer 默认 登录RSA版本号
var DefaultLoginRsaVer = LoginRsaVer

// Md5OfMachOHeader wechat的MachOHeader md5值 4c541f4fca66dd93a351d4239ecaf7ae
var Md5OfMachOHeader = string("d05a80a94b6c2e3c31424403437b6e18") //

// FileHelperWXID 文件传输助手微信ID
var FileHelperWXID = string("filehelper")

// HomeDIR 当前程序的工作路径
var HomeDIR string

// DeviceInfo 62设备信息
type DeviceInfo struct {
	UUIDOne            string `json:"uuidone"`
	UUIDTwo            string `json:"uuidtwo"`
	Imei               string `json:"imei"`
	DeviceID           []byte `json:"deviceid"`
	DeviceName         string `json:"devicename"`
	TimeZone           string `json:"timezone"`
	Language           string `json:"language"`
	DeviceBrand        string `json:"devicebrand"`
	RealCountry        string `json:"realcountry"`
	IphoneVer          string `json:"iphonever"`
	BundleID           string `json:"boudleid"`
	OsType             string `json:"ostype"`
	AdSource           string `json:"adsource"`
	OsTypeNumber       string `json:"ostypenumber"`
	CoreCount          uint32 `json:"corecount"`
	CarrierName        string `json:"carriername"`
	SoftTypeXML        string `json:"softtypexml"`
	ClientCheckDataXML string `json:"clientcheckdataxml"`
	// extInfo
	GUID2       string `json:"GUID2"`
	DeviceToken *wechat.TrustResp
}

func (d *DeviceInfo) SetDeviceId(deviceId string) {
	d.Imei = deviceId
	d.DeviceID = baseutils.HexStringToBytes(deviceId)
	d.DeviceID[0] = 0x49
}

// LoginDataInfo 62/16 数据登陆
type LoginDataInfo struct {
	Type     byte
	UserName string
	PassWord string
	//伪密码
	NewPassWord string
	//登录数据 62/A16
	LoginData string
	Ticket    string
	NewType   int
	Language  string
}

type SyncMsgKeyMgr struct {
	curKey *wechat.BufferT
	maxKey *wechat.BufferT
}

func (s *SyncMsgKeyMgr) MaxKey() *wechat.BufferT {
	return s.maxKey
}

func (s *SyncMsgKeyMgr) SetMaxKey(maxKey *wechat.BufferT) {
	s.maxKey = maxKey
}

func (s *SyncMsgKeyMgr) CurKey() *wechat.BufferT {
	return s.curKey
}

func (s *SyncMsgKeyMgr) SetCurKey(curKey *wechat.BufferT) {
	s.curKey = curKey
}

// UserInfo 用户信息
type UserInfo struct {
	LoginDataInfo     LoginDataInfo
	HostUrl           string `json:"hostUrl"`
	UUID              string `json:"uuid"`
	QrUuid            string `json:"qrUuid"`
	Uin               uint32 `json:"uin"`
	WxId              string `json:"wxid"`
	NickName          string `json:"nickname"`
	HeadURL           string `json:"headurl"`
	Session           []byte `json:"cookie"`
	SessionKey        []byte `json:"aeskey"`
	ShortHost         string `json:"shorthost"`
	LongHost          string `json:"longhost"`
	LongPort          string `json:"longPort"`
	EcPublicKey       []byte `json:"ecpukey"`
	EcPrivateKey      []byte `json:"ecprkey"`
	CheckSumKey       []byte `json:"checksumkey"`
	AutoAuthKey       []byte `json:"autoauthkey"`
	SyncKey           []byte `json:"synckey"`
	SyncHistoryKey    []byte `json:"syncHistorykey"`
	FavSyncKey        []byte `json:"favsynckey"`
	SnsSyncKey        []byte `json:"snssynckey"`
	HBAesKey          []byte `json:"hbaeskey"`
	HBAesKeyEncrypted string `json:"hbesKeyencrypted"`

	// CDNDns
	DNSInfo     *wechat.CDNDnsInfo `json:"dnsinfo"`
	SNSDnsInfo  *wechat.CDNDnsInfo `json:"snsdnsinfo"`
	APPDnsInfo  *wechat.CDNDnsInfo `json:"appdnsinfo"`
	FAKEDnsInfo *wechat.CDNDnsInfo `json:"fakednsinfo"`

	// ServerDns
	NetworkSect *wechat.NetworkSectResp

	// 设备信息62
	DeviceInfo *DeviceInfo
	//A16信息
	DeviceInfoA16  *AndroidDeviceInfo
	BalanceVersion uint32
	// Wifi信息
	WifiInfo *WifiInfo
	// MMTLS信息
	MMInfo *mmtls.MMInfo
	// 代理信息
	ProxyInfo *proxynet.WXProxyInfo
	// 代理
	Dialer proxy.Dialer

	//Mysql 参数
	loginState uint32

	// HybridKeyVer
	HybridLogin bool
	// 登录的Rsa 密钥版本
	LoginRsaVer uint32
	//lastAuthTime 上次登录时间
	lastAuthTime time.Time

	syncKeyMgr SyncMsgKeyMgr

	Ticket string
}

func (u *UserInfo) SyncKeyMgr() SyncMsgKeyMgr {
	return u.syncKeyMgr
}

func (u *UserInfo) SetSyncKeyMgr(syncKeyMgr SyncMsgKeyMgr) {
	u.syncKeyMgr = syncKeyMgr
}

// CheckCdn 检查cdn信息是否有空的
func (u *UserInfo) CheckCdn() bool {
	return u.DNSInfo == nil || u.SNSDnsInfo == nil || u.APPDnsInfo == nil || u.FAKEDnsInfo == nil
}

func (u *UserInfo) GetMMInfo() *mmtls.MMInfo {
	if u.MMInfo == nil {
		u.MMInfo = mmtls.InitMMTLSInfoShort(u.ShortHost, nil)
	}
	return u.MMInfo
}

// IntervalLastAuthTime 取上次与现在的间隔时间
func (u *UserInfo) IntervalLastAuthTime() time.Duration {
	return time.Now().Sub(u.lastAuthTime)
}

// UpdateLastAuthTime 更新二次登录时间
func (u *UserInfo) UpdateLastAuthTime() {
	u.lastAuthTime = time.Now()
}

func (u *UserInfo) SetProxy(proxy *proxynet.WXProxyInfo) {
	if proxy != nil {
		u.ProxyInfo = proxy
	}
}

// GetLoginRsaVer 获取登录密钥版本号
func (u *UserInfo) GetLoginRsaVer() uint32 {
	if u.LoginRsaVer == 0 {
		u.LoginRsaVer = DefaultLoginRsaVer
	}
	return u.LoginRsaVer
}

// SwitchRSACert  切换证书
func (u *UserInfo) SwitchRSACert() {
	if u.LoginRsaVer == DefaultLoginRsaVer {
		u.LoginRsaVer = XJLoginRSAVer //133
	} else {
		u.LoginRsaVer = DefaultLoginRsaVer // 135
	}
}

// SetLoginStatus 设置登录状态
func (u *UserInfo) SetLoginState(code uint32) {
	u.loginState = code
}

// GetLoginStatus
func (u *UserInfo) GetLoginState() uint32 {
	return u.loginState
}

// GetUserName 取用户账号信息
func (u *UserInfo) GetUserName() string {
	if u.WxId == "" {
		return u.LoginDataInfo.UserName
	} else {
		return u.WxId
	}
}

// SetWxId 设置WxId
func (u *UserInfo) SetWxId(s string) {
	u.WxId = s
}

// SetAutoKey 设置Token
func (u *UserInfo) SetAutoKey(key []byte) {
	if len(key) > 0 {
		u.AutoAuthKey = key
	}
}

// SetNetworkSect
func (u *UserInfo) SetNetworkSect(netWork *wechat.NetworkSectResp) {
	if netWork == nil {
		return
	}
	u.NetworkSect = netWork
}

// ConsultSessionKey 协商密钥并设置
func (u *UserInfo) ConsultSessionKey(ecServerPubKey, sessionKey []byte) {
	u.CheckSumKey = cecdh.ComputerECCKeyMD5(ecServerPubKey, u.EcPrivateKey)
	tmpAesKey, err := baseutils.AesDecryptByteKey(sessionKey, u.CheckSumKey)
	if err != nil {
		//如果密钥协商失败 使用返回的SessionKey
		u.SessionKey = sessionKey[:16]
		//log.Println("ConsultSessionKey 协商密钥失败使用SessionKey.", err)
	} else {
		u.SessionKey = tmpAesKey[:16]
	}
}

// GenHBKey  生成 HBAesKey 和 HBAesKeyEncrypted
func (u *UserInfo) GenHBKey() {
	u.HBAesKey = baseutils.RandomBytes(16)
	hbAesKeyBase64String := base64.StdEncoding.EncodeToString(u.HBAesKey)
	tmpEncKey, _ := baseutils.EncKeyRsaEncrypt([]byte(hbAesKeyBase64String))
	u.HBAesKeyEncrypted = base64.StdEncoding.EncodeToString(tmpEncKey)
}

// WifiInfo WifiInfo
type WifiInfo struct {
	Name      string
	WifiBssID string
}

// ModifyItem 修改用户信息项
type ModifyItem struct {
	CmdID uint32
	Len   uint32
	Data  []byte
}

// HeadImgItem 头像数据项
type HeadImgItem struct {
	ImgPieceData []byte
	TotalLen     uint32
	StartPos     uint32
	ImgHash      string
}

// RevokeMsgItem 撤回消息项
type RevokeMsgItem struct {
	FromUserName   string
	ToUserName     string
	NewClientMsgID uint32
	CreateTime     uint32
	SvrNewMsgID    uint64
	IndexOfRequest uint32
}

// DownMediaItem 下载图片/视频/文件项
type DownMediaItem struct {
	AesKey   string
	FileURL  string
	FileType uint32
}

// DownVoiceItem 下载音频信息项
type DownVoiceItem struct {
	TotalLength  uint32
	NewMsgID     uint64
	ChatRoomName string
	MasterBufID  uint64
}

// VerifyUserItem 添加好友/验证好友/打招呼 项
type VerifyUserItem struct {
	OpType           uint32 // 1免验证发送请求, 2发送验证申请, 3通过好友验证
	FromType         byte   // 1来源QQ，2来源邮箱，3来源微信号，14群聊，15手机号，18附近的人，25漂流瓶，29摇一摇，30二维码，13来源通讯录
	VerifyContent    string // 验证信息
	VerifyUserTicket string // 通过验证UserTicket(同步到的)
	AntispamTicket   string // searchcontact请求返回
	UserValue        string // searchcontact请求返回
	ChatRoomUserName string // 通过群来添加好友 需要设置此值为群id
	NeedConfirm      uint32 // 是否确认
}

// StatusNotifyItem 状态通知项
type StatusNotifyItem struct {
	Code         uint32
	ToUserName   string
	ClientMsgID  string
	FunctionName string
	FunctionArg  string
}

type CheckFavCdnItem struct {
	DataId         string
	DataSourceId   string
	DataSourceType uint32
	FullMd5        string
	FullSize       uint32
	Head256Md5     string
	IsThumb        uint32
}

// SnsLocationInfo 朋友圈地址项
type SnsLocationInfo struct {
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
type SnsMediaItem struct {
	EncKey        string
	EncValue      uint32
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

// SnsPostItem 发送朋友圈需要的信息
type SnsPostItem struct {
	Xml           bool   //Content 是否纯xml
	ContentStyle  uint32 // 纯文字/图文/引用/视频
	Description   string
	ContentUrl    string
	Privacy       uint32           // 是否仅自己可见
	Content       string           // 文本内容
	MediaList     []*SnsMediaItem  // 图片/视频列表
	WithUserList  []string         // 提醒好友看列表
	GroupUserList []string         // 可见好友列表
	BlackList     []string         // 不可见好友列表
	LocationInfo  *SnsLocationInfo // 发送朋友圈的位置信息
}

// SnsObjectOpItem SnsObjectOpItem
type SnsObjectOpItem struct {
	SnsObjID string // 朋友圈ID
	OpType   uint32 // 操作码
	DataLen  uint32 // 其它数据长度
	Data     []byte // 其它数据
	Ext      uint32
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
	OpType         uint32            // 操作类型：评论/点赞
	ItemID         uint64            // 朋友圈项ID
	ToUserName     string            // 好友微信ID
	Content        string            // 评论内容
	CreateTime     uint32            // 创建时间
	ReplyCommentID uint32            // 回复的评论ID
	ReplyItem      *ReplyCommentItem // 回覆的评论项
}

// GetLbsLifeListItem 获取地址列表项
type GetLbsLifeListItem struct {
	Opcode    uint32
	Buffer    []byte
	Longitude float32
	Latitude  float32
	KeyWord   string
}

// UploadVoiceItem 上传语音项
type UploadVoiceItem struct {
	ToUser      string
	Data        []byte
	VoiceLength uint32
	ClientMsgID string
	EndFlag     uint32
}

// LabelItem 标签项
type LabelItem struct {
	Name string
	ID   uint32
}

// UserLabelInfoItem 好友标签信息
type UserLabelInfoItem struct {
	UserName    string
	LabelIDList string
}

// ThumbItem 缩略图数据
type ThumbItem struct {
	Data   []byte
	Width  int32
	Height int32
}

// PackHeader 请求数据包头
type PackHeader struct {
	ReqData        []byte
	RetCode        int32
	Signature      byte
	HeadLength     byte
	CompressType   byte
	EncodeType     byte
	ServerVersion  uint32
	Uin            uint32
	Session        []byte
	SeqId          uint32
	URLID          uint32
	SrcLen         uint32
	ZipLen         uint32
	EncodeVersion  uint32
	HeadDeviceType byte
	CheckSum       uint32
	RunState       byte
	RqtCode        uint32
	EndFlag        byte
	Data           []byte
	HybridKeyVer   byte
}

func (p PackHeader) GetRetCode() int32 {
	return p.RetCode
}

func (p PackHeader) CheckSessionOut() bool {
	return p.RetCode == MMErrSessionTimeOut || p.RetCode == MMRequestRetSessionTimeOut
}

// ForwardImageItem 转发图片信息
type ForwardImageItem struct {
	ToUserName      string
	AesKey          string
	CdnMidImgUrl    string
	CdnMidImgSize   int32
	CdnThumbImgSize int32
}

// ForwardVideoItem 转发视频信息
type ForwardVideoItem struct {
	ToUserName     string
	AesKey         string
	CdnVideoUrl    string
	Length         int
	PlayLength     int
	CdnThumbLength int
}

type CheckLoginQrCodeResult struct {
	*wechat.LoginQRCodeNotify
	Ret                 int32  `json:"ret"`                 // 用户返回错误
	OthersInServerLogin bool   `json:"othersInServerLogin"` // 是否在其他服务器登录
	TargetServer        string `json:"tarGetServerIp"`      // 在其服务器登录的IP
	UUId                string `json:"uuId"`
}
