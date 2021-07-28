package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func new() <-chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}

func Handle(fn func()) {
	sig := new()
	for range sig {
		fn()
	}
}
