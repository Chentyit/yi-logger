package yi_logger

import (
	"github.com/Chentyit/yi-logger/logger"
	"testing"
)

func BenchmarkLog2Console(b *testing.B) {
	cfg := &logger.YiLogConfig{
		LogLevel:   logger.LogLevel.InfoLevel,
		Compress:   true,
		OutputWay:  logger.OutPut.Console,
		MaxSize:    50,
		DateFormat: logger.LogDateFormat.Compact,
		TimeFormat: logger.LogTimeFormat.Compact,
	}
	l := logger.BuildLogger(cfg)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("info message: %s", "This is a Benchmark Info.")
		}
	})
}

func BenchmarkLog2File(b *testing.B) {
	cfg := &logger.YiLogConfig{
		LogLevel:   logger.LogLevel.InfoLevel,
		Compress:   true,
		OutputWay:  logger.OutPut.File,
		File:       "./test.log",
		MaxSize:    50,
		DateFormat: logger.LogDateFormat.Compact,
		TimeFormat: logger.LogTimeFormat.Compact,
	}
	l := logger.BuildLogger(cfg)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("info message: %s", "This is a Benchmark Info.")
		}
	})
}
