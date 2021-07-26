package main

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
)

func main() {
	// build zap logger
	zap.NewDevelopment()
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// initialize rpio package and allocate memory
	if err := rpio.Open(); err != nil {
		log.Fatalf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(10)
	pin.Output()

	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second)
	}
}
