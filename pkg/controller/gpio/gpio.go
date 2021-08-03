package gpio

import (
	"context"
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"

	"github.com/AliRostami1/baagh/pkg/database"
)

type GPIO struct {
	db              *database.DB
	ctx             context.Context
	registeredItems map[string]*Item
}

type EventHandler func(item *Item)

func New(ctx context.Context, db *database.DB) (*GPIO, error) {
	if err := rpio.Open(); err != nil {
		return nil, fmt.Errorf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	gpio := &GPIO{
		db:              db,
		ctx:             ctx,
		registeredItems: make(map[string]*Item),
	}
	go gpio.cleanup()
	return gpio, nil
}

func (g *GPIO) GetItem(pin uint8) (*Item, error) {
	item, exists := g.registeredItems[makeKey(pin)]
	if !exists {
		return nil, NoController{pin: pin, key: makeKey(pin)}
	}
	return item, nil
}

func (g *GPIO) cleanup() {
	defer rpio.Close()
	<-g.ctx.Done()
	for _, item := range g.registeredItems {
		item.cleanup()
	}
}

type NoController struct {
	pin uint8
	key string
}

func (n NoController) Error() string {
	return fmt.Sprintf("there is no controller with pin number: %o, and corresponding key: %s", n.pin, n.key)
}

func makeKey(pin uint8) string {
	return fmt.Sprintf("pin_%o", pin)
}
