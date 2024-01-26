package model

type Chat struct {
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`
	Stream    bool   `json:"stream"`
}
