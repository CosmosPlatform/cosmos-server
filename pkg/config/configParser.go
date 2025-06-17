package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func readConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	err := viper.BindEnv("auth.jwtsecret", "JWT_SECRET")
	if err != nil {
		return nil, fmt.Errorf("failed to bind environment variable: %w", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf Config

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
