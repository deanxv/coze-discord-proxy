package discord

import (
	"bytes"
	"context"
	"coze-discord-proxy/common"
	"coze-discord-proxy/model"
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/h2non/filetype"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var cozeBotId = os.Getenv("COZE_BOT_ID")
var GuildId = os.Getenv("GUILD_ID")
var ChannelId = os.Getenv("CHANNEL_ID")
var ProxyUrl = os.Getenv("PROXY_URL")

var RepliesChans = make(map[string]chan model.ReplyResp)
var RepliesOpenAIChans = make(map[string]chan model.OpenAIChatCompletionResponse)
var RepliesOpenAIImageChans = make(map[string]chan model.OpenAIImagesGenerationResponse)

var ReplyStopChans = make(map[string]chan string)
var Session *discordgo.Session

func StartBot(ctx context.Context, token string) {
	var err error
	Session, err = discordgo.New("Bot " + token)

	if ProxyUrl != "" {
		client, err := NewProxyClient(ProxyUrl)
		if err != nil {
			common.FatalLog("error creating proxy client,", err)
		}
		Session.Client = client
		common.LogInfo(context.Background(), "Proxy Set Success")
	}

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

	common.LogInfo(ctx, "Bot is now running. Enjoy It.")

	go scheduleDailyMessage(ChannelId, "Hi!")

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
					replyOpenAIImageChan, exists := RepliesOpenAIImageChans[m.ReferencedMessage.ID]
					if exists {
						reply := processMessageForOpenAIImage(m)
						replyOpenAIImageChan <- reply
					} else {
						return
					}
				}
			}
			// data: {"id":"chatcmpl-8lho2xvdDFyBdFkRwWAcMpWWAgymJ","object":"chat.completion.chunk","created":1706380498,"model":"gpt-3.5-turbo-0613","system_fingerprint":null,"choices":[{"index":0,"delta":{"content":"？"},"logprobs":null,"finish_reason":null}]}
			// data :{"id":"1200873365351698694","object":"chat.completion.chunk","created":1706380922,"model":"COZE","choices":[{"index":0,"message":{"role":"assistant","content":"你好！有什么我可以帮您的吗？如果有任"},"logprobs":null,"finish_reason":"","delta":{"content":"吗？如果有任"}}],"usage":{"prompt_tokens":13,"completion_tokens":19,"total_tokens":32},"system_fingerprint":null}

			// 如果消息包含组件或嵌入，则发送停止信号
			if len(m.Message.Components) > 0 {
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

	if len(m.Embeds) != 0 {
		for _, embed := range m.Embeds {
			if embed.Image != nil && !strings.Contains(m.Content, embed.Image.URL) {
				if m.Content != "" {
					m.Content += "\n"
				}
				m.Content += fmt.Sprintf("%s\n![Image](%s)", embed.Image.URL, embed.Image.URL)
			}
		}
	}

	promptTokens := common.CountTokens(m.ReferencedMessage.Content)
	completionTokens := common.CountTokens(m.Content)

	return model.OpenAIChatCompletionResponse{
		ID:      m.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-3.5-turbo",
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

func processMessageForOpenAIImage(m *discordgo.MessageUpdate) model.OpenAIImagesGenerationResponse {
	var response model.OpenAIImagesGenerationResponse

	if len(m.Embeds) != 0 {
		for _, embed := range m.Embeds {
			if embed.Image != nil && !strings.Contains(m.Content, embed.Image.URL) {
				if m.Content != "" {
					m.Content += "\n"
				}
				response.Data = append(response.Data, struct {
					URL string `json:"url"`
				}{URL: embed.Image.URL})
			}
		}
	}

	return model.OpenAIImagesGenerationResponse{
		Created: time.Now().Unix(),
		Data:    response.Data,
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

func NewProxyClient(proxyUrl string) (*http.Client, error) {

	proxyParse, err := url.Parse(proxyUrl)
	if err != nil {
		common.FatalLog("代理地址设置有误")
	}

	if strings.HasPrefix(proxyParse.Scheme, "http") {
		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(proxyParse),
		}
		return &http.Client{
			Transport: httpTransport,
		}, nil
	} else if strings.HasPrefix(proxyParse.Scheme, "sock") {
		dialer, err := proxy.SOCKS5("tcp", proxyParse.Host, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, ", err)
		}

		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}

		// 使用该拨号器创建一个 HTTP 客户端
		httpClient := &http.Client{
			Transport: &http.Transport{
				DialContext: dialContext,
			},
		}

		return httpClient, nil
	} else {
		return nil, fmt.Errorf("仅支持sock和http代理！")
	}

}

func scheduleDailyMessage(channelID string, message string) {
	for {
		// 计算距离下一个晚上12点的时间间隔
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		delay := next.Sub(now)

		// 等待直到下一个间隔
		time.Sleep(delay)

		_, err := SendMessage(channelID, message)

		if err != nil {
			common.LogWarn(context.Background(), "活跃机器人任务消息发送异常!")
		} else {
			common.LogInfo(context.Background(), "活跃机器人任务消息发送成功!")
		}
	}
}

func UploadToDiscordAndGetURL(channelID string, base64Data string) (string, error) {

	// 获取";base64,"后的Base64编码部分
	dataParts := strings.Split(base64Data, ";base64,")
	if len(dataParts) != 2 {
		return "", fmt.Errorf("")
	}
	base64Data = dataParts[1]

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}
	// 创建一个新的文件读取器
	file := bytes.NewReader(data)

	kind, err := filetype.Match(data)

	if err != nil {
		return "", fmt.Errorf("无法识别的文件类型")
	}

	// 创建一个新的 MessageSend 结构
	m := &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   fmt.Sprintf("image-%s.%s", common.GetTimeString(), kind.Extension),
				Reader: file,
			},
		},
	}

	// 发送消息
	message, err := Session.ChannelMessageSendComplex(channelID, m)
	if err != nil {
		return "", err
	}

	// 检查消息中是否包含附件，并获取 URL
	if len(message.Attachments) > 0 {
		return message.Attachments[0].URL, nil
	}

	return "", fmt.Errorf("no attachment found in the message")
}
