package security

import "fmt"

type OptionError struct {
	Field string
	Value interface{}
}

type Option interface {
	applyOption(*Options) error
}

type Control struct {
	sensors   []int
	actuators []int
}

type Options struct {
	control map[string]Control
}

type ConfigOption struct {
	chip string
	Control
}

func (c ConfigOption) applyOption(o *Options) error {
	if c.actuators == nil {
		return OptionError{Field: "Actuators", Value: c.actuators}
	}
	if c.sensors == nil {
		return OptionError{Field: "Sensors", Value: c.sensors}

	}
	if c.chip == "" {
		return OptionError{Field: "Chip", Value: c.chip}

	}
	if con, ok := o.control[c.chip]; ok {
		o.control[c.chip] = Control{
			sensors:   append(con.sensors, c.sensors...),
			actuators: append(con.actuators, c.actuators...),
		}
		return nil
	}

	o.control[c.chip] = c.Control
	return nil
}

func WithConfig(chip string, sensors []int, actuators []int) ConfigOption {
	return ConfigOption{
		chip: chip,
		Control: Control{
			sensors:   sensors,
			actuators: actuators,
		},
	}
}

func (o OptionError) Error() string {
	return fmt.Sprintf("field %s can not be: %v", o.Field, o.Value)
}
