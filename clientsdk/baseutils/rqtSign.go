package baseutils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
)

var (
	t2 = [16]uint32{
		0x8C6C30E2, 0x1A30534F, 0x3A6956B1, 0x3D24767F,
		0xECC353E9, 0x50A19FE4, 0xE593AFE8, 0x93CD8473,
		0x92634626, 0xED0D416D, 0x4532D90D, 0x6F9DD7FA,
		0xF0CAB0B9, 0x6151E06B, 0x9FBA1597, 0xF597583E,
	}

	t1 = [16]uint32{
		0xBF5F03D1, 0x2903607C, 0x095A6582, 0x0E17454C,
		0xDFF060DA, 0x6392ACD7, 0xD6A09CDB, 0xA0FEB740,
		0xA1507515, 0xDE3E725E, 0x7601EA3E, 0x5CAEE4C9,
		0xC3F9838A, 0x5262D358, 0xAC8926A4, 0xC6A46B0D,
	}

	f1 = [64]uint32{
		0xA4C5F56C, 0x184B6267, 0x7ABB2204, 0x71E2D7D9, 0x1BAD3688, 0x23505A37, 0x9C672518, 0xAB6A91DB, 0x1530628E, 0x1ED1DA0F,
		0x0030C700, 0x7FBF0256, 0xBA92CC3E, 0xA58C197A, 0x0391BBB6, 0x536C4D6D, 0xA5015881, 0xAF449DFF, 0xBA30096B, 0xA0F1EEB5,
		0x25945EF0, 0x07C2DF83, 0xC78D76D3, 0x9553D1B3, 0x7E586FCE, 0x44CC0D68, 0xD7C7E552, 0x7A610BBD, 0x15A562D6, 0x05704D9F,
		0x62486C77, 0x7E632B3E, 0x677304B2, 0xC1BC402F, 0x0B13F16E, 0x20557559, 0x9377B4E4, 0xBA0AC06D, 0x6472A718, 0x6D3DAC04,
		0xC5FA3D28, 0x99E257C1, 0xFEF3D299, 0xAED264EF, 0xA2CD9FF6, 0x8E4250B1, 0x0A2C44D6, 0x5C1E3658, 0xAB339562, 0x1F8C79D6,
		0xD52B9487, 0x43A4FFDC, 0x5E35FD97, 0x095150C1, 0xC99FF457, 0xBFEF4DC7, 0x91B8AF9F, 0x587E0902, 0x2ACF386A, 0x5CB5A76C,
		0x8413F0AB, 0x0BF6BAEC, 0x2CDCA4DA, 0xBA735FDC,
	}

	f2 = [64]uint32{
		0x9ADDDA57, 0x47E0F9C0, 0x18D42EE5, 0xBF710F92, 0x426B07B3, 0x21F0A33E, 0x5284D73B, 0x38F2FE90, 0x798C7519, 0x9309273F,
		0x009B3210, 0xD5F036FB, 0xC0E40B74, 0xDA9A09D5, 0xEE0B9027, 0x7416E86C, 0x6D5307E0, 0x3D6650C2, 0x77BDDB1F, 0xB2531ECB,
		0xB35CBE6B, 0xF71E8563, 0x02A10281, 0xEFF1F827, 0x38C0D810, 0x074FBD7B, 0x5A384A4E, 0x294EBD9A, 0x71B7355D, 0xF70636F1,
		0xA407BF30, 0xF7C06247, 0x51688CDA, 0x86149715, 0xCB7CD5AA, 0x23F0F476, 0xE0642C8F, 0x4A9731DE, 0x13212825, 0x331CDFBF,
		0xB8C7B03E, 0xB214E4DD, 0x09FE275B, 0x3CCFB3EC, 0xB615AEEF, 0x83EB28BB, 0x77BF8EAB, 0x0882B107, 0x0D1692A7, 0x837E3B5B,
		0x5BBACD15, 0x935515A8, 0x9CDDC643, 0x7E47C0BA, 0x87BB3B22, 0x8B5A33E9, 0xD9EB84D3, 0x74EE6459, 0x03DE676A, 0xA3B330C0,
		0xFB210470, 0x1AD12206, 0x2461D9AA, 0xA244F1CB,
	}
)

func ROR(src, bits uint32) uint32 {
	return (src >> bits) | (src << (32 - bits))
}

