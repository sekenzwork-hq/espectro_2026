package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	Version       string
	RedisIP       string
	RedisPassword string
}

func Load(path string) (*Config, error) {
	err := godotenv.Load(path)

	if err != nil {

		return nil, err
	}

	return &Config{
		DBUrl:         os.Getenv("POSTGRES"),
		Version:       os.Getenv("VERSION"),
		RedisIP:       os.Getenv("REDIS_IP"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}, nil
}
