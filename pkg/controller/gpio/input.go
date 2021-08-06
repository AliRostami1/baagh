package gpio

import (
	"fmt"
	"log"
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
				Meta:  Meta{},
			},
			key: makeKey(pin),
			mu:  &sync.RWMutex{},
		},
	}

	handler := func(evt gpiod.LineEvent) {
		log.Println("movement detected")
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
