package android

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"feiyu.com/wx/clientsdk/android/mmproto"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"
)

const (
	HYBRID_ENC HYBRID_STATUS = 0
	HYBRID_DEC HYBRID_STATUS = 1
)

type Client struct {
	PubKey     []byte
	Privkey    []byte
	InitPubKey []byte
	Externkey  []byte

	Version    int
	DeviceType string

	clientHash hash.Hash
	serverHash hash.Hash

	curve elliptic.Curve

	Status HYBRID_STATUS
}

func (h *Client) Init(Model string, version int, deviceType string) {
	h.curve = elliptic.P256()
	h.clientHash = sha256.New()
	h.serverHash = sha256.New()
	if Model == "IOS" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = version
		h.DeviceType = deviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}
	if Model == "Android" {
		h.Status = HYBRID_ENC
		h.Version = version
		h.DeviceType = deviceType
		h.InitPubKey, _ = hex.DecodeString("0495BC6E5C1331AD172D0F35B1792C3CE63F91572ABD2DD6DF6DAC2D70195C3F6627CCA60307305D8495A8C38B4416C75021E823B6C97DFFE79C14CB7C3AF8A586")
	}

}

func GetECDH415Key() (privKey []byte, pubKey []byte) {
	privKey = nil
	pubKey = nil
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := &priv.PublicKey
	pubKey = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	privKey = priv.D.Bytes()
	return
}

func (h *Client) HybridEcdhPackAndroidEn(cmdid, cert, uin uint32, cookie, Data []byte) []byte {
	EnData := h.encryptAndroid(Data)
	inputlen := len(EnData)
	pack := append([]byte{}, cookie...)
	pack = proto.EncodeVarint(uint64(cmdid))
	pack = append(pack, proto.EncodeVarint(uint64(inputlen))...)
	pack = append(pack, proto.EncodeVarint(uint64(inputlen))...)
	pack = append(pack, proto.EncodeVarint(uint64(cert))...)
	pack = append(pack, 2)
	pack = append(pack, 0)
	pack = append(pack, 0xfe)
	pack = append(pack, proto.EncodeVarint(uint64(baseutils.CalcMsgCrcForData_807(EnData)))...)
	pack = append(pack, 0)
	headLen := len(pack) + 11
	headFlag := (12 << 12) | (len(cookie) << 8) | (headLen << 2) | 2
	var hybridpack = new(bytes.Buffer)
	hybridpack.WriteByte(0xbf)
	binary.Write(hybridpack, binary.LittleEndian, uint16(headFlag))
	binary.Write(hybridpack, binary.BigEndian, uint32(h.Version))
	binary.Write(hybridpack, binary.BigEndian, uint32(uin))
	hybridpack.Write(pack)
	hybridpack.Write(EnData)
	return hybridpack.Bytes()
}

func (h *Client) HybridEcdhPackAndroidUn(Data []byte) *PacketHeader {
	var ph PacketHeader
	readHeader := bytes.NewReader(Data)
	binary.Read(readHeader, binary.LittleEndian, &ph.PacketCryptType)
	binary.Read(readHeader, binary.LittleEndian, &ph.Flag)
	cookieLen := (ph.Flag >> 8) & 0x0f
	headerLen := (ph.Flag & 0xff) >> 2
	ph.Cookies = make([]byte, cookieLen)
	binary.Read(readHeader, binary.BigEndian, &ph.RetCode)
	binary.Read(readHeader, binary.BigEndian, &ph.UICrypt)
	binary.Read(readHeader, binary.LittleEndian, &ph.Cookies)
	ph.Data = h.decryptAndroid(Data[headerLen:])
	return &ph
}

