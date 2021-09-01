package main

import (
	"bufio"
	"log"
	"os"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/AliRostami1/baagh/pkg/logy"
	"github.com/AliRostami1/baagh/pkg/signal"
	"github.com/warthog618/gpiod"
	"go.uber.org/zap/zapcore"
)

func main() {
	defer core.Close()
	ctx, _ := signal.Gracefull()

	logger, err := logy.New(ctx, zapcore.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}
	core.SetLogger(logger)

	led, err := core.RequestItem(gpiod.Chips()[0], 10, core.AsOutput(core.StateInactive))
	if err != nil {
		log.Fatal(err)
	}
	// defer led.Close()
	log.Print("led registered")

	ledWatcher, err := core.NewWatcher(gpiod.Chips()[0], 10)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("led-watcher registered")
	// defer ledWatcher.Close()

	go func() {
		log.Print("type on/off")
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			if input.Text() == "on" {
				log.Print("turning led off")
				err = led.SetState(core.StateActive)
				if err != nil {
					log.Fatal(err)
				}
			} else if input.Text() == "off" {
				log.Print("turning led on")
				err = led.SetState(core.StateInactive)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	go func() {
		for ie := range ledWatcher.Watch() {
			log.Printf("%#v", ie)
		}
	}()

	<-ctx.Done()
}
