package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var Cache *redis.Client
var CacheChannel chan string

func SetUpRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
}

func SetupCacheChannel() {
	CacheChannel = make(chan string)

	go func(ch chan string) {
		for {
			time.Sleep(3 * time.Second)
			Cache.Del(context.Background(), <-ch)
		}
	}(CacheChannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		CacheChannel <- key
	}
}
