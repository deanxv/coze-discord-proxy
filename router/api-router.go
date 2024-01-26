package router

import (
	"coze-discord-proxy/controller"
	"coze-discord-proxy/middleware"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	router.Use(middleware.Auth())
	apiRouter := router.Group("/api")
	{
		chatRoute := apiRouter.Group("/chat")
		chatRoute.POST("/", controller.Chat)

		channelRoute := apiRouter.Group("/channel")
		channelRoute.POST("/add", controller.ChannelCreate)
		channelRoute.GET("/del/:Id", controller.ChannelDel)
	}

}
