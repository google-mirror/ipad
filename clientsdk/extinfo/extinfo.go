package extinfo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math/rand"
	"strings"
	"time"

	"feiyu.com/wx/clientsdk/android/mmproto"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/ccdata"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/go-resty/resty/v2"
	"github.com/gogf/guuid"
	"github.com/gogo/protobuf/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type GetCcDataRep struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

// 临时接口获取CCDATA
func GetCcData(deviceToken *wechat.TrustResp, deviceId string) (string, error) {
	client := resty.New()
	rep := new(GetCcDataRep)
	argsSend := make(map[string]interface{})
	println("devicetoken:%v", deviceToken.GetTrustResponseData().GetDeviceToken())
	argsSend["deviceToken"] = deviceToken.GetTrustResponseData().GetDeviceToken()
	argsSend["deviceId"] = deviceId
	_, err := client.R().SetBody(argsSend).
		SetHeader("Content-Type", "application/json").
		SetResult(rep).
		Post("http://127.0.0.1:8023/ccdata")
	if err != nil {
		return "", err
	}
	return rep.Data, nil
}

// 获取CCD
func GetCCDPbLib(iosVersion, deviceType, uuid1, uuid2, deviceName string, deviceToken *wechat.TrustResp, deviceId, userName, guid2 string, userInfo *baseinfo.UserInfo) ([]byte, error) {
	ccData1, err := GetNewSpamData(iosVersion, deviceType, uuid1, uuid2, deviceName, deviceToken, deviceId, userName, guid2, userInfo)
	println("原计算值：%v", hex.EncodeToString(ccData1))
	if err != nil {
		return nil, err
	}
	wcstf, err := GetWcstf(userName)
	if err != nil {
		return nil, err
	}
	wcste, err := GetWcste()
	if err != nil {
		return nil, err
	}
	deviceTokenObj := GetDeviceToken(deviceToken.GetTrustResponseData().GetDeviceToken())
	dt, err := proto.Marshal(deviceTokenObj)
	BehaviorId := []byte("login_" + GenGUId(deviceId, GetCid(0x098521236654)))

	if err != nil {
		return nil, err
	}
	wcExtInfo := &wechat.WCExtInfoNew{
		Wcstf: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(wcstf))),
			Buffer: wcstf,
		},
		Wcste: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(wcste))),
			Buffer: wcste,
		},
		CcData: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(ccData1))),
			Buffer: ccData1,
		},
		DeviceToken: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(dt))),
			Buffer: dt,
		},
		BehaviorId: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(BehaviorId))),
			Buffer: BehaviorId,
		},
	}
	return proto.Marshal(wcExtInfo)
}

