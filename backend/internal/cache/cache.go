package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func Connect(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return client
}
