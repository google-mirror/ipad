package android

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

type PacketHeader struct {
	PacketCryptType byte
	Flag            uint16
	RetCode         uint32
	UICrypt         uint32
	Uin             uint32
	Cookies         []byte
	Data            []byte
}

func DoZlibCompress(src []byte) []byte {

	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func PackHybridEcdh(cmdid, cert, uin uint32, cookie, input []byte) []byte {

	inputlen := len(input)

	crc := CalcMsgCrc(input)
	pack := append([]byte{}, cookie...)
	pack = proto.EncodeVarint(uint64(cmdid))
	pack = append(pack, proto.EncodeVarint(uint64(inputlen))...)
	pack = append(pack, proto.EncodeVarint(uint64(inputlen))...)
	pack = append(pack, proto.EncodeVarint(uint64(cert))...)
	pack = append(pack, 2)
	pack = append(pack, 0)
	pack = append(pack, 0xfe)
	pack = append(pack, proto.EncodeVarint(uint64(crc))...)
	pack = append(pack, 0)

	headLen := len(pack) + 11
	headFlag := (12 << 12) | (len(cookie) << 8) | (headLen << 2) | 2

	var hybridpack = new(bytes.Buffer)
	hybridpack.WriteByte(0xbf)
	binary.Write(hybridpack, binary.LittleEndian, uint16(headFlag))
	binary.Write(hybridpack, binary.BigEndian, uint32(0x27000b32))
	binary.Write(hybridpack, binary.BigEndian, uint32(uin))
	hybridpack.Write(pack)
	hybridpack.Write(input)

	return hybridpack.Bytes()
}

func UnpackHybridEcdh(input []byte) *PacketHeader {

	var ph PacketHeader
	readHeader := bytes.NewReader(input)
	binary.Read(readHeader, binary.LittleEndian, &ph.PacketCryptType)
	binary.Read(readHeader, binary.LittleEndian, &ph.Flag)

	cookieLen := (ph.Flag >> 8) & 0x0f
	headerLen := (ph.Flag & 0xff) >> 2
	ph.Cookies = make([]byte, cookieLen)
	binary.Read(readHeader, binary.BigEndian, &ph.RetCode)
	binary.Read(readHeader, binary.BigEndian, &ph.UICrypt)
	binary.Read(readHeader, binary.LittleEndian, &ph.Cookies)
	ph.Data = input[headerLen:]

	return &ph
}

func Pack(cmdid, cert, algo, uin uint32, cookies, authecdhkey, sesskey, input []byte) []byte {

	inputlen := len(input)

	crc := CalcMsgCrc(input)
	sign := GenSignature(uin, authecdhkey, input)

	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, uin)
	b.Write(cookies)
	b.Write(proto.EncodeVarint(uint64(cmdid)))
	b.Write(proto.EncodeVarint(uint64(inputlen)))
	b.Write(proto.EncodeVarint(uint64(inputlen)))
	b.Write(proto.EncodeVarint(uint64(cert)))
	b.Write(proto.EncodeVarint(uint64(2)))
	b.Write(proto.EncodeVarint(uint64(sign)))
	b.Write([]byte{0xfe})
	b.Write(proto.EncodeVarint(uint64(crc)))
	b.Write([]byte{0x00})

	var encData []byte
	var compress uint32
	if algo == 5 {
		encData = EncryptAES(input, sesskey)
		compress = 2
	}

	subHead := b.Bytes()
	flag := (algo << 12) | (uint32(len(cookies)) << 8) | ((7 + uint32(len(subHead))) << 2) | compress

	bb := new(bytes.Buffer)
	bb.Write([]byte{0xbf})
	binary.Write(bb, binary.LittleEndian, uint16(flag))
	binary.Write(bb, binary.BigEndian, uint32(0x27000b32))
	bb.Write(subHead)
	bb.Write(encData)

	return bb.Bytes()
}

func Unpack(input, key []byte) *PacketHeader {

	var ph PacketHeader
	readHeader := bytes.NewReader(input)
	binary.Read(readHeader, binary.LittleEndian, &ph.PacketCryptType)
	binary.Read(readHeader, binary.LittleEndian, &ph.Flag)

	cookieLen := (ph.Flag >> 8) & 0x0f
	headerLen := (ph.Flag & 0xff) >> 2
	algo := ph.Flag >> 12
	comp := ph.Flag & 3

	ph.Cookies = make([]byte, cookieLen)
	binary.Read(readHeader, binary.BigEndian, &ph.RetCode)
	binary.Read(readHeader, binary.BigEndian, &ph.UICrypt)
	binary.Read(readHeader, binary.LittleEndian, &ph.Cookies)
	Data := input[headerLen:]

	var DecData []byte
	if algo == 5 {
		DecData = DecryptAES(Data, key)
	}

	if comp == 1 {
		ph.Data = DoZlibUnCompress(DecData)
	} else {
		ph.Data = DecData
	}

	return &ph
}
