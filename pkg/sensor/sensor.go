package sensor

import (
	"context"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type SensorCallback = func(current bool)

type Pull uint8

const (
	PullOff Pull = iota
	PullDown
	PullUp
)

func SensorFunc(ctx context.Context, pin int, pull Pull, fn SensorCallback) {
	ch := make(chan bool)
	go Sensor(ctx, pin, pull, ch)
	for {
		if state, ok := <-ch; ok {
			fn(state)
		} else {
			break
		}
	}
}

func Sensor(ctx context.Context, pin int, pull Pull, ch chan<- bool) {
	defer close(ch)
	p := rpio.Pin(pin)
	p.Input()
	p.Pull(rpio.Pull(pull))

	prevState := false
infinite:
	for {

		time.Sleep(500 * time.Millisecond)
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
