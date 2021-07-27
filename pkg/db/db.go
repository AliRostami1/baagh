package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CallbackFn func(val interface{})

type Db struct {
	db     *redis.Client
	events map[string][]CallbackFn
	ctx    context.Context
}

// TODO: should add transaction to it
func (d *Db) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	status := d.db.Set(d.ctx, key, value, expiration)
	if fns, ok := d.events[key]; ok {
		for _, fn := range fns {
			fn(value)
		}
	}
	return status
}

func (d *Db) Get(key string) *redis.StringCmd {
	return d.db.Get(d.ctx, key)
}

func (d *Db) OnSet(key string, fn ...CallbackFn) {
	d.events[key] = append(d.events[key], fn...)
}

func New() *Db {
	return &Db{
		db: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		events: make(map[string][]CallbackFn),
		ctx:    context.Background(),
	}
}
