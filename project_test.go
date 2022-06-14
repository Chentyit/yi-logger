package yi_logger

import (
	"github.com/Chentyit/yi-logger/logger"
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	cfg := &logger.YiLogConfig{
		Compress:   true,
		OutputWay:  logger.OutPut.File,
		File:       "../test.log",
		MaxSize:    50,
		DateFormat: logger.LogDateFormat.Compact,
		TimeFormat: logger.LogTimeFormat.Compact,
	}
	l := logger.BuildLogger(cfg)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("info message: %s", "This is a Benchmark Info.")
			l.Trace("This is a Benchmark Trace.")
			l.Error("error message")
			l.Debug("Debug")
			l.Warn("warn message: %s", "This is a Benchmark Warn.")
		}
	})
}
