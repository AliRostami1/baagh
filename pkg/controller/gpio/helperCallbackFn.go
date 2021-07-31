package gpio

import (
	"time"

	"github.com/AliRostami1/baagh/pkg/debounce"
)

func (g *GPIO) OutputSync(pin int, key string) {
	g.Output(pin, &EventListeners{
		Key: key,
		Fn: func(p int, v bool) {
			g.Set(p, v)
		},
	})
}

func (g *GPIO) OutputRSync(pin int, key string) {
	g.Output(pin, &EventListeners{
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

func (g *GPIO) OutputAlarm(pin int, key string, delay time.Duration) {
	clear := debounce.Debounce(delay, func() {
		g.Set(pin, false)
	})

	g.Output(pin, &EventListeners{
		Key: key,
		Fn: func(p int, v bool) {
			if v {
				g.Set(p, true)
				clear()
			}
		},
	})
}
