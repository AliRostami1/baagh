package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CallbackFn func(key string, val *redis.StringCmd) error

type Db struct {
	db     *redis.Client
	events map[string][]CallbackFn
	ctx    context.Context
}

// TODO: should add transaction to it
func (d *Db) Set(key string, value interface{}, expiration time.Duration) (*redis.StatusCmd, error) {
	status := d.db.Set(d.ctx, key, value, expiration)
	if err := status.Err(); err != nil {
		return nil, err
	}
	if fns, ok := d.events[key]; ok {
		res, err := d.Get(key)
		if err != nil {
			return nil, err
		}
		for _, fn := range fns {
			fn(key, res)
		}
	}
	return status, nil
}

func (d *Db) Get(key string) (*redis.StringCmd, error) {
	status := d.db.Get(d.ctx, key)
	if err := status.Err(); err != nil {
		return nil, err
	}
	return status, nil
}

func (d *Db) On(key string, fn ...CallbackFn) {
	d.events[key] = append(d.events[key], fn...)
}

func (d *Db) Connected() error {
	if _, err := d.db.Ping(d.ctx).Result(); err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context, url string) (*Db, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db := &Db{
		db:     redis.NewClient(opt),
		events: make(map[string][]CallbackFn),
		ctx:    ctx,
	}

	if err := db.Connected(); err != nil {
		return nil, err
	}

	return db, nil
}
