package chip

import (
	"context"
	"fmt"

	"github.com/AliRostami1/baagh/pkg/log"
	"github.com/warthog618/gpiod"
)

type Chip struct {
	chip *gpiod.Chip
	ctx  context.Context

	logger   log.Logger
	backup   Setter
	name     string
	consumer string

	objects *objectRegistry
	events  *eventRegistry
}

// Setter is an interface to submit data KV-Data store
type Setter interface {
	Set(key string, value []byte) error
}

// ChipOption are Options for Chip!
type ChipOption struct {
	// Name of Chip to use, i.e. "gpiochip0"
	Name string

	// Consumer of the Chip, i.e. "myapp"
	Consumer string

	// Logger to use, ignore to disable logging
	Logger log.Logger

	// Store to save every registered lines state on disc,
	// ignore to disable data persistance.
	Store Setter
}

func (co ChipOption) validateAndApply(c *Chip) error {
	if co.Logger == nil {
		c.logger = &dummyLogger{}
	} else {
		c.logger = co.Logger
	}

	if co.Name == "" {
		return ChipOptionError{Field: "Name", Err: "empty"}
	}
	c.name = co.Name

	if co.Consumer == "" {
		return ChipOptionError{Field: "Consumer", Err: "empty"}
	}
	c.consumer = co.Consumer

	if co.Store == nil {
		c.logger.Warnf("data won't get saved: chip is initialized with no persistent store")
	}

	return nil
}

func New(ctx context.Context, options ChipOption) (c *Chip, err error) {
	err = options.validateAndApply(c)
	if err != nil {
		return
	}

	// init the chip
	chip, err := gpiod.NewChip(c.name, gpiod.WithConsumer(c.consumer))
	if err != nil {
		return nil, ChipInitError{Err: err}
	}
	// add the chip to Chip object
	c.chip = chip

	// add context
	c.ctx = ctx

	return
}

func (c *Chip) Cleanup() error {
	defer c.chip.Close()
	var err error
	c.registry.forEach(func(i int, obj *Object) {
		err = obj.Line.Close()
	})
	if err != nil {
		return ChipCleanupError{Err: err}
	}
	return nil
}

// ChipOptionError
type ChipOptionError struct {
	// the field that has originated the error
	Field string
	// can be "empty" or "bad" field
	Err string
}

func (g ChipOptionError) Error() string {
	return fmt.Sprintf("field %s is %s", g.Field, g.Err)
}

type ChipInitError struct {
	Err error
}

func (g ChipInitError) Error() string {
	return fmt.Sprintf("can't initialize chip controller: %v", g.Err)
}

type ChipCleanupError struct {
	Err error
}

func (g ChipCleanupError) Error() string {
	return fmt.Sprintf("a problem occurred while cleaning up the chip controller: %v", g.Err)
}
