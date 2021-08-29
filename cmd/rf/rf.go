package main

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/warthog618/gpiod"
)

func main() {
	chip, err := gpiod.NewChip(gpiod.Chips()[0])
	if err != nil {
		log.Fatalf("chip failed: %v", err)
	}
	reciever, err := gpio.NewReceiver(chip, 27)
	if err != nil {
		log.Fatalf("reciever failed: %v", err)
	}

	transmitter, err := gpio.NewTransmitter(chip, 17)
	if err != nil {
		log.Fatalf("transmitter failed: %v", err)
	}

	go func() {
		for c := range reciever.Receive() {
			log.Print(c)
		}
	}()

	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			code, err := strconv.ParseUint(input.Text(), 10, 64)
			if err != nil {
				log.Print("ERROR: can't parse the code: ", err)
			}
			transmitter.Transmit(code, gpio.DefaultProtocols[0], 174)
		}
	}()
}
