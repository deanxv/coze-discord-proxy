package middleware

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/common/config"
	"coze-discord-proxy/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func isValidSecret(secret string) bool {
	return config.ProxySecret != "" && !common.SliceContains(config.ProxySecrets, secret)
}

func authHelper(c *gin.Context) {
	secret := c.Request.Header.Get("proxy-secret")
	if isValidSecret(secret) {
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
	if isValidSecret(secret) {
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

	if config.ProxySecret == "" {
		c.Request.Header.Set("Authorization", "")
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
