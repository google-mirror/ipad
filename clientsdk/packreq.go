package clientsdk

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/xml"
	"feiyu.com/wx/clientsdk/android"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/clientsdk/cecdh"
	"feiyu.com/wx/clientsdk/extinfo"
	clientsdk "feiyu.com/wx/clientsdk/hybrid"
	"feiyu.com/wx/protobuf/wechat"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetLoginQRCodeReq 获取登陆二维码
func GetLoginQRCodeReq(userInfo *baseinfo.UserInfo) []byte {
	// 重新生成AesKey
	userInfo.SessionKey = baseutils.RandomBytes(16)

	// 构造请求
	var request wechat.LoginQRCodeRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// AESKey
	var aesKey wechat.AESKey
	var tmpAesKeyLen = uint32(16)
	aesKey.Len = &tmpAesKeyLen
	aesKey.Key = []byte(userInfo.SessionKey)
	request.Aes = &aesKey

	// OpCode
	var tmpOpcode = uint32(0)
	request.Opcode = &tmpOpcode
	// 发送请求
	src, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypeGetLoginQRCode, 7)
	return sendData
}

// SendManualAuth 发送ManualAuth请求
func GetManualAuthByAccountDataReq(userInfo *baseinfo.UserInfo, accountData []byte) []byte {
	// 打包加密数据
	subHeader := make([]byte, 0)
	tmpBytes := baseutils.Int32ToBytes(uint32(len(accountData)))
	subHeader = append(subHeader, tmpBytes[0:]...)

	// 获取 deviceDataRequest
	deviceData := GetManualAuthAesReqDataMarshal(userInfo)
	tmpBytes = baseutils.Int32ToBytes(uint32(len(deviceData)))
	subHeader = append(subHeader, tmpBytes[0:]...)

	// 加密压缩 accountData
	newAccountData := baseutils.CompressAndRsaByVer(accountData, userInfo.GetLoginRsaVer())
	tmpBytes = baseutils.Int32ToBytes(uint32(len(newAccountData)))
	subHeader = append(subHeader, tmpBytes[0:]...)
	subHeader = append(subHeader, newAccountData[0:]...)

	// 压缩加密 deviceData
	newDeviceData := baseutils.CompressAes(userInfo.SessionKey, deviceData)
	subHeader = append(subHeader, newDeviceData[0:]...)

	// 发送登陆请求
	sendData := Pack(userInfo, subHeader, baseinfo.MMRequestTypeManualAuth, 17)
	return sendData
}

// SendManualAuthA16发送A16登录请求
func GetManualAuthA16Req(userInfo *baseinfo.UserInfo, accountData []byte) []byte {
	/*req := &wechat.ManualAuthRequest{
		RsaReqData: GetManualAuthRsaReqDataProtobuf(userInfo,userInfo.LoginDataInfo.UserName,userInfo.LoginDataInfo.PassWord),
		AesReqData: GetManualAuthAesReqDataA16Protobuf(userInfo),
	}*/
	// 发送请求
	// 打包加密数据
	subHeader := make([]byte, 0)
	tmpBytes := baseutils.Int32ToBytes(uint32(len(accountData)))
	subHeader = append(subHeader, tmpBytes[0:]...)

	// 获取 deviceDataRequest
	deviceData := GetManualAuthAesReqDataA16Protobuf(userInfo)
	tmpBytes = baseutils.Int32ToBytes(uint32(len(deviceData)))
	subHeader = append(subHeader, tmpBytes[0:]...)

	// 加密压缩 accountData
	newAccountData := baseutils.CompressAndRsaByVer(accountData, userInfo.GetLoginRsaVer())
	tmpBytes = baseutils.Int32ToBytes(uint32(len(newAccountData)))
	subHeader = append(subHeader, tmpBytes[0:]...)
	subHeader = append(subHeader, newAccountData[0:]...)

	// 压缩加密 deviceData
	newDeviceData := baseutils.CompressAes(userInfo.SessionKey, deviceData)
	subHeader = append(subHeader, newDeviceData[0:]...)

	// 发送登陆请求
	sendData := Pack(userInfo, subHeader, baseinfo.MMRequestTypeManualAuth, 17)
	return sendData
}

// SendManualAuth 发送登陆请求
func GetManualAuthAccountDataReq(userInfo *baseinfo.UserInfo, newpass string, wxid string) ([]byte, error) {
	var tmpNid uint32 = 713
	userInfo.EcPublicKey, userInfo.EcPrivateKey = cecdh.GenerateEccKey()
	authRequest := &wechat.ManualAuthAccountRequest{}
	// aes_key
	var aesKey wechat.AESKey
	var tmpAesKeyLen uint32 = 16
	aesKey.Len = &tmpAesKeyLen
	aesKey.Key = []byte(userInfo.SessionKey)
	authRequest.Aes = &aesKey

	// 其它参数
	var ecdhKey wechat.ECDHKey
	var key wechat.SKBuiltinString_
	key.Buffer = userInfo.EcPublicKey
	var tmpLen = (uint32)(len(userInfo.EcPublicKey))
	key.Len = &tmpLen
	ecdhKey.Nid = &tmpNid
	ecdhKey.Key = &key
	authRequest.EcdhKey = &ecdhKey
	authRequest.UserName = &wxid
	userInfo.WxId = wxid
	//判断是否为iPad登录的伪密码
	if !strings.HasPrefix(newpass, "extdevnewpwd_") && !strings.HasPrefix(newpass, "strdm@") {
		newpass = baseutils.Md5Value(newpass)
	}
	authRequest.Password_1 = &newpass

	// 序列化
	accountData, err := proto.Marshal(authRequest)
	if err != nil {
		return nil, err
	}
	return accountData, err
}

// 获取DeviceToken IOS
func GetIosDeviceTokenReq(userInfo *baseinfo.UserInfo) ([]byte, *android.Client) {
	deviceIos := userInfo.DeviceInfo
	deviceId := Get62Key(userInfo.LoginDataInfo.LoginData)
	if deviceId[:2] != "49" {
		deviceId = "49" + deviceId[2:]
	}
	language := "id" //zh
	locale_country := "CN"
	if userInfo.LoginDataInfo.Language != "" {
		language = userInfo.LoginDataInfo.Language
		locale_country = language
	}
	uuid1, uuid2 := baseinfo.IOSUuid(deviceId)
	td := &wechat.TrustReq{
		Td: &wechat.TrustData{
			Tdi: []*wechat.TrustDeviceInfo{
				{Key: proto.String("deviceid"), Val: proto.String(deviceId)},
				{Key: proto.String("sdi"), Val: proto.String(extinfo.GetCidMd5(deviceId, extinfo.GetCid(0x0262626262626)))},
				{Key: proto.String("idfv"), Val: proto.String(uuid1)},
				{Key: proto.String("idfa"), Val: proto.String(uuid2)},
				{Key: proto.String("device_model"), Val: proto.String(deviceIos.IphoneVer)},
				{Key: proto.String("os_version"), Val: proto.String(deviceIos.OsType)},
				{Key: proto.String("core_count"), Val: proto.String("4")},
				{Key: proto.String("carrier_name"), Val: proto.String("")},
				{Key: proto.String("is_jailbreak"), Val: proto.String("0")},
				{Key: proto.String("device_name"), Val: proto.String(deviceIos.DeviceName)},
				{Key: proto.String("client_version"), Val: proto.String(fmt.Sprintf("%v", baseinfo.ClientVersion))},
				{Key: proto.String("plist_version"), Val: proto.String(fmt.Sprintf("%v", baseinfo.ClientVersion))},
				{Key: proto.String("language"), Val: proto.String(language)},
				{Key: proto.String("locale_country"), Val: proto.String(locale_country)},
				{Key: proto.String("screen_width"), Val: proto.String("834")},
				{Key: proto.String("screen_height"), Val: proto.String("1112")},
				{Key: proto.String("install_time"), Val: proto.String("1586355322")},
				{Key: proto.String("kern_boottime"), Val: proto.String("1586355519000")},
			},
		},
	}
	pb, _ := proto.Marshal(td)
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(pb)
	w.Close()

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(b.Bytes())
	randKey := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, randKey)
	fp := &wechat.FPFresh{
		BaseRequest: &wechat.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(0),
			DeviceId:      append([]byte(deviceId), 0),
			ClientVersion: proto.Uint32(baseinfo.ClientVersion),
			OsType:        &userInfo.DeviceInfo.OsType,
			Scene:         proto.Uint32(0),
		},
		SessKey: randKey,
		Ztdata: &wechat.ZTData{
			Version:   []byte("00000006"),
			Encrypted: proto.Uint64(1),
			Data:      encData,
			TimeStamp: proto.Int64(int64(time.Now().Unix())),
			OpType:    proto.Uint64(5),
			Uin:       proto.Uint64(0),
		},
	}
	reqData, _ := proto.Marshal(fp)

	hec := &android.Client{}
	hec.Init("IOS", int(baseinfo.ClientVersion), baseinfo.DeviceTypeIos)
	hecData := hec.HybridEcdhPackIosEn(3789, 0, nil, reqData)
	return hecData, hec
}

