package sensor

import "github.com/stianeikeland/go-rpio/v4"

type SensorCallback = func(current rpio.State)

func SensorFn(pin int, fn SensorCallback) {
	ch := make(chan rpio.State)
	go Sensor(pin, ch)
	for state := range ch {
		fn(state)
	}
}

func Sensor(pin int, ch chan<- rpio.State) {
	p := rpio.Pin(pin)
	p.Input()
	prevState := rpio.State(0)
	for {
		state := p.Read()
		if state != prevState {
			prevState = state
			ch <- state
		}
	}
}
