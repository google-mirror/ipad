package baseinfo

// Msg Msg
type Msg struct {
	APPMsg       APPMsg `xml:"appmsg"`
	FromUserName string `xml:"fromusername"`
}

// APPMsg APPMsg
type APPMsg struct {
	Des       string    `xml:"des"`
	URL       string    `xml:"url"`
	MsgType   uint32    `xml:"type"`
	Title     string    `xml:"title"`
	ThumbURL  string    `xml:"thumburl"`
	WCPayInfo WCPayInfo `xml:"wcpayinfo"`
}

// WCPayInfo WCPayInfo
type WCPayInfo struct {
	TemplatedID   string `xml:"templatedid"`
	URL           string `xml:"url"`
	IconURL       string `xml:"iconurl"`
	ReceiverTitle string `xml:"receivertitle"`
	SenderTitle   string `xml:"sendertitle"`
	SenderDes     string `xml:"senderdes"`
	ReceiverDes   string `xml:"receiverdes"`
	NativeURL     string `xml:"nativeurl"`
	SceneID       uint32 `xml:"sceneid"`
	InnerType     string `xml:"innertype"`
	PayMsgID      string `xml:"paymsgid"`
	SceneText     string `xml:"scenetext"`
	LocalLogoIcon string `xml:"locallogoicon"`
	InvalidTime   uint32 `xml:"invalidtime"`
	Broaden       string `xml:"broaden"`
}

// HongBaoURLItem 红包NativeURL项
type HongBaoURLItem struct {
	SendUserName   string
	ShowWxPayTitle string
	MsgType        string
	ChannelID      string
	SendID         string
	Ver            string
	Sign           string
	ShowSourceMac  string
}

// HongBaoReceiverItem 接收红包项
type HongBaoReceiverItem struct {
	CgiCmd         uint32
	Province       string
	City           string
	InWay          uint32
	NativeURL      string
	HongBaoURLItem *HongBaoURLItem
}

// HongBaoOpenItem 领取红包项
type HongBaoOpenItem struct {
	CgiCmd           uint32
	Province         string
	City             string
	HeadImg          string
	NativeURL        string
	NickName         string
	SessionUserName  string
	TimingIdentifier string
	Offset           int64
	Limit            int64
	HongBaoURLItem   *HongBaoURLItem
}

// HongBaoQryDetailItem 查询红包领取详情
type HongBaoQryDetailItem struct {
	CgiCmd         uint32
	Province       string
	City           string
	NativeURL      string
	HongBaoURLItem *HongBaoURLItem
}

// HongBaoQryListItem 查询领取的红包列表信息
type HongBaoQryListItem struct {
	CgiCmd   uint32
	Province string
	City     string
	Offset   uint32
	Limit    uint32
}

// SourceObject SourceObject
type SourceObject struct {
	CoverImage     string
	CoverImageMd5  string
	DetailImage    string
	DetailImageMd5 string
}

// ShowSourceRec ShowSourceRec
type ShowSourceRec struct {
	SubType      uint32
	SourceObject SourceObject
}

// HongBaoReceiverResp 接收红包响应项
type HongBaoReceiverResp struct {
	RetCode                 uint32
	RetMsg                  string
	SendID                  string
	Wishing                 string
	IsSender                uint32
	ReceiveStatus           uint32
	HBStatus                uint32
	StatusMess              string
	HBType                  uint32
	WaterMark               string
	ScenePicSwitch          uint32
	PreStrainFlag           uint32
	SendUserName            string
	TimingIdentifier        string
	ShowSourceRec           ShowSourceRec
	ShowYearExpression      uint32
	ExpressionMd5           string
	ShowRecNormalExpression uint32
}

// RecordItem 领取红包记录
type RecordItem struct {
	ReceiveAmount uint32
	ReceiveTime   string
	Answer        string
	ReceiveID     string
	State         uint32
	ReceiveOpenID string
	UserName      string
}

// RealnameInfo RealnameInfo
type RealnameInfo struct {
	GuideFlag uint32
}

