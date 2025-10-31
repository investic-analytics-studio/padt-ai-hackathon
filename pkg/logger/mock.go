package logger

import "go.uber.org/zap"

type loggerMock struct{}

func NewLoggerMock() *loggerMock {
	return &loggerMock{}
}

func (l loggerMock) Debug(msg string, fields ...zap.Field) {
}

func (l loggerMock) Info(msg string, fields ...zap.Field) {
}

func (l loggerMock) Warn(msg string, fields ...zap.Field) {
}

func (l loggerMock) Error(err error, fields ...zap.Field) {
}

func (l loggerMock) Panic(err error, fields ...zap.Field) {
}

func (l loggerMock) Fatal(err error, fields ...zap.Field) {
}

func (l loggerMock) Debugf(format string, v ...interface{}) {
}

func (l loggerMock) Infof(format string, v ...interface{}) {
}

func (l loggerMock) Warnf(format string, v ...interface{}) {
}

func (l loggerMock) Errorf(format string, v ...interface{}) {
}

func (l loggerMock) Panicf(format string, v ...interface{}) {
}

func (l loggerMock) Fatalf(format string, v ...interface{}) {
}

func (l loggerMock) BuildFields(args ...interface{}) []zap.Field {
	return []zap.Field{}
}

func (l *loggerMock) Close() {
}
