package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress     string
	DBURI          string
	AccrualAddress string
}

func MustLoad() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.RunAddress, "a", ":8000", "Address to run server")
	flag.StringVar(&cfg.DBURI, "d", "", "Datatbase connection string")
	flag.StringVar(&cfg.AccrualAddress, "r", "", "Accural system addres")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		cfg.RunAddress = envRunAddress
	}

	if dbURI := os.Getenv("DATABASE_URI"); dbURI != "" {
		cfg.DBURI = dbURI
	}

	if accrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualAddress != "" {
		cfg.AccrualAddress = accrualAddress
	}

	return &cfg, nil
}
