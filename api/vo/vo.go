package vo

import (
	"feiyu.com/wx/srv/srvconfig"
	"fmt"
	"github.com/gin-gonic/gin"
)

type DTO struct {

	//自定义提示状态码
	Code interface{}

	//数据展示
	Data interface{}

	//提示文本
	Text interface{}
}

func NewFail(errMsg string) DTO {
	return DTO{
		Code: FAIL,
		Data: nil,
		Text: errMsg,
	}
}
func NewFailOffline(errMsg string) DTO {
	return DTO{
		Code: FAIL_Offline,
		Data: nil,
		Text: errMsg,
	}
}

func NewSuccess(h gin.H, msg string) DTO {
	return DTO{
		Code: SUCCESS,
		Data: h,
		Text: msg,
	}
}

func NewSuccessObj(h interface{}, msg string) DTO {
	return DTO{
		Code: SUCCESS,
		Data: h,
		Text: msg,
	}
}

func NewFailUUId(uuid string) DTO {
	return DTO{
		Code: FAIL_UUID,
		Data: nil,
		Text: fmt.Sprintf("%s 该链接不存在！", uuid),
	}
}

func NewFAILDoesNotBelongToServer(uuid string) DTO {
	return DTO{
		Code: FAIL_DoesNotBelongToServer,
		Data: nil,
		Text: fmt.Sprintf("%s 不属于[%s]该服务器实例", uuid, srvconfig.GlobalSetting.TargetIp),
	}
}

const (
	FAIL_DATA                  = -1 //数据格式错误
	FAIL_UUID                  = -2 //UUID 不存在
	FAIL_Bound                 = -3 //账号被绑定
	FAIL_Offline               = -4 //账号掉线
	FAIL_DoesNotBelongToServer = -5 //UUId不需要该服务器实例

	SUCCESS = 200 // 执行成功
	FAIL    = 300 // 执行失败
)
