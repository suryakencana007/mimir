/*  http.go
*
* @Author:             Nanang Suryadi
* @Date:               November 21, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 21/11/19 12:28
 */

package mimir

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/suryakencana007/mimir/ruuto"
)

func Welkommen() string {
	return `
========================================================================================
   _     _     _     _     _  
  / \   / \   / \   / \   / \ 
 ( m ) ( i ) ( m ) ( i ) ( r )
  \_/   \_/   \_/   \_/   \_/ 
========================================================================================
- port    : %d
-----------------------------------------------------------------------------------------
`
}

type (
	WebPort       int
	WebTimeOut    int
	Https         bool
	ServerRunFunc func(context.Context) error
	ServeOpts     struct {
		Logger   Logging
		Port     WebPort
		Router   ruuto.Router
		TimeOut  WebTimeOut
		TLS      Https
		CertFile string
		KeyFile  string
	}
)

func ListenAndServe(opts ServeOpts) (ServerRunFunc, func(context.Context)) {
	logger := opts.Logger
	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", opts.Port),
		Handler:      opts.Router,
		ReadTimeout:  time.Duration(opts.TimeOut) * time.Second,
		WriteTimeout: time.Duration(opts.TimeOut) * time.Second,
	}

	cleanup := func(ctx context.Context) {
		logger.Info("I have to go...")
		logger.Info("Stopping server gracefully")
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Debug("Wait is over due to error")
			if err = httpServer.Close(); err != nil {
				logger.Debug(err.Error())
			}
		}
		logger.Info(fmt.Sprintf("Stop server at %s", httpServer.Addr))
	}

	return func(ctx context.Context) error {
		errChan := make(chan error)
		go func() {
			// Description µ micro service
			fmt.Println(
				fmt.Sprintf(
					Welkommen(),
					opts.Port,
				))
			logger.Info(fmt.Sprintf("Now serving at %s", httpServer.Addr))
			if opts.TLS {
				logger.Info("Secure with HTTPS")
				errChan <- httpServer.ListenAndServeTLS(opts.CertFile, opts.KeyFile)
			} else {
				errChan <- httpServer.ListenAndServe()
			}
		}()

		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			_ = httpServer.Close()
			return fmt.Errorf("server interrupted through context")
		}
	}, cleanup
}

type MediaType string

const (
	ApplicationJSON MediaType = "application/json"
	FormURLEncoded  MediaType = "application/x-www-form-urlencoded"
	MultipartForm   MediaType = "multipart/form-data"
	TextPlain       MediaType = "text/plain"
)

const (
	// 2xx
	StatusSuccess  = http.StatusOK
	StatusCreated  = http.StatusCreated
	StatusAccepted = http.StatusAccepted
	// 3xx
	StatusPermanentRedirect = http.StatusPermanentRedirect
	// 4xx
	StatusBadRequest            = http.StatusBadRequest
	StatusUnauthorized          = http.StatusUnauthorized
	StatusPaymentRequired       = http.StatusPaymentRequired
	StatusForbidden             = http.StatusForbidden
	StatusMethodNotAllowed      = http.StatusMethodNotAllowed
	StatusNotAcceptable         = http.StatusNotAcceptable
	StatusInvalidAuthentication = http.StatusProxyAuthRequired
	StatusRequestTimeout        = http.StatusRequestTimeout
	StatusUnsupportedMediaType  = http.StatusUnsupportedMediaType
	StatusUnProcess             = http.StatusUnprocessableEntity
	//5xx
	StatusInternalError           = http.StatusInternalServerError
	StatusBadGatewayError         = http.StatusBadGateway
	StatusServiceUnavailableError = http.StatusServiceUnavailable
	StatusGatewayTimeoutError     = http.StatusGatewayTimeout
)

var statusMap = map[int][]string{
	StatusSuccess:  {"STATUS_OK", "Success"},
	StatusCreated:  {"STATUS_CREATED", "Resource has been created"},
	StatusAccepted: {"STATUS_ACCEPTED", "Resource has been accepted"},

	StatusPermanentRedirect: {"STATUS_PERMANENT_REDIRECT", "The resource has moved to a new location"},

	StatusBadRequest:            {"STATUS_BAD_REQUEST", "Invalid data request"},
	StatusUnauthorized:          {"STATUS_UNAUTHORIZED", "Not authorized to access the service"},
	StatusPaymentRequired:       {"STATUS_PAYMENT_REQUIRED", "Payment need to be done"},
	StatusForbidden:             {"STATUS_FORBIDDEN", "Forbidden access the resource "},
	StatusMethodNotAllowed:      {"STATUS_METHOD_NOT_ALLOWED", "The method specified is not allowed"},
	StatusNotAcceptable:         {"STATUS_NOT_ACCEPTABLE", "Request cannot accepted"},
	StatusInvalidAuthentication: {"STATUS_INVALID_AUTHENTICATION", "The resource owner or authorization server denied the request"},
	StatusRequestTimeout:        {"STATUS_REQUEST_TIMEOUT", "Request Timeout"},
	StatusUnsupportedMediaType:  {"STATUS_UNSUPPORTED_MEDIA_TYPE", "Cannot understand request content"},
	StatusUnProcess:             {"STATUS_UNPROCESSABLE_ENTITY", "Unable to process the contained instructions"},

	StatusInternalError:           {"INTERNAL_SERVER_ERROR", "Oops something went wrong"},
	StatusBadGatewayError:         {"STATUS_BAD_GATEWAY_ERROR", "Oops something went wrong"},
	StatusServiceUnavailableError: {"STATUS_SERVICE_UNAVAILABLE_ERROR", "Service Unavailable"},
	StatusGatewayTimeoutError:     {"STATUS_GATEWAY_TIMEOUT_ERROR", "Gateway Timeout"},
}

func StatusCode(code int) string {
	return statusMap[code][0]
}

func StatusText(code int) string {
	return statusMap[code][1]
}
