package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/AliRostami1/baagh/internal/logy"
	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/AliRostami1/baagh/pkg/controller/rf"
	"github.com/warthog618/gpiod"
)

const (
	Open  = 0xdea921
	Close = 0xdea928
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logger, err := logy.New(ctx, logy.InfoLevel)
	if err != nil {
		log.Print(err)
	}
	core.SetLogger(logger)

	chipName := gpiod.Chips()[0]

	led, err := core.RequestItem(chipName, 10, core.AsOutput(core.StateInactive))
	if err != nil {
		logger.Fatalf("led failed: %v", err)
	}
	defer led.Close()
	logger.Info("led registered")

	reciever, err := rf.NewReceiver(chipName, 27)
	if err != nil {
		log.Fatalf("reciever failed: %v", err)
	}
	defer reciever.Close()
	logger.Info("reciever registered")

	transmitter, err := rf.NewTransmitter(chipName, 17)
	if err != nil {
		log.Fatalf("transmitter failed: %v", err)
	}
	defer transmitter.Close()
	logger.Info("transmitter registered")

	go func() {
		logger.Info("ready to recieve signals")
		for c := range reciever.Receive() {
			log.Printf("Signal Recieved: %#v", c)
			if c.Code == Open {
				led.SetState(core.StateActive)
			} else if c.Code == Close {
				led.SetState(core.StateInactive)
			}
		}
	}()

	go func() {
		input := bufio.NewScanner(os.Stdin)
		logger.Info("ready to send signals")
		// transmitter.Transmit(131564, rf.DefaultProtocols[0], 350)
		for input.Scan() {
			code, err := strconv.ParseUint(input.Text(), 10, 64)
			if err != nil {
				log.Print("ERROR: can't parse the code: ", err)
				continue
			}
			confirm := transmitter.Transmit(code, rf.DefaultProtocols[0], 350)
			<-confirm
			log.Printf("Signal Sent: %#v", code)
		}
	}()

	<-ctx.Done()
}
