package model

type ReplyResp struct {
	Content   string   `json:"content" swaggertype:"string" description:"回复内容"`
	EmbedUrls []string `json:"embedUrls" swaggertype:"array,string" description:"嵌入网址"`
}
