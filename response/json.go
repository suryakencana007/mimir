/*  response.go
*
* @Author:             Nanang Suryadi
* @Date:               October 01, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 01/10/18 00:21
 */

package response

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/suryakencana007/mimir/constant"
)

// Write writes the data to http response writer
func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(constant.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	_, err := w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
