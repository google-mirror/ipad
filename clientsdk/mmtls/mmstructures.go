package mmtls

import (
	"bufio"
	"crypto"
	"feiyu.com/wx/clientsdk/proxynet"
	"math/big"
	"net"

	"golang.org/x/net/proxy"
)

// AesGcmParam AesGcm加密解密参数
type AesGcmParam struct {
	AesKey []byte
	Nonce  []byte
}

// ClientEcdhKeys 客户端随机的两个EcdhKey私钥
type ClientEcdhKeys struct {
	PriKey1    crypto.PrivateKey
	PubKeyBuf1 []byte
	PriKey2    crypto.PrivateKey
	PubKeyBuf2 []byte
}

// HkdfKey28 HkdfKey28
type HkdfKey28 struct {
	AesKey []byte
	Nonce  []byte
}

// HkdfKey56 HkdfKey56
type HkdfKey56 struct {
	EncodeAesKey []byte
	EncodeNonce  []byte
	DecodeAesKey []byte
	DecodeNonce  []byte
}

// MMInfo MMInfo
type MMInfo struct {
	// 短链接 属性
	// mmtls 协议host 例如：hkextshort.weixin.qq.com，这个需要保存这数据库
	ShortHost string
	// mmtls路径 -- 例如：/mmtls/12345678(随机8位16进制字符串)，每次握手都随机一个
	ShortURL string
	// 短链接会话票据(服务端返回, 第一次握手不设置), 下一次握手选择其中一个发给服务器, 需要保存到数据库
	ShortPskList []*Psk
	// 握手扩展出来的用于后续加密的Key
	PskAccessKey []byte

	// 长链接 属性
	LongHost string
	LONGPort string

	// Deprecated:
	LONGClientSeq uint32 `json:"-"` // 不持久化
	// Deprecated:
	LONGServerSeq uint32 `json:"-"` // 不持久化
	// Deprecated:
	Conn   net.Conn `json:"-"` // 不持久化
	reader *bufio.Reader

	LongHdkfKey *HkdfKey56
	// ClientEcdhKeys
	ClientEcdhKeys *ClientEcdhKeys
	// 代理
	Dialer proxy.Dialer
	// 代理信息 http
	ProxyInfo *proxynet.WXProxyInfo
}

// EcdsaSignature 服务端传过来的校验数据
type EcdsaSignature struct {
	R, S *big.Int
}

// CipherSuiteInfo CipherSuiteInfo
type CipherSuiteInfo struct {
	SuiteCode uint16
	Clipher1  string
	Clipher2  string
	Clipher3  string
	Clipher4  string
	Clipher5  string
	Length1   uint32
	Length2   uint32
	Length3   uint32
}

// CipherSuite CipherSuite
type CipherSuite struct {
	SuiteCode uint16
	SuiteInfo *CipherSuiteInfo
}

// ClientKeyOffer ClientKeyOffer
type ClientKeyOffer struct {
	Version     uint32
	PublicValue []byte
}

// CertificateVerify CertificateVerify
type CertificateVerify struct {
	Signature []byte
}

// ClientKeyShareExtension ClientKeyShareExtension
type ClientKeyShareExtension struct {
	ClientKeyOfferList []*ClientKeyOffer
	CertificateVersion uint32
}

// EarlyEncryptDataExtension EarlyEncryptDataExtension
type EarlyEncryptDataExtension struct {
	ClientGmtTime uint32
}

// PreSharedKeyExtension PreSharedKeyExtension
type PreSharedKeyExtension struct {
	PskList []*Psk
}

// ServerKeyShareExtension ServerKeyShareExtension
type ServerKeyShareExtension struct {
	KeyOfferNameGroup uint32
	PublicValue       []byte
}

// Extension Extension
type Extension struct {
	ExtensionType uint16
	ExtensionData []byte
}

// EncryptedExtensions EncryptedExtensions
type EncryptedExtensions struct {
	ExtensionList []*Extension
}

// ClientHello ClientHello
type ClientHello struct {
	Version         uint16
	CipherSuiteList []*CipherSuite
	RandomBytes     []byte
	ClientGmtTime   uint32
	ExtensionList   []*Extension
}

// ServerHello ServerHello
type ServerHello struct {
	Version       uint16
	CipherSuite   *CipherSuite
	RandomBytes   []byte
	ExtensionList []*Extension
}

// Psk Psk
type Psk struct {
	Type                byte
	TicketKLifeTimeHint uint32
	MacValue            []byte
	KeyVersion          uint32
	Iv                  []byte
	EncryptedTicket     []byte
}

// ClientPsk CLientPsk
type ClientPsk struct {
	Psk            *Psk
	PskExpiredTime uint64
	PreSharedKey   []byte
}

// Finished Finished
type Finished struct {
	VerifyData []byte
}

// HTTPHandler HttpHandler
type HTTPHandler struct {
	URL   string
	Host  string
	MMPkg []byte
}

// KeyPair ECDH信息
type KeyPair struct {
	Version    uint32
	Nid        uint32
	PublicKey  []byte
	PrivateKey []byte
}

// NewSessionTicket NewSessionTicket
type NewSessionTicket struct {
	PskList []*Psk
}

// PskTicket PskTicket
type PskTicket struct {
	Version             byte
	MMTlsVersion        uint16
	CipherSuite         *CipherSuite
	KeyVersion          uint32
	TicketKLifeTimeHint uint32
	PreSharedKey        []byte
	MacKey              []byte
	ClientGmtTime       uint32
	ServerGmtTime       uint32
	EcdhVersion         uint32
	Valid               byte
}

// RecordHead RecordHead
type RecordHead struct {
	Type byte
	Tag  uint16
	Size uint16
}

// Alert Alert
type Alert struct {
	AlertLevel   byte
	AlertType    uint16
	FallBackURL  []byte
	SignatureURL []byte
}

// PackItem 包数量
type PackItem struct {
	RecordHead []byte
	PackData   []byte
}

// LongPackHeaderInfo 长链接请求包头部信息
type LongPackHeaderInfo struct {
	HeaderLen      uint16
	Version        uint16
	Operation      uint32
	SequenceNumber uint32
}

// LongRecvInfo 长链接接收信息
type LongRecvInfo struct {
	HeaderInfo *LongPackHeaderInfo
	RespData   []byte
}