// Operation Operation
type Operation struct {
	Name    string
	Type    string
	Content string
	Enable  uint32
	IconURL string
	OssKey  uint32
}

// ShowSourceOpen ShowSourceOpen
type ShowSourceOpen struct {
	Source    ShowSourceRec
	Operation Operation
}

// HongBaoOpenResp 打开红包响应项
type HongBaoOpenResp struct {
	RetCode                  uint32
	RetMsg                   string
	SendID                   string
	Amount                   uint32
	RecNum                   uint32
	RecAmount                uint32
	TotalNum                 uint32
	HasWriteAnswer           uint32
	HBType                   uint32
	IsSender                 uint32
	IsContinue               uint32
	ReceiveStatus            uint32
	HBStatus                 uint32
	StatusMess               string
	Wishing                  string
	ReceiveID                string
	HeadTitle                string
	CanShare                 uint32
	OperationHeader          []string
	Record                   []RecordItem
	WaterMark                string
	JumpChange               uint32
	ChangeWording            string
	SendUserName             string
	RealnameInfo             RealnameInfo
	SystemMsgContext         string
	SessionUserName          string
	JumpChangeType           uint32
	ChangeIconURL            string
	ShowSourceOpen           ShowSourceOpen
	ExpressionMd5            string
	ExpressionType           uint32
	ShowYearExpression       uint32
	ShowOpenNormalExpression uint32
	EnableAnswerByExpression uint32
	EnableAnswerBySelfie     uint32
}

// OperationTail OperationTail
type OperationTail struct {
	Enable uint32
}

// AtomicFunc AtomicFunc
type AtomicFunc struct {
	Enable uint32
}

// HongBaoQryDetailResp 红包领取详情响应
type HongBaoQryDetailResp struct {
	RetCode                    uint32
	RetMsg                     string
	RecNum                     uint32
	TotalNum                   uint32
	TotalAmount                uint32
	SendID                     string
	Amount                     uint32
	Wishing                    string
	IsSender                   uint32
	ReceiveID                  string
	HasWriteAnswer             uint32
	OperationHeader            []string
	HBType                     uint32
	IsContinue                 uint32
	HBStatus                   uint32
	ReceiveStatus              uint32
	StatusMess                 string
	HeadTitle                  string
	CanShare                   uint32
	HBKind                     uint32
	RecAmount                  uint32
	Record                     []RecordItem
	OperationTail              OperationTail
	AtomicFunc                 AtomicFunc
	JumpChange                 uint32
	ChangeWording              string
	SendUserName               string
	ChangeURL                  string
	JumpChangeType             uint32
	ShowSourceOpen             ShowSourceOpen
	ExpressionMd5              string
	ShowDetailNormalExpression uint32
	EnableAnswerByExpression   uint32
	EnableAnswerBySelfie       uint32
}

// QryRecordItem 领取过的红包项
type QryRecordItem struct {
	SendName      string
	ReceiveAmount uint32
	ReceiveTime   string
	HBType        uint32
	SendID        string
	HBKind        uint32
	SendUserName  string
	ReceiveID     string
}

// HongBaoQryListResp 查询领取的红包记录
type HongBaoQryListResp struct {
	RetCode        uint32
	RetMsg         string
	RecTotalNum    uint32
	RecTotalAmount uint32
	Years          string
	GameCount      uint32
	RecordYear     string
	IsContinue     uint32
	Record         []QryRecordItem
}

type RedPacket struct {
	RedType  uint32
	Username string
	From     uint32
	Count    uint32
	Amount   uint32
	Content  string
}

// HongBaoItem 红包项
type HongBaoItem struct {
	NativeURL string
	Limit     int64
	URLItem   *HongBaoURLItem
}

type GetRedPacketList struct {
	Offset      int64
	Limit       int64
	NativeURL   string
	HongBaoItem HongBaoURLItem
}

type CreatePreTransfer struct {
	ToUserName  string
	Fee         uint
	Description string
}

type ConfirmPreTransfer struct {
	BankType    string
	BankSerial  string
	ReqKey      string
	PayPassword string
}
