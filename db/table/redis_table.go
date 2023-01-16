package table

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	pb "feiyu.com/wx/protobuf/wechat"
	"github.com/gogo/protobuf/proto"
)

type SubMessageCheckLoginQrCode struct {
	Type             int
	TargetIp         string
	UUID             string
	CheckLoginResult *baseinfo.CheckLoginQrCodeResult
}

type SyncMessageResponse struct {
	Type            int
	TargetIp        string
	UUID            string
	UserName        string `json:"userName"`
	LoginState      uint32 `json:"loginState"`
	ModUserInfos    []*pb.ModUserInfo
	ModContacts     []*pb.ModContact
	DelContacts     []*pb.DelContact
	FunctionSwitchs []*pb.FunctionSwitch
	AddMsgs         []*pb.AddMsg
	ModUserImgs     []*pb.ModUserImg
	UserInfoExts    []*pb.UserInfoExt
	SnsObjects      []*pb.SnsObject
	SnsActionGroups []*pb.SnsActionGroup
	FavItem         *baseinfo.FavItem
	Key             *pb.SKBuiltinString_
	MsgIdRsp        string
	SendMsgResp     *pb.NewSendMsgResponse
	ErrorMsg        string
}

func (sync *SyncMessageResponse) GetContacts() []*pb.ModContact {
	return sync.ModContacts
}

func (sync *SyncMessageResponse) GetAddMsgs() []*pb.AddMsg {
	return sync.AddMsgs
}

func (sync *SyncMessageResponse) SetMessage(data []byte, cmdId int32) {
	switch cmdId {
	case 1:
		userInfo := new(pb.ModUserInfo)
		if err := proto.Unmarshal(data, userInfo); err != nil {
			//z.Errorf(err.Error())
			return
		}
		/*log.Printf("登录微信：[%s] 昵称 [%s] 手机 [%s] 别名 [%s]\n",
		userInfo.UserName.GetStr(),
		userInfo.NickName.GetStr(),
		userInfo.BindMobile.GetStr(),
		userInfo.GetAlias())*/
		sync.ModUserInfos = append(sync.ModUserInfos, userInfo)
	case 2:
		contact := new(pb.ModContact)
		if err := proto.Unmarshal(data, contact); err != nil {
			//z.Errorf(err.Error())
			return
		}
		sync.ModContacts = append(sync.ModContacts, contact)
	case 4: // CmdId = 4
		delContact := new(pb.DelContact)
		if err := proto.Unmarshal(data, delContact); err != nil {
			//z.Println(err)
			return
		}
		sync.DelContacts = append(sync.DelContacts, delContact)
	case 5:
		addMsg := new(pb.AddMsg)
		if err := proto.Unmarshal(data, addMsg); err != nil {
			//z.Println(err)
			return
		}
		sync.AddMsgs = append(sync.AddMsgs, addMsg)
	case 23: // CmdId = 23
		functionSwitch := new(pb.FunctionSwitch)
		if err := proto.Unmarshal(data, functionSwitch); err != nil {
			//z.Println(err)
			return
		}
		sync.FunctionSwitchs = append(sync.FunctionSwitchs, functionSwitch)
	case 35:
		userImg := new(pb.ModUserImg)
		if err := proto.Unmarshal(data, userImg); err != nil {
			//z.Println(err)
			return
		}
		sync.ModUserImgs = append(sync.ModUserImgs, userImg)
	case 44:
		userInfoExt := new(pb.UserInfoExt)
		if err := proto.Unmarshal(data, userInfoExt); err != nil {
			//z.Println(err)
			return
		}

		sync.UserInfoExts = append(sync.UserInfoExts, userInfoExt)
	case 45:
		snsObject := new(pb.SnsObject)
		if err := proto.Unmarshal(data, snsObject); err != nil {
			//z.Println(err)
			return
		}
		sync.SnsObjects = append(sync.SnsObjects, snsObject)
	case 46:
		snsActionGroup := new(pb.SnsActionGroup)
		if err := proto.Unmarshal(data, snsActionGroup); err != nil {
			//z.Println(err)
			return
		}
		sync.SnsActionGroups = append(sync.SnsActionGroups, snsActionGroup)
	default:
		/*empty := new(pb.EmptyMesssage)
		_ = proto.Unmarshal(data, empty)*/
		//logger.Printf("收到未处理类型：%d 的数据：%s\n", id, empty.String())
		//保存消息
	}

}
