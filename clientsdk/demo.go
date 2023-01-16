package clientsdk

//
//import (
//	"encoding/xml"
//	"errors"
//	"github.com/lunny/log"
//	"strings"
//	"time"
//
//	"feiyu.com/wx/clientsdk/baseinfo"
//	"feiyu.com/wx/clientsdk/baseutils"
//	"feiyu.com/wx/clientsdk/cecdh"
//	"feiyu.com/wx/clientsdk/mmtls"
//	"feiyu.com/wx/protobuf/wechat"
//	"github.com/golang/protobuf/proto"
//	"github.com/google/uuid"
//)
//
//// 二维码登录测试
//func qrcodeLoginDemo() *baseinfo.UserInfo {
//	userInfo := NewUserInfo(uuid.New().String(), "",nil)
//
//	// 长链接初始化MMTLS
//	dialer := GetDialer(userInfo)
//	tmpMMInfo, err := mmtls.InitMMTLSInfoLong(dialer, userInfo.LongHost, userInfo.ShortHost, nil)
//	if err != nil {
//		log.Println(err)
//		return nil
//	}
//	userInfo.MMInfo = tmpMMInfo
//
//	// 获取登录二维码
//	qrcodeResponse, err := getLoginQRCodeDemo(userInfo)
//	if err != nil {
//		log.Println(err)
//		return nil
//	}
//	// 二维码写入到文件
//	baseutils.WriteToFile(qrcodeResponse.Qrcode.Src, "./qrcode.png")
//	// 二维码信息
//	qrcodeUUID := qrcodeResponse.GetUuid()
//	qrcodeKey := qrcodeResponse.GetAes().GetKey()
//
//	// 检测扫码状态
//	bRet, err := checkQrcodeAndLoginDemo(userInfo, qrcodeUUID, qrcodeKey)
//	if !bRet {
//		log.Println("登录失败:", err)
//		return nil
//	}
//
//	log.Println("登录成功")
//	return userInfo
//}
//
//// 获取二维码请求示例，这里分开处理发送请求和解析响应，有利于统一长短连接处理数据的方式
//func getLoginQRCodeDemo(userInfo *baseinfo.UserInfo) (*wechat.LoginQRCodeResponse, error) {
//	// 发送获取二维码请求，得到响应数据
//	packHeader, err := SendLoginQRCodeRequest(userInfo)
//	if err != nil {
//		return nil, err
//	}
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeGetLoginQRCode {
//		return nil, errors.New("getLoginQRCodeDemo err: URLID != baseinfo.MMRequestTypeGetLoginQRCode")
//	}
//
//	// 解析数据
//	response := &wechat.LoginQRCodeResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//// checkQrcodeAndLoginDemo 检测扫码状态，然后登录
//func checkQrcodeAndLoginDemo(userInfo *baseinfo.UserInfo, qrcodeUUID string, qrcodeKey []byte) (bool, error) {
//	lgQrNotify, err := checkQrcodeDemo(userInfo, qrcodeUUID, qrcodeKey)
//	baseutils.ShowObjectValue(lgQrNotify)
//	if err != nil {
//		return false, err
//	}
//
//	// 没有检测到扫码状态 继续检测
//	if lgQrNotify.Wxid == nil && lgQrNotify.Wxnewpass == nil {
//		// 判断二维码是否过期
//		qrcodeValidTime := lgQrNotify.GetEffectiveTime()
//		if qrcodeValidTime <= 0 {
//			return false, errors.New("checkQrcodeAndLoginDemo err: qrcode invalid now")
//		}
//		// 暂停1秒
//		time.Sleep(time.Duration(1) * time.Second)
//		return checkQrcodeAndLoginDemo(userInfo, qrcodeUUID, qrcodeKey)
//	}
//
//	// 发送登陆请求
//	return checkManualAuthDemo(userInfo, lgQrNotify.GetWxnewpass(), lgQrNotify.GetWxid())
//}
//
//// 检测扫码状态, 这里要用长链接去检测，因为微信是采用长链接
//func checkQrcodeDemo(userInfo *baseinfo.UserInfo, qrcodeUUID string, qrcodeKey []byte) (*wechat.LoginQRCodeNotify, error) {
//	// 先获取要发送的数据
//	reqBytes, err := GetCheckLoginQRCodeReq(userInfo, qrcodeUUID, qrcodeKey)
//	if err != nil {
//		return nil, err
//	}
//
//	// 长链接发送
//	err = mmtls.MMTCPSendReq(userInfo.MMInfo, mmtls.MMLongOperationCheckQrcode, reqBytes)
//	if err != nil {
//		return nil, err
//	}
//
//	// 长链接接收数据, 应该用异步去接收，这里只是写个demo
//	recvInfo, err := mmtls.MMTCPRecvData(userInfo.MMInfo)
//	if err != nil {
//		return nil, err
//	}
//
//	// 系统推送,
//	if recvInfo.HeaderInfo.Operation < 1000000000 {
//		// 收到这个消息要发同步请求，去同步数据
//		return nil, errors.New("checkQrcodeDemo err: URLID != baseinfo.MMRequestTypeCheckLoginQRCode")
//	}
//
//	// 解析检测二维码状态响应数据
//	packHeader, err := DecodePackHeader(recvInfo.RespData, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeCheckLoginQRCode {
//		return nil, errors.New("checkQrcodeDemo err: URLID != baseinfo.MMRequestTypeCheckLoginQRCode")
//	}
//
//	// 解析数据
//	response := &wechat.CheckLoginQRCodeResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	if err != nil {
//		return nil, err
//	}
//
//	// 解密LoginQRCodeNotify
//	retBytes, err := baseutils.AesDecryptByteKey(response.LoginQrcodeNotifyPkg.NotifyData.Buffer, qrcodeKey)
//	if err != nil {
//		return nil, err
//	}
//	lgQrNotify := &wechat.LoginQRCodeNotify{}
//	err = proto.Unmarshal(retBytes, lgQrNotify)
//	return lgQrNotify, err
//}
//
//// 发送登录请求
//func checkManualAuthDemo(userInfo *baseinfo.UserInfo, newpass string, wxid string) (bool, error) {
//	authResponse, err := manualAuthDemo(userInfo, newpass, wxid)
//	if err != nil {
//		return false, err
//	}
//
//	// 重定向
//	if authResponse.GetBaseResponse().GetRet() == baseinfo.MMErrIdcRedirect {
//		// 修改服务器地址，然后重新发送登录包
//		userInfo.ShortHost = authResponse.GetDnsInfo().GetNewHostList().GetList()[1].GetSubstitute()
//		userInfo.LongHost = authResponse.GetDnsInfo().GetNewHostList().GetList()[0].GetSubstitute()
//		// 关闭长链接
//		if userInfo.MMInfo.Conn != nil {
//			userInfo.MMInfo.Conn.Close()
//		}
//		// 重新初始化MMTLS
//		dialer := GetDialer(userInfo)
//		tmpMMInfo, err := mmtls.InitMMTLSInfoLong(dialer, userInfo.LongHost, userInfo.ShortHost, nil)
//		if err != nil {
//			return false, err
//		}
//		userInfo.MMInfo = tmpMMInfo
//		return checkManualAuthDemo(userInfo, newpass, wxid)
//	} else if authResponse.GetBaseResponse().GetRet() == baseinfo.MMOk {
//		// 登录成功处理
//		myWxid := authResponse.AccountInfo.Wxid
//		if myWxid != nil {
//			userInfo.WxId = *myWxid
//		}
//		// 获取aesKey
//		ecServerPubKey := authResponse.AuthParam.EcdhKey.Key.GetBuffer()
//		userInfo.CheckSumKey = cecdh.ComputerECCKeyMD5(ecServerPubKey, userInfo.EcPrivateKey)
//		tmpAesKey, err := baseutils.AesDecryptByteKey(authResponse.AuthParam.SessionKey.Key, userInfo.CheckSumKey)
//		if err != nil {
//			return false, err
//		}
//		userInfo.SessionKey = tmpAesKey
//		// autoauthKey, token登录需要用到的key
//		userInfo.AutoAuthKey = authResponse.AuthParam.AutoAuthKey.Buffer
//		return true, nil
//	}
//
//	return false, errors.New("登录发生未知错误")
//}
//
//// 发送登录请求
//func manualAuthDemo(userInfo *baseinfo.UserInfo, newpass string, wxid string) (*wechat.ManualAuthResponse, error) {
//	packHeader, err := SendManualAuth(userInfo, newpass, wxid)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeManualAuth {
//		return nil, errors.New("manualAuthDemo err: URLID != baseinfo.MMRequestTypeManualAuth")
//	}
//
//	// 解析数据
//	response := &wechat.ManualAuthResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//// 获取账号信息
//func getProfileDemo(userInfo *baseinfo.UserInfo) (*wechat.GetProfileResponse, error) {
//	packHeader, err := SendGetProfileRequest(userInfo)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeGetProfile {
//		return nil, errors.New("getProfileDemo err: URLID != baseinfo.MMRequestTypeGetProfile")
//	}
//
//	// 解析数据
//	response := &wechat.GetProfileResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//// 同步信息，好友，消息，等等，任何信息都能同步到(例如：你在手机上给好友发个信息这里就可以同步到, 添加好友，接收到消息，等等)
//func NewSyncDemo(userInfo *baseinfo.UserInfo) {
//	for {
//		syncResp, err := newSyncRequestDemo(userInfo, 3, userInfo.SyncKey)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		// 跟新同步Key
//		userInfo.SyncKey = syncResp.GetKeyBuf().GetBuffer()
//
//		// 如果没有同步到数据则返回
//		cmdList := syncResp.GetCmdList()
//		syncCount := cmdList.GetCount()
//		log.Println("同步到消息条数--->", syncCount)
//		if syncCount <= 0 {
//			return
//		}
//
//		// 遍历同步到的每条信息
//		itemList := cmdList.GetItemList()
//		for index := uint32(0); index < syncCount; index++ {
//			item := itemList[index]
//			itemID := item.GetCmdId()
//			// 同步到联系人
//			if itemID == baseinfo.CmdIDModContact {
//				contact := new(wechat.ModContact)
//				err := proto.Unmarshal(item.CmdBuf.Data, contact)
//				if err != nil {
//					continue
//				}
//
//				log.Print(contact)
//				// 判断contact是否是群 == 0 不是群
//				if contact.GetChatroomVersion() == 0 {
//					continue
//				}
//
//				// 被移除群聊
//				if contact.GetChatRoomNotify() == 0 {
//					log.Println("消息免打扰群, 群wxid = ", contact.GetUserName().GetStr(), " 群昵称：", contact.GetNickName().GetStr())
//				} else {
//					log.Println("微信群, 群wxid = ", contact.GetUserName().GetStr(), " 群昵称：", contact.GetNickName().GetStr())
//				}
//			}
//			// 同步到消息
//			if itemID == baseinfo.CmdIDAddMsg {
//				var addMsg wechat.AddMsg
//				log.Println(addMsg)
//				// 判断接收人是不是群
//				userName := addMsg.GetToUserName().GetStr()
//				if !strings.HasSuffix(userName, "@chatroom") {
//					continue
//				}
//
//				// 首先类型：要是引用
//				if addMsg.GetMsgType() != baseinfo.MMAddMsgTypeRefer {
//					continue
//				}
//
//				// 解析引用的消息
//				tmpMsg := new(baseinfo.Msg)
//				err = xml.Unmarshal([]byte(addMsg.GetContent().GetStr()), tmpMsg)
//				if err != nil {
//					continue
//				}
//
//				// 判断是否支付类型
//				if tmpMsg.APPMsg.MsgType != baseinfo.MMAppMsgTypePayInfo {
//					continue
//				}
//
//				// 判断是否红包类型
//				if tmpMsg.APPMsg.WCPayInfo.SceneID != baseinfo.MMPayInfoSceneIDHongBao {
//					continue
//				}
//
//				// // 抢红包操作
//				// hbItem := new(baseinfo.HongBaoItem)
//				// hbItem.GroupWxid = userName
//				// hbItem.NativeURL = tmpMsg.APPMsg.WCPayInfo.NativeURL
//				// hongBaoURLItem, err := ParseHongBaoURL(hbItem.NativeURL)
//				// if err != nil {
//				// 	continue
//				// }
//				// hbItem.URLItem = hongBaoURLItem
//				// // 这里可以发打开红包请求, 后面添加
//			}
//
//			// 还有很多消息类型, 可以自己打印出来看看
//		}
//	}
//}
//
//// NewSyncDemo 同步请求
//func newSyncRequestDemo(userInfo *baseinfo.UserInfo, scene uint32, syncKey []byte) (*wechat.NewSyncResponse, error) {
//	// 发送同步消息
//	packHeader, err := SendNewSyncRequest(userInfo, scene)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeNewSync {
//		return nil, errors.New("getProfileDemo err: URLID != baseinfo.MMRequestTypeGetProfile")
//	}
//
//	// 解析数据
//	response := &wechat.NewSyncResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//// initContactListDemo 初始化通讯录列表
//func initContactListDemo(userInfo *baseinfo.UserInfo) {
//	// 第一步先获取所有的好友，公众号微信ID列表
//	userNameList := make([]string, 0)
//	contactSeq := uint32(0)
//	for {
//		// 先获取所有好友，群，公众号的微信ID列表
//		contactResp, err := initContactReqDemo(userInfo, contactSeq)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		contactSeq = contactResp.GetCurrentWxcontactSeq()
//
//		// 好友微信id列表
//		contactUserNameList := contactResp.GetContactUsernameList()
//		userNameList = append(userNameList, contactUserNameList[0:]...)
//		if len(contactUserNameList) < 100 {
//			break
//		}
//	}
//
//	// 分批获取所有 微信ID 对应的详细信息, 一次最多获取20个
//	userCount := len(userNameList)
//	offset := 0
//	for offset < userCount {
//		tmpCount := 20
//		if offset+tmpCount > userCount {
//			tmpCount = userCount - offset
//		}
//
//		// 批量(一次最多获取20个)根据微信号获取 对应的详细信息(好友、群、公众号)
//		tmpUserWxidList := userNameList[offset : offset+tmpCount]
//		briefResp, err := batchGetContactBriefInfoReqDemo(userInfo, tmpUserWxidList)
//		if err != nil {
//			log.Println(err)
//			break
//		}
//
//		// 遍历打印结果,
//		contactList := briefResp.GetContactList()
//		contactCount := len(contactList)
//		for tmpIndex := 0; tmpIndex < contactCount; tmpIndex++ {
//			contactItem := contactList[tmpIndex]
//			tmpContact := contactItem.GetContact()
//
//			tmpUserName := tmpContact.GetUserName().GetStr()
//			tmpNickName := tmpContact.GetNickName().GetStr()
//			// 暂时可以这样判断
//			if tmpContact.CustomizedInfo != nil {
//				log.Println("公众号：", tmpUserName, " ", tmpNickName)
//			} else {
//				log.Println("好友：", tmpUserName, " ", tmpNickName)
//			}
//		}
//		offset = offset + tmpCount
//	}
//}
//
//// 初始化通讯录
//func initContactReqDemo(userInfo *baseinfo.UserInfo, contactSeq uint32) (*wechat.InitContactResp, error) {
//	// 发送初始化通讯录请求
//	packHeader, err := SendInitContactReq(userInfo, contactSeq)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeInitContact {
//		return nil, errors.New("initContactReqDemo err: URLID != baseinfo.MMRequestTypeInitContact")
//	}
//
//	// 解析数据
//	response := &wechat.InitContactResp{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//func batchGetContactBriefInfoReqDemo(userInfo *baseinfo.UserInfo, userNameList []string) (*wechat.BatchGetContactBriefInfoResp, error) {
//	// 发送同步消息
//	packHeader, err := SendBatchGetContactBriefInfoReq(userInfo, userNameList)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeBatchGetContactBriefInfo {
//		return nil, errors.New("batchGetContactBriefInfoReqDemo err: URLID != baseinfo.MMRequestTypeBatchGetContactBriefInfo")
//	}
//	// 解析数据
//	response := &wechat.BatchGetContactBriefInfoResp{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//// token登录Demo代码
//func tokenLoginDemo(userInfo *baseinfo.UserInfo) error {
//	// 关闭长链接
//	if userInfo.MMInfo.Conn != nil {
//		userInfo.MMInfo.Conn.Close()
//	}
//
//	// 重新初始化MMTLS
//	dialer := GetDialer(userInfo)
//	tmpMMInfo, err := mmtls.InitMMTLSInfoLong(dialer, userInfo.LongHost, userInfo.ShortHost, nil)
//	if err != nil {
//		return err
//	}
//	userInfo.MMInfo = tmpMMInfo
//	authResp, err := autoAuthDemo(userInfo)
//	if err != nil {
//		return err
//	}
//
//	if authResp.GetBaseResponse().GetRet() == baseinfo.MMOk {
//		// 登录成功处理
//		myWxid := authResp.AccountInfo.Wxid
//		if myWxid != nil {
//			userInfo.WxId = *myWxid
//		}
//		// 获取aesKey
//		ecServerPubKey := authResp.AuthParam.EcdhKey.Key.GetBuffer()
//		userInfo.CheckSumKey = cecdh.ComputerECCKeyMD5(ecServerPubKey, userInfo.EcPrivateKey)
//		tmpAesKey, err := baseutils.AesDecryptByteKey(authResp.AuthParam.SessionKey.Key, userInfo.CheckSumKey)
//		if err != nil {
//			return err
//		}
//		userInfo.SessionKey = tmpAesKey
//		// autoauthKey, token登录需要用到的key
//		userInfo.AutoAuthKey = authResp.AuthParam.AutoAuthKey.Buffer
//		return nil
//	}
//	return errors.New("Token登录失败！！")
//}
//
//// 发送token登录请求
//func autoAuthDemo(userInfo *baseinfo.UserInfo) (*wechat.ManualAuthResponse, error) {
//	// 发送同步消息
//	packHeader, err := SendAutoAuthRequest(userInfo)
//	if err != nil {
//		return nil, err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeAutoAuth {
//		return nil, errors.New("batchGetContactBriefInfoReqDemo err: URLID != baseinfo.MMRequestTypeBatchGetContactBriefInfo")
//	}
//
//	// 解析数据
//	response := &wechat.ManualAuthResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	return response, err
//}
//
//func testGetCDNDnsInfo(userInfo *baseinfo.UserInfo) error {
//	packHeader, err := SendGetCDNDnsRequest(userInfo)
//	if err != nil {
//		return err
//	}
//
//	// 这里做判断, 这里肯定是相等的，因为知道调用的是这个请求，
//	// 这是给个示例处理方式，以后兼容长链接有用，根据URLID来判断是什么请求
//	if packHeader.URLID != baseinfo.MMRequestTypeGetCdnDNS {
//		return errors.New("SendGetCDNDnsRequest err: URLID != baseinfo.MMRequestTypeAutoAuth")
//	}
//
//	// 解析数据
//	response := &wechat.GetCDNDnsResponse{}
//	err = ParseResponseData(userInfo, packHeader, response)
//	if err != nil {
//		return err
//	}
//
//	userInfo.DNSInfo = response.GetDnsInfo()
//	userInfo.APPDnsInfo = response.GetAppDnsInfo()
//	userInfo.SNSDnsInfo = response.GetSnsDnsInfo()
//	userInfo.FAKEDnsInfo = response.GetFakeDnsInfo()
//	return nil
//}