func (h *Client) decryptAndroid(input []byte) []byte {

	if h.Status != HYBRID_DEC {
		return nil
	}

	var resp mmproto.HybridEcdhResp
	proto.Unmarshal(input, &resp)

	h.serverHash.Write(resp.GetGcm1())
	//	hServ := h.serverHash.Sum(nil)
	//	fmt.Printf("%x\n", hServ)

	//	ecdsa.Verify(h.clientEcdsaPub, resp.GetGcm2(), hServ)
	ecKey := Ecdh(h.curve, resp.GetSecKey().GetKey(), h.Privkey)

	h.clientHash.Write([]byte("415"))
	h.clientHash.Write(resp.GetSecKey().GetKey())
	h.clientHash.Write([]byte("1"))
	hCli := h.clientHash.Sum(nil)
	plain, _ := AesGcmDecryptWithUnCompress(ecKey[0:0x18], resp.GetGcm1(), hCli)
	return plain
}

func (h *Client) encryptAndroid(input []byte) []byte {
	if h.Status != HYBRID_ENC {
		return nil
	}

	priv, x, y, error := elliptic.GenerateKey(h.curve, rand.Reader)
	if error != nil {
		return nil
	}
	h.Privkey = priv
	h.PubKey = elliptic.Marshal(h.curve, x, y)

	ecdhKey := Ecdh(h.curve, h.InitPubKey, h.Privkey)

	//hash1
	h1 := sha256.New()
	h1.Write([]byte("1"))
	h1.Write([]byte("415"))
	h1.Write(h.PubKey)
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
	h2.Write(h.PubKey)
	h2.Write(gcm1)
	h2Sum := h2.Sum(nil)

	nonce2 := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce2)
	gcm2, _ := AesGcmEncryptWithCompress(hkdfKey[0:0x18], nonce2, input, h2Sum)
	var nid int32 = 415
	secKey := &mmproto.SecKey{
		Nid: &nid,
		Key: h.PubKey,
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

	h.Status = HYBRID_DEC

	return protoMsg
}

