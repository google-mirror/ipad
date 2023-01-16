package ccdata

import (
	"encoding/base64"
	"errors"
	"feiyu.com/wx/clientsdk/android"
	"hash/crc32"
	"strconv"
	"time"

	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"

	"github.com/gogo/protobuf/proto"
	"github.com/lunny/log"
)

// saeInfo ClientCheckData加密信息
var saeInfo *wechat.SaeInfo

func init() {
	initSaeInfo()
}

func initSaeInfo() {
	saeData, _ := base64.StdEncoding.DecodeString(android.SaeDat06)
	saeInfo = &wechat.SaeInfo{}
	unmarshalErr := proto.Unmarshal(saeData, saeInfo)
	if unmarshalErr != nil {
		log.Info(unmarshalErr)
	}
}

// GetSaeInfo GetSaeInfo
func GetSaeInfo() *wechat.SaeInfo {
	return saeInfo
}

// CircleShift CircleShift
func CircleShift(data []byte, offset uint32, pos uint32) []byte {
	retData := []byte{}
	retData = append(retData, data[0:]...)
	if pos == 1 {
		retData[offset+0] = data[offset+1]
		retData[offset+1] = data[offset+2]
		retData[offset+2] = data[offset+3]
		retData[offset+3] = data[offset+0]
	}

	if pos == 2 {
		retData[offset+0] = data[offset+2]
		retData[offset+2] = data[offset+0]
		retData[offset+1] = data[offset+3]
		retData[offset+3] = data[offset+1]
	}

	if pos == 3 {
		retData[offset+0] = data[offset+3]
		retData[offset+1] = data[offset+0]
		retData[offset+2] = data[offset+1]
		retData[offset+3] = data[offset+2]
	}

	return retData
}

// ShiftRows ShiftRows
func ShiftRows(data []byte) []byte {
	retData := CircleShift(data, 4, 1)
	retData = CircleShift(retData, 8, 2)
	retData = CircleShift(retData, 12, 3)
	return retData
}

// GetSecTable GetSecTable
func GetSecTable(secTable []byte, encryptRecordData []byte, secTableKey []byte, keyOffset int) []byte {
	// log.Println("keyOffset = ", keyOffset)
	tmpTableKeyOffset := 0
	for index := 0; index < 4; index++ {
		for secIndex := 0; secIndex < 4; secIndex++ {
			tmpOffset := index + secIndex*4
			tmpCount := 0
			recordIndex := 4*index + secIndex
			for threeIndex := 0; threeIndex < 64; threeIndex += 16 {
				tmpValue := 4 * int(encryptRecordData[recordIndex])
				tmpByte := secTableKey[keyOffset+tmpTableKeyOffset+tmpCount+tmpValue]
				secTable[tmpOffset+threeIndex] = tmpByte
				tmpCount++
			}
			tmpTableKeyOffset += 1024
		}
	}
	return secTable
}

// GetSecValue GetSecValue
func GetSecValue(encryptRecordData []byte, secTable []byte, secTableValue []byte, valueOffset int) []byte {
	// log.Println(secTable)
	secTableValueOffset := valueOffset
	for index := 0; index < 4; index++ {
		for secIndex := 0; secIndex < 4; secIndex++ {
			tmpValue := secTable[16*index+4*secIndex+3]
			tmpPtrOffset := 16*index + 4*secIndex + 2
			outBufferOffset := 4*index + secIndex
			for threeIndex := 0; threeIndex < 3; threeIndex++ {
				// 第一部分
				tmpHigh4Value := (secTable[tmpPtrOffset] & 0xF0) | (tmpValue&0xF0)>>4
				tmpValue12 := threeIndex * 0x100
				tmpValue14 := byte(secTableValue[secTableValueOffset+int(tmpHigh4Value&0x7F)+0x200-tmpValue12])
				if tmpHigh4Value&0x80 == 0 {
					tmpValue14 = tmpValue14 & 0x0F
				} else {
					tmpValue14 = tmpValue14 >> 4
				}

				// 第二部分
				tmpLow4Value := byte(tmpValue&0x0F | 16*secTable[tmpPtrOffset])
				tmpValue16 := byte(secTableValue[secTableValueOffset+int(tmpLow4Value&0x7F)+0x280-tmpValue12])
				if tmpLow4Value&0x80 == 0 {
					tmpValue16 = tmpValue16 & 0x0F
				} else {
					tmpValue16 = tmpValue16 >> 4
				}

				// 第三部分
				tmpValue = byte((tmpValue14 << 4) | tmpValue16)
				tmpPtrOffset = tmpPtrOffset - 1
				encryptRecordData[outBufferOffset] = tmpValue
			}
			secTableValueOffset = secTableValueOffset + 0x300
		}
	}

	return encryptRecordData
}

