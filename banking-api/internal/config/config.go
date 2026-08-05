package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {

	// Application
	AppName string
	AppEnv  string
	AppPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Logging
	LogLevel string

	// Feature Flag
	EnableTracing bool
	EnableMetrics bool

	// Build Info
	Version string
	Build   string
}

func Load() *Config {

	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return &Config{

		AppName: viper.GetString("APP_NAME"),
		AppEnv:  viper.GetString("APP_ENV"),
		AppPort: viper.GetString("APP_PORT"),

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSSLMode:  viper.GetString("DB_SSLMODE"),

		LogLevel: viper.GetString("LOG_LEVEL"),

		EnableTracing: viper.GetBool("ENABLE_TRACING"),
		EnableMetrics: viper.GetBool("ENABLE_METRICS"),

		Version: viper.GetString("APP_VERSION"),
		Build:   viper.GetString("APP_BUILD"),
	}
}
