package discord

import (
	"context"
	"coze-discord-proxy/common"
	"coze-discord-proxy/model"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cozeBotId = os.Getenv("COZE_BOT_ID")
var GuildId = os.Getenv("GUILD_ID")
var ChannelId = os.Getenv("CHANNEL_ID")

var RepliesChans = make(map[string]chan model.ReplyResp)
var RepliesOpenAIChans = make(map[string]chan model.OpenAIChatCompletionResponse)

var ReplyStopChans = make(map[string]chan string)
var Session *discordgo.Session

func StartBot(ctx context.Context, token string) {
	var err error
	Session, err = discordgo.New("Bot " + token)
	if err != nil {
		common.FatalLog("error creating Discord session,", err)
		return
	}

	// 注册消息处理函数
	Session.AddHandler(messageUpdate)

	// 打开websocket连接并开始监听
	err = Session.Open()
	if err != nil {
		common.FatalLog("error opening connection,", err)
		return
	}

	common.LogInfo(ctx, "Bot is now running. Press CTRL+C to exit.")

	go func() {
		<-ctx.Done()
		if err := Session.Close(); err != nil {
			common.FatalLog("error closing Discord session,", err)
		}
	}()

	// 等待信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// messageUpdate handles the updated messages in Discord.
func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// 提前检查参考消息是否为 nil
	if m.ReferencedMessage == nil {
		return
	}

	// 尝试获取 stopChan
	stopChan, exists := ReplyStopChans[m.ReferencedMessage.ID]
	if !exists {
		return
	}

	// 如果作者为 nil 或消息来自 bot 本身，则发送停止信号
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		stopChan <- m.ReferencedMessage.ID
		return
	}

	// 检查消息是否是对 bot 的回复
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			replyChan, exists := RepliesChans[m.ReferencedMessage.ID]
			if exists {
				reply := processMessage(m)
				replyChan <- reply
			} else {
				replyOpenAIChan, exists := RepliesOpenAIChans[m.ReferencedMessage.ID]
				if exists {
					reply := processMessageForOpenAI(m)
					replyOpenAIChan <- reply
				} else {
					return
				}
			}
			// data: {"id":"chatcmpl-8lho2xvdDFyBdFkRwWAcMpWWAgymJ","object":"chat.completion.chunk","created":1706380498,"model":"gpt-3.5-turbo-0613","system_fingerprint":null,"choices":[{"index":0,"delta":{"content":"？"},"logprobs":null,"finish_reason":null}]}
			// data :{"id":"1200873365351698694","object":"chat.completion.chunk","created":1706380922,"model":"COZE","choices":[{"index":0,"message":{"role":"assistant","content":"你好！有什么我可以帮您的吗？如果有任"},"logprobs":null,"finish_reason":"","delta":{"content":"吗？如果有任"}}],"usage":{"prompt_tokens":13,"completion_tokens":19,"total_tokens":32},"system_fingerprint":null}

			// 如果消息包含组件或嵌入，则发送停止信号
			if len(m.Embeds) > 0 || len(m.Message.Components) > 0 {
				replyOpenAIChan, exists := RepliesOpenAIChans[m.ReferencedMessage.ID]
				if exists {
					reply := processMessageForOpenAI(m)
					stopStr := "stop"
					reply.Choices[0].FinishReason = &stopStr
					replyOpenAIChan <- reply
				}
				stopChan <- m.ReferencedMessage.ID
			}

			return
		}
	}
}

// processMessage 提取并处理消息内容及其嵌入元素
func processMessage(m *discordgo.MessageUpdate) model.ReplyResp {
	var embedUrls []string
	for _, embed := range m.Embeds {
		if embed.Image != nil {
			embedUrls = append(embedUrls, embed.Image.URL)
		}
	}

	return model.ReplyResp{
		Content:   m.Content,
		EmbedUrls: embedUrls,
	}
}

func processMessageForOpenAI(m *discordgo.MessageUpdate) model.OpenAIChatCompletionResponse {
	var embedUrls []string
	for _, embed := range m.Embeds {
		if embed.Image != nil {
			embedUrls = append(embedUrls, embed.Image.URL)
		}
	}

	promptTokens := common.CountTokens(m.ReferencedMessage.Content)
	completionTokens := common.CountTokens(m.Content)

	return model.OpenAIChatCompletionResponse{
		ID:      m.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "COZE",
		Choices: []model.OpenAIChoice{
			{
				Index: 0,
				Message: model.OpenAIMessage{
					Role:    "assistant",
					Content: m.Content,
				},
			},
		},
		Usage: model.OpenAIUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}
}

func SendMessage(channelID, message string) (*discordgo.Message, error) {
	if Session == nil {
		return nil, fmt.Errorf("Discord session not initialized")
	}

	// 添加@机器人逻辑
	sentMsg, err := Session.ChannelMessageSend(channelID, fmt.Sprintf("<@%s> %s", cozeBotId, message))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %s", err)
	}
	return sentMsg, nil
}

func ChannelCreate(guildID, channelName string) (string, error) {
	// 创建新的频道
	st, err := Session.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("创建频道时异常 %s", err.Error()))
		return "", err
	}
	return st.ID, nil
}

func ChannelDel(channelId string) (string, error) {
	// 创建新的频道
	st, err := Session.ChannelDelete(channelId)
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("删除频道时异常 %s", err.Error()))
		return "", err
	}
	return st.ID, nil
}

func ChannelCreateComplex(guildID, parentId, channelName string) (string, error) {
	// 创建新的子频道
	st, err := Session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: parentId,
	})
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("创建子频道时异常 %s", err.Error()))
		return "", err
	}
	return st.ID, nil
}

func ThreadStart(channelId, threadName string, archiveDuration int) (string, error) {
	// 创建新的线程
	th, err := Session.ThreadStart(channelId, threadName, discordgo.ChannelTypeGuildText, archiveDuration)

	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("创建线程时异常 %s", err.Error()))
		return "", err
	}
	return th.ID, nil
}
