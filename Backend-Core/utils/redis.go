package utils

import (
	"context"
	"myproject/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() *redis.Client {
	cfg := config.GetConfig()

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		opt = &redis.Options{
			Addr:     cfg.RedisURL,
			Password: cfg.RedisPassword,
		}
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	RedisClient = client
	return client
}