// GetSecValueFinal GetSecValueFinal
func GetSecValueFinal(encryptRecordData []byte, saeTableFinal []byte) []byte {
	for index := 0; index < 4; index++ {
		for secIndex := 0; secIndex < 4; secIndex++ {
			recordIndex := index*4 + secIndex
			recordValue := int(encryptRecordData[recordIndex])
			tmpOffset := index*0x400 + secIndex*0x100
			tmpValue := saeTableFinal[tmpOffset+recordValue]
			encryptRecordData[recordIndex] = tmpValue
		}
	}

	return encryptRecordData
}

// EncodeZipData 加密ClientCheckData压缩后的数据
func EncodeZipData(data []byte, encodeType int) ([]byte, error) {
	retBytes := []byte{}
	tmpEncodeData := data
	saeInfo := GetSaeInfo()
	dataLen := len(data)
	if saeInfo == nil || dataLen <= 0 {
		return retBytes, errors.New("EncodeZipData err: saeInfo == nil || dataLen <= 0")
	}
	if encodeType != 0x3060 && encodeType != 0x4095 {
		return retBytes, errors.New("EncodeZipData err: encodeType != 0x3060 && encodeType != 0x4095")
	}

	// 先按16字节补齐
	lessLen := 16 - dataLen&0xF
	if lessLen < 16 {
		for index := 0; index < lessLen; index++ {
			tmpEncodeData = append(tmpEncodeData, byte(lessLen))
		}
	}

	// IV
	ivData := saeInfo.GetIv()
	lessEncodeLength := len(tmpEncodeData)
	secTable := make([]byte, 64)

	// 每次加密16字节
	count := lessEncodeLength / 16
	for index := 0; index < count; index++ {
		tmpOffset := index * 16
		outEncodeBuffer := make([]byte, 16)
		encryptRecordBuffer := make([]byte, 16)
		// 先跟IV异或
		for secIndex := 0; secIndex < 16; secIndex++ {
			outEncodeBuffer[secIndex] = tmpEncodeData[tmpOffset+secIndex] ^ ivData[secIndex]
		}

		// 第一次换算
		for secIndex := 0; secIndex < 4; secIndex++ {
			for threeIndex := 0; threeIndex < 4; threeIndex++ {
				encryptRecordBuffer[secIndex*4+threeIndex] = outEncodeBuffer[4*threeIndex+secIndex]
			}
		}

		// 行移位
		encryptRecordBuffer = ShiftRows(encryptRecordBuffer)

		// 下一步
		for secIndex := 0; secIndex < 9; secIndex++ {
			// 获取SecTable
			if (encodeType & 0x20) == 0x20 {
				secTable = GetSecTable(secTable, encryptRecordBuffer, saeInfo.GetTableKey(), secIndex*0x4000)
			}

			// 获取SecValue
			if (encodeType & 0x40) == 0x40 {
				encryptRecordBuffer = GetSecValue(encryptRecordBuffer, secTable, saeInfo.GetTableValue(), secIndex*0x3000)
			}
			encryptRecordBuffer = ShiftRows(encryptRecordBuffer)
		}

		// 获取最后的SecValue
		if (encodeType & 0x1000) == 0x1000 {
			encryptRecordBuffer = GetSecValueFinal(encryptRecordBuffer, saeInfo.GetUnknowValue18())
			ivData = outEncodeBuffer
			for secIndex := 0; secIndex < 4; secIndex++ {
				for threeIndex := 0; threeIndex < 4; threeIndex++ {
					outEncodeBuffer[secIndex+4*threeIndex] = encryptRecordBuffer[secIndex*4+threeIndex]
				}
			}
		}

		// 保存第Index次加密后的16字节数据
		retBytes = append(retBytes, outEncodeBuffer[0:]...)
	}
	return retBytes, nil
}

