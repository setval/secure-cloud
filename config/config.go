package config

import (
	"flag"
	"os"
)

// Config is a system configuration transmitted though env
type Config struct {
	Token string
	YadiskToken string
	IsSupportDrive bool
}

func NewConfig() *Config {
	isSupportDriveFlag := flag.Bool("drive", false, "Enable google drive")
	flag.Parse()

	var isSupportDrive bool
	if isSupportDriveFlag != nil {
		isSupportDrive = *isSupportDriveFlag
	}

	return &Config{
		Token: os.Getenv("SC_TOKEN"),
		YadiskToken: os.Getenv("YADISK_TOKEN"),
		IsSupportDrive: isSupportDrive,
	}
}
