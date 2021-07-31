package gpio

import "github.com/AliRostami1/baagh/pkg/sensor"

func (g *GPIO) Input(pin int, pull sensor.Pull) (ch chan error) {
	ch = make(chan error)
	go sensor.SensorFunc(g.ctx, pin, pull, func(s bool) {
		if err := g.Set(pin, s); err != nil {
			ch <- err
		}
	})
	return ch
}
