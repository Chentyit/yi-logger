package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"yi-go-logger/file_op"
)

// Level 日志等级
type Level byte

var LogLevel = struct {
	TraceLevel Level
	DebugLevel Level
	InfoLevel  Level
	WarnLevel  Level
	ErrorLevel Level
	PanicLevel Level
}{0, 1, 2, 3, 4, 5}

var logLevel = []string{
	0: "TRACE",
	1: "DEBUG",
	2: "INFO",
	3: "WARN",
	4: "ERROR",
	5: "PANIC",
}

// DateFormat 日期格式选项类型
type DateFormat string

// LogDateFormat 日志日期格式
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

// TimeFormat 时间格式选项类型
type TimeFormat string

// LogTimeFormat 日志时间格式
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

// OutPutWay 输出方式
type OutPutWay byte

// OutPut 日志输出方式
var OutPut = struct {
	File    OutPutWay
	Console OutPutWay
	Default OutPutWay
}{0, 1, 0}

// YiLogConfig
// @author Tianyi
// @description 日志基础配置
type YiLogConfig struct {
	showMs     bool       // 是否显示毫秒 (默认: false，不显示毫秒)
	logLevel   Level      // 日志等级 (默认: TraceLevel -> 0 打印所有类型日志)
	maxSize    int        // 每个日志最大容量 (默认: 10，单位: MB)
	maxBackups int        // 最多保存记录个数 (默认：5)
	maxAge     int        // 做多保存天数	(默认: 7)
	outputWay  OutPutWay  // 输出方式 (默认: 0 -> 输出到控制台)
	dateFormat DateFormat // 日期格式 (默认: yyyy-MM-dd)
	timeFormat TimeFormat // 时间格式 (默认: hh:HH:ss)
	file       string     // 日志保存文件 (默认: ./当前目录)
}

// yiLogEntry
// @author Tianyi
// @description 每行日志记录
type yiLogEntry struct {
	DateTime string `json:"time"`    // 日志记录时间
	Trace    string `json:"trace"`   // 文件路径
	Line     int    `json:"line"`    // 文件行数
	Level    string `json:"level"`   // 日志级别
	Message  string `json:"message"` // 日志信息
}

// yiLogger
// @author Tianyi
// @description 通过 yiLogger 进行操作（写，读，创建文件等）
type yiLogger struct {
	mu   *sync.Mutex     // 同步锁
	fo   *file_op.FileOp // 文件 IO
	date time.Time       // 日期，用于判断是否需要换文件
	cfg  *YiLogConfig    // logger config
}

// BuildLogger
// @author Tianyi
// @description 时间传参进行配置
func BuildLogger(cfg *YiLogConfig) *yiLogger {
	return buildLogger(cfg)
}

// BuildLoggerLink
// @author Tianyi
// @description 链式构建日志配置
func BuildLoggerLink() *YiLogConfig {
	return &YiLogConfig{}
}

// SetShowMs
// @author Tianyi
// @description 设置是否显示毫秒
func (cfg *YiLogConfig) SetShowMs(show bool) *YiLogConfig {
	cfg.showMs = show
	return cfg
}

// SetDateFormat
// @author Tianyi
// @description 设置日期格式
func (cfg *YiLogConfig) SetDateFormat(df DateFormat) *YiLogConfig {
	cfg.dateFormat = df
	return cfg
}

// SetTimeFormat
// @author Tianyi
// @description 设置时间格式
func (cfg *YiLogConfig) SetTimeFormat(tf TimeFormat) *YiLogConfig {
	cfg.timeFormat = tf
	return cfg
}

// SetLevel
// @author Tianyi
// @description 设置日志等级
func (cfg *YiLogConfig) SetLevel(level Level) *YiLogConfig {
	cfg.logLevel = level
	return cfg
}

// SetMaxSize
// @author Tianyi
// @description 设置最大容量
func (cfg *YiLogConfig) SetMaxSize(maxSize int) *YiLogConfig {
	cfg.maxSize = maxSize
	return cfg
}

// SetMaxBackups
// @author Tianyi
// @description 设置最大备份数量
func (cfg *YiLogConfig) SetMaxBackups(maxBackups int) *YiLogConfig {
	cfg.maxBackups = maxBackups
	return cfg
}

// SetMaxAge
// @author Tianyi
// @description 设置最大保存天数
func (cfg *YiLogConfig) SetMaxAge(maxAge int) *YiLogConfig {
	cfg.maxAge = maxAge
	return cfg
}

// SetFile
// @author Tianyi
// @description 设置保存日志文件
func (cfg *YiLogConfig) SetFile(file string) *YiLogConfig {
	cfg.file = file
	return cfg
}

// SetOutput
// @author Tianyi
// @description 设置输出方式
func (cfg *YiLogConfig) SetOutput(outputWay OutPutWay) *YiLogConfig {
	cfg.outputWay = outputWay
	return cfg
}

// Build
// @author Tianyi
// @description
func (cfg *YiLogConfig) Build() *yiLogger {
	return buildLogger(cfg)
}

