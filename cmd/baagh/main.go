package main

import (
	"log"
	"os"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/gpio"
	"github.com/AliRostami1/baagh/pkg/sensor"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}

	// initialize rpio package and allocate memory
	gpioController, err := gpio.New(app.Ctx, app.Db)
	if err != nil {
		app.Log.Fatalf("there was a problem initiating the gpio controller: %v", err)
	}

	gpioController.RegisterOutputPin(10, &gpio.EventListeners{
		Key: "9",
		Fn:  gpioController.Sync,
	})

	gpioController.RegisterInputPin(9, sensor.PullDown)

	<-app.Ctx.Done()
	os.Exit(0)
}
