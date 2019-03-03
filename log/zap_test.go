/*  zap_test.go
*
* @Author:             Nanang Suryadi
* @Date:               February 06, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-06 02:07 
 */

package log

import (
    "bytes"
    "fmt"
    "io"
    "net"
    "os"
    "strings"
    "testing"
    "time"

    "github.com/go-stack/stack"
    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func getZapStackCaller() string {
    strGoStack := "%n"
    stack := stack.Caller(2)
    return strings.Join([]string{
        fmt.Sprint(`"caller":"`, stack, `"`),
        fmt.Sprintf(`"function":"%s"`, fmt.Sprintf(strGoStack, stack)),
    }, ",")
}

func newZap(t *testing.T) (*zapLog, *testLogSpy) {
    ts := newTestLogSpy(t)
    writer := newTestingWriter(ts)
    return &zapLog{
        logger: zap.New(DefaultCore(writer)),
    }, ts
}

func newZapStdout(t *testing.T) *zapLog {
    return &zapLog{
        logger: zap.New(DefaultCore(nil)),
    }
}

func zapTestInit(t *testing.T) *testLogSpy {
    ts := newTestLogSpy(t)
    writer := newTestingWriter(ts)
    logger = NewZap(DefaultCore(writer))
    return ts
}

func TestZapInit(t *testing.T) {
    old := os.Stdout // keep backup of the real stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    ZapInit()
    logger.Info("received work order")
    logger.Debug("starting work")
    logger.Warn("work may fail")
    logger.Error("work failed")

    outC := make(chan string)
    go func() {
        var buf bytes.Buffer
        _, err := io.Copy(&buf, r)
        if err == nil {
            outC <- buf.String()
        }
    }()
    w.Close()
    os.Stdout = old
    out := <-outC
    expected := []string{
        `{"level":"info","msg":"received work order",` + getZapStackCaller() + `}`,
        `{"level":"debug","msg":"starting work",` + getZapStackCaller() + `}`,
        `{"level":"warn","msg":"work may fail",` + getZapStackCaller() + `}`,
        `{"level":"error","msg":"work failed",` + getZapStackCaller() + `}`,
    }
    assert.Equal(t,
        strings.Join(expected, "\n")+"\n",
        out)
    assert.Panics(t, func() {
        logger.Panic("failed to do work", Field("Error", "Should panic"))
    }, "log.Panic should panic")
}

func TestInternalZap(t *testing.T) {
    log, ts := newZap(t)
    log.Info("received work order")
    log.Debug("starting work")
    log.Warn("work may fail")
    log.Error("work failed")

    ts.AssertMessages(
        `{"level":"info","msg":"received work order",`+getZapStackCaller()+`}`,
        `{"level":"debug","msg":"starting work",`+getZapStackCaller()+`}`,
        `{"level":"warn","msg":"work may fail",`+getZapStackCaller()+`}`,
        `{"level":"error","msg":"work failed",`+getZapStackCaller()+`}`,
    )

    assert.Panics(t, func() {
        log.Panic("failed to do work", Field("Error", "Should panic"))
    }, "log.Panic should panic")

}

func TestInternalStdout(t *testing.T) {
    old := os.Stdout // keep backup of the real stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    log := newZapStdout(t)
    log.Info("received work order")
    log.Debug("starting work")
    log.Warn("work may fail")
    log.Error("work failed")

    outC := make(chan string)
    go func() {
        var buf bytes.Buffer
        _, err := io.Copy(&buf, r)
        if err == nil {
            outC <- buf.String()
        }
    }()
    w.Close()
    os.Stdout = old
    out := <-outC
    expected := []string{
        `{"level":"info","msg":"received work order",` + getZapStackCaller() + `}`,
        `{"level":"debug","msg":"starting work",` + getZapStackCaller() + `}`,
        `{"level":"warn","msg":"work may fail",` + getZapStackCaller() + `}`,
        `{"level":"error","msg":"work failed",` + getZapStackCaller() + `}`,
    }
    assert.Equal(t,
        strings.Join(expected, "\n")+"\n",
        out)

    assert.Panics(t, func() {
        log.Panic("failed to do work", Field("Error", "Should panic"))
    }, "log.Panic should panic")
}

func TestZapField(t *testing.T) {
    addr := net.ParseIP("1.2.3.4")
    name := username("phil")

    tests := []struct {
        name   string
        want   interface{}
        expect zap.Field
        ftype  zapcore.FieldType
    }{
        {name: "Binary", want: []byte("ab12"), expect: zap.Binary("k", []byte("ab12")), ftype: zapcore.BinaryType},
        {name: "Bool", want: true, expect: zap.Bool("k", true), ftype: zapcore.BoolType},
        {name: "Complex128", want: 1 + 2i, expect: zap.Complex128("k", 1+2i), ftype: zapcore.Complex128Type},
        {name: "Complex64", want: complex64(1 + 2i), expect: zap.Complex64("k", 1+2i), ftype: zapcore.Complex64Type},
        {name: "Duration", want: time.Duration(1), expect: zap.Duration("k", 1), ftype: zapcore.DurationType},
        {name: "Int", want: 1, expect: zap.Int("k", 1), ftype: zapcore.Int64Type},
        {name: "Int64", want: int64(1), expect: zap.Int64("k", 1), ftype: zapcore.Int64Type},
        {name: "Int32", want: int32(1), expect: zap.Int32("k", 1), ftype: zapcore.Int32Type},
        {name: "Int16", want: int16(1), expect: zap.Int16("k", 1), ftype: zapcore.Int16Type},
        {name: "Int8", want: int8(1), expect: zap.Int8("k", 1), ftype: zapcore.Int8Type},
        {name: "String", want: "foo", expect: zap.String("k", "foo"), ftype: zapcore.StringType},
        {name: "Time", want: time.Unix(0, 0), expect: zap.Time("k", time.Unix(0, 0).In(time.UTC)), ftype: zapcore.TimeType},
        {name: "Time", want: time.Unix(0, 1000), expect: zap.Time("k", time.Unix(0, 1000).In(time.UTC)), ftype: zapcore.TimeType},

        {name: "Uint", want: uint(1), expect: zap.Uint("k", 1), ftype: zapcore.Uint64Type},
        {name: "Uint64", want: uint64(1), expect: zap.Uint64("k", 1), ftype: zapcore.Uint64Type},
        {name: "Uint32", want: uint32(1), expect: zap.Uint32("k", 1), ftype: zapcore.Uint32Type},
        {name: "Uint16", want: uint16(1), expect: zap.Uint16("k", 1), ftype: zapcore.Uint16Type},
        {name: "Uint8", want: uint8(1), expect: zap.Uint8("k", 1), ftype: zapcore.Uint8Type},

        {name: "Uintptr", want: uintptr(10), expect: zap.Uintptr("k", 0xa), ftype: zapcore.UintptrType},
        {name: "Stringer", want: addr, expect: zap.Stringer("k", addr), ftype: zapcore.StringerType},
        {name: "object", want: name, expect: zap.Object("k", name), ftype: zapcore.ObjectMarshalerType},
        {name: "Reflect", want: []interface{}{}, expect: zap.Reflect("k", []interface{}{}), ftype: zapcore.ReflectType},
    }

    log := newZapStdout(t)
    for _, tt := range tests {
        fieldType := log.Field("k", tt.want).(zapcore.Field)
        assert.Equal(t, tt.expect.Type, fieldType.Type, "Expected output from field %+v - %+v.", fieldType.Type, tt.name)

        assert.True(t, fieldType.Type == tt.ftype, "Field does not equal itself %+v - %+v.", fieldType.Type, tt.name)
    }
}

type username string

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("username", string(n))
    return nil
}
