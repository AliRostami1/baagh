package gpio

import (
	"time"

	"github.com/AliRostami1/baagh/pkg/sensor"
	"github.com/stianeikeland/go-rpio/v4"
)

type InputController struct {
	*Item
	*sensor.Sensor

	OnErr func(err error, state State)
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

func (g *GPIO) Input(pin uint8, pull sensor.Pull) *InputController {
	input := InputController{
		Item:   &Item{GPIO: g, data: defaultItemData(pin, Input)},
		Sensor: sensor.New(g.ctx, pin, &sensor.Options{Pull: pull, TickDuration: 500 * time.Millisecond}),
		OnErr: func(err error, state State) {
		},
	}
	input.submitItem()

	input.Sensor.OnChange(func(state rpio.State) {
		if err := input.set(State(state)); err != nil {
			input.OnErr(err, State(state))
		}
	})
	return &input
}
