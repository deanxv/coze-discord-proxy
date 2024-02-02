package middleware

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func authHelper(c *gin.Context) {
	secret := c.Request.Header.Get("proxy-secret")
	if common.ProxySecret != "" && !common.SliceContains(strings.Split(common.ProxySecret, ","), secret) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无权进行此操作,未提供正确的 proxy-secret",
		})
		c.Abort()
		return
	}
	c.Next()
	return
}

func authHelperForOpenai(c *gin.Context) {
	secret := c.Request.Header.Get("Authorization")
	secret = strings.Replace(secret, "Bearer ", "", 1)

	if common.ProxySecret != "" && !common.SliceContains(strings.Split(common.ProxySecret, ","), secret) {
		c.JSON(http.StatusUnauthorized, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "authorization(proxy-secret)校验失败",
				Type:    "invalid_request_error",
				Code:    "invalid_authorization",
			},
		})
		c.Abort()
		return
	}
	c.Next()
	return
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c)
	}
}

func OpenAIAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelperForOpenai(c)
	}
}
