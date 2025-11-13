// This example sets up a button.Controller with the debounce duration set to 20ms, the long press duration set to 2s
// and the pin mode to machine.PinInputPulldown.
//
// When the registered button is clicked it will print "single click" to the monitor.
// When the registered button is pressed for 2s it will print "long press".
package main

import (
	"machine"
	"time"

	"github.com/gkits/button"
)

const buttonPin = machine.GP15

func main() {
	ctrl := button.NewController(
		button.WithDebounceDuration(20*time.Millisecond),
		button.WithLongPressDuration(2*time.Second),
		button.WithPinMode(machine.PinInputPulldown),
	)

	ctrl.Register(buttonPin, button.SingleClick, func() { println("single click") })
	ctrl.Register(buttonPin, button.LongPress, func() { println("long press") })

	for {
		ctrl.Update()
		time.Sleep(10 * time.Millisecond)
	}
}
