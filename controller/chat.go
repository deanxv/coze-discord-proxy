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
//func Chat(c *gin.Context) {
//
//	var chatModel model.ChatReq
//	err := json.NewDecoder(c.Request.Body).Decode(&chatModel)
//	if err != nil {
//		common.LogError(c.Request.Context(), err.Error())
//		c.JSON(http.StatusOK, gin.H{
//			"message": "无效的参数",
//			"success": false,
//		})
//		return
//	}
//
//	sendChannelId, calledCozeBotId, err := getSendChannelIdAndCozeBotId(c, false, chatModel)
//	if err != nil {
//		common.LogError(c.Request.Context(), err.Error())
//		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
//			OpenAIError: model.OpenAIError{
//				Message: "配置异常",
//				Type:    "invalid_request_error",
//				Code:    "discord_request_err",
//			},
//		})
//		return
//	}
//
//	sentMsg, err := discord.SendMessage(c, sendChannelId, calledCozeBotId, chatModel.Content)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"success": false,
//			"message": err.Error(),
//		})
//		return
//	}
//
//	replyChan := make(chan model.ReplyResp)
//	discord.RepliesChans[sentMsg.ID] = replyChan
//	defer delete(discord.RepliesChans, sentMsg.ID)
//
//	stopChan := make(chan model.ChannelStopChan)
//	discord.ReplyStopChans[sentMsg.ID] = stopChan
//	defer delete(discord.ReplyStopChans, sentMsg.ID)
//
//	timer, err := setTimerWithHeader(c, chatModel.Stream, common.RequestOutTimeDuration)
//	if err != nil {
//		common.LogError(c.Request.Context(), err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{
//			"success": false,
//			"message": "超时时间设置异常",
//		})
//		return
//	}
//
//	if chatModel.Stream {
//		c.Stream(func(w io.Writer) bool {
//			select {
//			case reply := <-replyChan:
//				timerReset(c, chatModel.Stream, timer, common.RequestOutTimeDuration)
//				urls := ""
//				if len(reply.EmbedUrls) > 0 {
//					for _, url := range reply.EmbedUrls {
//						urls += "\n" + fmt.Sprintf("![Image](%s)", url)
//					}
//				}
//				c.SSEvent("message", reply.Content+urls)
//				return true // 继续保持流式连接
//			case <-timer.C:
//				// 定时器到期时,关闭流
//				return false
//			case <-stopChan:
//				return false // 关闭流式连接
//			}
//		})
//	} else {
//		var replyResp model.ReplyResp
//		for {
//			select {
//			case reply := <-replyChan:
//				replyResp.Content = reply.Content
//				replyResp.EmbedUrls = reply.EmbedUrls
//			case <-timer.C:
//				c.JSON(http.StatusOK, gin.H{
//					"success": false,
//					"message": "request_out_time",
//				})
//				return
//			case <-stopChan:
//				c.JSON(http.StatusOK, gin.H{
//					"success": true,
//					"data":    replyResp,
//				})
//				return
//			}
//		}
//	}
//}

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
		common.LogError(c.Request.Context(), err.Error())
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	sendChannelId, calledCozeBotId, isNewChannel, err := getSendChannelIdAndCozeBotId(c, request.ChannelId, request.Model, true)

	if err != nil {
		common.LogError(c.Request.Context(), err.Error())
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "config error,check logs",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	if isNewChannel {
		defer func() {
			if discord.ChannelAutoDelTime != "" {
				delTime, _ := strconv.Atoi(discord.ChannelAutoDelTime)
				if delTime == 0 {
					discord.CancelChannelDeleteTimer(sendChannelId)
				} else if delTime > 0 {
					// 删除该频道
					discord.SetChannelDeleteTimer(sendChannelId, time.Duration(delTime)*time.Second)
				}
			} else {
				// 删除该频道
				discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
			}
		}()
	}

	content := "Hi！"
	messages := request.Messages

loop:
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		if message.Role == "user" {
			switch contentObj := message.Content.(type) {
			case string:
				if common.AllDialogRecordEnable == "1" || common.AllDialogRecordEnable == "" {
					messages[i] = model.OpenAIChatMessage{
						Role:    "user",
						Content: contentObj,
					}
				} else {
					content = contentObj
					break loop
				}
			case []interface{}:
				content, err = buildOpenAIGPT4VForImageContent(sendChannelId, contentObj)
				if err != nil {
					c.JSON(http.StatusOK, model.OpenAIErrorResponse{
						OpenAIError: model.OpenAIError{
							Message: "Image URL parsing error",
							Type:    "request_error",
							Code:    "500",
						},
					})
					return
				}
				if common.AllDialogRecordEnable == "1" || common.AllDialogRecordEnable == "" {
					messages[i] = model.OpenAIChatMessage{
						Role:    "user",
						Content: content,
					}
				} else {
					break loop
				}
			default:
				c.JSON(http.StatusOK, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: "Message format error",
						Type:    "request_error",
						Code:    "500",
					},
				})
				return

			}
			//break
		} else {
			messages[i] = model.OpenAIChatMessage{
				Role:    message.Role,
				Content: message.Content,
			}
		}
	}

	if common.AllDialogRecordEnable == "1" || common.AllDialogRecordEnable == "" {
		jsonData, err := json.Marshal(messages)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		content = string(jsonData)
	}

	//for i := len(messages) - 1; i >= 0; i-- {
	//	message := messages[i]
	//	if message.Role == "user" {
	//		switch contentObj := message.Content.(type) {
	//		case string:
	//			if common.AllDialogRecordEnable == "1" {
	//				content = contentObj
	//			} else {
	//				jsonData, err := json.Marshal(messages)
	//				if err != nil {
	//					c.JSON(http.StatusOK, gin.H{
	//						"success": false,
	//						"message": err.Error(),
	//					})
	//					return
	//				}
	//				content = string(jsonData)
	//			}
	//		case []interface{}:
	//			content, err = buildOpenAIGPT4VForImageContent(sendChannelId, contentObj)
	//			if err != nil {
	//				c.JSON(http.StatusOK, gin.H{
	//					"success": false,
	//					"message": err.Error(),
	//				})
	//				return
	//			}
	//		default:
	//			c.JSON(http.StatusOK, model.OpenAIErrorResponse{
	//				OpenAIError: model.OpenAIError{
	//					Message: "消息格式异常",
	//					Type:    "invalid_request_error",
	//					Code:    "discord_request_err",
	//				},
	//			})
	//			return
	//
	//		}
	//		break
	//	}
	//}

	sentMsg, userAuth, err := discord.SendMessage(c, sendChannelId, calledCozeBotId, content)
	if err != nil {
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: err.Error(),
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	replyChan := make(chan model.OpenAIChatCompletionResponse)
	discord.RepliesOpenAIChans[sentMsg.ID] = replyChan
	defer delete(discord.RepliesOpenAIChans, sentMsg.ID)

	stopChan := make(chan model.ChannelStopChan)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	timer, err := setTimerWithHeader(c, request.Stream, common.RequestOutTimeDuration)
	if err != nil {
		common.LogError(c.Request.Context(), err.Error())
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

				if common.SliceContains(common.CozeErrorMessages, reply.Choices[0].Message.Content) {
					if common.SliceContains(common.CozeDailyLimitErrorMessages, reply.Choices[0].Message.Content) {
						common.LogWarn(c, fmt.Sprintf("USER_AUTHORIZATION:%s DAILY LIMIT", userAuth))
						discord.UserAuthorizations = common.FilterSlice(discord.UserAuthorizations, userAuth)
					}
					//discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
					c.SSEvent("", " [DONE]")
					return false // 关闭流式连接
				}

				return true // 继续保持流式连接
			case <-timer.C:
				// 定时器到期时,关闭流
				//discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
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
				if common.SliceContains(common.CozeErrorMessages, reply.Choices[0].Message.Content) {
					if common.SliceContains(common.CozeDailyLimitErrorMessages, reply.Choices[0].Message.Content) {
						common.LogWarn(c, fmt.Sprintf("USER_AUTHORIZATION:%s DAILY LIMIT", userAuth))
						discord.UserAuthorizations = common.FilterSlice(discord.UserAuthorizations, userAuth)
					}
					//discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
					c.JSON(http.StatusOK, model.OpenAIErrorResponse{
						OpenAIError: model.OpenAIError{
							Message: reply.Choices[0].Message.Content,
							Type:    "request_error",
							Code:    "500",
						},
					})
					return
				}
				replyResp = reply
			case <-timer.C:
				//discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
				c.JSON(http.StatusOK, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: "Request timeout",
						Type:    "request_error",
						Code:    "500",
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

func buildOpenAIGPT4VForImageContent(sendChannelId string, objs []interface{}) (string, error) {
	var content string

	for i, obj := range objs {

		jsonData, err := json.Marshal(obj)
		if err != nil {
			return "", err
		}

		var req model.OpenAIGPT4VImagesReq
		err = json.Unmarshal(jsonData, &req)
		if err != nil {
			return "", err
		}

		if i == 0 && req.Type == "text" {
			content += req.Text
			continue
		} else if i != 0 && req.Type == "image_url" {
			if common.IsURL(req.ImageURL.URL) {
				content += fmt.Sprintf("\n%s ", req.ImageURL.URL)
			} else if common.IsImageBase64(req.ImageURL.URL) {
				url, err := discord.UploadToDiscordAndGetURL(sendChannelId, req.ImageURL.URL)
				if err != nil {
					return "", fmt.Errorf("文件上传异常")
				}
				content += fmt.Sprintf("\n%s ", url)
			} else {
				return "", fmt.Errorf("文件格式有误")
			}
		} else {
			return "", fmt.Errorf("消息格式错误")
		}
	}
	//if runeCount := len([]rune(content)); runeCount > 2000 {
	//	return "", fmt.Errorf("prompt最大为2000字符 [%v]", runeCount)
	//}
	return content, nil

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
		common.LogError(c.Request.Context(), err.Error())
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	if runeCount := len([]rune(request.Prompt)); runeCount > 2000 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("prompt最大为2000字符 [%v]", runeCount),
		})
		return
	}

	sendChannelId, calledCozeBotId, isNewChannel, err := getSendChannelIdAndCozeBotId(c, request.ChannelId, request.Model, true)
	if err != nil {
		common.LogError(c.Request.Context(), err.Error())
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "config error",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	if isNewChannel {
		defer func() {
			if discord.ChannelAutoDelTime != "" {
				delTime, _ := strconv.Atoi(discord.ChannelAutoDelTime)
				if delTime == 0 {
					discord.CancelChannelDeleteTimer(sendChannelId)
				} else if delTime > 0 {
					// 删除该频道
					discord.SetChannelDeleteTimer(sendChannelId, time.Duration(delTime)*time.Second)
				}
			} else {
				// 删除该频道
				discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
			}
		}()
	}

	sentMsg, userAuth, err := discord.SendMessage(c, sendChannelId, calledCozeBotId, request.Prompt)
	if err != nil {
		c.JSON(http.StatusOK, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: err.Error(),
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	replyChan := make(chan model.OpenAIImagesGenerationResponse)
	discord.RepliesOpenAIImageChans[sentMsg.ID] = replyChan
	defer delete(discord.RepliesOpenAIImageChans, sentMsg.ID)

	stopChan := make(chan model.ChannelStopChan)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	timer, err := setTimerWithHeader(c, false, common.RequestOutTimeDuration)
	if err != nil {
		common.LogError(c.Request.Context(), err.Error())
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
			if reply.DailyLimit {
				common.LogWarn(c, fmt.Sprintf("USER_AUTHORIZATION:%s DAILY LIMIT", userAuth))
				discord.UserAuthorizations = common.FilterSlice(discord.UserAuthorizations, userAuth)
				//discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
				c.JSON(http.StatusOK, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: "daily limit for sending messages",
						Type:    "request_error",
						Code:    "500",
					},
				})
				return
			}
			replyResp = reply
		case <-timer.C:
			//discord.SetChannelDeleteTimer(sendChannelId, 5*time.Second)
			c.JSON(http.StatusOK, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: "Request timed out, please try again later.",
					Type:    "request_error",
					Code:    "500",
				},
			})
			return
		case <-stopChan:
			if replyResp.Data == nil {
				c.JSON(http.StatusOK, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: "Failed to fetch image URL, please try again later.",
						Type:    "request_error",
						Code:    "500",
					},
				})
				return
			}
			c.JSON(http.StatusOK, replyResp)
			return
		}
	}

}

