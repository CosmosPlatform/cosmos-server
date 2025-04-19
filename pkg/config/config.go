package config

import "os"

type Config struct {
	Port string
}

func NewConfiguration() (Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		if conf, err := readConfig(); err != nil {
			return Config{}, err
		} else {
			return conf, nil
		}
	}

	return Config{}, nil
}
