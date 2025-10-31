package logger

import "go.uber.org/zap/zapcore"

func EncodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(LEVEL_DEBUG)
		case zapcore.InfoLevel:
			enc.AppendString(LEVEL_INFO)
		case zapcore.WarnLevel:
			enc.AppendString(LEVEL_WARN)
		case zapcore.ErrorLevel:
			enc.AppendString(LEVEL_ERROR)
		case zapcore.DPanicLevel:
			enc.AppendString(LEVEL_DPANIC)
		case zapcore.PanicLevel:
			enc.AppendString(LEVEL_PANIC)
		case zapcore.FatalLevel:
			enc.AppendString(LEVEL_FATAL)
		}
	}
}

func IsErrorLevel(lv zapcore.Level) bool {
	isError := lv == zapcore.ErrorLevel ||
		lv == zapcore.DPanicLevel ||
		lv == zapcore.PanicLevel ||
		lv == zapcore.FatalLevel
	return isError
}
