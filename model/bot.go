package model

type BotConfig struct {
	ProxySecret     string   `json:"proxySecret"`
	CozeBotId       string   `json:"cozeBotId"`
	Model           []string `json:"model"`
	ChannelId       string   `json:"channelId"`
	MessageMaxToken string   `json:"messageMaxToken"`
}

// FilterUniqueBotChannel 给定BotConfig切片,筛选出具有不同CozeBotId+ChannelId组合的元素
func FilterUniqueBotChannel(configs []*BotConfig) []*BotConfig {
	seen := make(map[string]struct{}) // 使用map来跟踪已见的CozeBotId+ChannelId组合
	var uniqueConfigs []*BotConfig

	for _, config := range configs {
		combo := config.CozeBotId + "+" + config.ChannelId // 创建组合键
		if _, exists := seen[combo]; !exists {
			seen[combo] = struct{}{} // 标记组合键为已见
			uniqueConfigs = append(uniqueConfigs, config)
		}
	}

	return uniqueConfigs
}
