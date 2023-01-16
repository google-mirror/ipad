package router

import (
	"feiyu.com/wx/srv/wxface"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account wxface.IWXAccount
		temp, _ := c.Get("account")
		account, _ = temp.(wxface.IWXAccount)
		ctx := &Context{
			c,
			account,
		}
		h(ctx)
	}
}

type Context struct {
	*gin.Context
	wxface.IWXAccount
}
