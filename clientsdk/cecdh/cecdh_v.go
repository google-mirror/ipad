package cecdh

/*
#cgo windows CFLAGS: -I./openssl/include -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -I./openssl/include -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -I./openssl/include -DCGO_OS_LINUX=1
#cgo windows LDFLAGS: -L./bin -lcrypto_windows -lstdc++
#cgo darwin LDFLAGS: -L./bin -lcrypto_darwin -lstdc++
#cgo linux LDFLAGS: -L./bin -lcrypto_linux -lstdc++

#include <stdio.h>
#include "internal/cryptlib.h"
#include <openssl/x509.h>
#include <openssl/ec.h>
#include <openssl/bn.h>
#include <openssl/cms.h>
#include <openssl/asn1t.h>
#include <openssl/md5.h>

int GenerateECKey(int nid, char* pubKey, char* priKey)
{
    char buf[68] = {0};

    EC_KEY* pKey = EC_KEY_new_by_curve_name(nid);
    if(pKey == NULL) { return 0;}
    const EC_GROUP* pGroup = EC_KEY_get0_group(pKey);
    if(pGroup == NULL) { EC_KEY_free(pKey); return 0;}

    int ret = EC_KEY_generate_key(pKey);
    if(ret == 0) { EC_KEY_free(pKey); return 0;}

    const EC_POINT* pPoint = EC_KEY_get0_public_key(pKey);
    if(pPoint == NULL) { EC_KEY_free(pKey); return 0;}

    size_t retSize = EC_POINT_point2oct(pGroup, pPoint, 4, (unsigned char *)buf, 64, 0);
    memcpy(pubKey, buf, retSize);

    unsigned char* outPriKey = NULL;
    retSize = i2d_ECPrivateKey(pKey, &outPriKey);
    memcpy(priKey, outPriKey, retSize);
    free(outPriKey);

    return retSize;
}

int ComputerECCKeyMD5(char* pub, int pubLen, char* pri, int priLen, char* eccKey)
{
    EC_KEY* pKey = NULL;
    d2i_ECPrivateKey(&pKey, (const unsigned char **)&pri, priLen);
    if(pKey == NULL) { return 0;}

    char outBuf[64] = {0};
    const EC_GROUP* pGroup = EC_KEY_get0_group(pKey);
    if(pGroup == NULL) { EC_KEY_free(pKey); return 0;}
    EC_POINT* pPoint = EC_POINT_new(pGroup);
    if(pPoint == NULL) { EC_KEY_free(pKey); return 0;}
    int ret = EC_POINT_oct2point(pGroup, pPoint, (const unsigned char *)pub, pubLen, NULL);
    if(ret != 1) { EC_KEY_free(pKey); return 0;}
    int retLen = ECDH_compute_key(outBuf, 64, pPoint, pKey, NULL);
    if(retLen <= 0) { return 0;}
    MD5((const unsigned char *)outBuf, retLen, (unsigned char *)eccKey);

    return retLen;
}
*/
/*
import "C"
import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"unsafe"
)

func getCharPointer(val []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&val[0]))
}

//GenerateEccKey GenerateEccKey2
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
}*/
