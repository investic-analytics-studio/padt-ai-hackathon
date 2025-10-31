package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var loggerObj Logger
var mu sync.Mutex

func InitGlobalLogger() {
	if loggerObj == nil {
		loggerObj = NewLogger()
	}
}

func Debug(msg string, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Warn(msg, fields...)
}

func Error(err error, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Error(err, fields...)
}

func Panic(err error, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Panic(err, fields...)
}

func Fatal(err error, fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Fatal(err, fields...)
}

func Debugf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Debug(fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Info(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Warn(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Errorf(format, v)
}

func Panicf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Errorf(format, v)
}

func Fatalf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	loggerObj.Errorf(format, v)
}
