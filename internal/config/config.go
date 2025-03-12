package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	BaseURL 			string
	LogLevel			string
}

func MustGetConfig() *Config {
	var c Config

	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "server start address")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "server address before the short URL")
	flag.StringVar(&c.LogLevel, "l", "info", "level logging")

	flag.Parse()

	if servAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ServerAddress = servAddress
	}
	if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = baseURL
	}
	if baseURL, ok := os.LookupEnv("LOG_LEVEL"); ok {
		c.BaseURL = baseURL
	}

	return &c
}