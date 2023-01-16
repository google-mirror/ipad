package clientsdk

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"feiyu.com/wx/clientsdk/android"
	"feiyu.com/wx/clientsdk/android/mmproto"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/ccdata"
	"feiyu.com/wx/clientsdk/extinfo"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
)

// GetClientCheckData 获取ClientCheckData
func GetClientCheckData(userInfo *baseinfo.UserInfo) []byte {
	retData, err := ccdata.CreateClientCheckData(userInfo.DeviceInfo.ClientCheckDataXML)
	if err != nil {
		log.Info(err)
		return retData
	}
	// DeviceRunningInfo
	deviceRuntimInfo := &wechat.DeviceRunningInfo{}
	mode := uint32(1)
	deviceRuntimInfo.Mode = &mode
	deviceRuntimInfo.Type = []byte("00000006")
	deviceRuntimInfo.Data = retData

	// 序列化
	finalData, err := proto.Marshal(deviceRuntimInfo)
	if err != nil {
		log.Info(err)
	}
	return finalData
}

// GetExtSpamInfoData 获取ExtSpamInfo
func GetExtSpamInfoData(userInfo *baseinfo.UserInfo) []byte {
	ccdata := GetClientCheckData(userInfo)

	// WCExtInfoLod
	extInfo := &wechat.WCExtInfo{}
	// ccDataBuffer
	var ccDataBuffer wechat.SKBuiltinString_
	tmpLen := uint32(len(ccdata))
	ccDataBuffer.Len = &tmpLen
	ccDataBuffer.Buffer = ccdata
	extInfo.CcData = &ccDataBuffer

	// 序列化
	retData, err := proto.Marshal(extInfo)
	if err != nil {
		log.Info(err)
	}
	return retData
}

func GetExtPBSpamInfoDataA16(userInfo *baseinfo.UserInfo, wxId ...string) []byte {
	wxId_ := ""
	if len(wxId) == 0 {
		wxId_ = userInfo.GetUserName()
	} else {
		wxId_ = wxId[0]
	}
	ccd1 := GetCCD1(wxId_)
	ccd1PB, _ := proto.Marshal(ccd1)
	ccd2 := GetCCD2()
	ccd2PB, _ := proto.Marshal(ccd2)

	ccd3 := GetCCD3(*userInfo.DeviceInfoA16)
	ccd3PB, _ := proto.Marshal(ccd3)
	devicetoken := GetDeviceToken(userInfo.DeviceInfoA16.DeviceIdStr)
	dtPB, _ := proto.Marshal(devicetoken)
	spamdatabody := &mmproto.SpamDataBody{
		Ccd1: &mmproto.SpamDataSubBody{
			Ilen:   proto.Uint32(uint32(len(ccd1PB))),
			Ztdata: ccd1,
		},
		Ccd2: &mmproto.SpamDataSubBody{
			Ilen:   proto.Uint32(uint32(len(ccd2PB))),
			Ztdata: ccd2,
		},
		Ccd3: &mmproto.SpamDataSubBody{
			Ilen:   proto.Uint32(uint32(len(ccd3PB))),
			Ztdata: ccd3,
		},
		Dt: &mmproto.DeviceTokenBody{
			Ilen:        proto.Uint32(uint32(len(dtPB))),
			DeviceToken: devicetoken,
		},
	}
	retData, err := proto.Marshal(spamdatabody)
	if err != nil {
		log.Info(err)
	}
	return retData
}

func GetCCD1(UserName string) *mmproto.ZTData {

	curtime := uint64(time.Now().UnixNano() / 1e6)
	contentlen := len(UserName)
	var ct []uint64
	ut := curtime
	for i := 0; i < contentlen; i++ {
		ut += uint64(rand.Intn(10000))
		ct = append(ct, ut)
	}
	ccd := &mmproto.Ccd1{
		StartTime: &curtime,
		CheckTime: &curtime,
		Count:     proto.Uint32(uint32(contentlen)),
		EndTime:   ct,
	}
	pb, _ := proto.Marshal(ccd)

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(pb)
	Ztdata := &mmproto.ZTData{
		Version:   proto.String("00000006\x00"),
		Encrypted: proto.Uint32(1),
		Data:      encData,
		TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
		Optype:    proto.Uint32(5),
		Uin:       proto.Uint32(0),
	}

	return Ztdata
}

