package config

type Config struct {
	ServerConfig  `mapstructure:"server"`
	StorageConfig `mapstructure:"storage"`
	AuthConfig    `mapstructure:"auth"`
	SystemConfig  `mapstructure:"system"`
	LogConfig     `mapstructure:"log"`
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

type DefaultAdmin struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

func NewConfiguration() (*Config, error) {
	if conf, err := readConfig("config/local.json"); err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}
