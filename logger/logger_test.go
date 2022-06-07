package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigBuildLink(t *testing.T) {
	ass := assert.New(t)

	logger := BuildLoggerLink().SetDateFormat(LogDateFormat.ShortLine).SetTimeFormat(LogTimeFormat.Default).Build()

	df := logger.cfg.dateFormat
	fmt.Println(df)
	fmt.Println(logger.cfg)
	ass.EqualValues(df, "yyyy-MM-dd", "日期格式错误")
}

func TestConfigBuild(t *testing.T) {
	ass := assert.New(t)

	cfg := &YiLogConfig{
		dateFormat: LogDateFormat.ShortLine,
		timeFormat: LogTimeFormat.Default,
	}

	logger := BuildLogger(cfg)

	df := logger.cfg.dateFormat
	fmt.Println(df)
	fmt.Println(logger.cfg)
	ass.EqualValues(df, "yyyy-MM-dd", "日期格式错误")
}

func TestBuildLogEntry(t *testing.T) {

	cfg := &YiLogConfig{
		dateFormat: LogDateFormat.Compact,
		timeFormat: LogTimeFormat.Compact,
	}

	entry := buildLogEntry(cfg, LogLevel.InfoLevel, "test")

	fmt.Println(entry.DateTime)
}

func TestInfo2Console(t *testing.T) {
	cfg := &YiLogConfig{
		dateFormat: LogDateFormat.Default,
		timeFormat: LogTimeFormat.Default,
	}

	logger := BuildLogger(cfg)
	logger.Info("logger\ntest")
}
