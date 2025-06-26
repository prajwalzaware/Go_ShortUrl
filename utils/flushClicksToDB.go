package utils

import (
	"context"
	"log"
	"strings"

	redisClinet "github.com/prajwalzaware/go-urlShortner/pkg/redis"
)

func FlushClicksToDB(ctx context.Context) {
	iter := redisClinet.Redis.Scan(ctx, 0, "click:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		shortCode := strings.TrimPrefix(key, "click:")
		count, err := redisClinet.Redis.Get(ctx, key).Int()
		if err != nil {
			log.Printf("❌ Failed to get click count: %v", err)
			continue
		}

		err = IncrementClickCountIbDb(shortCode, count)
		if err != nil {
			log.Printf("❌ Failed to update DB click count: %v", err)
			continue
		}

		redisClinet.Redis.Del(ctx, key)
	}

}
