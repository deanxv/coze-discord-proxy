package controller

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/discord"
	"coze-discord-proxy/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Chat 发送消息
// @Summary 发送消息
// @Description 发送消息
// @Tags chat
// @Accept json
// @Produce json
// @Param chatModel body model.ChatReq true "chatModel"
// @Param proxy-secret header string false "proxy-secret"
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
			"success": false,
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

	timer, err := setTimerWithHeader(chatModel.Stream, common.RequestOutTimeDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "超时时间设置异常",
		})
		return
	}

	if chatModel.Stream {
		c.Stream(func(w io.Writer) bool {
			select {
			case reply := <-replyChan:
				timerReset(chatModel.Stream, timer, common.RequestOutTimeDuration)
				urls := ""
				if len(reply.EmbedUrls) > 0 {
					for _, url := range reply.EmbedUrls {
						urls += "\n" + fmt.Sprintf("![Image](%s)", url)
					}
				}
				c.SSEvent("message", reply.Content+urls)
				return true // 继续保持流式连接
			case <-timer.C:
				// 定时器到期时，关闭流
				return false
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
			case <-timer.C:
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": "request_out_time",
				})
				return
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
// @Param Authorization header string false "Authorization"
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

	content := "Hi！"
	messages := request.Messages

	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		if message.Role == "user" {
			content = message.Content
			break
		}
	}

	sentMsg, err := discord.SendMessage(discord.ChannelId, content)
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

	timer, err := setTimerWithHeader(request.Stream, common.RequestOutTimeDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "超时时间设置异常",
		})
		return
	}

	if request.Stream {
		strLen := ""
		c.Stream(func(w io.Writer) bool {
			select {
			case reply := <-replyChan:
				timerReset(request.Stream, timer, common.RequestOutTimeDuration)
				newContent := strings.Replace(reply.Choices[0].Message.Content, strLen, "", 1)
				if newContent == "" && strings.HasSuffix(newContent, "[DONE]") {
					return true
				}
				reply.Choices[0].Delta.Content = newContent
				strLen += newContent
				reply.Object = "chat.completion.chunk"
				bytes, _ := common.Obj2Bytes(reply)
				c.SSEvent("", " "+string(bytes))
				return true // 继续保持流式连接
			case <-timer.C:
				// 定时器到期时，关闭流
				c.SSEvent("", " [DONE]")
				return false
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
			case <-timer.C:
				c.JSON(http.StatusOK, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: "请求超时",
						Type:    "request_error",
						Code:    "request_out_time",
					},
				})
				return
			case <-stopChan:
				c.JSON(http.StatusOK, replyResp)
				return
			}
		}
	}
}

func setTimerWithHeader(isStream bool, defaultTimeout time.Duration) (*time.Timer, error) {
	var outTimeStr string
	if isStream {
		outTimeStr = common.StreamRequestOutTime
	} else {
		outTimeStr = common.RequestOutTime
	}
	if outTimeStr != "" {
		outTime, err := strconv.ParseInt(outTimeStr, 10, 64)
		if err != nil {

			return nil, err
		}
		return time.NewTimer(time.Duration(outTime) * time.Second), nil
	}
	return time.NewTimer(defaultTimeout), nil
}

func timerReset(isStream bool, timer *time.Timer, defaultTimeout time.Duration) error {
	var outTimeStr string
	if isStream {
		outTimeStr = common.StreamRequestOutTime
	} else {
		outTimeStr = common.RequestOutTime
	}
	if outTimeStr != "" {
		outTime, err := strconv.ParseInt(outTimeStr, 10, 64)
		if err != nil {
			return err
		}
		timer.Reset(time.Duration(outTime) * time.Second)
		return nil
	}
	timer.Reset(defaultTimeout)
	return nil
}
