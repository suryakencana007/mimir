/*  uptime_test.go
*
* @Author:             Nanang Suryadi
* @Date:               April 05, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 05/04/20 04:54
 */

package web

import (
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUptimeHandler_ServeHTTP(t *testing.T) {

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/uptime", strings.NewReader(""))

	// Tested code
	DefaultUptimeHandler().ServeHTTP(resp, req)

	var response map[string]interface{}

	err := json.NewDecoder(resp.Body).Decode(&response)

	// Asserts
	assert.NoError(t, err)
	delta, err := strconv.ParseFloat(response["data"].(string), 64)
	assert.Nil(t, err)
	assert.InDelta(t, 0, delta, 0.1) // should be within 0.1 second of 0
}
