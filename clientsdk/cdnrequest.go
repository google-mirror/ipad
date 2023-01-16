package clientsdk

import (
	"encoding/hex"
	"errors"
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/clientsdk/mmtls"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"feiyu.com/wx/protobuf/wechat"
)

// ---------------------创建请求---------------------------------------------------------------

// CreateCDNDownloadRequest 创建Cdn资源下载请求
func CreateCDNDownloadRequest(userInfo *baseinfo.UserInfo, dnsInfo *wechat.CDNDnsInfo, aesKeyBytes []byte, imageURL string, fileType uint32) *baseinfo.CdnImageDownloadRequest {
	request := &baseinfo.CdnImageDownloadRequest{}
	request.Ver = 1
	request.WeiXinNum = dnsInfo.GetUin()
	request.Seq = 1
	request.ClientVersion = baseinfo.ClientVersion
	if userInfo.DeviceInfo == nil {
		request.ClientOsType = baseinfo.AndroidDeviceType
	} else {
		request.ClientOsType = userInfo.DeviceInfo.OsType
	}
	// authkey
	authKey := dnsInfo.GetAuthKey().GetBuffer()
	request.AuthKey = authKey
	request.NetType = 1
	request.AcceptDupack = 1
	request.RsaVer = 1
	rsaBytes, _ := baseutils.CdnRsaEncrypt(aesKeyBytes)
	request.RsaValue = rsaBytes
	request.FileType = fileType
	request.WxChatType = 0
	request.FileID = imageURL
	request.LastRetCode = 0
	request.IPSeq = 0
	request.CliQuicFlag = 0
	request.WxMsgFlag = nil
	request.WxAutoStart = 1
	request.DownPicFormat = 1
	request.Offset = 0
	request.LargesVideo = 0
	request.SourceFlag = 0
	return request
}

// CreateVideoUploadRequest 创建视频上传请求
func CreateVideoUploadRequest(userInfo *baseinfo.UserInfo, videoItem *baseinfo.UploadVideoItem) (*baseinfo.CdnVideoUploadRequest, error) {
	request := &baseinfo.CdnVideoUploadRequest{}
	request.Ver = 1
	request.WeiXinNum = videoItem.CDNDns.GetUin()
	request.Seq = videoItem.Seq
	request.ClientVersion = baseinfo.ClientVersion
	if userInfo.DeviceInfo == nil {
		request.ClientOSType = baseinfo.AndroidDeviceType
	} else {
		request.ClientOSType = userInfo.DeviceInfo.OsType
	}
	request.AutoKey = videoItem.CDNDns.AuthKey.GetBuffer()
	request.NetType = 1
	request.AcceptDuPack = 1
	request.SafeProto = 1
	request.FileType = 4
	request.WeChatType = 0
	request.LastRetCode = 0
	request.IpSeq = 0
	request.HastHumb = 1
	// ToUser
	retBytes := baseutils.AesEncrypt([]byte(videoItem.ToUser), []byte("wxusrname2016cdn"))
	request.ToUSerName = "@cdn2_" + baseutils.BytesToHexString(retBytes, false)
	request.CompressType = 0
	request.NoCheckAesKey = 1
	request.EnaBleHit = 1
	request.ExistAnceCheck = 0
	request.AppType = 0

	// 加密视频
	videoEncodeData := baseutils.AesEncryptECB(videoItem.VideoData, videoItem.AesKey)
	// 加密微缩图
	thumbEncodeData := baseutils.AesEncryptECB(videoItem.ThumbData, videoItem.AesKey)

	userMd5 := baseutils.Md5ValueByte([]byte(userInfo.GetUserName()+"-"+videoItem.ToUser), false)
	timeStamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000/1000)
	videoKey := timeStamp + "baed6285091"
	fileKey := fmt.Sprintf("aupvideo_%s_%s_%s", userMd5, timeStamp, videoKey) // guid
	request.FileKey = fileKey
	request.LocalName = baseutils.Md5ValueByte(videoItem.VideoData, false) + ".jpg"
	request.Offset = 0
	request.LargesVideo = 1
	request.SourceFlag = 0
	request.DropRateFlag = 1
	request.ClientRsaVer = 1
	//rsaValue, _ := baseutils.RsaEncryptByVer(videoItem.AesKey, userInfo.GetLoginRsaVer())
	rsaValue, _ := baseutils.CdnRsaEncrypt(videoItem.AesKey)
	request.ClientRsaVal = rsaValue
	request.AdVideoFlag = 0

	request.FileMd5 = baseutils.Md5ValueByte(videoItem.VideoData, false)
	request.RawFileMd5 = baseutils.Md5ValueByte(videoItem.VideoData, false)
	request.DataCheckSum = baseutils.Adler32(0, videoEncodeData)
	request.EncThumbCrc = baseutils.Adler32(0, thumbEncodeData)
	request.ThumbData = thumbEncodeData
	request.RawThumbMd5 = baseutils.Md5ValueByte(videoItem.ThumbData, false)

	request.TotalSize = uint32(len(videoEncodeData))
	request.RawTotalSize = uint32(len(videoItem.VideoData))
	request.ThumbTotalSize = uint32(len(thumbEncodeData))
	request.RawThumbSize = uint32(len(videoItem.ThumbData))
	request.FileCrc = baseutils.Adler32(0, videoItem.VideoData)
	request.FileData = videoEncodeData
	request.Mp4identify = "a79c98ca478c707db3c80d28766f89e0"

	return request, nil
}

