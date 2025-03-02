package utils

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Use the container name as the hostname
		Password: "",           // No password
		DB:       0,            // Default DB
	})

	val, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Printf("Redis Initiated!, %v", val)
}
