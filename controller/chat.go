package controller

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/discord"
	"coze-discord-proxy/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// Chat 发送消息
// @Summary 发送消息
// @Description 发送消息
// @Tags chat
// @Accept json
// @Produce json
// @Param chatModel body model.ChatReq true "chatModel"
// @Success 200 {object} model.ReplyResp "Successful response"
// @Router /api/chat [post]
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

	sentMsg, err := discord.SendMessage(chatModel.ChannelId, chatModel.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "discord发送消息异常",
		})
		return
	}

	replyChan := make(chan model.ReplyResp)
	discord.RepliesChans[sentMsg.ID] = replyChan
	defer delete(discord.RepliesChans, sentMsg.ID)

	stopChan := make(chan string)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	if chatModel.Stream {
		c.Stream(func(w io.Writer) bool {
			select {
			case reply := <-replyChan:
				urls := ""
				if len(reply.EmbedUrls) > 0 {
					for _, url := range reply.EmbedUrls {
						urls += "\n" + url
					}
				}
				c.SSEvent("message", reply.Content+urls)
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
				replyResp.Content = reply.Content
				replyResp.EmbedUrls = reply.EmbedUrls
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

// ChatForOpenAI 发送消息-openai
// @Summary 发送消息-openai
// @Description 发送消息-openai
// @Tags chat -openai
// @Accept json
// @Produce json
// @Param request body model.OpenAIChatCompletionRequest true "request"
// @Success 200 {object} model.OpenAIChatCompletionResponse "Successful response"
// @Router /v1/chat/completions [post]
func ChatForOpenAI(c *gin.Context) {

	var request model.OpenAIChatCompletionRequest
	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "无效的参数",
				Type:    "invalid_request_error",
				Code:    "invalid_parameter",
			},
		})
		return
	}

	sentMsg, err := discord.SendMessage(discord.ChannelId, request.Messages[0].Content)
	if err != nil {
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "discord发送消息异常",
				Type:    "invalid_request_error",
				Code:    "discord_request_err",
			},
		})
		return
	}

	replyChan := make(chan model.OpenAIChatCompletionResponse)
	discord.RepliesOpenAIChans[sentMsg.ID] = replyChan
	defer delete(discord.RepliesOpenAIChans, sentMsg.ID)

	stopChan := make(chan string)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	if request.Stream {
		strLen := ""
		c.Stream(func(w io.Writer) bool {
			select {
			case reply := <-replyChan:
				newContent := strings.Replace(reply.Choices[0].Message.Content, strLen, "", 1)
				if newContent == "" {
					return true
				}
				reply.Choices[0].Delta.Content = newContent
				strLen += newContent
				reply.Object = "chat.completion.chunk"
				bytes, _ := common.Obj2Bytes(reply)
				c.SSEvent("", " "+string(bytes))
				return true // 继续保持流式连接
			case <-stopChan:
				c.SSEvent("", " [DONE]")
				return false // 关闭流式连接
			}
		})
	} else {
		var replyResp model.OpenAIChatCompletionResponse
		for {
			select {
			case reply := <-replyChan:
				replyResp = reply
			case <-stopChan:
				c.JSON(http.StatusOK, replyResp)
				return
			}
		}
	}
}
