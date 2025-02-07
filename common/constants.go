package common

var Version = "v4.6.5" // this hard coding will be replaced automatically when building, no need to manually change

const (
	RequestIdKey = "X-Request-Id"
	OutTime      = "out-time"
)

var ImgGeneratePrompt = "Please adhere strictly to my instructions below for the drawing task. If I do not provide a specific directive for drawing, create an image that corresponds to the text I have provided:\\n"

var DefaultOpenaiModelList = []string{
	"gpt-3.5-turbo",
	"gpt-3.5-turbo-1106",
	"gpt-3.5-turbo-0125",
	"gpt-4",
	"gpt-4-0613",
	"gpt-4-32k",
	"gpt-4-32k-0613",
	"gpt-4-turbo",
	"gpt-4-turbo-preview",
	"gpt-4o",
	"gpt-4o-2024-05-13",
	"gpt-4o-2024-08-06",
	"gpt-4o-2024-11-20",
	"chatgpt-4o-latest",
	"gpt-4o-mini",
	"gpt-4o-mini-2024-07-18",
	"gpt-4-vision-preview",
	"gpt-4-turbo-2024-04-09",
	"gpt-4-1106-preview",
	"dall-e-3",
	"o1-mini",
	"o1-preview",
	"o3-mini",
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
	"Hi there! You've used up your free chat credits. To continue enjoying our service, please consider upgrading to our premium plan [Upgrade to Coze Premium to chat]",
	"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
	"Hi there! You've used up your credits for today. To continue enjoying our service, please try again tomorrow or consider upgrading to our premium plan.",
}

var CozeCreatorDailyLimitErrorMessages = []string{
	"The bot's usage is covered by the developer, but due to the developer's message credits being exhausted, the bot is temporarily unavailable.",
}

var CozeDailyLimitErrorMessages = append(CozeUserDailyLimitErrorMessages, CozeCreatorDailyLimitErrorMessages...)
