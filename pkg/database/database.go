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

func New(ctx context.Context, opt *Options) (*DB, error) {
	db, err := badger.Open(badger.DefaultOptions(opt.Path).WithLogger(opt.Logger))
	if err != nil {
		return nil, err
	}
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
			value = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return string(value), err
}

func (d *DB) Set(key string, value string) error {
	isDifferent := true
	err := d.Update(func(txn *badger.Txn) error {
		item, getErr := txn.Get([]byte(key))
		keyExists := true
		if getErr != nil {
			if getErr == badger.ErrKeyNotFound {
				keyExists = false
			} else {
				return getErr
			}
		}
		if keyExists {
			item.Value(func(val []byte) error {
				if string(val) == value {
					isDifferent = false
				}
				return nil
			})
		}
		var setErr error
		if isDifferent {
			setErr = txn.Set([]byte(key), []byte(value))
		}
		return setErr
	})
	if err != nil {
		return err
	}
	if isDifferent {
		go d.executeEvents(key, value)
	}
	return nil
}

func (d *DB) On(key string, fn ...EventListener) {
	d.eventListeners[key] = append(d.eventListeners[key], fn...)
}

func (d *DB) executeEvents(key string, value string) {
	for _, fn := range d.eventListeners[key] {
		fn(key, value)
	}
}
