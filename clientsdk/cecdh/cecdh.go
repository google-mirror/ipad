package cecdh

// import "C"
import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
)

/*//GenerateEccKey GenerateEccKey2
func GenerateEccKey() ([]byte, []byte) {
	pubKeyBuf := make([]byte, 57)
	priKeyBuf := make([]byte, 106)
	cPubKeyBuf := getCharPointer(pubKeyBuf)
	cPriKeyBuf := getCharPointer(priKeyBuf)
	//https://tools.ietf.org/html/rfc5480 Sect409k1.nid = 731
	C.GenerateECKey(713, cPubKeyBuf, cPriKeyBuf)
	return pubKeyBuf, priKeyBuf
}

//ComputerECCKeyMD5 ComputerECCKeyMD52
func ComputerECCKeyMD5(ecPubKey []byte, ecPriKey []byte) []byte {
	cEcPubKey := getCharPointer(ecPubKey)
	cEcPriKey := getCharPointer(ecPriKey)
	outAesKey := make([]byte, 16)
	cOutAesKey := getCharPointer(outAesKey)
	cPubKeyLen := C.int(len(ecPubKey))
	cPriKeyLen := C.int(len(ecPriKey))
	C.ComputerECCKeyMD5(cEcPubKey, cPubKeyLen, cEcPriKey, cPriKeyLen, cOutAesKey)

	return outAesKey
}
*/

//GenerateEccKey GenerateEccKey2
func GenerateEccKey() ([]byte, []byte) {
	newEllipticECDH := NewEllipticECDH(elliptic.P224())
	priKeyBuf, pubKeyBuf, err := newEllipticECDH.GenerateKey(rand.Reader)
	if err != nil {
		return []byte{}, []byte{}
	}

	return pubKeyBuf, priKeyBuf
}

//ComputerECCKeyMD5 ComputerECCKeyMD52
func ComputerECCKeyMD5(ecPubKey []byte, ecPriKey []byte) []byte {
	e := NewEllipticECDH(elliptic.P224())
	srvPubKey, ok := e.Unmarshal(ecPubKey)
	if !ok {
		return []byte{}
	}
	pwd, err := e.GenerateSharedSecret(ecPriKey, srvPubKey)
	if err != nil {
		return []byte{}
	}
	ctx := md5.New()
	_, _ = ctx.Write(pwd)
	return ctx.Sum(nil)
}
