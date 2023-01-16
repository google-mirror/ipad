package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"time"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	//实例化
	logger := glog.Async(true)
	//设置输出
	_ = logger.SetPath("./log/")
	//设置日志级别
	logger.SetLevel(glog.LEVEL_ALL)
	//设置日志格式
	//logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.Request.Host
		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
