package database

import "github.com/redis/go-redis/v9"

func NewRedis(ip string, password string) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: password,
	})

	return rdb
}
