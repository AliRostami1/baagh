package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func new() <-chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}

func Handle(fn func(os.Signal)) {
	sig := new()
	for s := range sig {
		fn(s)
	}
}

func ShutdownHandler(fn func(string)) {
	go Handle(func(s os.Signal) {
		fn(fmt.Sprintf("terminating: %v signal received", s))
	})
}
