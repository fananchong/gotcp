package gotcp

import "fmt"

// ILogger : 日志类接口
type ILogger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// DefaultLogger : 缺省日志类
type DefaultLogger struct {
}

// NewDefaultLogger : 缺省日志类的构造函数
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

// Info :
func (logger *DefaultLogger) Info(args ...interface{}) {
	fmt.Println(args...)
}

// Infof :
func (logger *DefaultLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Error :
func (logger *DefaultLogger) Error(args ...interface{}) {
	fmt.Println(args...)
}

// Errorf :
func (logger *DefaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

var (
	xlog ILogger = NewDefaultLogger()
)

// SetLogger : 设置日志类实例
func SetLogger(log ILogger) {
	xlog = log
}
