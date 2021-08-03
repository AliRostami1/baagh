package main

import (
	"log"
	"time"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/gpio/sensor"
	"github.com/AliRostami1/baagh/pkg/controller/gpio/state"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}
	defer app.Cleanup()

	pirSensor := app.Gpio.Input(9, sensor.PullDown)
	pirSensor.OnError(func(state state.State, err error) {
		app.Log.Fatalf("there was a problem while initiating pir sensor: %v", err)
	})

	_, _, err = app.Gpio.OutputAlarm(10, pirSensor.Key(), 7*time.Second)
	if err != nil {
		app.Log.Fatalf("there was a problem while initiating led light: %v", err)
	}

	<-app.Ctx.Done()
}
