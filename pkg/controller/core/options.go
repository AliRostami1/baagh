package core

import (
	"github.com/AliRostami1/baagh/pkg/errprim"
)

type ItemOption interface {
	applyItemOption(*ItemOptions) error
}

type ItemOptions struct {
	mode  Mode
	pull  Pull
	state State
}

type InputOption struct {
	pull Pull
}

func (i InputOption) applyItemOption(item *ItemOptions) (err error) {
	if err = i.pull.Check(); err != nil {
		return errprim.OptionError{Field: "pull", Value: i.pull}
	}
	item.mode = ModeInput
	item.pull = i.pull
	return nil
}

func AsInput(pull Pull) InputOption {
	return InputOption{
		pull: pull,
	}
}

type OutputOption struct {
	state State
}

func (o OutputOption) applyItemOption(item *ItemOptions) (err error) {
	if err = o.state.Check(); err != nil {
		return
	}
	item.state = o.state
	item.mode = ModeOutput
	return
}

func AsOutput(state State) OutputOption {
	return OutputOption{state: state}
}

type StateOption State

func (s StateOption) applyItemOption(item *ItemOptions) (err error) {
	state := State(s)
	if err = state.Check(); err != nil {
		return errprim.OptionError{Field: "state", Value: s}
	}
	item.state = State(s)
	return
}

func WithState(state State) StateOption {
	return StateOption(state)
}
