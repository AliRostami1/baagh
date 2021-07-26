package main

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
)

func main() {
	// viper stuff
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/baagh/")
	viper.AddConfigPath("$HOME/.baagh")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			log.Fatal(err)
		}
	}

	// build zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	// initialize rpio package and allocate memory
	if err := rpio.Open(); err != nil {
		sugar.Fatalf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(10)
	pin.Output()

	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second)
	}
}
