package sensor

import "github.com/stianeikeland/go-rpio/v4"

type SensorCallback = func(current bool)

func SensorFn(pin int, fn SensorCallback) {
	ch := make(chan bool)
	go Sensor(pin, ch)
	for state := range ch {
		fn(state)
	}
}

func Sensor(pin int, ch chan<- bool) {
	p := rpio.Pin(pin)
	p.Input()
	prevState := false
	for {
		state := rToB(p.Read())
		if state != prevState {
			prevState = state
			ch <- state
		}
	}
}

func rToB(s rpio.State) bool {
	return s == rpio.High
}