func genT1(s [16]uint32) [64]uint32 {

	var r [64]uint32
	for i := 0; i < 16; i++ {
		r[i] = ROR(s[i], 8)
	}

	for i := 16; i < 64; i++ {
		R3 := r[14+i-16]
		R5 := ROR(R3, 5)
		R6 := ROR(R3, 9)

		R7 := r[1+i-16]
		R8 := r[0+i-16]
		R9 := r[9+i-16]

		R3 >>= 2
		R3 ^= R5

		R5 = ROR(R7, 3)
		RA := ROR(R7, 11)
		R7 >>= 5
		R7 ^= RA
		R5 ^= R7
		R3 ^= R6

		R4 := R9
		R6 = R8
		R4 += R6
		R3 += R4
		R3 += R5

		r[i] = R3
	}

	return r
}

// RQTXHASH struct
type RQTXHASH struct {
	A1, A2, A3, A4, A5, A6, A7, A8, A9, A10 uint32

	mt [64]uint32
	ft [64]uint32

	state [8]uint32
}

func (r *RQTXHASH) init() {

	r.state[0] = 0xA195712E
	r.state[1] = 0x55C746B4
	r.state[2] = 0x40FCCA21
	r.state[3] = 0x5F96BC26
	r.state[4] = 0xD18A047B
	r.state[5] = 0xC54DC1C7
	r.state[6] = 0x9B113EE0
	r.state[7] = 0x52A98CF7
}

func (r *RQTXHASH) update(input [16]uint32) {

	r.A1 = r.state[0]
	r.A2 = r.state[1]
	r.A3 = r.state[2]
	r.A4 = r.state[3]
	r.A5 = r.state[4]
	r.A6 = r.state[5]
	r.A7 = r.state[6]
	r.A8 = r.state[7]

	r.mt = genT1(input)

	for i := 0; i < 64; i++ {

		if i%3 == 0 {
			r.ft = f1
		} else {
			r.ft = f2
		}

		R4 := ROR(r.A5, 7)
		R6 := ROR(r.A5, 13)

		R2 := R4 ^ R6

		R8 := r.A5 & r.A6
		R20 := r.A5 ^ 0xffffffff
		R20 &= r.A7

		R4 = ROR(r.A5, 17)
		R2 ^= R4
		R4 = R20 | R8
		R4 += r.A8

		R4 += r.mt[i]
		R2 += R4

		R1 := R2 + r.ft[i]

		R2 = ROR(r.A1, 5)
		R4 = ROR(r.A1, 11)
		R2 ^= R4

		R4 = ROR(r.A1, 19)
		R6 = r.A4 + R1

		r.A9 = R6

		R2 ^= R4
		R4 = r.A3 ^ r.A2
		R4 &= r.A1

		R6 = r.A3 & r.A2
		R4 ^= R6
		R2 += R4
		R1 += R2

		r.A10 = R1

		T1 := r.A1
		r.A1 = r.A10

		T2 := r.A2
		r.A2 = T1

		r.A8 = r.A7
		r.A7 = r.A6
		r.A6 = r.A5
		r.A5 = r.A9
		r.A4 = r.A3
		r.A3 = T2
	}

	r.state[0] += r.A8
	r.state[1] += r.A7
	r.state[2] += r.A6
	r.state[3] += r.A5
	r.state[4] += r.A4
	r.state[5] += r.A3
	r.state[6] += r.A2
	r.state[7] += r.A1
}

func (r *RQTXHASH) getState() []byte {

	rc := new(bytes.Buffer)
	binary.Write(rc, binary.LittleEndian, r.A1)
	binary.Write(rc, binary.LittleEndian, r.A2)
	binary.Write(rc, binary.LittleEndian, r.A3)
	binary.Write(rc, binary.LittleEndian, r.A4)
	binary.Write(rc, binary.LittleEndian, r.A5)
	binary.Write(rc, binary.LittleEndian, r.A6)
	binary.Write(rc, binary.LittleEndian, r.A7)
	binary.Write(rc, binary.LittleEndian, r.A8)
	return rc.Bytes()
}

