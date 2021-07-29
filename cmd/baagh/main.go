package main

import (
	"log"
	"os"
	"time"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/gpio"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}

	// initialize rpio package and allocate memory
	gpioController, err := gpio.New(app.Ctx, app.Db)
	if err != nil {
		app.Log.Errorf("there was a problem initiating the gpio controller: %v", err)
	}

	gpioController.RegisterOutputPin(10, &gpio.EventListeners{
		Key: "test",
		Fn:  gpioController.Sync,
	})

	app.Db.Set("test", false, 0)
	for {
		select {
		case _, ok := <-app.Ctx.Done():
			if !ok {
				os.Exit(1)
			}
		default: // pass
		}
		res, err := app.Db.Get("test").Bool()
		if err != nil {
			app.Log.Fatal("ddddaaymn")
		}
		app.Db.Set("test", !res, 0)
		time.Sleep(time.Second)
	}

}
