package utils

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis(host string, port string) {
	addr := host + ":" + port
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // No password by default
		DB:       0,  // Default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis:", addr)
}
