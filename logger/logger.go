package logger

import (
	"os"
	"sync"
	"time"
)

// Level 日志等级
type Level int

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
	Normal    DateFormat
	ShortLine DateFormat
	Compact   DateFormat
	Default   DateFormat
}{
	"yyyy/MM/dd",
	"yyyy-MM-dd",
	"yyyyMMdd",
	"yyyy/MM/dd",
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
	"hh:HH:ss",
	"hh/HH/ss",
	"hhHHss",
	"hh:HH:ss",
}

// OutPutWay 输出方式
type OutPutWay int

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
	dateFormat DateFormat // 日期格式 (默认: yyyy/MM/dd)
	timeFormat TimeFormat // 时间格式 (默认: hh:HH:ss)
	file       string     // 日志保存文件 (默认: ./当前目录)
	prefix     string     // 日志前缀 (默认: 空)
}

// yiLogEntry
// @author Tianyi
// @description 每行日志记录
type yiLogEntry struct {
	Time    time.Time `json:"time"`    // 日志记录时间
	File    string    `json:"file"`    // 文件路径
	Line    int       `json:"line"`    // 文件行数
	Level   string    `json:"level"`   // 日志级别
	Message string    `json:"message"` // 日志信息
}

// yiLogger
// @author Tianyi
// @description 通过 yiLogger 进行操作（写，读，创建文件等）
type yiLogger struct {
	mu   *sync.Mutex  // 同步锁
	size int          // 记录当期日志文件大小
	file *os.File     // 文件 IO
	date time.Time    // 日期，用于判断是否需要换文件
	cfg  *YiLogConfig // logger config
}

// BuildLoggerLink
// @author Tianyi
// @description 链式构建日志配置
func BuildLoggerLink() *YiLogConfig {
	return &YiLogConfig{}
}

// BuildLogger
// @author Tianyi
// @description 时间传参进行配置
func BuildLogger(cfg *YiLogConfig) *yiLogger {
	// 添加默认值
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

	var logger *yiLogger

	if cfg.outputWay == 0 {
		logger = &yiLogger{
			size: 0,
			file: nil,
			date: time.Now(),
			cfg:  cfg,
		}
	} else {
		f, err := os.OpenFile(cfg.file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic("file path not found: " + cfg.file)
		}
		logger = &yiLogger{
			size: 0,
			file: f,
			date: time.Now(),
			cfg:  cfg,
		}
	}

	return logger
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

// SetPrefix
// @author Tianyi
// @description 设置每行日志
func (cfg *YiLogConfig) SetPrefix(prefix string) *YiLogConfig {
	cfg.prefix = prefix
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

	var logger *yiLogger

	if cfg.outputWay == 0 {
		logger = &yiLogger{
			size: 0,
			file: nil,
			date: time.Now(),
			cfg:  cfg,
		}
	} else {
		f, err := os.OpenFile(cfg.file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic("file path not found: " + cfg.file)
		}
		logger = &yiLogger{
			size: 0,
			file: f,
			date: time.Now(),
			cfg:  cfg,
		}
	}

	return logger
}
