package clientsdk

import "github.com/gogo/protobuf/proto"

type MessageBase struct {
	ProtoData    []byte        `json:"-"`
	ProtoMessage proto.Message `json:"-"`
	CgiType      int
	EncType      byte
	CgiUrl       string
	JoinUin      bool
	JoinCookie   bool
	IsLongLink   bool
	ResetSecKey  bool
}

//请求体结构
type ReqMessageSession struct {
	*MessageBase
	UserName           string
	uin                uint32
	cookie             []byte
	userSessionAesKey  []byte
	clientVersion      uint32
	loginEcdhSecretKey []byte

	clientSessionKey []byte
	serverSessionKey []byte

	marshalAfterData []byte
	finalData        []byte
	secKeyMgr        *SecLoginKeyMgr
	rsaVersion       byte
}
