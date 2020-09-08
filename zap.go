/*  zap.go
** @Date:               November 21, 2019
** @Last Modified time: 21/11/19 07:35
 */

package mimir

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapField = zapcore.Field
type Core = zapcore.Core

func typeField(fields ...interface{}) []ZapField {
	f := make([]ZapField, 0)
	for _, v := range fields {
		val := v.([]interface{})
		if len(val) > 0 {
			for _, field := range val {
				f = append(f, field.(ZapField))
			}
		}
	}

	return f
}

// zapLog is the logger struct
type zapLog struct {
	logger *zap.Logger
}

func (z *zapLog) Debugf(msg string, args ...interface{}) {
	z.Validator().Debug(fmt.Sprintf(msg, args...))
}

func (z *zapLog) Infof(msg string, args ...interface{}) {
	z.Validator().Info(fmt.Sprintf(msg, args...))
}

func (z *zapLog) Warnf(msg string, args ...interface{}) {
	z.Validator().Warn(fmt.Sprintf(msg, args...))
}

func (z *zapLog) Errorf(msg string, args ...interface{}) {
	z.Validator().Error(fmt.Sprintf(msg, args...))
}

func (z *zapLog) With(fields ...interface{}) Logging {
	z.logger = z.logger.With(typeField(fields)...)
	return z
}

func (z *zapLog) Validator() *zap.Logger {
	if z.logger == nil {
		log.Fatal("please initiate properly")
	}
	return z.logger
}

func (z *zapLog) Debug(msg string, fields ...interface{}) {
	z.Validator().Debug(msg, typeField(fields)...)
}

func (z *zapLog) Info(msg string, fields ...interface{}) {
	z.Validator().Info(msg, typeField(fields)...)
}

func (z *zapLog) Warn(msg string, fields ...interface{}) {
	z.Validator().Warn(msg, typeField(fields)...)
}

func (z *zapLog) Error(msg string) {
	z.Validator().Error(msg)
}

func (z *zapLog) Fatal(msg string, fields ...interface{}) {
	z.Validator().Fatal(msg, typeField(fields)...)
}

func (z *zapLog) Panic(msg string, fields ...interface{}) {
	z.Validator().Panic(msg, typeField(fields)...)
}

func (z *zapLog) Field(key string, value interface{}) interface{} {
	switch val := value.(type) {
	case zapcore.ObjectMarshaler:
		return zap.Object(key, val)
	case zapcore.ArrayMarshaler:
		return zap.Array(key, val)
	case bool:
		return zap.Bool(key, val)
	case []bool:
		return zap.Bools(key, val)
	case complex128:
		return zap.Complex128(key, val)
	case []complex128:
		return zap.Complex128s(key, val)
	case complex64:
		return zap.Complex64(key, val)
	case []complex64:
		return zap.Complex64s(key, val)
	case float64:
		return zap.Float64(key, val)
	case []float64:
		return zap.Float64s(key, val)
	case float32:
		return zap.Float32(key, val)
	case []float32:
		return zap.Float32s(key, val)
	case int:
		return zap.Int(key, val)
	case []int:
		return zap.Ints(key, val)
	case int64:
		return zap.Int64(key, val)
	case []int64:
		return zap.Int64s(key, val)
	case int32:
		return zap.Int32(key, val)
	case []int32:
		return zap.Int32s(key, val)
	case int16:
		return zap.Int16(key, val)
	case []int16:
		return zap.Int16s(key, val)
	case int8:
		return zap.Int8(key, val)
	case []int8:
		return zap.Int8s(key, val)
	case string:
		return zap.String(key, val)
	case []string:
		return zap.Strings(key, val)
	case uint:
		return zap.Uint(key, val)
	case []uint:
		return zap.Uints(key, val)
	case uint64:
		return zap.Uint64(key, val)
	case []uint64:
		return zap.Uint64s(key, val)
	case uint32:
		return zap.Uint32(key, val)
	case []uint32:
		return zap.Uint32s(key, val)
	case uint16:
		return zap.Uint16(key, val)
	case []uint16:
		return zap.Uint16s(key, val)
	case uint8:
		return zap.Uint8(key, val)
	case []byte:
		return zap.Binary(key, val)
	case uintptr:
		return zap.Uintptr(key, val)
	case []uintptr:
		return zap.Uintptrs(key, val)
	case time.Time:
		return zap.Time(key, val)
	case []time.Time:
		return zap.Times(key, val)
	case time.Duration:
		return zap.Duration(key, val)
	case []time.Duration:
		return zap.Durations(key, val)
	case error:
		return zap.NamedError(key, val)
	case []error:
		return zap.Errors(key, val)
	case fmt.Stringer:
		return zap.Stringer(key, val)
	default:
		return zap.Reflect(key, val)
	}
}

// New creates new instance of zapLog
func NewZap(c Core) *zapLog {
	return &zapLog{
		logger: zap.New(c, zap.AddCaller(), zap.AddCallerSkip(3)),
	}
}

func NewZapProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func ProductionCore() Core {
	defaultEncoder := zapcore.NewJSONEncoder(NewZapProductionEncoderConfig())
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	debugging := zapcore.Lock(os.Stdout)
	errors := zapcore.Lock(os.Stderr)
	return zapcore.NewTee(
		zapcore.NewCore(defaultEncoder, errors, highPriority),
		zapcore.NewCore(defaultEncoder, debugging, lowPriority),
	)
}

func DefaultCore(out zapcore.WriteSyncer) Core {
	if out == nil {
		out = os.Stdout
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(NewZapProductionEncoderConfig()),
		out,
		zap.DebugLevel,
	)
}
