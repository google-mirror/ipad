package baseutils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gogf/guuid"
	"image/png"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 转换16进制
var numberHexSmall = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
var numberHexBig = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
var stringBytes = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
	'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '-', '+'}
var bcd2Bytes = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?'}

// RandomUUID 随机生成UUID
func RandomUUID() string {
	/*retAdSource := string("")
	retAdSource = retAdSource + RandomBigHexString(8) + "-"
	retAdSource = retAdSource + RandomBigHexString(4) + "-"
	retAdSource = retAdSource + RandomBigHexString(4) + "-"
	retAdSource = retAdSource + RandomBigHexString(4) + "-"
	retAdSource = retAdSource + RandomBigHexString(12)*/

	return guuid.New().String()
}

// RandomBSSID 随机生成RandomBSSID
func RandomBSSID() string {
	retAdSource := string("")
	retAdSource = retAdSource + RandomSmallHexString(2) + "-"
	retAdSource = retAdSource + RandomSmallHexString(2) + "-"
	retAdSource = retAdSource + RandomSmallHexString(2) + "-"
	retAdSource = retAdSource + RandomSmallHexString(2) + "-"
	retAdSource = retAdSource + RandomSmallHexString(2) + "-"
	retAdSource = retAdSource + RandomSmallHexString(2)

	return retAdSource
}

func BuildRandomMac() string {
	// 00:00:00:00:00:00
	macs := make([]string, 6)
	for i := range macs {
		macs[i] = fmt.Sprintf("%02X", rand.Intn(0xff))
	}

	return strings.Join(macs, ":")
}

//RandomBytes RandomBytes
func RandomBytes(length uint32) []byte {
	retBytes := make([]byte, length)
	tmpRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i uint32
	for i = 0; i < length; i++ {
		retBytes[i] = byte(tmpRand.Intn(256))
	}

	return retBytes
}

//RandomBigHexString RandomBigHexString
func RandomBigHexString(length uint32) string {
	retBytes := make([]byte, length)
	tmpRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i uint32
	for i = 0; i < length; i++ {
		tmpIndex := tmpRand.Intn(16)
		retBytes[i] = numberHexBig[tmpIndex]
	}

	return string(retBytes)
}

//RandomSmallHexString RandomSmallHexString
func RandomSmallHexString(length uint32) string {
	retBytes := make([]byte, length)
	tmpRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i uint32
	for i = 0; i < length; i++ {
		tmpIndex := tmpRand.Intn(16)
		retBytes[i] = numberHexSmall[tmpIndex]
	}

	return string(retBytes)
}

//RandomStringByLength 随机固定长度的字符串
func RandomStringByLength(length uint32) []byte {
	retBytes := make([]byte, length)
	tmpRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmpSize := len(stringBytes)
	var index = uint32(0)
	for index = 0; index < length; index++ {
		tmpIndex := tmpRand.Intn(tmpSize)
		retBytes[index] = stringBytes[tmpIndex]
	}

	return retBytes
}

//RandomString 随机长度的字符串
func RandomString(smallLength uint32, bigLength uint32) []byte {
	if smallLength >= bigLength || smallLength <= 0 {
		return []byte{}
	}

	tmpRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmpSize := bigLength - smallLength
	tmpLength := uint32(tmpRand.Intn(int(tmpSize))) + smallLength

	return RandomStringByLength(tmpLength)
}

//BytesToInt32 bytes 转 uint32
func BytesToInt32(bytesData []byte) uint32 {
	length := len(bytesData)
	if length < 4 {
		tmpBytes := make([]byte, 4-length)
		tmpBytes = append(tmpBytes, bytesData[0:]...)
		return binary.BigEndian.Uint32(tmpBytes)
	}
	return binary.BigEndian.Uint32(bytesData)
}

// BytesToInt32SmallEndian 小端整形
func BytesToInt32SmallEndian(bytesData []byte) uint32 {
	return binary.LittleEndian.Uint32(bytesData)
}

//Int32ToBytes uint32 转 bytes
func Int32ToBytes(int32Value uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, int32Value)
	return buf
}

// UInt64ToBytes uint64 转 bytes
func UInt64ToBytes(int64Value uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, int64Value)
	return buf
}

// BigEndianInt32ToLittleEndianInt32 大端整形 转成小端整形
func BigEndianInt32ToLittleEndianInt32(intValue uint32) uint32 {
	tmpBytes := Int32ToBytes(intValue)
	return binary.LittleEndian.Uint32(tmpBytes)
}

// Int32ToBytesLittleEndian uint32 转 bytes
func Int32ToBytesLittleEndian(int32Value uint32) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, int32Value)
	return buf
}

