package main

import (
	"bufio"
	"log"
	"os"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/warthog618/gpiod"
)

func main() {
	log.Print("Hello")
	chip, err := core.RequestChip(gpiod.Chips()[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Print("chip registered")

	led, err := chip.RequestItem(10, core.AsOutput(core.Inactive))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("led registered")

	go func() {
		log.Print("type on/off")
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			if input.Text() == "on" {
				log.Print("turning led off")
				led.SetState(core.Active)
			} else if input.Text() == "off" {
				log.Print("turning led on")
				led.SetState(core.Inactive)
			}
		}
	}()

}
