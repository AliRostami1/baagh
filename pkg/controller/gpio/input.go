package gpio

import (
	"fmt"
	"sync"

	"github.com/warthog618/gpiod"
)

type Pull int

type InputOption struct {
	Meta
	Pull gpiod.BiasOption
}

type InputObject struct {
	*Object
}

func (g *Gpio) Input(pin int, option InputOption) (*InputObject, error) {
	inputPin, err := g.chip.RequestLine(pin, gpiod.AsInput, option.Pull)
	if err != nil {
		return nil, fmt.Errorf("there was a problem with input controller: %v", err)
	}

	inputInfo, err := inputPin.Info()
	if err != nil {
		return nil, err
	}

	input := InputObject{
		Object: &Object{
			Gpio: g,
			Line: inputPin,
			data: &ObjectData{
				Info:  inputInfo,
				State: INACTIVE,
				Meta:  Meta{},
			},
			key: makeKey(pin),
			mu:  &sync.RWMutex{},
		},
	}

	handler := func(evt gpiod.LineEvent) {
		input.Object.set(func(trx *ObjectTrx) error {
			if evt.Type == gpiod.LineEventRisingEdge {
				trx.SetState(ACTIVE)
				return nil
			}
			if evt.Type == gpiod.LineEventFallingEdge {
				trx.SetState(INACTIVE)
				return nil
			}
			return fmt.Errorf("THIS IS WROOONGGGG")
		})
	}

	input.Reconfigure(gpiod.WithBothEdges, gpiod.WithEventHandler(handler))

	g.addItem(pin, input.Object)

	return &input, nil
}
