package gpio

import (
	"time"

	"github.com/AliRostami1/baagh/pkg/debounce"
	"github.com/stianeikeland/go-rpio/v4"
)

func (g *GPIO) Output(pin int, listen *EventListener) error {
	p := rpio.Pin(pin)
	p.Output()
	g.addOutputPins(p)
	if err := g.Set(pin, false); err != nil {
		return err
	}

	err := g.on(pin, listen)
	return err
}

func (g *GPIO) OutputSync(pin int, key string) error {
	return g.Output(pin, &EventListener{
		Key: key,
		Fn: func(p int, v bool) {
			g.Set(p, v)
		},
	})

}

func (g *GPIO) OutputRSync(pin int, key string) error {
	return g.Output(pin, &EventListener{
		Key: key,
		Fn: func(p int, v bool) {
			if v {
				g.Set(p, false)
			} else {
				g.Set(p, true)
			}
		},
	})
}

func (g *GPIO) OutputAlarm(pin int, key string, delay time.Duration) (cancel func(), err error) {
	fn, cancel := debounce.Debounce(delay, func() {
		g.Set(pin, false)
	})

	go func() {
		<-g.ctx.Done()
		cancel()
	}()

	return cancel, g.Output(pin, &EventListener{
		Key: key,
		Fn: func(p int, v bool) {
			if v {
				g.Set(p, true)
				fn()
			}
		},
	})
}
