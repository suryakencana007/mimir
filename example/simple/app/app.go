/*  app.go
*
* @Author:             Nanang Suryadi
* @Date:               April 04, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 04/04/20 22:11
 */

package app

import (
	"context"

	"github.com/spf13/viper"
	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/simple/config"
	"github.com/suryakencana007/mimir/example/simple/db"
	"github.com/suryakencana007/mimir/example/simple/web"
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
				Paths:    []string{".", "./config"},
			}, func(v *viper.Viper) error {
				return v.BindEnv("db.dsn_main")
			}); err != nil {
				return nil, nil, err
			}

			_, cleanupDB, err := db.PostgresDBConn(logger, cfg)
			if err != nil {
				return nil, nil, err
			}

			serve, cleanup := mimir.ListenAndServe(mimir.ServeOpts{
				Logger: logger,
				Port:   mimir.WebPort(cfg.App.Port),
				Router: web.MainRouter(
					web.Handlers{
						Uptime: web.DefaultUptimeHandler(),
					}),
			})

			runner := func(ctx context.Context) error {
				if err := serve(ctx); err != nil {
					return err
				}
				return nil
			}

			return runner, func() {
				cleanupDB()
				cleanup(ctx)
			}, nil
		})
}
