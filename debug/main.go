package main

import (
	req "feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"fmt"
	"github.com/lunny/log"

	"feiyu.com/wx/db"
	"feiyu.com/wx/db/mq"
	"feiyu.com/wx/srv/srvconfig"
	resty "github.com/go-resty/resty/v2"
)

var _VERSION_ = "20220803.01"

func main() {
	log.SetOutputLevel(log.Lerror)
	fmt.Printf("程序版本号:v%s\n", _VERSION_)
	c := resty.New()
	fmt.Println(c)
	srvconfig.ConfigSetUp()
	//初始化数据库连接
	db.InitDB()

	db.RedisSetup()

	mq.InitKafKa()

	// _ = api.WXServerGinHttpApiStart()

	uuid := "719ca14e-9a51-4f55-b6d0-4677d0cf54da"
	_ = uuid
	queryKey := "719ca14e-9a51-4f55-b6d0-4677d0cf54da"

	okUUID := "dad26a06-5726-4a5a-a320-e7708efab0b0"
	okKey := "dad26a06-5726-4a5a-a320-e7708efab0b0"
	_, _ = okUUID, okKey

	result := service.NewSyncHistoryMessageService(queryKey, req.SyncModel{
		Scene: 1,
	})
	//result := service.GetLoginStatusService(queryKey, false, true)
	fmt.Println(result)
}
