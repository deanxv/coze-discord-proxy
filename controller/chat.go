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
// @Param out-time header string false "out-time"
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

	timer, err := setTimerWithHeader(c, chatModel.Stream, common.RequestOutTimeDuration)
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
				timerReset(c, chatModel.Stream, timer, common.RequestOutTimeDuration)
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
// @Tags openai
// @Accept json
// @Produce json
// @Param request body model.OpenAIChatCompletionRequest true "request"
// @Param Authorization header string false "Authorization"
// @Param out-time header string false "out-time"
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

	sendChannelId, err := getSendChannelId(request)
	if err != nil {
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "未指定discord频道Id或未配置默认频道Id",
				Type:    "invalid_request_error",
				Code:    "discord_request_err",
			},
		})
		return
	}

	sentMsg, err := discord.SendMessage(sendChannelId, content)
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

	timer, err := setTimerWithHeader(c, request.Stream, common.RequestOutTimeDuration)
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
				timerReset(c, request.Stream, timer, common.RequestOutTimeDuration)

				// TODO 多张图片问题
				if !strings.HasPrefix(reply.Choices[0].Message.Content, strLen) {
					if len(strLen) > 3 && strings.HasPrefix(reply.Choices[0].Message.Content, "\\n1.") {
						strLen = strLen[:len(strLen)-2]
					} else {
						return true
					}
				}

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

// ImagesForOpenAI 图片生成-openai
// @Summary 图片生成-openai
// @Description 图片生成-openai
// @Tags openai
// @Accept json
// @Produce json
// @Param request body model.OpenAIImagesGenerationRequest true "request"
// @Param Authorization header string false "Authorization"
// @Param out-time header string false "out-time"
// @Success 200 {object} model.OpenAIImagesGenerationResponse "Successful response"
// @Router /v1/images/generations [post]
func ImagesForOpenAI(c *gin.Context) {

	var request model.OpenAIImagesGenerationRequest
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

	sendChannelId, err := getSendChannelId(request)
	if err != nil {
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "未指定discord频道Id或未配置默认频道Id",
				Type:    "invalid_request_error",
				Code:    "discord_request_err",
			},
		})
		return
	}

	sentMsg, err := discord.SendMessage(sendChannelId, request.Prompt)
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

	replyChan := make(chan model.OpenAIImagesGenerationResponse)
	discord.RepliesOpenAIImageChans[sentMsg.ID] = replyChan
	defer delete(discord.RepliesOpenAIImageChans, sentMsg.ID)

	stopChan := make(chan string)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	timer, err := setTimerWithHeader(c, false, common.RequestOutTimeDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "超时时间设置异常",
		})
		return
	}

	var replyResp model.OpenAIImagesGenerationResponse
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

func getSendChannelId(request model.ChannelIdentifier) (string, error) {
	var sendChannelId string
	channelId := request.GetChannelId()
	if channelId != nil && *channelId != "" {
		sendChannelId = *channelId
	} else {
		if discord.ChannelId != "" {
			sendChannelId = discord.ChannelId
		} else {
			return "", fmt.Errorf("未指定discord频道Id或未配置默认频道Id")
		}
	}
	return sendChannelId, nil
}

func setTimerWithHeader(c *gin.Context, isStream bool, defaultTimeout time.Duration) (*time.Timer, error) {

	outTimeStr := getOutTimeStr(c, isStream)

	if outTimeStr != "" {
		outTime, err := strconv.ParseInt(outTimeStr, 10, 64)
		if err != nil {
			return nil, err
		}
		return time.NewTimer(time.Duration(outTime) * time.Second), nil
	}
	return time.NewTimer(defaultTimeout), nil
}

func getOutTimeStr(c *gin.Context, isStream bool) string {
	var outTimeStr string
	if outTime := c.GetHeader(common.OutTime); outTime != "" {
		outTimeStr = outTime
	} else {
		if isStream {
			outTimeStr = common.StreamRequestOutTime
		} else {
			outTimeStr = common.RequestOutTime
		}
	}
	return outTimeStr
}

func timerReset(c *gin.Context, isStream bool, timer *time.Timer, defaultTimeout time.Duration) error {

	outTimeStr := getOutTimeStr(c, isStream)

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
