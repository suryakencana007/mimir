/*  middleware.go
*
* @Date:               April 22, 2020
* @Last Modified time: 22/04/20 11:41
 */

package http

import (
	"github.com/opentracing/opentracing-go"
	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/memorist/config"
	"github.com/suryakencana007/mimir/ruuto"
)

type Options struct {
	*config.Config
	Tracer opentracing.Tracer
}

func Middleware(opts Options) ruuto.Router {
	router := ruuto.NewChiRouter()
	router.Use(mimir.TracerServer(opts.Tracer, opts.App.Name))
	router.Use(mimir.Logger())

	return Router(Handlers{
		Config: opts.Config,
		Router: router,
		Health: Health(),
	})
}
