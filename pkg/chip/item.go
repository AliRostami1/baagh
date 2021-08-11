package chip

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/warthog618/gpiod"
)

type Meta struct {
	Name        string
	Description string
}

type ObjectData struct {
	Info  gpiod.LineInfo
	State State
	Meta  Meta
}

type Object struct {
	Chip *Chip
	*gpiod.Line
	data *ObjectData
	key  string
	mu   *sync.RWMutex
}

func (o *Object) set(fn func() error) error {
	err := fn()
	if err != nil {
		return err
	}

	data, err := o.Marshal()
	if err != nil {
		return err
	}

	err = o.commit(data)
	if err != nil {
		return err
	}

	return nil
}

func (o *Object) SetMeta(opt Meta) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.data.Meta = opt
}

func (o *Object) setState(state State) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.data.State = state
}

func (o *Object) setInfo(info gpiod.LineInfo) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.data.Info = info
}

func (o *Object) Inactive() {
	o.setState(ACTIVE)
}

func (o *Object) Active() {
	o.setState(INACTIVE)
}

func (i *Object) Data() ObjectData {
	i.mu.Lock()
	defer i.mu.Unlock()
	return *i.data
}

func (o *Object) Marshal() ([]byte, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	info := o.data
	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return jsonInfo, nil
}

// commit will update the data on Store if it is set
func (o *Object) commit(data []byte) error {
	if o.Chip.backup != nil {
		err := o.Chip.backup.Set(o.key, data)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// returns
func (o *Object) Key() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.key
}

type CircularDependency struct {
	key string
}

func (c CircularDependency) Error() string {
	return fmt.Sprintf("circular dependency: %s can't depend on itself", c.key)
}

type MultipleController struct {
	key string
}

func (m MultipleController) Error() string {
	return fmt.Sprintf("can't add 2 controllers for the same pin: %s", m.key)
}
