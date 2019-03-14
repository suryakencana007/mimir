/*  response.go
*
* @Author:             Nanang Suryadi
* @Date:               March 15, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-03-15 04:19 
 */

package response

import (
    "context"
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

func Status(r *http.Request, status int) {
    *r = *r.WithContext(context.WithValue(r.Context(), constant.StatusCtxKey, status))
}
