/*  logger.go
** @Date:               November 21, 2019
** @Last Modified time: 21/11/19 07:30
 */

package mimir

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	tag "github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap/zapcore"
)

type Logging interface {
	With(fields ...interface{}) Logging
	Debug(msg string, fields ...interface{})
	Debugf(msg string, args ...interface{})
	Info(msg string, fields ...interface{})
	Infof(msg string, args ...interface{})
	Warn(msg string, fields ...interface{})
	Warnf(msg string, args ...interface{})
	Error(msg string)
	Errorf(msg string, args ...interface{})
	Fatal(msg string, fields ...interface{})
	Panic(msg string, fields ...interface{})
	Field(key string, value interface{}) interface{} // for type fields
}

var (
	logger *zapLog
)

func Instance() *zapLog {
	return NewZap(ProductionCore())
}

func With(fields ...interface{}) Logging {
	return Instance().With(fields...)
}

func Field(key string, value interface{}) interface{} {
	return Instance().Field(key, value)
}

func Debug(msg string, fields ...interface{}) {
	Instance().Debug(msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	Instance().Debugf(msg, args...)
}

func Info(msg string, fields ...interface{}) {
	Instance().Info(msg, fields...)
}

func Infof(msg string, args ...interface{}) {
	Instance().Infof(msg, args...)
}

func Warn(msg string, fields ...interface{}) {
	Instance().Warn(msg, fields...)
}

func Warnf(msg string, args ...interface{}) {
	Instance().Warnf(msg, args...)
}

func Error(msg string) {
	Instance().Error(msg)
}

func Errorf(msg string, args ...interface{}) {
	Instance().Errorf(msg, args...)
}

func Fatal(msg string, fields ...interface{}) {
	Instance().Fatal(msg, fields...)
}

func Panic(msg string, fields ...interface{}) {
	Instance().Panic(msg, fields...)
}

type spanLogger struct {
	logger Logging
	span   opentracing.Span
}

func (sl spanLogger) Panic(msg string, fields ...interface{}) {
	panic("implement me")
}

func (sl spanLogger) Field(key string, value interface{}) interface{} {
	return sl.logger.Field(key, value)
}

func (sl spanLogger) Debug(msg string, fields ...interface{}) {
	sl.logToSpan("debug", msg, fields...)
	sl.logger.Debug(msg, fields...)
}

func (sl spanLogger) Debugf(msg string, args ...interface{}) {
	sl.logToSpan("debug", fmt.Sprintf(msg, args...))
	sl.logger.Debugf(msg, args...)
}

func (sl spanLogger) Infof(msg string, args ...interface{}) {
	sl.logToSpan("info", fmt.Sprintf(msg, args...))
	sl.logger.Debugf(msg, args...)
}

func (sl spanLogger) Warn(msg string, fields ...interface{}) {
	sl.logToSpan("warn", msg, fields...)
	sl.logger.Warn(msg, fields...)
}

func (sl spanLogger) Warnf(msg string, args ...interface{}) {
	sl.logToSpan("warn", fmt.Sprintf(msg, args...))
	sl.logger.Warnf(msg, args...)
}

func (sl spanLogger) Errorf(msg string, args ...interface{}) {
	sl.logToSpan("error", fmt.Sprintf(msg, args...))
	tag.Error.Set(sl.span, true)
	sl.logger.Errorf(msg, args...)
}

func (sl spanLogger) Info(msg string, fields ...interface{}) {
	sl.logToSpan("info", msg, fields...)
	sl.logger.Info(msg, fields...)
}

func (sl spanLogger) Error(msg string) {
	sl.logToSpan("error", msg)
	tag.Error.Set(sl.span, true)
	sl.logger.Error(msg)
}

func (sl spanLogger) Fatal(msg string, fields ...interface{}) {
	sl.logToSpan("fatal", msg, fields...)
	tag.Error.Set(sl.span, true)
	sl.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (sl spanLogger) With(fields ...interface{}) Logging {
	sl.logger.With(fields...)
	sl.logToSpan("With", "logging field for", fields...)
	return sl
}

func (sl spanLogger) logToSpan(level string, msg string, fields ...interface{}) {
	fa := fieldAdapter(make([]log.Field, 0, 2+len(fields)))
	fa = append(fa, log.String("event", msg))
	fa = append(fa, log.String("level", level))
	for _, field := range fields {
		field.(zapcore.Field).AddTo(&fa)
	}
	sl.span.LogFields(fa...)
}

type LogFactory struct {
	logger Logging
}

func For(ctx context.Context) Logging {
	b := LogFactory{logger: With()}
	if span := opentracing.SpanFromContext(ctx); span != nil {
		var tracerId string
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			tracerId = sc.TraceID().String()
		}

		return spanLogger{span: span, logger: b.logger.With(b.logger.Field("trace_id", tracerId))}
	}
	return b.Bg()
}

func (b LogFactory) Bg() Logging {
	return With(Field("", ""))
}
func (b LogFactory) With(fields ...interface{}) LogFactory {
	return LogFactory{logger: b.logger.With(fields...)}
}

type fieldAdapter []log.Field

func (fa *fieldAdapter) AddBool(key string, value bool) {
	*fa = append(*fa, log.Bool(key, value))
}

func (fa *fieldAdapter) AddFloat64(key string, value float64) {
	*fa = append(*fa, log.Float64(key, value))
}

func (fa *fieldAdapter) AddFloat32(key string, value float32) {
	*fa = append(*fa, log.Float64(key, float64(value)))
}

func (fa *fieldAdapter) AddInt(key string, value int) {
	*fa = append(*fa, log.Int(key, value))
}

func (fa *fieldAdapter) AddInt64(key string, value int64) {
	*fa = append(*fa, log.Int64(key, value))
}

func (fa *fieldAdapter) AddInt32(key string, value int32) {
	*fa = append(*fa, log.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt16(key string, value int16) {
	*fa = append(*fa, log.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt8(key string, value int8) {
	*fa = append(*fa, log.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint(key string, value uint) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint64(key string, value uint64) {
	*fa = append(*fa, log.Uint64(key, value))
}

func (fa *fieldAdapter) AddUint32(key string, value uint32) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint16(key string, value uint16) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint8(key string, value uint8) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUintptr(key string, value uintptr)                        {}
func (fa *fieldAdapter) AddArray(key string, marshaler zapcore.ArrayMarshaler) error { return nil }
func (fa *fieldAdapter) AddComplex128(key string, value complex128)                  {}
func (fa *fieldAdapter) AddComplex64(key string, value complex64)                    {}
func (fa *fieldAdapter) AddObject(key string, value zapcore.ObjectMarshaler) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	*fa = append(*fa, log.Object(key, b))
	return nil
}
func (fa *fieldAdapter) AddReflected(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	*fa = append(*fa, log.Object(key, string(b)))
	return nil
}
func (fa *fieldAdapter) OpenNamespace(key string) {}

func (fa *fieldAdapter) AddDuration(key string, value time.Duration) {
	// TODO inefficient
	*fa = append(*fa, log.String(key, value.String()))
}

func (fa *fieldAdapter) AddTime(key string, value time.Time) {
	// TODO inefficient
	*fa = append(*fa, log.String(key, value.String()))
}

func (fa *fieldAdapter) AddBinary(key string, value []byte) {
	*fa = append(*fa, log.Object(key, value))
}

func (fa *fieldAdapter) AddByteString(key string, value []byte) {
	*fa = append(*fa, log.Object(key, value))
}

func (fa *fieldAdapter) AddString(key, value string) {
	if key != "" && value != "" {
		*fa = append(*fa, log.String(key, value))
	}
}