func getSendChannelIdAndCozeBotId(c *gin.Context, channelId *string, model string, isOpenAIAPI bool) (sendChannelId string, calledCozeBotId string, isNewChannel bool, err error) {
	secret := ""
	if isOpenAIAPI {
		if secret = c.Request.Header.Get("Authorization"); secret != "" {
			secret = strings.Replace(secret, "Bearer ", "", 1)
		}
	} else {
		secret = c.Request.Header.Get("proxy-secret")
	}

	// botConfigs不为空
	if len(discord.BotConfigList) != 0 {

		botConfigs := discord.FilterConfigs(discord.BotConfigList, secret, model, nil)
		if len(botConfigs) != 0 {
			// 有值则随机一个
			botConfig, err := common.RandomElement(botConfigs)
			if err != nil {
				return "", "", false, err
			}

			if channelId != nil && *channelId != "" {
				return *channelId, botConfig.CozeBotId, false, nil
			}

			if discord.DefaultChannelEnable == "1" {
				return botConfig.ChannelId, botConfig.CozeBotId, false, nil
			} else {
				var sendChannelId string
				sendChannelId, err := discord.CreateChannelWithRetry(c, discord.GuildId, fmt.Sprintf("cdp-对话%s", c.Request.Context().Value(common.RequestIdKey)), 0)
				if err != nil {
					common.LogError(c, err.Error())
					return "", "", false, err
				}
				return sendChannelId, botConfig.CozeBotId, true, nil
			}

		}
		// 没有值抛出异常
		return "", "", false, fmt.Errorf("[proxy-secret]+[model]未匹配到有效bot")
	} else {

		if channelId != nil && *channelId != "" {
			return *channelId, discord.CozeBotId, false, nil
		}

		if discord.DefaultChannelEnable == "1" {
			return discord.ChannelId, discord.CozeBotId, false, nil
		} else {
			sendChannelId, err := discord.CreateChannelWithRetry(c, discord.GuildId, fmt.Sprintf("cdp-对话%s", c.Request.Context().Value(common.RequestIdKey)), 0)
			if err != nil {
				//common.LogError(c, err.Error())
				return "", "", false, err
			}
			return sendChannelId, discord.CozeBotId, true, nil
		}
	}
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
