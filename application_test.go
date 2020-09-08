/*  application_test.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 14:16
 */

package mimir

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	interrupt := make(chan os.Signal, 2)
	interrupt <- os.Interrupt

	var cleanupCalled, runnerInvoke bool
	cleanup, err := Application(interrupt, func(c context.Context) (AppRunner, func(), error) {
		runner := func(ctx context.Context) error {
			runnerInvoke = false
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return fmt.Errorf("server completed without interruption")
			}
		}
		return runner, func() {
			cleanupCalled = true
			runnerInvoke = true
		}, nil
	})

	assert.Error(t, err)
	assert.NotNil(t, cleanup)
	assert.False(t, runnerInvoke)
	cleanup()
	assert.True(t, cleanupCalled)
	assert.True(t, runnerInvoke)
}

func TestApplicationFatalError(t *testing.T) {
	interrupt := make(chan os.Signal, 2)

	var cleanupCalled, runnerInvoke bool
	cleanup, err := Application(interrupt, func(c context.Context) (AppRunner, func(), error) {
		runner := func(ctx context.Context) error {
			runnerInvoke = false
			return fmt.Errorf("server had an fatal error")
		}
		return runner, func() {
			cleanupCalled = true
			runnerInvoke = true
		}, nil
	})

	assert.NotNil(t, err)
	assert.NotNil(t, cleanup)
	assert.False(t, runnerInvoke)
	cleanup()
	assert.True(t, cleanupCalled)
	assert.True(t, runnerInvoke)
}
