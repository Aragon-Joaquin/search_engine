package main

import (
	"search_engine/internal/db"

	"github.com/redis/go-redis/v9"
)

var DBRedis *db.RedisClient

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // Use default DB
		Protocol: 2,
	})

	if client == nil {
		panic("redis server wasn't active")
	}

	DBRedis = &db.RedisClient{
		Db: client,
	}
}
