package general

import "fmt"

const (
	Sync  = "sync"
	RSync = "rsync"
	Alarm = "alarm"

	AllIn = "all-in"
	OneIn = "one-in"
)

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
	// kind can be "alarm", "sync" and "rsync"
	kind string
	// strategy is only relevant if kind is "sync" or "rsync"
	// it can either be "all-in" which means it will turn on
	// only when all inputs are active, and "one-in" which will
	// turn on when any of the inputs are active
	strategy string
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

type KindOption struct {
	kind     string
	strategy string
}

func (k KindOption) applyOption(o *Options) error {
	if k.kind == "" {
		k.kind = Sync
		k.strategy = AllIn
	} else if (k.kind == Sync || k.kind == RSync) && (k.strategy != AllIn && k.strategy != OneIn) {
		return OptionError{
			Field: "Strategy",
			Value: k,
		}
	} else if k.kind != Alarm && k.kind != Sync && k.kind != RSync {
		return OptionError{
			Field: "Kind",
			Value: k,
		}
	}
	o.kind = k.kind
	o.strategy = k.strategy
	return nil
}

func WithKind(kind string, strategy string) KindOption {
	return KindOption{
		kind:     kind,
		strategy: strategy,
	}
}

func AsAlarm() KindOption {
	return KindOption{
		kind:     Alarm,
		strategy: "",
	}
}

func AsSync(strategy string) KindOption {
	return KindOption{
		kind:     Sync,
		strategy: strategy,
	}
}

func AsRSync(strategy string) KindOption {
	return KindOption{
		kind:     RSync,
		strategy: strategy,
	}
}

type StrategyOption string

func (s StrategyOption) applyOption(o *Options) error {
	if s == "" {
		s = OneIn
	} else if s != AllIn && s != OneIn {
		return OptionError{
			Field: "Strategy",
			Value: s,
		}
	}
	o.strategy = string(s)
	return nil
}

func WithStrategy(strategy string) StrategyOption {
	return StrategyOption(strategy)
}

func AsAllIn() StrategyOption {
	return AllIn
}

func AsOneIn() StrategyOption {
	return OneIn
}

func (o OptionError) Error() string {
	return fmt.Sprintf("field %s can not be: %v", o.Field, o.Value)
}
