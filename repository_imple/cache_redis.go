package repositoryimple

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRedisRepo struct {
	rds *redis.Client
}

func NewCacheRedisRepo(rds *redis.Client) CacheRedisRepo {
	return CacheRedisRepo{rds: rds}
}

func (c CacheRedisRepo) StoreStrings(key string, strings []string) error {

	ctx := context.TODO()
	res := c.rds.RPush(ctx, key, strings)

	if res != nil {
		return res.Err()
	}
	resBool := c.rds.Expire(ctx, key, time.Minute)
	return resBool.Err()
}

func (c CacheRedisRepo) RetrieveStrings(key string, startIndex int, endIndex int) ([]string, error) {
	ctx := context.TODO()
	items, err := c.rds.LRange(ctx, key, int64(startIndex), int64(endIndex)).Result()
	return items, err
}
func (c CacheRedisRepo) ClearCache(key string) error {
	ctx := context.TODO()
	err := c.rds.Del(ctx, key).Err()
	return err
}
