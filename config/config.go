package config

import (
	"os"
	"time"
)

type Config struct {
	Host     string
	Addr     string
	DBconfig *DBconfig
}

type DBconfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// Load reads all configuration from environment variables.
// Call this once in main.go after loading the .env file.
func Load() *Config {
	return &Config{
		Addr:     ":4222",
		DBconfig: LoadDBConfig(),
	}
}

func LoadDBConfig() *DBconfig {
	return &DBconfig{
		URI:      os.Getenv("MONGODB_URI"),
		Database: os.Getenv("MONGODB_DATABASE"),
		Timeout:  10 * time.Second,
	}
}
