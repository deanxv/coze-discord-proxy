package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake"
)

// snowflakeGenerator 单例
var (
	generator *SnowflakeGenerator
	once      sync.Once
)

// SnowflakeGenerator 是雪花ID生成器的封装
type SnowflakeGenerator struct {
	flake *sonyflake.Sonyflake
}

// NextID 生成一个新的雪花ID
func NextID() (string, error) {
	once.Do(initGenerator)
	id, err := generator.flake.NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

// initGenerator 初始化生成器，只调用一次
func initGenerator() {
	st := sonyflake.Settings{
		StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	flake := sonyflake.NewSonyflake(st)
	if flake == nil {
		FatalLog("sonyflake not created")
	}
	generator = &SnowflakeGenerator{
		flake: flake,
	}
}
