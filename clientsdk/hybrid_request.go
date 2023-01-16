package clientsdk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"feiyu.com/wx/clientsdk/android"
	"feiyu.com/wx/clientsdk/android/mmproto"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/cecdh"
	"feiyu.com/wx/clientsdk/extinfo"
	clientsdk "feiyu.com/wx/clientsdk/hybrid"
	"feiyu.com/wx/clientsdk/mmtls"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/gogo/protobuf/proto"
	"github.com/lunny/log"
)

type SecLoginKeyMgr struct {
	WeChatPubKeyVersion byte
	WeChatPubKey        string `json:"-"`
	SourceData          []byte `json:"-"`
	PriKey              []byte `json:"-"`
	PubKey              []byte `json:"-"`
	FinalSha256         []byte `json:"-"`
}

func (Sec *SecLoginKeyMgr) Reset() {
	Sec.PubKey = []byte{}
	Sec.PriKey = []byte{}
	Sec.FinalSha256 = []byte{}
	Sec.SourceData = []byte{}
}

func (Sec *SecLoginKeyMgr) SetKey() {
	if 146 == Sec.WeChatPubKeyVersion {
		Sec.WeChatPubKeyVersion = 145
		Sec.WeChatPubKey = clientsdk.WeChatPubKey_145
	} else {
		Sec.WeChatPubKeyVersion = 146
		Sec.WeChatPubKey = clientsdk.WeChatPubKey_146
	}
}
func NewSecLoginKeyMgrByVer(ver byte) *SecLoginKeyMgr {
	sec := &SecLoginKeyMgr{}
	switch ver {
	case 146:
		sec.WeChatPubKeyVersion = 146
		sec.WeChatPubKey = clientsdk.WeChatPubKey_146
		break
	case 145:
		sec.WeChatPubKeyVersion = 145
		sec.WeChatPubKey = clientsdk.WeChatPubKey_145
		break
	}
	return sec
}

func NewSecLoginKeyMgr() *SecLoginKeyMgr {
	return &SecLoginKeyMgr{
		WeChatPubKeyVersion: 146,
		WeChatPubKey:        clientsdk.WeChatPubKey_146,
	}
}

// GetManualAuthAccountProtobuf 组用户登录基本信息
func GetManualAuthRsaReqDataProtobuf(userInfo *baseinfo.UserInfo, wxid string, newpass string) *wechat.ManualAuthRsaReqData {
	var tmpNid uint32 = 713
	userInfo.EcPublicKey, userInfo.EcPrivateKey = cecdh.GenerateEccKey()
	authRequest := &wechat.ManualAuthRsaReqData{}
	// aes_key
	var aesKey wechat.SKBuiltinString_
	var tmpAesKeyLen uint32 = 16
	aesKey.Len = &tmpAesKeyLen
	aesKey.Buffer = []byte(userInfo.SessionKey)
	authRequest.RandomEncryKey = &aesKey

	// 其它参数
	var ecdhKey wechat.ECDHKey
	var key wechat.SKBuiltinString_
	key.Buffer = userInfo.EcPublicKey
	var tmpLen = (uint32)(len(userInfo.EcPublicKey))
	key.Len = &tmpLen
	ecdhKey.Nid = &tmpNid
	ecdhKey.Key = &key
	authRequest.CliPubEcdhkey = &ecdhKey
	authRequest.UserName = &wxid
	//判断是否为iPad登录的伪密码
	if !strings.HasPrefix(newpass, "extdevnewpwd_") && !strings.HasPrefix(newpass, "strdm@") {
		newpass = baseutils.Md5Value(newpass)
	}
	authRequest.Pwd = &newpass
	return authRequest
}