// GetClientCheckDataInfo GetClientCheckDataInfo
func GetClientCheckDataInfo(deviceInfo *baseinfo.DeviceInfo) *baseinfo.ClientCheckDataInfo {
	wechatUUID := baseutils.RandomUUID()
	retInfo := &baseinfo.ClientCheckDataInfo{}
	retInfo.FileSafeAPI = "yes"
	retInfo.DylibSafeAPI = "yes"
	retInfo.OSVersion = deviceInfo.OsTypeNumber
	retInfo.Model = deviceInfo.DeviceName
	retInfo.CoreCount = deviceInfo.CoreCount
	retInfo.VendorID = baseutils.RandomUUID()
	retInfo.ADId = deviceInfo.AdSource
	retInfo.NetType = 1
	retInfo.IsJaiBreak = 0
	retInfo.BundleID = deviceInfo.BundleID
	retInfo.Device = deviceInfo.IphoneVer
	retInfo.DisplayName = "微信"
	retInfo.Version = baseinfo.ClientVersion
	retInfo.PListVersion = baseinfo.PlistVersion
	retInfo.USBState = 2
	retInfo.HasSIMCard = 2
	retInfo.LanguageNum = deviceInfo.Language
	retInfo.LocalCountry = deviceInfo.RealCountry
	retInfo.IsInCalling = 2
	retInfo.WechatUUID = "/var/mobile/Containers/Data/Application/" + baseutils.RandomUUID() + "/Documents"
	retInfo.APPState = 0
	retInfo.IllegalFileList = ""
	retInfo.EncryptStatusOfMachO = 1
	retInfo.Md5OfMachOHeader = baseinfo.Md5OfMachOHeader

	// DirUUID
	dirUUID := baseutils.RandomUUID()
	retInfo.DylibInfoList = make([]*baseinfo.DylibInfo, 0)

	// WeChat
	wechatDylibInfo := &baseinfo.DylibInfo{}
	wechatDylibInfo.S = "/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/WeChat"
	wechatDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, wechatDylibInfo)

	// andromeda
	andromedaDylibInfo := &baseinfo.DylibInfo{}
	andromedaDylibInfo.S = "/private/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/Frameworks/andromeda.framework/andromeda"
	andromedaDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, andromedaDylibInfo)

	// mars
	marsDylibInfo := &baseinfo.DylibInfo{}
	marsDylibInfo.S = "/private/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/Frameworks/mars.framework/mars"
	marsDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, marsDylibInfo)

	// marsbridgenetwork
	marsbridgenetworkDylibInfo := &baseinfo.DylibInfo{}
	marsbridgenetworkDylibInfo.S = "/private/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/Frameworks/marsbridgenetworkDylibInfo.framework/marsbridgenetworkDylibInfo"
	marsbridgenetworkDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, marsbridgenetworkDylibInfo)

	// matrixreport
	matrixreportDylibInfo := &baseinfo.DylibInfo{}
	matrixreportDylibInfo.S = "/private/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/Frameworks/matrixreport.framework/matrixreport"
	matrixreportDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, matrixreportDylibInfo)

	// OpenSSL
	openSSLDylibInfo := &baseinfo.DylibInfo{}
	openSSLDylibInfo.S = "/private/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/Frameworks/OpenSSL.framework/OpenSSL"
	openSSLDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, openSSLDylibInfo)

	// ProtobufLite
	protobufLiteDylibInfo := &baseinfo.DylibInfo{}
	protobufLiteDylibInfo.S = "/private/var/containers/Bundle/Application/" + dirUUID + "/WeChat.app/Frameworks/ProtobufLite.framework/ProtobufLite"
	protobufLiteDylibInfo.U = wechatUUID
	retInfo.DylibInfoList = append(retInfo.DylibInfoList, protobufLiteDylibInfo)

	return retInfo
}