// 获取DeviceToken
func GetAndroIdDeviceTokenReq(userInfo *baseinfo.UserInfo) ([]byte, *android.Client) {
	Android16 := userInfo.DeviceInfoA16
	td := &wechat.TrustReq{
		Td: &wechat.TrustData{
			Tdi: []*wechat.TrustDeviceInfo{
				{Key: proto.String("IMEI"), Val: proto.String(Android16.AndriodImei(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("AndroidID"), Val: proto.String(Android16.AndriodID(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("PhoneSerial"), Val: proto.String(Android16.AndriodPhoneSerial(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("cid"), Val: proto.String("")},
				{Key: proto.String("WidevineDeviceID"), Val: proto.String(Android16.AndriodWidevineDeviceID(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("WidevineProvisionID"), Val: proto.String(Android16.AndriodWidevineProvisionID(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("GSFID"), Val: proto.String("")},
				{Key: proto.String("SoterID"), Val: proto.String("")},
				{Key: proto.String("SoterUid"), Val: proto.String("")},
				{Key: proto.String("FSID"), Val: proto.String(Android16.AndriodFSID(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("BootID"), Val: proto.String("")},
				{Key: proto.String("IMSI"), Val: proto.String("")},
				{Key: proto.String("PhoneNum"), Val: proto.String("")},
				{Key: proto.String("WeChatInstallTime"), Val: proto.String("1556077144")},
				{Key: proto.String("PhoneModel"), Val: proto.String(Android16.AndroidPhoneModel(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("BuildBoard"), Val: proto.String("bullhead")},
				{Key: proto.String("BuildBootloader"), Val: proto.String(Android16.AndroidBuildBoard(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("SystemBuildDate"), Val: proto.String("Fri Sep 28 23:37:27 UTC 2019")},
				{Key: proto.String("SystemBuildDateUTC"), Val: proto.String("1538177847")},
				{Key: proto.String("BuildFP"), Val: proto.String(Android16.AndroidBuildFP(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("BuildID"), Val: proto.String(Android16.AndroidBuildID(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("BuildBrand"), Val: proto.String("google")},
				{Key: proto.String("BuildDevice"), Val: proto.String("bullhead")},
				{Key: proto.String("BuildProduct"), Val: proto.String("bullhead")},
				{Key: proto.String("Manufacturer"), Val: proto.String(Android16.AndroidManufacturer(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("RadioVersion"), Val: proto.String(Android16.AndroidRadioVersion(Android16.DeviceIdStr[:15]))},
				{Key: proto.String("AndroidVersion"), Val: proto.String(Android16.AndroidVersion())},
				{Key: proto.String("SdkIntVersion"), Val: proto.String("27")},
				{Key: proto.String("ScreenWidth"), Val: proto.String("1080")},
				{Key: proto.String("ScreenHeight"), Val: proto.String("1794")},
				{Key: proto.String("SensorList"), Val: proto.String("BMI160 accelerometer#Bosch#0.004788#1,BMI160 gyroscope#Bosch#0.000533#1,BMM150 magnetometer#Bosch#0.000000#1,BMP280 pressure#Bosch#0.005000#1,BMP280 temperature#Bosch#0.010000#1,RPR0521 Proximity Sensor#Rohm#1.000000#1,RPR0521 Light Sensor#Rohm#10.000000#1,Orientation#Google#1.000000#1,BMI160 Step detector#Bosch#1.000000#1,Significant motion#Google#1.000000#1,Gravity#Google#1.000000#1,Linear Acceleration#Google#1.000000#1,Rotation Vector#Google#1.000000#1,Geomagnetic Rotation Vector#Google#1.000000#1,Game Rotation Vector#Google#1.000000#1,Pickup Gesture#Google#1.000000#1,Tilt Detector#Google#1.000000#1,BMI160 Step counter#Bosch#1.000000#1,BMM150 magnetometer (uncalibrated)#Bosch#0.000000#1,BMI160 gyroscope (uncalibrated)#Bosch#0.000533#1,Sensors Sync#Google#1.000000#1,Double Twist#Google#1.000000#1,Double Tap#Google#1.000000#1,Device Orientation#Google#1.000000#1,BMI160 accelerometer (uncalibrated)#Bosch#0.004788#1")},
				{Key: proto.String("DefaultInputMethod"), Val: proto.String("com.google.android.inputmethod.latin")},
				{Key: proto.String("InputMethodList"), Val: proto.String("Google \345\215\260\345\272\246\350\257\255\351\224\256\347\233\230#com.google.android.apps.inputmethod.hindi,Google \350\257\255\351\237\263\350\276\223\345\205\245#com.google.android.googlequicksearchbox,Google \346\227\245\350\257\255\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.japanese,Google \351\237\251\350\257\255\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.korean,Gboard#com.google.android.inputmethod.latin,\350\260\267\346\255\214\346\213\274\351\237\263\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.pinyin")},
				{Key: proto.String("DeviceID"), Val: proto.String(Android16.DeviceIdStr[:15])},
				{Key: proto.String("OAID"), Val: proto.String("")},
			},
		},
	}
	pb, _ := proto.Marshal(td)

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(pb)
	randKey := make([]byte, 16)
	io.ReadFull(rand.Reader, randKey)
	fp := &wechat.FPFresh{
		BaseRequest: GetBaseRequest(userInfo),
		SessKey:     randKey,
		Ztdata: &wechat.ZTData{
			Version:   []byte("00000006"),
			Encrypted: proto.Uint64(1),
			Data:      encData,
			TimeStamp: proto.Int64(int64(time.Now().Unix())),
			OpType:    proto.Uint64(5),
			Uin:       proto.Uint64(0),
		},
	}
	reqData, _ := proto.Marshal(fp)
	hec := &android.Client{}
	hec.Init("Android", int(baseinfo.AndroidClientVersion), baseinfo.AndroidDeviceType)
	sendData := hec.HybridEcdhPackAndroidEn(3789, 10002, 0, nil, reqData)
	return sendData, hec
}

// 二次登录-new
func GetSecautouthReq(userInfo *baseinfo.UserInfo) ([]byte, *SecLoginKeyMgr, error) {
	userInfo.EcPublicKey, userInfo.EcPrivateKey = cecdh.GenerateEccKey()
	autoAuthKey := &wechat.AutoAuthKey{}
	err := proto.Unmarshal(userInfo.AutoAuthKey, autoAuthKey)
	if err != nil {
		return nil, nil, err
	}
	userInfo.SessionKey = autoAuthKey.EncryptKey.Buffer
	var tmpNid uint32 = 713
	var key wechat.SKBuiltinString_
	key.Buffer = userInfo.EcPublicKey
	var tmpLen = (uint32)(len(userInfo.EcPublicKey))
	key.Len = &tmpLen
	// ClientSeqId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	var strClientSeqID = string(userInfo.DeviceInfo.Imei + "-" + tmpTimeStr)
	// extSpamInfo
	var extSpamInfo wechat.SKBuiltinString_
	extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	req := &wechat.AutoAuthRequest{
		RsaReqData: &wechat.AutoAuthRsaReqData{
			AesEncryptKey: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(autoAuthKey.EncryptKey.Buffer))),
				Buffer: autoAuthKey.EncryptKey.Buffer,
			},
			PubEcdhKey: &wechat.ECDHKey{
				Nid: proto.Uint32(tmpNid),
				Key: &key,
			},
		},
		AesReqData: &wechat.AutoAuthAesReqData{
			BaseRequest: GetBaseRequest(userInfo),
			BaseReqInfo: &wechat.BaseAuthReqInfo{},
			AutoAuthKey: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(userInfo.AutoAuthKey))),
				Buffer: userInfo.AutoAuthKey,
			},
			Imei:         &userInfo.DeviceInfo.Imei,
			SoftType:     &userInfo.DeviceInfo.SoftTypeXML,
			BuiltinIpSeq: proto.Uint32(0),
			ClientSeqId:  &strClientSeqID,
			DeviceName:   proto.String(userInfo.DeviceInfo.DeviceName),
			SoftInfoXml:  proto.String("iPhone"),
			Language:     proto.String("zh_CN"),
			TimeZone:     proto.String("8.0"),
			ExtSpamInfo:  &extSpamInfo,
		},
	}
	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	secKeyMgr := NewSecLoginKeyMgrByVer(146)
	//加密
	encrypt, epKey, token, ecdhpairkey, err := clientsdk.HybridEncrypt(reqData, secKeyMgr.WeChatPubKey)
	if err != nil {
		return nil, nil, err
	}
	/*ecdhPacket := &wechat.EcdhPacket{
		Type: proto.Uint32(1),
		Key: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: ecdhpairkey.PubKey,
		},
		Token:        token,
		Url:          proto.String(""),
		ProtobufData: encrypt,
	}*/
	ecdhPacket := &wechat.HybridEcdhRequest{
		Type: proto.Int32(1),
		SecECDHKey: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: ecdhpairkey.PubKey,
		},
		Randomkeydata:       token,
		Randomkeyextenddata: epKey,
		Encyptdata:          encrypt,
	}
	secKeyMgr.PubKey = ecdhpairkey.PubKey
	secKeyMgr.PriKey = ecdhpairkey.PriKey
	secKeyMgr.SourceData = reqData
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, epKey[24:]...)
	secKeyMgr.FinalSha256 = append(secKeyMgr.FinalSha256, reqData...)
	ecdhDataPacket, err := proto.Marshal(ecdhPacket)
	if err != nil {
		return nil, nil, err
	}

	packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, 763, ecdhDataPacket, ecdhDataPacket, uint32(len(ecdhDataPacket)), 12, uint32(0x4e))
	//设置Hybrid 加密密钥版本
	packHeader.HybridKeyVer = secKeyMgr.WeChatPubKeyVersion
	//开始组头
	retData := PackHeaderSerialize(packHeader, false)
	return retData, secKeyMgr, nil
}

// Secautoauth二次登录
func GetSecautoauthReq(userInfo *baseinfo.UserInfo) ([]byte, *android.Client, error) {
	Autoauthkey := &wechat.AutoAuthKey{}
	_ = proto.Unmarshal(userInfo.AutoAuthKey, Autoauthkey)
	userInfo.SessionKey = Autoauthkey.EncryptKey.Buffer
	userInfo.EcPublicKey, userInfo.EcPrivateKey = cecdh.GenerateEccKey()
	//基础设备信息
	Imei := userInfo.DeviceInfo.Imei
	SoftType := userInfo.DeviceInfo.SoftTypeXML
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	ClientSeqId := string(userInfo.DeviceInfo.Imei + "-" + tmpTimeStr)
	WCExtInfoseq := GetExtPBSpamInfoData(userInfo)
	req := &wechat.AutoAuthRequest{
		RsaReqData: &wechat.AutoAuthRsaReqData{
			AesEncryptKey: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(Autoauthkey.EncryptKey.Buffer))),
				Buffer: Autoauthkey.EncryptKey.Buffer,
			},
			PubEcdhKey: &wechat.ECDHKey{
				Nid: proto.Uint32(713),
				Key: &wechat.SKBuiltinString_{
					Len:    proto.Uint32(uint32(len(userInfo.EcPublicKey))),
					Buffer: userInfo.EcPublicKey,
				},
			},
		},
		AesReqData: &wechat.AutoAuthAesReqData{
			BaseRequest: GetBaseRequest(userInfo),
			BaseReqInfo: &wechat.BaseAuthReqInfo{},
			AutoAuthKey: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(userInfo.AutoAuthKey))),
				Buffer: userInfo.AutoAuthKey,
			},
			Imei:         &Imei,
			SoftType:     &SoftType,
			BuiltinIpSeq: proto.Uint32(0),
			ClientSeqId:  &ClientSeqId,
			DeviceName:   proto.String(userInfo.DeviceInfo.DeviceName),
			SoftInfoXml:  proto.String("iPhone"),
			Language:     proto.String("zh_CN"),
			TimeZone:     proto.String("8.0"),
			ExtSpamInfo: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(WCExtInfoseq))),
				Buffer: WCExtInfoseq,
			},
		},
	}
	reqdata, err := proto.Marshal(req)
	if err != nil {
		return nil, nil, err
	}
	hec := &android.Client{}
	hec.Init("IOS", int(baseinfo.ClientVersion), userInfo.DeviceInfo.OsType)
	hecData := hec.HybridEcdhPackIosEn(763, userInfo.Uin, userInfo.Session, reqdata)
	return hecData, hec, nil
}

// GetAutoAuthReq 发送token登陆请求
func GetAutoAuthReq(userInfo *baseinfo.UserInfo) ([]byte, error) {
	userInfo.EcPublicKey, userInfo.EcPrivateKey = cecdh.GenerateEccKey()
	autoAuthKey := &wechat.AutoAuthKey{}
	err := proto.Unmarshal(userInfo.AutoAuthKey, autoAuthKey)
	if err != nil {
		return nil, err
	}
	userInfo.SessionKey = autoAuthKey.EncryptKey.Buffer
	// 获取AutoAuthRsaReqData
	rsaReqData := GetAutoAuthRsaReqDataMarshal(userInfo)
	aesReqData := GetAutoAuthAesReqDataMarshal(userInfo)

	// 开始打包数据
	// 加密压缩 rsaReqData
	rsaEncodeData := baseutils.CompressAndRsaByVer(rsaReqData, userInfo.GetLoginRsaVer())
	rsaAesEncodeData := baseutils.CompressAes(userInfo.SessionKey, rsaReqData)

	// 加密压缩 aesReqData
	aesEncodeData := baseutils.CompressAes(userInfo.SessionKey, aesReqData)

	body := make([]byte, 0)
	// rsaReqBuflen
	tmpBuf := baseutils.Int32ToBytes(uint32(len(rsaReqData)))
	body = append(body, tmpBuf[0:]...)

	// aesReqBuf len
	tmpBuf = baseutils.Int32ToBytes(uint32(len(aesReqData)))
	body = append(body, tmpBuf[0:]...)

	// rsaEncodeData len
	tmpBuf = baseutils.Int32ToBytes(uint32(len(rsaEncodeData)))
	body = append(body, tmpBuf[0:]...)

	// rsaAesEncodeData len
	tmpBuf = baseutils.Int32ToBytes(uint32(len(rsaAesEncodeData)))
	body = append(body, tmpBuf[0:]...)

	// rsaEncodeData
	body = append(body, rsaEncodeData[0:]...)
	body = append(body, rsaAesEncodeData[0:]...)
	body = append(body, aesEncodeData[0:]...)
	// 发送请求
	sendData := Pack(userInfo, body, baseinfo.MMRequestTypeAutoAuth, 9)
	return sendData, nil
}

// GetCheckLoginQRCodeReq 长链接：获取检测二维码状态-请求数据包
func GetCheckLoginQRCodeReq(userInfo *baseinfo.UserInfo, qrcodeUUID string) ([]byte, error) {
	var request wechat.CheckLoginQRCodeRequest
	// 重新生成AesKey
	//userInfo.SessionKey = baseutils.RandomBytes(16)

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// AESKey
	var aesKey wechat.AESKey
	var tmpAesKeyLen = uint32(16)
	aesKey.Len = &tmpAesKeyLen
	aesKey.Key = userInfo.SessionKey
	request.Aes = &aesKey

	// uuid
	request.Uuid = &qrcodeUUID

	// timeStamp 当前系统时间
	timeStamp := (uint32)(time.Now().Unix())
	request.TimeStamp = &timeStamp

	// OpCode
	var tmpOpcode = uint32(0)
	request.Opcode = &tmpOpcode

	// 发送请求
	src, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypeCheckLoginQRCode, 7)
	return sendData, nil
}

// SendPushQrLoginNotice 二维码二次登录
func GetPushQrLoginNoticeReq(userInfo *baseinfo.UserInfo) []byte {
	var request wechat.PushLoginURLRequest

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene

	request = wechat.PushLoginURLRequest{
		BaseRequest:    baseReq,
		UserName:       proto.String(userInfo.GetUserName()),
		AutoAuthTicket: proto.String(""),
		ClientID:       proto.String(fmt.Sprintf("iPad-Push-%s.110141", userInfo.DeviceInfo.DeviceID)),
		RandomEnCryptKey: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(userInfo.SessionKey))),
			Buffer: userInfo.SessionKey,
		},
		OPCode:     proto.Uint32(3),
		DeviceName: proto.String(userInfo.DeviceInfo.DeviceName),
		AutoAuthKey: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(userInfo.AutoAuthKey))),
			Buffer: userInfo.AutoAuthKey,
		},
	}
	// 发送请求
	src, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypePushQrLogin, 1)
	return sendData
}

