/*  http_test.go
*
* @Author:             Nanang Suryadi
* @Date:               November 21, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 21/11/19 17:41
 */

package mimir

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suryakencana007/mimir/ruuto"
)

var portSync sync.Mutex

func findOpenPort() (int, error) {
	portSync.Lock()
	defer portSync.Unlock()

	min := 10000
	max := 65535
	attempts := 10

	for i := 0; i < attempts; i++ {
		bg := big.NewInt(int64(max - min))
		n, err := rand.Int(rand.Reader, bg)
		if err != nil {
			continue
		}
		port := n.Int64() + int64(min)
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
			// Port unavailable
			continue
		} else if err := ln.Close(); err != nil {
			return 0, err
		}
		return int(port), nil
	}
	return 0, fmt.Errorf("could not find port to use for testing (%d attempts)", attempts)
}

func TestHttpListenAndServe(t *testing.T) {
	logger := With(
		Field("logger suki", "TestHttpListenAndServe"),
		Field("testing", "TestHttpListenAndServe"),
	)

	port, err := findOpenPort()
	if err != nil {
		assert.Fail(t, "could not find a testing port")
	}
	t.Log("Using port", port)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Testing
	runServer, cleanup := ListenAndServe(ServeOpts{
		Logger:  logger,
		Port:    WebPort(port),
		Router:  ruuto.NewChiRouter(),
		TimeOut: WebTimeOut(100),
	})
	serverErrChan := make(chan error)
	go func() {
		serverErrChan <- runServer(ctx)
	}()

	resp, queryErr := http.Get(fmt.Sprintf("http://localhost:%d", port))

	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	assert.Nil(t, queryErr)
	respStatus := resp.StatusCode
	assert.Equal(t, http.StatusNotFound, respStatus)
	cancel()

	select {
	case serverErr := <-serverErrChan:
		assert.EqualError(t, serverErr, "server interrupted through context")
		cleanup(ctx)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "server never quit")
	}
}
