package android

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

func AndriodDeviceID(DeviceId string) string {
	Md5Data := strconv.Itoa(BytesToInt([]byte(MD5ToLower(DeviceId + "SM10000011"))))
	return "A0" + Md5Data[0:14]
}

func AndriodImei(DeviceId string) string {
	Md5Data := strconv.Itoa(BytesToInt([]byte(MD5ToLower(DeviceId + "SM1000000"))))
	return "35" + Md5Data[0:13]
}

func AndriodID(DeviceId string) string {
	Md5Data := MD5ToLower(DeviceId + "SM1000001")
	return "06" + Md5Data[0:14]
}

func AndriodSerial(DeviceId string) string {
	Md5Data := MD5ToLower(DeviceId + "SM1000002")
	return "01" + Md5Data[0:14]
}

func AndriodWidevineDeviceID(DeviceId string) string {
	Md5DataA := MD5ToLower(DeviceId + "SM1000003")
	Md5DataB := MD5ToLower(DeviceId + "SM1000004")
	return "657" + Md5DataA[0:29] + Md5DataB
}

func AndriodWidevineProvisionID(DeviceId string) string {
	Md5DataA := MD5ToLower(DeviceId + "SM1000005")
	return "955" + Md5DataA[0:29]
}

func AndriodFSID(DeviceId string) string {
	Md5DataA := strconv.Itoa(BytesToInt([]byte(MD5ToLower(DeviceId + "SM1000012"))))
	Md5DataB := strconv.Itoa(BytesToInt([]byte(MD5ToLower(DeviceId + "SM1000006"))))
	return "37063" + Md5DataA[0:2] + "|37063" + Md5DataA[2:4] + "@" + Md5DataA[4:19] + "|" + MD5ToLower(DeviceId+"SM1000007") + "@" + Md5DataB[0:16] + MD5ToLower(DeviceId+"SM1000008")
}

func AndriodBssid(DeviceId string) string {
	Md5Data := MD5ToLower(DeviceId + "SM1000009")
	A := Md5Data[5:7] + ":"
	B := Md5Data[7:9] + ":"
	C := Md5Data[9:11] + ":"
	D := Md5Data[11:13] + ":"
	E := Md5Data[13:15] + ":"
	F := Md5Data[15:17]
	return A + B + C + D + E + F
}

func AndriodWLanAddress(DeviceId string) string {
	Md5Data := MD5ToLower(DeviceId + "SM1000009")
	B := Md5Data[7:9] + ":"
	C := Md5Data[9:11] + ":"
	D := Md5Data[11:13] + ":"
	E := Md5Data[13:15] + ":"
	F := Md5Data[15:17]
	return "00:" + B + C + D + E + F
}

func AndriodPackageSign(DeviceId string) string {
	Md5Data := MD5ToLower(DeviceId + "SM1000010")
	return "18" + Md5Data[0:30]
}

func MD5ToLower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}
