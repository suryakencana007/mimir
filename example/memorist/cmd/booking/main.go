/*  main.go
*
* @Date:               April 17, 2020
* @Last Modified time: 17/04/20 06:25
 */

package main

import (
	"fmt"
	"os"

	"github.com/suryakencana007/mimir/example/memorist/cmd/booking/app"
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
