package config

import (
	"flag"
	"os"
)

// Config is a system configuration transmitted though env
type Config struct {
	Token          string
	IsSupportDrive bool
	Address        string
}

func NewConfig() *Config {
	var cfg = &Config{
		Token: os.Getenv("SC_TOKEN"),
	}

	isSupportDriveFlag := flag.Bool("drive", false, "Enable google drive")
	address := flag.String("address", "[::]:80", "Server address")
	flag.Parse()

	if isSupportDriveFlag != nil {
		cfg.IsSupportDrive = *isSupportDriveFlag
	}
	if address != nil {
		cfg.Address = *address
	}

	return cfg
}
