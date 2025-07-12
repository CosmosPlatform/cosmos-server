package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type EnvVarConfig struct {
	Key   string
	Value string
}

func readConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	err := bindEnvVars([]EnvVarConfig{
		{"auth.jwt_secret", "JWT_SECRET"},
		{"system.default_admin.username", "DEFAULT_ADMIN_USERNAME"},
		{"system.default_admin.email", "DEFAULT_ADMIN_EMAIL"},
		{"system.default_admin.password", "DEFAULT_ADMIN_PASSWORD"},
	})
	if err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	fmt.Printf("All settings: %+v\n", viper.AllSettings())

	var conf Config

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	fmt.Printf("Configuration: %+v\n", conf)

	return &conf, nil
}

func bindEnvVars(configEnvVars []EnvVarConfig) error {
	for _, configEnvVar := range configEnvVars {
		if err := viper.BindEnv(configEnvVar.Key, configEnvVar.Value); err != nil {
			return fmt.Errorf("Error binding environment variable %s: %v\n", configEnvVar.Key, err)
		}
	}
	return nil
}
