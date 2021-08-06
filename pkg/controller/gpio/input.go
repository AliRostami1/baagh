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
	input := InputObject{
		Object: &Object{
			Gpio: g,
			Line: nil,
			data: &ObjectData{
				Info:  gpiod.LineInfo{},
				State: INACTIVE,
				Meta:  option.Meta,
			},
			key: makeKey(pin),
			mu:  &sync.RWMutex{},
		},
	}
	input.mu.Lock()
	defer input.mu.Unlock()

	handler := func(evt gpiod.LineEvent) {
		input.Object.set(func() error {
			if evt.Type == gpiod.LineEventRisingEdge {
				input.setState(ACTIVE)
				return nil
			}
			if evt.Type == gpiod.LineEventFallingEdge {
				input.setState(INACTIVE)
				return nil
			}
			return fmt.Errorf("THIS IS WROOONGGGG")
		})
	}

	inputLine, err := g.chip.RequestLine(pin, gpiod.AsInput, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)
	if err != nil {
		return nil, fmt.Errorf("there was a problem with input controller: %v", err)
	}

	inputInfo, err := inputLine.Info()
	if err != nil {
		return nil, err
	}

	input.Line = inputLine
	input.data.Info = inputInfo

	g.addItem(pin, input.Object)

	return &input, nil
}
