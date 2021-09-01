package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/warthog618/gpiod"
)

const (
	Open  = 0xdea921
	Close = 0xdea928
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	chip, err := gpiod.NewChip(gpiod.Chips()[0])
	if err != nil {
		log.Fatalf("chip failed: %v", err)
	}
	defer chip.Close()

	led, err := chip.RequestLine(10, gpiod.AsOutput(0))
	if err != nil {
		log.Fatalf("led failed: %v", err)
	}
	defer led.Close()

	reciever, err := gpio.NewReceiver(chip, 27)
	if err != nil {
		log.Fatalf("reciever failed: %v", err)
	}
	defer reciever.Close()

	transmitter, err := gpio.NewTransmitter(chip, 17)
	if err != nil {
		log.Fatalf("transmitter failed: %v", err)
	}
	defer transmitter.Close()

	go func() {
		for c := range reciever.Receive() {
			log.Printf("Signal Recieved: %#v", c)
			if c.Code == Open {
				led.SetValue(1)
			} else if c.Code == Close {
				led.SetValue(0)
			}
		}
	}()

	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			code, err := strconv.ParseUint(input.Text(), 10, 64)
			if err != nil {
				log.Print("ERROR: can't parse the code: ", err)
				continue
			}
			confirm := transmitter.Transmit(code, gpio.DefaultProtocols[0], 350)
			<-confirm
			log.Printf("Signal Sent: %#v", code)
		}
	}()

	<-ctx.Done()
}