// 705以上版本的计算extinfo
func GetNewSpamData(iosVersion string, deviceType, uuid1, uuid2, deviceName string, deviceToken *wechat.TrustResp, deviceId, userName, guid2 string, userInfo *baseinfo.UserInfo) ([]byte, error) {
	dateTimeSramp := time.Now().Unix()
	timeStamp := time.Now().Unix()
	xorKey := MakeXorKey(timeStamp)
	//xorKey := byte((dateTimeSramp * 0xffffffed) + 7)
	guid1 := guuid.New().String()
	if guid2 == "" {
		guid2 = guuid.New().String()
	}
	Lang := "zh"
	Country := "CN"
	if userInfo.DeviceInfo != nil {
		Lang = userInfo.DeviceInfo.Language
		Country = Lang
	}
	//guid2 := guuid.New().String()
	//24算法
	spamDataBody := wechat.SpamDataBody{
		UnKnown1:              proto.Int32(1),
		TimeStamp:             proto.Uint32(uint32(dateTimeSramp)),
		KeyHash:               proto.Int32(int32(MakeKeyHash(xorKey))),
		Yes1:                  proto.String(XorEncrypt("yes", xorKey)),
		Yes2:                  proto.String(XorEncrypt("yes", xorKey)),
		IosVersion:            proto.String(XorEncrypt(iosVersion, xorKey)),
		DeviceType:            proto.String(XorEncrypt(deviceType, xorKey)),
		UnKnown2:              proto.Int32(4), //cpu核数
		IdentifierForVendor:   proto.String(XorEncrypt(uuid1, xorKey)),
		AdvertisingIdentifier: proto.String(XorEncrypt(uuid2, xorKey)),
		Carrier:               proto.String(XorEncrypt("中国联通", xorKey)),
		BatteryInfo:           proto.Int32(1),
		NetworkName:           proto.String(XorEncrypt("en0", xorKey)),
		NetType:               proto.Int32(1),
		AppBundleId:           proto.String(XorEncrypt("com.tencent.xin", xorKey)),
		DeviceName:            proto.String(XorEncrypt(deviceName, xorKey)),
		UserName:              proto.String(XorEncrypt(userName, xorKey)),
		Unknown3:              proto.Int64(DeviceNumber(deviceId[:29] + "FFF")), //基带版本   77968568550229002
		Unknown4:              proto.Int64(DeviceNumber(deviceId[:29] + "OOO")), //基带通讯版本   77968568550228991
		Unknown5:              proto.Int32(0),                                   //IsJailbreak
		Unknown6:              proto.Int32(4),
		Lang:                  proto.String(XorEncrypt(Lang, xorKey)),    //zh
		Country:               proto.String(XorEncrypt(Country, xorKey)), //CN
		Unknown7:              proto.Int32(4),
		DocumentDir:           proto.String(XorEncrypt(fmt.Sprintf("/var/mobile/Containers/Data/Application/%s/Documents", guid1), xorKey)),
		Unknown8:              proto.Int32(0),
		Unknown9:              proto.Int32(0),
		HeadMD5:               proto.String(XorEncrypt(GetCidMd5(deviceId, GetCid(0x0262626262626)), xorKey)), //XorEncrypt("d13610700984cf481b7d3f5fa2011c30", xorKey)
		AppUUID:               proto.String(XorEncrypt(uuid1, xorKey)),
		SyslogUUID:            proto.String(XorEncrypt(uuid2, xorKey)),
		Unknown10:             proto.String(XorEncrypt(BuildRandomWifiSsid(), xorKey)),
		Unknown11:             proto.String(XorEncrypt(baseutils.BuildRandomMac(), xorKey)),
		AppName:               proto.String(XorEncrypt("微信", xorKey)),
		SshPath:               proto.String(XorEncrypt("/usr/bin/ssh", xorKey)),
		TempTest:              proto.String(XorEncrypt("/tmp/test.txt", xorKey)), //XorEncrypt("/tmp/test.txt", xorKey)
		DevMD5:                proto.String(""),
		DevUser:               proto.String(""),
		DevPrefix:             proto.String(""),
		AppFileInfo:           GetFileInfo(deviceId, xorKey, guid2),
		Unknown12:             proto.String(""),
		IsModify:              proto.Int32(0),
		ModifyMD5:             proto.String(baseutils.Md5Value("modify")),
		//RqtHash:               proto.Int64(baseutils.EncInt(int64(baseutils.CalcMsgCrc([]byte(XorEncrypt(baseutils.Md5Value("modify"), xorKey)))))), //288529533794259264
		RqtHash:   proto.Int64(baseutils.EncInt(int64(baseutils.CalcMsgCrcForString_807(baseutils.Md5Value("modify" + deviceId))))),
		Unknown53: proto.Uint64(1586355322),            //微信安装时间
		Unknown54: proto.Uint64(uint64(dateTimeSramp)), //微信启动时间  1586355519000
		Unknown55: proto.Uint64(0),                     //固定值0
		Unknown56: proto.Int64(baseutils.EncInt(int64(baseutils.CalcMsgCrcForString_807(baseutils.Md5Value("modify" + deviceId))))),
		//Unknown46:             proto.Int64(baseutils.EncInt(int64(baseuti ls.CalcMsgCrc([]byte(XorEncrypt(baseutils.Md5Value("modify"), xorKey)))))), //288529533794259264
		Unknown57: proto.Uint64(0), //固定值0
		Unknown58: proto.String(XorEncrypt(deviceId, xorKey)),
		//Unknown49: proto.String(strconv.FormatInt(baseutils.EncInt(int64(crc32.ChecksumIEEE([]byte(XorEncrypt(deviceId, xorKey))))), 10)),
		Unknown59: proto.String(fmt.Sprintf("%v", baseutils.EncInt(int64(crc32.ChecksumIEEE([]byte(XorEncrypt(deviceId, xorKey))))))),
		Unknown61: proto.String(XorEncrypt("2FFC7F6DFEEFFF2B3FFCA029", xorKey)),
		Unknown62: proto.Uint64(1175744137544159509),
		Unknown63: proto.String(XorEncrypt(baseutils.Md5Value(deviceId+"imsi"), xorKey)), //device + "imsi"  md5
	}
	appFileInfo := new(bytes.Buffer)
	appFile := make([]string, 0)
	filePath := []string{
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/WeChat",
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/Frameworks/mars.framework/mars",
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/Frameworks/andromeda.framework/andromeda",
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/Frameworks/OpenSSL.framework/OpenSSL",
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/Frameworks/ProtobufLite.framework/ProtobufLite",
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/Frameworks/marsbridgenetwork.framework/marsbridgenetwork",
		"/var/containers/Bundle/Application/" + GenGUId(deviceId, GetCid(0x098521236654)) + "/WeChat.app/Frameworks/matrixreport.framework/matrixreport",
	}
	fileUUid := []string{
		"3A1D0388-6BDB-350C-8706-80E3D15AA7C7",
		"A7DA401B-3FF6-3920-A30A-1B0FA8258202",
		"10F1245A-68FD-310D-98B3-0CFD51760BDE",
		"8FAE149B-602B-3B9D-A620-88EA75CE153F",
		"05BD590C-4DF6-3EDB-8316-0C9783928DD0",
		"CFED9A03-C881-3D50-B014-732D0A09879F",
		"1E7F06D2-DD36-31A8-AF3B-00D62054E1F9",
	}

	for i := 0; i < 7; i++ {
		appFile = append(appFile, filePath[i])
		appFile = append(appFile, fileUUid[i])
	}
	appFileInfo.WriteString(strings.Join(appFile, ""))
	encInt := baseutils.EncInt(int64(crc32.ChecksumIEEE(appFileInfo.Bytes())))
	spamDataBody.Unknown60 = proto.String(fmt.Sprintf("%v", encInt))
	spamDataBody.Unknown64 = proto.String(XorEncrypt(deviceToken.GetTrustResponseData().GetDeviceToken(), xorKey))
	data, err := proto.Marshal(&spamDataBody)
	if err != nil {
		return nil, err
	}

	newClientCheckData := &wechat.NewClientCheckData{
		C32CData:  proto.Int64(int64(crc32.ChecksumIEEE(data))),
		TimeStamp: proto.Int64(time.Now().Unix()),
		DataBody:  data,
	}

	ccData, err := proto.Marshal(newClientCheckData)
	if err != nil {
		return nil, err
	}

	afterCompressionCCData := baseutils.CompressByteArray(ccData)
	afterEnData, err := ccdata.EncodeZipData(afterCompressionCCData, 0x3060)
	if err != nil {
		return nil, err
	}
	deviceRunningInfo := &wechat.DeviceRunningInfoNew{
		Version:     []byte("00000006"),
		Type:        proto.Uint32(1),
		EncryptData: afterEnData,
		Timestamp:   proto.Uint32(uint32(timeStamp)),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	return proto.Marshal(deviceRunningInfo)
}

// 获取FilePathCrc
func GetFileInfo(deviceId string, xorKey byte, guid2 string) []*wechat.FileInfo {
	return []*wechat.FileInfo{
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/WeChat", xorKey)),
			Fileuuid: proto.String(XorEncrypt("3A1D0388-6BDB-350C-8706-80E3D15AA7C7", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/Frameworks/mars.framework/mars", xorKey)),
			Fileuuid: proto.String(XorEncrypt("A7DA401B-3FF6-3920-A30A-1B0FA8258202", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/Frameworks/andromeda.framework/andromeda", xorKey)),
			Fileuuid: proto.String(XorEncrypt("10F1245A-68FD-310D-98B3-0CFD51760BDE", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/Frameworks/OpenSSL.framework/OpenSSL", xorKey)),
			Fileuuid: proto.String(XorEncrypt("8FAE149B-602B-3B9D-A620-88EA75CE153F", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/Frameworks/ProtobufLite.framework/ProtobufLite", xorKey)),
			Fileuuid: proto.String(XorEncrypt("05BD590C-4DF6-3EDB-8316-0C9783928DD0", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/Frameworks/marsbridgenetwork.framework/marsbridgenetwork", xorKey)),
			Fileuuid: proto.String(XorEncrypt("CFED9A03-C881-3D50-B014-732D0A09879F", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/"+GenGUId(deviceId, GetCid(0x098521236654))+"/WeChat.app/Frameworks/matrixreport.framework/matrixreport", xorKey)),
			Fileuuid: proto.String(XorEncrypt("1E7F06D2-DD36-31A8-AF3B-00D62054E1F9", xorKey)),
		},
	}
}

func DeviceNumber(DeviceId string) int64 {
	ssss := []byte(baseutils.Md5Value(DeviceId))
	ccc := Hex2int(&ssss) >> 8
	ddd := ccc + 60000000000000000
	if ddd > 80000000000000000 {
		ddd = ddd - (80000000000000000 - ddd)
	}
	return int64(ddd)
}

func Hex2int(hexB *[]byte) uint64 {
	var retInt uint64
	hexLen := len(*hexB)
	for k, v := range *hexB {
		retInt += b2m_map[v] * exponent(16, uint64(2*(hexLen-k-1)))
	}
	return retInt
}

func exponent(a, n uint64) uint64 {
	result := uint64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}

func GenGUId(DeviceId, Cid string) string {
	Md5Data := baseutils.Md5Value(DeviceId + Cid)
	return fmt.Sprintf("%x-%x-%x-%x-%x", Md5Data[0:8], Md5Data[2:6], Md5Data[3:7], Md5Data[1:5], Md5Data[20:32])
}

func GetWcstf(userName string) ([]byte, error) {
	curtime := time.Now().UnixNano() / 1000000
	contentLen := len(userName)
	ct := make([]uint64, 0)
	ut := curtime
	for i := 0; i < contentLen; i++ {
		ut += rand.Int63n(10000)
		ct = append(ct, uint64(ut))
	}

	wcstf := &wechat.WCSTF{
		StartTime: proto.Uint64(uint64(curtime)),
		CheckTime: proto.Uint64(uint64(curtime)),
		Count:     proto.Uint32(uint32(contentLen)),
		EndTime:   ct,
	}
	ccData, err := proto.Marshal(wcstf)
	if err != nil {
		return nil, err
	}

	afterCompressionCCData := baseutils.CompressByteArray(ccData)
	afterEnData, err := ccdata.EncodeZipData(afterCompressionCCData, 0x3060)
	if err != nil {
		return nil, err
	}

	deviceRunningInfo := &wechat.DeviceRunningInfoNew{
		Version:     []byte("00000006"),
		Type:        proto.Uint32(1),
		EncryptData: afterEnData,
		Timestamp:   proto.Uint32(uint32(curtime)),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	return proto.Marshal(deviceRunningInfo)
}

func MakeXorKey(key int64) uint8 {
	var un int64 = int64(0xffffffed)
	xorKey := (uint8)(key*un + 7)
	return xorKey
}

func GetWcste() ([]byte, error) {
	curtime := time.Now().Unix()
	curNanoTime := time.Now().UnixNano() / 1000000
	wcste := &wechat.WCSTE{
		CheckId:   proto.String("<LoginByID>"),
		StartTime: proto.Uint32(uint32(curtime)),
		CheckTime: proto.Uint32(uint32(curtime)),
		Count1:    proto.Uint32(0),
		Count2:    proto.Uint32(1),
		Count3:    proto.Uint32(0),
		Const1:    proto.Uint64(384214787666497617),
		Const2:    proto.Uint64(uint64(curNanoTime)),
		Const3:    proto.Uint64(uint64(curNanoTime)),
		Const4:    proto.Uint64(uint64(curNanoTime)),
		Const5:    proto.Uint64(uint64(curNanoTime)),
		Const6:    proto.Uint64(384002236977512448),
	}
	ccData, err := proto.Marshal(wcste)
	if err != nil {
		return nil, err
	}

	afterCompressionCCData := baseutils.CompressByteArray(ccData)
	afterEnData, err := ccdata.EncodeZipData(afterCompressionCCData, 0x3060)
	if err != nil {
		return nil, err
	}

	deviceRunningInfo := &wechat.DeviceRunningInfoNew{
		Version:     []byte("00000006"),
		Type:        proto.Uint32(1),
		EncryptData: afterEnData,
		Timestamp:   proto.Uint32(uint32(curtime)),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	return proto.Marshal(deviceRunningInfo)
}

// 获取DeviceToken
func GetDeviceToken(deviceToken string) *mmproto.DeviceToken {
	curtime := uint32(time.Now().Unix())
	return &mmproto.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mmproto.SKBuiltinStringt{
			String_: proto.String(deviceToken),
		},
		TimeStamp: &curtime,
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
}

var wifiPrefix = []string{"TP_", "360_", "ChinaNet-", "MERCURY_", "DL-", "VF_", "HUAW-"}

func BuildRandomWifiSsid() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	i := r.Intn(len(wifiPrefix))
	randChar := make([]byte, 6)
	for x := 0; x < 6; x++ {
		randChar[x] = byte(r.Intn(26) + 65)
	}
	return wifiPrefix[i] + string(randChar)
}