// GetManualAuthAesReqProtobuf 生成自动登陆aesreq项
func GetManualAuthAesReqDataProtobuf(userInfo *baseinfo.UserInfo) *wechat.ManualAuthAesReqData {
	// if userInfo.DeviceInfoA16 != nil {
	// 	return GetManualAuthAesReqDataProtobufA16(userInfo)
	// }
	zeroUint32 := uint32(0)
	zeroInt32 := int32(0)
	emptyString := string("")

	var aesRequest wechat.ManualAuthAesReqData
	// BaseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene uint32 = 1
	baseReq.Scene = &tmpScene
	baseReq.SessionKey = []byte{}
	baseReq.Uin = proto.Uint32(0)
	aesRequest.BaseRequest = baseReq
	inputType := uint32(2)
	aesRequest.InputType = &inputType
	aesRequest.BaseReqInfo = &wechat.BaseAuthReqInfo{}
	if userInfo.Ticket != "" {
		aesRequest.BaseReqInfo = &wechat.BaseAuthReqInfo{
			AuthTicket: proto.String(userInfo.Ticket),
		}
		aesRequest.InputType = proto.Uint32(1)
	}

	// imei
	aesRequest.Imei = &userInfo.DeviceInfo.Imei
	aesRequest.TimeZone = &userInfo.DeviceInfo.TimeZone
	aesRequest.DeviceName = &userInfo.DeviceInfo.DeviceName
	aesRequest.DeviceType = &userInfo.DeviceInfo.DeviceName
	aesRequest.Channel = &zeroInt32
	aesRequest.Language = &userInfo.DeviceInfo.Language
	aesRequest.BuiltinIpseq = &zeroUint32
	aesRequest.Signature = &emptyString
	aesRequest.SoftType = &userInfo.DeviceInfo.SoftTypeXML
	aesRequest.DeviceBrand = &userInfo.DeviceInfo.DeviceBrand
	aesRequest.RealCountry = &userInfo.DeviceInfo.RealCountry
	aesRequest.BundleId = &userInfo.DeviceInfo.BundleID
	aesRequest.AdSource = &userInfo.DeviceInfo.AdSource

	// ClientSeqId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	var strClientSeqID = string(userInfo.DeviceInfo.Imei + "-" + tmpTimeStr)
	aesRequest.ClientSeqId = &strClientSeqID

	// TimeStamp
	tmpTime2 := uint32(time.Now().UnixNano() / 1000000000)
	aesRequest.TimeStamp = &tmpTime2

	// extSpamInfo
	var extSpamInfo wechat.SKBuiltinString_
	extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	aesRequest.ExtSpamInfo = &extSpamInfo

	return &aesRequest
}

type WCExtInfo struct {
	CcData      wechat.SKBuiltinString_
	DeviceToken wechat.SKBuiltinString_
	BehaviorID  wechat.SKBuiltinString_
}

// 新修改的安全登陸
func SendHybridManualAutoRequest(userInfo *baseinfo.UserInfo, newpass string, wxid string, keyVer byte) (*baseinfo.PackHeader, error) {
	if keyVer == 0 {
		keyVer = 146
	} else if keyVer == 145 {
		userInfo.HybridLogin = true
	}
	secKeyMgr := NewSecLoginKeyMgrByVer(keyVer)

	req := &wechat.ManualAuthRequest{
		RsaReqData: GetManualAuthRsaReqDataProtobuf(userInfo, wxid, newpass),
		AesReqData: GetManualAuthAesReqDataProtobuf(userInfo),
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	//加密
	encrypt, epKey, token, ecdhpairkey, err := clientsdk.HybridEncrypt(reqData, secKeyMgr.WeChatPubKey)
	if err != nil {
		return nil, err
	}
	ecdhPacket := &wechat.EcdhPacket{
		Type: proto.Uint32(1),
		Key: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: ecdhpairkey.PubKey,
		},
		Token:        token,
		Url:          proto.String(""),
		ProtobufData: encrypt,
	}

	secKeyMgr.PubKey = ecdhpairkey.PubKey
	secKeyMgr.PriKey = ecdhpairkey.PriKey
	secKeyMgr.SourceData = reqData
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, epKey[24:]...)
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, reqData...)

	ecdhDataPacket, err := proto.Marshal(ecdhPacket)
	if err != nil {
		return nil, err
	}

	packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, 252, ecdhDataPacket, ecdhDataPacket, uint32(len(ecdhDataPacket)), 12, uint32(0x4e))
	//设置Hybrid 加密密钥版本
	packHeader.HybridKeyVer = secKeyMgr.WeChatPubKeyVersion
	//开始组头
	retData := PackHeaderSerialize(packHeader, false)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secmanualauth", retData)
	if err != nil {
		return nil, err
	}
	packHeader, err = DecodePackHeader(resp, nil)
	if err != nil {
		return nil, err
	}
	packHeader.Data, err = clientsdk.HybridEcdhDecrypt(packHeader.Data, secKeyMgr.PriKey, secKeyMgr.PubKey, secKeyMgr.FinalSha256)
	if err != nil {
		return nil, err
	}
	return packHeader, err
}

