package baseinfo

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"feiyu.com/wx/clientsdk/android"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
	"fmt"
	"github.com/lunny/log"
	"regexp"
	"strconv"
	"strings"
)

type AndroidDeviceInfo struct {
	DeviceId            []byte
	DeviceIdStr         string
	DeviceToken         *wechat.TrustResp
	Imei                string
	AndriodId           string
	PhoneSerial         string
	WidevineDeviceID    string
	WidevineProvisionID string
	AndriodFsId         string
	AndriodBssId        string
	AndriodSsId         string
	WLanAddress         string
	PackageSign         string
	Androidversion      string
	RadioVersion        string
	Manufacturer        string
	BuildID             string
	BuildFP             string
	BuildBoard          string
	PhoneModel          string
	Hardware            string
	Features            string
	WifiName            string
	WifiFullName        string
	KernelReleaseNumber string
	Arch                string
	SfMD5               string
	SfArmMD5            string
	SfArm64MD5          string
	SbMD5               string
}

func (Info *AndroidDeviceInfo) AndriodImei(DeviceId string) string {
	if Info.Imei != "" && Info.Imei != "string" {
		return Info.Imei
	}
	Md5Data := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000000"))))
	return "35" + Md5Data[0:13]
}

func (Info *AndroidDeviceInfo) AndriodID(DeviceId string) string {
	if Info.AndriodId != "" && Info.AndriodId != "string" {
		return Info.AndriodId
	}
	Md5Data := baseutils.MD5ToLower(DeviceId + "SM1000001")
	return "06" + Md5Data[0:14]
}

func (Info *AndroidDeviceInfo) AndriodPhoneSerial(DeviceId string) string {
	if Info.PhoneSerial != "" && Info.PhoneSerial != "string" {
		return Info.PhoneSerial
	}
	//return "01c5cded725f4db6"
	Md5Data := baseutils.MD5ToLower(DeviceId + "SM1000002")
	return "01" + Md5Data[0:14]
}

func (Info *AndroidDeviceInfo) AndriodWidevineDeviceID(DeviceId string) string {
	if Info.WidevineDeviceID != "" && Info.WidevineDeviceID != "string" {
		return Info.WidevineDeviceID
	}
	Md5DataA := baseutils.MD5ToLower(DeviceId + "SM1000003")
	Md5DataB := baseutils.MD5ToLower(DeviceId + "SM1000004")
	return "657" + Md5DataA[0:29] + Md5DataB
}

func (Info *AndroidDeviceInfo) AndriodWidevineProvisionID(DeviceId string) string {
	if Info.WidevineProvisionID != "" && Info.WidevineProvisionID != "string" {
		return Info.WidevineProvisionID
	}
	Md5DataA := baseutils.MD5ToLower(DeviceId + "SM1000005")
	return "955" + Md5DataA[0:29]
}

func (Info *AndroidDeviceInfo) AndriodFSID(DeviceId string) string {
	if Info.AndriodFsId != "" && Info.AndriodFsId != "string" {
		return Info.AndriodFsId
	}
	Md5DataA := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000012"))))
	Md5DataB := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000006"))))
	return "37063" + Md5DataA[0:2] + "|37063" + Md5DataA[2:4] + "@" + Md5DataA[4:19] + "|" + baseutils.MD5ToLower(DeviceId+"SM1000007") + "@" + Md5DataB[0:16] + baseutils.MD5ToLower(DeviceId+"SM1000008")
}

func (Info *AndroidDeviceInfo) AndriodBssid(DeviceId string) string {
	if Info.AndriodBssId != "" && Info.AndriodBssId != "string" {
		return Info.AndriodBssId
	}
	return "02:00:00:00:00:00"
	/*Md5Data := baseutils.MD5ToLower(DeviceId + "SM1000009")
	A := Md5Data[5:7] + ":"
	B := Md5Data[7:9] + ":"
	C := Md5Data[9:11] + ":"
	D := Md5Data[11:13] + ":"
	E := Md5Data[13:15] + ":"
	F := Md5Data[15:17]
	return A + B + C + D + E + F*/
}

func (Info *AndroidDeviceInfo) AndriodSsid(DeviceId string) string {
	if Info.AndriodSsId != "" && Info.AndriodSsId != "string" {
		return Info.AndriodSsId
	}
	//return "02:00:00:00:00:00"
	Md5Data := baseutils.MD5ToLower(DeviceId + "SM10000026")
	A := Md5Data[5:7] + ":"
	B := Md5Data[7:9] + ":"
	C := Md5Data[9:11] + ":"
	D := Md5Data[11:13] + ":"
	E := Md5Data[13:15] + ":"
	F := Md5Data[15:17]
	return A + B + C + D + E + F
}