// CreateImageUploadRequest 创建上传高清图请求
func CreateImageUploadRequest(userInfo *baseinfo.UserInfo, imgItem *baseinfo.UploadImgItem) (*baseinfo.CdnImageUploadRequest, error) {
	request := &baseinfo.CdnImageUploadRequest{}
	request.Ver = 1
	request.WeiXinNum = uint32(imgItem.CDNDns.GetUin())
	request.Seq = imgItem.Seq
	if !strings.HasPrefix(userInfo.LoginDataInfo.LoginData, "A") && userInfo.DeviceInfo != nil {
		request.ClientOsType = userInfo.DeviceInfo.OsType
		request.ClientVersion = uint32(baseinfo.ClientVersion)
	} else {
		request.ClientOsType = baseinfo.AndroidDeviceType
		request.ClientVersion = baseinfo.AndroidClientVersion
	}
	request.AuthKey = imgItem.CDNDns.AuthKey.GetBuffer()
	request.NetType = 1
	request.AcceptDupack = 1
	request.SafeProto = 1
	request.FileType = 2
	request.WxChatType = 0
	request.LastRetCode = 0
	request.IPSeq = 0
	request.CliQuicFlag = 0
	request.CompressType = 0
	request.NoCheckAesKey = 1
	request.EnableHit = 1
	request.ExistAnceCheck = 0
	request.AppType = 1
	request.LocalName = imgItem.LocalName + imgItem.ExtName
	request.Offset = 0
	request.LargesVideo = 0
	request.SourceFlag = 0
	request.AdVideoFlag = 0

	// ToUser
	retBytes := baseutils.AesEncrypt([]byte(imgItem.ToUser), []byte("wxusrname2016cdn"))
	request.ToUser = "@cdn2_" + baseutils.BytesToHexString(retBytes, false)

	// ThumbData
	thumbItem := CreateThumbImage(imgItem.ImageData)
	if thumbItem == nil {
		return nil, errors.New("生成缩略图失败")
	}
	thumbEncodeData := baseutils.AesEncryptECB(thumbItem.Data, []byte(imgItem.AesKey))
	request.HasThumb = 1
	request.ThumbTotalSize = uint32(len(thumbEncodeData))
	request.RawThumbSize = uint32(len(thumbItem.Data))
	request.RawThumbMD5 = baseutils.Md5ValueByte(thumbItem.Data, false)
	request.EncThumbCRC = baseutils.Adler32(0, thumbEncodeData)
	request.ThumbData = thumbEncodeData

	// FileKey
	request.FileKey = "wxupload_" + imgItem.ToUser + imgItem.LocalName + "_" + strconv.Itoa(int(imgItem.CreateTime))

	// FileData
	fileEncodeData := baseutils.AesEncryptECB(imgItem.ImageData, []byte(imgItem.AesKey))
	request.TotalSize = uint32(len(fileEncodeData))
	request.RawTotalSize = uint32(len(imgItem.ImageData))
	request.FileMD5 = baseutils.Md5ValueByte(fileEncodeData, false)
	request.RawFileMD5 = baseutils.Md5ValueByte(imgItem.ImageData, false)
	request.FileCRC = baseutils.Adler32(0, imgItem.ImageData)
	request.DataCheckSum = baseutils.Adler32(0, fileEncodeData)
	request.FileData = fileEncodeData
	request.SetOfPicFormat = "011000"

	// prepareRequestItem
	var prepareRequestItem baseinfo.CDNUploadMsgImgPrepareRequestItem
	prepareRequestItem.ToUser = imgItem.ToUser
	prepareRequestItem.LocalName = imgItem.LocalName
	prepareRequestItem.CreateTime = imgItem.CreateTime
	prepareRequestItem.ThumbWidth = thumbItem.Width
	prepareRequestItem.ThumbHeight = thumbItem.Height
	prepareRequestItem.AesKey = imgItem.AesKey
	prepareRequestItem.Crc32 = request.FileCRC

	// SessionBuf
	request.SessionBuf = CreateCDNUploadMsgImgPrepareRequest(userInfo, &prepareRequestItem)
	return request, nil
}

