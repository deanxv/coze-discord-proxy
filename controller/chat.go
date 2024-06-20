package controller

import (
	"coze-discord-proxy/common"
	"coze-discord-proxy/common/config"
	"coze-discord-proxy/common/myerr"
	"coze-discord-proxy/discord"
	"coze-discord-proxy/model"
	"coze-discord-proxy/telegram"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
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
//	myerr := json.NewDecoder(c.Request.Body).Decode(&chatModel)
//	if myerr != nil {
//		common.LogError(c.Request.Context(), myerr.Error())
//		c.JSON(http.StatusOK, gin.H{
//			"message": "无效的参数",
//			"success": false,
//		})
//		return
//	}
//
//	sendChannelId, calledCozeBotId, myerr := getSendChannelIdAndCozeBotId(c, false, chatModel)
//	if myerr != nil {
//		common.LogError(c.Request.Context(), myerr.Error())
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
//	sentMsg, myerr := discord.SendMessage(c, sendChannelId, calledCozeBotId, chatModel.Content)
//	if myerr != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"success": false,
//			"message": myerr.Error(),
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
//	timer, myerr := setTimerWithHeader(c, chatModel.Stream, config.RequestOutTimeDuration)
//	if myerr != nil {
//		common.LogError(c.Request.Context(), myerr.Error())
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
//				timerReset(c, chatModel.Stream, timer, config.RequestOutTimeDuration)
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
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	if err := checkUserAuths(c); err != nil {
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: err.Error(),
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	sendChannelId, calledCozeBotId, isNewChannel, err := getSendChannelIdAndCozeBotId(c, request.ChannelId, request.Model, true)

	if err != nil {
		response := model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "config error,check logs",
				Type:    "request_error",
				Code:    "500",
			},
		}
		common.LogError(c.Request.Context(), err.Error())
		var myErr *myerr.ModelNotFoundError
		if errors.As(err, &myErr) {
			response.OpenAIError.Message = "model_not_found"
		}
		c.JSON(http.StatusInternalServerError, response)
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
				if config.AllDialogRecordEnable == "1" || config.AllDialogRecordEnable == "" {
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
					c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
						OpenAIError: model.OpenAIError{
							Message: "Image URL parsing error",
							Type:    "request_error",
							Code:    "500",
						},
					})
					return
				}
				if config.AllDialogRecordEnable == "1" || config.AllDialogRecordEnable == "" {
					messages[i] = model.OpenAIChatMessage{
						Role:    "user",
						Content: content,
					}
				} else {
					break loop
				}
			default:
				c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
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

	if config.AllDialogRecordEnable == "1" || config.AllDialogRecordEnable == "" {
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

	sentMsg, userAuth, err := discord.SendMessage(c, sendChannelId, calledCozeBotId, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: err.Error(),
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	replyChan := make(chan model.OpenAIChatCompletionResponse)
	discord.RepliesOpenAIChans[sentMsg.ID] = &model.OpenAIChatCompletionChan{
		Model:    request.Model,
		Response: replyChan,
	}
	defer delete(discord.RepliesOpenAIChans, sentMsg.ID)

	stopChan := make(chan model.ChannelStopChan)
	discord.ReplyStopChans[sentMsg.ID] = stopChan
	defer delete(discord.ReplyStopChans, sentMsg.ID)

	timer, err := setTimerWithHeader(c, request.Stream, config.RequestOutTimeDuration)
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
				timerReset(c, request.Stream, timer, config.RequestOutTimeDuration)

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
					c.SSEvent("", " [DONE]")
					return false // 关闭流式连接
				}

				return true // 继续保持流式连接
			case <-timer.C:
				// 定时器到期时,关闭流
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
					c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
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
				c.JSON(http.StatusOK, replyResp)
				return
			case <-stopChan:
				c.JSON(http.StatusOK, replyResp)
				return
			}
		}
	}
}

// OpenaiModels 模型列表-openai
// @Summary 模型列表-openai
// @Description 模型列表-openai
// @Tags openai
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization"
// @Success 200 {object} model.OpenaiModelListResponse "Successful response"
// @Router /v1/models [get]
func OpenaiModels(c *gin.Context) {
	var modelsResp []string

	secret := ""
	if len(discord.BotConfigList) != 0 {
		if secret = c.Request.Header.Get("Authorization"); secret != "" {
			secret = strings.Replace(secret, "Bearer ", "", 1)
		}

		botConfigs := discord.FilterConfigs(discord.BotConfigList, secret, "", nil)
		for _, botConfig := range botConfigs {
			modelsResp = append(modelsResp, botConfig.Model...)
		}

		modelsResp = lo.Uniq(modelsResp)
	} else {
		modelsResp = common.DefaultOpenaiModelList
	}

	var openaiModelListResponse model.OpenaiModelListResponse
	var openaiModelResponse []model.OpenaiModelResponse
	openaiModelListResponse.Object = "list"

	for _, modelResp := range modelsResp {
		openaiModelResponse = append(openaiModelResponse, model.OpenaiModelResponse{
			ID:     modelResp,
			Object: "model",
		})
	}
	openaiModelListResponse.Data = openaiModelResponse
	c.JSON(http.StatusOK, openaiModelListResponse)
	return
}

