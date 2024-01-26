package controller

import (
	"coze-discord-proxy/discord"
	"coze-discord-proxy/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChannelCreate(c *gin.Context) {

	var channelModel model.Channel
	err := json.NewDecoder(c.Request.Body).Decode(&channelModel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	channelId, err := discord.ChannelCreate(discord.GuildId, channelModel.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "discord创建服务器异常",
		})
	} else {
		channelModel.Id = channelId
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    channelModel,
		})
	}
	return
}

func ChannelDel(c *gin.Context) {
	channelId := c.Param("Id")

	if channelId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	channelId, err := discord.ChannelDel(channelId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "discord删除频道异常",
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "discord删除频道成功",
		})
	}
	return
}
