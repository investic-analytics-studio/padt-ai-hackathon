package logger

import (
	"fmt"
	"os"

	"github.com/quantsmithapp/datastation-backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger *zap.Logger
}

func NewLogger() Logger {
	cfg := config.GetConfig().Application
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = TIME_KEY
	encoderCfg.LevelKey = LEVEL_KEY
	encoderCfg.CallerKey = CALLER_KEY
	encoderCfg.MessageKey = MESSAGE_KEY
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = EncodeLevel()

	priorityLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.Level(cfg.LogLevel) && !IsErrorLevel(l)
	})

	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return IsErrorLevel(l)
	})

	encoder := NewLoggerEncoder(encoderCfg)
	stdoutCore := zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(os.Stdout)), priorityLevel)
	stderrCore := zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(os.Stdout)), errorLevel)
	core := zapcore.NewTee(stdoutCore, stderrCore)

	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &logger{logger: l}
}

func (l logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l logger) Error(err error, fields ...zap.Field) {
	l.logger.Error(err.Error(), fields...)
}

func (l logger) Panic(err error, fields ...zap.Field) {
	l.logger.Panic(err.Error(), fields...)
}

func (l logger) Fatal(err error, fields ...zap.Field) {
	l.logger.Fatal(err.Error(), fields...)
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Infof(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

func (l logger) Panicf(format string, v ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, v...))
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l logger) BuildFields(args ...interface{}) []zap.Field {
	fields := []zap.Field{}
	isEven := len(args)%2 == 0
	if !isEven {
		return fields
	}

	for i := 0; i < len(args); i += 2 {
		key, _ := args[i].(string)
		var field zapcore.Field
		switch v := args[i+1].(type) {
		case string:
			field = zap.String(key, v)
		case int:
			field = zap.Int(key, v)
		case int16:
			field = zap.Int16(key, v)
		case int32:
			field = zap.Int32(key, v)
		case int64:
			field = zap.Int64(key, v)
		case uint:
			field = zap.Uint(key, v)
		case uint16:
			field = zap.Uint16(key, v)
		case uint32:
			field = zap.Uint32(key, v)
		case uint64:
			field = zap.Uint64(key, v)
		case float32:
			field = zap.Float32(key, v)
		case float64:
			field = zap.Float64(key, v)
		case bool:
			field = zap.Bool(key, v)
		default:
			field = zap.Any(key, v)
		}

		fields = append(fields, field)
	}

	return fields
}

func (l *logger) Close() {
	l.logger.Sync()
}
