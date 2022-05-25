package logger

import (
	"os"
	"sync"
	"time"
)

const (
	TraceLevel = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	PanicLevel
)

// 日志日期格式
const (
	DateFormatNormal    = iota // yyyy/MM/dd
	DateFormatShortLine        // yyyy-MM-dd
	DateFormatCompact          // yyyyMMdd
)

const (
	TimeFormatNormal  = iota // hh:HH:ss
	TimeFormatSlash          // hh/HH/ss
	TimeFormatCompact        // hhHHss
)

// 日志输出方式
const (
	File    = 1 << iota // 输出到文件
	Console             // 输出到控制台
	Default = Console   // 默认输出到控制台
)

type yiLogConfig struct {
	DirPath    string // 日志保存目录
	Prefix     string // 日志前缀
	Level      int    // 日志等级
	Maxsize    int    // 每个日志最大容量
	MaxBackups int    // 最多保存记录个数
	MaxAge     int    // 做多保存天数
	DateFormat int    // 日期格式
	TimeFormat int    // 时间格式
	ShowMs     bool   // 是否显示毫秒
}

type yiLogEntry struct {
	Time    time.Time `json:"time"`    // 日志记录时间
	File    string    `json:"file"`    // 文件路径
	Line    int       `json:"line"`    // 文件行数
	Level   string    `json:"level"`   // 日志级别
	Message string    `json:"message"` // 日志信息
}

type logger struct {
	mu   *sync.Mutex
	file *os.File
}
