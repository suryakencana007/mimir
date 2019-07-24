/*  zap.go
*
* @Author:             Nanang Suryadi
* @Date:               February 12, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-12 11:39
 */

package log

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/go-stack/stack"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type ZapField = zapcore.Field
type Core = zapcore.Core

// zapLog is the logger struct
type zapLog struct {
    logger *zap.Logger
}

func ZapInit() {
    logger = NewZap(ProductionCore())
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

// New creates new instance of zapLog
func NewZap(c Core) logging {
    return &zapLog{
        logger: zap.New(c),
    }
}

func (l *zapLog) validator() {
    if l.logger == nil {

        log.Fatal("please initiate properly")
    }
}

// Debug add log entry with debug level
func (l *zapLog) Debug(msg string, fields ...interface{}) {
    l.validator()
    l.logger.Debug(msg, fieldType(fields)...)
}

// Info add log entry with info level
func (l *zapLog) Info(msg string, fields ...interface{}) {
    l.validator()
    l.logger.Info(msg, fieldType(fields)...)
}

// Warn add log entry with warn level
func (l *zapLog) Warn(msg string, fields ...interface{}) {
    l.validator()
    l.logger.Warn(msg, fieldType(fields)...)
}

// Error add log entry with error level
func (l *zapLog) Error(msg string, fields ...interface{}) {
    l.validator()
    l.logger.Error(msg, fieldType(fields)...)
}

// Fatal add log entry with fatal level
func (l *zapLog) Fatal(msg string, fields ...interface{}) {
    l.validator()
    l.logger.Fatal(msg, fieldType(fields)...)
}

// Panic add log entry with panic level
func (l *zapLog) Panic(msg string, fields ...interface{}) {
    l.validator()
    l.logger.Panic(msg, fieldType(fields)...)
}

func fieldType(fields ...interface{}) []ZapField {
    var f []ZapField

    for _, v := range fields {
        val := v.([]interface{})
        if len(val) > 0 {
            for _, field := range val {
                f = append(f, field.(ZapField))
            }
        }
    }
    strGoStack := "%n"
    stack := stack.Caller(3)
    f = append(f,
        ZapField{Key: "caller", Type: zapcore.ReflectType, Interface: stack},
        ZapField{Key: "function", Type: zapcore.StringType, String: fmt.Sprintf(strGoStack, stack)},
    )
    return f
}

func (l *zapLog) Field(key string, value interface{}) interface{} {
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
