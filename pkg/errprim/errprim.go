package errprim

import "fmt"

type OptionError struct {
	Field string
	Value interface{}
}

func (o OptionError) Error() string {
	return fmt.Sprintf("field %s can not be: %v", o.Field, o.Value)
}