// GetHeartBeatReq 长链接：获取心跳包-请求数据包
func GetHeartBeatReq(userInfo *baseinfo.UserInfo) ([]byte, error) {
	var request wechat.HeartBeatRequest

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// timeStamp 当前系统时间
	timeStamp := (uint32)(time.Now().UnixNano() / 1000000000)
	request.TimeStamp = &timeStamp

	// Scene
	request.Scene = &tmpScene

	// KeyBuf
	request.KeyBuf = &wechat.SKBuiltinString_{}
	request.KeyBuf.Buffer = userInfo.SessionKey
	request.KeyBuf.Len = proto.Uint32(uint32(len(userInfo.SessionKey)))

	// 打包数据
	src, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypeHeartBeat, 5)
	return sendData, nil
}

// 同步消息
func GetNewSyncHistoryMessageReq(userInfo *baseinfo.UserInfo, scene uint32, syncKey string) []byte {
	var Synckey wechat.SKBuiltinString_
	if userInfo.SyncHistoryKey == nil {
		userInfo.SyncHistoryKey = userInfo.SyncKey
	}
	Synckey = wechat.SKBuiltinString_{
		Len:    proto.Uint32(uint32(len(userInfo.SyncHistoryKey))),
		Buffer: userInfo.SyncHistoryKey,
	}
	request := &wechat.NewSyncRequest{
		Oplog: &wechat.CmdList{
			Count:    proto.Uint32(0),
			ItemList: nil,
		},
		Selector:      proto.Uint32(262151),
		KeyBuf:        &Synckey,
		Scene:         proto.Uint32(scene),
		SyncMsgDigest: proto.Uint32(3),
	}
	// DeviceType
	if userInfo.DeviceInfoA16 != nil {
		request.DeviceType = &baseinfo.AndroidDeviceType
	} else {
		request.DeviceType = &userInfo.DeviceInfo.OsType
	}
	// 发送请求
	src, _ := proto.Marshal(request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypeNewSync, 5)
	return sendData
}

