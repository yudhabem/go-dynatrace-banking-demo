package config

import "os"

type Config struct {
	BaseURL string
}

func Load() *Config {

	baseURL := os.Getenv("BANKING_API")

	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	return &Config{
		BaseURL: baseURL,
	}
}