func (Info *AndroidDeviceInfo) AndriodWLanAddress(DeviceId string) string {
	if Info.WLanAddress != "" && Info.WLanAddress != "string" {
		return Info.WLanAddress
	}
	//return "00:a0:07:86:17:18"
	Md5Data := baseutils.MD5ToLower(DeviceId + "SM1000009")
	B := Md5Data[7:9] + ":"
	C := Md5Data[9:11] + ":"
	D := Md5Data[11:13] + ":"
	E := Md5Data[13:15] + ":"
	F := Md5Data[15:17]
	return "00:" + B + C + D + E + F
}

func (Info *AndroidDeviceInfo) AndriodPackageSign(DeviceId string) string {
	if Info.PackageSign != "" && Info.PackageSign != "string" {
		return Info.PackageSign
	}
	/*Md5Data := baseutils.MD5ToLower(DeviceId + "SM1000010")
	return "18" + Md5Data[0:30]*/
	return "18c867f0717aa67b2ab7347505ba07ed"
}

func (Info *AndroidDeviceInfo) AndroidVersion() string {
	if Info.Androidversion != "" && Info.Androidversion != "string" {
		return Info.Androidversion
	}
	return "8.1.0"
}

func (Info *AndroidDeviceInfo) AndroidRadioVersion(DeviceId string) string {
	if Info.RadioVersion != "" && Info.RadioVersion != "string" {
		return Info.RadioVersion
	}
	return "M8994F-2.6.42.5.03"
	/*S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000013"))))
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000014")+baseutils.MD5ToLower(DeviceId+"SM2000014")+baseutils.MD5ToLower(DeviceId+"SM3000014"), -1))
	L := fmt.Sprintf("%v", M[:1])
	K := fmt.Sprintf("%v", M[:2])
	T := fmt.Sprintf("%v%v%v-2.5.%v.%v.%v", strings.ToUpper(L), S[:4], strings.ToUpper(K), S[5:7], S[7:8], S[8:10])
	return T*/
}

func (Info *AndroidDeviceInfo) AndroidManufacturer(DeviceId string) string {
	if Info.Manufacturer != "" && Info.Manufacturer != "string" {
		return Info.Manufacturer
	}
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000015")+baseutils.MD5ToLower(DeviceId+"SM2000015"), -1))
	L := fmt.Sprintf("%v", M[:3])
	return strings.ToUpper(L)
}

func (Info *AndroidDeviceInfo) AndroidBuildID(DeviceId string) string {
	if Info.BuildID != "" && Info.BuildID != "string" {
		return Info.BuildID
	}
	return "06a78780bc297bbd"
	/*S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000016"))))
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000017")+baseutils.MD5ToLower(DeviceId+"SM2000017"), -1))
	L := fmt.Sprintf("%v%v", M[:3], S[:1])
	N := fmt.Sprintf("%v.%v.%v", strings.ToUpper(L), S[2:8], S[9:12])
	return N*/
}

func (Info *AndroidDeviceInfo) AndroidBuildFP(DeviceId string) string {
	if Info.BuildFP != "" && Info.BuildFP != "string" {
		return Info.BuildFP
	}
	return "google/bullhead/bullhead:8.1.0/OPM7.181105.004/5038062:user/release-keys"
	/*S := Info.INCREMENTAL(DeviceId)
	return fmt.Sprintf("google/bullhead/bullhead:%v/%v/%v:user/release-keys", Info.AndroidVersion(), Info.AndroidBuildID(DeviceId), S[:7])*/
}

func (Info *AndroidDeviceInfo) AndroidBuildBoard(DeviceId string) string {
	if Info.BuildBoard != "" && Info.BuildBoard != "string" {
		return Info.BuildBoard
	}
	return "BHZ32c"
	/*S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000021"))))
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000022")+baseutils.MD5ToLower(DeviceId+"SM2000022"), -1))
	G := fmt.Sprintf("%v", M[:3])
	K := fmt.Sprintf("%v", M[3:4])
	return strings.ToUpper(G) + S[:2] + K*/
}

func (Info *AndroidDeviceInfo) AndroidPhoneModel(DeviceId string) string {
	if Info.PhoneModel != "" && Info.PhoneModel != "string" {
		return Info.PhoneModel
	}
	S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000019"))))
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000020")+baseutils.MD5ToLower(DeviceId+"SM2000020"), -1))
	return fmt.Sprintf("%v%v %v", strings.ToUpper(M[:1]), M[2:6], S[:2])
}

