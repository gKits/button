// This example sets up a button.Controller with its default configs.
//
// When the registered button is clicked it will print "single click" to the monitor.
// When the registered button is pressed for 1s it will print "long press".
package main

import (
	"machine"
	"time"

	"github.com/gkits/button"
)

func main() {
	ctrl := button.NewController()

	ctrl.Register(machine.GP15, button.SingleClick, func() { println("single click") })
	ctrl.Register(machine.GP15, button.LongPress, func() { println("long press") })

	// additional setup logic

	for {
		// additional runtime logic
		ctrl.Update()

		time.Sleep(10 * time.Millisecond)
	}
}
