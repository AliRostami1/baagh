package main

import (
	"log"
	"time"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/gpio"
	"github.com/warthog618/gpiod"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}
	defer app.Cleanup()

	pirSensor, err := app.Gpio.Input(9, gpio.InputOption{
		Meta: gpio.Meta{
			Name:        "pir_sensor",
			Description: "pir sensor for detecting movement",
		},
		Pull: gpiod.WithPullDown,
	})

	_, err = app.Gpio.OutputAlarm(10, pirSensor.Key(), 7*time.Second, gpio.OutputOption{
		Meta: gpio.Meta{
			Name:        "led_light",
			Description: "blue led light that turns on every time the pir_sensor senses something",
		},
	})
	if err != nil {
		app.Log.Fatalf("there was a problem while initiating led light: %v", err)
	}

	<-app.Ctx.Done()
}
