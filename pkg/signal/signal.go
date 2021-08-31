package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func New() <-chan os.Signal {
	// TODO use signal.NotifyContext instead
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}

func Handle(fn func(os.Signal)) {
	sig := New()
	for s := range sig {
		fn(s)
	}
}

func ShutdownHandler(fn func(string)) {
	go Handle(func(s os.Signal) {
		fn(fmt.Sprintf("terminating: %v signal received", s))
	})
}

func Gracefull() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}