// Int16ToBytesLittleEndian uint16 转 bytes
func Int16ToBytesLittleEndian(int16Value uint16) []byte {
	var buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, int16Value)
	return buf
}

// Int16ToBytesBigEndian uint16 转 bytes
func Int16ToBytesBigEndian(int16Value uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, int16Value)
	return buf
}

// BytesToUint16BigEndian bytes 转 uint16 BigEndian
func BytesToUint16BigEndian(bytesData []byte) uint16 {
	return binary.BigEndian.Uint16(bytesData)
}

// BytesToUint16LittleEndian bytes 转 uint16 LittleEndian
func BytesToUint16LittleEndian(bytesData []byte) uint16 {
	return binary.LittleEndian.Uint16(bytesData)
}

// Int2Byte int 转成 字节数组
func Int2Byte(data int) (ret []byte) {
	var len = unsafe.Sizeof(data)
	ret = make([]byte, len)
	var tmp = 0xff
	var index = uint(0)
	for index = 0; index < uint(len); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

//DecodeVByte32 获取变长整形值 和 所占用的字节数
// data ： 变长整数的字节数组
// current : 数组的起始坐标, 从这个坐标开启解析变长整数
// return 解密后的整形、占用的字节数
func DecodeVByte32(data []byte, current uint32) (uint32, uint32) {
	var index = uint32(0)
	for data[current+index] >= 0x80 {
		index++
	}

	var retLen = uint32(index + 1)
	var value = uint32(0)
	for index != 0 {
		value <<= 7
		value |= (uint32)(data[current+index] & 0x7f)
		index--
	}

	value <<= 7
	value |= (uint32)(data[current] & 0x7f)
	return value, retLen
}

//EncodeVByte32 uint32 转成 变长整形字节数组
func EncodeVByte32(num uint32) []byte {
	retBytes := make([]byte, 0)
	tmpValue := num
	for tmpValue >= 0x80 {
		var tmpByte = byte(0)
		tmpByte = byte(tmpValue%0x80 + 0x80)
		tmpValue = tmpValue / 0x80
		retBytes = append(retBytes, tmpByte)
	}

	var tmpByte = byte(tmpValue)
	retBytes = append(retBytes, tmpByte)

	return retBytes
}

//UInt32To16Bytes 整形转换成16进制字节数组
func UInt32To16Bytes(intValue uint32) []byte {
	retBytes := make([]byte, 0)

	tmpValue := intValue
	for tmpValue > 0 {
		modValue := tmpValue % 16
		retBytes = append(retBytes, numberHexSmall[modValue])
		tmpValue = tmpValue / 16
	}

	retBytes = reverse(retBytes)
	return retBytes
}

// 反转字节数组
func reverse(data []byte) []byte {
	retBytes := make([]byte, 0)
	dataLen := len(data)
	for dataLen > 0 {
		retBytes = append(retBytes, data[dataLen-1])
		dataLen = dataLen - 1
	}

	return retBytes
}

//HashCode 计算字符串hashcode
func HashCode(str string) uint32 {
	tmpBytes := []byte(str)
	tmpLen := len(tmpBytes)

	var retHashCode = uint32(0)
	for index := 0; index < tmpLen; index++ {
		var tmpValue = uint32(tmpBytes[index])
		retHashCode = 31*retHashCode + tmpValue
	}

	return retHashCode
}

// StringCut 截取字符串
// srcStr 待截取的字符串
// index 截取的开始索引
// length 截取的长度
func StringCut(srcStr string, index uint32, length uint32) string {
	srcBytes := []byte(srcStr)
	tmpBytes := make([]byte, 0)
	tmpBytes = append(tmpBytes, srcBytes[index:index+length]...)

	retString := string(tmpBytes)
	return retString
}

// HexStringToBytes 16进制字符串 转bytes
func HexStringToBytes(hexString string) []byte {
	retBytes := make([]byte, 0)
	count := len(hexString)
	for index := 0; index < count; index += 2 {
		tmpStr := StringCut(hexString, uint32(index), 2)
		value64, _ := strconv.ParseInt(tmpStr, 16, 16)
		retBytes = append(retBytes, byte(value64))
	}

	return retBytes
}

// BytesToHexString 字节数组转16进制String
func BytesToHexString(data []byte, isBig bool) string {
	changeBytes := numberHexSmall
	if isBig {
		changeBytes = numberHexBig
	}
	length := len(data)
	retBytes := make([]byte, length*2)
	for index := 0; index < length; index++ {
		tmpByte := data[index]
		highIndex := ((tmpByte & 0xf0) >> 4)
		lowIndex := tmpByte & 0x0f
		retBytes[index*2] = changeBytes[highIndex]
		retBytes[index*2+1] = changeBytes[lowIndex]
	}

	return string(retBytes)
}

// EscapeURL 对Url特殊字符转译
func EscapeURL(srcURL string) string {
	size := len(srcURL)
	retURL := []byte{}
	for index := 0; index < size; index++ {
		// 转译：%
		if srcURL[index] == '%' {
			retURL = append(retURL, []byte("%25")[0:]...)
			continue
		}
		// 转译：&
		if srcURL[index] == '&' {
			retURL = append(retURL, []byte("%26")[0:]...)
			continue
		}
		// 转译：+
		if srcURL[index] == '+' {
			retURL = append(retURL, []byte("%2B")[0:]...)
			continue
		}
		// 转译：/
		if srcURL[index] == '/' {
			retURL = append(retURL, []byte("%2F")[0:]...)
			continue
		}
		// 转译：
		if srcURL[index] == ':' {
			retURL = append(retURL, []byte("%3A")[0:]...)
			continue
		}
		// 转译：=
		if srcURL[index] == '=' {
			retURL = append(retURL, []byte("%3D")[0:]...)
			continue
		}
		// 转译：?
		if srcURL[index] == '?' {
			retURL = append(retURL, []byte("%3F")[0:]...)
			continue
		}
		retURL = append(retURL, srcURL[index])
	}

	return string(retURL)
}

// HongBaoStringToBytes HongBaoStringToBytes
func HongBaoStringToBytes(content string) string {
	retBytes := []byte{}
	tmpBytes := []byte(content)
	count := len(tmpBytes)
	for index := 0; index < count; index++ {
		byteValue := tmpBytes[index]
		if byteValue&0x80 != 0 {
			high4Value := byteValue >> 4
			low4Value := byteValue & 0xF
			retBytes = append(retBytes, '%')
			retBytes = append(retBytes, numberHexBig[high4Value])
			retBytes = append(retBytes, numberHexBig[low4Value])
		} else {
			retBytes = append(retBytes, byteValue)
		}
	}

	return string(retBytes)
}

// GetNumberString 获取前面的数字字符串
func GetNumberString(srcStr string) string {
	retBytes := make([]byte, 0)
	count := len(srcStr)
	for index := 0; index < count; index++ {
		if srcStr[index] >= '0' && srcStr[index] <= '9' {
			retBytes = append(retBytes, srcStr[index])
		} else {
			break
		}
	}
	return string(retBytes)
}

// BCD2ToASCII BCD2ToASCII
func BCD2ToASCII(ascString string) []byte {
	retBytes := []byte{}
	size := len(ascString)
	if size%2 == 1 {
		size = size + 1
		ascString = "0" + ascString
	}
	tmpAscData := []byte(ascString)
	for index := 0; index < size; index += 2 {
		tmpByte := byte(0)
		for index2 := byte(0); index2 < byte(16); index2++ {
			if tmpAscData[index] == numberHexBig[index2] {
				tmpByte = tmpByte + index2<<4
			}
			// 判断
			if tmpAscData[index+1] == numberHexBig[index2] {
				tmpByte = tmpByte + index2
			}
		}
		retBytes = append(retBytes, tmpByte)
	}
	return retBytes
}

// ASCIIToBCD2 ASCIIToBCD2
func ASCIIToBCD2(ascData []byte, bFlag bool) []byte {
	retBytes := []byte{}
	exchangeBytes := bcd2Bytes
	if bFlag {
		exchangeBytes = numberHexBig
	}

	size := len(ascData)
	for index := 0; index < size; index++ {
		high4 := ascData[index] >> 4
		low4 := ascData[index] & 0xF
		retBytes = append(retBytes, exchangeBytes[high4])
		retBytes = append(retBytes, exchangeBytes[low4])
	}
	return retBytes
}

// ParseInt 去掉除数字外的其他字符
func ParseInt(srcData string) uint32 {
	retBytes := make([]byte, 0)
	tmpLen := len(srcData)
	for index := 0; index < tmpLen; index++ {
		if srcData[index] >= '0' && srcData[index] <= '9' {
			retBytes = append(retBytes, srcData[index])
		}
	}

	retValue, err := strconv.Atoi(string(retBytes))
	if err != nil {
		return 0
	}
	return uint32(retValue)
}

// CreateQRCode 生成二维码
func CreateQRCode(data string) ([]byte, error) {
	qrCode, _ := qr.Encode(data, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, qrCode)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
