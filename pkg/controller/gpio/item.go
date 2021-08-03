package gpio

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/AliRostami1/baagh/pkg/controller/gpio/mode"
	"github.com/AliRostami1/baagh/pkg/controller/gpio/state"
	"github.com/stianeikeland/go-rpio/v4"
)

type Optional struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ItemData struct {
	Pin   rpio.Pin `json:"pin"`
	State string   `json:"state"`
	Mode  string   `json:"mode"`
	Optional
}

type Item struct {
	*GPIO

	key  string
	data *ItemData
	mu   *sync.RWMutex
}

func DefaultItem(g *GPIO, pin uint8, mode mode.Mode, state state.State) *Item {
	return &Item{
		GPIO: g,
		key:  makeKey(pin),
		data: &ItemData{
			Pin:      rpio.Pin(pin),
			State:    state.String(),
			Mode:     mode.String(),
			Optional: Optional{},
		},
		mu: &sync.RWMutex{},
	}
}

func (i *Item) SetMeta(opt Optional) {
	i.data.Optional = opt
}

func (i *Item) Marshal() (string, error) {
	data, err := json.Marshal(i.data)
	return string(data), err
}

func (i *Item) Commit() error {
	data, err := i.Marshal()
	if err != nil {
		return err
	}
	err = i.db.Set(i.key, data)
	return err
}

func (i *Item) State() state.State {
	i.mu.Lock()
	defer i.mu.Unlock()
	state, _ := state.FromString(i.data.State)
	return state
}

func (i *Item) Pin() uint8 {
	return uint8(i.data.Pin)
}

func (i *Item) Mode() mode.Mode {
	mode, _ := mode.FromString(i.data.Mode)
	return mode
}

func (i *Item) Key() string {
	return i.key
}

func (i *Item) cleanup() {
	if i.data.Mode == mode.OutputStr {
		i.data.Pin.Low()
	}
}

func (i *Item) submitItem() error {
	if _, exists := i.GPIO.registeredItems[i.key]; exists {
		return &MultipleController{pin: uint8(i.data.Pin)}
	}
	i.GPIO.registeredItems[i.key] = i
	return nil
}

type CircularDependency struct {
	pin uint8
}

func (c CircularDependency) Error() string {
	return fmt.Sprintf("circular dependency: %o can't depend on itself", c.pin)
}

type MultipleController struct {
	pin uint8
}

func (m MultipleController) Error() string {
	return fmt.Sprintf("can't add 2 controllers for the same pin: %o", m.pin)
}