func RQT(data []byte) int {

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

func HybridHkdfExpand(prikey []byte, salt []byte, info []byte, outLen int) []byte {
	h := hmac.New(sha256.New, prikey)
	h.Write(salt)
	T := h.Sum(nil)
	return HkdfExpand(sha256.New, T, info, outLen)
}

func Hkdf_Expand(h func() hash.Hash, prk, info []byte, outLen int) []byte {
	out := []byte{}
	T := []byte{}
	i := byte(1)
	for len(out) < outLen {
		block := append(T, info...)
		block = append(block, i)

		h := hmac.New(h, prk)
		h.Write(block)

		T = h.Sum(nil)
		out = append(out, T...)
		i++
	}
	return out[:outLen]
}

func HkdfExpand(h func() hash.Hash, prk, info []byte, outLen int) []byte {
	out := []byte{}
	T := []byte{}
	i := byte(1)
	for len(out) < outLen {
		block := append(T, info...)
		block = append(block, i)

		h := hmac.New(h, prk)
		h.Write(block)

		T = h.Sum(nil)
		out = append(out, T...)
		i++
	}
	return out[:outLen]
}

// ios组包
func (h *Client) HybridEcdhPackIosEn(Cgi, Uin uint32, Cookies, Data []byte) []byte {
	header := new(bytes.Buffer)
	header.Write([]byte{0xbf})
	header.Write([]byte{0x02}) //加密模式占坑,默认不压缩走12

	encryptdata := h.encryptoIOS(Data)

	cookielen := len(Cookies)
	header.Write([]byte{byte((12 << 4) + cookielen)})
	binary.Write(header, binary.BigEndian, int32(h.Version))
	if Uin != 0 {
		binary.Write(header, binary.BigEndian, int32(Uin))
	} else {
		header.Write([]byte{0x00, 0x00, 0x00, 0x00})
	}

	if len(Cookies) == 0xF {
		header.Write(Cookies)
	}

	header.Write(proto.EncodeVarint(uint64(Cgi)))
	header.Write(proto.EncodeVarint(uint64(len(encryptdata))))
	header.Write(proto.EncodeVarint(uint64(len(encryptdata))))
	header.Write([]byte{0x90, 0x4E, 0x0D, 0x00, 0xFF})
	header.Write(proto.EncodeVarint(uint64(RqtIOS(encryptdata))))
	header.Write([]byte{0x00})
	lens := len(header.Bytes())<<2 + 2
	header.Bytes()[1] = byte(lens)
	header.Write(encryptdata)
	return header.Bytes()
}

func (h *Client) HybridEcdhPackIosUn(Data []byte) *PacketHeader {
	var ph PacketHeader
	var body []byte
	var nCur int64
	var bfbit byte
	srcreader := bytes.NewReader(Data)
	binary.Read(srcreader, binary.BigEndian, &bfbit)
	if bfbit == byte(0xbf) {
		nCur += 1
	}
	nLenHeader := Data[nCur] >> 2
	nCur += 1
	nLenCookie := Data[nCur] & 0xf
	nCur += 1
	nCur += 4
	srcreader.Seek(nCur, io.SeekStart)
	binary.Read(srcreader, binary.BigEndian, &ph.Uin)
	nCur += 4
	cookie_temp := Data[nCur : nCur+int64(nLenCookie)]
	ph.Cookies = cookie_temp
	nCur += int64(nLenCookie)
	cgidata := Data[nCur:]
	_, nSize := proto.DecodeVarint(cgidata)
	nCur += int64(nSize)
	LenProtobufData := Data[nCur:]
	_, nLenProtobuf := proto.DecodeVarint(LenProtobufData)
	nCur += int64(nLenProtobuf)
	body = Data[nLenHeader:]
	protobufdata := h.decryptoIOS(body)
	ph.Data = protobufdata
	return &ph
}

func (h *Client) decryptoIOS(Data []byte) []byte {
	HybridEcdhResponse := &wechat.HybridEcdhResponse{}
	err := proto.Unmarshal(Data, HybridEcdhResponse)
	if err != nil {
		return nil
	}
	decrptecdhkey := DoECDH415Key(h.Privkey, HybridEcdhResponse.GetSecECDHKey().GetBuffer())
	m := sha256.New()
	m.Write(decrptecdhkey)
	decrptecdhkey = m.Sum(nil)
	h.serverHash.Write([]byte("415"))
	h.serverHash.Write(HybridEcdhResponse.GetSecECDHKey().GetBuffer())
	h.serverHash.Write([]byte("1"))
	mServerpubhashFinal_digest := h.serverHash.Sum(nil)

	outdata := AesGcmDecryptWithcompressZlib(decrptecdhkey[:24], HybridEcdhResponse.GetDecryptdata(), mServerpubhashFinal_digest)
	return outdata
}

func AesGcmDecryptWithcompressZlib(key []byte, ciphertext []byte, aad []byte) []byte {
	ciphertextinput := ciphertext[:len(ciphertext)-0x1c]
	endatanonce := ciphertext[len(ciphertext)-0x1c : len(ciphertext)-0x10]
	data := new(bytes.Buffer)
	data.Write(ciphertextinput)
	data.Write(ciphertext[len(ciphertext)-0x10 : len(ciphertext)])
	decrypt_data := NewAES_GCMDecrypter(key, data.Bytes(), endatanonce, aad)
	if len(decrypt_data) > 0 {
		return DoZlibUnCompress(decrypt_data)
	} else {
		return []byte{}
	}
}

func RqtIOS(srcdata []byte) int {
	h := md5.New()
	h.Write(srcdata)
	md5sign := hex.EncodeToString(h.Sum(nil))
	key, _ := hex.DecodeString("6a664d5d537c253f736e48273a295e4f")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(md5sign))
	my_sign := string(mac.Sum(nil))
	randvalue := 1
	index := 0
	temp0 := 0
	temp1 := 0
	temp2 := 0
	for index = 0; index+2 < 20; index++ {
		temp0 = (temp0&0xff)*0x83 + int(my_sign[index])
		temp1 = (temp1&0xff)*0x83 + int(my_sign[index+1])
		temp2 = (temp2&0xff)*0x83 + int(my_sign[index+2])

	}
	result := (temp2<<16)&0x7f0000 | temp0&0x7f | (randvalue&0x1f|0x20)<<24 | ((temp1 & 0x7f) << 8)
	return result

}

