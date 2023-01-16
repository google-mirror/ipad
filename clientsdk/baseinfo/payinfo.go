package baseinfo

// TenPayResp 响应
type TenPayResp struct {
	RetCode                 string                    `json:"retcode"`
	RetMsg                  string                    `json:"retmsg"`
	BindQueryScene          string                    `json:"bind_query_scene"`
	QueryCacheTime          uint32                    `json:"query_cache_time"`
	Array                   []Array                   `json:"Array"`
	VirtualCardArray        []VirtualCardArray        `json:"virtual_card_array"`
	UserInfo                TenPayUserInfo            `json:"user_info"`
	SwitchInfo              SwitchInfo                `json:"switch_info"`
	BalanceInfo             BalanceInfo               `json:"balance_info"`
	HistoryCardArray        []HistoryCardArray        `json:"history_card_array"`
	BalanceNotice           []BalanceNotice           `json:"balance_notice"`
	FetchNotice             []FetchNotice             `json:"fetch_notice"`
	QueryOrderTime          uint32                    `json:"query_order_time"`
	TimeStamp               uint32                    `json:"time_stamp"`
	PayMenuArray            []PayMenuArray            `json:"paymenu_array"`
	PayMenuUseNew           uint32                    `json:"paymenu_use_new"`
	WalletInfo              WalletInfo                `json:"wallet_info"`
	FavorComposeChannelInfo []FavorComposeChannelInfo `json:"favor_compose_channel_info"`
}

type GeneratePayQCodeResp struct {
	PayUrl  string `json:"pay_url"`
	RetCode string `json:"retcode"`
	RetMsg  string `json:"retmsg"`
}

type GetRedDetailsResp struct {
	RetCode                 string `json:"retcode"`
	RetMsg                  string `json:"retmsg"`
	SendId                  string `json:"sendId"`
	Wishing                 string `json:"wishing"`
	IsSender                string `json:"isSender"`
	ReceiveStatus           string `json:"receiveStatus"`
	HbStatus                string `json:"hbStatus"`
	StatusMess              string `json:"statusMess"`
	HbType                  string `json:"hbType"`
	Watermark               string `json:"watermark"`
	ScenePicSwitch          string `json:"scenePicSwitch"`
	PreStrainFlag           string `json:"preStrainFlag"`
	SendUserName            string `json:"sendUserName"`
	TimingIdentifier        string `json:"timingIdentifier"`
	ShowYearExpression      string `json:"showYearExpression"`
	ExpressionMd5           string `json:"expression_md5"`
	ShowRecNormalExpression string `json:"showRecNormalExpression"`
}

// Array Array
type Array struct {
	BankFlag        string `json:"bank_flag"`
	BankName        string `json:"bank_name"`
	BankType        string `json:"bank_type"`
	BindSerial      string `json:"bind_serial"`
	BankaccTypeName string `json:"bankacc_type_name"`
}

// VirtualCardArray VirtualCardArray
type VirtualCardArray struct {
}

// HistoryCardArray HistoryCardArray
type HistoryCardArray struct {
}

// BalanceNotice BalanceNotice
type BalanceNotice struct {
}

// FetchNotice FetchNotice
type FetchNotice struct {
}

// PayMenuArray PayMenuArray
type PayMenuArray struct {
}

// FavorComposeChannelInfo struct {
type FavorComposeChannelInfo struct {
}

// TouchInfo TouchInfo
type TouchInfo struct {
	ISOpenTouch string `json:"is_open_touch"`
	UseTouchPay string `json:"use_touch_pay"`
}

// TenPayUserInfo TenPayUserInfo
type TenPayUserInfo struct {
	ISReg              string    `json:"is_reg"`
	TrueName           string    `json:"true_name"`
	BindCardNum        string    `json:"bind_card_num"`
	ICardUserFlag      string    `json:"icard_user_flag"`
	CreName            string    `json:"cre_name"`
	CreType            string    `json:"cre_type"`
	TransferURL        string    `json:"transfer_url"`
	TouchInfo          TouchInfo `json:"touch_info"`
	LctWording         string    `json:"lct_wording"`
	LctURL             string    `json:"lct_url"`
	AuthenChannelState uint64    `json:"authen_channel_state"`
}

// SwitchInfo SwitchInfo
type SwitchInfo struct {
	SwitchBit uint32 `json:"switch_bit"`
}

// BalanceInfo BalanceInfo
type BalanceInfo struct {
	UseCftBalance     string `json:"use_cft_balance"`
	BalanceBankType   string `json:"balance_bank_type"`
	BalanceBindSerial string `json:"balance_bind_serial"`
	TotalBalance      string `json:"total_balance"`
	AvailBalance      string `json:"avail_balance"`
	FrozenBalance     string `json:"frozen_balance"`
	FetchBalance      string `json:"fetch_balance"`
	Mobile            string `json:"mobile"`
	SupportMicropay   string `json:"support_micropay"`
	BalanceListURL    string `json:"balance_list_url"`
	BalanceVersion    uint32 `json:"balance_version"`
	TimeOut           uint32 `json:"time_out"`
	BalanceLogoURL    string `json:"balance_logo_url"`
}

// WalletInfo WalletInfo
type WalletInfo struct {
	WalletBalance                    uint32 `json:"wallet_balance"`
	WalletEntranceBalanceWwitchState uint32 `json:"wallet_entrance_balance_switch_state"`
}

// BindQueryNewReq BindQueryNewReq
type BindQueryNewReq struct {
	BalanceVersion    uint32
	BindQueryScen     uint32
	BindTypeCond      string
	City              string
	ISDeviceOpenTouch uint32
	ISRoot            uint32
	Province          string
}

// TenPayReqItem TenPayReqItem
type TenPayReqItem struct {
	CgiCMD  uint32
	ReqText string
}

// PreTransferResp PreTransferResp
type PreTransferResp struct {
	ReqKey           string `json:"req_key"`
	TansferingStatus string `json:"tansfering_status"`
	RetCode          string `json:"retcode"`
	RetMsg           string `json:"retmsg"`
	//AppMsgContent     string `json:"appmsgcontent"`
	ReceiverTrueName string `json:"receiver_true_name"`
	TransferId       string `json:"transfer_id"`
	TransactionId    string `json:"transaction_id"`
	Fee              uint   `json:"fee"`
}
