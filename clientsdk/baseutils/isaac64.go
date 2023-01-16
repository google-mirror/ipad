package baseutils

const (
	msg = "a Top Secret secret"
	key = "this is my secret key"
)

// RandCtx64 64位随机数
type RandCtx64 struct {
	RandCnt uint64
	Seed    [256]uint64
	MM      [256]uint64
	AA      uint64
	BB      uint64
	CC      uint64
}

// CreateISAacInst 创建新的
func CreateISAacInst(encKey uint64) *RandCtx64 {
	randCtx64 := &RandCtx64{}
	randCtx64.RandCnt = 255
	randCtx64.AA = 0
	randCtx64.BB = 0
	randCtx64.CC = 0
	Rand64Init(randCtx64, encKey)
	return randCtx64
}

// Rand64Init 初始化
func Rand64Init(randCtx64 *RandCtx64, encKey uint64) {
	a := uint64(0x9e3779b97f4a7c13)
	b := uint64(0x9e3779b97f4a7c13)
	c := uint64(0x9e3779b97f4a7c13)
	d := uint64(0x9e3779b97f4a7c13)
	e := uint64(0x9e3779b97f4a7c13)
	f := uint64(0x9e3779b97f4a7c13)
	g := uint64(0x9e3779b97f4a7c13)
	h := uint64(0x9e3779b97f4a7c13)

	randCtx64.Seed[0] = encKey
	for index := 1; index < 256; index++ {
		randCtx64.Seed[index] = 0
	}

	for index := 0; index < 4; index++ {
		a, b, c, d, e, f, g, h = mix64(a, b, c, d, e, f, g, h)
	}

	for index := 0; index < 256; index += 8 {
		a += randCtx64.Seed[index]
		b += randCtx64.Seed[index+1]
		c += randCtx64.Seed[index+2]
		d += randCtx64.Seed[index+3]
		e += randCtx64.Seed[index+4]
		f += randCtx64.Seed[index+5]
		g += randCtx64.Seed[index+6]
		h += randCtx64.Seed[index+7]
		a, b, c, d, e, f, g, h = mix64(a, b, c, d, e, f, g, h)
		randCtx64.MM[index] = a
		randCtx64.MM[index+1] = b
		randCtx64.MM[index+2] = c
		randCtx64.MM[index+3] = d
		randCtx64.MM[index+4] = e
		randCtx64.MM[index+5] = f
		randCtx64.MM[index+6] = g
		randCtx64.MM[index+7] = h
	}

	for index := 0; index < 256; index += 8 {
		a += randCtx64.MM[index]
		b += randCtx64.MM[index+1]
		c += randCtx64.MM[index+2]
		d += randCtx64.MM[index+3]
		e += randCtx64.MM[index+4]
		f += randCtx64.MM[index+5]
		g += randCtx64.MM[index+6]
		h += randCtx64.MM[index+7]
		a, b, c, d, e, f, g, h = mix64(a, b, c, d, e, f, g, h)
		randCtx64.MM[index] = a
		randCtx64.MM[index+1] = b
		randCtx64.MM[index+2] = c
		randCtx64.MM[index+3] = d
		randCtx64.MM[index+4] = e
		randCtx64.MM[index+5] = f
		randCtx64.MM[index+6] = g
		randCtx64.MM[index+7] = h
	}
	isAAC64(randCtx64)
}

func mix64(a uint64, b uint64, c uint64, d uint64, e uint64, f uint64, g uint64, h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	a -= e
	f ^= h >> 9
	h += a
	b -= f
	g ^= a << 9
	a += b
	c -= g
	h ^= b >> 23
	b += c
	d -= h
	a ^= c << 15
	c += d
	e -= a
	b ^= d >> 14
	d += e
	f -= b
	c ^= e << 20
	e += f
	g -= c
	d ^= f >> 17
	f += g
	h -= d
	e ^= g << 14
	g += h
	return a, b, c, d, e, f, g, h
}

func isAAC64(randCtx64 *RandCtx64) {
	randCtx64.CC++
	randCtx64.BB += randCtx64.CC
	for i, x := range randCtx64.MM {
		switch i % 4 {
		case 0:
			randCtx64.AA = ^(randCtx64.AA ^ randCtx64.AA<<21)
		case 1:
			randCtx64.AA = randCtx64.AA ^ randCtx64.AA>>5
		case 2:
			randCtx64.AA = randCtx64.AA ^ randCtx64.AA<<12
		case 3:
			randCtx64.AA = randCtx64.AA ^ randCtx64.AA>>33
		}
		randCtx64.AA += randCtx64.MM[(i+128)%256]
		y := randCtx64.MM[(x>>3)%256] + randCtx64.AA + randCtx64.BB
		randCtx64.MM[i] = y
		randCtx64.BB = randCtx64.MM[(y>>11)%256] + x
		randCtx64.Seed[i] = randCtx64.BB
	}
}

// ISAacRandom 随机一个整数
func ISAacRandom(randCtx64 *RandCtx64) (r uint64) {
	r = randCtx64.Seed[randCtx64.RandCnt]
	if randCtx64.RandCnt == 0 {
		isAAC64(randCtx64)
		randCtx64.RandCnt = 255
	} else {
		randCtx64.RandCnt--
	}
	return
}
