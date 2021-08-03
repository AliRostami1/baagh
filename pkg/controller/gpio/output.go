package gpio

import (
	"time"

	"github.com/AliRostami1/baagh/pkg/controller/gpio/mode"
	"github.com/AliRostami1/baagh/pkg/controller/gpio/state"
	"github.com/AliRostami1/baagh/pkg/debounce"
	"github.com/stianeikeland/go-rpio/v4"
)

type OutputController struct {
	*Item
}

func (o *OutputController) Set(state state.State) error {
	err := state.Check()
	if err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.Item.data.State = state.String()
	err = o.Item.Commit()
	if err != nil {
		return err
	}

	o.data.Pin.Write(rpio.State(state))

	return nil
}

func (o *OutputController) SetHigh() {
	o.Set(state.High)
}

func (o *OutputController) SetLow() {
	o.Set(state.Low)
}

func (o *OutputController) On(key string, fns ...EventHandler) error {
	if key == o.key {
		return CircularDependency{pin: o.Pin()}
	}

	for _, fn := range fns {
		o.db.On(key, func(key, value string) {
			if item, ok := o.registeredItems[key]; ok {
				fn(item)
			}
		})
	}

	return nil
}

func (o *OutputController) OnItem(item *Item, fns ...EventHandler) error {
	return o.On(item.key, fns...)
}

func (o *OutputController) OnPin(pin uint8, fns ...EventHandler) error {
	item, err := o.GPIO.GetItem(pin)
	if err != nil {
		return err
	}
	return o.OnItem(item, fns...)
}

func (g *GPIO) Output(pin uint8) (*OutputController, error) {
	output := OutputController{
		Item: DefaultItem(g, pin, mode.Output, state.Low),
	}
	output.data.Pin.Output()
	err := output.submitItem()
	if err != nil {
		return nil, err
	}

	err = output.Set(state.Low)
	if err != nil {
		return nil, err
	}
	return &output, err
}

func (g *GPIO) OutputSync(pin uint8, key string) (*OutputController, error) {
	output, err := g.Output(pin)
	if err != nil {
		return nil, err
	}
	err = output.On(key, func(item *Item) {
		output.Set(item.State())
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (g *GPIO) OutputRSync(pin uint8, key string) (*OutputController, error) {
	output, err := g.Output(pin)
	if err != nil {
		return nil, err
	}
	err = output.On(key, func(item *Item) {
		if item.State() == state.High {
			output.Set(state.Low)
		} else {
			output.Set(state.High)
		}
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (g *GPIO) OutputAlarm(pin uint8, key string, delay time.Duration) (*OutputController, func(), error) {
	output, err := g.Output(pin)
	if err != nil {
		return nil, nil, err
	}
	fn, cancel := debounce.Debounce(delay, func() {
		output.Set(state.Low)
	})

	go func() {
		<-g.ctx.Done()
		cancel()
	}()

	err = output.On(key, func(item *Item) {
		if item.State() == state.High {
			output.Set(state.High)
			fn()
		}
	})
	if err != nil {
		return nil, nil, err
	}
	return output, cancel, nil
}
