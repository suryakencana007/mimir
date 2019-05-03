/*  logger_test.go
*
* @Author:             Nanang Suryadi
* @Date:               January 31, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-01-31 16:46
 */

package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/go-stack/stack"
	"github.com/stretchr/testify/assert"
)

func getStackCaller(line int, f string) string {
	strGoStack := "%n"
	stack := stack.Caller(1)
	return strings.Join([]string{
		fmt.Sprintf(`"caller":"%s:%d"`, stack, line),
		fmt.Sprintf(`"function":"%s%s"`, fmt.Sprintf(strGoStack, stack), f),
	}, ",")
}

func TestZapLogger(t *testing.T) {
	ts := zapTestInit(t)
	Info("received work order", Field("Key", "Key Field"))
	Debug("starting work", Field("Key", "Key Field"))
	Warn("work may fail", Field("Key", "Key Field"))
	Error("work failed", Field("Key", "Key Field"))

	assert.Panics(t, func() {
		Panic("failed to do work", Field("Error", "Should panic"))
	}, "log.Panic should panic")

	ts.AssertMessages(
		`{"level":"info","msg":"received work order","Key":"Key Field",`+getStackCaller(32, "")+`}`,
		`{"level":"debug","msg":"starting work","Key":"Key Field",`+getStackCaller(33, "")+`}`,
		`{"level":"warn","msg":"work may fail","Key":"Key Field",`+getStackCaller(34, "")+`}`,
		`{"level":"error","msg":"work failed","Key":"Key Field",`+getStackCaller(35, "")+`}`,
		`{"level":"panic","msg":"failed to do work","Error":"Should panic",`+getStackCaller(38, ".func1")+`}`,
	)
}

// testLogSpy is a testing.TB that captures logged messages.
type testLogSpy struct {
	testing.TB
	Messages []string
}

func newTestLogSpy(t testing.TB) *testLogSpy {
	return &testLogSpy{TB: t}
}

func (t *testLogSpy) Logf(format string, args ...interface{}) {
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *testLogSpy) AssertMessages(msgs ...string) {
	assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}

// testingWriter is a WriteSyncer that writes to the given testing.TB.
type testingWriter struct {
	t testing.TB

	// If true, the test will be marked as failed if this testingWriter is
	// ever used.
	markFailed bool
}

func newTestingWriter(t testing.TB) testingWriter {
	return testingWriter{t: t}
}

// WithMarkFailed returns a copy of this testingWriter with markFailed set to
// the provided value.
func (w testingWriter) WithMarkFailed(v bool) testingWriter {
	w.markFailed = v
	return w
}

func (w testingWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	// Strip trailing newline because want.Log always adds one.
	p = bytes.TrimRight(p, "\n")
	// Note: want.Log is safe for concurrent use.
	w.t.Logf("%s", p)
	if w.markFailed {
		w.t.Fail()
	}

	return n, nil
}

func (w testingWriter) Sync() error {
	return nil
}