// CreateSnsImageUploadRequest 创建Cdn上传朋友圈高清图请求
func CreateSnsImageUploadRequest(userInfo *baseinfo.UserInfo, snsImgItem *baseinfo.SnsUploadImgItem) (*baseinfo.CdnSnsImageUploadRequest, error) {
	request := &baseinfo.CdnSnsImageUploadRequest{}
	request.Ver = 1
	request.WeiXinNum = uint32(snsImgItem.CDNDns.GetUin())
	request.Seq = snsImgItem.Seq
	request.ClientVersion = uint32(baseinfo.ClientVersion)
	if userInfo.DeviceInfo == nil {
		request.ClientOsType = baseinfo.AndroidDeviceType
	} else {
		request.ClientOsType = userInfo.DeviceInfo.OsType
	}
	request.AuthKey = snsImgItem.CDNDns.AuthKey.GetBuffer()
	request.NetType = 1
	request.AcceptDupack = 1
	request.RsaVer = 1
	rsaValue, _ := baseutils.CdnRsaEncrypt(snsImgItem.AesKey)
	request.RsaValue = rsaValue
	request.FileType = 20201
	request.WxChatType = 0
	request.LastRetCode = 0
	request.IPSeq = 0
	request.CliQuicFlag = 0
	request.HasThumb = 0
	request.ToUser = ""
	request.CompressType = 1
	request.NoCheckAesKey = 1
	request.EnableHit = 1
	request.ExistAnceCheck = 0
	request.AppType = 108
	request.LocalName = "[TEMP]" + strconv.Itoa(int(snsImgItem.ImageID)) + "_" + strconv.Itoa(int(snsImgItem.CreateTime))
	request.FileKey = request.LocalName + "_" + strconv.Itoa(int(snsImgItem.CreateTime)+rand.Intn(1000))
	request.Offset = 0
	request.ThumbTotalSize = 0
	request.RawThumbSize = 0
	request.RawThumbMD5 = ""
	request.ThumbCRC = 0
	request.LargesVideo = 0
	request.SourceFlag = 0
	request.AdVideoFlag = 0

	// FileData
	fileEncodeData := baseutils.AesEncryptECB(snsImgItem.ImageData, []byte(snsImgItem.AesKey))
	request.TotalSize = uint32(len(snsImgItem.ImageData))
	request.RawTotalSize = request.TotalSize
	request.FileMD5 = baseutils.Md5ValueByte(snsImgItem.ImageData, false)
	request.RawFileMD5 = request.FileMD5
	request.FileCRC = baseutils.Adler32(0, snsImgItem.ImageData)
	request.DataCheckSum = baseutils.Adler32(0, fileEncodeData)
	request.FileData = snsImgItem.ImageData
	return request, nil
}

// CreateSnsVideoDownloadRequest 创建Cdn下载朋友圈视频请求
func CreateSnsVideoDownloadRequest(userInfo *baseinfo.UserInfo, snsVideoItem *baseinfo.SnsVideoDownloadItem) (*baseinfo.CdnSnsVideoDownloadRequest, error) {
	request := &baseinfo.CdnSnsVideoDownloadRequest{}
	request.Ver = 1
	request.WeiXinNum = uint32(snsVideoItem.CDNDns.GetUin())
	request.Seq = snsVideoItem.Seq
	request.ClientVersion = uint32(baseinfo.ClientVersion)
	if userInfo.DeviceInfo == nil {
		request.ClientOsType = baseinfo.AndroidDeviceType
	} else {
		request.ClientOsType = userInfo.DeviceInfo.OsType
	}
	request.AuthKey = snsVideoItem.CDNDns.AuthKey.GetBuffer()
	request.NetType = 1
	request.AcceptDupack = 1
	request.Signal = ""
	request.Scene = ""
	request.URL = snsVideoItem.URL
	request.RangeStart = snsVideoItem.RangeStart
	request.RangeEnd = snsVideoItem.RangeEnd
	request.LastRetCode = 0
	request.IPSeq = 0
	request.RedirectType = 0
	request.LastVideoFormat = 0
	request.VideoFormat = 2
	request.XSnsVideoFlag = snsVideoItem.XSnsVideoFlag
	return request, nil
}