func GetCCD2() *mmproto.ZTData {

	curtime := uint32(time.Now().Unix())
	curNanoTime := uint64(time.Now().UnixNano())

	ccd := &mmproto.Ccd2{
		Checkid:   proto.String("<LoginByID>"),
		StartTime: &curtime,
		CheckTime: &curtime,
		Count1:    proto.Uint32(0),
		Count2:    proto.Uint32(1),
		Count3:    proto.Uint32(0),
		Const1:    proto.Uint64(384214787666497617),
		Const2:    &curNanoTime,
		Const3:    &curNanoTime,
		Const4:    &curNanoTime,
		Const5:    &curNanoTime,
		Const6:    proto.Uint64(384002236977512448),
	}
	pb, _ := proto.Marshal(ccd)

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(pb)
	Ztdata := &mmproto.ZTData{
		Version:   proto.String("00000006\x00"),
		Encrypted: proto.Uint32(1),
		Data:      encData,
		TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
		Optype:    proto.Uint32(5),
		Uin:       proto.Uint32(0),
	}

	return Ztdata
}

// 算法
func GetCCD3(accoutInfo baseinfo.AndroidDeviceInfo) *mmproto.ZTData {

	curtime := uint32(time.Now().Unix())

	ccd3body := &mmproto.Ccd3Body{
		KernelReleaseNumber: proto.String(accoutInfo.KernelReleaseNumber),
		UsbState:            proto.Uint32(0),
		Sign:                proto.String(accoutInfo.PackageSign),
		PackageFlag:         proto.Uint32(14),
		AccessFlag:          proto.Uint32(364604),
		Unkonwn:             proto.Uint32(3),
		TbVersionCrc:        proto.Uint32(553983350),
		SfMD5:               proto.String("d001b450158a85142c953011c66d531d"),
		SfArmMD5:            proto.String("bf7f84d081f1dffd587803c233d4e235"),
		SfArm64MD5:          proto.String("85801b3939f277ad31c9f89edd9dd008"),
		SbMD5:               proto.String("683e7beb7a44017ca2e686e3acedfb9f"),
		SoterId2:            proto.String(""),
		TimeCheck:           proto.Uint32(0),
		NanoTime:            proto.Uint32(455583),
	}

	pb, _ := proto.Marshal(ccd3body)

	crc := crc32.ChecksumIEEE(pb)

	ccd3 := &mmproto.Ccd3{
		Crc:       &crc,
		TimeStamp: &curtime,
		Body:      nil,
	}

	pb, _ = proto.Marshal(ccd3)

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(pb)

	Ztdata := &mmproto.ZTData{
		Version:   proto.String("00000006\x00"),
		Encrypted: proto.Uint32(1),
		Data:      encData,
		TimeStamp: &curtime,
		Optype:    proto.Uint32(5),
		Uin:       proto.Uint32(0),
	}
	return Ztdata
}

func GetDeviceToken(devicetoken string) *mmproto.DeviceToken {
	curtime := uint32(time.Now().Unix())
	return &mmproto.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mmproto.SKBuiltinStringt{
			String_: &devicetoken,
		},
		TimeStamp: &curtime,
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
}

func GetExtPBSpamInfoData(userInfo *baseinfo.UserInfo, wxId ...string) []byte {
	wxId_ := ""
	if len(wxId) == 0 {
		wxId_ = userInfo.GetUserName()
	} else {
		wxId_ = wxId[0]
	}

	retData, err := extinfo.GetCCDPbLib(
		userInfo.DeviceInfo.OsTypeNumber,
		userInfo.DeviceInfo.OsType,
		userInfo.DeviceInfo.UUIDTwo,
		userInfo.DeviceInfo.UUIDTwo,
		userInfo.DeviceInfo.DeviceName,
		userInfo.DeviceInfo.DeviceToken,
		hex.EncodeToString(userInfo.DeviceInfo.DeviceID),
		wxId_,
		userInfo.DeviceInfo.GUID2,
		userInfo,
	)
	if err != nil {
		log.Info(err)
	}
	return retData
}

