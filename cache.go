package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

/*
Redis Cache Store, to store Username to UUIDs and Texture Paths
*/
var ctx = context.Background()
var redisClient *redis.Client

// Initialize Redis Client
func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		DB:   0,
	})
}

// Helper function to store skins in redis cache
func cacheSkin(uuid string, cachePath string) {
	redisClient.Set(ctx, "skin:"+uuid, cachePath, 2*time.Hour)
}

// Helper function to store usernames in redis cache
func cacheUsername(username string, uuid string) {
	redisClient.Set(ctx, "username:"+username, uuid, 2*time.Hour)
}
