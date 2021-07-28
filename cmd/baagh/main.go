package main

import (
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/signal"
)

func main() {

	app := application.New()

	// React to process signals
	exitSig := signal.HandleSignals()

	// initialize rpio package and allocate memory
	if err := rpio.Open(); err != nil {
		app.Log.Fatalf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(10)
	pin.Output()

	for {
		select {
		case sig := <-exitSig:
			pin.Low()
			app.Log.Info(sig)
			os.Exit(1)
		default:
		}
		pin.Toggle()
		time.Sleep(time.Second)
	}
}