// GetAutoAuthRsaReqDataMarshal 生成自动登陆rsareq项
func GetAutoAuthRsaReqDataMarshal(userInfo *baseinfo.UserInfo) []byte {
	var rsaReqData wechat.AutoAuthRsaReqData

	// AesEncyptKey
	var aesEncryptKey wechat.SKBuiltinString_
	aesEncryptKey.Buffer = []byte(userInfo.SessionKey)
	var aesKeyLen = uint32(len(userInfo.SessionKey))
	aesEncryptKey.Len = &aesKeyLen
	rsaReqData.AesEncryptKey = &aesEncryptKey

	// ecdh
	var ecdhKey wechat.ECDHKey
	var tmpNid = uint32(713)
	ecdhKey.Nid = &tmpNid
	// ecdhKey
	var ecdhKeyBuffer wechat.SKBuiltinString_
	var ecdhKeyLen = uint32(len(userInfo.EcPublicKey))
	ecdhKeyBuffer.Len = &ecdhKeyLen
	ecdhKeyBuffer.Buffer = userInfo.EcPublicKey
	ecdhKey.Key = &ecdhKeyBuffer
	rsaReqData.PubEcdhKey = &ecdhKey

	retData, err := proto.Marshal(&rsaReqData)
	if err != nil {
		log.Info("proto.Marshal AutoAuthRsaReqData failed: ", err)
	}
	return retData
}

// GetAutoAuthAesReqDataMarshal 生成自动登陆aesreq项
func GetAutoAuthAesReqDataMarshal(userInfo *baseinfo.UserInfo) []byte {
	var zeroUint32 = uint32(0)
	var emptyString = string("")
	var aesReqData wechat.AutoAuthAesReqData
	baseReq := GetBaseRequest(userInfo)
	var tmpScene uint32 = 2
	baseReq.Scene = &tmpScene
	baseReq.SessionKey = []byte{}
	aesReqData.BaseRequest = baseReq

	// autoauthkey
	tmpAuthKeyLen := uint32(len(userInfo.AutoAuthKey))
	var tmpAutoAuthKey wechat.SKBuiltinString_
	tmpAutoAuthKey.Buffer = userInfo.AutoAuthKey
	tmpAutoAuthKey.Len = &tmpAuthKeyLen
	aesReqData.AutoAuthKey = &tmpAutoAuthKey

	// BaseAuthReqInfo
	var baseReqInfo wechat.BaseAuthReqInfo
	baseReqInfo.AuthReqFlag = &zeroUint32
	baseReqInfo.AuthTicket = &emptyString
	aesReqData.BaseReqInfo = &baseReqInfo

	if userInfo.DeviceInfo == nil {
		//return GetManualAuthAesReqDataA16Protobuf(userInfo)
		aesReqData.Imei = proto.String(userInfo.DeviceInfoA16.AndriodImei(userInfo.DeviceInfoA16.DeviceIdStr))
		aesReqData.SoftType = proto.String(userInfo.DeviceInfoA16.AndriodGetSoftType(userInfo.DeviceInfoA16.DeviceIdStr))
		aesReqData.ClientSeqId = proto.String(fmt.Sprintf("%s_%d", userInfo.DeviceInfoA16.DeviceIdStr, (time.Now().UnixNano() / 1e6)))
		aesReqData.DeviceName = proto.String(userInfo.DeviceInfoA16.AndroidManufacturer(userInfo.DeviceInfoA16.DeviceIdStr) + "-" + userInfo.DeviceInfoA16.AndroidPhoneModel(userInfo.DeviceInfoA16.DeviceIdStr))
		aesReqData.Language = proto.String("Zh")
		//aesReqData.Language = proto.String(userInfo.DeviceInfo.Language)
		aesReqData.TimeZone = proto.String("8.0")
		aesReqData.Channel = &zeroUint32
		// TimeStamp
		aesReqData.Signature = proto.String(userInfo.DeviceInfoA16.AndriodPackageSign(userInfo.DeviceInfoA16.DeviceIdStr))
		aesReqData.BuiltinIpSeq = &zeroUint32
		ext, err := GetExtSpamInfoAndroid(userInfo)
		if err != nil {
			log.Error("Android extSpam err", err.Error())
		}
		// extSpamInfo
		var extSpamInfo wechat.SKBuiltinString_
		extSpamInfo.Buffer = ext
		extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
		extSpamInfo.Len = &extSpamInfoLen
		aesReqData.ExtSpamInfo = &extSpamInfo
		reqData, _ := proto.Marshal(&aesReqData)
		return reqData
	} else {
		// imei
		aesReqData.Imei = &userInfo.DeviceInfo.Imei
		aesReqData.TimeZone = &userInfo.DeviceInfo.TimeZone
		aesReqData.DeviceName = &userInfo.DeviceInfo.DeviceName
		aesReqData.Language = &userInfo.DeviceInfo.Language
		if userInfo.LoginDataInfo.Language != "" {
			aesReqData.Language = proto.String(userInfo.LoginDataInfo.Language)
		}
		aesReqData.BuiltinIpSeq = &zeroUint32
		aesReqData.Signature = &emptyString
		aesReqData.SoftType = &userInfo.DeviceInfo.SoftTypeXML

		tmpTime := int(time.Now().UnixNano() / 1000000000)
		tmpTimeStr := strconv.Itoa(tmpTime)
		var strClientSeqID = string(userInfo.DeviceInfo.Imei + "-" + tmpTimeStr)
		aesReqData.ClientSeqId = &strClientSeqID
		aesReqData.Channel = &zeroUint32

		// extSpamInfo
		var extSpamInfo wechat.SKBuiltinString_
		extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
		extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
		extSpamInfo.Len = &extSpamInfoLen
		aesReqData.ExtSpamInfo = &extSpamInfo

		retData, err := proto.Marshal(&aesReqData)
		if err != nil {
			log.Info("proto.Marshal AutoAuthAesReqData failed: ", err)
		}

		return retData
	}
}