func buildOpenAIGPT4VForImageContent(sendChannelId string, objs []interface{}) (string, error) {
	var content string
	var url string

	for _, obj := range objs {

		jsonData, err := json.Marshal(obj)
		if err != nil {
			return "", err
		}

		var req model.OpenAIGPT4VImagesReq
		err = json.Unmarshal(jsonData, &req)
		if err != nil {
			return "", err
		}

		if req.Type == "text" {
			content = req.Text
		} else if req.Type == "image_url" {
			if common.IsURL(req.ImageURL.URL) {
				url = fmt.Sprintf("%s ", req.ImageURL.URL)
			} else if common.IsImageBase64(req.ImageURL.URL) {
				imgUrl, err := discord.UploadToDiscordAndGetURL(sendChannelId, req.ImageURL.URL)
				if err != nil {
					return "", fmt.Errorf("文件上传异常")
				}
				url = fmt.Sprintf("\n%s ", imgUrl)
			} else {
				return "", fmt.Errorf("文件格式有误")
			}
		} else {
			return "", fmt.Errorf("消息格式错误")
		}
	}

	return fmt.Sprintf("%s\n%s", content, url), nil

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
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	if err := checkUserAuths(c); err != nil {
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: err.Error(),
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	sendChannelId, calledCozeBotId, isNewChannel, err := getSendChannelIdAndCozeBotId(c, request.ChannelId, request.Model, true)
	if err != nil {
		common.LogError(c.Request.Context(), err.Error())
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
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

	sentMsg, userAuth, err := discord.SendMessage(c, sendChannelId, calledCozeBotId, common.ImgGeneratePrompt+request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
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

	timer, err := setTimerWithHeader(c, false, config.RequestOutTimeDuration)
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
				c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: "daily limit for sending messages",
						Type:    "request_error",
						Code:    "500",
					},
				})
				return
			}
			if request.ResponseFormat == "b64_json" && reply.Data != nil && len(reply.Data) > 0 {
				for _, data := range reply.Data {
					if data.URL != "" {
						base64Str, err := getBase64ByUrl(data.URL)
						if err != nil {
							c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
								OpenAIError: model.OpenAIError{
									Message: err.Error(),
									Type:    "request_error",
									Code:    "500",
								},
							})
							return
						}
						data.B64Json = "data:image/webp;base64," + base64Str
					}
				}
			}
			replyResp = reply
		case <-timer.C:
			if replyResp.Data == nil {
				c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
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
		case <-stopChan:
			if replyResp.Data == nil {
				c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
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
				sendChannelId, err := discord.CreateChannelWithRetry(c, discord.GuildId, fmt.Sprintf("cdp-chat-%s", c.Request.Context().Value(common.RequestIdKey)), 0)
				if err != nil {
					common.LogError(c, err.Error())
					return "", "", false, err
				}
				return sendChannelId, botConfig.CozeBotId, true, nil
			}

		}
		// 没有值抛出异常
		return "", "", false, &myerr.ModelNotFoundError{
			ErrCode: 500,
			Message: fmt.Sprintf("[proxy-secret:%s]+[model:%s]未匹配到有效bot", secret, model),
		}
	} else {

		if channelId != nil && *channelId != "" {
			return *channelId, discord.CozeBotId, false, nil
		}

		if discord.DefaultChannelEnable == "1" {
			return discord.ChannelId, discord.CozeBotId, false, nil
		} else {
			sendChannelId, err := discord.CreateChannelWithRetry(c, discord.GuildId, fmt.Sprintf("cdp-chat-%s", c.Request.Context().Value(common.RequestIdKey)), 0)
			if err != nil {
				//common.LogError(c, myerr.Error())
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
			outTimeStr = config.StreamRequestOutTime
		} else {
			outTimeStr = config.RequestOutTime
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

func checkUserAuths(c *gin.Context) error {
	if len(discord.UserAuthorizations) == 0 {
		common.LogError(c, fmt.Sprintf("无可用的 user_auth"))
		// tg发送通知
		if !common.IsSameDay(discord.NoAvailableUserAuthPreNotifyTime, time.Now()) && telegram.NotifyTelegramBotToken != "" && telegram.TgBot != nil {
			go func() {
				discord.NoAvailableUserAuthChan <- "stop"
			}()
		}

		return fmt.Errorf("no_available_user_auth")
	}
	return nil
}

func getBase64ByUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Encode the image data to Base64
	base64Str := base64.StdEncoding.EncodeToString(imgData)
	return base64Str, nil
}
