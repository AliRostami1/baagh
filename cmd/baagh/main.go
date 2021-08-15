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

	core.SetLogger(app.Log)
	chip, err := core.RegisterChip(app.Ctx, core.WithName(gpiod.Chips()[0]), core.WithConsumer("baagh"))
	if err != nil {
		app.Log.Fatal(err)
	}
	defer core.Cleanup()

	led, err := chip.RegisterItem(10, core.AsOutput(), core.WithState(core.Inactive))
	if err != nil {
		app.Log.Fatal(err)
	}

	pir, err := chip.RegisterItem(9, core.AsInput(core.PullDown))
	if err != nil {
		app.Log.Fatalf("there was a problem while initiating pir sensor: %v", err)
	}
	pir.AddEventListener(func(event *core.ItemEvent) {
		err := led.SetState(event.Item.State())
		if err != nil {
			app.Log.Fatal(err)
		}
	})

	<-app.Ctx.Done()

}
