package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(string, ...zap.Field)
	Info(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Error(error, ...zap.Field)
	Panic(error, ...zap.Field)
	Fatal(error, ...zap.Field)
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
	BuildFields(...interface{}) []zap.Field
	Close()
}
