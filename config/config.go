package config

import (
	"flag"
	"os"
)

// Config is a system configuration transmitted though env
type Config struct {
	Token          string
	IsSupportDrive bool
	IsSupportS3    bool
	Address        string
	S3             S3Config
}

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Location  string
}

func NewConfig() *Config {
	var cfg = &Config{
		Token: os.Getenv("SC_TOKEN"),
		S3: S3Config{
			Endpoint:  os.Getenv("SC_S3_ENDPOINT"),
			AccessKey: os.Getenv("SC_S3_ACCESS_KEY"),
			SecretKey: os.Getenv("SC_S3_SECRET_KEY"),
			Bucket:    os.Getenv("SC_S3_BUCKET"),
			Location:  os.Getenv("SC_S3_LOCATION"),
		},
	}

	isSupportDriveFlag := flag.Bool("drive", false, "Enable google drive")
	isSupportS3Flag := flag.Bool("s3", false, "Enable s3 storage")
	address := flag.String("address", "[::]:80", "Server address")
	flag.Parse()

	if isSupportDriveFlag != nil {
		cfg.IsSupportDrive = *isSupportDriveFlag
	}
	if isSupportS3Flag != nil {
		cfg.IsSupportS3 = *isSupportS3Flag
	}
	if address != nil {
		cfg.Address = *address
	}

	return cfg
}
