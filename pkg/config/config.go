package config

import (
	"github.com/spf13/viper"
)

type ConfigOptions struct {
	ConfigName  string
	ConfigType  string
	ConfigPaths []string
	EnvPrefix   string
}

type Config struct {
	viper.Viper
}

func New(options *ConfigOptions) (*Config, error) {
	v := viper.New()
	v.SetConfigName(options.ConfigName)
	v.SetConfigType(options.ConfigType)
	for _, path := range options.ConfigPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvPrefix(options.EnvPrefix)
	v.AutomaticEnv()

	return &Config{*v}, nil
}
