package core

import (
	"encoding/json"
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

func (s *State) UnmarshalJSON(data []byte) error {
	var str string

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	return s.Set(str)
}

func (s State) MarshalJSON() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	str := s.String()
	data, err = json.Marshal(str)

	return
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
	ModeUnknown Mode = iota
	ModeInput
	ModeOutput
)

func (m *Mode) Set(value string) error {
	switch value {
	case "unknown":
		*m = ModeUnknown
	case "input":
		*m = ModeInput
	case "output":
		*m = ModeOutput
	default:
		return InvalidModeError{}
	}
	return nil
}

func (m *Mode) UnmarshalJSON(data []byte) error {
	var str string

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	return m.Set(str)
}

func (m Mode) MarshalJSON() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	str := m.String()
	data, err = json.Marshal(str)

	return
}

func (m Mode) String() string {
	switch m {
	case ModeUnknown:
		return "unknown"
	case ModeInput:
		return "input"
	case ModeOutput:
		return "output"
	default:
		panic(InvalidModeError{})
	}
}

func (m Mode) Type() string {
	return "mode"
}

func (m Mode) Check() error {
	if m >= ModeUnknown && m <= ModeOutput {
		return nil
	}
	return InvalidModeError{}
}

type InvalidModeError struct{}

func (u InvalidModeError) Error() string {
	return fmt.Sprintf("mode can't be any value other than %s, %s and %s", ModeUnknown, ModeOutput, ModeInput)
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

func (p *Pull) UnmarshalJSON(data []byte) error {
	var str string

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	return p.Set(str)
}

func (p Pull) MarshalJSON() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	str := p.String()
	data, err = json.Marshal(str)

	return
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
	if p >= PullUnknown && p <= PullUp {
		return nil
	}
	return InvalidPullError{}
}

type InvalidPullError struct{}

func (i InvalidPullError) Error() string {
	return fmt.Sprintf("pull can't be any value other than %s, %s, %s and %s", PullDisabled, PullUnknown, PullDown, PullUp)
}