// GetManualAuthAesReqDataMarshal 生成自动登陆aesreq项
func GetManualAuthAesReqDataMarshal(userInfo *baseinfo.UserInfo) []byte {
	zeroUint32 := uint32(0)
	zeroInt32 := int32(0)
	emptyString := string("")
	var aesRequest wechat.ManualAuthAesReqData
	baseReq := GetBaseRequest(userInfo)
	var tmpScene uint32 = 1
	baseReq.Scene = &tmpScene
	baseReq.SessionKey = []byte{}
	// ClientSeqId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	var strClientSeqID = string(userInfo.DeviceInfo.Imei + "-" + tmpTimeStr)
	// TimeStamp
	tmpTime2 := uint32(time.Now().UnixNano() / 1000000000)
	// extSpamInfo
	var extSpamInfo wechat.SKBuiltinString_
	extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	aesRequest = wechat.ManualAuthAesReqData{
		BaseRequest:  baseReq,
		Imei:         &userInfo.DeviceInfo.Imei,
		TimeZone:     &userInfo.DeviceInfo.TimeZone,
		DeviceName:   &userInfo.DeviceInfo.DeviceName,
		DeviceType:   &userInfo.DeviceInfo.DeviceName,
		Channel:      &zeroInt32,
		BuiltinIpseq: &zeroUint32,
		Signature:    &emptyString,
		SoftType:     &userInfo.DeviceInfo.SoftTypeXML,
		DeviceBrand:  &userInfo.DeviceInfo.DeviceBrand,
		RealCountry:  &userInfo.DeviceInfo.RealCountry,
		BundleId:     &userInfo.DeviceInfo.BundleID,
		AdSource:     &userInfo.DeviceInfo.AdSource,
		InputType:    proto.Uint32(uint32(2)),
		ClientSeqId:  &strClientSeqID,
		TimeStamp:    &tmpTime2,
		ExtSpamInfo:  &extSpamInfo,
		Language:     &userInfo.DeviceInfo.Language,
	}
	// BaseRequest
	if userInfo.Ticket != "" {
		aesRequest.BaseReqInfo = &wechat.BaseAuthReqInfo{
			AuthTicket: proto.String(userInfo.Ticket),
		}
		aesRequest.InputType = proto.Uint32(1)
	}
	retData, err := proto.Marshal(&aesRequest)
	if err != nil {
		log.Info("proto.Marshal AutoAuthAesReqData failed: ", err)
	}
	return retData
}

// ParseHongBaoURL 解析红包URL
func ParseHongBaoURL(hongBaoURL string, senderUserName string) (*baseinfo.HongBaoURLItem, error) {
	retURL, err := url.Parse(hongBaoURL)
	if err != nil {
		return nil, err
	}

	retMaps, err := url.ParseQuery(retURL.RawQuery)
	if err != nil {
		return nil, err
	}

	retItem := &baseinfo.HongBaoURLItem{}
	retItem.MsgType = retMaps["msgtype"][0]
	retItem.ChannelID = retMaps["channelid"][0]
	retItem.SendID = retMaps["sendid"][0]
	retItem.Ver = retMaps["ver"][0]
	retItem.Sign = retMaps["sign"][0]
	retItem.SendUserName = senderUserName
	return retItem, nil
}

