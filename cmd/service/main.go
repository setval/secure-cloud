package main

import (
	"github.com/DiscoreMe/SecureCloud/config"
	"github.com/DiscoreMe/SecureCloud/internal"
)

func main() {
	cfg := config.NewConfig()

	serv := internal.NewServer(internal.ServerConfig{
		ValidToken: cfg.Token,
	})

	serv.SetupAPI()

	if err := serv.EnableLocalStorage(); err != nil {
		panic(err)
	}

	if cfg.IsSupportDrive {
		if err := serv.EnableDriveStorage(); err != nil {
			if err := serv.UpdateDriveToken(); err != nil {
				panic(err)
			}
		}
	}

	if cfg.IsSupportS3 {
		if err := serv.EnableS3Storage(cfg.S3.Endpoint, cfg.S3.Bucket, cfg.S3.AccessKey, cfg.S3.SecretKey, cfg.S3.Location); err != nil {
			panic(err)
		}
	}

	if err := serv.Listen(cfg.Address); err != nil {
		panic(err)
	}
}
