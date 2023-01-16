/**
 * @author mii
 * @date 2020/2/29 0029
 */

package db

import (
	"errors"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/db/mq"
	"feiyu.com/wx/db/table"
	"feiyu.com/wx/protobuf/wechat"
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/lunny/log"
	"time"
)

const (
	REDIS_OPERATION_SET             = "SET"
	REDIS_OPERATION_GET             = "GET"
	REDIS_OPERATION_LIST_LLEN       = "llen"
	REDIS_OPERATION_LIST_LPSUH      = "LPUSH"
	REDIS_OPERATION_LIST_LPOP       = "lpop"
	REDIS_OPERATION_CHANNEL_PUBLISH = "publish"
	REDIS_OPERATION_EXISTS          = "EXISTS"
	REDIS_OPERATION_DELETE          = "DEL"
)

const (
	DEFAULT_GROUP_NAME = "default" // Default configuration group name.
	DEFAULT_REDIS_PORT = 6379      // Default redis port configuration if not passed.
)

type Obj struct {
	Value interface{}
}

func RedisSetup() {
	//配置redis
	redisConfig := &gredis.Config{
		Host:            srvconfig.GlobalSetting.RedisConfig.Host,
		Port:            srvconfig.GlobalSetting.RedisConfig.Port,
		Db:              srvconfig.GlobalSetting.RedisConfig.Db,
		Pass:            srvconfig.GlobalSetting.RedisConfig.Pass,
		MaxIdle:         srvconfig.GlobalSetting.RedisConfig.MaxIdle,
		MaxActive:       srvconfig.GlobalSetting.RedisConfig.MaxActive,
		IdleTimeout:     time.Duration(srvconfig.GlobalSetting.RedisConfig.IdleTimeout) * time.Millisecond,
		MaxConnLifetime: time.Duration(srvconfig.GlobalSetting.RedisConfig.MaxConnLifetime) * time.Millisecond,
		ConnectTimeout:  time.Duration(srvconfig.GlobalSetting.RedisConfig.ConnectTimeout) * time.Millisecond,
	}
	gredis.SetConfig(redisConfig, DEFAULT_GROUP_NAME)
	conn := g.Redis().Conn()
	err := conn.Err()
	if err != nil {
		panic(err)
	}
}

func getRedisConn() *gredis.Conn {
	return g.Redis().GetConn()
}

// 创建一个同步信息保存列表
func CreateSyncMsgList(exId string) bool {
	ok, err := LPUSH(exId, "撒大声地")
	if err != nil {
		return false
	}
	return ok
}

// 从列表中取出对象然后反序列化成对象
func LPOPObj(k string, i interface{}) error {
	_var, err := LPOP(k)
	if err != nil {
		return err
	}
	err = gjson.DecodeTo(_var.Bytes(), &i)
	if err != nil {
		return err
	}
	return nil
}

