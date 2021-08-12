package core

import (
	"fmt"

	"github.com/warthog618/gpiod"
)

type OptionError struct {
	Field string
	Value interface{}
}

type ChipOption interface {
	applyChipOption(*ChipOptions) error
}

type ChipOptions struct {
	name     string
	consumer string
}

type NameOption string

func (n NameOption) applyChipOption(c *ChipOptions) error {
	var chipExistsOnDevice bool
	for _, deviceChipName := range gpiod.Chips() {
		if string(n) == deviceChipName {
			chipExistsOnDevice = true
		}
	}
	if chipExistsOnDevice {
		c.name = string(n)
		return nil
	}
	return OptionError{Field: "name", Value: n}
}

func WithName(name string) NameOption {
	return NameOption(name)
}

type ConsumerOption string

func (n ConsumerOption) applyChipOption(c *ChipOptions) error {
	if string(n) == "" {
		n = "baagh"
	}
	c.consumer = string(n)
	return nil
}

func WithConsumer(consumer string) ConsumerOption {
	return ConsumerOption(consumer)
}

type ItemOption interface {
	applyItemOption(*ItemOptions) error
}

type ItemOptions struct {
	// mode can be "input" and "output"
	io struct {
		mode Mode
		pull Pull
	}
	state State
}

func (o OptionError) Error() string {
	return fmt.Sprintf("field %s can not be: %v", o.Field, o.Value)
}

type InputOption struct {
	mode Mode
	pull Pull
}

func (i InputOption) applyItemOption(item *ItemOptions) (err error) {
	if err = i.mode.Check(); err != nil {
		return OptionError{Field: "mode", Value: i.mode}
	}
	if err = i.pull.Check(); err != nil {
		return OptionError{Field: "pull", Value: i.pull}
	}
	item.io.mode = i.mode
	item.io.pull = i.pull
	return nil
}

func AsInput(pull Pull) InputOption {
	return InputOption{
		mode: Input,
		pull: pull,
	}
}

type OutputOption Mode

func (o OutputOption) applyItemOption(item *ItemOptions) (err error) {
	if err = Mode(o).Check(); err != nil {
		return OptionError{Field: "mode", Value: o}
	}
	item.io.mode = Mode(o)
	return
}

func AsOutput() OutputOption {
	return OutputOption(Output)
}

type StateOption State

func (s StateOption) applyItemOption(item *ItemOptions) (err error) {
	if err = State(s).Check(); err != nil {
		return OptionError{Field: "state", Value: s}
	}
	item.state = State(s)
	return
}

func WithState(state State) StateOption {
	return StateOption(state)
}
