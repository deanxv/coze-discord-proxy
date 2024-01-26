package middleware

import (
	"coze-discord-proxy/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func authHelper(c *gin.Context) {
	secret := c.Request.Header.Get("proxy-secret")
	if secret != common.ProxySecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无权进行此操作，未提供正确的 proxy-secret",
		})
		c.Abort()
		return
	}
	c.Next()
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c)
	}
}
