/*  router_test.go
*
* @Author:             Nanang Suryadi
* @Date:               April 05, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 05/04/20 04:47
 */

package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	// Setup
	var uptimeCalled, counterCalled bool
	stubUptimeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uptimeCalled = true
		counterCalled = true
	})

	handlers := Handlers{
		Uptime: stubUptimeHandler,
	}
	router := MainRouter(handlers)

	// Tested code 1
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/uptime", strings.NewReader("")))

	// Asserts 1
	assert.True(t, uptimeCalled)
	assert.True(t, counterCalled)
}
