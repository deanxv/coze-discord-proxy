package common

var Version = "v4.4.8" // this hard coding will be replaced automatically when building, no need to manually change

const (
	RequestIdKey = "X-Request-Id"
	OutTime      = "out-time"
)

var ImgGeneratePrompt = "Please adhere strictly to my instructions below for the drawing task. If I do not provide a specific directive for drawing, create an image that corresponds to the text I have provided:\\n"

var CozeErrorMessages = []string{"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
	"Something wrong occurs, please retry. If the error persists, please contact the support team.",
	"There are too many users now. Please try again a bit later.",
	"I'm sorry, but I can't assist with that.",
	"We've detected unusual traffic from your network, so Coze is temporarily unavailable."}

var CozeDailyLimitErrorMessages = []string{"You have exceeded the daily limit for sending messages to the bot. Please try again later."}
