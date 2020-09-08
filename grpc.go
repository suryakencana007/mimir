/*  grpc.go
*
* @Date:               April 17, 2020
* @Last Modified time: 17/04/20 06:56
 */

package mimir

import (
	"context"
	"fmt"
	"net"

	rpc "google.golang.org/grpc"
)

func GRPCkommen() string {
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
	GRPCPort     int
	GRPCCallback func(*rpc.Server) error
	GRPCRunFunc  func(context.Context, GRPCCallback) error
	GRPCOpts     struct {
		Logger Logging
		Port   GRPCPort
		Opts   []rpc.ServerOption
	}
)

func RemoteCallProc(opts GRPCOpts) (GRPCRunFunc, func()) {
	logger := opts.Logger
	s := rpc.NewServer(opts.Opts...) // GRpc Server
	cleanup := func() {
		logger.Info("I have to go...")
		logger.Info("Stopping server gracefully")
		if s != nil {
			logger.Infof("Stop server at :%d", opts.Port)
			s.GracefulStop()
		}
	}
	return func(ctx context.Context, callback GRPCCallback) error {
		errChan := make(chan error)
		go func() {
			n, err := net.Listen("tcp", fmt.Sprintf(":%v", opts.Port))
			if err != nil {
				logger.With(
					logger.Field("port", opts.Port),
					logger.Field("error", err.Error()),
				).Error("failed to listen:")
				errChan <- err
			}

			if err := callback(s); err != nil {
				errChan <- err
			}
			// Description µ micro service
			fmt.Println(
				fmt.Sprintf(
					GRPCkommen(),
					opts.Port,
				))
			logger.Info(fmt.Sprintf("Now serving at %v", s.GetServiceInfo()))
			errChan <- s.Serve(n)
		}()

		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			if s != nil {
				logger.Infof("Stop server at :%d", opts.Port)
				s.Stop()
			}
			return fmt.Errorf("server interrupted through context")
		}
	}, cleanup
}