// 将对象序列化成json保存
func LPUSHObj(k string, i interface{}) (bool, error) {
	iData, err := gjson.Encode(i)
	if err != nil {
		return false, err
	}
	ok, err := LPUSH(k, iData)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

func LPUSH(k string, i interface{}) (bool, error) {
	redisConn := getRedisConn()
	defer redisConn.Close()
	result, err := redisConn.Do(REDIS_OPERATION_LIST_LPSUH, k, i)
	if err != nil {
		//logger.Errorln(err)
		return false, err
	}
	return gconv.String(result) == "OK", nil
}

// 从列表中取出一个值
func LPOP(k string) (*g.Var, error) {
	redisConn := getRedisConn()
	defer redisConn.Close()
	r, err := redisConn.DoVar(REDIS_OPERATION_LIST_LPOP, k)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// 取列表长度
func LLEN(k string) (int32, error) {
	redisConn := getRedisConn()
	defer redisConn.Close()
	_var, err := redisConn.DoVar(REDIS_OPERATION_LIST_LLEN, k)
	if err != nil {
		//logger.Errorln(err)
		return 0, err
	}
	return _var.Int32(), nil
}

// 发布消息
func PUBLISH(k string, i interface{}) error {
	log.Debug("消息内容：", i)
	iData, err := gjson.Encode(i)
	if err != nil {
		return err
	}
	//PushRocketMq(k, iData)
	//_ = PublishRabbitMq(k, iData)
	//是否开启kafka,开启不存redis
	if srvconfig.GlobalSetting.Kafka {
		mq.SendKafKaMsg(k, iData)
		return nil
	}
	if srvconfig.GlobalSetting.RabbitMq {
		_ = mq.PublishRabbitMq(k, iData)
		return nil
	}
	/*redisConn := getRedisConn()
	defer redisConn.Close()
	_, err = redisConn.Do(REDIS_OPERATION_CHANNEL_PUBLISH, k, string(iData))
	if err != nil {
		return err
	}*/
	return PushQueue(k, i)
}
func SETExpirationObj(k string, i interface{}, expiration int64) error {
	//redisConn := getRedisConn()
	redisConn := getRedisConn()
	defer redisConn.Close()
	iData, err := gjson.Encode(i)
	if err != nil {
		return err
	}
	var result interface{}
	if expiration > 0 {
		result, err = redisConn.Do(REDIS_OPERATION_SET, k, iData, "EX", expiration)
	} else {
		result, err = redisConn.Do(REDIS_OPERATION_SET, k, iData)
	}

	if err != nil {
		//logger.Errorln(err)
		return err
	}
	if gconv.String(result) == "OK" {
		return nil
	}
	return errors.New(gconv.String(result))
}

func SETObj(k string, i interface{}) error {
	redisConn := getRedisConn()
	defer redisConn.Close()
	iData, err := gjson.Encode(i)
	if err != nil {
		return err
	}
	result, err := redisConn.Do(REDIS_OPERATION_SET, k, iData)
	if err != nil {
		//logger.Errorln(err)
		return err
	}
	if gconv.String(result) == "OK" {
		return nil
	}
	return errors.New(gconv.String(result))
}

func GETObj(k string, i interface{}) error {
	redisConn := getRedisConn()
	defer redisConn.Close()
	_var, err := redisConn.Do(REDIS_OPERATION_GET, k)
	if err != nil {
		return err
	}
	err = gjson.DecodeTo(_var, &i)
	if err != nil {
		return err
	}
	return nil
}

func DelObj(k string) error {
	redisConn := getRedisConn()
	defer redisConn.Close()
	_, err := redisConn.Do(REDIS_OPERATION_DELETE, k)
	if err != nil {
		return err
	}
	return nil
}

func Exists(k string) (bool, error) {
	//检查是否存在key值
	redisConn := getRedisConn()
	defer redisConn.Close()
	exists, err := redisConn.Do(REDIS_OPERATION_EXISTS, k)
	if err != nil {
		log.Println("illegal exception")
		return false, err
	}
	//log.Printf("exists or not: %v \n", exists)
	if exists.(int64) == 1 {
		return true, nil
	}
	return false, nil
}

// PublishSyncMsgWxMessage 发布微信消息
func PublishSyncMsgWxMessage(userInfo *baseinfo.UserInfo, response table.SyncMessageResponse) error {
	if userInfo == nil {
		return errors.New("PublishSyncMsgWxMessage userInfo == nil")
	}
	response.UUID = userInfo.UUID
	response.UserName = userInfo.GetUserName()
	response.Type = table.RedisPushSyncTypeWxMsg
	response.TargetIp = srvconfig.GlobalSetting.TargetIp
	if len(response.GetAddMsgs()) != 0 || len(response.GetContacts()) != 0 {
		log.Info("用户昵称=", userInfo.NickName+"--微信wxid="+userInfo.GetUserName()+"--发布同步信息消息中")
		if srvconfig.GlobalSetting.NewsSynWxId {
			return PUBLISH(response.UUID+"_wx_sync_msg_topic", &response)
		}
		return PUBLISH(srvconfig.GlobalSetting.Topic, &response)
	}
	return nil
}

// 异步接口返回推送
func PublishTxtImagePush(userInfo *baseinfo.UserInfo, sendMsgResp *wechat.NewSendMsgResponse, MsgIdRsp string) error {
	if userInfo == nil {
		return errors.New("PublishTxtImagePush userInfo == nil")
	}
	response := table.SyncMessageResponse{}
	response.UUID = userInfo.UUID
	response.UserName = userInfo.GetUserName()
	response.TargetIp = srvconfig.GlobalSetting.TargetIp
	response.MsgIdRsp = MsgIdRsp
	response.SendMsgResp = sendMsgResp
	response.Type = table.RedisPushTxtImageOk
	// 缓存reids
	if srvconfig.GlobalSetting.NewsSynWxId {
		return PUBLISH(response.UUID+"_wx_sync_msg_topic", &response)
	}
	return PUBLISH(srvconfig.GlobalSetting.Topic, &response)
}

// PublishSyncMsgLoginState 微信状态
func PublishSyncMsgLoginState(wxid string, state uint32) error {
	//log.Println("推送->wxid=="+wxid+"--->", state)
	response := table.SyncMessageResponse{}
	response.TargetIp = srvconfig.GlobalSetting.TargetIp
	response.Type = table.RedisPushSyncTypeLoginState
	response.UserName = wxid
	response.LoginState = state
	// 缓存reids
	if srvconfig.GlobalSetting.NewsSynWxId {
		return PUBLISH(response.UUID+"_wx_sync_msg_topic", &response)
	}
	return PUBLISH(srvconfig.GlobalSetting.Topic, &response)
}

// 初始化完成
func PublishWxInItOk(UUID string, state uint32) error {
	if UUID == "" || state == 0 {
		return errors.New("uuid || state ==  nil")
	}
	response := table.SyncMessageResponse{}
	response.TargetIp = srvconfig.GlobalSetting.TargetIp
	response.Type = table.RedisPushWxInItOk
	response.LoginState = state
	response.UUID = UUID
	response.UserName = "初始化完成!"
	// 缓存reids
	if srvconfig.GlobalSetting.NewsSynWxId {
		_, _ = LPUSHObj(response.UUID+"_wx_sync_msg_topic", &response)
		return PUBLISH(response.UUID+"_wx_sync_msg_topic", &response)
	}
	return PUBLISH(srvconfig.GlobalSetting.Topic, &response)
}

// PublishSyncMsgCheckLogin 扫码结果
func PublishSyncMsgCheckLogin(UUID string, result *baseinfo.CheckLoginQrCodeResult) error {
	if UUID == "" || result == nil {
		return errors.New("uuid || result ==  nil")
	}
	response := table.SubMessageCheckLoginQrCode{}
	response.TargetIp = srvconfig.GlobalSetting.TargetIp
	response.Type = table.RedisPushSyncTypeCheckLogin
	response.CheckLoginResult = result
	response.UUID = UUID
	// 缓存reids
	if srvconfig.GlobalSetting.NewsSynWxId {
		_, _ = LPUSHObj(response.UUID+"_wx_sync_msg_topic", &response)
		return PUBLISH(response.UUID+"_wx_sync_msg_topic", &response)
	}
	return PUBLISH(srvconfig.GlobalSetting.Topic, &response)
}

// 获取指定号缓存在redis的消息
func GETSyncMsg(uuid string) (int, []*table.SyncMessageResponse, error) {
	_len, err := LLEN(uuid + "_syncMsg")
	if err != nil {
		return 0, nil, err
	}
	if _len == 0 {
		return 0, nil, nil
	}
	syncMsgLen := _len
	if syncMsgLen > 10 {
		syncMsgLen = 10
	}
	opVals := []*table.SyncMessageResponse{}
	for i := int32(0); i <= syncMsgLen; i++ {
		opVal := &table.SyncMessageResponse{}
		err = LPOPObj(uuid+"_wx_sync_msg_topic", opVal)
		if err != nil {
			continue
		}
		opVals = append(opVals, opVal)
	}
	ContinueFlag := 0
	if (_len % 10) > 0 {
		ContinueFlag = 1
	}

	return ContinueFlag, opVals, nil
}

/***
* 发送队例
 */
func PushQueue(queueName string, i interface{}) (err error) {
	if i == nil {
		return err
	}
	con := getRedisConn()
	defer con.Close()
	_, err = con.Do("rpush", queueName, gconv.String(i))
	if err != nil {
		return fmt.Errorf("PushQueue redis error: %s", err)
	}
	log.Debug("----redis队列发布完成----!")
	return nil
}
