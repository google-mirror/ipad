package baseutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/lunny/log"
	"hash/crc32"
	"math/big"
)

const (
	public_key_99      = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDf5W7uZQbl+XlrTxLDpIEhuE5U\njpmZ2DTiwDfjzSdunEorF1jFgqZ/bRKJXOVSXd5R0LktMri+eyyFgncpw1cdzBS1\ngYd7xjS8x/naOCXJeiWzQaZClQmDA8S1hOxXnsynyLlngvZdZQA57noHcsGV2+/E\nSIvfsLmljFwFjjqwTQIDAQAB\n-----END PUBLIC KEY-----\n"
	public_key_189     = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7Czt3VSBlM6KjpJYBfVE\nR9G7t9XeRwHxvx/NFLCIuJqL3KoiEVwfVW2ZRyKcrYR8y92hC9ojyk6Xm/O8WVVK\n+Ly9ruBjSjWG7nG1DrPZIGGuJI2Uc9yXvg7vSBmoNo/JH6WwVGQPsWknOe/jnAzs\n4OuTRi0A5DeSJclENnVth6Xi9lke50M9eS1S3Xy9GPSEbLg5CGG/E9r04p/SvP4g\n9UcBYyj21EvpIrWei6QyC2b9wbJosc+r8opF5LzTab45NM2suvTNIPdJVXVjGaWo\n+ZN1tK9PRKHTIabJPmS3eYyIl56KR1o7IRd71iQGN/omcJKfjlCcF3x2Gyto0o+E\nWwIDAQAB\n-----END PUBLIC KEY-----\n"
	public_key_190     = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA21WeytqtK3+fBi3AistJ\n4Zt9WIDaMadUSEPT+nmUYH8KtaWYZnl99fnAbs+hhcsf7ziE57KKfPW9VbMk2D5j\nKnOld18BcWIWb/obhGtJoTceA6J/cXKvEXkD1SMM0ZpadVK12rJlNgozV2hc3TsV\n0wteK9EYennlQot/pjIniAOSyUzeROmMvkzf53v3syzwuA3CQipdoAO3pNN86urk\ntiRb+HtQbKjeahHAEIws4XV4okSRGrf9DQeLzObxaFl6rpqxSBMl+HfZEL3rBvIW\ns9+dSb5Y56Kq32MMVfKkQXZ+awU+DyCyv7KbAyaDX/LCzWP7USTWsgcQsE+BoxMM\nSwIDAQAB\n-----END PUBLIC KEY-----\n"
	public_key_133     = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvWpUR3ZA8MCyCdt3RxJo\nlrJ/trIZq5vJxM2WYfQi4UOnWrLDTquI9EcZ2NLg1XzslxN0i/gh7CAU35ewHM4m\nLyfKJPTYlJL5ncjBpBTQuOdg2BXfU6kR1dgHyvaCcIS76CWknBu5NpZ1xL5DVZdW\nW1xCIgkCNfalWVAD1dX6Z4Dr1Rzqx20D2Oufl7RSmXGffDUrLvMkSeD90JtWK6Ax\ndBi2b8CFPqn1/6heq4oU4nhcArDKxq/UUO5aaXHCIOcv5vpLeBI1850gZzTJl0En\n42nkeb8yVf/4xfpLEzxkKlZWqOXxdkcsWj/hjYgW5A5Yq8KkoyugVusLUEyG2uBZ\nBwIDAQAB\n-----END PUBLIC KEY-----"
	public_key_133_hex = "BD6A54477640F0C0B209DB7747126896B27FB6B219AB9BC9C4CD9661F422E143A75AB2C34EAB88F44719D8D2E0D57CEC9713748BF821EC2014DF97B01CCE262F27CA24F4D89492F99DC8C1A414D0B8E760D815DF53A911D5D807CAF6827084BBE825A49C1BB9369675C4BE435597565B5C4222090235F6A5595003D5D5FA6780EBD51CEAC76D03D8EB9F97B45299719F7C352B2EF32449E0FDD09B562BA0317418B66FC0853EA9F5FFA85EAB8A14E2785C02B0CAC6AFD450EE5A6971C220E72FE6FA4B781235F39D206734C9974127E369E479BF3255FFF8C5FA4B133C642A5656A8E5F176472C5A3FE18D8816E40E58ABC2A4A32BA056EB0B504C86DAE05907"
	public_key_135     = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtXkUc/36zOQmBYQBthJa\nPW/t12x90bBCanPYpBgrKeptBfT16NmaTT0cPlzzyMs83fk1ZDyU04kniBsUTQTz\nEPEzB9GuY6EAonl6cUwNHipaDvd5/D1vfTwzlidr8n2m1m4mlqZVfv1LYZDHJolN\nNc5VnhR5abrASv67DjojWyx5WsapgY4UozpEaPj/ar6KVKdBgAQr8P04Qn9wtoG5\nQxoJnndGGNRV8U0fdRIVd9rmbDhToqqcTw+cIhpm9kpG1faLDVDyLH5PoNhASLL5\nF59LhkQqJyDI/ie8aMXGOE3MM2+XkU8niLkF5f6Yxbt1RIiw9rCUIbsnv/UY7w6S\nmQIDAQAB\n-----END PUBLIC KEY-----"
	public_key_135_hex = "B5791473FDFACCE426058401B6125A3D6FEDD76C7DD1B0426A73D8A4182B29EA6D05F4F5E8D99A4D3D1C3E5CF3C8CB3CDDF935643C94D38927881B144D04F310F13307D1AE63A100A2797A714C0D1E2A5A0EF779FC3D6F7D3C3396276BF27DA6D66E2696A6557EFD4B6190C726894D35CE559E147969BAC04AFEBB0E3A235B2C795AC6A9818E14A33A4468F8FF6ABE8A54A74180042BF0FD38427F70B681B9431A099E774618D455F14D1F75121577DAE66C3853A2AA9C4F0F9C221A66F64A46D5F68B0D50F22C7E4FA0D84048B2F9179F4B86442A2720C8FE27BC68C5C6384DCC336F97914F2788B905E5FE98C5BB754488B0F6B09421BB27BFF518EF0E9299"
	public_key_125_hex = "D8D2AE73FF601B93B1471B35870A1B59D7649EEA815CDD8CE5496BBD0C6CFE19C0E082F4E513B615C6030CCFCE3153E25AA00E8156D0311AF72ABBB9BBEC8B1D3751592234B1A621CA774E2EC50047A93FA0BC60DF0C10E8A65C3B29D13167EC217FC6A29034494870705CBF4AC929FBA0E1E656A8F8B50E779AD89BB4EEF6FF"
	public_key_182_hex = "C8930AB6E688F68513682FA555E1A3C175867CDC4AAE1D054F75134D1553D9E4A1BBC846FBAEE947E1515363365185AEB39C9DD5B76BF8ADE21233E27728BC0ED8C465CCC7DBABC7EBE08B1FA23A89098D7730C31FBE375745A9AA717D7F3DE5FB8126B6D6B2B9EB1643346F00EBE3A2AA915A417B263E4026FEF4BFA91B81B035DE224857E87FB292FB9AFBAC45725D968068385963E3CDB0162C901D0921030515D1CA1E079129DD585969EF6CBABBE72E287D9A9F757FAE91543F5AADD96777BE49D1CEF58669250DB4B992C01AFE22AF5DDDB0147BD5C8B373F39381A1E914078050239D490B6FC02E68B61A81BCE0ED9710BC84481273E32ED89ACD5211"
)