// CreateClientCheckDataXML 创建ClientCheckDataXML
func CreateClientCheckDataXML(deviceInfo *baseinfo.DeviceInfo) string {
	clientCheckDataInfo := GetClientCheckDataInfo(deviceInfo)

	retString := "<clientCheckData>"
	retString = retString + "<fileSafeAPI>" + clientCheckDataInfo.FileSafeAPI + "</fileSafeAPI>"
	retString = retString + "<dylibSafeAPI>" + clientCheckDataInfo.DylibSafeAPI + "</dylibSafeAPI>"
	retString = retString + "<OSVersion>" + clientCheckDataInfo.OSVersion + "</OSVersion>"
	retString = retString + "<table>" + clientCheckDataInfo.Model + "</table>"
	retString = retString + "<coreCount>" + strconv.Itoa(int(clientCheckDataInfo.CoreCount)) + "</coreCount>"
	retString = retString + "<vendorID>" + clientCheckDataInfo.VendorID + "</vendorID>"
	retString = retString + "<adID>" + clientCheckDataInfo.ADId + "</adID>"
	retString = retString + "<netType>" + strconv.Itoa(int(clientCheckDataInfo.NetType)) + "</netType>"
	retString = retString + "<isJailbreak>" + strconv.Itoa(int(clientCheckDataInfo.IsJaiBreak)) + "</isJailbreak>"
	retString = retString + "<bundleID>" + clientCheckDataInfo.BundleID + "</bundleID>"
	retString = retString + "<device>" + clientCheckDataInfo.Device + "</device>"
	retString = retString + "<displayName>" + clientCheckDataInfo.DisplayName + "</displayName>"
	retString = retString + "<version>" + strconv.Itoa(int(clientCheckDataInfo.Version)) + "</version>"
	retString = retString + "<plistVersion>" + strconv.Itoa(int(clientCheckDataInfo.PListVersion)) + "</plistVersion>"
	retString = retString + "<USBState>" + strconv.Itoa(int(clientCheckDataInfo.USBState)) + "</USBState>"
	retString = retString + "<HasSIMCard>" + strconv.Itoa(int(clientCheckDataInfo.HasSIMCard)) + "</HasSIMCard>"
	retString = retString + "<languageNum>" + clientCheckDataInfo.LanguageNum + "</languageNum>"
	retString = retString + "<localeCountry>" + clientCheckDataInfo.LocalCountry + "</localeCountry>"
	retString = retString + "<isInCalling>" + strconv.Itoa(int(clientCheckDataInfo.IsInCalling)) + "</isInCalling>"
	retString = retString + "<weChatUUID>" + clientCheckDataInfo.WechatUUID + "</weChatUUID>"
	retString = retString + "<AppState>" + strconv.Itoa(int(clientCheckDataInfo.APPState)) + "</AppState>"
	retString = retString + "<illegalFileList>" + clientCheckDataInfo.IllegalFileList + "</illegalFileList>"
	retString = retString + "<encryptStatusOfMachO>" + strconv.Itoa(int(clientCheckDataInfo.EncryptStatusOfMachO)) + "</encryptStatusOfMachO>"
	retString = retString + "<md5OfMachOHeader>" + clientCheckDataInfo.Md5OfMachOHeader + "</md5OfMachOHeader>"
	retString = retString + "<dylibInfo>"

	count := len(clientCheckDataInfo.DylibInfoList)
	for index := 0; index < count; index++ {
		dylibInfo := clientCheckDataInfo.DylibInfoList[index]
		retString = retString + "<i>"
		retString = retString + "<s>"
		retString = retString + dylibInfo.S
		retString = retString + "</s>"
		retString = retString + "<u>"
		retString = retString + dylibInfo.U
		retString = retString + "</u>"
		retString = retString + "</i>"
	}

	retString = retString + "</dylibInfo>"
	retString = retString + "</clientCheckData>"
	return retString
}

// CreateClientCheckData 生成ClientCheckData
func CreateClientCheckData(clientCheckDataXML string) ([]byte, error) {
	finalXML := clientCheckDataXML
	crc32Value := crc32.ChecksumIEEE([]byte(finalXML))
	currentTime := int(time.Now().UnixNano() / 1000000000)
	finalXML = finalXML + "<ccdcc>" + strconv.Itoa(int(crc32Value)) + "</ccdcc>"
	finalXML = finalXML + "<ccdts>" + strconv.Itoa(currentTime) + "</ccdts>"

	// 压缩
	zipData := baseutils.CompressByteArray([]byte(finalXML))
	retData, err := EncodeZipData(zipData, 0x3060)
	if err != nil {
		return []byte{}, err
	}

	return retData, nil
}
