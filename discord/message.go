package discord

import (
	"bytes"
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
func SendMsgByAuthorization(c *gin.Context, content, channelId string) (string, error) {
	postUrl := "https://discord.com/api/v9/channels/%s/messages"
	content = strings.Replace(content, `\u0026`, "&", -1)
	// 构造请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"content": content,
	})
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(postUrl, channelId), bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// 设置请求头-部分请求头不传没问题，但目前仍有被discord检测异常的风险
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", UserAuthorization)
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
		fmt.Println("Error sending request:", err)
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
		common.LogError(c.Request.Context(), fmt.Sprintf("result:%s", bodyString))
		return "", fmt.Errorf("ID is not a string")
	}
	return id, nil
}
