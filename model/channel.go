package model

type ChannelResp struct {
	Id   string `json:"id" swaggertype:"string" description:"频道ID"`
	Name string `json:"name" swaggertype:"string" description:"频道名称"`
}

type ChannelStopChan struct {
	Id    string `json:"id" `
	IsNew bool   `json:"IsNew"`
}

type ChannelReq struct {
	ParentId string `json:"parentId" swaggertype:"string" description:"父频道Id,为空时默认为创建父频道"`
	Type     int    `json:"type" swaggertype:"number" description:"类型:[0:文本频道,4:频道分类](其它枚举请查阅discord-api文档)"`
	Name     string `json:"name" swaggertype:"string" description:"频道名称"`
}
