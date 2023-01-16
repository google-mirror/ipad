package srvconfig

import (
	"encoding/json"
	"errors"
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/clientsdk/baseutils"
	"fmt"
	"github.com/lunny/log"
	"net"
	"os"
	"strconv"
)

// GlobalSetting 全局设置
var GlobalSetting Setting

// TaskExecWaitTimes 任务执行间隔时间 500毫秒
var TaskExecWaitTimes = uint32(500)

// Redis configuration.
type RedisConfig struct {
	Host            string
	Port            int
	Db              int
	Pass            string // Password for AUTH.
	MaxIdle         int    // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int    // Maximum number of connections limit (default is 0 means no limit).
	IdleTimeout     int    // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	MaxConnLifetime int    // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	ConnectTimeout  int    // Dial connection timeout.
}

// Setting 设置
type Setting struct {
	Host             string `json:"host"`
	Port             string `json:"port"`
	WorkerPoolSize   uint32 `json:"workerpoolsize"`
	MaxWorkerTaskLen uint32 `json:"maxworkertasklen"`
	WebDomain        string `json:"webdomain"`
	WebTaskName      string `json:"webtaskname"`
	WebTaskAppNumber string `json:"webtaskappnumber"`
	//Redis
	RedisConfig RedisConfig `json:"redisConfig"`
	//消息同步 是否按微信id 发布消息
	NewsSynWxId     bool   `json:"newsSynWxId"`
	Dt              bool   `json:"dt"`
	Pprof           bool   `json:"pprof"`
	MysqlConnectStr string `json:"mySqlConnectStr"`
	//队例名
	Topic string `json:"topic"`
	//rocketMq是否开启
	RocketMq bool `json:"rocketMq"`
	//rocketMq地址
	RocketMqHost string `json:"rocketMqHost"`
	//rm key
	RocketAccessKey string `json:"rocketAccessKey"`
	RocketSecretKey string `json:"rocketSecretKey"`
	//rabbitMq是否开启
	RabbitMq bool `json:"rabbitMq"`
	//rabbitMqUrl 的url包含账号密码
	RabbitMqUrl string `json:"rabbitMqUrl"`
	//kafka是否开启
	Kafka         bool   `json:"kafka"`
	KafkaUrl      string `json:"kafkaUrl"`
	KafkaUsername string `json:"kafkaUsername"`
	KafkaPassword string `json:"kafkaPassword"`
	// 当前服务器外网IP地址
	TargetIp string
	LogLevel string `json:"logLevel"`
	DbType   string
}

// getExternal 请求获取外网ip
func getExternal() string {
	//resp, err := http.Get("https://api.ipify.org/")
	//if err != nil {
	//	return fmt.Sprintf("%s:%s", GlobalSetting.Host, GlobalSetting.Port)
	//}
	//defer resp.Body.Close()
	//content, _ := ioutil.ReadAll(resp.Body)
	//if content != nil {
	//	return fmt.Sprintf("%s:%s", content, GlobalSetting.Port)
	//}
	return fmt.Sprintf("%s:%s", GlobalSetting.Host, GlobalSetting.Port)

}

func ConfigSetUp() {
	tmpHomeDir, err := os.Getwd()
	if err != nil {
		baseutils.PrintLog(err.Error())
		return
	}
	baseinfo.HomeDIR = tmpHomeDir

	// 读取配置文件
	configData, err := baseutils.ReadFile(baseinfo.HomeDIR + "/assets/setting.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(configData, &GlobalSetting)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(GlobalSetting.RedisConfig.Db, "GlobalSetting.RedisConfig.Db")

	// 获取当前服务器外网IP地址
	GlobalSetting.TargetIp = getExternal()
	//判断Ip
	/*ip, _ := getClientIp()
	if ip != "172.27.36.211" {
		fmt.Printf(ip + "->ip不一致，退出程序!---->")
		os.Exit(1)
	}*/
	if GlobalSetting.LogLevel != "" {
		level, err := strconv.Atoi(GlobalSetting.LogLevel)
		if err == nil {
			log.SetOutputLevel(level)
		}
	}
}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}
	return "", errors.New("Can not find the client ip address!")
}
