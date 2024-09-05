package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress     string
	DbUri          string
	AccrualAddress string
}

func MustLoad() *Config {
	var cfg Config
	flag.StringVar(&cfg.RunAddress, "a", ":8000", "Address to run server")
	flag.StringVar(&cfg.DbUri, "d", "", "Datatbase connection string")
	flag.StringVar(&cfg.AccrualAddress, "r", "", "Accural system addres")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		cfg.RunAddress = envRunAddress
	}

	if dbUri := os.Getenv("DATABASE_URI"); dbUri != "" {
		cfg.DbUri = dbUri
	}

	if accrualAddress := os.Getenv("RUN_ADDRESS"); accrualAddress != "" {
		cfg.AccrualAddress = accrualAddress
	}

	return &cfg
}
