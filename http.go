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

func Velkommen() string {
	return `
========================================================================================
   _     _     _     _     _  
  / \   / \   / \   / \   / \ 
 ( s ) ( u ) ( k ) ( i ) ( ~ )
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
					Velkommen(),
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
	StatusSuccess               = http.StatusOK
	StatusErrorForm             = http.StatusBadRequest
	StatusErrorUnknown          = http.StatusBadGateway
	StatusInternalError         = http.StatusInternalServerError
	StatusUnauthorized          = http.StatusUnauthorized
	StatusCreated               = http.StatusCreated
	StatusAccepted              = http.StatusAccepted
	StatusForbidden             = http.StatusForbidden
	StatusInvalidAuthentication = http.StatusProxyAuthRequired
	StatusUnProcess             = http.StatusUnprocessableEntity
	StatusPaymentRequired       = http.StatusPaymentRequired
	StatusMethodNotAllowed      = http.StatusMethodNotAllowed
	StatusNotAcceptable         = http.StatusNotAcceptable
	StatusUnsupportedMediaType  = http.StatusUnsupportedMediaType
	StatusPermanentRedirect     = http.StatusPermanentRedirect
)

var statusMap = map[int][]string{
	StatusSuccess:               {"STATUS_OK", "Success"},
	StatusErrorForm:             {"STATUS_BAD_REQUEST", "Invalid data request"},
	StatusErrorUnknown:          {"STATUS_BAG_GATEWAY", "Oops something went wrong"},
	StatusInternalError:         {"INTERNAL_SERVER_ERROR", "Oops something went wrong"},
	StatusUnauthorized:          {"STATUS_UNAUTHORIZED", "Not authorized to access the service"},
	StatusCreated:               {"STATUS_CREATED", "Resource has been created"},
	StatusAccepted:              {"STATUS_ACCEPTED", "Resource has been accepted"},
	StatusForbidden:             {"STATUS_FORBIDDEN", "Forbidden access the resource "},
	StatusInvalidAuthentication: {"STATUS_INVALID_AUTHENTICATION", "The resource owner or authorization server denied the request"},
	StatusUnProcess:             {"STATUS_UNPROCESSABLE_ENTITY", "Unable to process the contained instructions"},
	StatusPaymentRequired:       {"STATUS_PAYMENT_REQUIRED", "Payment need to be done"},
	StatusMethodNotAllowed:      {"STATUS_METHOD_NOT_ALLOWED", "The method specified is not allowed"},
	StatusNotAcceptable:         {"STATUS_NOT_ACCEPTABLE", "Request cannot accepted"},
	StatusUnsupportedMediaType:  {"STATUS_UNSUPPORTED_MEDIA_TYPE", "Cannot understand request content"},
	StatusPermanentRedirect:     {"STATUS_PERMANENT_REDIRECT", "The resource has moved to a new location"},
}

func StatusCode(code int) string {
	return statusMap[code][0]
}

func StatusText(code int) string {
	return statusMap[code][1]
}
