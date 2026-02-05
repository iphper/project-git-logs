package utils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// 日志格式化结构
type LogFormatter struct{}

// 日志格式化逻辑
func (f *LogFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(
		fmt.Sprintf("[%s][%s] %s\n",
			e.Level,
			e.Time.Format("2006-01-02 15:04:05"),
			strings.ReplaceAll(e.Message, "\n", ""),
		),
	), nil
}

func init() {
	// 设置自定义的格式化逻辑
	logrus.SetFormatter(&LogFormatter{})
}
