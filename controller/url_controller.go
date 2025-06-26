package controllers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prajwalzaware/go-urlShortner/config"
	"github.com/prajwalzaware/go-urlShortner/models"
	"github.com/prajwalzaware/go-urlShortner/utils"
)

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

func ShortenURLHandler(c *fiber.Ctx) error {
	type ShortenRequest struct {
		OriginalURL string `json:"original_url"`
	}

	var req ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if !utils.IsValidURL(req.OriginalURL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL format",
		})
	}

	shortCode, err := models.GenerateUniqueShortCode(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate short code",
		})
	}

	email := c.Locals("email").(string)
	user, err := models.GetUserByEmail(context.Background(), email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}
	log.Println("User:", user)
	if user.TenantID == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User does not belong to any tenant",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		utils.CacheUrl(cacheCtx, shortCode, req.OriginalURL)
	}()

	url, err := models.InsertURL(ctx, shortCode, req.OriginalURL, user.ID, *user.TenantID, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store URL",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"short_url": config.GetEnv("BASE_URL", "http://localhost:8080/") + url.ShortCode,
		"original":  url.OriginalURL,
	})
}

func RedirectHandler(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		clickCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		utils.IncrementClick(clickCtx, shortCode)
	}()

	redisURL, err := utils.GetCacheUrl(ctx, shortCode)
	if err == nil && redisURL != "" {
		log.Println("✅ Found in Redis")
		return c.Redirect(redisURL, fiber.StatusMovedPermanently)
	}

	originalURL, err := models.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}
	go utils.CacheUrl(ctx, shortCode, originalURL)

	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}

func GetStats(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	user, err := models.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Println("❌ Failed to fetch user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User lookup failed",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tenantID, err := utils.GetTenantIDFromUser(ctx, user.ID)
	if err != nil {
		log.Println("❌ Error retrieving tenant ID:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tenant ID",
		})
	}
	if tenantID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User does not belong to any tenant",
		})
	}

	var stats []utils.Statastics
	var shortCodes []string

	if role == "admin" {
		stats, err = utils.GetAllStats(ctx, tenantID)
		if err != nil {
			log.Println("❌ Failed to fetch admin stats:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch stats",
			})
		}
	} else if role == "user" {
		stats, err = utils.GetUserStats(ctx, user.ID)
		if err != nil {
			log.Println("❌ Failed to fetch user stats:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch stats",
			})
		}
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized role",
		})
	}

	// 🔁 Fill shortCodes from stats (for both roles)
	for i := range stats {
		shortCodes = append(shortCodes, stats[i].ShortCode)
	}

	// 🧠 Merge Redis Clicks
	redisClicks, err := utils.GetClicksFromRedis(ctx, shortCodes)
	if err != nil {
		log.Println("❌ Failed to get clicks from Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get clicks from Redis",
		})
	}

	for i := range stats {
		for j := range redisClicks {
			if stats[i].ShortCode == redisClicks[j].ShortCode {
				stats[i].Click_count += redisClicks[j].Count
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

func GetAllUrls(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{

			"error": "Unauthorized access",
		})
	}
	user, err := models.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Printf("❌ Error fetching user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	AllUrls, err := utils.GetAllUrlsFromDB(ctx, user.TenantID)
	if err != nil {
		log.Printf("❌ Error fetching URLs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch URLs",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"urls":    AllUrls,
		"message": "URLs fetched successfully",
	})

}

func DeleteShortCode(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	email := c.Locals("email").(string)
	shortCode := c.FormValue("short_code")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := models.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("❌ Error fetching user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	count, err := utils.CheckShortCodeExists(ctx, shortCode)
	if err != nil {
		log.Printf("❌ Error checking short code existence: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check short code existence",
		})
	}

	if !count {
		log.Printf("❌ Short code %s does not exist for user %d", shortCode, user.ID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short code does not exist for user",
		})
	}
	if role == "user" {

		err := utils.DeleteUrlForUser(ctx, shortCode, user.ID)
		if err != nil {
			log.Printf("❌ Error deleting short code for user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete short code",
			})
		}
	} else if role == "admin" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := utils.DeleteUrlForAdmin(ctx, shortCode, user.ID)
		if err != nil {
			log.Printf("❌ Error deleting short code for admin: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete short code",
			})
		}
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Short code deleted successfully",
	})

}
