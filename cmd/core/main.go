package main

import (
	"bufio"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/warthog618/gpiod"
)

func main() {
	go func() {
		http.ListenAndServe(":1234", nil)
	}()

	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}

	chipName := gpiod.Chips()[0]

	core.SetLogger(app.Log)
	if err != nil {
		app.Log.Fatal(err)
	}

	led10, err := core.RequestItem(chipName, 10, core.AsOutput(core.StateActive))
	if err != nil {
		log.Fatalf("there was a problem registering led on pin 10, %v", err)
	}
	defer led10.Close()

	ledWatcher, err := core.NewWatcher(chipName, 10, core.AsOutput(core.StateActive))
	if err != nil {
		log.Fatalf("there was a problem registering led-watcher on pin 10, %v", err)
	}
	defer ledWatcher.Close()

	rfWatcher, err := core.NewInputWatcher(chipName, 9)
	if err != nil {
		log.Fatalf("there was a problem registering rf-watcher on pin 9, %v", err)
	}
	defer rfWatcher.Close()

	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			switch input.Text() {
			case "on":
				led10.SetState(core.StateActive)
			case "off":
				led10.SetState(core.StateInactive)
			}
		}
	}()
	go func() {
		for c := range ledWatcher.Watch() {
			log.Printf("from led-watcher: %v", c)
		}
	}()
	go func() {
		for range rfWatcher.Watch() {
			// log.Printf("from rf-watcher: %v", c)
		}
	}()

	go func() {
		for {
			app.Log.Sync()
		}
	}()

	<-app.Ctx.Done()
}
