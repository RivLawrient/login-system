package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewAccess() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if err := redis.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return redis
}