func SendHybridManualAutoRequest000(userInfo *baseinfo.UserInfo, newpass string, wxid string, keyVer byte) (*baseinfo.PackHeader, error) {
	if keyVer == 0 {
		keyVer = 146
	} else if keyVer == 145 {
		userInfo.HybridLogin = true
	}
	secKeyMgr := NewSecLoginKeyMgrByVer(keyVer)

	req := &wechat.ManualAuthRequest{
		RsaReqData: GetManualAuthRsaReqDataProtobuf(userInfo, wxid, newpass),
		AesReqData: GetManualAuthAesReqDataProtobuf(userInfo),
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	//加密
	encrypt, epKey, token, ecdhpairkey, err := clientsdk.HybridEncrypt(reqData, secKeyMgr.WeChatPubKey)
	if err != nil {
		return nil, err
	}
	ecdhPacket := &wechat.EcdhPacket{
		Type: proto.Uint32(1),
		Key: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: ecdhpairkey.PubKey,
		},
		Token:        token,
		Url:          proto.String(""),
		ProtobufData: encrypt,
	}

	secKeyMgr.PubKey = ecdhpairkey.PubKey
	secKeyMgr.PriKey = ecdhpairkey.PriKey
	secKeyMgr.SourceData = reqData
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, epKey[24:]...)
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, reqData...)

	ecdhDataPacket, err := proto.Marshal(ecdhPacket)
	if err != nil {
		return nil, err
	}

	packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, 252, ecdhDataPacket, ecdhDataPacket, uint32(len(ecdhDataPacket)), 12, uint32(0x4e))
	//设置Hybrid 加密密钥版本
	packHeader.HybridKeyVer = secKeyMgr.WeChatPubKeyVersion
	//开始组头
	retData := PackHeaderSerialize(packHeader, false)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secmanualauth", retData)
	if err != nil {
		return nil, err
	}
	packHeader, err = DecodePackHeader(resp, nil)
	if err != nil {
		return nil, err
	}
	packHeader.Data, err = clientsdk.HybridEcdhDecrypt(packHeader.Data, secKeyMgr.PriKey, secKeyMgr.PubKey, secKeyMgr.FinalSha256)
	if err != nil {
		return nil, err
	}
	return packHeader, err
}

