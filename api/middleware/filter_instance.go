package middleware

import (
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/db"
	"feiyu.com/wx/srv/bizcgi"
	"feiyu.com/wx/srv/wxmgr"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"net/http"
	"strconv"
	"strings"
)

var NotFilterReqRouterArray = []string{
	"/v1/login/GetLoginStatus",
	"/v1/login/WakeUpLogin",
	"/v1/login/A16Login",
	"/v1/login/DeviceLogin",
	"/v1/login/GetLoginQrCode",
	"/v1/login/GetLoginQrCodeNew",
	"/v1/login/CheckLoginStatus",
}

// FilterInstanceMiddleware 过滤不同服务器不同的实例 避免出现串号问题
func FilterInstanceMiddleware(ctx *gin.Context) {
	fullPath := ctx.FullPath()
	if fullPath == "" {
		ctx.JSON(http.StatusOK, vo.NewFail("路径不正确"))
		ctx.Abort()
		return
	}
	// 放行指定请求
	for _, s := range NotFilterReqRouterArray {
		if s == ctx.FullPath() {
			ctx.Next() //放行请求
			return
		}
	}

	// 是否带Key
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		ctx.Next() //放行请求
		return
	}
	//// 从数据库取该key 的用户信息
	//userInfoEntity := db.GetUserInfoEntity(queryKey)
	//if userInfoEntity == nil {
	//	ctx.Next()
	//	return
	//}
	//if userInfoEntity.TargetIp == "" {
	//	ctx.Next()
	//	return
	//}
	//// 判断是否为 同个服务器
	//if userInfoEntity.TargetIp != srvconfig.GlobalSetting.TargetIp {
	//	ctx.JSON(http.StatusOK, vo.NewFAILDoesNotBelongToServer(queryKey))
	//	ctx.Abort()
	//	return
	//}
	checkExIdPerformNoCreateConnect(queryKey, ctx)
}

// 检查实例Id是否存在 链接不存在返回错误不创建新链接
func checkExIdPerformNoCreateConnect(queryKey string, ctx *gin.Context) {
	//查询的queryKey为空创建一个链接实例
	if queryKey == "" {
		ctx.JSON(http.StatusOK, vo.NewFailUUId(queryKey))
		ctx.Abort()
		return
	}

	//查询该链接是否存在
	iwxAccount := wxmgr.WxAccountMgr.GetWXAccountByUserInfoUUID(queryKey)
	if iwxAccount == nil {
		//如果链接管理器不存在该链接查询数据库是否存在
		dbUserInfo := db.GetUSerInfoByUUID(queryKey)
		//数据库存在该链接数据 重新实例化一个链接对象
		if dbUserInfo == nil {
			ctx.JSON(http.StatusOK, vo.NewFailUUId(queryKey))
			ctx.Abort()
			return
		}
		//创建一个用户信息
		iwxAccount = service.CreateWXAccountByQueryKey(queryKey, "", dbUserInfo)
		//设置用户信息
		iwxAccount.SetUserInfo(dbUserInfo)
	} else {
		log.Debugf("GET Connection locfree success by %s", queryKey)
	}
	//取基本信息
	//loginState := iwxAccount.GetLoginState()
	//判断在线情况
	//if loginState == baseinfo.MMLoginStateNoLogin {
	//	ctx.AbortWithStatusJSON(http.StatusOK, &RespVo{
	//		Code:    1401,
	//		Message: "该账号需要重新登录！loginState == MMLoginStateNoLogin ",
	//	})
	//} else
	if !bizcgi.CheckOnLineStatus(iwxAccount) && !bizcgi.RecoverOnLineStatus(iwxAccount) {
		ctx.AbortWithStatusJSON(http.StatusOK, &RespVo{
			Code:    1402,
			Message: "账号离线,自动上线失败！loginState == " + strconv.Itoa(int(iwxAccount.GetLoginState())),
		})
	}
	ctx.Set("account", iwxAccount)
}
