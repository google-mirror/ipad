package main

import (
	"fmt"
	"github.com/lunny/log"

	"feiyu.com/wx/api"
	"feiyu.com/wx/db"
	"feiyu.com/wx/db/mq"
	"feiyu.com/wx/srv/srvconfig"
)

var _VERSION_ = "20220803.01"

func main() {
	log.SetOutputLevel(log.Lerror)
	fmt.Printf("程序版本号:v%s\n", _VERSION_)
	srvconfig.ConfigSetUp()
	//初始化数据库连接
	db.InitDB()

	db.RedisSetup()

	mq.InitKafKa()

	_ = api.WXServerGinHttpApiStart()
}
