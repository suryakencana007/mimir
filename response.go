/*  response.go
*
* @Author:             Nanang Suryadi
* @Date:               November 21, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 21/11/19 22:02
 */

package mimir

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type ctxKeyVersion struct {
	Name string
}

func (r *ctxKeyVersion) String() string {
	return "context value " + r.Name
}

type ctxKeyResponse struct {
	Name string
}

func (r *ctxKeyResponse) String() string {
	return "context value " + r.Name
}

var (
	CtxResponse = ctxKeyResponse{Name: "context response"}
	CtxVersion  = ctxKeyVersion{Name: "context version"}
)

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type ErrorValidator struct {
	Type    string `json:"error_type,omitempty"`
	Tag     string `json:"error_tag,omitempty"`
	Field   string `json:"error_field,omitempty"`
	Value   string `json:"error_value,omitempty"`
	Message string `json:"error_message,omitempty"`
}

type Meta struct {
	Code    string `json:"code,omitempty"`
	Type    string `json:"error_type,omitempty"`
	Message string `json:"error_message,omitempty"`
}

type Version struct {
	Label  string `json:"label,omitempty"`
	Number string `json:"number,omitempty"`
}

type response struct {
	Version    interface{} `json:"version,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func Response(r *http.Request) *response {
	null := make(map[string]interface{})
	resp := &response{
		Version: Version{
			Label:  "v1",
			Number: "0.1.0",
		},
		Meta:       null,
		Data:       null,
		Pagination: null,
	}
	if ver, ok := r.Context().Value(CtxVersion).(Version); ok {
		resp.Version = ver
	}
	return resp
}

func (r *response) Errors(err ...Meta) *response {
	r.Meta = err
	return r
}

func (r *response) Success(code int) *response {
	r.Meta = Meta{Code: StatusText(code)}
	return r
}

func (r *response) Body(body interface{}) {
	r.Data = body
}

func (r *response) Page(p Pagination) {
	r.Pagination = p
}

// APIStatusSuccess for standard request api status success
func (r *response) APIStatusSuccess(w http.ResponseWriter, req *http.Request) *responseWriter {
	r.Success(StatusSuccess)
	return Status(w, req, StatusSuccess, r)
}

// APIStatusCreated
func (r *response) APIStatusCreated(w http.ResponseWriter, req *http.Request) *responseWriter {
	r.Success(StatusCreated)
	return Status(w, req, StatusCreated, r)
}

// APIStatusAccepted
func (r *response) APIStatusAccepted(w http.ResponseWriter, req *http.Request) *responseWriter {
	r.Success(StatusAccepted)
	return Status(w, req, StatusAccepted, r)
}

// APIStatusPermanentRedirect
func (r *response) APIStatusPermanentRedirect(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusPermanentRedirect),
		Type:    StatusCode(StatusPermanentRedirect),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusPermanentRedirect), err.Error()),
	})
	return Status(w, req, StatusPermanentRedirect, r)
}

// APIStatusBadRequest
func (r *response) APIStatusBadRequest(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusBadRequest),
		Type:    StatusCode(StatusBadRequest),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusBadRequest), err.Error()),
	})
	return Status(w, req, StatusBadRequest, r)
}

// APIStatusUnauthorized
func (r *response) APIStatusUnauthorized(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusUnauthorized),
		Type:    StatusCode(StatusUnauthorized),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusUnauthorized), err.Error()),
	})
	return Status(w, req, StatusUnauthorized, r)
}

// APIStatusPaymentRequired
func (r *response) APIStatusPaymentRequired(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusPaymentRequired),
		Type:    StatusCode(StatusPaymentRequired),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusPaymentRequired), err.Error()),
	})
	return Status(w, req, StatusPaymentRequired, r)
}

// APIStatusForbidden
func (r *response) APIStatusForbidden(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusForbidden),
		Type:    StatusCode(StatusForbidden),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusForbidden), err.Error()),
	})
	return Status(w, req, StatusForbidden, r)
}

// APIStatusMethodNotAllowed
func (r *response) APIStatusMethodNotAllowed(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusMethodNotAllowed),
		Type:    StatusCode(StatusMethodNotAllowed),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusMethodNotAllowed), err.Error()),
	})
	return Status(w, req, StatusMethodNotAllowed, r)
}

// APIStatusNotAcceptable
func (r *response) APIStatusNotAcceptable(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusNotAcceptable),
		Type:    StatusCode(StatusNotAcceptable),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusNotAcceptable), err.Error()),
	})
	return Status(w, req, StatusNotAcceptable, r)
}

// APIStatusInvalidAuthentication
func (r *response) APIStatusInvalidAuthentication(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusInvalidAuthentication),
		Type:    StatusCode(StatusInvalidAuthentication),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusInvalidAuthentication), err.Error()),
	})
	return Status(w, req, StatusInvalidAuthentication, r)
}

// APIStatusRequestTimeout
func (r *response) APIStatusRequestTimeout(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusRequestTimeout),
		Type:    StatusCode(StatusRequestTimeout),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusRequestTimeout), err.Error()),
	})
	return Status(w, req, StatusRequestTimeout, r)
}

// APIStatusUnsupportedMediaType
func (r *response) APIStatusUnsupportedMediaType(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    StatusCode(StatusUnsupportedMediaType),
		Type:    StatusCode(StatusUnsupportedMediaType),
		Message: err.Error(),
	})
	return Status(w, req, StatusUnsupportedMediaType, r)
}

// APIStatusUnProcess
func (r *response) APIStatusUnProcess(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusUnProcess),
		Type:    StatusCode(StatusUnProcess),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusUnProcess), err.Error()),
	})
	return Status(w, req, StatusUnProcess, r)
}

// APIStatusInternalError
func (r *response) APIStatusInternalError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusInternalError),
		Type:    StatusCode(StatusInternalError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusInternalError), err.Error()),
	})
	return Status(w, req, StatusInternalError, r)
}

// APIStatusBadGatewayError
func (r *response) APIStatusBadGatewayError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusBadGatewayError),
		Type:    StatusCode(StatusBadGatewayError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusBadGatewayError), err.Error()),
	})
	return Status(w, req, StatusBadGatewayError, r)
}

// APIStatusServiceUnavailableError
func (r *response) APIStatusServiceUnavailableError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusServiceUnavailableError),
		Type:    StatusCode(StatusServiceUnavailableError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusServiceUnavailableError), err.Error()),
	})
	return Status(w, req, StatusServiceUnavailableError, r)
}

// APIStatusGatewayTimeoutError
func (r *response) APIStatusGatewayTimeoutError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusGatewayTimeoutError),
		Type:    StatusCode(StatusGatewayTimeoutError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusGatewayTimeoutError), err.Error()),
	})
	return Status(w, req, StatusGatewayTimeoutError, r)
}

type responseWriter struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Response *response
}

func Status(w http.ResponseWriter, r *http.Request, status int, v *response) *responseWriter {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxResponse, status))
	return &responseWriter{
		Request:  r,
		Writer:   w,
		Response: v,
	}
}

func SemanticVersion(r *http.Request, label string, version string) {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxVersion, Version{
		Label:  label,
		Number: version,
	}))
}
