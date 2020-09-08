/*  router_test.go
*
* @Date:               April 22, 2020
* @Last Modified time: 22/04/20 01:36
 */

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suryakencana007/mimir/example/memorist/config"
	"github.com/suryakencana007/mimir/ruuto"
)

func TestRouter(t *testing.T) {
	// setup
	var uptimeCalled, healthCalled bool
	stubHealthzHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		healthCalled = true
		uptimeCalled = false
	})

	handlers := Handlers{
		Config: &config.Config{
			Main: struct {
				Prefix string `mapstructure:"prefix"`
			}{Prefix: "/v1"},
			API: struct {
				Prefix string `mapstructure:"prefix"`
			}{Prefix: "/api"},
		},
		Router: ruuto.NewChiRouter(),
		Health: stubHealthzHandler,
	}
	router := Router(handlers)

	// Tested code 1
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/v1/api/healthz", strings.NewReader("")))

	// Asserts 1
	assert.True(t, healthCalled)
	assert.False(t, uptimeCalled)
}

func TestHealthHandler_ServeHTTP(t *testing.T) {

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/api/healthz", strings.NewReader(""))

	// Tested code
	Health().ServeHTTP(resp, req)

	var response map[string]interface{}

	err := json.NewDecoder(resp.Body).Decode(&response)
	// Asserts
	assert.NoError(t, err)
	assert.Nil(t, err)
}