// GetNewSyncReq 发送同步信息请求
func GetNewSyncReq(userInfo *baseinfo.UserInfo, scene uint32, short bool) []byte {
	var request wechat.NewSyncRequest
	zeroValue32 := uint32(0)

	// Oplog
	var opLog wechat.CmdList
	opLog.Count = &zeroValue32
	opLog.ItemList = make([]*wechat.CmdItem, 0)
	request.Oplog = &opLog

	// Selector
	tmpSelector := uint32(262151)
	request.Selector = &tmpSelector

	// keyBuf
	var keyBuf wechat.SKBuiltinString_
	keyBuf.Buffer = userInfo.SyncKey
	tmpLen := uint32(len(keyBuf.Buffer))
	keyBuf.Len = &tmpLen
	request.KeyBuf = &keyBuf

	// Scene
	request.Scene = &scene
	// DeviceType
	if userInfo.DeviceInfoA16 != nil {
		request.DeviceType = &baseinfo.AndroidDeviceType
	} else {
		request.DeviceType = &userInfo.DeviceInfo.OsType
	}
	// syncMsgDigest : 短链接同步
	syncMsgDigest := baseinfo.MMSyncMsgDigestTypeShortLink
	if !short {
		syncMsgDigest = baseinfo.MMSyncMsgDigestTypeLongLink
	}
	request.SyncMsgDigest = &syncMsgDigest

	// 发送请求
	src, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypeNewSync, 5)
	return sendData
}

// 同步消息
func GetWxSyncMsgReq(userInfo *baseinfo.UserInfo, key string) []byte {
	zeroValue32 := uint32(0)
	// Oplog
	var opLog wechat.CmdList
	opLog.Count = &zeroValue32
	opLog.ItemList = make([]*wechat.CmdItem, 0)
	//
	var keyBuf wechat.SKBuiltinString_
	keyBuf.Buffer = userInfo.SyncKey
	if key != "" {
		keyBuf.Buffer = []byte(key)
	}
	tmpLen := uint32(len(keyBuf.Buffer))
	keyBuf.Len = &tmpLen
	request := wechat.NewSyncRequest{
		Oplog:         &opLog,
		Selector:      proto.Uint32(262151),
		Scene:         proto.Uint32(4),
		DeviceType:    proto.String(userInfo.DeviceInfo.OsType),
		SyncMsgDigest: proto.Uint32(baseinfo.MMSyncMsgDigestTypeShortLink),
		KeyBuf:        &keyBuf,
	}
	// 发送请求
	src, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, src, baseinfo.MMRequestTypeNewSync, 5)
	return sendData
}

// GetProfileReq 发送获取帐号所有信息请求
func GetProfileReq(userInfo *baseinfo.UserInfo) []byte {
	var request wechat.GetProfileRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpZero = uint32(0)
	baseReq.Scene = &tmpZero
	request.BaseRequest = baseReq

	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetProfile, 5)
	return sendData
}

