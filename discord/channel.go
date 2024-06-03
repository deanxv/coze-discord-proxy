package discord

import (
	"context"
	"coze-discord-proxy/common"
	"coze-discord-proxy/common/config"
	"coze-discord-proxy/telegram"
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
	if err == nil && !strings.HasPrefix(channel.Name, "cdp-chat-") {
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
		//common.SysError(fmt.Sprintf("频道无定时删除:%s", channelId))
	}
}

func ChannelCreate(guildID, channelName string, channelType int) (string, error) {
	// 创建新的频道
	st, err := Session.GuildChannelCreate(guildID, channelName, discordgo.ChannelType(channelType))
	if err != nil {
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

	for attempt := 1; attempt <= 3; attempt++ {
		resultChan := make(chan channelCreateResult, 1)

		go func() {
			id, err := ChannelCreate(guildID, channelName, channelType)
			resultChan <- channelCreateResult{
				ID:  id,
				Err: err,
			}
		}()

		// 设置超时时间为60秒
		select {
		case result := <-resultChan:
			if result.Err != nil {
				if strings.Contains(result.Err.Error(), "Maximum number of server channels reached") {

					var ok bool
					var err error

					if config.MaxChannelDelType == "ALL" {
						// 频道已满 删除所有频道
						ok, err = ChannelDelAllForCdp(c)
					} else if config.MaxChannelDelType == "OLDEST" {
						// 频道已满 删除最旧频道
						ok, err = ChannelDelOldestForCdp(c)
					} else {
						return "", result.Err
					}
					if err != nil {
						return "", err
					}
					if ok {
						return ChannelCreate(guildID, channelName, channelType)
					} else {
						return "", fmt.Errorf("当前Discord服务器频道已满，且无CDP临时频道可删除")
					}
				} else {
					return "", result.Err
				}
			}
			// 成功创建频道，返回结果
			return result.ID, nil
		case <-time.After(60 * time.Second):
			common.LogWarn(c, fmt.Sprintf("Create channel timed out, retrying...%v", attempt))
		}
	}
	// tg发送通知
	if !common.IsSameDay(CreateChannelRiskPreNotifyTime, time.Now()) && telegram.NotifyTelegramBotToken != "" && telegram.TgBot != nil {
		go func() {
			CreateChannelRiskChan <- "stop"
		}()
	}
	// 所有尝试后仍失败，返回最后的错误
	return "", fmt.Errorf("failed after 3 attempts due to timeout, please reset BOT_TOKEN")
}

func ChannelDelAllForCdp(c *gin.Context) (bool, error) {
	// 获取服务器内所有频道的信息
	channels, err := Session.GuildChannels(GuildId)
	if err != nil {
		common.LogError(c, fmt.Sprintf("服务器Id查询频道失败 %s", err.Error()))
		return false, err
	}

	flag := false

	// 遍历所有频道
	for _, channel := range channels {

		// 过滤掉配置中的频道id
		for _, config := range BotConfigList {
			if config.ChannelId == channel.ID {
				continue
			}
		}

		if ChannelId == channel.ID {
			continue
		}

		// 检查频道名是否以"cdp-"开头
		if strings.HasPrefix(channel.Name, "cdp-chat-") {
			// 删除该频道
			_, err := Session.ChannelDelete(channel.ID)
			if err != nil {
				common.LogError(c, fmt.Sprintf("频道数量已满-删除频道异常(可能原因:对话请求频道已被自动删除) %s", err.Error()))
				return false, err
			}
			common.LogWarn(c, fmt.Sprintf("频道数量已满-自动删除频道Id %s", channel.ID))
			flag = true
		}
	}
	return flag, nil
}

func ChannelDelOldestForCdp(c *gin.Context) (bool, error) {
	// 获取服务器内所有频道的信息
	channels, err := Session.GuildChannels(GuildId)
	if err != nil {
		common.LogError(c, fmt.Sprintf("服务器Id查询频道失败 %s", err.Error()))
		return false, err
	}

	// 遍历所有频道
	for _, channel := range channels {

		// 过滤掉配置中的频道id
		for _, config := range BotConfigList {
			if config.ChannelId == channel.ID {
				continue
			}
		}

		if ChannelId == channel.ID {
			continue
		}

		// 检查频道名是否以"cdp-"开头
		if strings.HasPrefix(channel.Name, "cdp-chat-") {
			// 删除该频道
			_, err := Session.ChannelDelete(channel.ID)
			if err != nil {
				common.LogError(c, fmt.Sprintf("频道数量已满-删除频道异常(可能原因:对话请求频道已被自动删除) %s", err.Error()))
				return false, err
			}
			common.LogWarn(c, fmt.Sprintf("频道数量已满-自动删除频道Id %s", channel.ID))
			return true, nil
		}
	}
	return false, nil
}
