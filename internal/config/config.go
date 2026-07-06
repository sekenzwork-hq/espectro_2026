package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
}

func Load(path string) (*Config, error) {
	err := godotenv.Load(path)

	if err != nil {

		return nil, err
	}

	return &Config{
		DBUrl: os.Getenv("POSTGRES"),
	}, nil
}
