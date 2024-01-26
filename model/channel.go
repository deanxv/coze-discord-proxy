package model

type ChannelResp struct {
	Id   string `json:"id" swaggertype:"string" description:"频道ID"`
	Name string `json:"name" swaggertype:"string" description:"频道名称"`
}

type ChannelReq struct {
	Name string `json:"name" swaggertype:"string" description:"频道名称"`
}
