/*  app.go
*
* @Date:               April 22, 2020
* @Last Modified time: 22/04/20 01:53
 */

package app

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/memorist/config"
	rpc "google.golang.org/grpc"
)

func InitializeApplication() (func(), error) {
	return mimir.Application(
		mimir.InterruptChannelFunc(),
		func(ctx context.Context) (mimir.AppRunner, func(), error) {
			logger := mimir.With(mimir.Field("Headless", "listen and serve rpc"))
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

			server, cleanup := mimir.RemoteCallProc(mimir.GRPCOpts{
				Logger: logger,
				Port:   mimir.GRPCPort(cfg.GRPC.Port),
				Opts:   nil,
			})

			runner := func(ctx context.Context) error {
				serverSpan := opentracing.StartSpan("Pool Testing PG")
				defer serverSpan.Finish()

				ctxSpan := opentracing.ContextWithSpan(ctx, serverSpan)

				if err := server(ctxSpan, func(s *rpc.Server) error {
					// pb handler
					return nil
				}); err != nil {
					return err
				}
				return nil
			}

			return runner, func() {
				cleanup()
			}, nil
		})
}
