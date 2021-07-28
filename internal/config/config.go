package config

import (
	"github.com/spf13/viper"
)

type Config struct {
}

func New() (conf *Config, err error) {
	// get environment variables, flags and ...

	// parse config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/baagh/")
	viper.AddConfigPath("$HOME/.baagh")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return conf, nil
}
