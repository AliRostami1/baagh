package gpio

import (
	"time"

	"github.com/AliRostami1/baagh/pkg/controller/gpio/mode"
	"github.com/AliRostami1/baagh/pkg/controller/gpio/sensor"
	"github.com/AliRostami1/baagh/pkg/controller/gpio/state"
	"github.com/stianeikeland/go-rpio/v4"
)

type ErrorFunction func(state state.State, err error)

type InputController struct {
	*Item
	sensor *sensor.Sensor

	errFn ErrorFunction
}

func (i *InputController) set(state state.State) error {
	err := state.Check()
	if err != nil {
		return err
	}

	i.SetState(state)
	err = i.Item.Commit()
	return err
}

func (i *InputController) OnError(errFn ErrorFunction) {
	i.errFn = errFn
}

func (g *GPIO) Input(pin uint8, pull sensor.Pull) *InputController {
	input := InputController{
		Item:   DefaultItem(g, pin, mode.Input, state.Low),
		sensor: sensor.New(g.ctx, pin, &sensor.Options{Pull: pull, TickDuration: 500 * time.Millisecond}),
		errFn: func(state state.State, err error) {
		},
	}
	g.addItem(pin, input.Item)

	input.sensor.OnChange(func(s rpio.State) {
		if err := input.set(state.State(s)); err != nil {
			input.errFn(state.State(s), err)
		}
	})

	input.sensor.Start()

	return &input
}
