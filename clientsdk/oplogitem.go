package clientsdk

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
)

// CreateModifyUserInfoField 创建修改 账号信息 项
func CreateModifyUserInfoField(modUserInfo *wechat.ModUserInfo, initFlag uint32, nickName string) *baseinfo.ModifyItem {
	// set ModUserInfo
	var bitFlag = initFlag
	// NickName
	if modUserInfo.NickName != nil &&
		modUserInfo.NickName.Str != nil &&
		len(*modUserInfo.NickName.Str) > 0 {
		bitFlag |= 0x2
	}
	// BindEmail
	if modUserInfo.BindEmail != nil &&
		modUserInfo.BindEmail.Str != nil &&
		len(*modUserInfo.BindEmail.Str) > 0 {
		bitFlag |= 0x8
	}

	// PersonalCard
	if modUserInfo.PersonalCard != nil &&
		*modUserInfo.PersonalCard != 0 {
		bitFlag |= 0x80
	}
	modUserInfo.BitFlag = &bitFlag
	data, marshalErr := proto.Marshal(modUserInfo)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}
	return &baseinfo.ModifyItem{
		CmdID: uint32(1),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateDeleteFriendField 创建删掉好友项
func CreateDeleteFriendField(modContact *wechat.ModContact) *baseinfo.ModifyItem {
	zeroValue32 := uint32(0)
	var emptySKString wechat.SKBuiltinString
	emptySKString.Str = nil
	var emptySKBuffer wechat.SKBuiltinString_
	emptySKBuffer.Len = &zeroValue32

	// 构造新的ModContact
	var tmpModContact wechat.ModContact
	tmpModContact.UserName = modContact.UserName
	tmpModContact.NickName = modContact.NickName
	tmpModContact.Pyinitial = &emptySKString
	tmpModContact.QuanPin = modContact.QuanPin
	tmpModContact.Sex = modContact.Sex
	tmpModContact.ImgBuf = &emptySKBuffer
	tmpModContact.ImgFlag = &zeroValue32

	// bitVal
	bitVal := modContact.GetBitVal()&modContact.GetBitMask() | 2
	bitVal = bitVal & 0xFFFFFFFE
	tmpModContact.BitVal = &bitVal
	// bitMask
	bitMask := uint32(0xFFFFFFFF)
	tmpModContact.BitMask = &bitMask
	tmpModContact.Remark = modContact.Remark
	tmpModContact.RemarkPyinitial = modContact.RemarkPyinitial
	tmpModContact.RemarkQuanPin = modContact.RemarkQuanPin
	tmpModContact.ContactType = modContact.ContactType
	tmpModContact.ChatRoomNotify = modContact.ChatRoomNotify
	tmpModContact.AddContactScene = &zeroValue32
	tmpModContact.DeleteContactScene = &zeroValue32

	data, marshalErr := proto.Marshal(&tmpModContact)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(2),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateModifyFriendField 创建修改好友 备注名项
func CreateModifyFriendField(modContact *wechat.ModContact) *baseinfo.ModifyItem {
	zeroValue32 := uint32(0)
	emptyString := string("")
	var emptySKString wechat.SKBuiltinString
	emptySKString.Str = &emptyString
	var emptySKBuffer wechat.SKBuiltinString_
	emptySKBuffer.Len = &zeroValue32

	// 构造新的ModContact
	var tmpModContact wechat.ModContact
	tmpModContact.UserName = modContact.UserName
	tmpModContact.NickName = modContact.NickName
	tmpModContact.Pyinitial = &emptySKString
	tmpModContact.QuanPin = modContact.QuanPin
	tmpModContact.Sex = modContact.Sex
	tmpModContact.ImgBuf = &emptySKBuffer
	tmpModContact.ImgFlag = &zeroValue32

	// bitVal
	bitVal := modContact.GetBitVal() & modContact.GetBitMask()
	bitVal = bitVal | 0x5
	tmpModContact.BitVal = &bitVal
	// bitMask
	bitMask := uint32(0xFFFFFFFF)
	tmpModContact.BitMask = &bitMask
	tmpModContact.Remark = modContact.Remark
	tmpModContact.RemarkPyinitial = modContact.RemarkPyinitial
	tmpModContact.RemarkQuanPin = modContact.RemarkQuanPin
	tmpModContact.ContactType = modContact.ContactType
	tmpModContact.RoomInfoCount = &zeroValue32
	tmpModContact.RoomInfoList = make([]*wechat.RoomInfo, 0)
	tmpModContact.ChatRoomNotify = modContact.ChatRoomNotify
	tmpModContact.AddContactScene = &zeroValue32
	tmpModContact.DeleteContactScene = &zeroValue32

	data, marshalErr := proto.Marshal(&tmpModContact)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(2),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateSaveGroupToAddressBookField 创建保存群聊到通讯录项
func CreateSaveGroupToAddressBookField(modContact *wechat.ModContact, bSafeToAddressBook bool) *baseinfo.ModifyItem {
	zeroValue32 := uint32(0)
	emptyString := string("")
	var emptySKString wechat.SKBuiltinString
	emptySKString.Str = &emptyString
	var emptySKBuffer wechat.SKBuiltinString_
	emptySKBuffer.Len = &zeroValue32

	// 构造新的ModContact
	var tmpModContact wechat.ModContact
	tmpModContact.UserName = modContact.UserName
	tmpModContact.NickName = modContact.NickName
	tmpModContact.Pyinitial = &emptySKString
	tmpModContact.QuanPin = modContact.QuanPin
	tmpModContact.Sex = modContact.Sex
	tmpModContact.ImgBuf = &emptySKBuffer
	tmpModContact.ImgFlag = &zeroValue32

	// bitVal
	bitVal := modContact.GetBitVal() & modContact.GetBitMask()
	if bSafeToAddressBook {
		bitVal = bitVal | 0x1
	} else {
		bitVal = bitVal & 0xFFFFFFFE
	}
	tmpModContact.BitVal = &bitVal
	// bitMask
	bitMask := uint32(0xFFFFFFFF)
	tmpModContact.BitMask = &bitMask
	tmpModContact.Remark = modContact.Remark
	tmpModContact.RemarkPyinitial = modContact.RemarkPyinitial
	tmpModContact.RemarkQuanPin = modContact.RemarkQuanPin
	tmpModContact.ContactType = modContact.ContactType
	tmpModContact.RoomInfoCount = &zeroValue32
	tmpModContact.RoomInfoList = make([]*wechat.RoomInfo, 0)
	tmpModContact.ChatRoomNotify = modContact.ChatRoomNotify
	tmpModContact.AddContactScene = &zeroValue32
	tmpModContact.DeleteContactScene = &zeroValue32
	data, marshalErr := proto.Marshal(&tmpModContact)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(2),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateModifyFriendBlack 创建修改好友项：拉黑
func CreateModifyFriendBlack(modContact *wechat.ModContact) *baseinfo.ModifyItem {
	zeroValue32 := uint32(0)
	emptyString := string("")
	var emptySKString wechat.SKBuiltinString
	emptySKString.Str = &emptyString
	var emptySKBuffer wechat.SKBuiltinString_
	emptySKBuffer.Len = &zeroValue32

	// 构造新的ModContact wxid_z9qu97pz6em212  wxid_9ozlmtdbfs8h12
	var tmpModContact wechat.ModContact
	tmpModContact.UserName = modContact.UserName
	tmpModContact.NickName = modContact.NickName
	tmpModContact.Pyinitial = &emptySKString
	tmpModContact.QuanPin = modContact.QuanPin
	tmpModContact.Sex = modContact.Sex
	tmpModContact.ImgBuf = &emptySKBuffer
	tmpModContact.ImgFlag = &zeroValue32

	// bitVal
	bitVal := modContact.GetBitVal() & modContact.GetBitMask()
	bitVal = bitVal | 0xF
	tmpModContact.BitVal = &bitVal
	// bitMask
	bitMask := uint32(0xFFFFFFFF)
	tmpModContact.BitMask = &bitMask
	tmpModContact.Remark = modContact.Remark
	tmpModContact.RemarkPyinitial = modContact.RemarkPyinitial
	tmpModContact.RemarkQuanPin = modContact.RemarkQuanPin
	tmpModContact.ContactType = modContact.ContactType
	tmpModContact.RoomInfoCount = &zeroValue32
	tmpModContact.RoomInfoList = make([]*wechat.RoomInfo, 0)
	// ChatRoomNotify
	chatRoomNotify := uint32(0)
	if modContact.GetChatroomStatus() != 1 {
		chatRoomNotify = 1
	}
	tmpModContact.ChatRoomNotify = &chatRoomNotify
	tmpModContact.AddContactScene = &zeroValue32
	tmpModContact.DeleteContactScene = &zeroValue32

	data, marshalErr := proto.Marshal(&tmpModContact)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(2),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateDelContactField 创建联系人项
func CreateDelContactField(userName string) *baseinfo.ModifyItem {
	var delContact wechat.DelContact

	// UserName
	var skUserName wechat.SKBuiltinString
	skUserName.Str = &userName
	delContact.UserName = &skUserName

	// DeleteContactScene
	delContactSecne := uint32(0)
	delContact.DeleteContactScene = &delContactSecne

	data, marshalErr := proto.Marshal(&delContact)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(7),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateQutiChatRoomItem 创建退出群聊项
// chatRoomName: 群微信ID
// userName: 自己的微信ID
func CreateQutiChatRoomItem(chatRoomName string, userName string) *baseinfo.ModifyItem {
	req := wechat.QuitChatRoom{
		ChatRoomName: &wechat.SKBuiltinString{
			Str: proto.String(chatRoomName),
		},
		UserName: &wechat.SKBuiltinString{
			Str: proto.String(userName),
		},
	}
	data, marshalErr := proto.Marshal(&req)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}
	return &baseinfo.ModifyItem{
		CmdID: uint32(16),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateFunctionSwitchItem 创建开关项
// funcID：功能ID
// switchType：功能开关值
func CreateFunctionSwitchItem(funcID uint32, switchType uint32) *baseinfo.ModifyItem {
	var functionSwitch wechat.FunctionSwitch
	functionSwitch.FunctionId = &funcID
	functionSwitch.SwitchValue = &switchType
	data, marshalErr := proto.Marshal(&functionSwitch)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(23),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateModifyGroupNameField 创建修改群名称(群主题名称) 项
// groupWxid: 群微信id
// topicName: 新的群主题名称
func CreateModifyGroupNameField(groupWxid string, topicName string) *baseinfo.ModifyItem {
	var modChatRoomTopic wechat.ModChatRoomTopic
	// 群号
	var groupWxidBuffer wechat.SKBuiltinString
	groupWxidBuffer.Str = &groupWxid
	modChatRoomTopic.ChatRoomName = &groupWxidBuffer
	// 新的昵称
	var topicNameBuffer wechat.SKBuiltinString
	topicNameBuffer.Str = &topicName
	modChatRoomTopic.ChatRoomTopic = &topicNameBuffer

	data, marshalErr := proto.Marshal(&modChatRoomTopic)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(27),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateModifyGroupNickNameField 创建修改群的个人昵称 项
// groupWxid: 群微信ID
// wxID: 自己的微信ID
// groupNickName: 新的群个人昵称
func CreateModifyGroupNickNameField(groupWxid string, wxID string, groupNickName string) *baseinfo.ModifyItem {
	var chatRoomMemberDisplayName wechat.ModChatRoomMemberDisplayName
	// 群号
	chatRoomMemberDisplayName.ChatRoomName = &groupWxid
	// 自己的微信ID
	chatRoomMemberDisplayName.UserName = &wxID
	// 新的个人群昵称
	chatRoomMemberDisplayName.DisplayName = &groupNickName

	data, marshalErr := proto.Marshal(&chatRoomMemberDisplayName)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(48),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateSnsShowTimeItem 允许朋友查看朋友的时间范围(三天、一个月、半年、无限制)
// snsUserInfo：通过getProfile请求可以获取到
// hours：时间范围小时(72-三天，720-一个月、4320-半年、4294967295-全部)，目前只能是这几个值
func CreateSnsShowTimeItem(snsUserInfo *wechat.SnsUserInfo, hours uint32) *baseinfo.ModifyItem {
	snsUserInfo.SnsPrivacyRecent = &hours
	data, marshalErr := proto.Marshal(snsUserInfo)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(51),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateSnsShowTenLineOfStrangerItem 允许陌生人查看十条朋友圈，设置
// snsUserInfo：通过getProfile请求可以获取到
// bShowTenLines: true-代表允许，false-代表不允许
func CreateSnsShowTenLineOfStrangerItem(snsUserInfo *wechat.SnsUserInfo, bShowTenLines bool) *baseinfo.ModifyItem {
	tmpSnsFlagEx := snsUserInfo.GetSnsFlagex()
	if bShowTenLines {
		tmpSnsFlagEx = tmpSnsFlagEx & 0xFFFFFFFE
	} else {
		tmpSnsFlagEx = tmpSnsFlagEx | 1
	}
	snsUserInfo.SnsFlagex = &tmpSnsFlagEx
	data, marshalErr := proto.Marshal(snsUserInfo)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(51),
		Len:   uint32(len(data)),
		Data:  data,
	}
}

// CreateBlackSnsItem 创建Oplog项：别人添加我是否需要验证
// needVerify: 0-不需要验证，1-需要验证
func CreateBlackSnsItem(friendWxid string, modType uint32) *baseinfo.ModifyItem {
	/*var snsBlackList wechat.ModSnsBlackList
	snsBlackList.ContactUsername = &friendWxid
	snsBlackList.ModType = &modType
	data, marshalErr := proto.Marshal(&snsBlackList)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModUserInfo failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}*/
	req := &wechat.FunctionSwitch{
		FunctionId:  proto.Uint32(4),
		SwitchValue: proto.Uint32(modType),
	}
	reqBuf, _ := proto.Marshal(req)
	return &baseinfo.ModifyItem{
		CmdID: uint32(23),
		Len:   uint32(len(reqBuf)),
		Data:  reqBuf,
	}
}

// CreateModifyNickNameField 创建修改 昵称 请求项
func CreateModifyNickNameField(newNickName string) *baseinfo.ModifyItem {
	var singleField wechat.ModSingleField

	// set singlefield
	var opType = uint32(1)
	singleField.OpType = &opType
	singleField.Value = &newNickName

	singleData, marshalErr := proto.Marshal(&singleField)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModSingleField failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}
	return &baseinfo.ModifyItem{
		CmdID: uint32(64),
		Len:   uint32(len(singleData)),
		Data:  singleData,
	}
}

// CreateModifySignatureField 创建修改 个性签名 请求项
func CreateModifySignatureField(newSignature string) *baseinfo.ModifyItem {
	var singleField wechat.ModSingleField

	// set singlefield
	var opType = uint32(2)
	singleField.OpType = &opType
	singleField.Value = &newSignature

	singleData, marshalErr := proto.Marshal(&singleField)
	if marshalErr != nil {
		log.Info("proto.Marshal wechat.ModSingleField failed: ", marshalErr)
		return &baseinfo.ModifyItem{}
	}

	return &baseinfo.ModifyItem{
		CmdID: uint32(64),
		Len:   uint32(len(singleData)),
		Data:  singleData,
	}
}
