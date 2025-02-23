package config

import "flag"

type Config struct {
	ServerAddress string
	PublicAddress string
}

func MustGetConfig() *Config {
	var c Config

	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "server start address")
	flag.StringVar(&c.PublicAddress, "b", "http://localhost:8080/", "server address before the short URL")

	flag.Parse()

	return &c
}