package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl   string
	Version string
}

func Load(path string) (*Config, error) {
	err := godotenv.Load(path)

	if err != nil {

		return nil, err
	}

	return &Config{
		DBUrl:   os.Getenv("POSTGRES"),
		Version: os.Getenv("VERSION"),
	}, nil
}
