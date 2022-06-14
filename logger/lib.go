package logger

import (
	"fmt"
	"runtime"
	"strings"
)

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
