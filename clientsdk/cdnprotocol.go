package clientsdk

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/lunny/log"
	"strconv"

	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
)

// DecodeCdnData 解析Cdn加密的数据
func DecodeCdnData(data []byte) {
	log.Info("DecodeCdnData ----------------------------------------- in")
	headerLength := uint32(25)
	// 解析头部
	// 总长度
	totalLength := baseutils.BytesToInt32(data[1:5])
	log.Info("totalLength = ", totalLength)
	// flag
	flag := baseutils.BytesToInt32(data[5:7])
	log.Info("flag = ", flag)
	// weixinnum
	winxinnum := baseutils.BytesToInt32SmallEndian(data[7:11])
	log.Info("weixinnum = ", winxinnum)
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	log.Info("bodyLength = ", bodyLength)

	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		if offset+4 >= uint32(len(data)-25) {
			log.Info(fieldName, " = ...(", valueSize, ")")
			break
		}
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		if fieldName == "authkey" ||
			fieldName == "sessionbuf" ||
			fieldName == "thumbdata" ||
			fieldName == "rsavalue" {
			log.Debug(fieldName, " = ", len(value))
			// baseutils.PrintBytesHex([]byte(value), fieldName)
		} else if fieldName == "filedata" {
			log.Info("fileDataLen = ", len(value))
		} else {
			log.Info(fieldName, " = ", value)
		}

		// if fieldName == "sessionbuf" ||
		// 	fieldName == "skeybuf" {
		// 	userInfo := NewUserInfo()
		// 	// userInfo.SessionKey = []byte{0x65, 0x2b, 0x79, 0x25, 0x68, 0x24, 0x2e, 0x40, 0x4c, 0x52, 0x57, 0x2a, 0x2a, 0x68, 0x6a, 0x53}
		// 	userInfo.SessionKey = []byte{0x7b, 0x2e, 0x21, 0x2a, 0x68, 0x62, 0x3f, 0x63, 0x6b, 0x29, 0x61, 0x70, 0x6a, 0x77, 0x55, 0x4c}
		// 	response := &wechat.CDNUploadMsgImgPrepareResponse{}
		// 	ParseResponseData(userInfo, []byte(value), response)
		// 	ShowObjectValue(response)
		// }
	}
	log.Info("DecodeCdnData ----------------------------------------- out\n")
}

// GetErrStringByRetCode 获取错误信息
func GetErrStringByRetCode(retCode uint32) string {
	if retCode == 4289864094 {
		return "大小超过限制"
	}
	return strconv.Itoa(int(retCode))
}

// CreateID 加密生成ID
func CreateID(data []byte) uint32 {
	length := len(data)
	if length < 1 {
		return 0
	}

	tmpTotalLength := uint32(length)
	if length>>2 > 0 {
		tmpLen := length>>2 + 1

		index := 0
		for tmpLen > 1 {
			value0 := uint32(data[index])
			value1 := uint32(data[index+1])
			value2 := uint32(data[index+2])
			value3 := uint32(data[index+3])

			v5 := (value0 | (value1 << 8)) + tmpTotalLength
			v6 := value2 | (value3 << 8)
			tmpValue := (v5 ^ (v5 << 16) ^ (v6 << 11))
			tmpTotalLength = tmpValue + (tmpValue >> 11)
			index = index + 4
			tmpLen = tmpLen - 1
		}
	}

	caseValue := length & 3
	if caseValue == 1 {
		tmpValue0 := uint32(data[0])
		tmpValue := tmpTotalLength + tmpValue0
		tmpValue2 := tmpValue ^ (tmpValue << 10)
		tmpTotalLength = tmpValue2 + (tmpValue2 >> 1)
	}

	if caseValue == 2 {
		value0 := uint32(data[0])
		value1 := uint32(data[1])
		tmpValue0 := value0 | (value1 << 8)
		tmpValue := tmpTotalLength + tmpValue0
		tmpValue2 := tmpValue ^ (tmpValue << 11)
		tmpTotalLength = tmpValue2 + (tmpValue2 >> 17)
	}

	if caseValue == 3 {
		value0 := uint32(data[0])
		value1 := uint32(data[1])
		value2 := uint32(data[2])
		tmpValue0 := (value0 | (value1 << 8)) + tmpTotalLength
		tmpValue1 := tmpValue0 ^ (value2 << 18)
		tmpValue2 := tmpValue1 ^ (tmpValue0 << 16)
		tmpTotalLength = tmpValue2 + (tmpValue2 >> 11)
	}

	tmpValue0 := tmpTotalLength ^ (8 * tmpTotalLength)
	tmpValue1 := tmpValue0 + (tmpValue0 >> 5)
	tmpValue2 := tmpValue1 ^ (16 * tmpValue1)
	tmpValue3 := tmpValue2 + (tmpValue2 >> 17)
	tmpValue4 := tmpValue3 ^ (tmpValue3 << 25)
	tmpValue5 := tmpValue4 + (tmpValue4 >> 6)

	return tmpValue5
}

