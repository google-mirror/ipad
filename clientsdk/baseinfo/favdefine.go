package baseinfo

// FavItem 收藏项
type FavItem struct {
	Type             uint32           `xml:"type,attr"`
	CtrlFlag         uint32           `xml:"ctrlflag"`
	FavItemID        uint32           `xml:"favItemID"`
	Source           Source           `xml:"source"`
	Desc             string           `xml:"desc"`
	DataList         DataList         `xml:"datalist"`
	RecommendTagList RecommendTagList `xml:"recommendtaglist"`
}

// Source FavItem属性
type Source struct {
	SourceID   string `xml:"sourceid,attr"`
	SourceType uint32 `xml:"sourcetype,attr"`
	CreateTime uint32 `xml:"createtime"`
	FromUsr    string `xml:"fromusr"`
	EventID    string `xml:"eventid"`
}

// DataList 数据列表
type DataList struct {
	Count    uint32     `xml:"count,attr"`
	DataItem []DataItem `xml:"dataitem"`
}

// DataItem 数据项
type DataItem struct {
	HTMLID          string `xml:"htmlid,attr"`
	DataID          string `xml:"dataid,attr"`
	DataType        uint32 `xml:"datatype,attr"`
	SubType         uint32 `xml:"subtype,attr"`
	DataSourceID    string `xml:"datasourceid,attr"`
	DataIllegalType uint32 `xml:"dataillegaltype,attr"`
	ThumbFullSize   uint32 `xml:"thumbfullsize"`
	SourceThumbPath string `xml:"sourcethumbpath"`
	FullMD5         string `xml:"fullmd5"`
	ThumbHead256MD5 string `xml:"thumbhead256md5"`
	CdnThumbURL     string `xml:"cdn_thumburl"`
	SourceDataPath  string `xml:"sourcedatapath"`
	CdnDataKey      string `xml:"cdn_datakey"`
	FullSize        uint32 `xml:"fullsize"`
	Head256MD5      string `xml:"head256md5"`
	CdnThumbKey     string `xml:"cdn_thumbkey"`
	CdnDataURL      string `xml:"cdn_dataurl"`
	ThumbFullMD5    string `xml:"thumbfullmd5"`
	DataDesc        string `xml:"datadesc"`
	DataTitle       string `xml:"datatitle"`
}

// RecommendTagList 推荐标签列表
type RecommendTagList struct {
}

// CardInfo 名片信息
type CardInfo struct {
	BigHeadImgURL           string `xml:"bigheadimgurl,attr"`
	SmallHeadImgURL         string `xml:"smallheadimgurl,attr"`
	UserName                string `xml:"username,attr"`
	NickName                string `xml:"nickname,attr"`
	FullPY                  string `xml:"fullpy,attr"`
	ShortPY                 string `xml:"shortpy,attr"`
	Alias                   string `xml:"alias,attr"`
	ImageStatus             uint32 `xml:"imagestatus,attr"`
	Scene                   uint32 `xml:"scene,attr"`
	Province                string `xml:"province,attr"`
	City                    string `xml:"city,attr"`
	Sign                    string `xml:"sign,attr"`
	Sex                     uint32 `xml:"sex,attr"`
	CertFlag                uint32 `xml:"certflag,attr"`
	CertInfo                string `xml:"certinfo,attr"`
	BrandIconURL            string `xml:"brandIconUrl,attr"`
	BrandHomeURL            string `xml:"brandHomeUrl,attr"`
	BrandSubscriptConfigURL string `xml:"brandSubscriptConfigUrl,attr"`
	BrandFlags              uint32 `xml:"brandFlags,attr"`
	RegionCode              string `xml:"regionCode,attr"`
}