// 获取设备
func GetSafetyInfoReq(userInfo *baseinfo.UserInfo) []byte {
	baseRequest := GetBaseRequest(userInfo)
	//baseRequest.Scene = proto.Uint32(1)
	var req = wechat.GetSafetyInfoRequest{
		BaseRequest: baseRequest,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 850, 5)
	return sendEncodeData
}

// 删除设备
func GetDelSafeDeviceReq(userInfo *baseinfo.UserInfo, uuid string) []byte {
	baseRequest := GetBaseRequest(userInfo)
	//baseRequest.Scene = proto.Uint32(1)
	var req = wechat.DelSafeDeviceRequest{
		BaseRequest: baseRequest,
		Uuid:        proto.String(uuid),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 362, 5)
	return sendEncodeData
}

// 检测微信登录环境
func GetCheckCanSetAliasReq(userInfo *baseinfo.UserInfo) []byte {
	baseRequest := GetBaseRequest(userInfo)
	var req = wechat.CheckCanSetAliasReq{
		BaseRequest: baseRequest,
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, 926, 5)
	return sendEncodeData
}

// 扫码登录新设备
func GetExtDeviceLoginConfirmGetReq(userInfo *baseinfo.UserInfo, url string) []byte {
	Url := strings.Replace(url, "https", "http", -1)
	req := &wechat.ExtDeviceLoginConfirmGetRequest{
		LoginUrl:   proto.String(Url),
		DeviceName: proto.String(baseinfo.DeviceTypeIos),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, 971, 5)
	return sendEncodeData
}

// 新设备扫码确认登录
func GetExtDeviceLoginConfirmOkReq(userInfo *baseinfo.UserInfo, url string) []byte {
	Url := strings.Replace(url, "https", "http", -1)
	req := &wechat.ExtDeviceLoginConfirmOKRequest{
		LoginUrl:    proto.String(Url),
		SessionList: proto.String(""),
		SyncMsg:     proto.Uint64(1),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, 972, 5)
	return sendEncodeData
}

// SendInitContactReq 初始化联系人列表
func GetInitContactReq(userInfo *baseinfo.UserInfo, contactSeq uint32) []byte {
	var request wechat.InitContactReq

	// Username
	request.Username = &userInfo.WxId
	// CurrentWxcontactSeq
	request.CurrentWxcontactSeq = &contactSeq
	// CurrentChatRoomContactSeq
	roomContactSeq := uint32(0)
	request.CurrentChatRoomContactSeq = &roomContactSeq

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeInitContact, 5)
	return sendEncodeData
}

func GetContactListPageReq(userInfo *baseinfo.UserInfo, CurrentWxcontactSeq uint32, CurrentChatRoomContactSeq uint32) []byte {
	var request wechat.InitContactReq

	// Username
	request.Username = &userInfo.WxId
	// CurrentWxcontactSeq
	request.CurrentWxcontactSeq = &CurrentWxcontactSeq
	request.CurrentChatRoomContactSeq = &CurrentChatRoomContactSeq

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeInitContact, 5)
	return sendEncodeData
}

// SendBatchGetContactBriefInfoReq 批量获取联系人信息
func GetBatchGetContactBriefInfoReq(userInfo *baseinfo.UserInfo, userNameList []string) []byte {
	var request wechat.BatchGetContactBriefInfoReq
	request.ContactUsernameList = userNameList
	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeBatchGetContactBriefInfo, 5)
	return sendEncodeData
}

func GetFriendRelationReq(userInfo *baseinfo.UserInfo, userName string) []byte {
	var request wechat.MMBizJsApiGetUserOpenIdRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(1)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	request.AppId = proto.String("wx7c8d593b2c3a7703")
	request.UserName = proto.String(userName)
	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	//获取好友关系状态
	sendEncodeData := Pack(userInfo, srcData, 1177, 5)
	return sendEncodeData
}

// SendGetContactRequest 获取指定微信号信息请求, userWxID:联系人ID  roomWxID：群ID
func GetContactReq(userInfo *baseinfo.UserInfo, userWxIDList []string, antisPanTicketList []string, roomWxIDList []string) []byte {
	var request wechat.GetContactRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// userCount
	var userCount = uint32(len(userWxIDList))
	request.UserCount = &userCount
	// UserNameList
	userNameList := make([]*wechat.SKBuiltinString, userCount)
	// 遍历
	for index := uint32(0); index < userCount; index++ {
		userNameItem := new(wechat.SKBuiltinString)
		userNameItem.Str = &userWxIDList[index]
		userNameList[index] = userNameItem
	}
	request.UserNameList = userNameList

	// AntispamTicketCount
	antispamTicketCount := uint32(len(antisPanTicketList))
	request.AntispamTicketCount = &antispamTicketCount
	// AntispamTicket
	tmpAntispamTicketList := make([]*wechat.SKBuiltinString, antispamTicketCount)
	for index := uint32(0); index < antispamTicketCount; index++ {
		antispamTicket := new(wechat.SKBuiltinString)
		antispamTicket.Str = &antisPanTicketList[index]
		tmpAntispamTicketList[index] = antispamTicket
	}
	request.AntispamTicket = tmpAntispamTicketList

	// FromChatRoomCount
	fromChatRoomCount := uint32(len(roomWxIDList))
	request.FromChatRoomCount = &fromChatRoomCount
	// FromChatRoom
	fromChatRoomList := make([]*wechat.SKBuiltinString, fromChatRoomCount)
	for index := uint32(0); index < fromChatRoomCount; index++ {
		fromChatRoom := new(wechat.SKBuiltinString)
		fromChatRoom.Str = &roomWxIDList[index]
		fromChatRoomList[index] = fromChatRoom
	}
	request.FromChatRoom = fromChatRoomList

	// GetContactScene
	var getContactScene = uint32(0)
	request.GetContactScene = &getContactScene

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetContact, 5)
	return sendEncodeData
}

// 创建红包
func GetWXCreateRedPacketReq(userInfo *baseinfo.UserInfo, hbItem *baseinfo.RedPacket) []byte {
	var request wechat.HongBaoReq
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// CgiCmd
	request.CgiCmd = proto.Uint32(0)
	// OutPutType
	request.OutPutType = proto.Uint32(0)
	// ReqText
	strReqText := string("")
	strReqText = strReqText + "city=Guangzhou&"
	strReqText = strReqText + "hbType=" + strconv.Itoa(int(hbItem.RedType)) + "&"
	strReqText = strReqText + "headImg=" + "&"
	strReqText = strReqText + "inWay=" + strconv.Itoa(int(hbItem.From)) + "&"
	strReqText = strReqText + "needSendToMySelf=0" + "&"
	strReqText = strReqText + "nickName=" + url.QueryEscape(userInfo.NickName) + "&"
	strReqText = strReqText + "perValue=" + strconv.Itoa(int(hbItem.Amount)) + "&"
	strReqText = strReqText + "province=Guangdong" + "&"
	strReqText = strReqText + "sendUserName=" + userInfo.WxId + "&"
	strReqText = strReqText + "totalAmount=" + strconv.Itoa(int(hbItem.Amount*hbItem.Count)) + "&"
	strReqText = strReqText + "totalNum=" + strconv.Itoa(int(hbItem.Count)) + "&"
	strReqText = strReqText + "username=" + hbItem.Username + "&"
	strReqText = strReqText + "wishing=" + url.QueryEscape(hbItem.Content)
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText
	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, 1575, 5)
	return sendEncodeData
}