// ParseCdnResponseDataLength 解析cdn响应数据的总长度
func ParseCdnResponseDataLength(data []byte) uint32 {
	totalLength := baseutils.BytesToInt32(data[1:5])
	return totalLength
}

// PackCdnRequestElementUint32Pointer 请求字段*uint32元素
func PackCdnRequestElementUint32Pointer(fieldName string, value *uint32) []byte {
	retData := make([]byte, 0)

	// 写入字段名长度
	fieldNameLength := uint32(len(fieldName))
	fieldNameLengthData := baseutils.Int32ToBytes(fieldNameLength)
	retData = append(retData, fieldNameLengthData[0:]...)

	// 写入字段名称
	retData = append(retData, ([]byte(fieldName))[0:]...)

	// 字段值转成string
	valueString := string("")
	if value != nil {
		valueString = strconv.Itoa(int(*value))
	}
	// 写入字段值字符串 长度
	valueStringLength := uint32(len(valueString))
	valueStringLengththData := baseutils.Int32ToBytes(valueStringLength)
	retData = append(retData, valueStringLengththData[0:]...)
	// 写入字段值字符串
	retData = append(retData, ([]byte(valueString))[0:]...)

	return retData
}

// PackCdnRequestElementUint32 请求字段uint32元素
func PackCdnRequestElementUint32(fieldName string, value uint32) []byte {
	retData := make([]byte, 0)

	// 写入字段名长度
	fieldNameLength := uint32(len(fieldName))
	fieldNameLengthData := baseutils.Int32ToBytes(fieldNameLength)
	retData = append(retData, fieldNameLengthData[0:]...)

	// 写入字段名称
	retData = append(retData, ([]byte(fieldName))[0:]...)

	// 字段值转成string
	valueString := strconv.Itoa(int(value))
	// 写入字段值字符串 长度
	valueStringLength := uint32(len(valueString))
	valueStringLengththData := baseutils.Int32ToBytes(valueStringLength)
	retData = append(retData, valueStringLengththData[0:]...)
	// 写入字段值字符串
	retData = append(retData, ([]byte(valueString))[0:]...)

	return retData
}

// PackCdnRequestElementData 请求字段[]byte元素
func PackCdnRequestElementData(fieldName string, value []byte) []byte {
	retData := make([]byte, 0)

	// 写入字段名长度
	fieldNameLength := uint32(len(fieldName))
	fieldNameLengthData := baseutils.Int32ToBytes(fieldNameLength)
	retData = append(retData, fieldNameLengthData[0:]...)

	// 写入字段名称
	retData = append(retData, ([]byte(fieldName))[0:]...)

	// 写入字段值字符串 长度
	valueLength := uint32(len(value))
	valueLengththData := baseutils.Int32ToBytes(valueLength)
	retData = append(retData, valueLengththData[0:]...)
	// 写入字段值字符串
	retData = append(retData, value[0:]...)

	return retData
}

// PackCdnImageDownloadRequest 对Cdn图片下载请求进行打包
func PackCdnImageDownloadRequest(request *baseinfo.CdnImageDownloadRequest) []byte {
	// 打包请求包体
	bodyData := make([]byte, 0)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ver", request.Ver)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("weixinnum", request.WeiXinNum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("seq", request.Seq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientversion", request.ClientVersion)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientostype", []byte(request.ClientOsType))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("authkey", request.AuthKey)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nettype", request.NetType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("acceptdupack", request.AcceptDupack)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rsaver", request.RsaVer)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rsavalue", request.RsaValue)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filetype", request.FileType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("wxchattype", request.WxChatType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("fileid", []byte(request.FileID))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastretcode", request.LastRetCode)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ipseq", request.IPSeq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("cli-quic-flag", request.CliQuicFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32Pointer("wxmsgflag", request.WxMsgFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("wxautostart", request.WxAutoStart)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("downpicformat", request.DownPicFormat)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("offset", request.Offset)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("largesvideo", request.LargesVideo)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("sourceflag", request.SourceFlag)[0:]...)

	retData := make([]byte, 0)
	// 包头
	retData = append(retData, 0xab)
	// 总长度
	totalLength := uint32(25 + len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(totalLength)[0:]...)
	// Flag标志
	flag := uint16(20000)
	retData = append(retData, baseutils.Int16ToBytesBigEndian(flag)[0:]...)
	// weixinnum
	retData = append(retData, baseutils.Int32ToBytesLittleEndian(request.WeiXinNum)[0:]...)
	// 固定为空的数据
	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	retData = append(retData, zeroBytes[0:]...)
	// bodyLength
	bodyLength := uint32(len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(bodyLength)[0:]...)
	// BodyData
	retData = append(retData, bodyData[0:]...)

	return retData
}