func (Info *AndroidDeviceInfo) AndroidHardware(DeviceId string) string {
	if Info.Hardware != "" && Info.Hardware != "string" {
		return Info.Hardware
	}
	return "Qualcomm Technologies, Inc MSM8992"
	/*S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000023"))))
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000024")+baseutils.MD5ToLower(DeviceId+"SM1000025"), -1))
	G := fmt.Sprintf("%v", M[:1])
	H := fmt.Sprintf("%v", M[2:5])
	O := fmt.Sprintf("%v", M[6:9])
	return strings.ToUpper(G) + H + " Technologies, Inc " + O + S[:4]*/
}

func (Info *AndroidDeviceInfo) AndroidFeatures() string {
	if Info.Features != "" && Info.Features != "string" {
		return Info.Features
	}
	return "half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt evtstrm aes pmull sha1 sha2 crc32"
}

func (Info *AndroidDeviceInfo) AndroidWifiName(DeviceId string) string {
	if Info.WifiName != "" && Info.WifiName != "string" {
		return Info.WifiName
	}
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(DeviceId+"SM1000027")+baseutils.MD5ToLower(DeviceId+"SM1000028"), -1))
	return "Chinanet-" + fmt.Sprintf("%v", M[0:5])
}

func (Info *AndroidDeviceInfo) AndroidWifiFullName(DeviceId string) string {
	if Info.WifiFullName != "" && Info.WifiFullName != "string" {
		return Info.WifiFullName
	}
	return fmt.Sprintf("&quot;%v&quot;", Info.AndroidWifiName(DeviceId))
}

func (Info *AndroidDeviceInfo) AndroidKernelReleaseNumber(DeviceId string) string {
	if Info.KernelReleaseNumber != "" && Info.KernelReleaseNumber != "string" {
		return Info.KernelReleaseNumber
	}
	S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000029"))))
	M := baseutils.MD5ToLower(DeviceId + "SM1000027")
	return fmt.Sprintf("%v.%v.%v-%v", S[:1], S[1:3], S[4:6], M[:13])
}

func (Info *AndroidDeviceInfo) AndroidArch(DeviceId string) string {
	if Info.Arch != "" && Info.Arch != "string" {
		return Info.Arch
	}
	M := baseutils.MD5ToLower(DeviceId + "SM1000030")
	return fmt.Sprintf("armeabi-%v", M[:3])
}

func (Info *AndroidDeviceInfo) AndriodSfMD5(DeviceId string) string {
	if Info.SfMD5 != "" && Info.SfMD5 != "string" {
		return Info.SfMD5
	}
	return baseutils.MD5ToLower(DeviceId + "SM1000031")
}

func (Info *AndroidDeviceInfo) AndriodSfArmMD5(DeviceId string) string {
	if Info.SfArmMD5 != "" && Info.SfArmMD5 != "string" {
		return Info.SfArmMD5
	}
	return baseutils.MD5ToLower(DeviceId + "SM1000032")
}

func (Info *AndroidDeviceInfo) AndriodSfArm64MD5(DeviceId string) string {
	if Info.SfArm64MD5 != "" && Info.SfArm64MD5 != "string" {
		return Info.SfArm64MD5
	}
	return baseutils.MD5ToLower(DeviceId + "SM1000033")
}

func (Info *AndroidDeviceInfo) AndriodSbMD5(DeviceId string) string {
	if Info.SbMD5 != "" && Info.SbMD5 != "string" {
		return Info.SbMD5
	}
	return baseutils.MD5ToLower(DeviceId + "SM1000034")
}

func (Info *AndroidDeviceInfo) AndriodAccessFlag(DeviceId string) int {
	S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000035"))))
	B, _ := strconv.Atoi(S[:6])
	return B
}

func (Info *AndroidDeviceInfo) AndriodTbVersionCrc(DeviceId string) int {
	S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000036"))))
	B, _ := strconv.Atoi(S[:9])
	return B
}

func (Info *AndroidDeviceInfo) AndriodNanoTime(DeviceId string) int {
	S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000037"))))
	B, _ := strconv.Atoi(S[:6])
	return B
}

func (Info *AndroidDeviceInfo) INCREMENTAL(DeviceId string) string {
	S := strconv.Itoa(baseutils.BytesToInt([]byte(baseutils.MD5ToLower(DeviceId + "SM1000018"))))
	return S[:7]
}

