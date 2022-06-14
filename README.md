# Yi-Logger

## Introduction

**Yi-logger** is an easy-to-use log library written in **Golang**.

The output format of this log library is JSON:

~~~json
{"time":"20220614 161129","trace":"/Users/chentianyi/Program/Goland-workplace/yi-logger/logger/logger_test.go","line":71,"level":"INFO","message":"info message"}
~~~

The format is JSON to ensure consistency when parsing logs.

## Log DateFormat

- **logger.LogDateFormat.Slash** - 2006/01/02
- **logger.LogDateFormat.Default** - 2006/01/02
- **logger.LogDateFormat.ShortLine** - 2006-01-02
- **logger.LogDateFormat.Compact** - 20060102

~~~golang
var LogDateFormat = struct {
    Slash     DateFormat
    ShortLine DateFormat
    Compact   DateFormat
    Default   DateFormat
}{
    "2006/01/02",
    "2006-01-02",
    "20060102",
    "2006-01-02",
}
~~~

## Log TimeFormat

- **logger.LogTimeFormat.Normal** - 15:04:05
- **logger.LogTimeFormat.Default** - 15:04:05
- **logger.LogTimeFormat.Slash** - 15/04/05
- **logger.LogTimeFormat.Compact** - 150405

~~~golang
var LogTimeFormat = struct {
    Normal  TimeFormat
    Slash   TimeFormat
    Compact TimeFormat
    Default TimeFormat
}{
    "15:04:05",
    "15/04/05",
    "150405",
    "15:04:05",
}
~~~

## Log Output Way

- logger.OutPut.Console
- logger.OutPut.File
- logger.OutPut.Default

~~~golang
var OutPut = struct {
    Console OutPutWay
    File    OutPutWay
    Default OutPutWay
}{0, 1, 0}
~~~

## Log Level

- TRACE
- DEBUG
- INFO
- WARN
- ERROR
- PANIC: This level will cause the program to **terminate and panic**

~~~golang
var LogLevel = struct {
	TraceLevel Level
	DebugLevel Level
	InfoLevel  Level
	WarnLevel  Level
	ErrorLevel Level
	PanicLevel Level
}{0, 1, 2, 3, 4, 5}
~~~

## Usage

### Method 1

~~~golang
import (
    "github.com/Chentyit/yi-logger/logger"
)

func TestBuildLoggerConfig(t *testing.T) {
    cfg := &logger.YiLogConfig{
        Compress:   true, 												// 是否压缩
        OutputWay:  logger.OutPut.File, 					// 输出方式
        File:       "../test.log", 								// 日志保存位置
        MaxSize:    20, 													// 日志文件大小上限
        DateFormat: logger.LogDateFormat.Default, // 日志中的日期格式 "2006-01-02"
        TimeFormat: logger.LogTimeFormat.Default, // 日志中的时间格式 "15:04:05"
    }
    logger := logger.BuildLogger(cfg)
  
    logger.Info("info message")
    logger.Info("info: %s", "this is a info message")
  
    logger.Trace("trace message")
    logger.Error("error message")
    logger.Debug("debug message")
    logger.Warn("warn message")
}
~~~

### Method 2

~~~golang
func TestBuildLoggerLink(t *testing.T) {
    logger := logger.BuildLoggerLink()
	            .SetCompress(true)
	            .SetOutput(logger.OutPut.File)
	            .SetFile("./test.log")
	            .SetMaxSize(20)
	            .SetDateFormat(logger.LogDateFormat.Default)
	            .SetTimeFormat(logger.LogTimeFormat.Default)
	            .Build()

    logger.Info("info message")
    logger.Info("info: %s", "this is a info message")
	
    logger.Trace("trace message")
    logger.Error("error message")
    logger.Debug("debug message")
    logger.Warn("warn message")
}
~~~

## Benchamark Test

### Output to console

~~~bash
❯ go test -bench=Console -run=none -count=3 -cpu=1,2,4,8 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/Chentyit/yi-logger
BenchmarkLog2Console              984158              1157 ns/op             600 B/op          9 allocs/op
BenchmarkLog2Console             1000000              1161 ns/op             600 B/op          9 allocs/op
BenchmarkLog2Console             1000000              1156 ns/op             600 B/op          9 allocs/op
BenchmarkLog2Console-2           1843352               638.2 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-2           1838167               642.5 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-2           1891980               637.2 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-4           3234051               369.3 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-4           3176193               371.7 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-4           3222770               371.2 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-8           3718717               321.2 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-8           3745462               319.3 ns/op           600 B/op          9 allocs/op
BenchmarkLog2Console-8           3733071               322.1 ns/op           600 B/op          9 allocs/op
PASS
ok      github.com/Chentyit/yi-logger   19.399s
~~~

### Output to file

~~~bash
❯ go test -bench=File -run=none -count=3 -cpu=1,2,4,8 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/Chentyit/yi-logger
BenchmarkLog2File         339774              3895 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File         332946              3768 ns/op             829 B/op         11 allocs/op
BenchmarkLog2File         351229              3737 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File-2       418248              2872 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File-2       405909              2959 ns/op             828 B/op         11 allocs/op
BenchmarkLog2File-2       411133              2849 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File-4       419563              2710 ns/op             828 B/op         11 allocs/op
BenchmarkLog2File-4       436137              2693 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File-4       431779              2712 ns/op             828 B/op         11 allocs/op
BenchmarkLog2File-8       428808              2569 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File-8       409362              2589 ns/op             826 B/op         11 allocs/op
BenchmarkLog2File-8       424374              2603 ns/op             826 B/op         11 allocs/op
PASS
ok      github.com/Chentyit/yi-logger   16.653s
~~~

