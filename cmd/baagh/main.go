package main

import (
	"log"

	"github.com/warthog618/gpiod"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/core"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}
	defer app.Cleanup()

	chipName := gpiod.Chips()[0]

	_, err = core.RegisterChip(app.Ctx, core.WithName(chipName), core.WithConsumer("baagh"))

	led, err := core.RegisterItem(chipName, 10, core.AsOutput(), core.WithState(core.Inactive))

	pir, err := core.RegisterItem(chipName, 9, core.AsInput(core.PullDown))
	if err != nil {
		app.Log.Fatalf("there was a problem while initiating pir sensor: %v", err)
	}
	pir.AddEventListener(func(event *core.ItemEvent) {
		led.SetState(event.Item.State())
	})

	if err != nil {
		app.Log.Fatalf("there was a problem while initiating led light: %v", err)
	}

	<-app.Ctx.Done()
}
