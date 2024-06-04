package common

var Version = "v4.5.1" // this hard coding will be replaced automatically when building, no need to manually change

const (
	RequestIdKey = "X-Request-Id"
	OutTime      = "out-time"
)

var ImgGeneratePrompt = "Please adhere strictly to my instructions below for the drawing task. If I do not provide a specific directive for drawing, create an image that corresponds to the text I have provided:\\n"

var DefaultOpenaiModelList = []string{
	"gpt-3.5-turbo", "gpt-3.5-turbo-0301", "gpt-3.5-turbo-0613", "gpt-3.5-turbo-1106", "gpt-3.5-turbo-0125",
	"gpt-3.5-turbo-16k", "gpt-3.5-turbo-16k-0613",
	"gpt-3.5-turbo-instruct",
	"gpt-4", "gpt-4-0314", "gpt-4-0613", "gpt-4-1106-preview", "gpt-4-0125-preview",
	"gpt-4-32k", "gpt-4-32k-0314", "gpt-4-32k-0613",
	"gpt-4-turbo-preview", "gpt-4-turbo", "gpt-4-turbo-2024-04-09",
	"gpt-4o", "gpt-4o-2024-05-13",
	"gpt-4-vision-preview",
	"dall-e-3",
}

var CozeErrorMessages = []string{
	"Something wrong occurs, please retry. If the error persists, please contact the support team.",
	"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
	"Some error occurred. Please try again or contact the support team in our communities.",
	"We've detected unusual traffic from your network, so Coze is temporarily unavailable.",
	"There are too many users now. Please try again a bit later.",
	"I'm sorry, but I can't assist with that.",
}

var CozeDailyLimitErrorMessages = []string{
	"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
}
