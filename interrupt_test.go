/*  interrupt_test.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 14:24
 */

package mimir

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterruptChannelFunc(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Interrupt test not supported on windows")
	}

	// Setup
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		assert.Fail(t, "failed to find my process: %v", err)
	}

	// Tested Code
	interrupt := InterruptChannelFunc()
	if err := process.Signal(os.Interrupt); err != nil {
		assert.Fail(t, "failed to send interrupt signal: %v", err)
	}

	// Asserts
	assert.Equal(t, os.Interrupt, <-interrupt)
}