// buildLogEntry
// @author Tianyi
// @description 构建每行日志记录
func buildLogEntry(cfg *YiLogConfig, level Level, msg string) *yiLogEntry {

	parser := fmt.Sprintf("%v %v", cfg.dateFormat, cfg.timeFormat)
	dateTime := time.Now().Format(parser)
	// 定位调用目标
	trace, line := getTraceAndLine()

	return &yiLogEntry{
		DateTime: dateTime,
		Trace:    trace,
		Line:     line,
		Level:    logLevel[level],
		Message:  msg,
	}
}

// buildLogger
// @author Tianyi
// @description 构建 Logger 对象
func buildLogger(cfg *YiLogConfig) *yiLogger {
	// 配置默认值
	if cfg.maxSize == 0 {
		cfg.maxSize = 10
	}

	if cfg.maxBackups == 0 {
		cfg.maxBackups = 5
	}

	if cfg.maxAge == 0 {
		cfg.maxAge = 7
	}

	if len(cfg.dateFormat) == 0 {
		cfg.dateFormat = "yyyy/MM/dd"
	}

	if len(cfg.timeFormat) == 0 {
		cfg.timeFormat = "hh:HH:ss"
	}

	if cfg.outputWay == OutPut.File && len(cfg.file) == 0 {
		cfg.file = "./"
	}

	if cfg.outputWay == OutPut.Default {
		return &yiLogger{
			fo:   nil,
			date: time.Now(),
			cfg:  cfg,
		}
	} else {
		fo := file_op.CreateFileOp(cfg.file, cfg.maxSize, cfg.maxAge, cfg.maxBackups)
		return &yiLogger{
			fo:   fo,
			date: time.Now(),
			cfg:  cfg,
		}
	}
}

func (logger *yiLogger) Trace(format string, a ...any) {
	if LogLevel.TraceLevel <= logger.cfg.logLevel {
		return
	}

	log := createLog(logger.cfg, LogLevel.TraceLevel, format, a...)

	logger.output(log)
}

func (logger *yiLogger) Debug(format string, a ...any) {
	if LogLevel.DebugLevel <= logger.cfg.logLevel {
		return
	}

	log := createLog(logger.cfg, LogLevel.DebugLevel, format, a...)

	logger.output(log)
}

func (logger *yiLogger) Info(format string, a ...any) {
	// 如果 Log 配置的等级大于当前等级，则输入当前等级日志
	if LogLevel.InfoLevel <= logger.cfg.logLevel {
		return
	}

	log := createLog(logger.cfg, LogLevel.InfoLevel, format, a...)

	logger.output(log)
}

func (logger *yiLogger) Warn(format string, a ...any) {
	if LogLevel.WarnLevel <= logger.cfg.logLevel {
		return
	}

	log := createLog(logger.cfg, LogLevel.WarnLevel, format, a...)

	logger.output(log)
}

func (logger *yiLogger) Error(format string, a ...any) {
	if LogLevel.ErrorLevel <= logger.cfg.logLevel {
		return
	}

	log := createLog(logger.cfg, LogLevel.ErrorLevel, format, a...)

	logger.output(log)
}

// Panic
// @author Tianyi
// @description 该日志级别会直接让整个程序退出，慎用
func (logger *yiLogger) Panic(format string, a ...any) {
	if LogLevel.PanicLevel <= logger.cfg.logLevel {
		return
	}

	log := createLog(logger.cfg, LogLevel.PanicLevel, format, a...)

	logger.output(log)

	os.Exit(1)
}

// createLog
// @author Tianyi
// @description 生成日志内容
func createLog(cfg *YiLogConfig, logLevel Level, format string, a ...any) []byte {
	// 格式化 msg
	msg := formatMsg(format, a...)
	// 构建日志每行信息
	entry := buildLogEntry(cfg, logLevel, msg)
	log, _ := json.Marshal(entry)
	return log
}

// getTraceAndLine
// @author Tianyi
// @description 获取调用栈信息
func getTraceAndLine() (string, int) {
	_, trace, line, ok := runtime.Caller(4)
	if !ok {
		trace = "???"
		line = 0
	}
	return trace, line
}

// formatMsg
// @author Tianyi
// @description 格式化日志信息，因为该日志框架使用的是 Json 保存，保证每行日志都是
//				一个 Json 字符串，如果存在 '\n' 和 '\r'，就会导致 Json 字符串换
//				行，无法统一数据格式，以后也不方便扩展日志框架功能
func formatMsg(format string, a ...any) string {
	msg := fmt.Sprintf(format, a...)
	msg = strings.Replace(msg, "\n", " ", -1)
	msg = strings.Replace(msg, "\r", " ", -1)
	return msg
}

// output
// @author Tianyi
// @description 根据配置输出到文件或者控制台
func (logger *yiLogger) output(log []byte) {
	if logger.cfg.outputWay == OutPut.Default || logger.cfg.outputWay == OutPut.Console {
		fmt.Println(string(log))
	} else {
		_ = logger.fo.Write(log)
	}
}
