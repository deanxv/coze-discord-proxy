package common

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"html/template"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Println(err)
	}
}

func GetIp() (ip string) {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return ip
	}

	for _, a := range ips {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				if strings.HasPrefix(ip, "10") {
					return
				}
				if strings.HasPrefix(ip, "172") {
					return
				}
				if strings.HasPrefix(ip, "192.168") {
					return
				}
				ip = ""
			}
		}
	}
	return
}

var sizeKB = 1024
var sizeMB = sizeKB * 1024
var sizeGB = sizeMB * 1024

func Bytes2Size(num int64) string {
	numStr := ""
	unit := "B"
	if num/int64(sizeGB) > 1 {
		numStr = fmt.Sprintf("%.2f", float64(num)/float64(sizeGB))
		unit = "GB"
	} else if num/int64(sizeMB) > 1 {
		numStr = fmt.Sprintf("%d", int(float64(num)/float64(sizeMB)))
		unit = "MB"
	} else if num/int64(sizeKB) > 1 {
		numStr = fmt.Sprintf("%d", int(float64(num)/float64(sizeKB)))
		unit = "KB"
	} else {
		numStr = fmt.Sprintf("%d", num)
	}
	return numStr + " " + unit
}

func Seconds2Time(num int) (time string) {
	if num/31104000 > 0 {
		time += strconv.Itoa(num/31104000) + " 年 "
		num %= 31104000
	}
	if num/2592000 > 0 {
		time += strconv.Itoa(num/2592000) + " 个月 "
		num %= 2592000
	}
	if num/86400 > 0 {
		time += strconv.Itoa(num/86400) + " 天 "
		num %= 86400
	}
	if num/3600 > 0 {
		time += strconv.Itoa(num/3600) + " 小时 "
		num %= 3600
	}
	if num/60 > 0 {
		time += strconv.Itoa(num/60) + " 分钟 "
		num %= 60
	}
	time += strconv.Itoa(num) + " 秒"
	return
}

func Interface2String(inter interface{}) string {
	switch inter.(type) {
	case string:
		return inter.(string)
	case int:
		return fmt.Sprintf("%d", inter.(int))
	case float64:
		return fmt.Sprintf("%f", inter.(float64))
	}
	return "Not Implemented"
}

func UnescapeHTML(x string) interface{} {
	return template.HTML(x)
}

func IntMax(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func GetUUID() string {
	code := uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	return code
}

func Max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, length)
	for i := 0; i < length; i++ {
		key[i] = keyChars[rand.Intn(len(keyChars))]
	}
	return string(key)
}

const keyChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func GetTimeString() string {
	now := time.Now()
	return fmt.Sprintf("%s%d", now.Format("20060102150405"), now.UnixNano()%1e9)
}

func Obj2Bytes(obj interface{}) ([]byte, error) {
	// 创建一个jsonIter的Encoder
	configCompatibleWithStandardLibrary := jsoniter.ConfigCompatibleWithStandardLibrary
	// 将结构体转换为JSON文本并保持顺序
	bytes, err := configCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// IsImageBase64 判断给定的字符串是否可能是 Base64 编码
func IsImageBase64(s string) bool {
	// 检查字符串是否符合数据URL的格式
	//if !strings.HasPrefix(s, "data:image/") || !strings.Contains(s, ";base64,") {
	//	return false
	//}

	if !strings.Contains(s, ";base64,") {
		return false
	}

	// 获取";base64,"后的Base64编码部分
	dataParts := strings.Split(s, ";base64,")
	if len(dataParts) != 2 {
		return false
	}
	base64Data := dataParts[1]

	// 尝试Base64解码
	_, err := base64.StdEncoding.DecodeString(base64Data)
	return err == nil
}

// IsURL 判断给定的字符串是否可能是 URL
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "ftp://")
}

// Contains checks if a string is present in a slice of strings.
func SliceContains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// RandomElement 返回给定切片中的随机元素
func RandomElement[T any](slice []T) (T, error) {
	if len(slice) == 0 {
		var zero T
		return zero, fmt.Errorf("empty slice")
	}

	// 确保每次随机都不一样
	rand.Seed(time.Now().UnixNano())

	// 随机选择一个索引
	index := rand.Intn(len(slice))
	return slice[index], nil
}

func ReverseSegment(s string, segLen int) []string {
	var result []string
	runeSlice := []rune(s) // 将字符串转换为rune切片，以正确处理多字节字符

	// 从字符串末尾开始切片
	for i := len(runeSlice); i > 0; i -= segLen {
		// 检查是否到达或超过字符串开始
		if i-segLen < 0 {
			// 如果超过，直接从字符串开始到当前位置的所有字符都添加到结果切片中
			result = append([]string{string(runeSlice[0:i])}, result...)
		} else {
			// 否则，从i-segLen到当前位置的子切片添加到结果切片中
			result = append([]string{string(runeSlice[i-segLen : i])}, result...)
		}
	}
	return result
}

func FilterSlice(slice []string, filter string) []string {
	var result []string
	for _, value := range slice {
		if value != filter {
			result = append(result, value)
		}
	}
	return result
}

// isSameDay 检查两个时间是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
