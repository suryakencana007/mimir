/*  aws.go
*
* @Author:             Nanang Suryadi
* @Date:               February 13, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-13 11:23
 */

package response

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/suryakencana007/mimir/constant"
	"github.com/suryakencana007/mimir/log"
)

type ErrorData struct {
	Code    string `json:"code"`
	Key     string `json:"key"`
	Message string `json:"message"`
}

type ErrorContext struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
	Data    []ErrorData `json:"data"`
}

type AWSResponse struct {
	Errors     interface{} `json:"errors"` // ErrorContext
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}

func NewAWSResponse() AWSResponse {
	null := make(map[string]interface{})
	return AWSResponse{
		Data:       null,
		Errors:     null,
		Pagination: null,
	}
}

func (a *AWSResponse) SetErrors(err ErrorContext) {
	a.Errors = err
}

func (a *AWSResponse) SetData(v interface{}) {
	a.Data = v
}

type Response events.APIGatewayProxyResponse

func NewResponse() Response {
	return Response{
		StatusCode:      constant.StatusSuccess,
		IsBase64Encoded: false,
	}
}

func (r *Response) SetCode(code int) {
	r.StatusCode = code
}

func (r *Response) SetHeader(headers map[string]string) {
	headers["X-COMPANY-FUNCTION"] = "warungpintar.co"
	r.Headers = headers
}

func (r *Response) SetBody(w string) {
	r.Body = w
}

/**
  ErrorContext{
      Data:    make([]ErrorData, 0),
      Message: make(map[string]interface{}, 0),
  }
*/

// Write writes the data to http response writer
func Write(v interface{}) Response {
	response := NewResponse()

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		log.Error("Error Response", err.Error())
		response.SetCode(constant.StatusInternalError)
		response.SetBody(buf.String())
		return response
	}
	log.Info("Response")
	response.SetBody(buf.String())
	return response
}
