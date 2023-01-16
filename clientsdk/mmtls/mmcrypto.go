package mmtls

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"

	"feiyu.com/wx/clientsdk/baseutils"
)

var verifyPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8uOhBSSfVijKin+SZO/0IXUrmf8l
9sa7VgqOIH/AO3Xb1MF4Xm25bBSb5znHsknQsNPSye3vVo80NUi2gEHw8g==
-----END PUBLIC KEY-----`)

// Sha256 Sha256
func Sha256(data []byte) []byte {
	hash256 := sha256.New()
	hash256.Write(data)
	hashRet := hash256.Sum(nil)

	return hashRet
}

// HmacHash256 HmacHash256
func HmacHash256(key []byte, data []byte) []byte {
	hmacTool := hmac.New(sha256.New, key)
	hmacTool.Write(data)
	return hmacTool.Sum(nil)
}

// HkdfExpand HkdfExpand
func HkdfExpand(key, message []byte, outLen int) []byte {
	result := make([]byte, 0)
	count := outLen / 32
	if outLen%32 != 0 {
		count = count + 1
	}
	for i := 1; i <= count; i++ {
		h := hmac.New(sha256.New, key)
		tmp := append(message, byte(i))
		tmp = append(result, tmp...)
		h.Write(tmp)
		result = append(result, h.Sum(nil)...)
	}
	return result[:outLen]
}

// AesGcmEncrypt AesGcmEncrypt
func AesGcmEncrypt(key, nonce, aad, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, nil
	}
	ciphertext := aesgcm.Seal(nil, nonce, data, aad)
	return ciphertext, nil
}

// AesGcmDecrypt AesGcmDecrypt
func AesGcmDecrypt(key, nonce, aad, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plain, err := aesgcm.Open(nil, nonce, data, aad)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

// GetNonce GetNonce
func GetNonce(data []byte, seq uint32) []byte {
	ret := make([]byte, len(data))
	copy(ret, data)
	seqBytes := baseutils.Int32ToBytes(seq)
	baseOffset := 8
	for index := 0; index < 4; index++ {
		ret[baseOffset+index] = ret[baseOffset+index] ^ byte(seqBytes[index])
	}
	return ret
}

// ECDSAVerifyData 校验服务端握手数据
func ECDSAVerifyData(message []byte, signature []byte) (bool, error) {
	block, _ := pem.Decode(verifyPubKey)
	publicStream, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	//接口转换成公钥
	publicKey := publicStream.(*ecdsa.PublicKey)

	// 反序列化ecdsaSignature
	ecdsaSignature := &EcdsaSignature{}
	_, err = asn1.Unmarshal(signature, ecdsaSignature)
	if err != nil {
		return false, err
	}
	flag := ecdsa.Verify(publicKey, Sha256(message), ecdsaSignature.R, ecdsaSignature.S)

	return flag, nil
}