// DecodeImageDownloadResponse 解析下载图片响应数据
func DecodeImageDownloadResponse(data []byte) (*baseinfo.CdnDownloadResponse, error) {
	if len(data) < 25 {
		return nil, errors.New("DecodeImageDownloadResponse err: len(data) < 25")
	}
	response := &baseinfo.CdnDownloadResponse{}

	// 头的总长度 固定为25个字节
	headerLength := uint32(25)
	// 解析头部
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		// Ver
		if fieldName == "ver" {
			ver, _ := strconv.Atoi(value)
			response.Ver = uint32(ver)
		}

		// Seq
		if fieldName == "seq" {
			seq, _ := strconv.Atoi(value)
			response.Seq = uint32(seq)
		}

		// VideoFormat
		if fieldName == "videoformat" {
			videoformat, _ := strconv.Atoi(value)
			response.VideoFormat = uint32(videoformat)
		}

		// RspPicFormat
		if fieldName == "rsppicformat" {
			rsppicformat, _ := strconv.Atoi(value)
			response.RspPicFormat = uint32(rsppicformat)
		}

		// RangeStart
		if fieldName == "rangestart" {
			rangestart, _ := strconv.Atoi(value)
			response.RangeStart = uint32(rangestart)
		}

		// RangeEnd
		if fieldName == "rangeend" {
			rangeend, _ := strconv.Atoi(value)
			response.RangeEnd = uint32(rangeend)
		}

		// TotalSize
		if fieldName == "totalsize" {
			totalsize, _ := strconv.Atoi(value)
			response.TotalSize = uint32(totalsize)
		}

		// SrcSize
		if fieldName == "srcsize" {
			srcsize, _ := strconv.Atoi(value)
			response.SrcSize = uint32(srcsize)
		}

		// RetCode
		if fieldName == "retcode" {
			retcode, _ := strconv.Atoi(value)
			response.RetCode = uint32(retcode)
		}

		// SubStituteFType
		if fieldName == "substituteftype" {
			substituteftype, _ := strconv.Atoi(value)
			response.SubStituteFType = uint32(substituteftype)
		}

		// RetrySec
		if fieldName == "retrysec" {
			retrysec, _ := strconv.Atoi(value)
			response.RetrySec = uint32(retrysec)
		}

		// IsRetry
		if fieldName == "isretry" {
			isretry, _ := strconv.Atoi(value)
			response.IsRetry = uint32(isretry)
		}

		// IsOverLoad
		if fieldName == "isoverload" {
			isoverload, _ := strconv.Atoi(value)
			response.IsOverLoad = uint32(isoverload)
		}

		// IsGetCdn
		if fieldName == "isgetcdn" {
			isgetcdn, _ := strconv.Atoi(value)
			response.IsGetCdn = uint32(isgetcdn)
		}

		// XClientIP
		if fieldName == "x-ClientIp" {
			response.XClientIP = value
		}

		// FileData
		if fieldName == "filedata" {
			response.FileData = []byte(value)
		}
	}

	return response, nil
}

