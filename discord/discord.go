package discord

import (
	"context"
	"coze-discord-proxy/common"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var cozeBotId = os.Getenv("COZE_BOT_ID")
var GuildId = os.Getenv("GUILD_ID")

var RepliesChans = make(map[string]chan string)
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

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 检查消息是否是对我们的回复
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			if replyChan, exists := RepliesChans[m.ReferencedMessage.ID]; exists {
				replyChan <- m.Content
				if len(m.Message.Components) > 0 {
					stopChan := ReplyStopChans[m.ReferencedMessage.ID]
					stopChan <- m.ReferencedMessage.ID
				}
			}
			break
		}
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
	fmt.Println("频道创建成功")
	return st.ID, nil
}

func ChannelDel(channelId string) (string, error) {
	// 创建新的频道
	st, err := Session.ChannelDelete(channelId)
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("删除频道时异常 %s", err.Error()))
		return "", err
	}
	fmt.Println("删除成功")
	return st.ID, nil
}
