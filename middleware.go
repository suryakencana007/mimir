/*  middleware.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 05:01
 */

package mimir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/opentracing/opentracing-go"
	opExt "github.com/opentracing/opentracing-go/ext"
)

func Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			With(
				Field("method", r.Method),
				Field("path", r.URL.Path),
			).Debug("Started Request")
			m := httpsnoop.CaptureMetrics(next, w, r)
			With(
				Field("code", m.Code),
				Field("duration", int(m.Duration/time.Millisecond)),
				Field("duration-fmt", m.Duration.String()),
				Field("method", r.Method),
				Field("host", r.Host),
				Field("request", r.RequestURI),
				Field("remote-addr", r.RemoteAddr),
				Field("referer", r.Referer()),
				Field("user-agent", r.UserAgent()),
			).Info("Completed handling request")
		})
	}
}

func Recovery() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					resp := Response(r)
					With(
						Field("method", r.Method),
						Field("path", r.URL.Path),
					).Error("Internal server error handled")
					switch internalErr := err.(type) {
					case error:
						resp.APIStatusInternalError(w, r, internalErr).WriteJSON()
					default:
						resp.APIStatusInternalError(w, r, fmt.Errorf("%v", err)).WriteJSON()
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func TracerServer(tracer opentracing.Tracer, operationName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			serverSpan := opentracing.SpanFromContext(ctx)
			if serverSpan == nil {
				// All we can do is create a new root span.
				serverSpan = tracer.StartSpan(operationName)
			} else {
				serverSpan.SetOperationName(operationName)
			}
			defer serverSpan.Finish()

			opExt.SpanKindRPCServer.Set(serverSpan)
			opExt.HTTPMethod.Set(serverSpan, r.Method)
			opExt.HTTPUrl.Set(serverSpan, r.URL.String())

			// There's nothing we can do with any errors here.
			if err := tracer.Inject(
				serverSpan.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header),
			); err != nil {
				For(ctx).Warnf("tracing err %s", err)
			}

			ctx = opentracing.ContextWithSpan(ctx, serverSpan)
			log := For(ctx)

			// check content length
			if r.ContentLength > 0 {
				// Request
				var buf []byte
				if r.Body != nil { // Read
					buf, _ = ioutil.ReadAll(r.Body)
				}

				r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				mediaBody := string(buf)

				// b := http.MaxBytesReader(w, r.Body, 1048576)
				resp := Response(r)
				// get content-type
				s := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))

				response := make(map[string]interface{}, 0)

				switch MediaType(s) {
				case TextPlain:
				case FormURLEncoded:
					if err := r.ParseForm(); err != nil {
						log.Errorf("Request body contains badly-formed form-urlencoded Error: %v", err)
						resp.APIStatusBadRequest(w, r, err).WriteJSON()
						return
					}

					log.Info("request payload", log.Field("content-type", s), log.Field("body", mediaBody))
				case MultipartForm:
				case ApplicationJSON:
					// b := http.MaxBytesReader(w, b, 1048576)
					body := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer([]byte(mediaBody))))
					body.DisallowUnknownFields()

					if err := body.Decode(&response); err != nil {
						log.Errorf("Request body contains badly-formed JSON Error: %v", err)
						resp.APIStatusBadRequest(w, r, err).WriteJSON()
						return
					}

					log.Info("request payload", log.Field("content-type", s), log.Field("body", response))
				default:
					log.Info("request payload", log.Field("content-type", s), log.Field("body", mediaBody))
				}
			}

			log.Infof("tracing form middleware endpoint %s", r.URL.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
