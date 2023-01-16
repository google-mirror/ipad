package baseutils

import (
	"github.com/lunny/log"
	"time"
)

// PrintLog 打印日志
func PrintLog(logStr string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	log.Println(currentTime + ": " + logStr)
}
