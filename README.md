# Yi-Logger

## 简介

**Yi-Logger** 是 Golang 语言编写一款简单易用的日志库

## 使用方法

### 方法一

~~~golang
func TestBuildLoggerConfig(t *testing.T) {
    cfg := &YiLogConfig{
        Compress:   true, // 是否压缩
        OutputWay:  OutPut.File, // 输出方式
        File:       "../test.log", // 日志保存位置
        MaxSize:    20, // 日志文件大小上限
        DateFormat: LogDateFormat.Default, // 日志中的日期格式 "2006-01-02"
        TimeFormat: LogTimeFormat.Default, // 日志中的时间格式 "15:04:05"
    }
    logger := BuildLogger(cfg)
  
    logger.Info("info message")
    logger.Info("info: %s", "this is a info message")
  
    logger.Trace("trace message")
    logger.Error("error message")
    logger.Debug("debug message")
    logger.Warn("warn message")
}
~~~

### 方法二

~~~golang
func TestBuildLoggerLink(t *testing.T) {
    logger := BuildLoggerLink()
                .SetCompress(true)
                .SetOutput(OutPut.File)
                .SetFile("../test.log")
                .SetMaxSize(20)
                .SetDateFormat(LogDateFormat.Default)
                .SetTimeFormat(LogTimeFormat.Default)
                .Build()
  
    logger.Info("info message")
    logger.Info("info: %s", "this is a info message")
	
    logger.Trace("trace message")
    logger.Error("error message")
    logger.Debug("debug message")
    logger.Warn("warn message")
}
~~~

