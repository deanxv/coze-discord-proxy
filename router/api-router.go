package router

import (
	"coze-discord-proxy/common/config"
	"coze-discord-proxy/controller"
	_ "coze-discord-proxy/docs"
	"coze-discord-proxy/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetApiRouter(router *gin.Engine) {
	router.Use(middleware.CORS())
	//router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.RequestRateLimit())
	if config.SwaggerEnable == "" || config.SwaggerEnable == "1" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	if config.OnlyOpenaiApi != "1" {
		apiRouter := router.Group("/api")
		apiRouter.Use(middleware.Auth())
		{
			//chatRoute := apiRouter.Group("/chat")
			//chatRoute.POST("", controller.Chat)

			channelRoute := apiRouter.Group("/channel")
			channelRoute.POST("/create", controller.ChannelCreate)
			channelRoute.GET("/del/:id", controller.ChannelDel)
			channelRoute.GET("/del/all/cdp", controller.ChannelDelAllCdp)

			threadRoute := apiRouter.Group("/thread")
			threadRoute.POST("/create", controller.ThreadCreate)
		}
	}

	//https://api.openai.com/v1/images/generations
	v1Router := router.Group("/v1")
	v1Router.Use(middleware.OpenAIAuth())
	v1Router.POST("/chat/completions", controller.ChatForOpenAI)
	v1Router.POST("/images/generations", controller.ImagesForOpenAI)

}
