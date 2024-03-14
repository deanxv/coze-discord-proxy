package middleware

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/common/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

var timeFormat = "2006-01-02T15:04:05.000Z"

var inMemoryRateLimiter common.InMemoryRateLimiter

func memoryRateLimiter(c *gin.Context, maxRequestNum int, duration int64, mark string) {
	key := mark + c.ClientIP()
	if !inMemoryRateLimiter.Request(key, maxRequestNum, duration) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success": false,
			"message": "请求过于频繁,请稍后再试",
		})
		c.Abort()
		return
	}
}

func rateLimitFactory(maxRequestNum int, duration int64, mark string) func(c *gin.Context) {
	// It's safe to call multi times.
	inMemoryRateLimiter.Init(config.RateLimitKeyExpirationDuration)
	return func(c *gin.Context) {
		memoryRateLimiter(c, maxRequestNum, duration, mark)
	}

}

func RequestRateLimit() func(c *gin.Context) {
	return rateLimitFactory(config.RequestRateLimitNum, config.RequestRateLimitDuration, "REQUEST_RATE_LIMIT")
}