// Rsa的publicKey -- 135
var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtXkUc/36zOQmBYQBthJa
PW/t12x90bBCanPYpBgrKeptBfT16NmaTT0cPlzzyMs83fk1ZDyU04kniBsUTQTz
EPEzB9GuY6EAonl6cUwNHipaDvd5/D1vfTwzlidr8n2m1m4mlqZVfv1LYZDHJolN
Nc5VnhR5abrASv67DjojWyx5WsapgY4UozpEaPj/ar6KVKdBgAQr8P04Qn9wtoG5
QxoJnndGGNRV8U0fdRIVd9rmbDhToqqcTw+cIhpm9kpG1faLDVDyLH5PoNhASLL5
F59LhkQqJyDI/ie8aMXGOE3MM2+XkU8niLkF5f6Yxbt1RIiw9rCUIbsnv/UY7w6S
mQIDAQAB
-----END PUBLIC KEY-----`)

var cdnRsaPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC/7f+16ihQn5yJ7YP6f92oiBQ1
1ETphNU6mK2OlBDxFF7dU3iQ4QRWGQsi5uUAZFXvxsEuQf2phfOPu8chPsuBDjBT
1LjXT/vHC0YAq9coICMir84UBgRmMSYb1e49RHIQgv6rdDQNc2RdwNAqKTuWK51H
5KZBAL11JN4A2dO1wQIDAQAB
-----END PUBLIC KEY-----`)

