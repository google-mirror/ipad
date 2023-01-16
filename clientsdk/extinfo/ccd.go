package extinfo

import (
	"feiyu.com/wx/clientsdk/android"
	"feiyu.com/wx/clientsdk/android/mmproto"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/golang/protobuf/proto"
	"hash/crc32"
	"math/rand"
	"time"
)

// 安卓23算法
func AndroidCcData(DeviceId string, info *baseinfo.AndroidDeviceInfo, DeviceToken *wechat.TrustResp, T int64) *wechat.ZTData {
	microseconds3 := time.Now().UnixNano()%1000 + time.Now().Unix()%1000*1000
	tvSec := time.Now().Unix() / 1000
	tvUsec := time.Now().UnixNano()%1000 + time.Now().Unix()%1000*1000
	//microseconds4:=time.Now().UnixNano()%1000+time.Now().Unix()%1000*1000
	defaultUid := 10200 + rand.Int63n(2000)

	ccd3body := &wechat.SpamAndroidBody{
		Loc:                  proto.Uint32(0),
		Root:                 proto.Uint32(0),
		Debug:                proto.Uint32(0),
		PackageSign:          proto.String(info.AndriodPackageSign(DeviceId)),
		RadioVersion:         proto.String(info.AndroidRadioVersion(DeviceId)),
		BuildVersion:         proto.String(info.AndroidVersion()),
		DeviceId:             proto.String(info.AndriodImei(DeviceId)),
		AndroidId:            proto.String(info.AndriodID(DeviceId)),
		SerialId:             proto.String(info.AndriodPhoneSerial(DeviceId)),
		Model:                proto.String(info.AndroidPhoneModel(DeviceId)),
		CpuCount:             proto.Uint32(6),
		CpuBrand:             proto.String(info.AndroidHardware(DeviceId)),
		CpuExt:               proto.String(info.AndroidFeatures()),
		WlanAddress:          proto.String(info.AndriodWLanAddress(DeviceId)),
		Ssid:                 proto.String(info.AndriodSsid(DeviceId)),
		Bssid:                proto.String(info.AndriodBssid(DeviceId)),
		SimOperator:          proto.String(""),
		WifiName:             proto.String(info.AndroidWifiName(DeviceId)),
		BuildFP:              proto.String(info.AndroidBuildFP(DeviceId)),
		BuildBoard:           proto.String("bullhead"),
		BuildBootLoader:      proto.String(info.AndroidBuildBoard(DeviceId)),
		BuildBrand:           proto.String("google"),
		BuildDevice:          proto.String("bullhead"),
		DataDir:              proto.String("/data/user/0/com.tencent.mm/"),
		NetType:              proto.String("wifi"),
		PackageName:          proto.String("com.tencent.mm"),
		Task:                 proto.Uint64(0),
		GsmSimOperatorNumber: proto.String(""),
		SoterId:              proto.String(""),
		KernelReleaseNumber:  proto.String(info.AndroidKernelReleaseNumber(DeviceId)),
		UsbState:             proto.Uint64(0),
		Sign:                 proto.String(info.AndriodPackageSign(DeviceId)),
		PackageFlag:          proto.Uint64(14),
		AccessFlag:           proto.Uint64(uint64(info.AndriodAccessFlag(DeviceId))),
		Unkonwn:              proto.Uint64(3),
		TbVersionCrc:         proto.Uint64(uint64(info.AndriodTbVersionCrc(DeviceId))),
		SfMD5:                proto.String(info.AndriodSfMD5(DeviceId)),
		SfArmMD5:             proto.String(info.AndriodSfArmMD5(DeviceId)),
		SfArm64MD5:           proto.String(info.AndriodSfArm64MD5(DeviceId)),
		SbMD5:                proto.String(info.AndriodSbMD5(DeviceId)),
		SoterId2:             proto.String(""),
		WidevineDeviceID:     proto.String(info.AndriodWidevineDeviceID(DeviceId)),
		FSID:                 proto.String(info.AndriodFSID(DeviceId)),
		Oaid:                 proto.String(""),
		TimeCheck:            proto.Uint64(0),
		NanoTime:             proto.Uint64(uint64(info.AndriodNanoTime(DeviceId))),
		Refreshtime:          proto.Uint64(DeviceToken.GetTrustResponseData().GetTimestamp()),
		SoftConfig:           proto.String(DeviceToken.GetTrustResponseData().GetSoftData().GetSoftConfig()),
		SoftData:             DeviceToken.GetTrustResponseData().GetSoftData().GetSoftData(),
		DebugFlags:           proto.Uint64(uint64(microseconds3)),
		RouteIFace:           proto.String("eth0"), //this.disableWifi ? "eth0" : "wlan0"
		TvSec:                proto.Uint64(uint64(tvSec)),
		TvUsec:               proto.Uint64(uint64(tvUsec)),
		//TvCheck:			  proto.Uint64(0), //???
		//PkgHash3Encrypted:	  proto.EncodeVarint(0),
		Uid: proto.Uint64(uint64(defaultUid)),
	}
	//

	pb, _ := proto.Marshal(ccd3body)

	crc := crc32.ChecksumIEEE(pb)

	curtime := uint32(T)

	ccd3 := &mmproto.Ccd3{
		Crc:       &crc,
		TimeStamp: &curtime,
		Body:      ccd3body,
	}

	pb, _ = proto.Marshal(ccd3)

	zt := new(android.ZT)
	zt.Init()
	encData := zt.WBAesEncrypt(pb)

	Ztdata := &wechat.ZTData{
		Version:   []byte("00000006"),
		Encrypted: proto.Uint64(1),
		Data:      encData,
		TimeStamp: &T,
		OpType:    proto.Uint64(5),
		Uin:       proto.Uint64(0),
	}
	return Ztdata
}
