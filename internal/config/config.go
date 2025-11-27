package config

import "os"

type Config struct {
	Port string
	DSN  string
}

func LoadConfig() *Config {
	return &Config{
		Port: os.Getenv("PORT"),
		DSN:  os.Getenv("POSTGRES_DSN"),
	}
}