func (r *RQTXHASH) final() []byte {

	rc := new(bytes.Buffer)
	binary.Write(rc, binary.BigEndian, r.state[0])
	binary.Write(rc, binary.BigEndian, r.state[1])
	binary.Write(rc, binary.BigEndian, r.state[2])
	binary.Write(rc, binary.BigEndian, r.state[3])
	binary.Write(rc, binary.BigEndian, r.state[4])
	binary.Write(rc, binary.BigEndian, r.state[5])
	binary.Write(rc, binary.BigEndian, r.state[6])
	binary.Write(rc, binary.BigEndian, r.state[7])
	return rc.Bytes()
}

//

func trans(input []byte) []byte {
	var BB [16]uint32
	salt := [4]uint32{0x7BA2E3FC, 0xE381CF0C, 0x0C1C991C, 0xACC7E8E4}
	for i := 0; i < 16; i++ {
		BB[i] = uint32(input[4*i+2]) | (uint32(input[4*i+1]) << 8) | (uint32(input[4*i]) << 16) | (uint32(input[4*i+3]) << 24)
	}

	R8 := salt[0] + BB[0] + (salt[1] ^ salt[2] ^ salt[3]) + 0xa4bdc6e9
	R8 = ROR(R8, 7)
	R8 += salt[1]

	R9 := salt[1] ^ salt[2]
	R9 ^= R8

	RA := salt[3] + BB[1]
	R9 += RA
	R9 += 0x554e80f9
	R9 = ROR(R9, 12)
	R9 += R8

	RB := salt[2] + BB[2]
	RC := R9 ^ R8
	RD := RC ^ salt[1]
	RB += RD
	RB += 0xd3d89f33
	RB = ROR(RB, 17)
	RB += R9

	RC ^= RB
	RD = salt[1] + BB[3]
	RC += RD
	RC += 0x15a2af88
	RC = ROR(RC, 22)

	R8 += BB[4]
	RC += RB
	RE := RC ^ RB
	RF := RE ^ R9
	R8 += RF
	R8 += 0x924dac27
	R8 = ROR(R8, 7)

	R9 += BB[5]
	R8 += RC
	RE ^= R8
	R9 += RE
	R9 += 0x7e563339
	R9 = ROR(R9, 12)

	RB += BB[6]
	R9 += R8
	RE = R9 ^ R8
	RF = RE ^ RC
	RB += RF
	RB += 0xf6bae427
	RB = ROR(RB, 17)

	RC += BB[7]
	RB += R9
	RE ^= RB
	RC += RE
	RC += 0x12132474
	RC = ROR(RC, 22)

	R8 += BB[8]
	RC += RB
	RE = RC ^ RB
	RF = RE ^ R9
	R8 += RF
	R8 += 0x9436b72e
	R8 = ROR(R8, 7)

	R9 += BB[9]
	R8 += RC
	RE ^= R8
	R9 += RE
	R9 += 0x8B366D12
	R9 = ROR(R9, 12)

	RB += BB[10]
	R9 += R8
	RE = R8 ^ RC
	RE ^= R9
	RB += RE
	RB += 0x23e22dcc
	RB = ROR(RB, 17)

	RC += BB[11]
	RB += R9
	RE = R9 ^ R8
	RE ^= RB
	RC += RE
	RC += 0x69ab94c8
	RC = ROR(RC, 22)

	R8 += BB[12]
	RC += RB
	RE = RB ^ R9
	RE ^= RC
	R8 += RE
	R8 += 0xcfdb953b
	R8 = ROR(R8, 7)

	R9 += BB[13]
	R8 += RC
	RE = RC ^ RB
	RE ^= R8
	R9 += RE
	R9 += 0xe3c702cf
	R9 = ROR(R9, 12)

	RA = RB + BB[14]
	R9 += R8
	RB = R8 ^ RC
	RB ^= R9
	RA += RB
	RA += 0x48233422
	RA = ROR(RA, 17)

	RB = RC + BB[15]
	RA += R9
	RC = R9 ^ R8
	RC ^= RA
	RB += RC
	RB += 0xaad28857
	RB = ROR(RB, 22)

	//round 2
	RB += RA
	RC = RB ^ 0xFFFFFFFF
	RC &= RA
	RD = RB & R9
	RC |= RD
	R8 += BB[1]
	R8 += RC
	R8 += 0x8de587b3
	R8 = ROR(R8, 5)

	R8 += RB
	RC = R8 ^ 0xFFFFFFFF
	RC &= RB
	RD = R8 & RA
	RC |= RD
	R9 += BB[6]
	R9 += RC
	R9 += 0x8f209322
	R9 = ROR(R9, 9)

	R9 += R8
	RD = R9 ^ 0xFFFFFFFF
	RD &= R8
	RE = R9 & RB
	RD |= RE
	RA += BB[11]
	RA += RD
	RA += 0xaed47066
	RA = ROR(RA, 14)
	RA += R9

	RE = ^RA
	RE &= R9
	RF = RA & R8
	RE |= RF
	RB += BB[0]
	RB += RE
	RB += 0x5e211c2c
	RB = ROR(RB, 20)
	RB += RA

	RF = ^RB
	RF &= RA
	R10 := RB & R9
	RF |= R10
	R8 += BB[5]
	R8 += RF
	R8 += 0xac5e94b2
	R8 = ROR(R8, 5)
	R8 += RB

	RF = ^R8
	RF &= RB
	R10 = R8 & RA
	RF |= R10
	R9 += BB[10]
	R9 += RF
	R9 += 0xfa747ca2
	R9 = ROR(R9, 9)
	R9 += R8

	RF = ^R9
	RF &= R8
	R10 = R9 & RB
	RF |= R10
	RA += BB[15]
	RA += RF
	RA += 0x674b87b8
	RA = ROR(RA, 14)
	RA += R9

	RF = ^RA
	RF &= R9
	R10 = RA & R8
	RF |= R10
	RB += BB[4]
	RB += RF
	RB += 0xc4d2faa2
	RB = ROR(RB, 20)
	RB += RA

	RF = ^RB
	RF &= RA
	R10 = RB & R9
	RF |= R10
	R8 += BB[9]
	R8 += RF
	R8 += 0x2f1b9b0e
	R8 = ROR(R8, 5)
	R8 += RB

	RF = ^R8
	RF &= RB
	R10 = R8 & RA
	RF |= R10
	R9 += BB[14]
	R9 += RF
	R9 += 0x1a55a83a
	R9 = ROR(R9, 9)
	R9 += R8

	RF = ^R9
	RF &= R8
	R10 = R9 & RB
	RF |= R10
	RA += BB[3]
	RA += RF
	RA += 0x91bf8078
	RA = ROR(RA, 14)
	RA += R9

	RF = ^RA
	RF &= R9
	R10 = RA & R8
	RF |= R10
	RB += BB[8]
	RB += RF
	RB += 0x86b23cfc
	RB = ROR(RB, 20)
	RB += RA

	RF = ^RB
	RF &= RA
	R10 = RB & R9
	RF |= R10
	R8 += BB[13]
	R8 += RF
	R8 += 0x7c24dd24
	R8 = ROR(R8, 5)
	R8 += RB

	RF = ^R8
	RF &= RB
	R10 = R8 & RA
	RF |= R10
	R9 += BB[2]
	R9 += RF
	R9 += 0x25ac16a7
	R9 = ROR(R9, 9)
	R9 += R8

	RC = R9 & RB
	RF = ^R9
	R10 = R8 & RF
	RC |= R10
	RA += BB[7]
	RA += RC
	RA += 0xd88f1be9
	RA = ROR(RA, 14)
	RA += R9

	RC = RA & R8
	RD = ^RA
	R10 = R9 & RD
	RC |= R10
	RB += BB[12]
	RB += RC
	RB += 0x8a19fa0b
	RB = ROR(RB, 20)
	RB += RA

	//ROUND3
	RC = RB & RF
	RE = RA & R9
	RC |= RE
	R8 += BB[5]
	R8 += RC
	R8 += 0x813c91ac
	R8 = ROR(R8, 4)
	R8 += RB

	RC = R8 & RD
	RD = RB & RA
	RC |= RD
	R9 += BB[8]
	R9 += RC
	R9 += 0xcf6eb59d
	R9 = ROR(R9, 11)
	R9 += R8

	RC = ^RB
	RC &= R9
	RD = R8 & RB
	RC |= RD
	RA += BB[11]
	RA += RC
	RA += 0xafdbd5e2
	RA = ROR(RA, 16)
	RA += R9

	RD = ^R8
	RD &= RA
	RE = R9 & R8
	RD |= RE
	RB += BB[14]
	RB += RD
	RB += 0x80e7c964
	RB = ROR(RB, 23)
	RB += RA

	RE = ^R9
	RE &= RB
	RF = RA & R9
	RE |= RF
	R8 += BB[1]
	R8 += RE
	R8 += 0xe15e5e35
	R8 = ROR(R8, 4)
	R8 += RB

	RE = ^RA
	RE &= R8
	RF = RB & RA
	RE |= RF
	R9 += BB[4]
	R9 += RE
	R9 += 0xfbd6b4af
	R9 = ROR(R9, 11)
	R9 += R8

	RE = ^RB
	RE &= R9
	RF = R8 & RB
	RE |= RF
	RA += BB[7]
	RA += RE
	RA += 0xbc162b09
	RA = ROR(RA, 16)
	RA += R9

	RE = ^R8
	RE &= RA
	RF = R9 & R8
	RE |= RF
	RB += BB[10]
	RB += RE
	RB += 0x35f562aa
	RB = ROR(RB, 23)
	RB += RA

	RE = ^R9
	RE &= RB
	RF = RA & R9
	RE |= RF
	R8 += BB[13]
	R8 += RE
	R8 += 0x27ef30d3
	R8 = ROR(R8, 4)
	R8 += RB

	RE = ^RA
	RE &= R8
	RF = RB & RA
	RE |= RF
	R9 += BB[0]
	R9 += RE
	R9 += 0x5ea6b3b5
	R9 = ROR(R9, 11)
	R9 += R8

	RE = ^RB
	RE &= R9
	RF = R8 & RB
	RE |= RF
	RA += BB[3]
	RA += RE
	RA += 0x299cca2a
	RA = ROR(RA, 16)
	RA += R9

	RE = ^R8
	RE &= RA
	RF = R9 & R8
	RE |= RF
	RB += BB[6]
	RB += RE
	RB += 0xc56e9540
	RB = ROR(RB, 23)
	RB += RA

	RE = ^R9
	RE &= RB
	RF = RA & R9
	RE |= RF
	R8 += BB[9]
	R8 += RE
	R8 += 0x1f269216
	R8 = ROR(R8, 4)
	R8 += RB

	RE = ^RA
	RE &= R8
	RF = RB & RA
	RE |= RF
	R9 += BB[12]
	R9 += RE
	R9 += 0x98ca2340
	R9 = ROR(R9, 11)
	R9 += R8

	RE = ^RB
	RE &= R9
	RF = R8 & RB
	RE |= RF
	RA += BB[15]
	RA += RE
	RA += 0xd02ead4d
	RA = ROR(RA, 16)
	RA += R9

	RC = R8 + BB[0]

	RE = ^R8
	R8 &= R9
	RE &= RA
	R8 |= RE
	RB += BB[2]
	R8 += RB
	R8 += 0x5bbad257
	R8 = ROR(R8, 23)
	R8 += RA

	RB = R9 + BB[7]

	//ROUND4
	RD = ^R8
	R9 |= RD
	R9 ^= RA
	R9 += RC
	R9 += 0xe6b95e33
	R9 = ROR(R9, 6)
	R9 += R8

	R5 := RA + BB[14]

	RD = ^R9
	RA |= RD
	RA ^= R8
	RA += RB
	RA += 0xedecac7a
	RA = ROR(RA, 10)
	RA += R9

	RD = R8 + BB[5]

	RE = ^RA
	R8 |= RE
	R8 ^= R9
	R5 += R8
	R5 += 0xedb098d7
	R5 = ROR(R5, 15)
	R5 += RA

	RE = R9 + BB[12] //?D238

	RF = ^R5
	R9 |= RF
	R9 ^= RA
	R9 += RD
	R9 += 0x21149206
	R9 = ROR(R9, 21)
	R9 += R5

	RF = RA + BB[3]

	R10 = ^R9
	RA |= R10
	RA ^= R5
	RA += RE
	RA += 0x4dd44b9a
	RA = ROR(RA, 6)
	RA += R9

	RE = R5 + BB[10] //D2C0

	R10 = ^RA
	R5 |= R10
	R5 ^= R9
	R5 += RF
	R5 += 0x1d14b76d
	R5 = ROR(R5, 10)
	R5 += RA

	R4 := R9 + BB[1]

	RF = ^R5
	R9 |= RF
	R9 ^= RA
	R9 += RE
	R9 += 0xa971bdae
	R9 = ROR(R9, 15)
	R9 += R5

	RE = RA + BB[8]

	RF = ^R9
	RA |= RF
	RA ^= R5
	R4 += RA
	R4 += 0x41499ee4
	R4 = ROR(R4, 21)
	R4 += R9

	R2 := R5 + BB[15] //D374

	RA = ^R4
	R5 |= RA
	R5 ^= R9
	R5 += RE
	R5 += 0xbecda772
	R5 = ROR(R5, 6)
	R5 += R4

	R7 := R9 + BB[6]

	RA = ^R5
	R9 |= RA
	R9 ^= R4
	R2 += R9
	R2 += 0x1a5dcda7
	R2 = ROR(R2, 10)
	R2 += R5

	R9 = R4 + BB[13] //D3EC

	RA = ^R2
	R4 |= RA
	R4 ^= R5
	R4 += R7
	R4 += 0x89adf877
	R4 = ROR(R4, 15)
	R4 += R2

	R3 := R5 + BB[4]

	R7 = ^R4
	R5 |= R7
	R5 ^= R2
	R5 += R9
	R5 += 0xbeab3f78
	R5 = ROR(R5, 21)
	R5 += R4

	R7 = R2 + BB[11] //D464

	R9 = ^R5
	R2 |= R9
	R2 ^= R4
	R2 += R3
	R2 += 0xc9579d4f
	R2 = ROR(R2, 6)
	R2 += R5

	R3 = ^R2
	R3 |= R4
	R3 ^= R5
	R3 += R7
	R3 += 0xc708d7b8
	R3 = ROR(R3, 10)
	R3 += R2

	R7 = R3 + salt[0]

	R1 := R4 + BB[2]
	R4 = R5 + BB[9] //D4F4

	R6 := ^R3
	R5 |= R6
	R5 ^= R2
	R1 += R5
	R1 += 0xa2dba837
	R1 = ROR(R1, 15)
	R1 += R3

	R5 = R2 + salt[3]

	R6 = ^R1
	R2 |= R6

	R6 = R1 + salt[2]
	R1 = R1 + salt[1]

	R2 ^= R3
	R2 += R4
	R2 += 0x5cbd8717
	R2 = ROR(R2, 21)
	R1 += R2

	dd := new(bytes.Buffer)
	binary.Write(dd, binary.LittleEndian, R7)
	binary.Write(dd, binary.LittleEndian, R1)
	binary.Write(dd, binary.LittleEndian, R6)
	binary.Write(dd, binary.LittleEndian, R5)

	return dd.Bytes()
}

