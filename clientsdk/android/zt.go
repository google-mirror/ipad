package android

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/golang/protobuf/proto"
)

// ZT SS
type ZT struct {
	ver       string
	initKey   []byte
	totalSize int32
	xorKey1   []byte
	key1      []byte
	xorKey2   []byte
	key2      []byte
	key3      []byte
}

// Init s
func (z *ZT) Init() {
	saeData, _ := base64.StdEncoding.DecodeString(SaeDat06)
	saePB := new(wechat.SaeInfoAndroid)
	proto.Unmarshal(saeData, saePB)
	z.ver = saePB.GetVer()
	z.initKey = saePB.GetInitKey()
	z.totalSize = saePB.GetTotalSize()
	z.xorKey1 = saePB.GetXorKey1()
	z.key1 = saePB.GetKey1()
	z.xorKey2 = saePB.GetXorKey2()
	z.key2 = saePB.GetKey2()
	z.key3 = saePB.GetKey3()
}

func (z *ZT) chooseKey(in, key []byte) []byte {

	var randKey [4][4][4]byte

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				randKey[k][j][i] = key[i*0x1000+j*0x400+int(in[i*4+j])*4+k]
			}
		}
	}

	var ret []byte
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				ret = append(ret, randKey[i][j][k])
			}
		}
	}

	return ret
}

func (z *ZT) chooseKey2Sub(a, b byte, key []byte) byte {
	var keySub1 = (a & 0xf0) | (b >> 4)
	if (keySub1 & 0x80) != 0 {
		keySub1 = key[keySub1&0x7f] >> 4
	} else {
		keySub1 = key[keySub1] & 0xf
	}

	var keySub2 = ((a & 0xf) << 4) | (b & 0xf)
	if (keySub2 & 0x80) != 0 {
		keySub2 = key[keySub2&0x7f+0x80] >> 4
	} else {
		keySub2 = key[keySub2+0x80] & 0x0f
	}

	return ((keySub1 & 0xf) << 4) | (keySub2 & 0x0f)
}

func (z *ZT) chooseKey2(keyA, keyB []byte) []byte {

	result := make([]byte, 16)

	for k := 0; k < 4; k++ {
		for j := 0; j < 4; j++ {

			result[4*k+j] = keyA[16*k+4*j+3]
			offset := 0
			for i := 2; i != -1; i-- {

				result[4*k+j] = z.chooseKey2Sub(keyA[16*k+j*4+i], result[4*k+j], keyB[(k*0xc00+j*0x300+0x200-offset*0x100):])
				offset++
			}

		}
	}

	return result
}

func (z *ZT) chooseKey3(in, key []byte) []byte {

	result := make([]byte, 16)
	for k := 0; k < 4; k++ {
		for j := 0; j < 4; j++ {
			result[k*4+j] = key[uint(in[k*4+j])+uint(j)*0x100+uint(k)*0x400]
		}
	}

	return result
}

func (z *ZT) shiftKey(in [4][4]byte) []byte {

	var ret [4][4]byte
	ret[0][0] = in[0][0]
	ret[0][1] = in[0][1]
	ret[0][2] = in[0][2]
	ret[0][3] = in[0][3]

	ret[1][0] = in[1][1]
	ret[1][1] = in[1][2]
	ret[1][2] = in[1][3]
	ret[1][3] = in[1][0]

	ret[2][0] = in[2][2]
	ret[2][1] = in[2][3]
	ret[2][2] = in[2][0]
	ret[2][3] = in[2][1]

	ret[3][0] = in[3][3]
	ret[3][1] = in[3][0]
	ret[3][2] = in[3][1]
	ret[3][3] = in[3][2]

	var result []byte

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result = append(result, ret[i][j])
		}
	}

	return result
}

func (z *ZT) reAssemble(in []byte) []byte {

	result := make([]byte, 16)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i*4+j] = in[j*4+i]
		}
	}
	return result
}

func (z *ZT) byte2Array(in []byte) [4][4]byte {
	var r [4][4]byte
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			r[i][j] = in[i*4+j]
		}
	}

	return r
}

func zipCCData(in []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(in)
	w.Close()

	return b.Bytes()
}

// WBAesEncrypt s
func (z *ZT) WBAesEncrypt(pb []byte) []byte {

	in := zipCCData(pb)

	size := len(in)
	pad := 16 - (size % 16)
	for i := 0; i < pad; i++ {
		in = append(in, byte(pad))
	}

	initKey := z.initKey
	totalRound := len(in) / 16

	result := make([]byte, len(in))

	for i := 0; i < totalRound; i++ {
		//step1
		var step1 [4][4]byte
		for j := 0; j < 16; j++ {
			step1[j/4][j%4] = in[i*16+j] ^ initKey[j]
		}

		//step2
		var step2 [4][4]byte
		for k := 0; k < 4; k++ {
			for m := 0; m < 4; m++ {
				step2[k][m] = step1[m][k]
			}
		}

		//step3
		for l := 0; l < 9; l++ {
			step3 := z.shiftKey(step2)

			step4 := z.chooseKey(step3, z.key1[0x4000*l:])

			step5 := z.chooseKey2(step4, z.key2[0x3000*l:])

			step2 = z.byte2Array(step5)
		}

		step6 := z.shiftKey(step2)
		step7 := z.chooseKey3(step6, z.key3)
		step8 := z.reAssemble(step7)

		copy(result[i*16:], step8)
		initKey = step8
	}
	return result
}