func (h *Client) encryptoIOS(Data []byte) []byte {
	ecdhkey := DoECDH415Key(h.Privkey, h.InitPubKey)
	m := sha256.New()
	m.Write(ecdhkey)
	ecdhkey = m.Sum(nil)
	mClientpubhash := sha256.New()
	mClientpubhash.Write([]byte("1"))
	mClientpubhash.Write([]byte("415"))
	mClientpubhash.Write(h.PubKey)
	mClientpubhash_digest := mClientpubhash.Sum(nil)

	mRandomEncryptKey := make([]byte, 32)
	io.ReadFull(rand.Reader, mRandomEncryptKey)
	mNonce := make([]byte, 12)
	io.ReadFull(rand.Reader, mNonce)

	mEncryptdata := AesGcmEncryptWithCompressZlib(ecdhkey[:24], mRandomEncryptKey, mNonce, mClientpubhash_digest)
	var mExternEncryptdata []byte
	if len(h.Externkey) == 0x20 {
		mExternEncryptdata = AesGcmEncryptWithCompressZlib(h.Externkey[:24], mRandomEncryptKey, mNonce, mClientpubhash_digest)
	}
	hkdfexpand_security_key := HybridHkdfExpand([]byte("security hdkf expand"), mRandomEncryptKey, mClientpubhash_digest, 56)

	mClientpubhashFinal := sha256.New()
	mClientpubhashFinal.Write([]byte("1"))
	mClientpubhashFinal.Write([]byte("415"))
	mClientpubhashFinal.Write(h.PubKey)
	mClientpubhashFinal.Write(mEncryptdata)
	mClientpubhashFinal.Write(mExternEncryptdata)
	mClientpubhashFinal_digest := mClientpubhashFinal.Sum(nil)

	mEncryptdataFinal := AesGcmEncryptWithCompressZlib(hkdfexpand_security_key[:24], Data, mNonce, mClientpubhashFinal_digest)

	h.clientHash.Write(mEncryptdataFinal)

	h.serverHash.Write(hkdfexpand_security_key[24:56])
	h.serverHash.Write(Data)

	HybridEcdhRequest := &wechat.HybridEcdhRequest{
		Type: proto.Int32(1),
		SecECDHKey: &wechat.BufferT{
			ILen:   proto.Uint32(415),
			Buffer: h.PubKey,
		},
		Randomkeydata:       mEncryptdata,
		Randomkeyextenddata: mExternEncryptdata,
		Encyptdata:          mEncryptdataFinal,
	}
	reqdata, _ := proto.Marshal(HybridEcdhRequest)
	return reqdata
}

func DoECDH415Key(privD, pubData []byte) []byte {
	X, Y := elliptic.Unmarshal(elliptic.P256(), pubData)
	if X == nil || Y == nil {
		return nil
	}
	x, _ := elliptic.P256().ScalarMult(X, Y, privD)
	return x.Bytes()
}

func AesGcmEncryptWithCompressZlib(key []byte, plaintext []byte, nonce []byte, aad []byte) []byte {
	compressData := DoZlibCompress(plaintext)
	//nonce := []byte(randSeq(12)) //获取随机密钥
	encrypt_data := NewAES_GCMEncrypter(key, compressData, nonce, aad)
	outdata := encrypt_data[:len(encrypt_data)-16]
	retdata := new(bytes.Buffer)
	retdata.Write(outdata)
	retdata.Write(nonce)
	retdata.Write(encrypt_data[len(encrypt_data)-16:])
	return retdata.Bytes()
}

func NewAES_GCMEncrypter(key []byte, plaintext []byte, nonce []byte, aad []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, aad)
	return ciphertext
}

func NewAES_GCMDecrypter(key []byte, ciphertext []byte, nonce []byte, aad []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil
	}
	return plaintext
}
