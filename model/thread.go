package model

type ThreadResp struct {
	Id   string `json:"id" swaggertype:"string" description:"线程ID"`
	Name string `json:"name" swaggertype:"string" description:"线程名称"`
}

type ThreadReq struct {
	ChannelId       string `json:"channelId" swaggertype:"string" description:"频道Id"`
	Name            string `json:"name" swaggertype:"string" description:"线程名称"`
	ArchiveDuration int    `json:"archiveDuration" swaggertype:"number" description:"线程存档时间[分钟]"`
}
