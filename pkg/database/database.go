package database

import (
	"context"

	"github.com/dgraph-io/badger/v3"
)

type EventListener func(key string, value string)

type DB struct {
	*badger.DB
	eventListeners map[string][]EventListener
}

type Options struct {
	Path   string
	Logger badger.Logger
}

// TODO add our custom logger to this
func New(ctx context.Context, opt *Options) (*DB, error) {
	db, err := badger.Open(badger.DefaultOptions(opt.Path).WithLogger(opt.Logger))
	if err != nil {
		return nil, err
	}
	go func() {
		defer db.Close()
		<-ctx.Done()
	}()
	return &DB{
		DB:             db,
		eventListeners: make(map[string][]EventListener),
	}, nil
}

func (d *DB) Get(key string) (string, error) {
	var value []byte
	err := d.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			copy(value, val)
			return nil
		})
		return err
	})
	return string(value), err
}

func (d *DB) Set(key string, value string) error {
	isDifferent := false
	err := d.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			if string(val) != value {
				isDifferent = true
			}
			return nil
		})
		if isDifferent {
			err = txn.Set([]byte(key), []byte(value))
		}
		return err
	})
	if isDifferent {
		d.executeEvents(key, value)
	}
	return err
}

func (d *DB) On(key string, fn ...EventListener) {
	d.eventListeners[key] = append(d.eventListeners[key], fn...)
}

func (d *DB) executeEvents(key string, value string) {
	for _, fn := range d.eventListeners[key] {
		fn(key, value)
	}
}
