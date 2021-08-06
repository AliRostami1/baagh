package gpio

import (
	"context"
	"fmt"
	"sync"

	"github.com/warthog618/gpiod"

	"github.com/AliRostami1/baagh/pkg/database"
)

type Gpio struct {
	chip *gpiod.Chip
	db   *database.DB
	ctx  context.Context

	*ItemRegistry
}

type GpioOption struct {
	ChipName string
	Ctx      context.Context
	DB       *database.DB
	Consumer string
}

func (g GpioOption) validate() error {
	if g.Ctx == nil {
		return GpioOptionError{missingField: "Ctx"}
	}
	if g.DB == nil {
		return GpioOptionError{missingField: "DB"}
	}
	if g.Consumer == "" {
		return GpioOptionError{missingField: "Consumer"}
	}
	return nil
}

func New(options GpioOption) (*Gpio, error) {
	err := options.validate()
	if err != nil {
		return nil, err
	}

	chip, err := gpiod.NewChip(gpiod.Chips()[0], gpiod.WithConsumer(options.Consumer))
	if err != nil {
		return nil, GpioInitializationError{err: err}
	}

	gpio := &Gpio{
		chip: chip,
		db:   options.DB,
		ctx:  options.Ctx,
		ItemRegistry: &ItemRegistry{
			registry: make(map[string]*Object),
			RWMutex:  &sync.RWMutex{},
		},
	}
	return gpio, nil
}

func (g *Gpio) Cleanup() error {
	defer g.chip.Close()
	var err error
	g.ItemRegistry.forEach(func(obj *Object) {
		err = obj.Line.Close()
	})
	if err != nil {
		return GpioCleanupError{err: err}
	}
	return nil
}

func makeKey(pin int) string {
	return fmt.Sprintf("pin_%o", pin)
}

type GpioOptionError struct {
	missingField string
}

func (g GpioOptionError) Error() string {
	return fmt.Sprintf("field %s can't be empty", g.missingField)
}

type GpioInitializationError struct {
	err error
}

func (g GpioInitializationError) Error() string {
	return fmt.Sprintf("can't initialize gpio controller: %v", g.err)
}

type GpioCleanupError struct {
	err error
}

func (g GpioCleanupError) Error() string {
	return fmt.Sprintf("a problem occurred while cleaning up the gpio controller: %v", g.err)
}