// CreateMediaItemXML 将mediaItem转换成XML
func CreateMediaItemXML(mediaItem *baseinfo.SnsMediaItem) string {
	retString := string("<media>")
	if mediaItem.EncKey != "" {
		retString = retString + "<enc key=\"" + mediaItem.EncKey + "\">" + strconv.Itoa(int(mediaItem.EncValue)) + "</enc>"
	}
	retString = retString + "<id>" + strconv.Itoa(int(mediaItem.ID)) + "</id>"
	retString = retString + "<type>" + strconv.Itoa(int(mediaItem.Type)) + "</type>"

	// title
	titleString := "<title/>"
	if len(mediaItem.Title) > 0 {
		titleString = "<title>" + mediaItem.Title + "</title>"
	}
	retString = retString + titleString

	// description
	descriptionString := "<description/>"
	if len(mediaItem.Description) > 0 {
		descriptionString = "<description>" + mediaItem.Description + "</description>"
	}
	retString = retString + descriptionString

	// private
	retString = retString + "<private>" + strconv.Itoa(int(mediaItem.Private)) + "</private>"

	// userData
	userDataString := "<userData/>"
	if len(mediaItem.UserData) > 0 {
		userDataString = "<userData>" + mediaItem.UserData + "</userData>"
	}
	retString = retString + userDataString

	// subType
	retString = retString + "<subType>" + strconv.Itoa(int(mediaItem.SubType)) + "</subType>"

	// videoSize
	retString = retString + "<videoSize width=\"" + mediaItem.VideoWidth + "\" height=\"" + mediaItem.VideoHeight + "\"/>"

	// url
	retString = retString + "<url type=\"" + mediaItem.URLType + "\" md5=\"" + mediaItem.MD5 + "\" videomd5=\"" + mediaItem.VideoMD5 + "\">"
	retString = retString + mediaItem.URL + "</url>"

	// thumb
	retString = retString + "<thumb type=\"" + mediaItem.ThumType + "\">"
	retString = retString + mediaItem.Thumb + "</thumb>"

	// size
	retString = retString + "<size width=\"" + mediaItem.SizeWidth + "\" height=\"" + mediaItem.SizeHeight + "\" totalSize=\"" + mediaItem.TotalSize + "\"/>"

	// videoDuration 如果是视频
	if mediaItem.Type == baseinfo.MMSNSMediaTypeVideo {
		// 格式化
		tmpValue := strconv.FormatFloat(mediaItem.VideoDuration, 'f', 6, 64)
		retString = retString + "<videoDuration>" + tmpValue + "</videoDuration>"
	}

	retString = retString + "</media>"
	return retString
}

