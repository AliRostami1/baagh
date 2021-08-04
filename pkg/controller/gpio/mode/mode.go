package mode

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

type Mode rpio.Mode

const (
	Input Mode = iota
	Output
	// we do not yet support these:
	// Clock
	// Pwm
	// Spi
)

const (
	InputStr  = "input"
	OutputStr = "output"
)

type InvalidModeError struct{}

func (u InvalidModeError) Error() string {
	return fmt.Sprintf("gpio mode can't be any value other than %s=%o and %s=%o", InputStr, Input, OutputStr, Output)
}

func FromString(mode string) (Mode, error) {
	if mode == InputStr {
		return Input, nil
	} else if mode == OutputStr {
		return Output, nil
	}
	return 255, InvalidModeError{}
}

func (m Mode) String() string {
	if m == Input {
		return InputStr
	} else if m == Output {
		return OutputStr
	}
	panic(InvalidModeError{})
}

func (m Mode) Set(val string) error {
	if val == OutputStr {
		m = Output
		return nil
	} else if val == InputStr {
		m = Input
		return nil
	}
	return InvalidModeError{}
}

func (m Mode) Check() error {
	if m == Output || m == Input {
		return nil
	}
	return InvalidModeError{}
}
