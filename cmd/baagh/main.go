package main

import (
	"bufio"
	"log"
	"os"

	"github.com/warthog618/gpiod"

	"github.com/AliRostami1/baagh/internal/application"
	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/AliRostami1/baagh/pkg/controller/general"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("there was a problem initiating the application: %v", err)
	}
	defer app.Cleanup()

	chipName := gpiod.Chips()[0]

	core.SetLogger(app.Log)
	_, err = core.RegisterChip(app.Ctx, core.WithName(chipName), core.WithConsumer("baagh"))
	if err != nil {
		app.Log.Fatal(err)
	}
	defer core.Cleanup()

	alarm, err := general.Register("security-system", general.AsRSync(general.OneIn), general.WithConfig(chipName, []int{9}, []int{10}))
	if err != nil {
		return
	}

	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			if input.Text() == "turn off" {
				alarm.TurnOff()
			}
		}
	}()

	<-app.Ctx.Done()

}
