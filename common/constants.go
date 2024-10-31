package common

var Version = "v4.6.3" // this hard coding will be replaced automatically when building, no need to manually change

const (
	RequestIdKey = "X-Request-Id"
	OutTime      = "out-time"
)

var ImgGeneratePrompt = "Please adhere strictly to my instructions below for the drawing task. If I do not provide a specific directive for drawing, create an image that corresponds to the text I have provided:\\n"

var DefaultOpenaiModelList = []string{
	"gpt-3.5-turbo",
	"gpt-3.5-turbo-1106",
	"gpt-3.5-turbo-0125",
	"gpt-3.5-turbo-16k",
	"gpt-3.5-turbo-16k-0613",
	"gpt-3.5-turbo-instruct",
	"gpt-4",
	"gpt-4-0613",
	"gpt-4-1106-preview",
	"gpt-4-0125-preview",
	"gpt-4-turbo-preview",
	"gpt-4-turbo",
	"gpt-4-turbo-2024-04-09",
	"gpt-4o",
	"gpt-4o-2024-05-13",
	"gpt-4o-mini",
	"gpt-4o-mini-2024-07-18",
	"text-embedding-ada-002",
	"text-embedding-3-small",
	"text-embedding-3-large",
	"text-moderation-latest",
	"text-moderation-stable",
	"davinci-002",
	"babbage-002",
	"dall-e-2",
	"dall-e-3",
	"whisper-1",
	"tts-1",
	"tts-1-1106",
	"tts-1-hd",
	"tts-1-hd-1106",
	"gpt-4o-2024-08-06",
	"o1-mini",
	"o1-preview",
	"o1-preview-2024-09-12",
	"chatgpt-4o-latest",
	"o1-mini-2024-09-12",
}

var CozeErrorMessages = append(append(CozeOtherErrorMessages, CozeUserDailyLimitErrorMessages...), CozeCreatorDailyLimitErrorMessages...)

var CozeOtherErrorMessages = []string{
	"Something wrong occurs, please retry. If the error persists, please contact the support team.",
	"Some error occurred. Please try again or contact the support team in our communities.",
	"We've detected unusual traffic from your network, so Coze is temporarily unavailable.",
	"There are too many users now. Please try again a bit later.",
	"I'm sorry, but I can't assist with that.",
}

var CozeUserDailyLimitErrorMessages = []string{
	"Hi there! You've used up your free chat credits. To continue enjoying our service, please consider upgrading to our premium plan [Upgrade to Coze Premium to chat](https://www.coze.com/premium?connectID=10000028&botID=7376964308913422354)",
	"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
}

var CozeCreatorDailyLimitErrorMessages = []string{
	"The bot's usage is covered by the developer, but due to the developer's message credits being exhausted, the bot is temporarily unavailable.",
}

var CozeDailyLimitErrorMessages = append(CozeUserDailyLimitErrorMessages, CozeCreatorDailyLimitErrorMessages...)