func CalcCRC(input []byte) uint32 {

	R1 := uint32(input[1])
	R2 := uint32(input[0])

	R2 *= 0x85
	R2 += R1

	R3 := uint32(input[2])
	R3 &= 0xff
	R1 *= 0x85
	R1 += R3
	R1 *= 0x85
	R2 &= 0xff
	R2 *= 0x85
	R2 += R3
	R2 *= 0x85

	R4 := uint32(input[3])
	R4 &= 0xff
	R2 += R4
	R1 += R4
	R3 *= 0x85
	R3 += R4
	R3 *= 0x85
	R1 &= 0xFF
	R1 *= 0x85

	R4 = uint32(input[4])
	R4 &= 0xFF
	R1 += R4
	R3 += R4
	R2 *= 0x85
	R2 += R4
	R2 &= 0xFF
	R2 *= 0x85
	R3 &= 0xFF
	R3 *= 0x85

	R4 = uint32(input[5])
	R4 &= 0xFF
	R3 += R4
	R2 += R4
	R1 *= 0x85
	R1 += R4

	R4 = uint32(input[6])
	R4 &= 0xFF
	R1 *= 0x85
	R2 *= 0x85
	R2 += R4
	R1 += R4
	R3 *= 0x85
	R3 += R4

	R4 = uint32(input[7])
	R4 &= 0xFF
	R3 *= 0x85
	R1 &= 0xFF
	R1 *= 0x85
	R1 += R4
	R3 += R4
	R2 *= 0x85
	R2 += R4

	R4 = uint32(input[8])
	R4 &= 0xFF
	R2 &= 0xFF
	R2 *= 0x85
	R3 &= 0xFF
	R3 *= 0x85
	R3 += R4
	R2 += R4
	R1 *= 0x85
	R1 += R4

	R4 = uint32(input[9])
	R4 &= 0xFF
	R1 *= 0x85
	R2 *= 0x85
	R2 += R4
	R1 += R4
	R3 *= 0x85
	R3 += R4

	R4 = uint32(input[10])
	R4 &= 0xFF
	R3 *= 0x85
	R1 &= 0xFF
	R1 *= 0x85
	R1 += R4
	R3 += R4
	R2 *= 0x85
	R2 += R4

	R4 = uint32(input[11])
	R4 &= 0xFF
	R2 &= 0xFF
	R2 *= 0x85
	R3 &= 0xFF
	R3 *= 0x85
	R1 *= 0x85
	R1 += R4
	R3 += R4
	R2 += R4

	R6 := uint32(input[12])
	R6 &= 0xFF

	R5 := uint32(input[13])
	R5 &= 0xFF

	R4 = uint32(input[14])
	R4 &= 0xFF

	R1 *= 0x85
	R1 += R6
	R1 &= 0xFF
	R1 *= 0x85
	R1 += R5
	R1 *= 0x85
	R1 += R4

	R1 <<= 8
	R1 &= 0x6F00

	R2 *= 0x85
	R2 += R6
	R3 *= 0x85
	R3 += R6
	R3 *= 0x85
	R3 += R5
	R2 *= 0x85
	R2 += R5
	R2 &= 0x3F
	R1 |= R2

	R2 = R3 & 0xFF
	R2 *= 0x85
	R2 += R4

	R0 := uint32(input[15])
	R2 *= 0x85
	R0 += R2
	R0 <<= 16
	R0 &= 0x7B0000

	R0 |= R1
	R0 |= 0x42000000

	return R0
}

