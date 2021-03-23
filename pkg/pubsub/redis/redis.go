package redis

import "github.com/go-redis/redis/v8"

type redisPubSub struct {
	redis *redis.Client
}
