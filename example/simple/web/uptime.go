package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/suryakencana007/mimir"
)

func DefaultUptimeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := mimir.Response(r)
		timeDelta := time.Now().Sub(time.Now()).Seconds()
		resp.Body(fmt.Sprintf("%.2f", timeDelta))
		resp.APIStatusSuccess(w, r).WriteJSON()
	}
}
