package api

import (
	"feiyu.com/wx/api/middleware"
	"feiyu.com/wx/api/router"
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func WXServerGinHttpApiStart() error {

	app := router.SetUpRouter(func(engine *gin.Engine) {
		//获取系统
		sysType := runtime.GOOS
		if sysType == "linux" {
			engine.Use(middleware.BasicAuth())
		}
		//中间件
		engine.Use(middleware.FilterInstanceMiddleware)
		engine.Use(middleware.LoggerToFile())

		//ginpprof.Wrap(engine)
		//中间件需要再创建接口之前完成
		//异常处理防止程序奔溃
		engine.Use(gin.Recovery())
	}, false)
	addr := ":" + srvconfig.GlobalSetting.Port
	fmt.Println("启动GIN服务成功！", addr)
	if srvconfig.GlobalSetting.Pprof {
		go func() {
			_ = http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	err := app.Run(addr)
	if err != nil {
		return err
	}
	return nil
}
