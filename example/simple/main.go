/*  main.go
*
* @Author:             Nanang Suryadi
* @Date:               April 04, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 04/04/20 02:53
 */

package main

import (
	"fmt"
	"os"

	"github.com/suryakencana007/mimir/example/simple/app"
)

func main() {
	if cleanup, err := app.InitializeApplication(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during dependency injection: %v\n", err)
		if cleanup != nil {
			cleanup()
		}
		os.Exit(1)
	}
}
