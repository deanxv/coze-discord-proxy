// @title COZE-DISCORD-PROXY
// @version 1.0.0
// @description COZE-DISCORD-PROXY 代理服务
// @BasePath
package main

import (
	"context"
	"coze-discord-proxy/common"
	"coze-discord-proxy/common/config"
	"coze-discord-proxy/discord"
	"coze-discord-proxy/middleware"
	"coze-discord-proxy/router"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go discord.StartBot(ctx, discord.BotToken)

	common.SetupLogger()
	common.SysLog("COZE-DISCORD-PROXY " + common.Version + " started")
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	if config.DebugEnabled {
		common.SysLog("running in debug mode")
	}

	// Initialize HTTP server
	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetUpLogger(server)
	router.SetApiRouter(server)

	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			common.FatalLog("failed to start HTTP server: " + err.Error())
		}
	}()

	// 等待中断信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// 收到信号后取消 context
	cancel()

	// 给 HTTP 服务器一些时间来关闭
	ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutDown()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		common.FatalLog("HTTP server Shutdown failed:" + err.Error())
	}
}
