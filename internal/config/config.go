package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
}

func GetConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/baagh/")
	viper.AddConfigPath("$HOME/.baagh")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			log.Fatalf("can't read the config: %v", err)
		}
	}

	return Config{}
}
