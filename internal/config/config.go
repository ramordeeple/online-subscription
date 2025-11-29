package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Printf("Warning: no .env file found at %s", path)
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