func SendHybridManualAutoRequest2(userInfo *baseinfo.UserInfo, newpass string, wxid string, keyVer byte) (*android.PacketHeader, error) {
	if keyVer == 0 {
		keyVer = 146
	} else if keyVer == 145 {
		userInfo.HybridLogin = true
	}
	//secKeyMgr := NewSecLoginKeyMgrByVer(keyVer)

	req := &wechat.ManualAuthRequest{
		RsaReqData: GetManualAuthRsaReqDataProtobuf(userInfo, wxid, newpass),
		AesReqData: GetManualAuthAesReqDataProtobuf(userInfo),
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	hec := &android.Client{}
	hec.Init("IOS", int(baseinfo.ClientVersion), baseinfo.DeviceTypeIos)
	hecData := hec.HybridEcdhPackIosEn(252, 0, nil, reqData)
	recvData, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secmanualauth", hecData)
	if err != nil {
		return nil, err
	}
	ph := hec.HybridEcdhPackIosUn(recvData)
	return ph, nil

	//加密
	/*encrypt, epKey, token, ecdhpairkey, err := clientsdk.HybridEncrypt(reqData, secKeyMgr.WeChatPubKey)
	if err != nil {
		return nil, err
	}
	ecdhPacket := &wechat.EcdhPacket{
		Type: proto.Uint32(1),
		Key: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: ecdhpairkey.PubKey,
		},
		Token:        token,
		Url:          proto.String(""),
		ProtobufData: encrypt,
	}

	secKeyMgr.PubKey = ecdhpairkey.PubKey
	secKeyMgr.PriKey = ecdhpairkey.PriKey
	secKeyMgr.SourceData = reqData
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, epKey[24:]...)
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, reqData...)

	ecdhDataPacket, err := proto.Marshal(ecdhPacket)
	if err != nil {
		return nil, err
	}

	packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, 252, ecdhDataPacket, ecdhDataPacket, uint32(len(ecdhDataPacket)), 12, uint32(0x4e))
	//设置Hybrid 加密密钥版本
	packHeader.HybridKeyVer = secKeyMgr.WeChatPubKeyVersion
	//开始组头
	retData := PackHeaderSerialize(packHeader, false)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/secmanualauth", retData)
	if err != nil {
		return nil, err
	}
	packHeader, err = DecodePackHeader(resp, nil)
	if err != nil {
		return nil, err
	}
	packHeader.Data, err = clientsdk.HybridEcdhDecrypt(packHeader.Data, secKeyMgr.PriKey, secKeyMgr.PubKey, secKeyMgr.FinalSha256)
	if err != nil {
		return nil, err
	}
	return packHeader, err*/
}

func GetManualAuthAesReqDataProtobufA16(userInfo *baseinfo.UserInfo) *wechat.ManualAuthAesReqData {
	AndroidId16 := userInfo.DeviceInfoA16
	zeroUint32 := uint32(0)
	zeroInt32 := int32(0)
	var aesRequest wechat.ManualAuthAesReqData
	// BaseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene uint32 = 1
	baseReq.Scene = &tmpScene
	baseReq.SessionKey = []byte{}
	baseReq.DeviceId = AndroidId16.DeviceId //A0cb9c716a607ed
	baseReq.ClientVersion = &baseinfo.AndroidClientVersion
	baseReq.OsType = &baseinfo.AndroidDeviceType
	aesRequest.BaseRequest = baseReq
	aesRequest.Imei = proto.String(AndroidId16.AndriodImei(AndroidId16.DeviceIdStr[:15])) //353978478595717
	//<softtype><lctmoc>0</lctmoc><level>0</level><k1>0 </k1><k2>B7366BF-2.5.27.9.56</k2><k3>8.1.0</k3><k4>357149293324373</k4><k5></k5><k6></k6><k7>06bef21034caa37f</k7><k8>01cf44dd6b7af383</k8><k9>Eeedf 37</k9><k10>8</k10><k11>Dbbf Technologies, Inc bcf7378</k11><k12></k12><k13></k13><k14>73:05:71:93:21:1e</k14><k15></k15><k16>half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt evtstrm aes pmull sha1 sha2 crc32</k16><k18>189df3a2df2ce70c0fdffb2b56ae884c</k18><k21>Chinanet-cbfda</k21><k22></k22><k24>9f:5e:d7:79:a5:88</k24><k26>0</k26><k30>&quot;Chinanet-cbfda&quot;</k30><k33>com.tencent.mm</k33><k34>google/bullhead/bullhead:8.1.0/DFA3.345967.885/3472334:user/release-keys</k34><k35>bullhead</k35><k36>CCD72f</k36><k37>google</k37><k38>bullhead</k38><k39>bullhead</k39><k40>bullhead</k40><k41>0</k41><k42>BDC</k42><k43>null</k43><k44>0</k44><k45></k45><k46></k46><k47>wifi</k47><k48>357149293324373</k48><k49>data/user/0/com.tencent.mm/</k49><k52>0</k52><k53>1</k53><k57>1640</k57><k58></k58><k59>3</k59><k60></k60><k61>true</k61><k62></k62><k63>A0cb9c716a607edb</k63><k64>30303566-6534-6636-6561-31643565393834313334383364313462643335653537</k64><k65></k65></softtype>
	aesRequest.SoftType = proto.String(AndroidId16.AndriodGetSoftType(AndroidId16.DeviceIdStr))
	aesRequest.BuiltinIpseq = &zeroUint32
	aesRequest.ClientSeqId = proto.String(fmt.Sprintf("%s_%d", AndroidId16.DeviceIdStr, (time.Now().UnixNano() / 1e6)))
	aesRequest.DeviceName = proto.String(AndroidId16.AndroidManufacturer(AndroidId16.DeviceIdStr) + "-" + AndroidId16.AndroidPhoneModel(AndroidId16.DeviceIdStr))
	//<AndroidDeviceInfo><MANUFACTURER name="BDC"><MODEL name="Eeedf 37"><VERSION_RELEASE name="8.1.0"><VERSION_INCREMENTAL name="3472334"><DISPLAY name="DFA3.345967.885"></DISPLAY></VERSION_INCREMENTAL></VERSION_RELEASE></MODEL></MANUFACTURER></AndroidDeviceInfo>
	aesRequest.DeviceType = proto.String(AndroidId16.AndriodDeviceType(AndroidId16.DeviceIdStr))
	aesRequest.Language = proto.String("Zh")
	aesRequest.TimeZone = proto.String("8.0")
	aesRequest.Channel = &zeroInt32
	//
	aesRequest.Ostype = &baseinfo.AndroidDeviceType
	aesRequest.Signature = proto.String(AndroidId16.AndriodPackageSign(AndroidId16.DeviceIdStr)) //189df3a2df2ce70c0fdffb2b56ae884c
	// TimeStamp
	tmpTime2 := uint32(time.Now().UnixNano() / 1000000000)
	aesRequest.TimeStamp = &tmpTime2
	aesRequest.DeviceBrand = proto.String("google")
	//Eeedf 37armeabi-10d
	aesRequest.DeviceModel = proto.String(AndroidId16.AndroidPhoneModel(AndroidId16.DeviceIdStr) + AndroidId16.AndroidArch(AndroidId16.DeviceIdStr))
	aesRequest.RealCountry = proto.String("cn")
	aesRequest.BundleId = proto.String("com.tencent.mm")
	inputType := uint32(2)
	aesRequest.InputType = &inputType
	ext, err := GetExtSpamInfoAndroid(userInfo)
	if err != nil {
		log.Error("Android extSpam err", err.Error())
	}
	// extSpamInfo
	var extSpamInfo wechat.SKBuiltinString_
	extSpamInfo.Buffer = ext
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	aesRequest.ExtSpamInfo = &extSpamInfo
	return &aesRequest
}

// GetManualAuthAesReqProtobuf 生成自动登陆aesreq项
func GetManualAuthAesReqDataA16Protobuf(userInfo *baseinfo.UserInfo) []byte {
	AndroidId16 := userInfo.DeviceInfoA16
	zeroUint32 := uint32(0)
	zeroInt32 := int32(0)
	var emptyString = string("")
	var aesRequest wechat.ManualAuthAesReqData
	// BaseRequest
	baseReq := GetBaseRequest(userInfo)
	// BaseAuthReqInfo
	var baseReqInfo wechat.BaseAuthReqInfo
	baseReqInfo.AuthReqFlag = &zeroUint32
	baseReqInfo.AuthTicket = &emptyString
	aesRequest.BaseReqInfo = &baseReqInfo

	var tmpScene uint32 = 1
	baseReq.Scene = &tmpScene
	baseReq.SessionKey = []byte{}
	baseReq.DeviceId = AndroidId16.DeviceId
	baseReq.ClientVersion = &baseinfo.AndroidClientVersion
	baseReq.OsType = &baseinfo.AndroidDeviceType
	aesRequest.BaseRequest = baseReq
	aesRequest.Imei = proto.String(AndroidId16.AndriodImei(AndroidId16.DeviceIdStr))
	aesRequest.SoftType = proto.String(AndroidId16.AndriodGetSoftType(AndroidId16.DeviceIdStr))
	aesRequest.BuiltinIpseq = &zeroUint32
	aesRequest.ClientSeqId = proto.String(fmt.Sprintf("%s_%d", AndroidId16.DeviceIdStr, (time.Now().UnixNano() / 1e6)))
	aesRequest.DeviceName = proto.String(AndroidId16.AndroidManufacturer(AndroidId16.DeviceIdStr) + "-" + AndroidId16.AndroidPhoneModel(AndroidId16.DeviceIdStr))
	aesRequest.DeviceType = proto.String(AndroidId16.AndriodDeviceType(AndroidId16.DeviceIdStr))
	aesRequest.Language = proto.String("Zh")
	if userInfo.LoginDataInfo.Language != "" {
		aesRequest.Language = proto.String(userInfo.LoginDataInfo.Language)
	}
	aesRequest.TimeZone = proto.String("8.0")
	aesRequest.Channel = &zeroInt32
	// TimeStamp
	tmpTime2 := uint32(time.Now().UnixNano() / 1000000000)
	aesRequest.TimeStamp = &tmpTime2
	aesRequest.DeviceBrand = proto.String("google")
	aesRequest.DeviceModel = proto.String(AndroidId16.AndroidPhoneModel(AndroidId16.DeviceIdStr)) // + AndroidId16.AndroidArch(AndroidId16.DeviceIdStr)
	aesRequest.RealCountry = proto.String("cn")
	aesRequest.BundleId = proto.String("com.tencent.mm")
	inputType := uint32(2)
	aesRequest.InputType = &inputType

	aesRequest.Ostype = &baseinfo.AndroidDeviceType
	aesRequest.Signature = proto.String(AndroidId16.AndriodPackageSign(AndroidId16.DeviceIdStr))
	//
	ext, err := GetExtSpamInfoAndroid(userInfo)
	if err != nil {
		log.Error("Android extSpam err", err.Error())
	}
	// extSpamInfo
	var extSpamInfo wechat.SKBuiltinString_
	extSpamInfo.Buffer = ext
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	aesRequest.ExtSpamInfo = &extSpamInfo
	reqData, _ := proto.Marshal(&aesRequest)
	return reqData
}

func GetExtSpamInfoAndroid(userInfo *baseinfo.UserInfo) ([]byte, error) {
	deviceToken := userInfo.DeviceInfoA16.DeviceToken
	T := time.Now().Unix()
	Wcstf, _ := extinfo.GetWcstf(userInfo.WxId)
	Wcste, _ := extinfo.GetWcste()
	AndroidCcData := extinfo.AndroidCcData(userInfo.DeviceInfoA16.DeviceIdStr, userInfo.DeviceInfoA16, deviceToken, T)
	CcData3PB, _ := proto.Marshal(AndroidCcData)
	curtime := uint32(T)
	DeviceTokenCCD := &mmproto.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mmproto.SKBuiltinStringt{
			String_: proto.String(deviceToken.GetTrustResponseData().GetDeviceToken()),
		},
		TimeStamp: &curtime,
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
	DeviceTokenCCDPB, _ := proto.Marshal(DeviceTokenCCD)

	WCExtInfo := &wechat.WCExtInfoNew{
		Wcstf: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(Wcstf))),
			Buffer: Wcstf,
		},
		Wcste: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(Wcste))),
			Buffer: Wcste,
		},
		CcData: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(CcData3PB))),
			Buffer: CcData3PB,
		},
		DeviceToken: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(DeviceTokenCCDPB))),
			Buffer: DeviceTokenCCDPB,
		},
	}
	WCExtInfoPB, _ := proto.Marshal(WCExtInfo)
	return WCExtInfoPB, nil
}
