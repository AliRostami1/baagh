package sensor

import (
	"context"

	"github.com/stianeikeland/go-rpio/v4"
)

type SensorCallback = func(current bool)

func SensorFn(ctx context.Context, pin int, fn SensorCallback) {
	ch := make(chan bool)
	go Sensor(ctx, pin, ch)
	for {
		if state, ok := <-ch; ok {
			fn(state)
		} else {
			break
		}
	}
}

func Sensor(ctx context.Context, pin int, ch chan<- bool) {
	defer close(ch)
	p := rpio.Pin(pin)
	p.Input()
	p.PullDown()

	prevState := false
infinite:
	for {
		select {
		case _, ok := <-ctx.Done():
			if !ok {
				break infinite
			}
		default:
			state := rToB(p.Read())
			if state != prevState {
				prevState = state
				ch <- state
			}
		}
	}
}

func rToB(s rpio.State) bool {
	return s == rpio.High
}
