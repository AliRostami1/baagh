package gpio

import (
	"context"
	"fmt"
	"sync"

	"github.com/stianeikeland/go-rpio/v4"

	"github.com/AliRostami1/baagh/pkg/database"
)

type GPIO struct {
	db  *database.DB
	ctx context.Context

	*ItemRegistery
}

type EventHandler func(item *Item)

func New(ctx context.Context, db *database.DB) (*GPIO, error) {
	if err := rpio.Open(); err != nil {
		return nil, fmt.Errorf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	gpio := &GPIO{
		db:  db,
		ctx: ctx,
		ItemRegistery: &ItemRegistery{
			registry: make(map[string]*Item),
			RWMutex:  &sync.RWMutex{},
		},
	}
	return gpio, nil
}

func (g *GPIO) Cleanup() {
	defer rpio.Close()
	g.ItemRegistery.forEach(func(item *Item) {
		item.cleanup()
	})
}

func makeKey(pin uint8) string {
	return fmt.Sprintf("pin_%o", pin)
}
