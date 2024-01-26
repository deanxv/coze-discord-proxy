package model

type ChatReq struct {
	ChannelID string `json:"channelID"  swaggertype:"string" description:"频道ID"`
	Content   string `json:"content" swaggertype:"string" description:"消息内容"`
	Stream    bool   `json:"stream"  swaggertype:"boolean" description:"是否流式返回"`
}
