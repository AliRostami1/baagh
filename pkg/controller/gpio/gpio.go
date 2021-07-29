package gpio

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stianeikeland/go-rpio/v4"

	"github.com/AliRostami1/baagh/pkg/db"
	"github.com/AliRostami1/baagh/pkg/sensor"
)

type GPIO struct {
	db         *db.Db
	ctx        context.Context
	outputPins []rpio.Pin
}

func New(ctx context.Context, db *db.Db) (*GPIO, error) {
	if err := rpio.Open(); err != nil {
		return nil, fmt.Errorf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}

	gpio := &GPIO{
		db:  db,
		ctx: ctx,
	}
	go gpio.cleanup()

	return gpio, nil
}

type EventHandler func(pin int, val bool)

type EventListeners struct {
	Key string
	Fn  EventHandler
}

func (g *GPIO) RegisterOutputPin(pin int, listen *EventListeners) (err error) {
	p := rpio.Pin(pin)
	p.Output()
	g.addOutputPins(p)
	g.Set(pin, false)

	err = g.On(pin, append([]*EventListeners{}, listen))

	return err
}

func (g *GPIO) RegisterInputPin(pin int) {
	g.Set(pin, false)
	go sensor.SensorFn(pin, func(s bool) {
		g.Set(pin, s)
	})
}

func (g *GPIO) On(pin int, listen []*EventListeners) error {
	for _, ev := range listen {
		if ev.Key == fmt.Sprint(pin) {
			return fmt.Errorf("circular dependency: pin%[1]o can't depend on pin%[1]o", pin)
		}

		g.db.On(ev.Key, func(key string, val *redis.StringCmd) error {
			v, err := val.Bool()
			if err != nil {
				return fmt.Errorf("can't sync ")
			}
			ev.Fn(pin, v)
			return nil
		})
	}
	return nil
}

func (g *GPIO) Sync(pin int, val bool) {
	g.Set(pin, val)
}

func (g *GPIO) ReverseSync(pin int, val bool) {
	if val {
		g.Set(pin, false)
	} else {
		g.Set(pin, true)
	}
}

func (g *GPIO) Set(pin int, val bool) {
	p := rpio.Pin(pin)
	g.db.Set(fmt.Sprint(pin), val, 0)
	if val {
		p.Write(rpio.High)
	} else {
		p.Write(rpio.Low)
	}
}

func (g *GPIO) cleanup() {
	defer rpio.Close()
	<-g.ctx.Done()
	for _, pin := range g.outputPins {
		pin.Low()
	}
}

func (g *GPIO) addOutputPins(pin rpio.Pin) error {
	for _, p := range g.outputPins {
		if p == pin {
			return fmt.Errorf("can't add 2 controllers for the same pin")
		}
	}
	g.outputPins = append(g.outputPins, pin)
	return nil
}
