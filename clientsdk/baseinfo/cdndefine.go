package baseinfo

import "feiyu.com/wx/protobuf/wechat"

// CDNUploadMsgImgPrepareRequestItem 请求项
type CDNUploadMsgImgPrepareRequestItem struct {
	ToUser      string //  接受人微信ID
	LocalName   string //  本地名称(不包含扩展名)
	CreateTime  uint32 // 创建时间
	ThumbWidth  int32  // 缩略图宽
	ThumbHeight int32  // 缩略图高
	AesKey      []byte // SessionKey
	Crc32       uint32 // 源图片Crc32
}

type UploadVideoItem struct {
	ToUser     string             // 视频接收人
	AesKey     []byte             // 加密用的AesKey
	Seq        uint32             // 代表第几个请求
	VideoID    uint32             // ID
	CreateTime uint32             // 创建时间
	VideoData  []byte             // 视频数据
	ThumbData  []byte             // ThumbData
	CDNDns     *wechat.CDNDnsInfo // DNS信息
}

// UploadImgItem 上传图片项
type UploadImgItem struct {
	ToUser     string             // 图片接收人
	Seq        uint32             // 发送图片的序号代表 今天第几张
	LocalName  string             // 图片本地名称（可以随机，不包含扩展名）
	ExtName    string             // 图片扩展名
	AesKey     []byte             // 加密用的AesKey
	ImageData  []byte             // 源图片数据
	CreateTime uint32             // 发送时间
	CDNDns     *wechat.CDNDnsInfo // DNS信息
}

// SnsUploadImgItem 朋友圈上传图片项
type SnsUploadImgItem struct {
	AesKey     []byte             // 加密用的AesKey
	Seq        uint32             // 代表第几个请求
	ImageData  []byte             // 源图片数据
	ImageID    uint32             // 图片ID
	CreateTime uint32             // 发送时间
	CDNDns     *wechat.CDNDnsInfo // DNS信息
}

// SnsVideoDownloadItem 朋友圈视频下载项
type SnsVideoDownloadItem struct {
	Seq           uint32             // 代表第几个请求
	URL           string             // 视频加密地址
	RangeStart    uint32             // 起始地址
	RangeEnd      uint32             // 结束地址
	XSnsVideoFlag string             // 视频标志
	CDNDns        *wechat.CDNDnsInfo // DNS信息
}

// SnsVideoUploadItem 朋友圈视频上传项
type SnsVideoUploadItem struct {
	AesKey     []byte             // 加密用的AesKey
	Seq        uint32             // 代表第几个请求
	VideoID    uint32             // ID
	CreateTime uint32             // 创建时间
	VideoData  []byte             // 视频数据
	ThumbData  []byte             // ThumbData
	CDNDns     *wechat.CDNDnsInfo // DNS信息
}

// CdnImageDownloadRequest 高清图片下载请求
type CdnImageDownloadRequest struct {
	Ver           uint32
	WeiXinNum     uint32
	Seq           uint32
	ClientVersion uint32
	ClientOsType  string
	AuthKey       []byte
	NetType       uint32
	AcceptDupack  uint32
	RsaVer        uint32
	RsaValue      []byte
	FileType      uint32
	WxChatType    uint32
	FileID        string
	LastRetCode   uint32
	IPSeq         uint32
	CliQuicFlag   uint32
	WxMsgFlag     *uint32
	WxAutoStart   uint32
	DownPicFormat uint32
	Offset        uint32
	LargesVideo   uint32
	SourceFlag    uint32
}

// CdnDownloadResponse Cdn下载响应
type CdnDownloadResponse struct {
	Ver             uint32
	Seq             uint32
	VideoFormat     uint32
	RspPicFormat    uint32
	RangeStart      uint32
	RangeEnd        uint32
	TotalSize       uint32
	SrcSize         uint32
	RetCode         uint32
	SubStituteFType uint32
	RetrySec        uint32
	IsRetry         uint32
	IsOverLoad      uint32
	IsGetCdn        uint32
	XClientIP       string
	FileData        []byte
}