// CreateCdnSnsVideoUploadRequest 上传朋友圈视频
func CreateCdnSnsVideoUploadRequest(userInfo *baseinfo.UserInfo, videoUploadItem *baseinfo.SnsVideoUploadItem) (*baseinfo.CdnSnsVideoUploadRequest, error) {
	request := &baseinfo.CdnSnsVideoUploadRequest{}
	request.Ver = 1
	request.WeiXinNum = uint32(videoUploadItem.CDNDns.GetUin())
	request.Seq = videoUploadItem.Seq
	request.ClientVersion = uint32(baseinfo.ClientVersion)
	if userInfo.DeviceInfo == nil {
		request.ClientOsType = baseinfo.AndroidDeviceType
	} else {
		request.ClientOsType = userInfo.DeviceInfo.OsType
	}
	request.AuthKey = videoUploadItem.CDNDns.AuthKey.GetBuffer()
	request.NetType = 1
	request.AcceptDupack = 1
	request.RsaVer = 1
	rsaValue, _ := baseutils.CdnRsaEncrypt(videoUploadItem.AesKey)
	request.RsaValue = rsaValue
	request.FileType = 20303 //20202
	request.WxChatType = 0
	request.LastRetCode = 0
	request.IPSeq = 0
	request.CliQuicFlag = 0
	request.IsStoreVideo = 0
	request.NoCheckAesKey = 1
	request.EnableHit = 1
	request.ExistAnceCheck = 0
	request.AppType = 102
	totalSize := uint32(len(videoUploadItem.VideoData))
	request.TotalSize = totalSize
	request.RawTotalSize = totalSize
	tmpLocalNameNoExt := "[TEMP]" + strconv.Itoa(int(videoUploadItem.VideoID)) + "_" + strconv.Itoa(int(videoUploadItem.CreateTime))
	request.LocalName = tmpLocalNameNoExt + ".mp4"
	request.FileKey = tmpLocalNameNoExt + "_" + strconv.Itoa(int(videoUploadItem.CreateTime)+rand.Intn(1000000000))
	request.Offset = 0

	// 暂时不设置Thumb数据
	thumbDataLen := uint32(len(videoUploadItem.ThumbData))
	request.HasThumb = 1
	request.ThumbTotalSize = thumbDataLen
	request.RawThumbSize = thumbDataLen
	request.RawThumbMD5 = baseutils.Md5ValueByte(videoUploadItem.ThumbData, false)
	request.ThumbCRC = baseutils.Adler32(0, videoUploadItem.ThumbData)
	request.ThumbData = videoUploadItem.ThumbData

	request.LargesVideo = 80
	request.SourceFlag = 0
	request.AdVideoFlag = 0
	fileEncodeData := baseutils.AesEncryptECB(videoUploadItem.VideoData, videoUploadItem.AesKey)
	request.Mp4Identify = baseutils.Md5ValueByte(fileEncodeData, false)
	md5Value := baseutils.Md5ValueByte(videoUploadItem.VideoData, false)
	request.FileMD5 = md5Value
	request.RawFileMD5 = md5Value
	request.FileCRC = baseutils.Adler32(0, videoUploadItem.VideoData)
	request.DataCheckSum = baseutils.Adler32(0, fileEncodeData)
	request.FileData = videoUploadItem.VideoData
	//request.UserLargeFileApi=true
	return request, nil
}

// ---------------------发送请求---------------------------------------------------------------

