package model

type ChannelResp struct {
	Id   string `json:"id" swaggertype:"string" description:"频道ID"`
	Name string `json:"name" swaggertype:"string" description:"频道名称"`
}

type ChannelReq struct {
	ParentId string `json:"parentId" swaggertype:"string" description:"父频道Id,为空时默认为创建父频道"`
	Name     string `json:"name" swaggertype:"string" description:"频道名称"`
}
