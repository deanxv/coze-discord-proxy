package common

var Version = "v4.4.5" // this hard coding will be replaced automatically when building, no need to manually change

const (
	RequestIdKey = "X-Request-Id"
	OutTime      = "out-time"
)

var CozeErrorMessages = []string{"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
	"Something wrong occurs, please retry. If the error persists, please contact the support team.",
	"There are too many users now. Please try again a bit later.",
	"I'm sorry, but I can't assist with that."}

var CozeDailyLimitErrorMessages = []string{"You have exceeded the daily limit for sending messages to the bot. Please try again later."}
