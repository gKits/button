# Button

A library for tinygo to handle button inputs

## Usage

1. Import the library in your `tinygo` project.

```bash
go get github.com/gkits/button
```

2. Register pins as buttons with their assgined actions.

```go
ctrl := button.NewController()

ctrl.Register(machine.GP15, button.SingleClick, func() { println("single click") })
ctrl.Register(machine.GP15, button.LongPress, func() { println("long press") })
```

3. Update the controller in your programms main loop.

```go
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
```
