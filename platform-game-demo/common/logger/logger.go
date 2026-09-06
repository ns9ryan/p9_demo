package logger

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// 日志级别前缀（高亮）
const (
	DebugPrefix = "🔍 [DEBUG]"
	InfoPrefix  = "ℹ️  [INFO]"
	WarnPrefix  = "⚠️  [WARN]"
	ErrorPrefix = "❌ [ERROR]"
)

// Log 日志管理器（单例模式）
type Log struct {
	logger logx.Logger
}

var defaultLog = &Log{
	logger: logx.WithContext(nil),
}

// Init 初始化日志器（可选，根据需要配置）
func Init(logger logx.Logger) {
	if logger != nil {
		defaultLog.logger = logger
	}
}

// Debug 输出调试日志
func Debug(msg string, fields ...interface{}) {
	defaultLog.logger.Debugf("%s %s", DebugPrefix, formatMsg(msg, fields...))
}

// Debugf 输出格式化的调试日志
func Debugf(format string, args ...interface{}) {
	defaultLog.logger.Debugf("%s %s", DebugPrefix, fmt.Sprintf(format, args...))
}

// Info 输出信息日志
func Info(msg string, fields ...interface{}) {
	defaultLog.logger.Infof("%s %s", InfoPrefix, formatMsg(msg, fields...))
}

// Infof 输出格式化的信息日志
func Infof(format string, args ...interface{}) {
	defaultLog.logger.Infof("%s %s", InfoPrefix, fmt.Sprintf(format, args...))
}

// Warn 输出警告日志
func Warn(msg string, fields ...interface{}) {
	defaultLog.logger.Infof("%s %s", WarnPrefix, formatMsg(msg, fields...))
}

// Warnf 输出格式化的警告日志
func Warnf(format string, args ...interface{}) {
	defaultLog.logger.Infof("%s %s", WarnPrefix, fmt.Sprintf(format, args...))
}

// Error 输出错误日志
func Error(msg string, fields ...interface{}) {
	defaultLog.logger.Error(fmt.Sprintf("%s %s", ErrorPrefix, formatMsg(msg, fields...)))
}

// Errorf 输出格式化的错误日志
func Errorf(format string, args ...interface{}) {
	defaultLog.logger.Errorf("%s %s", ErrorPrefix, fmt.Sprintf(format, args...))
}

// WithFields 返回包含额外字段的日志器
func WithFields(fields ...interface{}) logx.Logger {
	return defaultLog.logger
}

// 内部方法：格式化消息和字段
func formatMsg(msg string, fields ...interface{}) string {
	if len(fields) == 0 {
		return msg
	}
	return fmt.Sprintf("%s %v", msg, fields)
}
