package gpio

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/stianeikeland/go-rpio/v4"
)

type ItemData struct {
	Pin   rpio.Pin `json:"pin"`
	State State    `json:"state"`
	Mode  Mode     `json:"mode"`

	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	*GPIO
	data *ItemData
	mu   *sync.RWMutex
}

func (i *ItemData) withName(name string) *ItemData {
	i.Name = name
	return i
}

func (i *ItemData) withDescription(description string) *ItemData {
	i.Description = description
	return i
}

func defaultItemData(pin uint8, mode Mode) *ItemData {
	return &ItemData{
		Pin:  rpio.Pin(pin),
		Mode: mode,
		Key:  makeKey(pin),
	}
}

func (i *Item) Commit() error {
	data, err := json.Marshal(i.data)
	if err != nil {
		return err
	}
	err = i.db.Set(i.data.Key, string(data))
	return err
}

func (i *Item) State() State {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.data.State
}

func (i *Item) Data() *ItemData {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.data
}

func (i *Item) Pin() uint8 {
	return uint8(i.data.Pin)
}

func (i *Item) Mode() string {
	return fmt.Sprint(i.data.Mode)
}

func (i *Item) Key() string {
	return i.data.Key
}

func (i *Item) cleanup() {
	if i.data.Mode == Output {
		i.data.Pin.Low()
	}
}

func (i *Item) submitItem() error {
	if _, exists := i.GPIO.registeredItems[i.data.Key]; exists {
		return &MultipleController{pin: uint8(i.data.Pin)}
	}
	i.GPIO.registeredItems[i.data.Key] = i
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
