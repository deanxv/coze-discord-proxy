package router

import (
	"coze-discord-proxy/controller"
	_ "coze-discord-proxy/docs"
	"coze-discord-proxy/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetApiRouter(router *gin.Engine) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.Auth())
	apiRouter := router.Group("/api")
	{
		chatRoute := apiRouter.Group("/chat")
		chatRoute.POST("/", controller.Chat)

		channelRoute := apiRouter.Group("/channel")
		channelRoute.POST("/create", controller.ChannelCreate)
		channelRoute.GET("/del/:id", controller.ChannelDel)

		threadRoute := apiRouter.Group("/thread")
		threadRoute.POST("/create", controller.ThreadCreate)
		//threadRoute.GET("/del/:id", controller.ChannelDel)
	}

}
