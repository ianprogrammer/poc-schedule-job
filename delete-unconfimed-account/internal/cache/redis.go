package cache

import (
	"delete-unconfirmed-account/internal/configuration"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type CacheClient struct {
	RedisClient *redis.Client
}

func NewRedisClient(conf configuration.RedisConfig) *CacheClient {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Password: conf.Password,
		},
	)
	return &CacheClient{
		RedisClient: rdb,
	}
}

func (c *CacheClient) Add(key string, value interface{}, exp time.Duration) error {

	err := c.RedisClient.Set(key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}
