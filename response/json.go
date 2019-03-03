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

type ErrorValidation struct {
    Errors interface{} `json:"errors"`
}

// APIResponse defines attributes for api Response
type APIResponse struct {
    HTTPCode   int         `json:"-"`
    Code       int         `json:"code"`
    Message    interface{} `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    Pagination interface{} `json:"pagination,omitempty"`
}

type Pagination struct {
    Page  int `json:"page"`
    Size  int `json:"size"`
    Total int `json:"total"`
}

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
    w.Write(buf.Bytes())
}
