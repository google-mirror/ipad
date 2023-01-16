package clientsdk

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"feiyu.com/wx/clientsdk/cecdh"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/hkdf"
	"io"
	"log"
	rnd "math/rand"
	"time"
)

const (
	//145 新疆号可以登录
	WeChatPubKey_145 = "0493b4723be07a56d81e6f1994a55597ade831b78c60b8c8a9a5ad656bc10534a6ecfcd6b2504ccbc9682dc9c6f5d213bc3cc3e98a5d759747d4f9c7c46f0aa6fa"
	//146
	//WeChatPubKey_146 = "047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5"
	WeChatPubKey_146 = "0495bc6e5c1331ad172d0f35b1792c3ce63f91572abd2dd6df6dac2d70195c3f6627cca60307305d8495a8c38b4416c75021e823b6c97dffe79c14cb7c3af8a586"
)

type EcdhKeyPair struct {
	PriKey []byte
	PubKey []byte
}

func HybridEcdhDecrypt(data, priKey, pubKey, FinalShaData []byte) ([]byte, error) {
	hybridEcdhDecryptResp := &wechat.HybridDecryptResponse{}
	if err := proto.Unmarshal(data, hybridEcdhDecryptResp); err != nil {
		return nil, err
	}
	//进行ecdh
	secretKey, err := ECDH(priKey, hybridEcdhDecryptResp.Key.Buffer)
	if err != nil {
		return nil, err
	}

	/*logger.Debugln(hex.EncodeToString(priKey))
	logger.Debugln(hex.EncodeToString(hybridEcdhDecryptResp.Key.Buffer))
	logger.Debugln("HybridEcdhDecrypt -> secretKey: ", hex.EncodeToString(secretKey))*/

	//pubKey,_ = hex.DecodeString(WeChatPubKey)
	h1 := sha256.New()
	h1.Write(FinalShaData)
	h1.Write([]byte{0x34, 0x31, 0x35})
	h1.Write(hybridEcdhDecryptResp.Key.Buffer)
	h1.Write([]byte{0x31})
	h1Sha256RetData := h1.Sum(nil)

	deProtoBuf, err := AesGcmDecryptWithUnCompress(secretKey[:24], h1Sha256RetData, hybridEcdhDecryptResp.ProtobufData)
	if err != nil {
		return nil, err
	}
	return deProtoBuf, nil
}

func HybridEncrypt(data []byte, WeChatPubKey string) (protoEnData []byte, epKey []byte, token []byte, ecdhpairkey *EcdhKeyPair, err error) {
	//pubKey2,_ := hex.DecodeString("046ECE0D01D24E9360397CDB0B44B07FC94312E0DDCEB0C671C410A475BEF4C115844C8C98C78C6C17AA8A547B5EDBCA62DD40B3ABBDFC1A7B4AE5C17B144C6D31")
	//priKey2,_ := hex.DecodeString("30770201010420F8D96C0C7DD8474D5CA4A5CC374825DBC9D5BC144D90DA6A90F07509AD05ADDBA00A06082A8648CE3D030107A144034200046ECE0D01D24E9360397CDB0B44B07FC94312E0DDCEB0C671C410A475BEF4C115844C8C98C78C6C17AA8A547B5EDBCA62DD40B3ABBDFC1A7B4AE5C17B144C6D31")
	//生成公私密钥
	ecdhKeyPair, err := GenEcdhKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	/*	logger.Debugln("GenEcdhKeyPair -> PriKey: ", ecdhKeyPair.PriKey)
		logger.Debugln("GenEcdhKeyPair -> PubKey: ", ecdhKeyPair.PubKey)
	*/
	//进行ECDH
	pubKey, _ := hex.DecodeString(WeChatPubKey)
	secretKey, err := ECDH(ecdhKeyPair.PriKey, pubKey)
	if err != nil {
		log.Fatal(err)
	}
	//ecdh成功后取前24 为作为AesKey
	//logger.Debugln("SecretKey -> ", hex.EncodeToString(secretKey), len(secretKey))
	h1 := sha256.New()
	// 31343135 进行sha256
	h1.Write([]byte{0x31, 0x34, 0x31, 0x35})
	// 对生成的pubKey 进行sha256
	h1.Write(ecdhKeyPair.PubKey)
	h1SHA256RetData := h1.Sum(nil)
	//生成32个字节的随机数
	randomBytes := GenRandomBytes(32)
	//对随机生成32个字节进行加密
	enData, err := AesGcmEncryptWithCompress(secretKey[:24], h1SHA256RetData, randomBytes)
	if err != nil {
		log.Fatal(err)
	}
	//logger.Debugln(enData)
	// 对随机生成的 32 个字节和 sha256 的结果进行密钥扩展HKDF
	epKey = hkdfEP(randomBytes, h1SHA256RetData, []byte("security hdkf expand"))
	//log.Println("securityHdkfExpand -> ", hex.EncodeToString(epKey))
	//再次进行Sha256

	/*h2 := sha256.New()
	//1 + 415 字符串进行sha256
	h2.Write([]byte{0x31,0x34,0x31,0x35})
	//再次对生成PubKey 进行Sha256
	h2.Write(pubKey2)
	//对上次AEsGCm加密后结果进行sha256
	h2.Write(enData)*/
	//生成32个字节的随机数 ToKen 进行sha256
	/*token = GenRandomBytes(32)
	h2.Write(token)*/
	h1.Write(enData)
	h2SHA256RetData := h1.Sum(nil)

	//对ProtoBuf进行加密
	//取HKDF 后的数据取前24个字节
	protoEnData, err = AesGcmEncryptWithCompress(epKey[:24], h2SHA256RetData, data)
	if err != nil {
		log.Fatal(err)
	}

	/*ecdhKeyPair.PubKey = pubKey2
	ecdhKeyPair.PriKey = priKey2*/

	//需要返回的数据 epKey ecdhKeyPair token protoEnData
	//最后还需要对 原data 和 epKey 进行最后的 hash256
	return protoEnData, epKey, enData, ecdhKeyPair, nil
}

