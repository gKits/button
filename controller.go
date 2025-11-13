package button

import (
	"machine"
	"time"
)

// ActionType defines different types of click actions that can be used when registering a button action.
type ActionType uint8

const (
	// SingleClick actions are triggered once the assigned button is pressed once for a short time.
	SingleClick ActionType = iota
	// LongPress actions are triggered once the assigned button is pressed and held down for the configured long press
	// duration of the controller [default: 1s].
	LongPress
)

type Controller struct {
	buttons      []button
	actions      map[ActionType]map[int]func()
	debounceDur  time.Duration
	longPressDur time.Duration
	pinMode      machine.PinMode
}

// NewController returns a new button controller with the optionally provided options applied.
//
// By default a new controller uses:
//   - a debounce duration of 50ms
//   - a long press duration of 1s
//   - the machine.PinInputPullup pin mode for the pins that are registered as buttons
//
// You can adjust the values by using the provided ControllerOptions
func NewController(opts ...ControllerOption) *Controller {
	ctrl := &Controller{
		actions:      make(map[ActionType]map[int]func()),
		debounceDur:  50 * time.Millisecond,
		longPressDur: 1 * time.Second,
		pinMode:      machine.PinInputPullup,
	}

	for _, opt := range opts {
		opt(ctrl)
	}

	return ctrl
}

// Register configures the given pin as a button and adds a new action that is triggered once the
// actionTypes requirements are met.
func (c *Controller) Register(pin machine.Pin, actionType ActionType, action func()) {
	pin.Configure(machine.PinConfig{Mode: c.pinMode})
	id := len(c.buttons)
	c.buttons = append(c.buttons, button{
		pin:        pin,
		recordedAt: time.Now(),
		pullup:     c.pinMode == machine.PinInputPullup,
	})
	if _, ok := c.actions[actionType]; !ok {
		c.actions[actionType] = make(map[int]func())
	}
	c.actions[actionType][id] = action
}

// Update proceeds the state of the controller and its buttons. Make sure to run this method in
// the main loop of your tinygo programm.
func (c *Controller) Update() {
	now := time.Now()
	for id, btn := range c.buttons {
		isPressed := btn.readState()

		if isPressed != btn.isPressed && now.Sub(btn.recordedAt) > c.debounceDur {
			btn.isPressed = isPressed
			btn.recordedAt = now

			if isPressed {
				btn.pressedAt = now
			} else if action, ok := c.actions[SingleClick][id]; ok && now.Sub(btn.pressedAt) < c.longPressDur {
				go action()
			}
		}

		if btn.isPressed && now.Sub(btn.pressedAt) >= c.longPressDur {
			if !btn.longPressed {
				btn.longPressed = true
				if action, ok := c.actions[LongPress][id]; ok {
					go action()
				}
			}
		} else if !btn.longPressed {
			btn.longPressed = false
		}

		c.buttons[id] = btn
	}
}

type button struct {
	pin         machine.Pin
	pullup      bool
	isPressed   bool
	longPressed bool
	recordedAt  time.Time
	pressedAt   time.Time
}

func (btn *button) readState() bool {
	return btn.pullup == btn.pin.Get()
}

// ControllerOption is a type used to configure specific details of a Controller.
type ControllerOption func(ctrl *Controller)

// WithDebounceDuration sets the debounce duration of the controller.
func WithDebounceDuration(debounce time.Duration) ControllerOption {
	return func(ctrl *Controller) {
		ctrl.debounceDur = debounce
	}
}

// WithLongPressDuration set the long press duration of the controller.
func WithLongPressDuration(longPress time.Duration) ControllerOption {
	return func(ctrl *Controller) {
		ctrl.longPressDur = longPress
	}
}

// WithPinMode sets the pin mode the controller uses when registering a new button.
func WithPinMode(pinMode machine.PinMode) ControllerOption {
	return func(ctrl *Controller) {
		ctrl.pinMode = pinMode
	}
}