// PackCdnImageUploadRequest 打包上传图片数据
func PackCdnImageUploadRequest(request *baseinfo.CdnImageUploadRequest) []byte {
	bodyData := make([]byte, 0)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ver", request.Ver)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("weixinnum", request.WeiXinNum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("seq", request.Seq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientversion", request.ClientVersion)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientostype", []byte(request.ClientOsType))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("authkey", request.AuthKey[0:])...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nettype", request.NetType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("acceptdupack", request.AcceptDupack)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("safeproto", request.SafeProto)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filetype", request.FileType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("wxchattype", request.WxChatType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastretcode", request.LastRetCode)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ipseq", request.IPSeq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("cli-quic-flag", request.CliQuicFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("hasthumb", request.HasThumb)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("touser", []byte(request.ToUser))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("compresstype", request.CompressType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nocheckaeskey", request.NoCheckAesKey)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("enablehit", request.EnableHit)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("existancecheck", request.ExistAnceCheck)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("apptype", request.AppType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filekey", []byte(request.FileKey))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("totalsize", request.TotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawtotalsize", request.RawTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("localname", []byte(request.LocalName))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("sessionbuf", request.SessionBuf)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("offset", request.Offset)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("thumbtotalsize", request.ThumbTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawthumbsize", request.RawThumbSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawthumbmd5", []byte(request.RawThumbMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("encthumbcrc", request.EncThumbCRC)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("thumbdata", request.ThumbData)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("largesvideo", request.LargesVideo)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("sourceflag", request.SourceFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("advideoflag", request.AdVideoFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filemd5", []byte(request.FileMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawfilemd5", []byte(request.RawFileMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("datachecksum", request.DataCheckSum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filecrc", request.FileCRC)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("setofpicformat", []byte(request.SetOfPicFormat))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filedata", request.FileData)[0:]...)

	retData := make([]byte, 0)
	// 包头
	retData = append(retData, 0xab)
	// 总长度
	totalLength := uint32(25 + len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(totalLength)[0:]...)
	// Flag标志
	flag := uint16(10000)
	retData = append(retData, baseutils.Int16ToBytesBigEndian(flag)[0:]...)
	// weixinnum
	retData = append(retData, baseutils.Int32ToBytesLittleEndian(request.WeiXinNum)[0:]...)
	// 固定为空的数据
	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	retData = append(retData, zeroBytes[0:]...)
	// bodyLength
	bodyLength := uint32(len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(bodyLength)[0:]...)
	// BodyData
	retData = append(retData, bodyData[0:]...)

	return retData
}

// PackCdnVideoUploadRequest 视频上传请求
func PackCdnVideoUploadRequest(request *baseinfo.CdnVideoUploadRequest) []byte {
	bodyData := make([]byte, 0)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ver", request.Ver)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("weixinnum", request.WeiXinNum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("seq", request.Seq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientversion", request.ClientVersion)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientostype", []byte(request.ClientOSType))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("authkey", request.AutoKey[0:])...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nettype", request.NetType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("acceptdupack", request.AcceptDuPack)[0:]...)
	/*	bodyData = append(bodyData, PackCdnRequestElementUint32("rsaver", request.ClientRsaVer)[0:]...)
		bodyData = append(bodyData, PackCdnRequestElementData("rsavalue", request.ClientRsaVal)[0:]...)*/
	bodyData = append(bodyData, PackCdnRequestElementUint32("safeproto", request.SafeProto)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filetype", request.FileType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("wxchattype", request.WeChatType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastretcode", request.LastRetCode)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ipseq", request.IpSeq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("hasthumb", request.HastHumb)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("touser", []byte(request.ToUSerName))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("compresstype", request.CompressType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nocheckaeskey", request.NoCheckAesKey)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("enablehit", request.EnaBleHit)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("existancecheck", request.ExistAnceCheck)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("apptype", request.AppType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filekey", []byte(request.FileKey))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("totalsize", request.TotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawtotalsize", request.RawTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("localname", []byte(request.LocalName))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("offset", request.Offset)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("thumbtotalsize", request.ThumbTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawthumbsize", request.RawThumbSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawthumbmd5", []byte(request.RawThumbMd5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("encthumbcrc", request.EncThumbCrc)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("thumbdata", request.ThumbData)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("largesvideo", request.LargesVideo)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("sourceflag", request.SourceFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("advideoflag", request.AdVideoFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("mp4identify", []byte(request.Mp4identify))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("droprateflag", request.DropRateFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientrsaver", request.ClientRsaVer)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientrsaval", request.ClientRsaVal)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filemd5", []byte(request.FileMd5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawfilemd5", []byte(request.RawFileMd5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("datachecksum", request.DataCheckSum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filecrc", request.FileCrc)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filedata", request.FileData)[0:]...)
	retData := make([]byte, 0)
	// 包头
	retData = append(retData, 0xab)
	// 总长度
	totalLength := uint32(25 + len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(totalLength)[0:]...)
	/*retData = append(retData, []byte{0x75, 0x30, 0x10, 0xa4, 0x65, 0x9a,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)*/
	// Flag标志
	flag := uint16(30000)
	retData = append(retData, baseutils.Int16ToBytesBigEndian(flag)[0:]...)
	// weixinnum
	retData = append(retData, baseutils.Int32ToBytesLittleEndian(request.WeiXinNum)[0:]...)
	// 固定为空的数据
	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	retData = append(retData, zeroBytes[0:]...)
	// bodyLength
	bodyLength := uint32(len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(bodyLength)[0:]...)
	log.Debug(hex.EncodeToString(retData))
	// BodyData
	retData = append(retData, bodyData[0:]...)

	return retData
}

// DecodeImageUploadResponse 解析上传图片响应
func DecodeImageUploadResponse(data []byte) (*baseinfo.CdnImageUploadResponse, error) {
	if len(data) < 25 {
		return nil, errors.New("DecodeImageUploadResponse err: len(data) < 25")
	}

	response := &baseinfo.CdnImageUploadResponse{}

	// 头的总长度 固定为25个字节
	headerLength := uint32(25)
	// 解析头部
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		// Ver
		if fieldName == "ver" {
			ver, _ := strconv.Atoi(value)
			response.Ver = uint32(ver)
		}

		// Seq
		if fieldName == "seq" {
			seq, _ := strconv.Atoi(value)
			response.Seq = uint32(seq)
		}

		// RetCode
		if fieldName == "retcode" {
			retcode, _ := strconv.Atoi(value)
			response.RetCode = uint32(retcode)
		}

		// FileKey
		if fieldName == "filekey" {
			response.FileKey = value
		}

		// RecvLen
		if fieldName == "recvlen" {
			recvlen, _ := strconv.Atoi(value)
			response.RecvLen = uint32(recvlen)
		}

		// SKeyResp
		if fieldName == "skeyresp" {
			skeyresp, _ := strconv.Atoi(value)
			response.SKeyResp = uint32(skeyresp)
		}

		// SKeyBuf
		if fieldName == "skeybuf" {
			response.SKeyBuf = []byte(value)
		}

		// FileID
		if fieldName == "fileid" {
			response.FileID = value
		}

		// ExistFlag
		if fieldName == "existflag" {
			existflag, _ := strconv.Atoi(value)
			response.ExistFlag = uint32(existflag)
		}

		// hittype
		if fieldName == "hittype" {
			hittype, _ := strconv.Atoi(value)
			response.HitType = uint32(hittype)
		}

		// RetrySec
		if fieldName == "retrysec" {
			retrysec, _ := strconv.Atoi(value)
			response.RetrySec = uint32(retrysec)
		}

		// IsRetry
		if fieldName == "isretry" {
			isretry, _ := strconv.Atoi(value)
			response.IsRetry = uint32(isretry)
		}

		// IsOverLoad
		if fieldName == "isoverload" {
			isoverload, _ := strconv.Atoi(value)
			response.IsOverLoad = uint32(isoverload)
		}

		// IsGetCdn
		if fieldName == "isgetcdn" {
			isgetcdn, _ := strconv.Atoi(value)
			response.IsGetCDN = uint32(isgetcdn)
		}

		// XClientIP
		if fieldName == "x-ClientIp" {
			response.XClientIP = value
		}
	}

	return response, nil
}

// PackCdnSnsImageUploadRequest 打包上传朋友圈图片数据
func PackCdnSnsImageUploadRequest(request *baseinfo.CdnSnsImageUploadRequest) []byte {
	bodyData := make([]byte, 0)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ver", request.Ver)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("weixinnum", request.WeiXinNum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("seq", request.Seq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientversion", request.ClientVersion)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientostype", []byte(request.ClientOsType))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("authkey", request.AuthKey[0:])...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nettype", request.NetType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("acceptdupack", request.AcceptDupack)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rsaver", request.RsaVer)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rsavalue", request.RsaValue)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filetype", request.FileType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("wxchattype", request.WxChatType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastretcode", request.LastRetCode)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ipseq", request.IPSeq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("cli-quic-flag", request.CliQuicFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("hasthumb", request.HasThumb)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("touser", []byte(request.ToUser))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("compresstype", request.CompressType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nocheckaeskey", request.NoCheckAesKey)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("enablehit", request.EnableHit)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("existancecheck", request.ExistAnceCheck)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("apptype", request.AppType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filekey", []byte(request.FileKey))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("totalsize", request.TotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawtotalsize", request.RawTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("localname", []byte(request.LocalName))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("offset", request.Offset)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("thumbtotalsize", request.ThumbTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawthumbsize", request.RawThumbSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawthumbmd5", []byte(request.RawThumbMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("thumbcrc", request.ThumbCRC)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("largesvideo", request.LargesVideo)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("sourceflag", request.SourceFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("advideoflag", request.AdVideoFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filemd5", []byte(request.FileMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawfilemd5", []byte(request.RawFileMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("datachecksum", request.DataCheckSum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filecrc", request.FileCRC)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filedata", request.FileData)[0:]...)

	retData := make([]byte, 0)
	// 包头
	retData = append(retData, 0xab)
	// 总长度
	totalLength := uint32(25 + len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(totalLength)[0:]...)
	// Flag标志
	flag := uint16(10002)
	retData = append(retData, baseutils.Int16ToBytesBigEndian(flag)[0:]...)
	// weixinnum
	retData = append(retData, baseutils.Int32ToBytesLittleEndian(request.WeiXinNum)[0:]...)
	// 固定为空的数据
	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	retData = append(retData, zeroBytes[0:]...)
	// bodyLength
	bodyLength := uint32(len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(bodyLength)[0:]...)
	// BodyData
	retData = append(retData, bodyData[0:]...)

	return retData
}

// PackCdnSnsVideoUploadRequest 打包上传朋友圈视频请求
func PackCdnSnsVideoUploadRequest(request *baseinfo.CdnSnsVideoUploadRequest) []byte {
	bodyData := make([]byte, 0)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ver", request.Ver)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("weixinnum", request.WeiXinNum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("seq", request.Seq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientversion", request.ClientVersion)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientostype", []byte(request.ClientOsType))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("authkey", request.AuthKey[0:])...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nettype", request.NetType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("acceptdupack", request.AcceptDupack)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rsaver", request.RsaVer)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rsavalue", request.RsaValue)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filetype", request.FileType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("wxchattype", request.WxChatType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastretcode", request.LastRetCode)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ipseq", request.IPSeq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("cli-quic-flag", request.CliQuicFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("isstorevideo", request.IsStoreVideo)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("hasthumb", request.HasThumb)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nocheckaeskey", request.NoCheckAesKey)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("enablehit", request.EnableHit)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("existancecheck", request.ExistAnceCheck)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("apptype", request.AppType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filekey", []byte(request.FileKey))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("totalsize", request.TotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawtotalsize", request.RawTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("localname", []byte(request.LocalName))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("offset", request.Offset)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("thumbtotalsize", request.ThumbTotalSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rawthumbsize", request.RawThumbSize)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawthumbmd5", []byte(request.RawThumbMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("thumbcrc", request.ThumbCRC)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("thumbdata", request.ThumbData)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("largesvideo", request.LargesVideo)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("sourceflag", request.SourceFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("advideoflag", request.AdVideoFlag)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("mp4identify", []byte(request.Mp4Identify))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filemd5", []byte(request.FileMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("rawfilemd5", []byte(request.RawFileMD5))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("datachecksum", request.DataCheckSum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("filecrc", request.FileCRC)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("filedata", request.FileData)[0:]...)

	retData := make([]byte, 0)
	// 包头
	retData = append(retData, 0xab)
	// 总长度
	totalLength := uint32(25 + len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(totalLength)[0:]...)
	// Flag标志
	flag := uint16(10002)
	retData = append(retData, baseutils.Int16ToBytesBigEndian(flag)[0:]...)
	// weixinnum
	retData = append(retData, baseutils.Int32ToBytesLittleEndian(request.WeiXinNum)[0:]...)
	// 固定为空的数据
	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	retData = append(retData, zeroBytes[0:]...)
	// bodyLength
	bodyLength := uint32(len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(bodyLength)[0:]...)
	// BodyData
	retData = append(retData, bodyData[0:]...)

	return retData
}

// DecodeSnsVideoUploadResponse 解析上传朋友圈视频响应
func DecodeSnsVideoUploadResponse(data []byte) (*baseinfo.CdnSnsVideoUploadResponse, error) {
	if len(data) < 25 {
		return nil, errors.New("DecodeSnsVideoUploadResponse err: len(data) < 25")
	}

	response := &baseinfo.CdnSnsVideoUploadResponse{}

	// 头的总长度 固定为25个字节
	headerLength := uint32(25)
	// 解析头部
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		// Ver
		if fieldName == "ver" {
			ver, _ := strconv.Atoi(value)
			response.Ver = uint32(ver)
		}

		// Seq
		if fieldName == "seq" {
			seq, _ := strconv.Atoi(value)
			response.Seq = uint32(seq)
		}

		// RetCode
		if fieldName == "retcode" {
			retcode, _ := strconv.Atoi(value)
			log.Debug(retcode)
			response.RetCode = uint32(retcode)
		}

		// FileKey
		if fieldName == "filekey" {
			response.FileKey = value
		}

		// FileURL
		if fieldName == "fileurl" {
			response.FileURL = value
		}

		// ThumbURL
		if fieldName == "thumburl" {
			response.ThumbURL = value
		}

		// FileID
		if fieldName == "fileid" {
			response.FileID = value
		}

		// RecvLen
		if fieldName == "recvlen" {
			recvlen, _ := strconv.Atoi(value)
			response.RecvLen = uint32(recvlen)
		}

		// RetrySec
		if fieldName == "retrysec" {
			retrysec, _ := strconv.Atoi(value)
			response.RetrySec = uint32(retrysec)
		}

		// IsRetry
		if fieldName == "isretry" {
			isretry, _ := strconv.Atoi(value)
			response.IsRetry = uint32(isretry)
		}

		// IsOverLoad
		if fieldName == "isoverload" {
			isoverload, _ := strconv.Atoi(value)
			response.IsOverLoad = uint32(isoverload)
		}

		// IsGetCdn
		if fieldName == "isgetcdn" {
			isgetcdn, _ := strconv.Atoi(value)
			response.IsGetCDN = uint32(isgetcdn)
		}

		// XClientIP
		if fieldName == "x-ClientIp" {
			response.XClientIP = value
		}

	}
	dataJson, _ := json.Marshal(response)
	log.Info(string(dataJson))
	return response, nil
}

// DecodeSnsImageUploadResponse 解析朋友圈上传图片响应
func DecodeSnsImageUploadResponse(data []byte) (*baseinfo.CdnSnsImageUploadResponse, error) {
	if len(data) < 25 {
		return nil, errors.New("DecodeSnsImageUploadResponse err: len(data) < 25")
	}

	response := &baseinfo.CdnSnsImageUploadResponse{}
	// 头的总长度 固定为25个字节
	headerLength := uint32(25)
	// 解析头部
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		// Ver
		if fieldName == "ver" {
			ver, _ := strconv.Atoi(value)
			response.Ver = uint32(ver)
		}

		// Seq
		if fieldName == "seq" {
			seq, _ := strconv.Atoi(value)
			response.Seq = uint32(seq)
		}

		// RetCode
		if fieldName == "retcode" {
			retcode, _ := strconv.Atoi(value)
			response.RetCode = uint32(retcode)
		}

		// FileKey
		if fieldName == "filekey" {
			response.FileKey = value
		}

		// RecvLen
		if fieldName == "recvlen" {
			recvlen, _ := strconv.Atoi(value)
			response.RecvLen = uint32(recvlen)
		}

		// FileURL
		if fieldName == "fileurl" {
			response.FileURL = value
		}

		// ThumbURL
		if fieldName == "thumburl" {
			response.ThumbURL = value
		}

		// EnableQuic
		if fieldName == "enablequic" {
			enablequic, _ := strconv.Atoi(value)
			response.EnableQuic = uint32(enablequic)
		}

		// RetrySec
		if fieldName == "retrysec" {
			retrysec, _ := strconv.Atoi(value)
			response.RetrySec = uint32(retrysec)
		}

		// IsRetry
		if fieldName == "isretry" {
			isretry, _ := strconv.Atoi(value)
			response.IsRetry = uint32(isretry)
		}

		// IsOverLoad
		if fieldName == "isoverload" {
			isoverload, _ := strconv.Atoi(value)
			response.IsOverLoad = uint32(isoverload)
		}

		// IsGetCdn
		if fieldName == "isgetcdn" {
			isgetcdn, _ := strconv.Atoi(value)
			response.IsGetCDN = uint32(isgetcdn)
		}

		// XClientIP
		if fieldName == "x-ClientIp" {
			response.XClientIP = value
		}
	}

	return response, nil
}

// PackCdnSnsVideoDownloadRequest 对Cdn朋友圈视频下载请求打包
func PackCdnSnsVideoDownloadRequest(request *baseinfo.CdnSnsVideoDownloadRequest) []byte {
	// 打包请求包体
	bodyData := make([]byte, 0)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ver", request.Ver)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("weixinnum", request.WeiXinNum)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("seq", request.Seq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("clientversion", request.ClientVersion)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("clientostype", []byte(request.ClientOsType))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("authkey", request.AuthKey)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("nettype", request.NetType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("acceptdupack", request.AcceptDupack)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("signal", []byte(request.Signal))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("scene", []byte(request.Scene))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("url", []byte(request.URL))[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rangestart", request.RangeStart)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("rangeend", request.RangeEnd)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastretcode", request.LastRetCode)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("ipseq", request.IPSeq)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("redirect_type", request.RedirectType)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("lastvideoformat", request.LastVideoFormat)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementUint32("videoformat", request.VideoFormat)[0:]...)
	bodyData = append(bodyData, PackCdnRequestElementData("X-snsvideoflag", []byte(request.XSnsVideoFlag))[0:]...)

	retData := make([]byte, 0)
	// 包头
	retData = append(retData, 0xab)
	// 总长度
	totalLength := uint32(25 + len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(totalLength)[0:]...)
	// Flag标志
	flag := uint16(10005)
	retData = append(retData, baseutils.Int16ToBytesBigEndian(flag)[0:]...)
	// weixinnum
	retData = append(retData, baseutils.Int32ToBytesLittleEndian(request.WeiXinNum)[0:]...)
	// 固定为空的数据
	zeroBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	retData = append(retData, zeroBytes[0:]...)
	// bodyLength
	bodyLength := uint32(len(bodyData))
	retData = append(retData, baseutils.Int32ToBytes(bodyLength)[0:]...)
	// BodyData
	retData = append(retData, bodyData[0:]...)

	return retData
}

// DecodeSnsVideoDownloadResponse 解析朋友圈视频下载响应
func DecodeSnsVideoDownloadResponse(data []byte) (*baseinfo.CdnSnsVideoDownloadResponse, error) {
	if len(data) < 25 {
		return nil, errors.New("DecodeSnsVideoDownloadResponse err: len(data) < 25")
	}

	response := &baseinfo.CdnSnsVideoDownloadResponse{}
	// 头的总长度 固定为25个字节
	headerLength := uint32(25)

	// 解析头部
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		// Ver
		if fieldName == "ver" {
			ver, _ := strconv.Atoi(value)
			response.Ver = uint32(ver)
		}

		// Seq
		if fieldName == "seq" {
			seq, _ := strconv.Atoi(value)
			response.Seq = uint32(seq)
		}

		// RetCode
		if fieldName == "retcode" {
			retcode, _ := strconv.Atoi(value)
			response.RetCode = uint32(retcode)
		}

		// RangeStart
		if fieldName == "rangestart" {
			rangeStart, _ := strconv.Atoi(value)
			response.RangeStart = uint32(rangeStart)
		}

		// RangeEnd
		if fieldName == "rangeend" {
			rangeEnd, _ := strconv.Atoi(value)
			response.RangeEnd = uint32(rangeEnd)
		}

		// TotalSize
		if fieldName == "totalsize" {
			totalSize, _ := strconv.Atoi(value)
			response.TotalSize = uint32(totalSize)
		}

		// EnableQuic
		if fieldName == "enablequic" {
			enablequic, _ := strconv.Atoi(value)
			response.EnableQuic = uint32(enablequic)
		}

		// IsRetry
		if fieldName == "isretry" {
			isretry, _ := strconv.Atoi(value)
			response.IsRetry = uint32(isretry)
		}

		// IsOverLoad
		if fieldName == "isoverload" {
			isoverload, _ := strconv.Atoi(value)
			response.IsOverLoad = uint32(isoverload)
		}

		// IsGetCdn
		if fieldName == "isgetcdn" {
			isgetcdn, _ := strconv.Atoi(value)
			response.IsGetCdn = uint32(isgetcdn)
		}

		// XClientIP
		if fieldName == "x-ClientIp" {
			response.XClientIP = value
		}

		// XEncFlag
		if fieldName == "X-encflag" {
			encFlag, _ := strconv.Atoi(value)
			response.XEncFlag = uint32(encFlag)
		}

		// XEncLen
		if fieldName == "X-enclen" {
			encLen, _ := strconv.Atoi(value)
			response.XEncLen = uint32(encLen)
		}

		// XSnsVideoFlag
		if fieldName == "X-snsvideoflag" {
			response.XSnsVideoFlag = value
		}

		// XSnsVideoTicket
		if fieldName == "X-snsvideoticket" {
			response.XSnsVideoTicket = value
		}

		// FileData
		if fieldName == "filedata" {
			response.FileData = []byte(value)
		}
	}

	return response, nil
}

// DecodeVideoUploadResponse 解析上传视频响应
func DecodeVideoUploadResponse(data []byte) (*baseinfo.CdnMsgVideoUploadResponse, error) {
	if len(data) < 25 {
		return nil, errors.New("DecodeSnsVideoUploadResponse err: len(data) < 25")
	}

	response := &baseinfo.CdnMsgVideoUploadResponse{}

	// 头的总长度 固定为25个字节
	headerLength := uint32(25)
	// 解析头部
	// bodyLength
	bodyLength := baseutils.BytesToInt32(data[21:25])
	// 解析包体
	body := data[headerLength:]
	offset := uint32(0)
	for offset < bodyLength {
		fieldNameSize := baseutils.BytesToInt32(body[offset : offset+4])
		fieldName := string(body[offset+4 : offset+4+fieldNameSize])
		offset = offset + fieldNameSize + 4

		valueSize := baseutils.BytesToInt32(body[offset : offset+4])
		value := string(body[offset+4 : offset+4+valueSize])
		offset = offset + valueSize + 4

		// Ver
		if fieldName == "ver" {
			ver, _ := strconv.Atoi(value)
			response.Ver = uint32(ver)
		}

		// Seq
		if fieldName == "seq" {
			seq, _ := strconv.Atoi(value)
			response.Seq = uint32(seq)
		}

		// RetCode
		if fieldName == "retcode" {
			retcode, _ := strconv.Atoi(value)
			log.Debug(retcode)
			response.RetCode = uint32(retcode)
		}

		// FileKey
		if fieldName == "filekey" {
			response.FileKey = value
		}

		// FileURL
		if fieldName == "fileurl" {
			response.FileURL = value
		}

		// ThumbURL
		if fieldName == "thumburl" {
			response.ThumbURL = value
		}

		// FileID
		if fieldName == "fileid" {
			response.FileID = value
		}

		// RecvLen
		if fieldName == "recvlen" {
			recvlen, _ := strconv.Atoi(value)
			response.RecvLen = uint32(recvlen)
		}

		// RetrySec
		if fieldName == "retrysec" {
			retrysec, _ := strconv.Atoi(value)
			response.RetrySec = uint32(retrysec)
		}

		// IsRetry
		if fieldName == "isretry" {
			isretry, _ := strconv.Atoi(value)
			response.IsRetry = uint32(isretry)
		}

		// IsOverLoad
		if fieldName == "isoverload" {
			isoverload, _ := strconv.Atoi(value)
			response.IsOverLoad = uint32(isoverload)
		}

		// IsGetCdn
		if fieldName == "isgetcdn" {
			isgetcdn, _ := strconv.Atoi(value)
			response.IsGetCDN = uint32(isgetcdn)
		}

		// XClientIP
		if fieldName == "x-ClientIp" {
			response.XClientIP = value
		}

	}
	dataJson, _ := json.Marshal(response)
	log.Info(string(dataJson))
	return response, nil
}