// CreateSnsPostItemXML 转成xml字符串，字节数组
func CreateSnsPostItemXML(userName string, postItem *baseinfo.SnsPostItem) []byte {
	// createTime
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)

	// start
	retString := string("<TimelineObject>")
	retString = retString + "<id><![CDATA[0]]></id>"
	retString = retString + "<username><![CDATA[" + userName + "]]></username>"
	retString = retString + "<createTime><![CDATA[" + tmpTimeStr + "]]></createTime>"
	retString = retString + "<contentDesc><![CDATA[" + postItem.Content + "]]></contentDesc>"
	retString = retString + "<contentDescShowType><![CDATA[0]]></contentDescShowType>"
	retString = retString + "<contentDescScene><![CDATA[" + strconv.Itoa(0) + "]]></contentDescScene>"
	retString = retString + "<private>" + strconv.Itoa(int(postItem.Privacy)) + "</private>"
	retString = retString + "<sightFolded>" + strconv.Itoa(0) + "</sightFolded>"
	retString = retString + "<showFlag>" + strconv.Itoa(0) + "</showFlag>"
	// location

	// appInfo
	retString = retString + "<appInfo>"
	retString = retString + "<id/>"
	retString = retString + "<version/>"
	retString = retString + "<appName/>"
	retString = retString + "<installUrl/>"
	retString = retString + "<fromUrl/>"
	retString = retString + "<isForceUpdate>" + strconv.Itoa(0) + "</isForceUpdate>"
	retString = retString + "</appInfo>"

	retString = retString + "<sourceUserName/>"
	retString = retString + "<sourceNickName/>"
	retString = retString + "<statisticsData/>"
	retString = retString + "<statExtStr/>"
	// ContentObject
	/*retString = retString + "<ContentObject>"
	retString = retString + "<contentStyle>" + strconv.Itoa(int(postItem.ContentStyle)) + "</contentStyle>"
	retString = retString + "<title>&#x0A;&#x0A;&#x0A;习近平<title/>"
	retString = retString + "<description/>"*/

	// location
	if postItem.LocationInfo != nil {
		// float, int 转成string
		longitudeStr := postItem.LocationInfo.Longitude
		latitudeStr := postItem.LocationInfo.Latitude
		poiScaleStr := strconv.Itoa(int(postItem.LocationInfo.PoiScale))
		poiClassfyTypeStr := strconv.Itoa(int(postItem.LocationInfo.PoiClassifyType))
		poiClickableStatusStr := strconv.Itoa(int(postItem.LocationInfo.PoiClickableStatus))
		// 增加 LocationXml
		retString = retString + "<location "
		retString = retString + "city = \"" + postItem.LocationInfo.City + "\" "
		retString = retString + "longitude = \"" + longitudeStr + "\" "
		retString = retString + "latitude = \"" + latitudeStr + "\" "
		retString = retString + "poiName = \"" + postItem.LocationInfo.PoiName + "\" "
		retString = retString + "poiAddress = \"" + postItem.LocationInfo.PoiAddress + "\" "
		retString = retString + "poiScale = \"" + poiScaleStr + "\" "
		retString = retString + "poiInfoUrl = \"" + postItem.LocationInfo.PoiInfoURL + "\" "
		retString = retString + "poiClassifyId = \"" + postItem.LocationInfo.PoiClassifyID + "\" "
		retString = retString + "poiClassifyType = \"" + poiClassfyTypeStr + "\" "
		retString = retString + "poiClickableStatus = \"" + poiClickableStatusStr + "\" "
		retString = retString + "></location>"
	}

	retString = retString + "<ContentObject>"
	if postItem.ContentUrl != "" {
		retString = retString + "<contentUrl><![CDATA[" + postItem.ContentUrl + "]]></contentUrl>"
	}
	retString = retString + "<contentStyle><![CDATA[" + strconv.Itoa(int(postItem.ContentStyle)) + "]]></contentStyle>"
	if len(postItem.MediaList) > 0 && postItem.ContentUrl != "" {
		retString = retString + "<title><![CDATA[" + postItem.MediaList[0].Title + "]]>&#x0A;&#x0A;&#x0A;习近平--习大大</title>" //习近平
	} else {
		retString = retString + "<title>&#x0A;&#x0A;&#x0A;习近平;习大大</title>" //习近平
	}
	if postItem.Description != "" {
		retString = retString + "<description>![CDATA[" + postItem.Description + "]]></description>"
	} else {
		retString = retString + "<description></description>"
	}
	// mediaList
	mediaListString := "<mediaList/>"
	mediaCount := len(postItem.MediaList)
	if mediaCount > 0 {
		mediaListString = "<mediaList>"
		for index := 0; index < mediaCount; index++ {
			mediaListString = mediaListString + CreateMediaItemXML(postItem.MediaList[index])
		}
		mediaListString = mediaListString + "</mediaList>"
	}
	retString = retString + mediaListString
	retString = retString + "</ContentObject>"
	// end
	retString = retString + "</TimelineObject>"
	log.Debug(retString)
	return []byte(retString)
}

