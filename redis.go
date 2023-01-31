package main

import (
	"context"
	"log"

	redis "github.com/redis/go-redis/v9"
)

var MainCtx context.Context
var RedisClient *redis.Client

func InitRedis(config ServerConfig) {
	MainCtx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{Addr: config.RedisUrl})
	_, err := RedisClient.Ping(MainCtx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func RedisGet(key string) (string, bool) {
	value, err := RedisClient.Get(MainCtx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false
		} else {
			log.Fatal(err)
		}
	}
	return value, true
}

func RedisSet(key, value string) {
	err := RedisClient.Set(MainCtx, key, value, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
