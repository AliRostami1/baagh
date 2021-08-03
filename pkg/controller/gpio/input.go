package gpio

import (
	"sync"
	"time"

	"github.com/AliRostami1/baagh/pkg/sensor"
	"github.com/stianeikeland/go-rpio/v4"
)

type ErrorFunction func(state State, err error)

type InputController struct {
	*Item
	sensor *sensor.Sensor

	errFn ErrorFunction
}

func (i *InputController) set(state State) error {
	err := state.Check()
	if err != nil {
		return err
	}
	i.mu.Lock()
	defer i.mu.Unlock()

	i.Item.data.State = state
	err = i.Item.Commit()
	return err
}

func (i *InputController) OnError(errFn ErrorFunction) {
	i.errFn = errFn
}

func (g *GPIO) Input(pin uint8, pull sensor.Pull) *InputController {
	input := InputController{
		Item:   &Item{GPIO: g, data: defaultItemData(pin, Input), mu: &sync.RWMutex{}},
		sensor: sensor.New(g.ctx, pin, &sensor.Options{Pull: pull, TickDuration: 500 * time.Millisecond}),
		errFn: func(state State, err error) {
		},
	}
	input.submitItem()

	input.sensor.OnChange(func(state rpio.State) {
		if err := input.set(State(state)); err != nil {
			input.errFn(State(state), err)
		}
	})

	input.sensor.Start()
	return &input
}