// CreateSnsMediaItem 创建媒体项
// privacy：公开/不公开
// description：描述
func CreateSnsMediaItem(snsImgResponse *baseinfo.CdnSnsImageUploadResponse, privacy uint32, description string) *baseinfo.SnsMediaItem {
	mediaItem := &baseinfo.SnsMediaItem{}
	mediaItem.ID = 0
	mediaItem.Type = baseinfo.MMSNSMediaTypeImage
	mediaItem.Title = ""
	mediaItem.Description = description
	mediaItem.Private = privacy
	mediaItem.UserData = ""
	mediaItem.SubType = 0
	mediaItem.URL = snsImgResponse.FileURL
	mediaItem.URLType = "1"
	mediaItem.Thumb = snsImgResponse.ThumbURL
	mediaItem.ThumType = "1"
	mediaItem.MD5 = snsImgResponse.ImageMD5
	mediaItem.VideoMD5 = ""
	mediaItem.VideoWidth = strconv.Itoa(int(snsImgResponse.ImageWidth))
	mediaItem.VideoHeight = strconv.Itoa(int(snsImgResponse.ImageHeight))
	tmpWidth := float64(snsImgResponse.ImageWidth)
	tmpHeight := float64(snsImgResponse.ImageHeight)
	mediaItem.SizeWidth = strconv.FormatFloat(tmpWidth, 'f', 6, 64)
	mediaItem.SizeHeight = strconv.FormatFloat(tmpHeight, 'f', 6, 64)
	mediaItem.TotalSize = "0"

	return mediaItem
}

// CreateSnsCommentLikeItem 创建评论项：朋友圈点赞
func CreateSnsCommentLikeItem(itemID uint64, toUserName string) *baseinfo.SnsCommentItem {
	retItem := &baseinfo.SnsCommentItem{}
	retItem.OpType = baseinfo.MMSnsCommentTypeLike
	retItem.ItemID = itemID
	retItem.ToUserName = toUserName
	retItem.CreateTime = uint32(time.Now().UnixNano() / 1000000000)
	retItem.Content = ""
	retItem.ReplyCommentID = 0
	retItem.ReplyItem = nil

	return retItem
}

// CreateSnsCommentItem 创建评论项：朋友圈评论
func CreateSnsCommentItem(itemID uint64, toUserName string, content string, replyComment *wechat.SnsCommentInfo) *baseinfo.SnsCommentItem {
	retItem := &baseinfo.SnsCommentItem{}
	retItem.OpType = baseinfo.MMSnsCommentTypeComment
	retItem.ItemID = itemID
	retItem.ToUserName = toUserName
	retItem.CreateTime = uint32(time.Now().UnixNano() / 1000000000)
	retItem.Content = content
	retItem.ReplyCommentID = 0
	retItem.ReplyItem = nil

	// 回复
	if replyComment != nil {
		retItem.ReplyItem = &baseinfo.ReplyCommentItem{}
		retItem.ReplyItem.UserName = replyComment.GetUsername()
		retItem.ReplyItem.NickName = replyComment.GetNickname()
		retItem.ReplyCommentID = replyComment.GetCommentId()
		retItem.ReplyItem.OpType = replyComment.GetType()
		retItem.ReplyItem.Source = replyComment.GetSource()
	}
	return retItem
}

// CreateSnsLocationInfo 创建朋友圈项: 发送朋友圈地址
func CreateSnsLocationInfo(lbsLife *wechat.LbsLife, cityName string, longitude string, latitude string) *baseinfo.SnsLocationInfo {
	retLocationInfo := &baseinfo.SnsLocationInfo{}

	retLocationInfo.City = cityName
	retLocationInfo.Longitude = longitude
	retLocationInfo.Latitude = latitude
	retLocationInfo.PoiName = lbsLife.GetTitle()
	retLocationInfo.PoiAddress = ""
	retLocationInfo.PoiScale = 11.0
	retLocationInfo.PoiInfoURL = lbsLife.GetPoiUrl()
	retLocationInfo.PoiClassifyID = lbsLife.GetBid()
	retLocationInfo.PoiClassifyType = 1
	if lbsLife.GetType() == 1 {
		retLocationInfo.PoiClassifyType = 2
	}
	retLocationInfo.PoiClickableStatus = 0
	return retLocationInfo
}

// CreateGetLbsLifeListItem 创建GetLbsLifeListItem项
func CreateGetLbsLifeListItem(longitude float64, latitude float64, buff []byte, keyWord string) *baseinfo.GetLbsLifeListItem {
	// 获取Lbs地址列表
	lbsLifeListItem := &baseinfo.GetLbsLifeListItem{}
	lbsLifeListItem.Opcode = baseinfo.MMLbsLifeOpcodeNormal
	lbsLifeListItem.Buffer = buff
	lbsLifeListItem.Longitude = float32(longitude)
	lbsLifeListItem.Latitude = float32(latitude)
	lbsLifeListItem.KeyWord = keyWord

	return lbsLifeListItem
}