func (Info *AndroidDeviceInfo) AndriodUUID(DeviceId string) string {
	S := baseutils.MD5ToLower(DeviceId + "SM1000018")
	return fmt.Sprintf("%x-%x-%x-%x-%x", S[:4], S[4:6], S[6:8], S[8:10], S[10:])
}

func (Info *AndroidDeviceInfo) AndriodDeviceType(DeviceId string) string {
	MANUFACTURER := Info.AndroidManufacturer(DeviceId)
	MODEL := Info.AndroidPhoneModel(DeviceId)
	RELEASE := Info.AndroidVersion()
	INCREMENTAL := Info.INCREMENTAL(DeviceId)
	DISPLAY := Info.AndroidBuildID(DeviceId)
	M := fmt.Sprintf("<AndroidDeviceInfo><MANUFACTURER name=\"%v\"><MODEL name=\"%v\"><VERSION_RELEASE name=\"%v\"><VERSION_INCREMENTAL name=\"%v\"><DISPLAY name=\"%v\"></DISPLAY></VERSION_INCREMENTAL></VERSION_RELEASE></MODEL></MANUFACTURER></AndroidDeviceInfo>", MANUFACTURER, MODEL, RELEASE, INCREMENTAL, DISPLAY)
	return M
}

func (Info *AndroidDeviceInfo) AndriodGetSoftType(DeviceId string) string {
	softType := "<softtype><lctmoc>"
	softType += fmt.Sprintf("%d", 0)
	softType += "</lctmoc><level>"
	softType += fmt.Sprintf("%d", 0)
	softType += "</level><k1>"
	softType += "0 "
	softType += "</k1><k2>"
	softType += Info.AndroidRadioVersion(DeviceId)
	softType += "</k2><k3>"
	softType += Info.AndroidVersion()
	softType += "</k3><k4>"
	softType += Info.AndriodImei(DeviceId)
	softType += "</k4><k5>"
	softType += ""
	softType += "</k5><k6>"
	softType += ""
	softType += "</k6><k7>"
	softType += Info.AndriodID(DeviceId)
	softType += "</k7><k8>"
	softType += Info.AndriodPhoneSerial(DeviceId)
	softType += "</k8><k9>"
	softType += Info.AndroidPhoneModel(DeviceId)
	softType += "</k9><k10>"
	softType += fmt.Sprintf("%d", 8)
	softType += "</k10><k11>"
	softType += Info.AndroidHardware(DeviceId)
	softType += "</k11><k12>"
	softType += ""
	softType += "</k12><k13>"
	softType += ""
	softType += "</k13><k14>"
	softType += Info.AndriodSsid(DeviceId)
	softType += "</k14><k15>"
	softType += ""
	softType += "</k15><k16>"
	softType += Info.AndroidFeatures()
	softType += "</k16><k18>"
	softType += Info.AndriodPackageSign(DeviceId)
	softType += "</k18><k21>"
	softType += Info.AndroidWifiName(DeviceId)
	softType += "</k21><k22>"
	softType += ""
	softType += "</k22><k24>"
	softType += Info.AndriodBssid(DeviceId)
	softType += "</k24><k26>"
	softType += fmt.Sprintf("%d", 0)
	softType += "</k26><k30>"
	softType += Info.AndroidWifiFullName(DeviceId)
	softType += "</k30><k33>"
	softType += "com.tencent.mm"
	softType += "</k33><k34>"
	softType += Info.AndroidBuildFP(DeviceId)
	softType += "</k34><k35>"
	softType += "bullhead"
	softType += "</k35><k36>"
	softType += Info.AndroidBuildBoard(DeviceId)
	softType += "</k36><k37>"
	softType += "google"
	softType += "</k37><k38>"
	softType += "bullhead"
	softType += "</k38><k39>"
	softType += "bullhead"
	softType += "</k39><k40>"
	softType += "bullhead"
	softType += "</k40><k41>"
	softType += fmt.Sprintf("%d", 0)
	softType += "</k41><k42>"
	softType += Info.AndroidManufacturer(DeviceId)
	//43 "89884a87498ef44f" setting
	//44 -> 0
	softType += "</k42><k43>null</k43><k44>0</k44><k45>"
	softType += ""
	softType += "</k45><k46>"
	softType += ""
	softType += "</k46><k47>"
	softType += "wifi"
	softType += "</k47><k48>"
	softType += Info.AndriodImei(DeviceId)
	softType += "</k48><k49>"
	softType += "data/user/0/com.tencent.mm/"
	softType += "</k49><k52>"
	softType += fmt.Sprintf("%d", 0)
	softType += "</k52><k53>"
	softType += fmt.Sprintf("%d", 1)
	softType += "</k53><k57>"
	softType += fmt.Sprintf("%d", 1640)
	//58 apkseccode
	softType += "</k57><k58></k58><k59>"
	softType += fmt.Sprintf("%d", 3)
	softType += "</k59><k60>"
	softType += ""
	//61 true
	softType += "</k60><k61>true</k61><k62>"
	softType += ""
	softType += "</k62><k63>"
	softType += string([]byte(DeviceId))
	softType += "</k63><k64>"
	softType += Info.AndriodUUID(DeviceId)
	softType += "</k64><k65>"
	softType += ""
	softType += "</k65></softtype>"
	return softType
}

