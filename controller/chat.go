package controller

import (
	"coze-discord-proxy/discord"
	"coze-discord-proxy/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// Chat 发送消息
// @Summary 发送消息
// @Description 发送消息
// @Tags chat
// @Accept json
// @Produce json
// @Param chatModel body model.ChatReq true "chatModel"
// @Success 200 {object} model.ReplyResp "Successful response"
// @Router /chat [post]
func Chat(c *gin.Context) {

	var chatModel model.ChatReq
	err := json.NewDecoder(c.Request.Body).Decode(&chatModel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "无效的参数",
			"success": false,
		})
		return
	}

	sentMsg, err := discord.SendMessage(chatModel.ChannelID, chatModel.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "discord发送消息异常",
		})
	}

	replyChan := make(chan string)
	discord.RepliesChans[sentMsg.ID] = replyChan
	defer delete(discord.RepliesChans, sentMsg.ID)

	stopChan := make(chan string)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	if chatModel.Stream {
		c.Stream(func(w io.Writer) bool {
			select {
			case reply := <-replyChan:
				c.SSEvent("message", reply)
				return true // 继续保持流式连接
			case <-stopChan:
				return false // 关闭流式连接
			}
		})
	} else {
		var replyResp model.ReplyResp
		for {
			select {
			case reply := <-replyChan:
				replyResp.Content = reply
			case <-stopChan:
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data":    replyResp,
				})
				return
			}
		}
	}

}
