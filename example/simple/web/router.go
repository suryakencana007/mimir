/*  router.go
*
* @Author:             Nanang Suryadi
* @Date:               April 05, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 05/04/20 04:34
 */

package web

import (
	"net/http"

	"github.com/suryakencana007/mimir/ruuto"
)

type Handlers struct {
	Uptime http.HandlerFunc
}

func MainRouter(handlers Handlers) ruuto.Router {
	router := ruuto.NewChiRouter()
	router.GET("/uptime", handlers.Uptime)
	return router
}
