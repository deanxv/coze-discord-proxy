package model

type ChatReq struct {
	ChannelId *string `json:"channelId"  swaggertype:"string" description:"频道ID/线程ID"`
	Content   string  `json:"content" swaggertype:"string" description:"消息内容"`
	Stream    bool    `json:"stream"  swaggertype:"boolean" description:"是否流式返回"`
}

func (request ChatReq) GetChannelId() *string {
	return request.ChannelId
}
