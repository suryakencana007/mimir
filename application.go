/*  application.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 07:06
 */

package mimir

import (
	"context"
	"fmt"
)

type AppRunner func(context.Context) error
type ApplicationFunc func(context.Context) (AppRunner, func(), error)

func Application(interrupt InterruptChannel, app ApplicationFunc) (func(), error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)

	runner, cleanup, err := app(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		errChan <- runner(ctx)
	}()

	select {
	case <-interrupt:
		cancel()
		return cleanup, fmt.Errorf("interrupt received, shutting down")
	case err := <-errChan:
		return cleanup, err
	}
}
