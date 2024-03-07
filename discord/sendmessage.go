package discord

import (
	"bytes"
	"context"
	"coze-discord-proxy/common"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 用户端发送消息 注意 此为临时解决方案 后续会优化代码
func SendMsgByAuthorization(c *gin.Context, userAuth, content, channelId string) (string, error) {
	var ctx context.Context
	if c == nil {
		ctx = context.Background()
	} else {
		ctx = c.Request.Context()
	}

	postUrl := "https://discord.com/api/v9/channels/%s/messages"

	// 构造请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"content": content,
	})
	if err != nil {
		common.LogError(ctx, fmt.Sprintf("Error encoding request body:%s", err))
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(postUrl, channelId), bytes.NewBuffer(requestBody))
	if err != nil {
		common.LogError(ctx, fmt.Sprintf("Error creating request:%s", err))
		return "", err
	}

	// 设置请求头-部分请求头不传没问题，但目前仍有被discord检测异常的风险
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", userAuth)
	req.Header.Set("Origin", "https://discord.com")
	req.Header.Set("Referer", fmt.Sprintf("https://discord.com/channels/%s/%s", GuildId, channelId))
	if UserAgent != "" {
		req.Header.Set("User-Agent", UserAgent)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	}

	// 发起请求
	client := &http.Client{}
	if ProxyUrl != "" {
		proxyURL, _ := url.Parse(ProxyUrl)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client = &http.Client{
			Transport: transport,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		common.LogError(ctx, fmt.Sprintf("Error sending request:%s", err))
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 将响应体转换为字符串
	bodyString := string(bodyBytes)

	// 使用map来解码JSON
	var result map[string]interface{}

	// 解码JSON到map中
	err = json.Unmarshal([]byte(bodyString), &result)
	if err != nil {
		return "", err
	}

	// 类型断言来获取id的值
	id, ok := result["id"].(string)

	if !ok {
		// 401
		if errMessage, ok := result["message"].(string); ok {
			if strings.Contains(errMessage, "401: Unauthorized") ||
				strings.Contains(errMessage, "You need to verify your account in order to perform this action.") {
				common.LogWarn(ctx, fmt.Sprintf("USER_AUTHORIZATION:%s EXPIRED", userAuth))
				return "", &common.DiscordUnauthorizedError{
					ErrCode: 401,
					Message: "discord 鉴权未通过",
				}
			}
		}
		common.LogError(ctx, fmt.Sprintf("user_auth:%s result:%s", userAuth, bodyString))
		return "", fmt.Errorf("/api/v9/channels/%s/messages response err", channelId)
	} else {
		return id, nil
	}
}
