package cache

import (
	"delete-unconfirmed-account/internal/configuration"
	"fmt"

	"github.com/go-redis/redis/v7"
)

func NewRedisClient(conf configuration.RedisConfig) *redis.Client {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Password: conf.Password,
		},
	)
	return rdb
}
