package myerr

import "fmt"

type ModelNotFoundError struct {
	Message string
	ErrCode int
}

// 实现 error 接口的 Error 方法
func (e *ModelNotFoundError) Error() string {
	return fmt.Sprintf("errCode: %v, message: %v", e.ErrCode, e.Message)
}
