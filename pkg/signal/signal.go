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

func Handle(fn func(string)) {
	sig := new()
	for s := range sig {
		fn(fmt.Sprintf("%v signal recieved, terminating", s))
	}
}