// CdnImageUploadRequest 高清图片上传请求
type CdnImageUploadRequest struct {
	Ver            uint32 // 1
	WeiXinNum      uint32 //
	Seq            uint32 // 6
	ClientVersion  uint32
	ClientOsType   string
	AuthKey        []byte
	NetType        uint32 // 1
	AcceptDupack   uint32 // 1
	SafeProto      uint32 // 1
	FileType       uint32 // 2
	WxChatType     uint32 // 1
	LastRetCode    uint32 // 0
	IPSeq          uint32 // 0
	CliQuicFlag    uint32 // 0
	HasThumb       uint32 // 1
	ToUser         string // @cdn2_9887af1554e6f59f5e0489e399439cffe8fd07b9009032161122cee11c8537dd
	CompressType   uint32 // 0
	NoCheckAesKey  uint32 // 1
	EnableHit      uint32 // 1
	ExistAnceCheck uint32 // 0
	AppType        uint32 // 1
	FileKey        string // wxupload_21533455325@chatroom29_1572079793
	TotalSize      uint32 // 53440
	RawTotalSize   uint32 // 53425
	LocalName      string // 29.wxgf
	SessionBuf     []byte // CDNUploadMsgImgPrepareRequest
	Offset         uint32 // 0
	ThumbTotalSize uint32 // 4496
	RawThumbSize   uint32 // 4487
	RawThumbMD5    string // 0d29df2b74d29efa46dd6fa1e75e71ba
	EncThumbCRC    uint32 // 2991702343
	ThumbData      []byte // 缩略图加密后数据
	LargesVideo    uint32 // 0
	SourceFlag     uint32 // 0
	AdVideoFlag    uint32 // 0
	FileMD5        string // e851e118f524b4219928bed3f3bd0d24
	RawFileMD5     string // e851e118f524b4219928bed3f3bd0d24
	DataCheckSum   uint32 // 737909102
	FileCRC        uint32 // 2444306137
	SetOfPicFormat string // 001010
	FileData       []byte // 文件数据
}

// CdnVideoUploadRequest 视频上传请求
type CdnVideoUploadRequest struct {
	Ver            uint32
	WeiXinNum      uint32
	Seq            uint32
	ClientVersion  uint32
	ClientOSType   string
	AutoKey        []byte
	NetType        uint32
	AcceptDuPack   uint32
	SafeProto      uint32
	FileType       uint32
	WeChatType     uint32
	LastRetCode    uint32
	IpSeq          uint32
	HastHumb       uint32
	ToUSerName     string
	CompressType   uint32
	NoCheckAesKey  uint32
	EnaBleHit      uint32
	ExistAnceCheck uint32
	AppType        uint32
	FileKey        string
	TotalSize      uint32
	RawTotalSize   uint32
	LocalName      string
	Offset         uint32
	ThumbTotalSize uint32
	RawThumbSize   uint32
	RawThumbMd5    string
	EncThumbCrc    uint32
	ThumbData      []byte
	LargesVideo    uint32
	SourceFlag     uint32
	AdVideoFlag    uint32
	Mp4identify    string
	DropRateFlag   uint32
	ClientRsaVer   uint32
	ClientRsaVal   []byte
	FileMd5        string
	RawFileMd5     string
	DataCheckSum   uint32
	FileCrc        uint32
	FileData       []byte
}

// CdnImageUploadResponse 高清图片上传响应
type CdnImageUploadResponse struct {
	Ver        uint32
	Seq        uint32
	RetCode    uint32
	FileKey    string
	RecvLen    uint32
	SKeyResp   uint32
	SKeyBuf    []byte
	FileID     string
	ExistFlag  uint32
	HitType    uint32
	RetrySec   uint32
	IsRetry    uint32
	IsOverLoad uint32
	IsGetCDN   uint32
	XClientIP  string
}

// CdnSnsImageUploadRequest 朋友圈图片上传请求
type CdnSnsImageUploadRequest struct {
	Ver            uint32 // 1
	WeiXinNum      uint32 //
	Seq            uint32 // 6
	ClientVersion  uint32
	ClientOsType   string
	AuthKey        []byte
	NetType        uint32 // 1
	AcceptDupack   uint32 // 1
	RsaVer         uint32 // 1
	RsaValue       []byte
	FileType       uint32 // 2
	WxChatType     uint32 // 1
	LastRetCode    uint32 // 0
	IPSeq          uint32 // 0
	CliQuicFlag    uint32 // 0
	HasThumb       uint32 // 1
	ToUser         string // @cdn2_9887af1554e6f59f5e0489e399439cffe8fd07b9009032161122cee11c8537dd
	CompressType   uint32 // 0
	NoCheckAesKey  uint32 // 1
	EnableHit      uint32 // 1
	ExistAnceCheck uint32 // 0
	AppType        uint32 // 1
	FileKey        string // wxupload_21533455325@chatroom29_1572079793
	TotalSize      uint32 // 53440
	RawTotalSize   uint32 // 53425
	LocalName      string // 29.wxgf
	Offset         uint32 // 0
	ThumbTotalSize uint32 // 4496
	RawThumbSize   uint32 // 4487
	RawThumbMD5    string // 0d29df2b74d29efa46dd6fa1e75e71ba
	ThumbCRC       uint32 // 2991702343
	LargesVideo    uint32 // 0
	SourceFlag     uint32 // 0
	AdVideoFlag    uint32 // 0
	FileMD5        string // e851e118f524b4219928bed3f3bd0d24
	RawFileMD5     string // e851e118f524b4219928bed3f3bd0d24
	DataCheckSum   uint32 // 737909102
	FileCRC        uint32 // 2444306137
	FileData       []byte // 文件数据
}