func IOSImei(DeviceId string) string {
	return DeviceId
}

func SoftType_iPad(DeviceId string) string {
	uuid1, uuid2 := IOSUuid(DeviceId)
	return "<softtype><k3>13.5</k3><k9>iPad</k9><k10>6</k10><k19>" + uuid1 + "</k19><k20>" + uuid2 + "</k20><k22>(null)</k22><k33>微信</k33><k47>1</k47><k50>1</k50><k51>com.tencent.xin</k51><k54>iPad11,3</k54><k61>2</k61></softtype>"
}

func SoftType_iPhone(DeviceId string) string {
	uuid1, uuid2 := IOSUuid(DeviceId)
	return "<softtype><k3>13.5</k3><k9>iPhone</k9><k10>2</k10><k19>" + uuid1 + "</k19><k20>" + uuid2 + "</k20><k22>中国移动</k22><k33>微信</k33><k47>1</k47><k50>1</k50><k51>com.tencent.xin</k51><k54>iPhone9,1</k54><k61>2</k61></softtype>"
}

func IOSUuid(DeviceId string) (uuid1 string, uuid2 string) {
	Md5DataA := baseutils.MD5ToLower(DeviceId + "SM2020032204320")
	//"b58b9b87-5124-4907-8c92-d8ad59ee0430"
	log.Println("Md5DataA", string(Md5DataA))
	uuid1 = fmt.Sprintf("%s-%s-%s-%s-%s", Md5DataA[0:8], Md5DataA[2:6], Md5DataA[3:7], Md5DataA[1:5], Md5DataA[20:32])
	Md5DataB := baseutils.MD5ToLower(DeviceId + "BM2020032204321")
	uuid2 = fmt.Sprintf("%s-%s-%s-%s-%s", Md5DataB[0:8], Md5DataB[2:6], Md5DataB[3:7], Md5DataB[1:5], Md5DataB[20:32])
	return
}

func IOSMac(DeviceId string) string {
	Md5Data := baseutils.MD5ToLower(DeviceId + "CP2020032204321")
	return fmt.Sprintf("3C:2E:F9:%v:%v:%v", Md5Data[5:7], Md5Data[7:9], Md5Data[10:12])
}

func IOSGetCid(s int) string {
	M := inttobytes(s >> 12)
	return hex.EncodeToString(M)
}

func IOSGetCidUUid(DeviceId, Cid string) string {
	Md5Data := baseutils.MD5ToLower(DeviceId + Cid)
	return fmt.Sprintf("%x-%x-%x-%x-%x", Md5Data[0:8], Md5Data[2:6], Md5Data[3:7], Md5Data[1:5], Md5Data[20:32])
}

func IOSGetCidMd5(DeviceId, Cid string) string {
	Md5Data := baseutils.MD5ToLower(DeviceId + Cid)
	return "A136" + Md5Data[5:]
}

func inttobytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func IOSDeviceNumber(DeviceId string) int64 {
	ssss := []byte(baseutils.MD5ToLower(DeviceId))
	ccc := android.Hex2int(&ssss) >> 8
	ddd := ccc + 60000000000000000
	if ddd > 80000000000000000 {
		ddd = ddd - (80000000000000000 - ddd)
	}
	return int64(ddd)
}

func CreateWIFIinfo(DeviceId, Spare string) (BssID string, Name string) {
	Md5Data := baseutils.MD5ToLower(baseutils.MD5ToLower(DeviceId+Spare) + "WIFI")
	reg := regexp.MustCompile("[a-zA-Z]+")
	M := baseutils.ALLGather(reg.FindAllString(baseutils.MD5ToLower(Md5Data+"A")+baseutils.MD5ToLower(Md5Data+"B"), -1))
	Name = "Chinanet-" + fmt.Sprintf("%v", M[0:5])
	BssID = IOSMac(DeviceId)
	return
}
