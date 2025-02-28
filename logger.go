package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	env "github.com/2997215859/goenv"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func init() {
	InitDailyLogger()
}

var Logger *logrus.Logger

type formatter struct{}

func (m *formatter) Format(entry *logrus.Entry) ([]byte, error) {

	b := entry.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")

	//newLog := fmt.Sprintf("%s [%s] %s:%d %s. %s\n", timestamp, entry.Level, path.Base(entry.Caller.File), entry.Caller.Line, entry.Caller.Function, entry.Message)
	//newLog := fmt.Sprintf("%s [%s] %s:%d %s. %s\n", timestamp, entry.Level, path.Base(entry.Data["file"].(string)), entry.Data["line"], entry.Data["function"], entry.Message)
	newLog := fmt.Sprintf("%s [%s] %s:%d %s. %s\n", timestamp, entry.Level, path.Base(entry.Data["file"].(string)), entry.Data["line"], path.Base(entry.Data["function"].(string)), entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func createLogFile() *os.File {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Panicf("os getwd error: %s", err)
	}
	dir := currentDir + "/logs"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			logrus.Panicf("make dir(%s) error: %s", dir, err)
		}
	}

	date := time.Now().Format("20060102")
	logPath := fmt.Sprintf("%s/%s.log", dir, date)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Panicf("open file error: %s", err)
	}
	return file
}

const SkipKey = "@skip"

func InitDailyLogger() {
	InitLogger("./logs/app.access_log.%Y_%m_%d",
		rotatelogs.WithLinkName("./logs/app.access_log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
}

func InitHourLogger() {
	InitLogger("./logs/app.access_log.%Y_%m_%d_%H_%M_%S",
		rotatelogs.WithLinkName("./logs/app.access_log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
}

func IsStdout() bool {
	if os.Getenv("log_stdout") == "true" {
		return true
	}
	return false
}

func InitLogger(pattern string, options ...rotatelogs.Option) {
	Logger = logrus.New()

	rotateWriter, err := rotatelogs.New(pattern, options...)
	if err != nil {
		panic("rotatelogs.New error: %s")
	}

	writers := []io.Writer{rotateWriter}
	if env.ENV() == env.ENV_DEV || IsStdout() {
		writers = append(writers, os.Stdout)
	}
	mw := io.MultiWriter(writers...)

	Logger.SetOutput(mw)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&formatter{})
	Logger.AddHook(NewHook(WithSkipKey(SkipKey)))
}

func Error(format string, args ...interface{}) {
	Logger.WithField(SkipKey, 1).Errorf(format, args...)
}

func Info(format string, args ...interface{}) {
	Logger.WithField(SkipKey, 1).Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	Logger.WithField(SkipKey, 1).Warnf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	Logger.WithField(SkipKey, 1).Fatalf(format, args...)
}
