/*  app.go
*
* @Date:               April 22, 2020
* @Last Modified time: 22/04/20 01:48
 */

package app

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/memorist/config"
	"github.com/suryakencana007/mimir/example/memorist/internal/http"
)

func InitializeApplication() (func(), error) {
	return mimir.Application(
		mimir.InterruptChannelFunc(),
		func(ctx context.Context) (mimir.AppRunner, func(), error) {
			logger := mimir.With(mimir.Field("Headless", "listen and serve"))
			cfg := &config.Config{}
			if err := mimir.Config(mimir.ConfigOpts{
				Config:   cfg,
				Filename: "app.config",
				Paths:    []string{"./config"},
			}, func(v *viper.Viper) error {
				return v.BindEnv("db.dsn_main")
			}); err != nil {
				return nil, nil, err
			}

			tracer, cleanupTrace, err := mimir.Tracer(cfg.App.Name, cfg.App.Version, logger)
			if err != nil {
				logger.Errorf("tracing is disconnected: %s", err)
				return nil, nil, err
			}

			opentracing.SetGlobalTracer(tracer)

			// http router
			router := http.Middleware(http.Options{
				Config: cfg,
				Tracer: tracer,
			})

			serve, cleanup := mimir.ListenAndServe(mimir.ServeOpts{
				Logger: logger,
				Port:   mimir.WebPort(cfg.App.Port),
				Router: router,
			})

			runner := func(ctx context.Context) error {
				if err := serve(ctx); err != nil {
					return err
				}
				return nil
			}

			return runner, func() {
				cleanupTrace()
				cleanup(ctx)
			}, nil
		})
}
