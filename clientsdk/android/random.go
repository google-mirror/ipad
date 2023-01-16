package android

import (
	"bytes"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/lunny/log"
	"hash/adler32"
	"io"
	"time"
)

func RandBytes(size int) []byte {

	r := make([]byte, size)
	n, err := io.ReadFull(rand.Reader, r)
	if err != nil || n != size {
		log.Println("Gen rand faile")
		return nil
	}
	return r
}

func Gen713Key() ([]byte, []byte) {

	priv, x, y, _ := elliptic.GenerateKey(elliptic.P224(), rand.Reader)

	pub := elliptic.Marshal(elliptic.P224(), x, y)

	return pub, priv
}

func Do713Ecdh(pub, priv []byte) []byte {
	curve := elliptic.P224()
	x, y := elliptic.Unmarshal(curve, pub)
	if x == nil {
		return nil
	}

	xShared, _ := curve.ScalarMult(x, y, priv)
	sharedKey := make([]byte, (curve.Params().BitSize+7)>>3)
	xBytes := xShared.Bytes()
	copy(sharedKey[len(sharedKey)-len(xBytes):], xBytes)

	dh := md5.Sum(sharedKey)
	return dh[:]
}

func GenSignature(uiCryptin uint32, salt, data []byte) uint32 {

	var b1 bytes.Buffer
	binary.Write(&b1, binary.BigEndian, uiCryptin)

	h1 := md5.New()
	h1.Write(b1.Bytes())
	h1.Write(salt)
	sum1 := h1.Sum(nil)

	dataSize := len(data)
	var b2 bytes.Buffer
	binary.Write(&b2, binary.BigEndian, dataSize)

	h2 := md5.New()
	h2.Write(b2.Bytes())
	h2.Write(salt)
	h2.Write(sum1)
	sum2 := h2.Sum(nil)

	a := adler32.New()
	a.Write(nil)
	a.Write(sum2)
	a.Write(data)

	return a.Sum32()
}

func GenUUID() string {

	randomData := make([]byte, 0x100)
	io.ReadFull(rand.Reader, randomData)

	h := md5.New()
	h.Write(randomData)

	nanoTime := time.Now().UnixNano()
	secTime := time.Now().Unix()

	h.Write([]byte(fmt.Sprintf("%v_%v", secTime, nanoTime)))
	sum := h.Sum(nil)

	return fmt.Sprintf("%x-%x-%x-%x-%x", sum[:4], sum[4:6], sum[6:8], sum[8:10], sum[10:])
}
