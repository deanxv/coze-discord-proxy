package model

type OpenAIChatCompletionRequest struct {
	Model    string              `json:"model"`
	Stream   bool                `json:"stream"`
	Messages []OpenAIChatMessage `json:"messages"`
	OpenAIChatCompletionExtraRequest
}

type OpenAIChatCompletionExtraRequest struct {
	ChannelId *string `json:"channelId"`
}

type OpenAIChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type OpenAIErrorResponse struct {
	OpenAIError OpenAIError `json:"error"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type OpenAIChatCompletionChan struct {
	Model    string `json:"model"`
	Response chan OpenAIChatCompletionResponse
}

type OpenAIChatCompletionResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []OpenAIChoice `json:"choices"`
	Usage             OpenAIUsage    `json:"usage"`
	SystemFingerprint *string        `json:"system_fingerprint"`
	Suggestions       []string       `json:"suggestions"`
}

type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	LogProbs     *string       `json:"logprobs"`
	FinishReason *string       `json:"finish_reason"`
	Delta        OpenAIDelta   `json:"delta"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIDelta struct {
	Content string `json:"content"`
}

type OpenAIImagesGenerationRequest struct {
	OpenAIChatCompletionExtraRequest
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	ResponseFormat string `json:"response_format"`
}

type OpenAIImagesGenerationResponse struct {
	Created     int64                                 `json:"created"`
	DailyLimit  bool                                  `json:"dailyLimit"`
	Data        []*OpenAIImagesGenerationDataResponse `json:"data"`
	Suggestions []string                              `json:"suggestions"`
}

type OpenAIImagesGenerationDataResponse struct {
	URL           string `json:"url"`
	RevisedPrompt string `json:"revised_prompt"`
	B64Json       string `json:"b64_json"`
}

type OpenAIGPT4VImagesReq struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	ImageURL struct {
		URL string `json:"url"`
	} `json:"image_url"`
}

// Model represents a model with its properties.
type OpenaiModelResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	//Created time.Time `json:"created"`
	//OwnedBy string    `json:"owned_by"`
}

// ModelList represents a list of models.
type OpenaiModelListResponse struct {
	Object string                `json:"object"`
	Data   []OpenaiModelResponse `json:"data"`
}

//type ChannelIdentifier interface {
//	GetChannelId() *string
//}
//
//func (request OpenAIChatCompletionRequest) GetChannelId() *string {
//	return request.ChannelId
//}
//
//func (request OpenAIImagesGenerationRequest) GetChannelId() *string {
//	return request.ChannelId
//}
