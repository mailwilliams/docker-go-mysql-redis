package database

import "github.com/go-redis/redis/v8"

var Cache *redis.Client

func SetUpRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}
