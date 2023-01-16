package table

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	//发送文本消息
	MYSQL_BUSINESS_TYPE_SENDTEXTMSG = "SendTextMessage"
	//发送图片信息
	MYSQL_BUSINESS_TYPE_SENDIMGMSG = "SendImageMessage"
)

type LocalTime struct {
	time.Time
}

// MarshalJSON on LocalTime format Time field with %Y-%m-%d %H:%M:%S
func (t LocalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *LocalTime) Scan(v interface{}) error {
	switch v.(type) {
	case []byte:
		timeBytes, ok := v.([]byte)
		if ok {
			todayZero, _ := time.ParseInLocation("2006-01-02 15:04:05", string(timeBytes), time.Local)
			*t = LocalTime{Time: todayZero}
			return nil
		}
	case time.Time:
		value, ok := v.(time.Time)
		if ok {
			*t = LocalTime{Time: value}
			return nil
		}

	}

	return fmt.Errorf("can not convert %v to timestamp", v)
}

func NewLocalTime() *LocalTime {
	return &LocalTime{time.Now()}
}

type CdnSnsImageInfo struct {
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
	ImageMD5    string `gorm:"primary_key"`
	ImageWidth  uint32
	ImageHeight uint32
}

type MysqlBase struct {
	TargetIp string `gorm:"column:targetIp"`
}

// USerBusinessLog 用户行为日志
type UserBusinessLog struct {
	Id            uint   `gorm:"primary_key,AUTO_INCREMENT"`               //自增Id
	UUID          string `gorm:"column:uuid" json:"uuid"`                  //用户链接Id
	UserName      string `gorm:"column:user_name" json:"userName"`         //登录的WXID
	BusinessType  string `gorm:"column:business_type" json:"businessType"` //业务类型 所调用的接口
	ExecuteResult string `gorm:"column:ex_result" json:"executeResult"`    //执行结果
}

//用户登录日志
type UserLoginLog struct {
	MysqlBase
	Id        uint `gorm:"primary_key,AUTO_INCREMENT" json:"id"`
	UUId      string
	UserName  string
	NickName  string
	LoginType string
	UpdatedAt LocalTime `json:"loginTime"`
	RetCode   int32     `gorm:"column:ret_code"`
	ErrMsg    string    `gorm:"column:err_msg;type:text"`
}

// UserInfoEntity 用户信息
type UserInfoEntity struct {
	MysqlBase
	UUID         string `gorm:"column:uuid" json:"uuid"`
	Uin          uint32 `gorm:"column:uin" json:"uin"`
	WxId         string `gorm:"column:wxId;primary_key" json:"wxId"`
	NickName     string `gorm:"column:nickname" json:"nickname"`
	UserName     string `gorm:"column:userName" json:"user_name"`
	Password     string `gorm:"column:password" json:"password"`
	HeadURL      string `gorm:"column:headurl" json:"headurl"`
	Session      []byte `gorm:"column:cookie" json:"cookie"`
	SessionKey   []byte `gorm:"column:sessionKey" json:"sessionKey"`
	ShortHost    string `gorm:"column:shorthost" json:"shorthost"`
	LongHost     string `gorm:"column:longhost" json:"longhost"`
	EcPublicKey  []byte `gorm:"column:ecpukey" json:"ecpukey"`
	EcPrivateKey []byte `gorm:"column:ecprkey" json:"ecprkey"`
	CheckSumKey  []byte `gorm:"column:checksumkey" json:"checksumkey"`
	AutoAuthKey  string `gorm:"column:autoauthkey;type:varchar(2048)" json:"autoauthkey"`
	State        int32  `gorm:"column:state" json:"state"`
	SyncKey      string `gorm:"column:synckey;type:varchar(1024)" json:"synckey"`
	FavSyncKey   string `gorm:"column:favsynckey;type:varchar(100)" json:"favsynckey"`
	// 登录的Rsa 密钥版本
	LoginRsaVer uint32
	ErrMsg      string `gorm:"type:text"`
}

// DeviceInfoEntity 设备信息
type DeviceInfoEntity struct {
	WxId               string `gorm:"column:wxid;primary_key" json:"wxid"`
	UUIDOne            string `gorm:"column:uuidone" json:"uuidone"`
	UUIDTwo            string `gorm:"column:uuidtwo" json:"uuidtwo"`
	Imei               string `gorm:"column:imei" json:"imei"`
	DeviceID           []byte `gorm:"column:deviceid" json:"deviceid"`
	DeviceName         string `gorm:"column:devicename" json:"devicename"`
	TimeZone           string `gorm:"column:timezone" json:"timezone"`
	Language           string `gorm:"column:language" json:"language"`
	DeviceBrand        string `gorm:"column:devicebrand" json:"devicebrand"`
	RealCountry        string `gorm:"column:realcountry" json:"realcountry"`
	IphoneVer          string `gorm:"column:iphonever" json:"iphonever"`
	BundleID           string `gorm:"column:boudleid" json:"boudleid"`
	OsType             string `gorm:"column:ostype" json:"ostype"`
	AdSource           string `gorm:"column:adsource" json:"adsource"`
	OsTypeNumber       string `gorm:"column:ostypenumber" json:"ostypenumber"`
	CoreCount          uint32 `gorm:"column:corecount" json:"corecount"`
	CarrierName        string `gorm:"column:carriername" json:"carriername"`
	SoftTypeXML        string `gorm:"column:softtypexml;type:varchar(2048)" json:"softtypexml"`
	ClientCheckDataXML string `gorm:"column:clientcheckdataxml;type:varchar(4096)" json:"clientcheckdataxml"`
	// extInfo
	GUID2 string `json:"GUID2"`
}
