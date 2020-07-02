package config

import "os"

// Config is a system configuration transmitted though env
type Config struct {
	Token string
}

func NewConfig() *Config {
	return &Config{
		Token: os.Getenv("SC_TOKEN"),
	}
}
