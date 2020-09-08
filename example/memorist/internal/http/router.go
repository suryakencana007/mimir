/*  router.go
*
* @Date:               April 22, 2020
* @Last Modified time: 22/04/20 01:36
 */

package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/memorist/config"
	"github.com/suryakencana007/mimir/ruuto"
)

type Handlers struct {
	*config.Config
	Router ruuto.Router
	Health http.HandlerFunc
}

func Router(handlers Handlers) ruuto.Router {
	router := handlers.Router

	api := router.Group(fmt.Sprintf("%s%s", handlers.Main.Prefix, handlers.API.Prefix))
	api.GET("/healthz", handlers.Health)

	return router
}

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "health.handler")
		defer span.Finish()

		log := mimir.For(ctx)
		resp := mimir.Response(r)

		if mimir.MediaType(r.Header.Get("Content-Type")) == mimir.ApplicationJSON {

			var response map[string]interface{}
			b := http.MaxBytesReader(w, r.Body, 1048576)
			body := json.NewDecoder(b)
			body.DisallowUnknownFields()

			if err := body.Decode(&response); err != nil {
				log.Errorf("Request body contains badly-formed JSON Error: %v", err)
				resp.APIStatusBadRequest(w, r, err).WriteJSON()
				return
			}

			log.Info("request health handler", log.Field("body", response))
		}

		log.Info("hahahaha", mimir.Field("name", "health"), mimir.Field("body", map[string]interface{}{
			"Status": "is okay",
			"Method": r.Method,
			"Body":   "awesome awesome some",
		}))

		resp.Body(map[string]interface{}{
			"Status": "Health is Okay",
		})

		resp.APIStatusSuccess(w, r).WriteJSON()
	}
}
