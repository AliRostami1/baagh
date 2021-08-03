package gpio

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

type State rpio.State

const (
	Low State = iota
	High
)

const (
	HighStr = "on"
	LowStr  = "off"
)

type InvalidStateError struct{}

func (u InvalidStateError) Error() string {
	return fmt.Sprintf("gpio state can't be any value other than %s=%o and %s=%o", HighStr, High, LowStr, Low)
}

func (s State) String() string {
	if s == High {
		return HighStr
	} else if s == Low {
		return LowStr
	}
	panic(InvalidStateError{})
}

func (s State) Set(val string) error {
	if val == HighStr {
		s = High
		return nil
	} else if val == LowStr {
		s = Low
		return nil
	}
	return InvalidStateError{}

}

func (s State) Check() error {
	if s == High || s == Low {
		return nil
	}
	return InvalidStateError{}
}
