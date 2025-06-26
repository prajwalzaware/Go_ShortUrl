package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	redisClinet "github.com/prajwalzaware/go-urlShortner/pkg/redis"
	"github.com/redis/go-redis/v9"
)

func CacheUrl(ctx context.Context, shortCode, originalURL string) {
	err := redisClinet.Redis.Set(ctx, shortCode, originalURL, 24*time.Hour).Err()
	if err != nil {
		log.Printf("❌ Failed to cache in Redis: %v", err)
	}
	log.Println("set in redis")
}

func GetCacheUrl(ctx context.Context, shortCode string) (string, error) {
	originalURL, err := redisClinet.Redis.Get(ctx, shortCode).Result()
	if err != nil {
		log.Printf("❌ url not found in redis: %v", err)
		return "", err
	}
	return originalURL, nil
}

func IncrementClick(ctx context.Context, shortCode string) {
	key := fmt.Sprintf("click:%s", shortCode)
	err := redisClinet.Redis.Incr(ctx, key).Err()
	if err != nil {
		log.Printf("❌ Failed to increment click count: %v", err)
	}
}

func GetRedisTemporaryClick(ctx context.Context, shortCode string) (int, error) {
	key := fmt.Sprintf("click:%s", shortCode)
	temporaryRedisCount, err := redisClinet.Redis.Get(ctx, key).Result()
	log.Printf("Temporary Click Count from Redis: %s", temporaryRedisCount)
	if err != nil {
		// If the error is redis.Nil, treat as zero clicks (key does not exist)
		if err == redis.Nil {
			log.Printf("Click count not found in Redis, returning 0")
			return 0, nil
		}
		log.Printf("❌ Failed to get click count from redis: %v", err)
		return 0, err
	}
	strconvCount, err := strconv.Atoi(temporaryRedisCount)
	if err != nil {
		log.Printf("❌ Failed to convert click count to int: %v", err)
		return 0, err

	}
	return strconvCount, nil

}

type ClickCount struct {
	ShortCode string `json:"short_code"`
	Count     int    `json:"count"`
}

func GetClicksFromRedis(ctx context.Context, shortCodes []string) ([]ClickCount, error) {
	var clickCounts []ClickCount
	for _, shortCode := range shortCodes {
		count, err := GetRedisTemporaryClick(ctx, shortCode)
		if err != nil {
			log.Printf("❌ Error getting click count for %s: %v", shortCode, err)
			return nil, err
		}
		clickCounts = append(clickCounts, ClickCount{ShortCode: shortCode, Count: count})
	}
	
	return clickCounts, nil
}
