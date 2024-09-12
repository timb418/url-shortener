package config

import (
	"flag"
)

type Config struct {
	Address string
	BaseURL string
}

var config Config

func NewConfig() *Config {
	flag.Parse()

	return &config
}

func init() {
	flag.StringVar(&config.Address, "a", "localhost:8080", "Address to run the HTTP server on")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
}
