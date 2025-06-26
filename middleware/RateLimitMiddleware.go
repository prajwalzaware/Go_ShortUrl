package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	redisClinet "github.com/prajwalzaware/go-urlShortner/pkg/redis"
)

const (
	RateLimit       = 5                // Max 5 requests
	RateLimitWindow = 60 * time.Second // Per 60 seconds
)

func RateLimitMiddleware(c *fiber.Ctx) error {
	ip := c.IP()
	key := fmt.Sprintf("rate_limit:%s", ip)

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	count, err := redisClinet.Redis.Get(ctx, key).Int()
	if err != nil && err.Error() != "redis: nil" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if count >= RateLimit {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Rate limit exceeded. Try again later.",
		})
	}

	pipe := redisClinet.Redis.TxPipeline()
	pipe.Incr(ctx, key)
	if count == 0 {
		pipe.Expire(ctx, key, RateLimitWindow)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("❌ Redis pipeline error: %v", err)
	}
	return c.Next()

}
