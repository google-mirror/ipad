package cecdh

import (
	"crypto"
	"crypto/elliptic"
	"io"
	"math/big"
)

// ECDH 秘钥交换算法的主接口
type ECDH interface {
	GenerateKey(io.Reader) ([]byte, []byte, error)
	Marshal(crypto.PublicKey) []byte
	Unmarshal([]byte) (crypto.PublicKey, bool)
	GenerateSharedSecret([]byte, crypto.PublicKey) ([]byte, error)
}

type ellipticECDH struct {
	ECDH
	curve elliptic.Curve
}

type ellipticPublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

// NewEllipticECDH 指定一种椭圆曲线算法用于创建一个ECDH的实例
// 关于椭圆曲线算法标准库里面实现了4种: 见crypto/elliptic
func NewEllipticECDH(curve elliptic.Curve) ECDH {
	return &ellipticECDH{
		curve: curve,
	}
}

// GenerateKey 基于标准库的NIST椭圆曲线算法生成秘钥对
func (e *ellipticECDH) GenerateKey(rand io.Reader) ([]byte, []byte, error) {
	var d []byte
	var x, y *big.Int
	var pub *ellipticPublicKey
	var err error
	d, x, y, err = elliptic.GenerateKey(e.curve, rand)
	if err != nil {
		return nil, nil, err
	}
	pub = &ellipticPublicKey{
		Curve: e.curve,
		X:     x,
		Y:     y,
	}
	return d, e.Marshal(pub), nil
}

// Marshal用于公钥的序列化
func (e *ellipticECDH) Marshal(p crypto.PublicKey) []byte {
	pub := p.(*ellipticPublicKey)
	return elliptic.Marshal(e.curve, pub.X, pub.Y)
}

// Unmarshal用于公钥的反序列化
func (e *ellipticECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	var key *ellipticPublicKey
	var x, y *big.Int
	x, y = elliptic.Unmarshal(e.curve, data)
	if x == nil || y == nil {
		return key, false
	}
	key = &ellipticPublicKey{
		Curve: e.curve,
		X:     x,
		Y:     y,
	}
	return key, true
}

// GenerateSharedSecret 通过自己的私钥和对方的公钥协商一个共享密码
func (e *ellipticECDH) GenerateSharedSecret(privKey []byte, pubKey crypto.PublicKey) ([]byte, error) {
	pub := pubKey.(*ellipticPublicKey)
	x, _ := e.curve.ScalarMult(pub.X, pub.Y, privKey)
	return x.Bytes(), nil
}
