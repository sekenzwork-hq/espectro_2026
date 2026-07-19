package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ip string) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: ip,
	})

	status := rdb.Ping(context.TODO())
	fmt.Println(status)

	return rdb
}
