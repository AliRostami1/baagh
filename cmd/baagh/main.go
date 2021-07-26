package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
)

func main() {
	// build zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

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
			sugar.Fatal(err)
		}
	}

	// React to process signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// initialize rpio package and allocate memory
	if err := rpio.Open(); err != nil {
		sugar.Fatalf("can't open and memory map GPIO memory range from /dev/mem: %v", err)
	}
	defer rpio.Close()

	pin := rpio.Pin(10)
	pin.Output()

	for x := 0; x < 20; x++ {
		select {
		case sig := <-sigs:
			pin.Low()
			sugar.Info(sig)
			os.Exit(1)
		default:
		}
		pin.Toggle()
		time.Sleep(time.Second)
	}
}
