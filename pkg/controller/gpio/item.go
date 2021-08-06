package gpio

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/warthog618/gpiod"
)

type ObjectTrx struct {
	old      Object
	new      Object
	discard  bool
	newState int
}

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
	Gpio *Gpio
	*gpiod.Line
	data *ObjectData
	key  string
	mu   *sync.RWMutex
}

func (o *Object) set(fn func(trx *ObjectTrx) error) error {
	trx := &ObjectTrx{
		old:      *o,
		new:      *o,
		discard:  false,
		newState: -1,
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	err := fn(trx)
	if err != nil {
		return err
	}
	if !trx.discard {
		if trx.newState != -1 {
			err := o.SetValue(trx.newState)
			if err != nil {
				return err
			}
		}
		o = &trx.new
		err = o.commitToDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ObjectTrx) Discard() {
	o.discard = true
}

func (o *ObjectTrx) SetMeta(opt Meta) {
	o.new.data.Meta = opt
}

func (o *ObjectTrx) SetState(state State) error {
	if o.new.data.Info.Config.Direction == gpiod.LineDirectionOutput {
		o.newState = int(state)
	}
	o.new.data.State = state
	return nil
}

func (o *ObjectTrx) SetInfo(info gpiod.LineInfo) {
	o.new.data.Info = info
}

func (o *ObjectTrx) Inactive() {
	o.SetState(ACTIVE)
}

func (o *ObjectTrx) Active() {
	o.SetState(INACTIVE)
}

func (i *Object) Data() ObjectData {
	i.mu.Lock()
	defer i.mu.Unlock()
	return *i.data
}

func (i *Object) Marshal() (string, error) {
	info := i.data
	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(jsonInfo), nil
}

func (i *Object) commitToDB() error {
	data, err := i.Marshal()
	if err != nil {
		return err
	}
	err = i.Gpio.db.Set(i.key, data)
	if err != nil {
		return err
	}
	return nil
}

func (i *Object) Key() string {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.key
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
