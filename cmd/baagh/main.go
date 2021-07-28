package main

import (
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"

	"github.com/AliRostami1/baagh/internal/application"
)

func main() {
	app := application.New()

	// initialize rpio package and allocate memory
	if err := rpio.Open(); err != nil {
		app.Log.Fatalf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(10)
	pin.Output()

	for {
		select {
		case _, ok := <-app.Ctx.Done():
			if !ok {
				pin.Low()
				app.Log.Info(app.Ctx.Err())
				os.Exit(1)
			}
		default: // pass
		}
		pin.Toggle()
		time.Sleep(time.Second)
	}

}
