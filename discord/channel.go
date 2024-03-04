package discord

import (
	"context"
	"coze-discord-proxy/common"
	"coze-discord-proxy/telegram"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"time"
)

var (
	channelTimers sync.Map // 用于存储频道ID和对应的定时器
)

// SetChannelDeleteTimer 设置或重置频道的删除定时器
func SetChannelDeleteTimer(channelId string, duration time.Duration) {

	channel, err := Session.Channel(channelId)
	// 非自动生成频道不删除
	if err == nil && !strings.HasPrefix(channel.Name, "cdp-对话") {
		return
	}

	// 过滤掉配置中的频道id
	for _, config := range BotConfigList {
		if config.ChannelId == channelId {
			return
		}
	}

	if ChannelId == channelId {
		return
	}

	// 检查是否已存在定时器
	if timer, ok := channelTimers.Load(channelId); ok {
		if timer.(*time.Timer).Stop() {
			// 仅当定时器成功停止时才从映射中删除
			channelTimers.Delete(channelId)
		}
	}

	// 设置新的定时器
	newTimer := time.AfterFunc(duration, func() {
		ChannelDel(channelId)
		// 删除完成后从map中移除
		channelTimers.Delete(channelId)
	})
	// 存储新的定时器
	channelTimers.Store(channelId, newTimer)
}

// CancelChannelDeleteTimer 取消频道的删除定时器
func CancelChannelDeleteTimer(channelId string) {
	// 尝试从映射中获取定时器
	if timer, ok := channelTimers.Load(channelId); ok {
		// 如果定时器存在，尝试停止它
		if timer.(*time.Timer).Stop() {
			// 定时器成功停止后，从映射中移除
			channelTimers.Delete(channelId)
		} else {
			common.SysError(fmt.Sprintf("定时器无法停止或已触发，频道可能已被删除:%s", channelId))
		}
	} else {
		common.SysError(fmt.Sprintf("频道无定时删除:%s", channelId))
	}
}

func ChannelCreate(guildID, channelName string, channelType int) (string, error) {
	// 创建新的频道
	st, err := Session.GuildChannelCreate(guildID, channelName, discordgo.ChannelType(channelType))
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("创建频道时异常 %s", err.Error()))
		return "", err
	}
	return st.ID, nil
}

func ChannelDel(channelId string) (string, error) {
	// 删除频道
	st, err := Session.ChannelDelete(channelId)
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("删除频道时异常 %s", err.Error()))
		return "", err
	}
	return st.ID, nil
}

func ChannelCreateComplex(guildID, parentId, channelName string, channelType int) (string, error) {
	// 创建新的子频道
	st, err := Session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelType(channelType),
		ParentID: parentId,
	})
	if err != nil {
		common.LogError(context.Background(), fmt.Sprintf("创建子频道时异常 %s", err.Error()))
		return "", err
	}
	return st.ID, nil
}

type channelCreateResult struct {
	ID  string
	Err error
}

func CreateChannelWithRetry(c *gin.Context, guildID, channelName string, channelType int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		resultCh := make(chan channelCreateResult)

		go func() {
			channelID, err := ChannelCreate(guildID, channelName, channelType)
			resultCh <- channelCreateResult{ID: channelID, Err: err}
		}()

		select {
		case result := <-resultCh:
			if result.Err != nil {
				common.LogWarn(c, fmt.Sprintf("Failed to create channel, error: %v", result.Err))
				continue
			}
			return result.ID, nil
		case <-ctx.Done():
			common.LogWarn(c, "Create channel timed out, retrying...")
		}
	}
	// tg发送通知
	if telegram.NotifyTelegramBotToken != "" && telegram.TgBot != nil {
		go func() {
			CreateChannelRiskChan <- "stop"
		}()
	}
	return "", errors.New("failed to create channel after 3 attempts, please reset BOT_TOKEN")
}
