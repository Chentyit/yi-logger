package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigBuildLink(t *testing.T) {
	ass := assert.New(t)

	logger := BuildLoggerLink().SetDateFormat(LogDateFormat.ShortLine).SetTimeFormat(LogTimeFormat.Default).Build()

	df := logger.cfg.DateFormat
	fmt.Println(df)
	fmt.Println(logger.cfg)
	ass.EqualValues(df, "2006-01-02", "日期格式错误")
}

func TestConfigBuild(t *testing.T) {
	ass := assert.New(t)

	cfg := &YiLogConfig{
		DateFormat: LogDateFormat.ShortLine,
		TimeFormat: LogTimeFormat.Default,
	}

	logger := BuildLogger(cfg)

	df := logger.cfg.DateFormat
	fmt.Println(df)
	fmt.Println(logger.cfg)
	ass.EqualValues(df, "2006-01-02", "日期格式错误")
}

func TestBuildLogEntry(t *testing.T) {

	cfg := &YiLogConfig{
		DateFormat: LogDateFormat.Compact,
		TimeFormat: LogTimeFormat.Compact,
	}

	entry := buildLogEntry(cfg, LogLevel.InfoLevel, "test")

	fmt.Println(entry.DateTime)
}

func TestInfo2Console(t *testing.T) {
	cfg := &YiLogConfig{
		DateFormat: LogDateFormat.Default,
		TimeFormat: LogTimeFormat.Default,
	}

	logger := BuildLogger(cfg)
	logger.Info("logger\ntest")
}

func TestWriteBigLog(t *testing.T) {
	cfg := &YiLogConfig{
		Compress:   true,
		OutputWay:  OutPut.File,
		File:       "../test.log",
		MaxSize:    20,
		DateFormat: LogDateFormat.Compact,
		TimeFormat: LogTimeFormat.Compact,
	}
	logger := BuildLogger(cfg)
	for true {
		logger.Info("info message")
		logger.Trace("trace message")
		logger.Error("error message")
		logger.Debug("debug message")
		logger.Warn("warn message")
	}
}
