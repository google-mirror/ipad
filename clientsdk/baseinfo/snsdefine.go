package baseinfo

// TimelineObject TimelineObject
type TimelineObject struct {
	ID                  uint64        `xml:"id"`
	UserName            string        `xml:"username"`
	CreateTime          uint32        `xml:"createTime"`
	ContentDesc         string        `xml:"contentDesc"`
	ContentDescShowType uint32        `xml:"contentDescShowType"`
	ContentDescScene    uint32        `xml:"contentDescScene"`
	Private             uint32        `xml:"private"`
	SightFolded         uint32        `xml:"sightFolded"`
	ShowFlag            uint32        `xml:"showFlag"`
	AppInfo             AppInfo       `xml:"appInfo"`
	SourceUserName      string        `xml:"sourceUserName"`
	SourceNickName      string        `xml:"sourceNickName"`
	StatisticsData      string        `xml:"statisticsData"`
	StatExtStr          string        `xml:"statExtStr"`
	ContentObject       ContentObject `xml:"ContentObject"`
	ActionInfo          ActionInfo    `xml:"actionInfo"`
	Location            Location      `xml:"location"`
	PublicUserName      string        `xml:"publicUserName"`
	StreamVideo         StreamVideo   `xml:"streamvideo"`
}

// AppInfo AppInfo
type AppInfo struct {
	ID            string `xml:"id"`
	Version       string `xml:"version"`
	AppName       string `xml:"appName"`
	InstallURL    string `xml:"installUrl"`
	FromURL       string `xml:"fromUrl"`
	IsForceUpdate uint32 `xml:"isForceUpdate"`
}

// ContentObject ContentObject
type ContentObject struct {
	ContentStyle uint32    `xml:"contentStyle"`
	Title        string    `xml:"title"`
	Description  string    `xml:"description"`
	MediaList    MediaList `xml:"mediaList"`
	ContentURL   string    `xml:"contentUrl"`
}

// MediaList MediaList
type MediaList struct {
	Media []Media `xml:"media"`
}

// Media Media
type Media struct {
	Enc           Enc       `xml:"enc"`
	ID            uint64    `xml:"id"`
	Type          uint32    `xml:"type"`
	Title         string    `xml:"title"`
	Description   string    `xml:"description"`
	Private       uint32    `xml:"private"`
	UserData      string    `xml:"userData"`
	SubType       uint32    `xml:"subType"`
	VideoSize     VideoSize `xml:"videoSize"`
	URL           URL       `xml:"url"`
	Thumb         Thumb     `xml:"thumb"`
	Size          Size      `xml:"size"`
	VideoDuration float64   `xml:"videoDuration"`
}

// Enc Enc
type Enc struct {
	Key   string `xml:"key,attr"`
	Value uint32 `xml:",chardata"`
}

// VideoSize 视频大小
type VideoSize struct {
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}

// URL URL
type URL struct {
	Type     string `xml:"type,attr"`
	Token    string `xml:"token,attr"`
	Key      string `xml:"key,attr"`
	EncIdx   string `xml:"enc_idx,attr"`
	MD5      string `xml:"md5,attr"`
	VideoMD5 string `xml:"videomd5,attr"`
	Value    string `xml:",chardata"`
}

// Thumb Thumb
type Thumb struct {
	Type   string `xml:"type,attr"`
	Token  string `xml:"token,attr"`
	Key    string `xml:"key,attr"`
	EncIdx string `xml:"enc_idx,attr"`
	Value  string `xml:",chardata"`
}

// Size Size
type Size struct {
	Width     string `xml:"width,attr,omitempty"`
	Height    string `xml:"height,attr,omitempty"`
	TotalSize string `xml:"totalSize,attr"`
}

// ActionInfo ActionInfo
type ActionInfo struct {
	AppMsg AppMsg `xml:"appMsg"`
}

// AppMsg AppMsg
type AppMsg struct {
	MessageAction string `xml:"messageAction"`
}

// Location Location
type Location struct {
	PoiClassifyID   string `xml:"poiClassifyId,attr"`
	PoiName         string `xml:"poiName,attr"`
	PoiAddress      string `xml:"poiAddress,attr"`
	PoiClassifyType uint32 `xml:"poiClassifyType,attr"`
	City            string `xml:"city,attr"`
	Latitude        string `xml:"latitude,attr"`
	Longitude       string `xml:"longitude,attr"`
}

// StreamVideo StreamVideo
type StreamVideo struct {
	StreamVideoURL      string `xml:"streamvideourl"`
	StreamVideoThumbURL string `xml:"streamvideothumburl"`
	StreamVideoWebURL   string `xml:"streamvideoweburl"`
}

// FriendTransItem 同步转发的朋友项
type FriendTransItem struct {
	FriendWXID   string
	FirstPageMd5 string
	CreateTime   uint32
	IsInited     bool
}
