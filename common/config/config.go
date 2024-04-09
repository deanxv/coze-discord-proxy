package config

import (
	"coze-discord-proxy/common/env"
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
var MaxChannelDelType = os.Getenv("MAX_CHANNEL_DEL_TYPE")

var DebugEnabled = os.Getenv("DEBUG") == "true"

var RateLimitKeyExpirationDuration = 20 * time.Minute

var RequestOutTimeDuration = 5 * time.Minute

var (
	RequestRateLimitNum            = env.Int("REQUEST_RATE_LIMIT", 60)
	RequestRateLimitDuration int64 = 1 * 60
)
