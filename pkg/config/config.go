package config

import (
	"fmt"
	"time"
)

type Config struct {
	ServerConfig   `mapstructure:"server"`
	StorageConfig  `mapstructure:"storage"`
	AuthConfig     `mapstructure:"auth"`
	SystemConfig   `mapstructure:"system"`
	TokenConfig    `mapstructure:"token"`
	LogConfig      `mapstructure:"log"`
	SentinelConfig `mapstructure:"sentinel"`
	MailConfig     `mapstructure:"mail"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type StorageConfig struct {
	DatabaseURL string `mapstructure:"database_url"`
}

type SystemConfig struct {
	DefaultAdmin `mapstructure:"default_admin"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type TokenConfig struct {
	EncryptionKey string `mapstructure:"encryption_key"`
}

type SentinelConfig struct {
	DefaultEnabled         bool   `mapstructure:"default_enabled"`
	DefaultInterval        string `mapstructure:"default_interval"`
	MinInterval            string `mapstructure:"min_interval"`
	MaxInterval            string `mapstructure:"max_interval"`
	DefaultIntervalSeconds int
	MinIntervalSeconds     int
	MaxIntervalSeconds     int
	SentinelWorkers        *int `mapstructure:"sentinel_workers"`
}

type DefaultAdmin struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

type MailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SenderEmail  string `mapstructure:"sender_email"`
	SMTPPassword string `mapstructure:"smtp_password"`
}

func NewConfiguration() (*Config, error) {
	if conf, err := readConfig("config/local.json"); err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}

func (sc *SentinelConfig) ValidateAndSetDefaults() error {
	defaultInterval, err := time.ParseDuration(sc.DefaultInterval)
	if err != nil {
		return err
	}
	sc.DefaultIntervalSeconds = int(defaultInterval.Seconds())

	minInterval, err := time.ParseDuration(sc.MinInterval)
	if err != nil {
		return err
	}
	sc.MinIntervalSeconds = int(minInterval.Seconds())

	maxInterval, err := time.ParseDuration(sc.MaxInterval)
	if err != nil {
		return err
	}
	sc.MaxIntervalSeconds = int(maxInterval.Seconds())

	if sc.MinIntervalSeconds <= 0 {
		return fmt.Errorf("min_interval must be greater than 0")
	}

	if sc.MaxIntervalSeconds < sc.MinIntervalSeconds {
		return fmt.Errorf("max_interval must be greater than or equal to min_interval")
	}

	if sc.DefaultIntervalSeconds < sc.MinIntervalSeconds || sc.DefaultIntervalSeconds > sc.MaxIntervalSeconds {
		return fmt.Errorf("default_interval must be between min_interval and max_interval")
	}

	return nil
}

func (mc *MailConfig) Validate() error {
	if mc.SMTPPort <= 0 || mc.SMTPPort > 65535 {
		return fmt.Errorf("invalid SMTP port: %d", mc.SMTPPort)
	}

	return nil
}