// CreateCDNUploadMsgImgPrepareRequest 创建 CreateCDNUploadMsgImgPrepareRequest请求
// prepareRequestItem: 请求项
func CreateCDNUploadMsgImgPrepareRequest(userInfo *baseinfo.UserInfo, prepareRequestItem *baseinfo.CDNUploadMsgImgPrepareRequestItem) []byte {
	var request wechat.CDNUploadMsgImgPrepareRequest
	var emptyString = string("")
	var zeroValue32 = int32(0)

	// FromUserName
	request.FromUserName = &userInfo.WxId

	// ToUserName
	request.ToUserName = &prepareRequestItem.ToUser

	// ClientImgId
	clientImgID := prepareRequestItem.ToUser + prepareRequestItem.LocalName + "_" + strconv.Itoa(int(prepareRequestItem.CreateTime))
	request.ClientImgId = &clientImgID

	// ThumbWidth
	request.ThumbWidth = &prepareRequestItem.ThumbWidth

	// ThumbHeight
	request.ThumbHeight = &prepareRequestItem.ThumbHeight

	// encryVer
	request.EncryVer = &zeroValue32

	// Scene
	request.Scene = &zeroValue32

	// Crc32
	request.Crc32 = &prepareRequestItem.Crc32

	// Aeskey
	aesKey := baseutils.BytesToHexString(prepareRequestItem.AesKey, false)
	request.Aeskey = &aesKey

	// MsgForwardType
	msgForwardType := uint32(1)
	request.MsgForwardType = &msgForwardType

	// AttachedContent
	request.AttachedContent = &emptyString

	// Longitude
	longitude := float32(0.0)
	request.Longitude = &longitude
	// Latitude
	request.Latitude = &longitude

	// Source
	source := uint32(2)
	request.Source = &source

	// Appid
	request.Appid = &emptyString

	// MessageAction
	request.MessageAction = &emptyString

	// MessageExt
	request.MeesageExt = &emptyString

	// MediaTagName
	request.MediaTagName = &emptyString

	// 打包
	srcData, _ := proto.Marshal(&request)
	retData := Pack(userInfo, srcData, 625, 5)

	return retData
}

// CreateSendEmojiMsgXMl 生成表情xml
func CreateSendEmojiMsgXMl(md5 string, totallen int32) string {
	retString := strings.Builder{}
	retString.WriteString("<appmsg appid=\"\"  sdkver=\"0\">")
	retString.WriteString("<title></title>")
	retString.WriteString("<des></des>")
	retString.WriteString("<action></action>")
	retString.WriteString("<type>8</type>") // 发送表情 type：8
	retString.WriteString("<showtype>0</showtype>")
	retString.WriteString("<soundtype>0</soundtype>")
	retString.WriteString("<mediatagname></mediatagname>")
	retString.WriteString("<messageext></messageext>")
	retString.WriteString("<messageaction></messageaction>")
	retString.WriteString("<content></content>")
	retString.WriteString("<contentattr>0</contentattr>")
	retString.WriteString("<url></url>")
	retString.WriteString("<lowurl></lowurl>")
	retString.WriteString("<dataurl></dataurl>")
	retString.WriteString("<lowdataurl></lowdataurl>")
	retString.WriteString("<songalbumurl></songalbumurl>")
	retString.WriteString("<songlyric></songlyric>")
	retString.WriteString("<appattach>")
	retString.WriteString(fmt.Sprintf("<totallen>%d</totallen>", totallen))
	retString.WriteString(fmt.Sprintf("<attachid>0:0:%s</attachid>", md5))
	retString.WriteString(fmt.Sprintf("<emoticonmd5>%s</emoticonmd5>", md5))
	retString.WriteString("<fileext>pic</fileext>")
	retString.WriteString("<cdnthumbaeskey></cdnthumbaeskey>" +
		"<aeskey></aeskey>" +
		"</appattach>" +
		"<extinfo></extinfo>" +
		"<sourceusername></sourceusername>" +
		"<sourcedisplayname></sourcedisplayname>" +
		"<thumburl></thumburl>" +
		"<md5></md5>" +
		"<statextstr></statextstr>" +
		"<directshare>0</directshare>" +
		"</appmsg>" +
		"<fromusername></fromusername>")
	return retString.String()
}
