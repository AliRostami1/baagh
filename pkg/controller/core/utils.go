package core

import (
	"fmt"
)

type State int

const (
	Inactive State = iota
	Active
)

func (s State) String() string {
	switch s {
	case Active:
		return "active"
	case Inactive:
		return "inactive"
	default:
		return ""
	}
}

func (s State) Check() error {
	if s == Active || s == Inactive {
		return nil
	}
	return InvalidStateError{}
}

type InvalidStateError struct{}

func (u InvalidStateError) Error() string {
	return fmt.Sprintf("state can't be any value other than %s and %s", Active, Inactive)
}

type Mode int

const (
	_ Mode = iota
	Input
	Output
)

func (m Mode) String() string {
	switch m {
	case Input:
		return "input"
	case Output:
		return "output"
	default:
		return ""
	}
}

func (m Mode) Check() error {
	if m == Input || m == Output {
		return nil
	}
	return InvalidModeError{}
}

type InvalidModeError struct{}

func (u InvalidModeError) Error() string {
	return fmt.Sprintf("mode can't be any value other than %s and %s", Output, Input)
}

type Pull int

const (
	PullUnknown Pull = iota
	PullDisabled
	PullDown
	PullUp
)

func (p Pull) String() string {
	switch p {
	case PullUnknown:
		return "unknown"
	case PullDisabled:
		return "disabled"
	case PullDown:
		return "down"
	case PullUp:
		return "up"
	default:
		return "ERROR"
	}
}

func (p Pull) Check() error {
	if p == PullDisabled || p == PullDown || p == PullUp {
		return nil
	}
	return InvalidPullError{}
}

type InvalidPullError struct{}

func (i InvalidPullError) Error() string {
	return fmt.Sprintf("mode can't be any value other than %s and %s", Output, Input)
}
