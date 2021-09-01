package core

import (
	"fmt"
)

type State int

const (
	StateInactive State = iota
	StateActive
)

func (s State) String() string {
	switch s {
	case StateActive:
		return "active"
	case StateInactive:
		return "inactive"
	default:
		return ""
	}
}

func (s State) Check() error {
	if s == StateActive || s == StateInactive {
		return nil
	}
	return InvalidStateError{}
}

type InvalidStateError struct{}

func (u InvalidStateError) Error() string {
	return fmt.Sprintf("state can't be any value other than %s and %s", StateActive, StateInactive)
}

type Mode int

const (
	_ Mode = iota
	ModeInput
	ModeOutput
)

func (m Mode) String() string {
	switch m {
	case ModeInput:
		return "input"
	case ModeOutput:
		return "output"
	default:
		return ""
	}
}

func (m Mode) Check() error {
	if m == ModeInput || m == ModeOutput {
		return nil
	}
	return InvalidModeError{}
}

type InvalidModeError struct{}

func (u InvalidModeError) Error() string {
	return fmt.Sprintf("mode can't be any value other than %s and %s", ModeOutput, ModeInput)
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
	return fmt.Sprintf("mode can't be any value other than %s and %s", ModeOutput, ModeInput)
}
