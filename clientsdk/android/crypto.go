// crypto.go
package android

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"feiyu.com/wx/clientsdk/android/mmproto"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"

	proto "github.com/golang/protobuf/proto"
)

type HYBRID_STATUS int32

type HybridEcdhClient struct {
	hybridStatus HYBRID_STATUS

	clientHash hash.Hash
	serverHash hash.Hash

	clientStaticPub []byte
	clientEcdsaPub  []byte

	genClientPub  []byte
	genClientPriv []byte

	curve elliptic.Curve
}

func AesGcmEncrypt(key, nonce, input, additional []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		//fmt.Println("cipher init faile....")
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("gcm init faile....")
		return nil, err
	}

	result := aesgcm.Seal(nil, nonce, input, additional)

	return result, nil
}

func AesGcmDecrypt(key, nonce, input, additional []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println("cipher init faile....")
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("gcm init faile....")
		return nil, err
	}

	return aesgcm.Open(nil, nonce, input, additional)
}

func AesGcmEncryptWithCompress(key, nonce, input, additional []byte) ([]byte, error) {

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(input)
	w.Close()

	data, _ := AesGcmEncrypt(key, nonce, b.Bytes(), additional)

	encData := data[:len(data)-16]
	tag := data[len(data)-16:]

	totalData := []byte{}
	totalData = append(totalData, encData...)
	totalData = append(totalData, nonce...)
	totalData = append(totalData, tag...)
	return totalData, nil
}

func AesGcmDecryptWithUnCompress(key, input, additional []byte) ([]byte, error) {

	inputSize := len(input)

	nonce := make([]byte, 12)
	copy(nonce, input[inputSize-28:inputSize-16])

	tag := make([]byte, 16)
	copy(tag, input[inputSize-16:])

	cipherText := make([]byte, inputSize-28)
	copy(cipherText, input[:inputSize-28])
	cipherText = append(cipherText, tag...)

	result, _ := AesGcmDecrypt(key, nonce, cipherText, additional)

	b := bytes.NewReader(result)

	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	r.Close()

	return out.Bytes(), nil
}

func Ecdh(curve elliptic.Curve, pub, priv []byte) []byte {

	x, y := elliptic.Unmarshal(curve, pub)
	if x == nil {
		return nil
	}

	xShared, _ := curve.ScalarMult(x, y, priv)
	sharedKey := make([]byte, (curve.Params().BitSize+7)>>3)
	xBytes := xShared.Bytes()
	copy(sharedKey[len(sharedKey)-len(xBytes):], xBytes)

	dh := sha256.Sum256(sharedKey)
	return dh[:]
}

func (h *HybridEcdhClient) Init() {

	h.hybridStatus = HYBRID_ENC

	h.clientHash = sha256.New()
	h.serverHash = sha256.New()

	h.clientStaticPub, _ = hex.DecodeString("0495BC6E5C1331AD172D0F35B1792C3CE63F91572ABD2DD6DF6DAC2D70195C3F6627CCA60307305D8495A8C38B4416C75021E823B6C97DFFE79C14CB7C3AF8A586")
	h.clientEcdsaPub, _ = hex.DecodeString("2D2D2D2D2D424547494E205055424C4943204B45592D2D2D2D2D0A4D466B77457759484B6F5A497A6A3043415159494B6F5A497A6A3044415163445167414552497979694B33533950374854614B4C654750314B7A6243435139490A4C537845477861465645346A6E5A653646717777304A6877356D41716266574C4B364E69387075765356364371432B44324B65533373767059773D3D0A2D2D2D2D2D454E44205055424C4943204B45592D2D2D2D2D0A")

	h.curve = elliptic.P256()
}

func (h *HybridEcdhClient) Encrypt(input []byte) []byte {
	if h.hybridStatus != HYBRID_ENC {
		return nil
	}

	priv, x, y, error := elliptic.GenerateKey(h.curve, rand.Reader)
	if error != nil {
		return nil
	}
	h.genClientPriv = priv
	h.genClientPub = elliptic.Marshal(h.curve, x, y)

	ecdhKey := Ecdh(h.curve, h.clientStaticPub, h.genClientPriv)

	//hash1
	h1 := sha256.New()
	h1.Write([]byte("1"))
	h1.Write([]byte("415"))
	h1.Write(h.genClientPub)
	h1Sum := h1.Sum(nil)

	//Random
	random := make([]byte, 32)
	io.ReadFull(rand.Reader, random)

	nonce1 := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce1)
	gcm1, _ := AesGcmEncryptWithCompress(ecdhKey[0:0x18], nonce1, random, h1Sum)
	//hkdf
	salt, _ := hex.DecodeString("73656375726974792068646B6620657870616E64")
	hkdfKey := make([]byte, 56)
	hkdf.New(sha256.New, random, salt, h1Sum).Read(hkdfKey)

	//hash2
	h2 := sha256.New()
	h2.Write([]byte("1"))
	h2.Write([]byte("415"))
	h2.Write(h.genClientPub)
	h2.Write(gcm1)
	h2Sum := h2.Sum(nil)

	nonce2 := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce2)
	gcm2, _ := AesGcmEncryptWithCompress(hkdfKey[0:0x18], nonce2, input, h2Sum)

	var nid int32 = 415
	secKey := &mmproto.SecKey{
		Nid: &nid,
		Key: h.genClientPub,
	}

	var ver int32 = 1
	he := &mmproto.HybridEcdhReq{
		Version: &ver,
		SecKey:  secKey,
		Gcm1:    gcm1,
		Autokey: []byte{},
		Gcm2:    gcm2,
	}

	protoMsg, _ := proto.Marshal(he)

	// update client
	h.clientHash.Write(hkdfKey[0x18:0x38])
	h.clientHash.Write(input)

	// update server
	h.serverHash.Write(gcm2)

	h.hybridStatus = HYBRID_DEC

	return protoMsg
}

func (h *HybridEcdhClient) Decrypt(input []byte) []byte {

	if h.hybridStatus != HYBRID_DEC {
		return nil
	}

	var resp mmproto.HybridEcdhResp
	proto.Unmarshal(input, &resp)

	h.serverHash.Write(resp.GetGcm1())
	//	hServ := h.serverHash.Sum(nil)
	//	fmt.Printf("%x\n", hServ)

	//	ecdsa.Verify(h.clientEcdsaPub, resp.GetGcm2(), hServ)
	ecKey := Ecdh(h.curve, resp.GetSecKey().GetKey(), h.genClientPriv)

	h.clientHash.Write([]byte("415"))
	h.clientHash.Write(resp.GetSecKey().GetKey())
	h.clientHash.Write([]byte("1"))
	hCli := h.clientHash.Sum(nil)

	plain, _ := AesGcmDecryptWithUnCompress(ecKey[0:0x18], resp.GetGcm1(), hCli)
	return plain
}
