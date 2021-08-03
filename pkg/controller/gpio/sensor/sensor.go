package sensor

import (
	"context"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type EventHandler = func(state rpio.State)

type Pull uint8

const (
	PullOff Pull = iota
	PullDown
	PullUp
)

type Options struct {
	Pull         Pull
	TickDuration time.Duration
}

type Internal struct {
	ctx context.Context
	pin uint8

	changeHandlers []EventHandler
	tickHandlers   []EventHandler
	changeChannel  []chan rpio.State
	tickChannel    []chan rpio.State
}

type Sensor struct {
	*Internal
	*Options
}

func (s *Sensor) OnChange(fn ...EventHandler) {
	s.changeHandlers = append(s.changeHandlers, fn...)
}

func (s *Sensor) OnTick(fn ...EventHandler) {
	s.tickHandlers = append(s.tickHandlers, fn...)
}

func (s *Sensor) Tick() <-chan rpio.State {
	ch := make(chan rpio.State)
	s.tickChannel = append(s.tickChannel, ch)
	return ch
}

func (s *Sensor) Change() <-chan rpio.State {
	ch := make(chan rpio.State)
	s.changeChannel = append(s.changeChannel, ch)
	return ch
}

func New(ctx context.Context, pin uint8, opt *Options) *Sensor {
	return &Sensor{
		Internal: &Internal{
			ctx: ctx,
			pin: pin,
		},
		Options: opt,
	}
}

func (s *Sensor) Start() {
	go func() {
		defer s.cleanup()
		p := rpio.Pin(s.pin)
		p.Input()
		p.Pull(rpio.Pull(s.Pull))

		prevState := rpio.Low
	infinite:
		for {
			time.Sleep(s.TickDuration)
			select {
			case _, ok := <-s.ctx.Done():
				if !ok {
					break infinite
				}
			default:
				state := p.Read()
				go s.runOnTick(state)
				if state != prevState {
					prevState = state
					go s.runOnChange(state)
				}
			}
		}
	}()
}

func (s *Sensor) runOnTick(state rpio.State) {
	for _, tf := range s.tickHandlers {
		tf(state)
	}
	for _, tc := range s.tickChannel {
		tc <- state
	}
}

func (s *Sensor) runOnChange(state rpio.State) {
	for _, cf := range s.changeHandlers {
		cf(state)
	}
	for _, cc := range s.changeChannel {
		cc <- state
	}
}

func (s *Sensor) cleanup() {
	for _, cc := range s.changeChannel {
		close(cc)
	}
	for _, tc := range s.tickChannel {
		close(tc)
	}
}
