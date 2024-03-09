package common

import (
	"os"
	"strings"
	"time"
)

var ProxySecret = os.Getenv("PROXY_SECRET")
var ProxySecrets = strings.Split(os.Getenv("PROXY_SECRET"), ",")
var AllDialogRecordEnable = os.Getenv("ALL_DIALOG_RECORD_ENABLE")
var RequestOutTime = os.Getenv("REQUEST_OUT_TIME")
var StreamRequestOutTime = os.Getenv("STREAM_REQUEST_OUT_TIME")
var SwaggerEnable = os.Getenv("SWAGGER_ENABLE")
var OnlyOpenaiApi = os.Getenv("ONLY_OPENAI_API")

var DebugEnabled = os.Getenv("DEBUG") == "true"

var Version = "v4.3.9" // this hard coding will be replaced automatically when building, no need to manually change

const (
	RequestIdKey = "X-Request-Id"
	OutTime      = "out-time"
)

// Shouldn't larger then RateLimitKeyExpirationDuration
var (
	GlobalApiRateLimitNum            = 60
	GlobalApiRateLimitDuration int64 = 3 * 60

	GlobalWebRateLimitNum            = 60
	GlobalWebRateLimitDuration int64 = 3 * 60

	UploadRateLimitNum            = 10
	UploadRateLimitDuration int64 = 60

	DownloadRateLimitNum            = 10
	DownloadRateLimitDuration int64 = 60

	CriticalRateLimitNum            = 20
	CriticalRateLimitDuration int64 = 20 * 60
)

var RateLimitKeyExpirationDuration = 20 * time.Minute

var RequestOutTimeDuration = 5 * time.Minute

var CozeErrorMessages = []string{"You have exceeded the daily limit for sending messages to the bot. Please try again later.",
	"Something wrong occurs, please retry. If the error persists, please contact the support team.",
	"There are too many users now. Please try again a bit later."}

var CozeDailyLimitErrorMessages = []string{"You have exceeded the daily limit for sending messages to the bot. Please try again later."}