// encAesKeyPublicKey 这个是用来加密EncryptUserInfo Key的rsa公钥
var encAesKeyPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDoqn5y8C1D0g8LVr/LSmCBiNJC
4TdPsHa+qWBuq/xEhv01MVzIoN5Vc8UZ0qLjDeWh462Ou2Ps6E6O1C+VSTBxtLrp
zBOxB/iAUKgw5we4G8kNijeVumVro/AovUeqUJRrMVkSBMbp/O1WZfiK7bQN0UxF
zpQe0j8HZ9hwGWnjSQIDAQAB
-----END PUBLIC KEY-----`)

// AesEncryptKeyBytes Aes加密
func AesEncryptKeyBytes(orig []byte, key []byte) []byte {
	// 转成字节数组
	// origData := []byte(orig)
	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("AesEncryptKeyBytes - ", err.Error())
		return []byte{}
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	orig = PKCS7Padding(orig, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 创建数组
	cryted := make([]byte, len(orig))
	// 加密
	blockMode.CryptBlocks(cryted, orig)
	// return base64.StdEncoding.EncodeToString(cryted)
	return cryted
}

// AesEncrypt Aes加密
func AesEncrypt(orig []byte, key []byte) []byte {
	// 转成字节数组
	// origData := []byte(orig)
	if len(key) > 16 {
		key = key[:16]
	}
	return AesEncryptKeyBytes(orig, key)
}

// AesDecryptByteKey Aes CBC解密
func AesDecryptByteKey(cryted []byte, key []byte) ([]byte, error) {
	if len(cryted) <= 0 || len(key) <= 0 {
		return []byte{}, errors.New("AesDecryptByteKey err: len(cryted) <= 0 || len(key) <= 0")
	}
	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 创建数组
	orig := make([]byte, len(cryted))
	// 解密
	blockMode.CryptBlocks(orig, cryted)
	// 去补全码
	return PKCS7UnPadding(orig)
}

// AesDecrypt aes CBC解密
func AesDecrypt(cryted []byte, key []byte) ([]byte, error) {
	return AesDecryptByteKey(cryted, key)
}

// PKCS7Padding 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if unpadding > length {
		fmt.Println(hex.EncodeToString(origData))
		return []byte{}, errors.New("PKCS7UnPadding err: unpadding > length")
	}

	return origData[:(length - unpadding)], nil
}

// RsaEncryptByVer 根据RsaKey 版本进行加密
func RsaEncryptByVer(origData []byte, ver uint32) ([]byte, error) {
	var publicKey []byte
	switch ver {
	case 133: //新疆号使用
		publicKey = []byte(public_key_133)
	case 135: //国内国外
		publicKey = []byte(public_key_135)
	case 190:
		publicKey = []byte(public_key_190)
	case 189:
		publicKey = []byte(public_key_189)
	case 99:
		publicKey = []byte(public_key_99)

	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RsaEncrypt 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// EncKeyRsaEncrypt 加密
func EncKeyRsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(encAesKeyPublicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// CdnRsaEncrypt Cdn,Aes加密
func CdnRsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(cdnRsaPubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// Md5Value 计算字符串的MD5值
func Md5Value(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func MD5ToLower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ALLGather(Data []string) string {
	var M string
	for i := 0; i < len(Data); i++ {
		M += Data[i]
	}
	return M
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

// Md5ValueByte 计算字节数组的MD5值
func Md5ValueByte(data []byte, bUpper bool) string {
	has := md5.Sum(data)
	tmpString := "%x"
	if bUpper {
		tmpString = "%X"
	}
	md5str := fmt.Sprintf(tmpString, has)
	return md5str
}

// Md5Value16 计算16为的md5值
func Md5Value16(data []byte) [16]byte {
	return md5.Sum(data)
}

// Adler32 计算Adler32
func Adler32(adler uint32, data []byte) uint32 {
	s1 := adler & 0xffff
	s2 := (adler >> 16) & 0xffff

	length := len(data)
	for index := 0; index < length; index++ {
		s1 = (s1 + uint32(data[index])) % 65521
		s2 = (s2 + s1) % 65521
	}

	return (s2 << 16) + s1
}

// AesEncryptECB Aes ECB模式加密
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(key)
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// AesDecryptECB AesDecryptECB
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(key)
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

// Sha1 计算Sha值
func Sha1(data []byte) []byte {
	sha1 := sha1.New()
	sha1.Write(data)
	return sha1.Sum([]byte(""))
}
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

// 根据Key 版本进行加密
func RsaEncryptByVerNew(origData []byte, ver uint32) ([]byte, error) {
	publicKey := ""

	switch ver {
	case 133:
		publicKey = public_key_133_hex
		break
	case 135:
		publicKey = public_key_135_hex
		break
	case 125:
		publicKey = public_key_125_hex
	case 182:
		publicKey = public_key_182_hex

	}

	bigInt := new(big.Int)
	n, ok := bigInt.SetString(publicKey, 16)
	if !ok {
		return nil, errors.New("转换失败")
	}

	pub := &rsa.PublicKey{
		N: n,
		E: 65537,
	}

	partLen := pub.N.BitLen()/8 - 11
	chunks := split(origData, partLen)
	buffer := bytes.NewBuffer(nil)
	for _, chunk := range chunks {
		bs, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bs)
	}

	return buffer.Bytes(), nil
}

// NoCompressRsaByVer NoCompressRsa 根据Key 版本进行加密
// todo：2020/6/17 12:44 增加
func NoCompressRsaByVer(data []byte, ver uint32) []byte {

	retData, err := RsaEncryptByVerNew(data, ver)
	if err != nil {
		fmt.Println("NoCompressRsaByVer err == ", err.Error())
	}
	return retData
}

// CompressAndRsaByVer 压缩后进行rsa加密  根据Key 版本进行加密
// todo：2020/6/17 12:44 增加
func CompressAndRsaByVer(data []byte, ver uint32) []byte {
	newData := CompressByteArray(data)
	retData, err := RsaEncryptByVerNew(newData, ver)
	if err != nil {
		fmt.Println("CompressAndRsaByVer err == ", err.Error())
	}
	return retData
}

// NoCompressRsa NoCompressRsa
func NoCompressRsa(data []byte) []byte {
	rsaKeySize := 2048
	rsaLen := rsaKeySize/8 - 12
	dataLen := len(data)

	retData := make([]byte, 0)
	if dataLen <= rsaLen {
		retData, err := RsaEncrypt(data)
		if err != nil {
			log.Info("CompressAndRsa err = ", err)
		}

		return retData
	}

	// 分块加密
	blockCnt := dataLen / rsaLen
	if dataLen%rsaLen != 0 {
		blockCnt = blockCnt + 1
	}
	var tmpBlockCnt = blockCnt
	for tmpBlockCnt > 0 {
		startPos := (blockCnt - tmpBlockCnt) * rsaLen
		blockSize := startPos + rsaLen
		if tmpBlockCnt == 1 {
			blockSize = dataLen
		}

		tmpBuf := make([]byte, 0)
		tmpBuf = append(tmpBuf, data[startPos:blockSize]...)
		encodedData, err := RsaEncrypt(tmpBuf)
		if err != nil {
			log.Info("CompressAndRsa err = ", err)
		}
		retData = append(retData, encodedData[0:]...)
		tmpBlockCnt = tmpBlockCnt - 1
	}

	return retData
}

// CompressAndRsa 压缩后进行rsa加密
func CompressAndRsa(data []byte) []byte {
	newData := CompressByteArray(data)
	rsaKeySize := 2048
	rsaLen := rsaKeySize/8 - 12
	newDataLen := len(newData)

	if newDataLen <= rsaLen {
		newData, err := RsaEncrypt(newData)
		if err != nil {
			log.Info("CompressAndRsa err = ", err)
		}
		return newData
	}

	retData := make([]byte, 0)
	// 分块加密
	blockCnt := newDataLen / rsaLen
	if newDataLen%rsaLen != 0 {
		blockCnt = blockCnt + 1
	}
	var tmpBlockCnt = blockCnt
	for tmpBlockCnt > 0 {
		startPos := (blockCnt - tmpBlockCnt) * rsaLen
		blockSize := startPos + rsaLen
		if tmpBlockCnt == 1 {
			blockSize = newDataLen
		}

		tmpBuf := make([]byte, 0)
		tmpBuf = append(tmpBuf, newData[startPos:blockSize]...)
		encodedData, err := RsaEncrypt(tmpBuf)
		if err != nil {
			PrintLog(err.Error())
		}
		retData = append(retData, encodedData[0:]...)
		tmpBlockCnt = tmpBlockCnt - 1
	}

	return retData
}

// CompressAes 压缩然后aes加密
func CompressAes(aesKey []byte, data []byte) []byte {
	newData := CompressByteArray(data)
	newData = AesEncrypt(newData, aesKey)

	return newData
}

// DecryptSnsVideoData 解密视频
func DecryptSnsVideoData(data []byte, encLen uint32, tmpKey uint64) []byte {
	aacInst := CreateISAacInst(tmpKey)
	// 解密数据
	for index := uint32(0); index < encLen; index += 8 {
		randNumber := ISAacRandom(aacInst)
		tmpBytes := UInt64ToBytes(randNumber)
		for tmpIndex := uint32(0); tmpIndex < 8; tmpIndex++ {
			realIndex := index + tmpIndex
			if realIndex >= encLen {
				return data
			}
			data[index+tmpIndex] ^= tmpBytes[tmpIndex]
		}
	}
	return data
}

func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

// PKCS5Padding PKCS5Padding
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS5UnPadding PKCS5UnPadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// EncryptDESECB ECB加密
func EncryptDESECB(data []byte, keyByte []byte) ([]byte, error) {
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// DecryptDESECB DesEC解密
func DecryptDESECB(data []byte, key []byte) ([]byte, error) {
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("DecryptDES crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	// out = PKCS5UnPadding(out)
	return out, nil
}

// Encrypt3DES 3DES加密
// src: 8字节
// key: 16字节
func Encrypt3DES(srcData []byte, key []byte) ([]byte, error) {
	if len(srcData) != 8 || len(key) != 16 {
		return nil, errors.New("Encrypt3DES err: srcLen != 8 || keyLen != 16")
	}
	tmpSrcData := make([]byte, 0)
	tmpSrcData = append(tmpSrcData, srcData...)
	encData, err := EncryptDESECB(tmpSrcData, key[0:8])
	if err != nil {
		return nil, err
	}
	decData, err := DecryptDESECB(encData, key[8:])
	if err != nil {
		return nil, err
	}
	encData, err = EncryptDESECB(decData[0:8], key[0:8])
	if err != nil {
		return nil, err
	}
	return encData[0:8], err
}

// Decrypt3DES Encrypt3DES 3DES解密
func Decrypt3DES(src []byte, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}

// AesEncryptECBTest Aes ECB模式加密
func AesEncryptECBTest(origData []byte, key []byte, pad byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(key)
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// CalcMsgCrc calc msg hash
func CalcMsgCrc(data []byte) int {
	salt1 := [16]byte{0x5c, 0x50, 0x7b, 0x6b, 0x65, 0x4a, 0x13, 0x09, 0x45, 0x58, 0x7e, 0x11, 0x0c, 0x1f, 0x68, 0x79}
	salt2 := [16]byte{0x36, 0x3a, 0x11, 0x01, 0x0f, 0x20, 0x79, 0x63, 0x2f, 0x32, 0x14, 0x7b, 0x66, 0x75, 0x02, 0x13}
	pad1 := [0x30]byte{
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
	}
	pad2 := [0x30]byte{
		0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c,
		0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c,
		0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c, 0x5c,
	}
	md5hash := md5.Sum(data)
	hashstr := hex.EncodeToString(md5hash[:])
	hash1Data := make([]byte, 16)
	copy(hash1Data, salt1[:])
	hash1Data = append(hash1Data, pad1[:]...)
	hash1 := sha1.New()
	hash1.Write(hash1Data)
	hash1.Write([]byte(hashstr))
	h1 := hash1.Sum(nil)
	hash2Data := make([]byte, 16)
	copy(hash2Data, salt2[:])
	hash2Data = append(hash2Data, pad2[:]...)
	hash2 := sha1.New()
	hash2.Write(hash2Data)
	hash2.Write(h1)
	h2 := hash2.Sum(nil)
	var b1, b2, b3 byte
	size := len(h2)
	for i := 0; i < size-2; i++ {
		b1 = h2[i+0] - b1*0x7d
		b2 = h2[i+1] - b2*0x7d
		b3 = h2[i+2] - b3*0x7d
	}
	return (int(0x21) << 24) | ((int(b3) & 0x7f) << 16) | ((int(b2) & 0x7f) << 8) | (int(b1) & 0x7f)
}

func GetCRC32String(msg string) uint32 {
	v := []byte(msg)
	return GetCRC32(v)
}

func GetCRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// 移位操作
func EncInt(d int64) int64 {
	a, b := int64(0), int64(0)
	for i := 0; i < 16; i++ {
		a |= ((1 << (2 * i)) & d) << (2 * i)
		b |= ((1 << (2*i + 1)) & d) << (2*i + 1)
	}
	return a | b
}

func IOSGetCidMd5(DeviceId, Cid string) string {
	Md5Data := MD5ToLower(DeviceId + Cid)
	return "A136" + Md5Data[5:]
}

func IOSGetCid(s int) string {
	M := inttobytes(s >> 12)
	return hex.EncodeToString(M)
}

func inttobytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func IOSUuid(DeviceId string) (uuid1 string, uuid2 string) {
	Md5DataA := MD5ToLower(DeviceId + "SM2020032204320")
	//"b58b9b87-5124-4907-8c92-d8ad59ee0430"
	fmt.Println("Md5DataA", string(Md5DataA))
	uuid1 = fmt.Sprintf("%s-%s-%s-%s-%s", Md5DataA[0:8], Md5DataA[2:6], Md5DataA[3:7], Md5DataA[1:5], Md5DataA[20:32])
	Md5DataB := MD5ToLower(DeviceId + "BM2020032204321")
	uuid2 = fmt.Sprintf("%s-%s-%s-%s-%s", Md5DataB[0:8], Md5DataB[2:6], Md5DataB[3:7], Md5DataB[1:5], Md5DataB[20:32])
	return
}
