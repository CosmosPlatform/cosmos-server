package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type EnvVarConfig struct {
	Key   string
	Value string
}

func readConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); err == nil {
		viper.SetConfigFile(path)
		viper.SetConfigType("json")

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	err := bindEnvVars([]EnvVarConfig{
		{"storage.database_url", "DATABASE_URL"},
		{"auth.jwt_secret", "JWT_SECRET"},
		{"system.default_admin.username", "DEFAULT_ADMIN_USERNAME"},
		{"system.default_admin.email", "DEFAULT_ADMIN_EMAIL"},
		{"system.default_admin.password", "DEFAULT_ADMIN_PASSWORD"},
		{"server.port", "PORT"},
		{"log.level", "LOG_LEVEL"},
		{"sentinel.default_enabled", "SENTINEL_DEFAULT_ENABLED"},
		{"sentinel.default_interval", "SENTINEL_DEFAULT_INTERVAL"},
		{"sentinel.min_interval", "SENTINEL_MIN_INTERVAL"},
		{"sentinel.max_interval", "SENTINEL_MAX_INTERVAL"},
		{"sentinel.sentinel_workers", "SENTINEL_WORKERS"},
		{"token.encryption_key", "TOKEN_ENCRYPTION_KEY"},
		{"mail.smtp_host", "MAIL_SMTP_HOST"},
		{"mail.smtp_port", "MAIL_SMTP_PORT"},
		{"mail.sender_email", "MAIL_SENDER_EMAIL"},
		{"mail.smtp_password", "MAIL_SMTP_PASSWORD"},
	})
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	if err := validateConfig(&conf); err != nil {
		return nil, err
	}

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

func validateConfig(conf *Config) error {
	requiredFields := map[string]string{
		"DATABASE_URL":              conf.StorageConfig.DatabaseURL,
		"JWT_SECRET":                conf.AuthConfig.JWTSecret,
		"DEFAULT_ADMIN_USERNAME":    conf.SystemConfig.DefaultAdmin.Username,
		"DEFAULT_ADMIN_EMAIL":       conf.SystemConfig.DefaultAdmin.Email,
		"DEFAULT_ADMIN_PASSWORD":    conf.SystemConfig.DefaultAdmin.Password,
		"SERVER_PORT":               conf.ServerConfig.Port,
		"LOG_LEVEL":                 conf.LogConfig.Level,
		"SENTINEL_DEFAULT_INTERVAL": conf.SentinelConfig.DefaultInterval,
		"SENTINEL_MIN_INTERVAL":     conf.SentinelConfig.MinInterval,
		"SENTINEL_MAX_INTERVAL":     conf.SentinelConfig.MaxInterval,
		"SENTINEL_WORKERS":          fmt.Sprintf("%d", conf.SentinelConfig.SentinelWorkers),
		"TOKEN_ENCRYPTION_KEY":      conf.TokenConfig.EncryptionKey,
		"MAIL_SMTP_HOST":            conf.MailConfig.SMTPHost,
		"MAIL_SMTP_PORT":            fmt.Sprintf("%d", conf.MailConfig.SMTPPort),
		"MAIL_SENDER_EMAIL":         conf.MailConfig.SenderEmail,
		"MAIL_SMTP_PASSWORD":        conf.MailConfig.SMTPPassword,
	}

	var missingFields []string
	for fieldName, value := range requiredFields {
		if value == "" {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required configuration values: %v", missingFields)
	}

	if err := conf.SentinelConfig.ValidateAndSetDefaults(); err != nil {
		return err
	}

	if err := conf.MailConfig.Validate(); err != nil {
		return err
	}

	return nil
}
