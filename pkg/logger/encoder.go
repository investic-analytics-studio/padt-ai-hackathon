package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type loggerEncoder struct {
	encoder zapcore.Encoder
}

func NewLoggerEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	encoder := zapcore.NewJSONEncoder(config)
	return &loggerEncoder{encoder: encoder}
}

func (c *loggerEncoder) Clone() zapcore.Encoder {
	return c.encoder.Clone()
}

func (c *loggerEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// |=> manipulate logging message
	fullPath := entry.Caller.File
	line := entry.Caller.Line
	message := entry.Message
	if message == "" {
		message = "-"
	}

	entry.Message = fmt.Sprintf("[%v:%v]: %v", trimPath(fullPath), line, message)
	return c.encoder.EncodeEntry(entry, fields)
}

// Logging-specific marshalers.
func (c *loggerEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	return c.encoder.AddArray(key, marshaler)
}

func (c *loggerEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	return c.encoder.AddObject(key, marshaler)
}

// Built-in types.
func (c *loggerEncoder) AddBinary(key string, value []byte) { // for arbitrary bytes
	c.encoder.AddBinary(key, value)
}

func (c *loggerEncoder) AddByteString(key string, value []byte) { // for UTF-8 encoded bytes
	c.encoder.AddByteString(key, value)
}

func (c *loggerEncoder) AddBool(key string, value bool) { c.encoder.AddBool(key, value) }
func (c *loggerEncoder) AddComplex128(key string, value complex128) {
	c.encoder.AddComplex128(key, value)
}
func (c *loggerEncoder) AddComplex64(key string, value complex64) { c.encoder.AddComplex64(key, value) }
func (c *loggerEncoder) AddDuration(key string, value time.Duration) {
	c.encoder.AddDuration(key, value)
}
func (c *loggerEncoder) AddFloat64(key string, value float64) { c.encoder.AddFloat64(key, value) }
func (c *loggerEncoder) AddFloat32(key string, value float32) { c.encoder.AddFloat32(key, value) }
func (c *loggerEncoder) AddInt(key string, value int)         { c.encoder.AddInt(key, value) }
func (c *loggerEncoder) AddInt64(key string, value int64)     { c.encoder.AddInt64(key, value) }
func (c *loggerEncoder) AddInt32(key string, value int32)     { c.encoder.AddInt32(key, value) }
func (c *loggerEncoder) AddInt16(key string, value int16)     { c.encoder.AddInt16(key, value) }
func (c *loggerEncoder) AddInt8(key string, value int8)       { c.encoder.AddInt8(key, value) }
func (c *loggerEncoder) AddString(key, value string)          { c.encoder.AddString(key, value) }
func (c *loggerEncoder) AddTime(key string, value time.Time)  { c.encoder.AddTime(key, value) }
func (c *loggerEncoder) AddUint(key string, value uint)       { c.encoder.AddUint(key, value) }
func (c *loggerEncoder) AddUint64(key string, value uint64)   { c.encoder.AddUint64(key, value) }
func (c *loggerEncoder) AddUint32(key string, value uint32)   { c.encoder.AddUint32(key, value) }
func (c *loggerEncoder) AddUint16(key string, value uint16)   { c.encoder.AddUint16(key, value) }
func (c *loggerEncoder) AddUint8(key string, value uint8)     { c.encoder.AddUint8(key, value) }
func (c *loggerEncoder) AddUintptr(key string, value uintptr) { c.encoder.AddUintptr(key, value) }

// AddReflected uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (c *loggerEncoder) AddReflected(key string, value interface{}) error {
	return c.encoder.AddReflected(key, value)
}

// OpenNamespace opens an isolated namespace where all subsequent fields will
// be added. Applications can use namespaces to prevent key collisions when
// injecting loggers into sub-components or third-party libraries.
func (c *loggerEncoder) OpenNamespace(key string) { c.encoder.OpenNamespace(key) }

func trimPath(path string) string {
	parts := strings.Split(path, "/")
	size := len(parts)

	if size < 2 {
		return path
	}

	return "/" + strings.Join(parts[size-2:], "/")
}
