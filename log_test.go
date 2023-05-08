package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	InitLogger(
		"./logs/app.access_log.%Y_%m_%d_%H_%M_%S",
		rotatelogs.WithLinkName("./logs/app.access_log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	Logger.Info("asdfasdfasdfasd")
}