// rqt
func CalcMsgCrcForData_807(s []byte) uint32 {

	h := md5.Sum(s)
	return CalcMsgCrcForString_807(hex.EncodeToString(h[:]))
}

// CalcRQTX calc in hash
// @in: input string
func CalcMsgCrcForString_807(instr string) uint32 {
	if len(instr) != 32 {
		return 0
	}

	r := new(RQTXHASH)
	r.init()
	r.update(t1)
	var T1 [16]uint32

	input := bytes.NewBuffer([]byte(instr))
	binary.Read(input, binary.LittleEndian, &T1[0])
	binary.Read(input, binary.LittleEndian, &T1[1])
	binary.Read(input, binary.LittleEndian, &T1[2])
	binary.Read(input, binary.LittleEndian, &T1[3])
	binary.Read(input, binary.LittleEndian, &T1[4])
	binary.Read(input, binary.LittleEndian, &T1[5])
	binary.Read(input, binary.LittleEndian, &T1[6])
	binary.Read(input, binary.LittleEndian, &T1[7])
	T1[8] = 0x000000FC
	T1[9] = 0x00000000
	T1[10] = 0x00000000
	T1[11] = 0x00000000
	T1[12] = 0x00000000
	T1[13] = 0x00000000
	T1[14] = 0xFFFFFFFF
	T1[15] = 0xFFFCFFFF

	r.update(T1)
	hash1 := r.final()

	r = new(RQTXHASH)
	r.init()
	r.update(t2)

	var T2 [16]uint32
	input = bytes.NewBuffer(hash1)
	binary.Read(input, binary.LittleEndian, &T2[0])
	binary.Read(input, binary.LittleEndian, &T2[1])
	binary.Read(input, binary.LittleEndian, &T2[2])
	binary.Read(input, binary.LittleEndian, &T2[3])
	binary.Read(input, binary.LittleEndian, &T2[4])
	binary.Read(input, binary.LittleEndian, &T2[5])
	binary.Read(input, binary.LittleEndian, &T2[6])
	binary.Read(input, binary.LittleEndian, &T2[7])
	T2[8] = 0x000000FC
	T2[9] = 0x00000000
	T2[10] = 0x00000000
	T2[11] = 0x00000000
	T2[12] = 0x00000000
	T2[13] = 0x00000000
	T2[14] = 0xFFFFFFFF
	T2[15] = 0xFFFCFFFF

	r.update(T2)
	hash2 := r.final()

	FH := append(hash2,
		[]byte{
			0xDA, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xFB, 0x0F, 0x13, 0x64, 0x00, 0x00,
			0x00, 0x00}...,
	)
	return CalcCRC(trans(FH))
}