// SendReceiveWxHB 发送接收红包请求
func GetReceiveWxHBReq(userInfo *baseinfo.UserInfo, hongBaoReceiverItem *baseinfo.HongBaoReceiverItem) []byte {
	var request wechat.HongBaoReq

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &hongBaoReceiverItem.CgiCmd

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	strReqText := string("")
	strReqText = strReqText + "agreeDuty=0&"
	strReqText = strReqText + "channelId=" + hongBaoReceiverItem.HongBaoURLItem.ChannelID + "&"
	strReqText = strReqText + "city=" + hongBaoReceiverItem.City + "&"
	strReqText = strReqText + "encrypt_key=" + baseutils.EscapeURL(userInfo.HBAesKeyEncrypted) + "&"
	strReqText = strReqText + "encrypt_userinfo=" + baseutils.EscapeURL(GetEncryptUserInfo(userInfo)) + "&"
	strReqText = strReqText + "inWay=" + strconv.Itoa(int(hongBaoReceiverItem.InWay)) + "&"
	strReqText = strReqText + "msgType=" + hongBaoReceiverItem.HongBaoURLItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoReceiverItem.NativeURL) + "&"
	strReqText = strReqText + "province=" + hongBaoReceiverItem.Province + "&"
	strReqText = strReqText + "sendId=" + hongBaoReceiverItem.HongBaoURLItem.SendID
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeReceiveWxHB, 5)
	return sendEncodeData
}

// SendOpenWxHB 发送领取红包请求
func GetOpenWxHBReq(userInfo *baseinfo.UserInfo, hongBaoOpenItem *baseinfo.HongBaoOpenItem) []byte {
	var request wechat.HongBaoReq

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &hongBaoOpenItem.CgiCmd

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	strReqText := string("")
	strReqText = strReqText + "channelId=" + hongBaoOpenItem.HongBaoURLItem.ChannelID + "&"
	strReqText = strReqText + "city=" + hongBaoOpenItem.City + "&"
	strReqText = strReqText + "encrypt_key=" + baseutils.EscapeURL(userInfo.HBAesKeyEncrypted) + "&"
	strReqText = strReqText + "encrypt_userinfo=" + baseutils.EscapeURL(GetEncryptUserInfo(userInfo)) + "&"
	strReqText = strReqText + "headImg=" + baseutils.EscapeURL(hongBaoOpenItem.HeadImg) + "&"
	strReqText = strReqText + "msgType=" + hongBaoOpenItem.HongBaoURLItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoOpenItem.NativeURL) + "&"
	strReqText = strReqText + "nickName=" + baseutils.HongBaoStringToBytes(hongBaoOpenItem.NickName) + "&"
	strReqText = strReqText + "province=" + hongBaoOpenItem.Province + "&"
	strReqText = strReqText + "sendId=" + hongBaoOpenItem.HongBaoURLItem.SendID + "&"
	strReqText = strReqText + "sessionUserName=" + hongBaoOpenItem.HongBaoURLItem.SendUserName + "&"
	strReqText = strReqText + "timingIdentifier=" + hongBaoOpenItem.TimingIdentifier
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeOpenWxHB, 5)
	return sendEncodeData
}

// SendOpenWxHB 发送查看红包请求
func GetRedEnvelopeWxHBReq(userInfo *baseinfo.UserInfo, hongBaoOpenItem *baseinfo.HongBaoOpenItem) []byte {
	var request wechat.HongBaoReq

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// CgiCmd
	request.CgiCmd = &hongBaoOpenItem.CgiCmd

	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType

	// ReqText
	strReqText := string("")
	strReqText = strReqText + "agreeDuty=1" + "&"
	strReqText = strReqText + "inWay=1" + "&"
	strReqText = strReqText + "channelId=" + hongBaoOpenItem.HongBaoURLItem.ChannelID + "&"
	strReqText = strReqText + "msgType=" + hongBaoOpenItem.HongBaoURLItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoOpenItem.NativeURL) + "&"
	strReqText = strReqText + "sendId=" + hongBaoOpenItem.HongBaoURLItem.SendID + "&"
	strReqText = strReqText + "sessionUserName=" + hongBaoOpenItem.HongBaoURLItem.SendUserName + "&"
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText
	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeOpenWxHB, 5)
	return sendEncodeData
}

// 查看红包领取列表
func GetRedPacketListReq(userInfo *baseinfo.UserInfo, hongBaoOpenItem *baseinfo.GetRedPacketList) []byte {
	var request wechat.HongBaoReq
	if hongBaoOpenItem.Limit == 0 {
		hongBaoOpenItem.Limit = 11
	}
	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// CgiCmd
	request.CgiCmd = proto.Uint32(5)
	// OutPutType
	outputType := baseinfo.MMTenPayReqOutputTypeJSON
	request.OutPutType = &outputType
	// ReqText
	strReqText := string("")
	strReqText = strReqText + "channelId=" + hongBaoOpenItem.HongBaoItem.ChannelID + "&"
	strReqText = strReqText + "msgType=" + hongBaoOpenItem.HongBaoItem.MsgType + "&"
	strReqText = strReqText + "nativeUrl=" + baseutils.EscapeURL(hongBaoOpenItem.NativeURL) + "&province=&"
	strReqText = strReqText + "sendId=" + hongBaoOpenItem.HongBaoItem.SendID + "&"
	strReqText = strReqText + "limit=" + strconv.FormatInt(hongBaoOpenItem.Limit, 10) + "&"
	strReqText = strReqText + "offset=" + strconv.FormatInt(hongBaoOpenItem.Offset, 10)
	var reqText wechat.SKBuiltinString_
	reqText.Buffer = []byte(strReqText)
	tmpLen := uint32(len(reqText.Buffer))
	reqText.Len = &tmpLen
	request.ReqText = &reqText
	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeQryDetailWxHB, 5)
	return sendEncodeData
}

// SendTextMsg 发送文本消息请求 toWxid：接受人微信id，content：消息内容，atWxIDList：@用户微信id列表（toWxid只能是群的wxid，content应为：@用户昵称 @用户昵称 消息内容）
func GetTextMsReq(userInfo *baseinfo.UserInfo, toWxid string, content string, atWxIDList []string, ContentType int) []byte {
	// 构造请求
	var request wechat.NewSendMsgRequest
	var count uint32 = 1
	request.MsgCount = &count
	var msgRequestNewList = make([]*wechat.MicroMsgRequestNew, count)
	var msgRequestNew wechat.MicroMsgRequestNew
	var recvierString wechat.SKBuiltinString
	recvierString.Str = &toWxid
	msgRequestNew.ToUserName = &recvierString // 设置接收人wxid
	msgRequestNew.Content = &content          // 设置发送内容
	if ContentType == 0 {
		ContentType = 1
	}
	var tmpType = uint32(ContentType)
	msgRequestNew.Type = &tmpType // 发送的类型
	currentTime := time.Now()
	misSecond := currentTime.UnixNano() / 1000000
	var seconds = uint32(misSecond / 1000)
	msgRequestNew.CreateTime = &seconds // 设置时间 秒为单位
	seqID := time.Now().UnixNano() / int64(time.Millisecond)
	var tmpCheckCode = WithSeqidCalcCheckCode(toWxid, seqID)
	msgRequestNew.ClientMsgId = &tmpCheckCode // 设置校验码

	// atUserList
	var atUserStr = ""
	size := len(atWxIDList)
	if size > 0 {
		atUserStr = atUserStr + "<msgsource><atuserlist>"
		for index := int(0); index < size; index++ {
			atUserStr = atUserStr + atWxIDList[index]
			if index < size-1 {
				atUserStr = atUserStr + ","
			}
		}
		atUserStr = atUserStr + "</atuserlist></msgsource>"
		log.Println(atUserStr)
		msgRequestNew.MsgSource = &atUserStr
	}
	msgRequestNewList[0] = &msgRequestNew
	request.ChatSendList = msgRequestNewList

	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeNewSendMsg, 5)
	return sendData
}

