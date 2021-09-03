package grace

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Channel(sigs ...os.Signal) <-chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, sigs...)
	return sigCh
}

func Context(sigs ...os.Signal) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), sigs...)
}

func Shutdown(sigs ...os.Signal) (context.Context, context.CancelFunc) {
	sigs = append(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return Context(sigs...)
}
