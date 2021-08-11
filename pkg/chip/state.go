package chip

import (
	"fmt"
)

type State int

const (
	INACTIVE State = iota
	ACTIVE
)

func (s State) String() string {
	if s == ACTIVE {
		return "active"
	} else if s == INACTIVE {
		return "inactive"
	}
	panic(InvalidStateError{})
}

func (s State) Check() error {
	if s == ACTIVE || s == INACTIVE {
		return nil
	}
	return InvalidStateError{}
}

type InvalidStateError struct{}

func (u InvalidStateError) Error() string {
	return fmt.Sprintf("gpio state can't be any value other than %s and %s", ACTIVE, INACTIVE)
}