// 发送图片v1.1 todo 待定
func GetUploadImageNewReq(userInfo *baseinfo.UserInfo, imgData []byte, toUserName string) []byte {
	// 构造请求
	imgStream := bytes.NewBuffer(imgData)
	Startpos := 0
	datalen := 50000
	datatotalength := imgStream.Len()
	ClientImgId := fmt.Sprintf("%v_%v", userInfo.WxId, time.Now().Unix())
	I := 0
	for {
		Startpos = I * datalen
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}
		Databuff := make([]byte, count)
		_, _ = imgStream.Read(Databuff)
		request := &wechat.UploadMsgImgRequest{
			BaseRequest: GetBaseRequest(userInfo),
			ClientImgId: &wechat.SKBuiltinString{
				Str: proto.String(ClientImgId),
			},
			SenderWxid: &wechat.SKBuiltinString{
				Str: proto.String(userInfo.WxId),
			},
			RecvWxid: &wechat.SKBuiltinString{
				Str: proto.String(toUserName),
			},
			TotalLen: proto.Uint32(uint32(datatotalength)),
			StartPos: proto.Uint32(uint32(Startpos)),
			DataLen:  proto.Uint32(uint32(len(Databuff))),
			Data: &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			MsgType:    proto.Uint32(3),
			EncryVer:   proto.Uint32(0),
			ReqTime:    proto.Uint32(uint32(time.Now().Unix())),
			MessageExt: proto.String("png"),
		}
		//序列化
		srcData, _ := proto.Marshal(request)
		sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeForwardCdnImage, 5)
		I++
		return sendData
	}
	return nil
}

// 发送企业oplog
func GetQWOpLogReq(userInfo *baseinfo.UserInfo, cmdId int64, value []byte) []byte {
	var request wechat.QYOpLogRequest
	request.Type = proto.Int64(cmdId)
	request.V = value
	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, 806, 5)
	return sendData
}

// SendOplogRequest 发送修改帐号信息请求
func GetOplogReq(userInfo *baseinfo.UserInfo, modifyItems []*baseinfo.ModifyItem) []byte {
	var request wechat.OplogRequest
	// CmdList
	var oplog wechat.CmdList
	count := uint32(len(modifyItems))
	oplog.Count = &count

	// ItemList
	cmdItemList := make([]*wechat.CmdItem, count)
	var index = uint32(0)
	for ; index < count; index++ {
		//Item
		cmdItem := &wechat.CmdItem{}
		cmdItem.CmdId = &modifyItems[index].CmdID

		cmdBuf := &wechat.DATA{}
		cmdBuf.Len = &modifyItems[index].Len
		cmdBuf.Data = modifyItems[index].Data
		cmdItem.CmdBuf = cmdBuf
		cmdItemList[index] = cmdItem
	}
	oplog.ItemList = cmdItemList
	request.Oplog = &oplog

	// 发送请求
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeOplog, 5)
	return sendData
}

// SendGetQRCodeRequest 获取二维码
func GetQRCodeReq(userInfo *baseinfo.UserInfo, userName string) []byte {
	var request wechat.GetQRCodeRequest

	// baserequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// opCode
	opcode := uint32(0)
	request.Opcode = &opcode

	// style
	style := uint32(0)
	request.Style = &style

	// UserName
	var userNameSKBuffer wechat.SKBuiltinString
	userNameSKBuffer.Str = &userName
	request.UserName = &userNameSKBuffer

	// 打包发送数据
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeGetQrCode, 5)
	return sendData
}

// SendLogOutRequest 发送登出请求
func GetLogOutReq(userInfo *baseinfo.UserInfo) []byte {
	var request wechat.LogOutRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq
	// 打包数据
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, 282, 5)
	return sendData
}

// SendSnsPostRequest 发送朋友圈
func GetSnsPostReq(userInfo *baseinfo.UserInfo, postItem *baseinfo.SnsPostItem) []byte {
	var request wechat.SnsPostRequest
	zeroValue32 := uint32(0)
	zeroValue64 := uint64(0)

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// ObjectDesc
	objectDescData := CreateSnsPostItemXML(userInfo.WxId, postItem)
	if postItem.Xml {
		objectDescData = []byte(postItem.Content)
	}

	length := uint32(len(objectDescData))
	var objectDesc wechat.SKBuiltinString_
	objectDesc.Len = &length
	objectDesc.Buffer = objectDescData
	request.ObjectDesc = &objectDesc

	// WithUserListCount
	withUserListCount := uint32(len(postItem.WithUserList))
	request.WithUserListCount = &withUserListCount
	// WithUserList
	request.WithUserList = make([]*wechat.SKBuiltinString, withUserListCount)
	index := uint32(0)
	for ; index < withUserListCount; index++ {
		withUser := &wechat.SKBuiltinString{}
		withUser.Str = &postItem.WithUserList[index]
		request.WithUserList[index] = withUser
	}

	// BlackListCount
	blackListCount := uint32(len(postItem.BlackList))
	request.BlackListCount = &blackListCount
	// BlackList
	request.BlackList = make([]*wechat.SKBuiltinString, blackListCount)
	index = uint32(0)
	for ; index < blackListCount; index++ {
		blackUser := &wechat.SKBuiltinString{}
		blackUser.Str = &postItem.BlackList[index]
		request.BlackList[index] = blackUser
	}

	// GroupUserCount
	groupUserCount := uint32(len(postItem.GroupUserList))
	request.GroupUserCount = &groupUserCount
	// GroupUser
	request.GroupUser = make([]*wechat.SKBuiltinString, groupUserCount)
	index = uint32(0)
	for ; index < groupUserCount; index++ {
		groupUser := &wechat.SKBuiltinString{}
		groupUser.Str = &postItem.GroupUserList[index]
		request.GroupUser[index] = groupUser
	}

	// otherFields
	bgImageType := uint32(1)
	request.PostBgimgType = &bgImageType
	request.ObjectSource = &zeroValue32
	request.ReferId = &zeroValue64
	request.Privacy = &postItem.Privacy
	request.SyncFlag = &zeroValue32

	// ClientId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	clientID := string("sns_post_")
	clientID = clientID + userInfo.WxId + "_" + tmpTimeStr + "_0"
	request.ClientId = &clientID

	// groupCount
	request.GroupCount = &zeroValue32
	request.GroupIds = make([]*wechat.SnsGroup, zeroValue32)

	// mediaInfoCount MediaInfo
	mediaInfoCount := uint32(len(postItem.MediaList))
	request.MediaInfoCount = &mediaInfoCount
	request.MediaInfo = make([]*wechat.MediaInfo, mediaInfoCount)
	for index := uint32(0); index < mediaInfoCount; index++ {
		mediaInfo := &wechat.MediaInfo{}
		source := uint32(2)
		mediaInfo.Source = &source

		// MediaType
		mediaType := uint32(1)
		if postItem.MediaList[index].Type == baseinfo.MMSNSMediaTypeImage {
			mediaType = 1
		}
		mediaInfo.MediaType = &mediaType

		// VideoPlayLength
		mediaInfo.VideoPlayLength = &zeroValue32

		// SessionId
		currentTime := int(time.Now().UnixNano() / 1000000)
		sessionID := "memonts-" + strconv.Itoa(currentTime)
		mediaInfo.SessionId = &sessionID

		// startTime
		startTime := uint32(time.Now().UnixNano() / 1000000000)
		mediaInfo.StartTime = &startTime

		request.MediaInfo[index] = mediaInfo
	}

	// SnsPostOperationFields
	var postOperationFields wechat.SnsPostOperationFields
	postOperationFields.ContactTagCount = &zeroValue32
	postOperationFields.TempUserCount = &zeroValue32
	request.SnsPostOperationFields = &postOperationFields

	// clientcheckdata
	var extSpamInfo wechat.SKBuiltinString_
	if userInfo.DeviceInfo != nil {
		extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	} else {
		extSpamInfo.Buffer = GetExtPBSpamInfoDataA16(userInfo)
	}
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	request.ExtSpamInfo = &extSpamInfo

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsPost, 5)
	return sendEncodeData
}

