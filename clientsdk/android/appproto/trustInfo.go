package appproto

import (
	"crypto/rand"
	"feiyu.com/wx/clientsdk/android"
	"feiyu.com/wx/clientsdk/android/mmproto"
	"feiyu.com/wx/clientsdk/baseinfo"
	proto "github.com/golang/protobuf/proto"
	"io"
	"time"
)

type TrustInfo struct {
	uri   string
	cmdid uint32
}

//MakeTrustDeviceReq s

func (t *TrustInfo) GetUri() string {
	return "/cgi-bin/micromsg-bin/fpfreshnl"
}

func (t *TrustInfo) GetCmdid() uint32 {
	return 0
}

func (t *TrustInfo) MakeTrustDeviceReq() []byte {

	trustInfo := &mmproto.TrustReq{
		Td: &mmproto.TrustData{
			Tdi: []*mmproto.TrustDeviceInfo{
				{Key: proto.String("IMEI"), Val: proto.String("353627078088849")},
				{Key: proto.String("AndroidID"), Val: proto.String("06a78780bc297bbd")},
				{Key: proto.String("PhoneSerial"), Val: proto.String("01c5cded725f4db6")},
				{Key: proto.String("cid"), Val: proto.String("")},
				{Key: proto.String("WidevineDeviceID"), Val: proto.String("657a6b657b4f79614563447b54796c526f67724c466b79644564634c45675600")},
				{Key: proto.String("WidevineProvisionID"), Val: proto.String("955e20f6b905cbadbe67a580129b8f36")},
				{Key: proto.String("GSFID"), Val: proto.String("")},
				{Key: proto.String("SoterID"), Val: proto.String("")},
				{Key: proto.String("SoterUid"), Val: proto.String("")},
				{Key: proto.String("FSID"), Val: proto.String("3706372|3706398@8c829c5f2697bfed|2fe1cc4100a798d0b60909e0ea1090e7@d3609fe804970d6b|2ee5e2decd893fffee73ceab66e30640")},
				{Key: proto.String("BootID"), Val: proto.String("")},
				{Key: proto.String("IMSI"), Val: proto.String("")},
				{Key: proto.String("PhoneNum"), Val: proto.String("")},
				{Key: proto.String("WeChatInstallTime"), Val: proto.String("1515061151")},
				{Key: proto.String("PhoneModel"), Val: proto.String("Nexus 5X")},
				{Key: proto.String("BuildBoard"), Val: proto.String("bullhead")},
				{Key: proto.String("BuildBootloader"), Val: proto.String("BHZ32c")},
				{Key: proto.String("SystemBuildDate"), Val: proto.String("Fri Sep 28 23:37:27 UTC 2018")},
				{Key: proto.String("SystemBuildDateUTC"), Val: proto.String("1538177847")},
				{Key: proto.String("BuildFP"), Val: proto.String("google/bullhead/bullhead:8.1.0/OPM7.181105.004/5038062:user/release-keys")},
				{Key: proto.String("BuildID"), Val: proto.String("OPM7.181105.004")},
				{Key: proto.String("BuildBrand"), Val: proto.String("google")},
				{Key: proto.String("BuildDevice"), Val: proto.String("bullhead")},
				{Key: proto.String("BuildProduct"), Val: proto.String("bullhead")},
				{Key: proto.String("Manufacturer"), Val: proto.String("LGE")},
				{Key: proto.String("RadioVersion"), Val: proto.String("M8994F-2.6.42.5.03")},
				{Key: proto.String("AndroidVersion"), Val: proto.String("8.1.0")},
				{Key: proto.String("SdkIntVersion"), Val: proto.String("27")},
				{Key: proto.String("ScreenWidth"), Val: proto.String("1080")},
				{Key: proto.String("ScreenHeight"), Val: proto.String("1794")},
				{Key: proto.String("SensorList"), Val: proto.String("BMI160 accelerometer#Bosch#0.004788#1,BMI160 gyroscope#Bosch#0.000533#1,BMM150 magnetometer#Bosch#0.000000#1,BMP280 pressure#Bosch#0.005000#1,BMP280 temperature#Bosch#0.010000#1,RPR0521 Proximity Sensor#Rohm#1.000000#1,RPR0521 Light Sensor#Rohm#10.000000#1,Orientation#Google#1.000000#1,BMI160 Step detector#Bosch#1.000000#1,Significant motion#Google#1.000000#1,Gravity#Google#1.000000#1,Linear Acceleration#Google#1.000000#1,Rotation Vector#Google#1.000000#1,Geomagnetic Rotation Vector#Google#1.000000#1,Game Rotation Vector#Google#1.000000#1,Pickup Gesture#Google#1.000000#1,Tilt Detector#Google#1.000000#1,BMI160 Step counter#Bosch#1.000000#1,BMM150 magnetometer (uncalibrated)#Bosch#0.000000#1,BMI160 gyroscope (uncalibrated)#Bosch#0.000533#1,Sensors Sync#Google#1.000000#1,Double Twist#Google#1.000000#1,Double Tap#Google#1.000000#1,Device Orientation#Google#1.000000#1,BMI160 accelerometer (uncalibrated)#Bosch#0.004788#1")},
				{Key: proto.String("DefaultInputMethod"), Val: proto.String("com.google.android.inputmethod.latin")},
				{Key: proto.String("InputMethodList"), Val: proto.String("Google \345\215\260\345\272\246\350\257\255\351\224\256\347\233\230#com.google.android.apps.inputmethod.hindi,Google \350\257\255\351\237\263\350\276\223\345\205\245#com.google.android.googlequicksearchbox,Google \346\227\245\350\257\255\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.japanese,Google \351\237\251\350\257\255\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.korean,Gboard#com.google.android.inputmethod.latin,\350\260\267\346\255\214\346\213\274\351\237\263\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.pinyin")},
				{Key: proto.String("DeviceID"), Val: proto.String("A0e4a76905e8f67b")},
				{Key: proto.String("OAID"), Val: proto.String("")},
			},
		},
		Md: proto.String("e05ac1f886668063fabe3231fd78a2cb"),
	}

	pb, _ := proto.Marshal(trustInfo)
	//log.Printf("%x\n", pb)

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(pb)

	//log.Printf("ZT: %x\n", encData)

	randKey := make([]byte, 16)
	io.ReadFull(rand.Reader, randKey)

	fp := &mmproto.FPFresh{
		BaseReq: &mmproto.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint64(0),
			DeviceID:      append([]byte("A0e4a76905e8f67"), 0),
			ClientVersion: proto.Int32(int32(baseinfo.AndroidClientVersion)),
			DeviceType:    proto.String("android-27"),
			Scene:         proto.Uint32(0),
		},
		SessKey: randKey,
		Ztdata: &mmproto.ZTData{
			Version:   proto.String("00000003\x00"),
			Encrypted: proto.Uint32(1),
			Data:      encData,
			TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
			Optype:    proto.Uint32(5),
			Uin:       proto.Uint32(0),
		},
	}

	fpPB, _ := proto.Marshal(fp)
	return fpPB
}
