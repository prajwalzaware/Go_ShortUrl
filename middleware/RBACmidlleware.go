package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	redisClinet "github.com/prajwalzaware/go-urlShortner/pkg/redis"
)

func RequireRoleAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Missing role in context"})
		}
		routePath := c.Route().Path
		routePath = strings.TrimSpace(routePath)
		log.Printf("🔍 Checking access for role: %s on route: %s", role, routePath)
		ctx := context.Background()

		key := "rbac:" + role

		allowed, err := redisClinet.Redis.SIsMember(ctx, key, routePath).Result()
		log.Printf("🔍 Redis check for role %s on route %s: %v", role, routePath, allowed)
		if err != nil {
			log.Printf("❌ Redis error checking role access: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "RBAC Redis error"})
		}
		if !allowed {
			log.Printf("❌ Access denied for role %s on route %s", role, routePath)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		log.Printf("✅ Access Granted. Role %s can access %s", role, routePath)
		return c.Next()
	}
}