func GenEcdhKeyPair() (*EcdhKeyPair, error) {
	ecdh := cecdh.NewEllipticECDH(elliptic.P256())
	priKey, pubKey, err := ecdh.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &EcdhKeyPair{
		PriKey: priKey,
		PubKey: pubKey,
	}, nil
}

func ECDH(priKey, pubKey []byte) (SecretKey []byte, err error) {
	ecdh := cecdh.NewEllipticECDH(elliptic.P256())
	publicKey, ok := ecdh.Unmarshal(pubKey)
	if !ok {
		return nil, errors.New("序列化失败！")
	}
	SecretKey, err = ecdh.GenerateSharedSecret(priKey, publicKey)
	if err != nil {
		return
	}

	h := sha256.New()
	h.Write(SecretKey)
	SecretKey = h.Sum(nil)
	return
}

//Aes-Gcm 压缩加密
func AesGcmEncryptWithCompress(key, aad, data []byte) ([]byte, error) {
	compressData := DoZlibCompress(data)
	//生成12个字节的nonce
	nonce := GenRandomBytes(12)
	//logger.Debugln("nonce -> ",hex.EncodeToString(nonce))
	//logger.Debugln("aad   -> ",hex.EncodeToString(aad))
	enData, err := AesGcmEncrypt(key, nonce, aad, compressData)
	if err != nil {
		return nil, err
	}
	newData := append([]byte{}, enData[:len(enData)-16]...)
	newData = append(newData, nonce...)
	newData = append(newData, enData[len(enData)-16:]...)
	//logger.Debugln("enData  -> ", hex.EncodeToString(enData))
	//logger.Debugln("newData -> ", hex.EncodeToString(newData))
	return newData, nil
}

func AesGcmDecryptWithUnCompress(key, aad, data []byte) ([]byte, error) {
	newData := append([]byte{}, data[:len(data)-28]...)
	newData = append(newData, data[len(data)-16:]...)
	//logger.Debugln("Data    -> ", hex.EncodeToString(data))
	//logger.Debugln("newData -> ", hex.EncodeToString(newData))
	nonce := data[len(data)-(12+16) : len(data)-16]
	//logger.Debugln("nonce -> ",hex.EncodeToString(nonce))
	deData, err := AesGCMDecrypt(key, nonce, aad, newData)
	if err != nil {
		return nil, err
	}
	return DoZlibUnCompress(deData)
}

// 进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, _ = w.Write(src)
	_ = w.Close()
	return in.Bytes()
}

func GenRandomBytes(length int) []byte {
	randomStr := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	var result []byte
	r := rnd.New(rnd.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, randomStr[r.Intn(len(randomStr))])
	}
	return result
}

func hkdfEP(randData, sha256Data, info []byte) []byte {
	ex := hkdf.Extract(sha256.New, randData, []byte("security hdkf expand"))
	expand := hkdf.Expand(sha256.New, ex, sha256Data)
	expandBytes := make([]byte, 56)
	_, _ = expand.Read(expandBytes)
	return expandBytes
}

/*func hkdfEP(secret, salt, info []byte) []byte {
	reader := hkdf.New(sha256.New, secret, salt,info)
	ex := make([]byte, 56)
	reader.Read(ex)
	return ex
}*/

func AesGCMDecrypt(key, nonce, aad, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, data, aad)
}

func AesGcmEncrypt(key, nonce, aad, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, data, aad), nil
}

// 进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
