package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerConfig  `json:"server"`
	StorageConfig `json:"storage"`
}

type ServerConfig struct {
	Host string
	Port string
}

type StorageConfig struct {
	Host string
	Port string
}

func NewConfiguration() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		if conf, err := readConfig("config/local.json"); err != nil {
			return nil, err
		} else {
			return conf, nil
		}
	}

	return nil, fmt.Errorf("unsupported environment: %s", env)
}