// SendCdnDownloadReuqest 发送CDN下载请求
func SendCdnDownloadReuqest(userInfo *baseinfo.UserInfo, downItem *baseinfo.DownMediaItem) (*baseinfo.CdnDownloadResponse, error) {
	dnsInfo := userInfo.DNSInfo
	aesKeyBytes := baseutils.HexStringToBytes(downItem.AesKey)
	request := CreateCDNDownloadRequest(userInfo, dnsInfo, aesKeyBytes, downItem.FileURL, downItem.FileType)
	sendData := PackCdnImageDownloadRequest(request)

	// 连接Cdn服务器
	serverIP := dnsInfo.FrontIplist[0].GetStr()
	serverPort := dnsInfo.FrontIpportList[0].PortList[0]
	conn, err := ConnectCdnServer(serverIP, serverPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 发送请求数据
	conn.Write(sendData)
	retData := CDNRecvData(conn)

	// 返回解密图片数据
	response, err := DecodeImageDownloadResponse(retData)
	if err != nil {
		return nil, err
	}
	response.FileData = baseutils.AesDecryptECB(response.FileData, aesKeyBytes)
	return response, nil
}

// 获取图片
func GetMsgBigImg(userInfo *baseinfo.UserInfo, m model.GetMsgBigImgModel) (*baseinfo.PackHeader, error) {
	datalen := m.Datatotalength
	if datalen > 65535 {
		datalen = 65535
	}
	datatotalength := m.Datatotalength
	Startpos := 0
	count := 0
	if datatotalength-Startpos > datalen {
		count = datalen
	} else {
		count = datatotalength - Startpos
	}

	req := wechat.GetMsgImgRequest{
		BaseRequest: GetBaseRequest(userInfo),
		MsgId:       proto.Uint32(m.MsgId),
		FromUserName: &wechat.SKBuiltinString{
			Str: proto.String(m.ToWxid),
		},
		ToUserName: &wechat.SKBuiltinString{
			Str: proto.String(userInfo.WxId),
		},
		TotalLen:     proto.Uint32(uint32(datatotalength)),
		StartPos:     proto.Uint32(uint32(Startpos)),
		DataLen:      proto.Uint32(uint32(count)),
		CompressType: proto.Uint32(0),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(&req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeForwardCdnImage, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/getmsgimg", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

func SendCdnUploadVideoRequest(userInfo *baseinfo.UserInfo, toUser string, imgData string, videoData []byte) (*baseinfo.CdnMsgVideoUploadResponse, error) {
	videoItem := &baseinfo.UploadVideoItem{}
	videoItem.ThumbData = []byte(imgData)
	videoItem.VideoData = videoData
	videoItem.ToUser = toUser
	videoItem.CDNDns = userInfo.DNSInfo
	videoItem.Seq = uint32(rand.Intn(10))
	videoItem.AesKey, _ = hex.DecodeString("1f22e78fd07d46a68b889aa222e93563")
	request, err := CreateVideoUploadRequest(userInfo, videoItem)
	if err != nil {
		return nil, err
	}
	sendData := PackCdnVideoUploadRequest(request)
	// 连接Cdn服务器
	serverIP := videoItem.CDNDns.FrontIplist[0].GetStr()
	serverPort := videoItem.CDNDns.FrontIpportList[0].PortList[0]
	conn, err := ConnectCdnServer(serverIP, serverPort)
	if err != nil {
		return nil, err
	}

	// 发送数据
	conn.Write(sendData)
	defer conn.Close()

	retryCount := uint32(0)
	// 接收响应信息
	for {
		// 接收响应信息，解析
		retData := CDNRecvData(conn)
		response, err := DecodeVideoUploadResponse(retData)
		if err != nil {
			if retryCount < 3 {
				retryCount++
				continue
			}
			return nil, err
		}

		retryCount = 0
		// 判断错误码
		if response.RetCode != 0 {
			return nil, errors.New("上传视频失败: ErrCode = " + GetErrStringByRetCode(response.RetCode))
		}

		// 判断 服务器是否接收完毕
		if response.FileID != "" {
			response.FileAesKey = "1f22e78fd07d46a68b889aa222e93563"
			response.ThumbDataSize = uint32(len(videoItem.ThumbData))
			response.VideoDataSize = uint32(len(videoItem.VideoData))
			response.Mp4identify = "a79c98ca478c707db3c80d28766f89e0"
			response.VideoDataMD5 = request.RawFileMd5
			response.ThumbWidth = 200
			response.ThumbHeight = 200
			return response, nil
		}

		// 设置请求数据
		//response.ReqData = request
	}
}

// SendCdnUploadImageReuqest 发送CDN上传图片请求
func SendCdnUploadImageReuqest(userInfo *baseinfo.UserInfo, toUser string, imgData []byte) (bool, error) {
	imgItem := &baseinfo.UploadImgItem{}
	imgItem.ToUser = toUser
	imgItem.Seq = uint32(rand.Intn(10))
	imgItem.LocalName = strconv.Itoa(rand.Intn(100))
	imgItem.ExtName = ".pic"
	// 随机生成AesKey
	imgItem.AesKey = []byte(baseutils.RandomStringByLength(16))
	imgItem.ImageData = imgData
	imgItem.CreateTime = uint32(time.Now().UnixNano() / 1000000000)
	// 设置Dns路由信息
	imgItem.CDNDns = userInfo.FAKEDnsInfo

	// 创建上传图片请求
	request, err := CreateImageUploadRequest(userInfo, imgItem)
	if err != nil {
		return false, err
	}

	// 打包请求
	sendData := PackCdnImageUploadRequest(request)
	// 连接Cdn服务器
	serverIP := imgItem.CDNDns.FrontIplist[0].GetStr()
	serverPort := imgItem.CDNDns.FrontIpportList[0].PortList[0]
	conn, err := ConnectCdnServer(serverIP, serverPort)
	if err != nil {
		return false, err
	}

	// 发送数据
	conn.Write(sendData)
	defer conn.Close()

	retryCount := uint32(0)
	// 接收响应信息
	sendImgLen := uint32(len(request.FileData))
	for {
		// 接收响应信息，解析
		retData := CDNRecvData(conn)
		response, err := DecodeImageUploadResponse(retData)
		if err != nil {
			if retryCount < 3 {
				retryCount++
				continue
			}
			return false, err
		}

		retryCount = 0
		// 判断错误码
		if response.RetCode != 0 {
			return false, errors.New("发送图片失败: ErrCode = " + GetErrStringByRetCode(response.RetCode))
		}

		// 判断 服务器是否接收完毕
		if response.RecvLen < sendImgLen {
			continue
		}

		// SKeyResp CDNUploadMsgImgPrepareResponse的错误码
		if response.SKeyResp != 0 {
			return false, errors.New("发送图片失败: SKeyResp = " + strconv.Itoa(int(response.SKeyResp)))
		}

		//  SKeyBuf
		skbufLen := len(response.SKeyBuf)
		if skbufLen <= 0 {
			continue
		}
		break
	}

	return true, nil
}

// SendCdnSnsUploadImageReuqest 发送CDN朋友圈上传图片请求
func SendCdnSnsUploadImageReuqest(userInfo *baseinfo.UserInfo, imgData []byte) (*baseinfo.CdnSnsImageUploadResponse, error) {
	// 生产SnsImgItem
	snsImgItem := &baseinfo.SnsUploadImgItem{}
	snsImgItem.Seq = uint32(rand.Intn(50))
	snsImgItem.AesKey = []byte(baseutils.RandomStringByLength(16))
	snsImgItem.ImageData = imgData
	snsImgItem.ImageID = CreateID(imgData)
	snsImgItem.CreateTime = uint32(time.Now().UnixNano() / 1000000000)
	if userInfo.SNSDnsInfo == nil {
		return nil, errors.New("SendCdnSnsUploadImageReuqest err:userInfo.SNSDnsInfo == nil")
	}
	snsImgItem.CDNDns = userInfo.SNSDnsInfo

	// 创建上传图片请求
	request, err := CreateSnsImageUploadRequest(userInfo, snsImgItem)
	if err != nil {
		return nil, err
	}

	// 打包请求
	sendData := PackCdnSnsImageUploadRequest(request)
	// 连接Cdn服务器
	serverIP := snsImgItem.CDNDns.FrontIplist[0].GetStr()
	serverPort := snsImgItem.CDNDns.FrontIpportList[0].PortList[0]
	conn, err := ConnectCdnServer(serverIP, serverPort)
	if err != nil {
		return nil, err
	}

	// 发送数据
	conn.Write(sendData)
	defer conn.Close()

	retryCount := uint32(0)
	// 接收响应信息
	sendImgLen := uint32(len(request.FileData))
	for {
		// 接收响应信息，解析
		retData := CDNRecvData(conn)
		response, err := DecodeSnsImageUploadResponse(retData)
		if err != nil {
			if retryCount < 3 {
				retryCount++
				continue
			}
			return nil, err
		}

		retryCount = 0
		// 判断错误码
		if response.RetCode != 0 {
			return nil, errors.New("发送图片失败: ErrCode = " + GetErrStringByRetCode(response.RetCode))
		}

		// 判断 服务器是否接收完毕
		if response.RecvLen < sendImgLen {
			continue
		}

		response.ImageMD5 = request.FileMD5
		width, height := GetImageBounds(request.FileData)
		response.ImageWidth = width
		response.ImageHeight = height
		return response, nil
	}
}

// SendCdnSnsVideoUploadReuqest 发送CDN朋友圈上传视频请求
func SendCdnSnsVideoUploadReuqest(userInfo *baseinfo.UserInfo, videoData []byte, thumbData []byte) (*baseinfo.CdnSnsVideoUploadResponse, error) {
	// 生产SnsImgItem
	snsVideoItem := &baseinfo.SnsVideoUploadItem{}
	snsVideoItem.Seq = uint32(rand.Intn(10))
	snsVideoItem.AesKey = []byte(baseutils.RandomStringByLength(16))
	snsVideoItem.VideoData = videoData
	snsVideoItem.ThumbData = thumbData
	snsVideoItem.VideoID = CreateID(videoData)
	snsVideoItem.CreateTime = uint32(time.Now().UnixNano() / 1000000000)
	snsVideoItem.CDNDns = userInfo.SNSDnsInfo

	// 创建上传朋友圈视频请求
	request, err := CreateCdnSnsVideoUploadRequest(userInfo, snsVideoItem)
	if err != nil {
		return nil, err
	}

	// 打包请求
	sendData := PackCdnSnsVideoUploadRequest(request)
	// 连接Cdn服务器
	serverIP := snsVideoItem.CDNDns.FrontIplist[0].GetStr()
	serverPort := snsVideoItem.CDNDns.FrontIpportList[0].PortList[0]
	conn, err := ConnectCdnServer(serverIP, serverPort)
	if err != nil {
		return nil, err
	}

	// 发送数据
	conn.Write(sendData)
	defer conn.Close()

	retryCount := uint32(0)
	// 接收响应信息
	for {
		// 接收响应信息，解析
		retData := CDNRecvData(conn)
		response, err := DecodeSnsVideoUploadResponse(retData)
		if err != nil {
			if retryCount < 3 {
				retryCount++
				continue
			}
			return nil, err
		}

		retryCount = 0
		// 判断错误码
		if response.RetCode != 0 {
			return nil, errors.New("上传朋友圈视频失败: ErrCode = " + GetErrStringByRetCode(response.RetCode))
		}

		// 判断 服务器是否接收完毕
		if response.RecvLen < request.TotalSize {
			continue
		}

		// 设置请求数据
		response.ReqData = request
		return response, nil
	}
}

// SendCdnSnsVideoDownloadReuqestPiece 分片下载
func SendCdnSnsVideoDownloadReuqestPiece(userInfo *baseinfo.UserInfo, snsVideoItem *baseinfo.SnsVideoDownloadItem) (*baseinfo.CdnSnsVideoDownloadResponse, error) {
	// 创建朋友圈视频下载请求
	request, err := CreateSnsVideoDownloadRequest(userInfo, snsVideoItem)
	if err != nil {
		return nil, err
	}

	// 打包请求
	sendData := PackCdnSnsVideoDownloadRequest(request)
	// 连接Cdn服务器
	serverIP := snsVideoItem.CDNDns.FrontIplist[0].GetStr()
	serverPort := snsVideoItem.CDNDns.FrontIpportList[0].PortList[0]
	conn, err := ConnectCdnServer(serverIP, serverPort)
	if err != nil {
		return nil, err
	}

	// 发送数据
	conn.Write(sendData)
	defer conn.Close()

	// 接收响应信息
	// 接收响应信息，解析
	retData := CDNRecvData(conn)
	response, err := DecodeSnsVideoDownloadResponse(retData)
	if err != nil {
		return nil, err
	}

	// 判断错误码
	if response.RetCode != 0 {
		return nil, errors.New("下载朋友圈视频失败: ErrCode = " + GetErrStringByRetCode(response.RetCode))
	}
	return response, nil
}

// SendCdnSnsVideoDownloadReuqest 发送CDN朋友圈视频下载请求
func SendCdnSnsVideoDownloadReuqest(userInfo *baseinfo.UserInfo, encKey uint64, tmpURL string) ([]byte, error) {
	retFileData := []byte{}
	lessLength := uint32(2000000)
	encLen := uint32(0)
	videoFlag := string("V2")

	retryCount := uint32(0)
	for {
		// 生产SnsImgItem
		var snsVideoItem baseinfo.SnsVideoDownloadItem
		snsVideoItem.Seq = uint32(rand.Intn(10))
		snsVideoItem.URL = tmpURL
		snsVideoItem.RangeStart = uint32(len(retFileData))
		snsVideoItem.RangeEnd = snsVideoItem.RangeStart + lessLength
		snsVideoItem.XSnsVideoFlag = videoFlag
		snsVideoItem.CDNDns = userInfo.SNSDnsInfo

		// 发送分片下载请求
		response, err := SendCdnSnsVideoDownloadReuqestPiece(userInfo, &snsVideoItem)
		if err != nil {
			if retryCount < 3 {
				retryCount++
				continue
			}
			return nil, err
		}

		retryCount = 0
		// 判断错误码
		if response.RetCode != 0 {
			return nil, errors.New("SendCdnSnsVideoDownloadReuqest err: response.RetCode != 0")
		}

		// 设置加密的字节数
		if encLen == 0 {
			encLen = response.XEncLen
		}

		// 合并数据
		retFileData = append(retFileData, response.FileData[0:]...)
		currentLen := uint32(len(retFileData))
		if currentLen >= response.TotalSize {
			break
		}

		// 如果没有读取完
		lessLength = response.TotalSize - currentLen
		videoFlag = response.XSnsVideoFlag
	}

	// 解密数据
	retFileData = baseutils.DecryptSnsVideoData(retFileData, encLen, encKey)
	return retFileData, nil
}

// ForwardCdnImageRequest 转发Cdn图片
func ForwardCdnImageRequest(userInfo *baseinfo.UserInfo, item baseinfo.ForwardImageItem) (*baseinfo.PackHeader, error) {
	req := &wechat.UploadMsgImgRequest{
		BaseRequest:    GetBaseRequest(userInfo),
		AesKey:         proto.String(item.AesKey),
		CdnThumbAesKey: proto.String(item.AesKey),

		CdnMidImgSize:   proto.Uint32(uint32(item.CdnMidImgSize)),
		CdnThumbImgSize: proto.Uint32(uint32(item.CdnThumbImgSize)),

		CdnSmallImgUrl: proto.String(item.CdnMidImgUrl),
		CdnThumbImgUrl: proto.String(item.CdnMidImgUrl),

		TotalLen: proto.Uint32(uint32(item.CdnMidImgSize)),
		DataLen:  proto.Uint32(uint32(item.CdnMidImgSize)),

		ClientImgId: &wechat.SKBuiltinString{
			Str: proto.String(func() string {
				//logger.Println(cli.GetAcctSect().GetUserName())
				userMd5 := strings.ToLower(baseutils.Md5Value(userInfo.GetUserName() + "-" + item.ToUserName)[:16])
				filekey := fmt.Sprintf("aupimg_%s_%d", userMd5, time.Now().UnixNano()/10000)
				return filekey
			}()),
		},
		RecvWxid: &wechat.SKBuiltinString{
			Str: proto.String(item.ToUserName),
		},
		MsgType: proto.Uint32(3),
		SenderWxid: &wechat.SKBuiltinString{
			Str: proto.String(userInfo.GetUserName()),
		},
		StartPos: proto.Uint32(0),
		ReqTime:  proto.Uint32(uint32(time.Now().Unix())),
		EncryVer: proto.Uint32(1),
		Data: &wechat.SKBuiltinString_{
			Len:    proto.Uint32(0),
			Buffer: []byte{},
		},
		MessageExt: proto.String("png"),
	}

	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeForwardCdnImage, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadmsgimg", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// ForwardCdnVideoRequest 转发Cdn视频
func ForwardCdnVideoRequest(userInfo *baseinfo.UserInfo, item baseinfo.ForwardVideoItem) (*baseinfo.PackHeader, error) {
	req := &wechat.UploadVideoRequest{
		BaseRequest: GetBaseRequest(userInfo),
		ClientMsgId: proto.String(func() string {
			userMd5 := strings.ToLower(baseutils.Md5Value(userInfo.GetUserName() + "-" + item.ToUserName)[:16])
			videokey := strconv.FormatInt(time.Now().UnixNano()/10000, 10) + "baed6285091"
			filekey := fmt.Sprintf("aupvideo_%s_%d_%s", userMd5, time.Now().UnixNano()/10000, videokey)
			return filekey
		}()),
		FromUserName: proto.String(userInfo.GetUserName()),
		ToUserName:   proto.String(item.ToUserName),

		AESKey:            proto.String(item.AesKey),
		CDNThumbAESKey:    proto.String(item.AesKey),
		CDNThumbImgHeight: proto.Int32(120),
		CDNThumbImgWidth:  proto.Int32(120),
		CDNThumbUrl:       proto.String(item.CdnVideoUrl),
		CDNVideoUrl:       proto.String(item.CdnVideoUrl),
		EncryVer:          proto.Int32(1),
		ThumbData: &wechat.BufferT{
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		VideoData: &wechat.BufferT{
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		MsgForwardType: proto.Uint32(1),
		CameraType:     proto.Uint32(2),
		MsgSource:      proto.String(""),
		FuncFlag:       proto.Uint32(0),
		ThumbTotalLen:  proto.Uint32(uint32(item.CdnThumbLength)),
		VideoTotalLen:  proto.Uint32(uint32(item.Length)),
		VideoStartPos:  proto.Uint32(0),
		ThumbStartPos:  proto.Uint32(0),
		PlayLength:     proto.Uint32(uint32(item.PlayLength)),
		VideoFrom:      proto.Int32(0),
		NetworkEnv:     proto.Uint32(1),
		ReqTime:        proto.Uint32(uint32(time.Now().Unix())),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	sendEncodeData := Pack(userInfo, srcData, baseinfo.MMRequestTypeForwardCdnVideo, 5)
	resp, err := mmtls.MMHTTPPostData(userInfo.GetMMInfo(), "/cgi-bin/micromsg-bin/uploadvideo", sendEncodeData)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}
