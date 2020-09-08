/*  interupt.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 14:22
 */

package mimir

import (
	"os"
	"os/signal"
	"syscall"
)

type InterruptChannel <-chan os.Signal

func InterruptChannelFunc() InterruptChannel {
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	return interrupt
}
