package middleware

import (
	"bytes"
	"feiyu.com/wx/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

/*
**
认证
*/
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		//if auth == "" {
		//	RsqK(c, "授权码不能为空，请Telegram联系https://t.me/h106548564")
		//	return
		//}
		modelAuth := &Auth{Code: 0, Data: 200, Msg: "成功"}
		var key = fmt.Sprintf("api:auth:%s", auth)
		//modelAuth = GetCacheAuth(key)
		//if modelAuth == nil {
		//	rsq := Get("http://119.45.28.143:8000/api/v1/auth?key=" + auth)
		//	if err := json.Unmarshal([]byte(rsq), &modelAuth); err == nil {
		//		if modelAuth.Data == 200 {
		//			CacheAuthAdd(key, modelAuth)
		//			c.Next()
		//			return
		//		}
		//		RsqK(c, modelAuth.Msg)
		//		return
		//	}
		//}
		if modelAuth.Data == 200 {
			CacheAuthAdd(key, modelAuth)
			c.Next()
			return
		}
		RsqK(c, modelAuth.Msg)
		return
	}
}

func RsqK(c *gin.Context, msg string) {
	c.AbortWithStatus(http.StatusCreated)
	respVo := &RespVo{
		Code:    http.StatusOK,
		Message: msg,
		Data:    http.StatusCreated,
	}
	c.JSON(http.StatusOK, respVo)
}

type Auth struct {
	Code int32
	Data int32
	Msg  string
}

// {"code":0,"message":"","data":""}
type RespVo struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string) string {
	// 超时时间：15秒
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}

// 添加授权缓存
func CacheAuthAdd(key string, val *Auth) {
	cache := db.GetCaChe()
	isContains, _ := cache.Contains(key)
	if !isContains {
		cache.Set(key, val, time.Hour*1)
		return
	}
}

// 查询授权缓存
func GetCacheAuth(key string) *Auth {
	cache := db.GetCaChe()
	isContains, _ := cache.Contains(key)
	if !isContains {
		return nil
	}
	obj, _ := cache.Get(key)
	return obj.(*Auth)
}
