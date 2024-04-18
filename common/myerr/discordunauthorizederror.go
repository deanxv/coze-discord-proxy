package myerr

import "fmt"

// 自定义错误类型
type DiscordUnauthorizedError struct {
	Message string
	ErrCode int
}

// 实现 error 接口的 Error 方法
func (e *DiscordUnauthorizedError) Error() string {
	return fmt.Sprintf("errCode: %v, message: %v", e.ErrCode, e.Message)
}
