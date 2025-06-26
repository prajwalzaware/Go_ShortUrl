package redisClinet

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() {
	var ctx = context.Background()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("❌ Failed to parse Redis URL: %v", err)
	}

	Redis = redis.NewClient(opt)
	pong, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	fmt.Println("✅ Connected to Redis:", pong)

}
