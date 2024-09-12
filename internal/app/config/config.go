package config

import (
	"flag"
	"os"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

var config Config

func NewConfig() *Config {
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.Address = envRunAddr
	}

	if envbaseUrl := os.Getenv("BASE_URL"); envbaseUrl != "" {
		config.Address = envbaseUrl
	}

	return &config
}

func init() {
	flag.StringVar(&config.Address, "a", "localhost:8080", "Address to run the HTTP server on")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
}
