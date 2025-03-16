package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	BaseURL 			string
	LogLevel			string
	FilePath			string
}

func MustGetConfig() *Config {
	var c Config

	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "server start address")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "server address before the short URL")
	flag.StringVar(&c.LogLevel, "l", "info", "level logging")
	flag.StringVar(&c.FilePath, "f", "temp_storage.txt", "file storage path")

	flag.Parse()

	if servAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ServerAddress = servAddress
	}
	if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = baseURL
	}
	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		c.LogLevel = lvl
	}
	if path, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.FilePath = path
	}

	return &c
}