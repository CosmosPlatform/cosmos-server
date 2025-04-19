package config

import "github.com/spf13/viper"

func readConfig() (Config, error) {
	viper.SetConfigFile("config/local.json")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var conf Config

	if err := viper.Unmarshal(&conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
