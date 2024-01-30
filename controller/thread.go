package controller

import (
	"coze-discord-proxy/discord"
	"coze-discord-proxy/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ThreadCreate 创建线程
// @Summary 创建线程
// @Description 创建线程
// @Tags thread
// @Accept json
// @Produce json
// @Param threadModel body model.ThreadReq true "threadModel"
// @Success 200 {object} model.ThreadResp "Successful response"
// @Router /api/thread/create [post]
func ThreadCreate(c *gin.Context) {
	var threadModel model.ThreadReq
	err := json.NewDecoder(c.Request.Body).Decode(&threadModel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	if threadModel.ArchiveDuration != 60 && threadModel.ArchiveDuration != 1440 && threadModel.ArchiveDuration != 4320 && threadModel.ArchiveDuration != 10080 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "线程创建时间只可为[60,1440,4320,10080]",
		})
		return
	}

	threadId, err := discord.ThreadStart(threadModel.ChannelId, threadModel.Name, threadModel.ArchiveDuration)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "discord创建线程异常",
		})
	} else {
		var thread model.ThreadResp
		thread.Id = threadId
		thread.Name = threadModel.Name
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    thread,
		})
	}
	return
}