// CdnSnsImageUploadResponse 高清图片上传响应
type CdnSnsImageUploadResponse struct {
	Ver         uint32
	Seq         uint32
	RetCode     uint32
	FileKey     string
	RecvLen     uint32
	FileURL     string
	ThumbURL    string
	EnableQuic  uint32
	RetrySec    uint32
	IsRetry     uint32
	IsOverLoad  uint32
	IsGetCDN    uint32
	XClientIP   string
	ImageMD5    string
	ImageWidth  uint32
	ImageHeight uint32
}

// CdnSnsVideoDownloadRequest 朋友圈视频下载请求
type CdnSnsVideoDownloadRequest struct {
	Ver             uint32
	WeiXinNum       uint32
	Seq             uint32
	ClientVersion   uint32
	ClientOsType    string
	AuthKey         []byte
	NetType         uint32
	AcceptDupack    uint32
	Signal          string
	Scene           string
	URL             string
	RangeStart      uint32
	RangeEnd        uint32
	LastRetCode     uint32
	IPSeq           uint32
	RedirectType    uint32
	LastVideoFormat uint32
	VideoFormat     uint32
	XSnsVideoFlag   string
}

// CdnSnsVideoDownloadResponse 朋友圈视频下载响应
type CdnSnsVideoDownloadResponse struct {
	Ver             uint32
	Seq             uint32
	RangeStart      uint32
	RangeEnd        uint32
	TotalSize       uint32
	RetCode         uint32
	EnableQuic      uint32
	IsRetry         uint32
	IsOverLoad      uint32
	IsGetCdn        uint32
	XClientIP       string
	XSnsVideoFlag   string
	XSnsVideoTicket string
	XEncFlag        uint32
	XEncLen         uint32
	FileData        []byte
}

// CdnSnsVideoUploadRequest 朋友圈视频上传请求
type CdnSnsVideoUploadRequest struct {
	Ver              uint32 // 1
	WeiXinNum        uint32 //
	Seq              uint32 // 6
	ClientVersion    uint32
	ClientOsType     string
	AuthKey          []byte
	NetType          uint32 // 1
	AcceptDupack     uint32 // 1
	RsaVer           uint32 // 1
	RsaValue         []byte
	FileType         uint32 // 2
	WxChatType       uint32 // 1
	LastRetCode      uint32 // 0
	IPSeq            uint32 // 0
	CliQuicFlag      uint32 // 0
	HasThumb         uint32 // 1
	NoCheckAesKey    uint32 // 1
	EnableHit        uint32 // 1
	ExistAnceCheck   uint32 // 0
	AppType          uint32 // 1
	FileKey          string // wxupload_21533455325@chatroom29_1572079793
	TotalSize        uint32 // 53440
	RawTotalSize     uint32 // 53425
	LocalName        string // 29.wxgf
	Offset           uint32 // 0
	ThumbTotalSize   uint32 // 4496
	RawThumbSize     uint32 // 4487
	RawThumbMD5      string // 0d29df2b74d29efa46dd6fa1e75e71ba
	ThumbCRC         uint32 // 2991702343
	IsStoreVideo     uint32
	ThumbData        []byte
	LargesVideo      uint32 // 0
	SourceFlag       uint32 // 0
	AdVideoFlag      uint32 // 0
	Mp4Identify      string
	FileMD5          string // e851e118f524b4219928bed3f3bd0d24
	RawFileMD5       string // e851e118f524b4219928bed3f3bd0d24
	DataCheckSum     uint32 // 737909102
	FileCRC          uint32 // 2444306137
	FileData         []byte // 文件数据
	UserLargeFileApi bool
}

// CdnSnsVideoUploadResponse 上传朋友圈视频响应
type CdnSnsVideoUploadResponse struct {
	Ver        uint32
	Seq        uint32
	RetCode    uint32
	FileKey    string
	RecvLen    uint32
	FileURL    string
	ThumbURL   string
	FileID     string
	EnableQuic uint32
	RetrySec   uint32
	IsRetry    uint32
	IsOverLoad uint32
	IsGetCDN   uint32
	XClientIP  string
	ReqData    *CdnSnsVideoUploadRequest
}

// CdnMsgVideoUploadResponse 上传视频
type CdnMsgVideoUploadResponse struct {
	Ver           uint32
	Seq           uint32
	RetCode       uint32
	FileKey       string
	RecvLen       uint32
	FileURL       string
	ThumbURL      string
	FileID        string
	EnableQuic    uint32
	RetrySec      uint32
	IsRetry       uint32
	IsOverLoad    uint32
	IsGetCDN      uint32
	XClientIP     string
	FileAesKey    string
	ThumbDataSize uint32
	VideoDataSize uint32
	VideoDataMD5  string
	Mp4identify   string
	ThumbWidth    uint32
	ThumbHeight   uint32
	ReqData       *CdnVideoUploadRequest
}