// SendSnsPostRequestByXML 通过XML的信息来发送朋友圈
func GetSnsPostReqByXML(userInfo *baseinfo.UserInfo, timeLineObj *baseinfo.TimelineObject, blackList []string) ([]byte, error) {
	var request wechat.SnsPostRequest
	zeroValue32 := uint32(0)
	zeroValue64 := uint64(0)

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// WithUserListCount
	withUserListCount := uint32(0)
	request.WithUserListCount = &withUserListCount
	// WithUserList
	request.WithUserList = make([]*wechat.SKBuiltinString, withUserListCount)

	// BlackListCount
	tmpCount := uint32(len(blackList))
	request.BlackListCount = &tmpCount
	request.BlackList = make([]*wechat.SKBuiltinString, tmpCount)
	// BlackList
	for index := uint32(0); index < tmpCount; index++ {
		tmpSKBuiltinString := &wechat.SKBuiltinString{}
		tmpSKBuiltinString.Str = &blackList[index]
		request.BlackList[index] = tmpSKBuiltinString
	}

	// GroupUserCount
	groupUserCount := uint32(0)
	request.GroupUserCount = &groupUserCount
	// GroupUser
	request.GroupUser = make([]*wechat.SKBuiltinString, groupUserCount)

	// otherFields
	bgImageType := uint32(1)
	request.PostBgimgType = &bgImageType
	request.ObjectSource = &zeroValue32
	request.ReferId = &zeroValue64
	request.Privacy = &timeLineObj.Private
	request.SyncFlag = &zeroValue32

	// ClientId
	tmpTime := int(time.Now().UnixNano() / 1000000000)
	tmpTimeStr := strconv.Itoa(tmpTime)
	clientID := string("sns_post_")
	clientID = clientID + userInfo.WxId + "_" + tmpTimeStr + "_0"
	request.ClientId = &clientID

	// groupCount
	request.GroupCount = &zeroValue32
	request.GroupIds = make([]*wechat.SnsGroup, zeroValue32)

	// mediaInfoCount MediaInfo
	mediaInfoCount := uint32(len(timeLineObj.ContentObject.MediaList.Media))
	request.MediaInfoCount = &mediaInfoCount
	request.MediaInfo = make([]*wechat.MediaInfo, mediaInfoCount)
	for index := uint32(0); index < mediaInfoCount; index++ {
		tmpMediaItem := timeLineObj.ContentObject.MediaList.Media[index]

		// 解析Source
		mediaInfo := &wechat.MediaInfo{}
		tmpSource := baseutils.ParseInt(tmpMediaItem.URL.Type)
		mediaInfo.Source = &tmpSource
		// MediaType
		mediaType := tmpMediaItem.Type - 1
		mediaInfo.MediaType = &mediaType
		// VideoPlayLength
		playLength := uint32(tmpMediaItem.VideoDuration)
		mediaInfo.VideoPlayLength = &playLength
		// SessionId
		currentTime := int(time.Now().UnixNano() / 1000000)
		sessionID := "memonts-" + strconv.Itoa(currentTime)
		mediaInfo.SessionId = &sessionID

		// startTime
		startTime := uint32(time.Now().UnixNano() / 1000000000)
		mediaInfo.StartTime = &startTime
		request.MediaInfo[index] = mediaInfo
	}

	// ID和UserName置为0
	timeLineObj.UserName = userInfo.WxId
	timeLineObj.CreateTime = uint32(int(time.Now().UnixNano() / 1000000000))
	// ObjectDesc
	objectDescData, err := xml.Marshal(timeLineObj)
	if err != nil {
		return nil, err
	}
	str := string(objectDescData)
	str = strings.ReplaceAll(str, "token=\"\"", "")
	str = strings.ReplaceAll(str, "key=\"\"", "")
	str = strings.ReplaceAll(str, "enc_idx=\"\"", "")
	str = strings.ReplaceAll(str, "md5=\"\"", "")
	str = strings.ReplaceAll(str, "videomd5=\"\"", "")
	str = strings.ReplaceAll(str, "video", "")
	objectDescData = []byte(str)
	length := uint32(len(objectDescData))
	var objectDesc wechat.SKBuiltinString_
	objectDesc.Len = &length
	objectDesc.Buffer = objectDescData
	request.ObjectDesc = &objectDesc

	// SnsPostOperationFields
	var postOperationFields wechat.SnsPostOperationFields
	postOperationFields.ContactTagCount = &zeroValue32
	postOperationFields.TempUserCount = &zeroValue32
	request.SnsPostOperationFields = &postOperationFields

	// clientcheckdata
	var extSpamInfo wechat.SKBuiltinString_
	if userInfo.DeviceInfo != nil {
		extSpamInfo.Buffer = GetExtPBSpamInfoData(userInfo)
	} else {
		extSpamInfo.Buffer = GetExtPBSpamInfoDataA16(userInfo)
	}
	extSpamInfoLen := uint32(len(extSpamInfo.Buffer))
	extSpamInfo.Len = &extSpamInfoLen
	request.ExtSpamInfo = &extSpamInfo

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsPost, 5)
	return sendEncodeData, nil
}

// SendSnsObjectOpRequest 发送朋友圈操作
func GetSnsObjectOpReq(userInfo *baseinfo.UserInfo, opItems []*baseinfo.SnsObjectOpItem) []byte {
	var request wechat.SnsObjectOpRequest

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// OpCount
	opCount := uint32(len(opItems))
	request.OpCount = &opCount

	// OpList
	request.OpList = make([]*wechat.SnsObjectOp, opCount)
	index := uint32(0)
	for ; index < opCount; index++ {
		snsObject := &wechat.SnsObjectOp{}
		id, _ := strconv.ParseUint(opItems[index].SnsObjID, 0, 64)
		snsObject.Id = &id
		snsObject.OpType = &opItems[index].OpType
		if opItems[index].DataLen > 0 {
			skBuffer := &wechat.SKBuiltinString_{}
			skBuffer.Len = &opItems[index].DataLen
			skBuffer.Buffer = opItems[index].Data
		}
		if opItems[index].Ext != 0 {
			extInfo := &wechat.SnsObjectOpExt{
				Id: &opItems[index].Ext,
			}
			CommnetId, _ := proto.Marshal(extInfo)
			snsObject.Ext = &wechat.SKBuiltinString_{
				Len:    proto.Uint32(uint32(len(CommnetId))),
				Buffer: CommnetId,
			}
		}
		request.OpList[index] = snsObject
	}

	// 打包数据 发送
	srcData, _ := proto.Marshal(&request)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsObjectOp, 5)
	return sendEncodeData
}

// SendSnsUserPageRequest 发送 获取朋友圈信息 请求
func GetSnsUserPageReq(userInfo *baseinfo.UserInfo, userName string, firstPageMd5 string, maxID uint64) []byte {
	var request wechat.SnsUserPageRequest
	var zeroValue64 = uint64(0)
	var zeroValue32 = uint32(0)

	// baseRequest
	baseReq := GetBaseRequest(userInfo)
	var tmpScene = uint32(0)
	baseReq.Scene = &tmpScene
	request.BaseRequest = baseReq

	// 其它参数
	request.Username = &userName
	request.FirstPageMd5 = &firstPageMd5
	request.MaxId = &maxID
	request.MinFilterId = &zeroValue64
	request.LastRequestTime = &zeroValue32
	request.FilterType = &zeroValue32

	// 打包数据
	srcData, _ := proto.Marshal(&request)
	sendData := Pack(userInfo, srcData, baseinfo.MMRequestTypeMMSnsUserPage, 5)

	return sendData
}
