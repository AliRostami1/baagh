package core

import (
	"fmt"
)

type State int

const (
	StateInactive State = iota
	StateActive
)

func (s *State) Set(value string) error {
	switch value {
	case "active":
		*s = StateActive
	case "inactive":
		*s = StateInactive
	default:
		return InvalidStateError{}
	}
	return nil
}

func (s State) String() string {
	switch s {
	case StateActive:
		return "active"
	case StateInactive:
		return "inactive"
	default:
		panic(InvalidStateError{})
	}
}

func (s State) Type() string {
	return "state"
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

func (m *Mode) Set(value string) error {
	switch value {
	case "input":
		*m = ModeInput
	case "output":
		*m = ModeOutput
	default:
		return InvalidModeError{}
	}
	return nil
}

func (m Mode) String() string {
	switch m {
	case ModeInput:
		return "input"
	case ModeOutput:
		return "output"
	default:
		panic("")
	}
}

func (m Mode) Type() string {
	return "mode"
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

func (p *Pull) Set(value string) error {
	switch value {
	case "unknown":
		*p = PullUnknown
	case "disabled":
		*p = PullDisabled
	case "down":
		*p = PullDown
	case "up":
		*p = PullUp
	default:
		return InvalidPullError{}
	}
	return nil
}

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
		panic(InvalidPullError{})
	}
}

func (p Pull) Type() string {
	return "pull"
}

func (p Pull) Check() error {
	if p == PullDisabled || p == PullDown || p == PullUp {
		return nil
	}
	return InvalidPullError{}
}

type InvalidPullError struct{}

func (i InvalidPullError) Error() string {
	return fmt.Sprintf("pull can't be any value other than %s, %s, %s and %s", PullDisabled, PullUnknown, PullDown, PullUp)
}
