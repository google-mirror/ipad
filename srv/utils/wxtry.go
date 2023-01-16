package utils

import (
	"fmt"
	"runtime/debug"
)

// 异常处理
func TryE(userName string) {
	errs := recover()
	if errs == nil {
		return
	}
	_, _ = fmt.Println(fmt.Sprintf("%s,%vrn", userName, errs)) //输出panic信息
	_, _ = fmt.Println(string(debug.Stack()))                  //输出堆栈信息
}
