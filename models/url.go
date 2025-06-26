package models

import (
	"context"
	"errors"
	"log"
	"time"

	// "github.com/gofiber/fiber"

	"github.com/google/uuid"
	"github.com/prajwalzaware/go-urlShortner/pkg/db"
)

type URL struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedBy   int       `json:"created_by"`
	TenantID    int       `json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// url, err := models.InsertURL(ctx, shortCode, req.OriginalURL, user.ID, user.TenantID)

func InsertURL(ctx context.Context, shortCode, originalURL string, createdBy int, tenantID int, createdAt time.Time) (*URL, error) {
	query := `
		INSERT INTO urls (short_code, original_url, created_by, tenant_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, short_code, original_url, click_count, created_by, tenant_id, created_at
	`

	row := db.Pool.QueryRow(ctx, query, shortCode, originalURL, createdBy, tenantID, createdAt)

	var url URL
	err := row.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.Clicks, &url.CreatedBy, &url.TenantID, &url.CreatedAt)
	if err != nil {
		return nil, err
	}
	log.Println("✅ URL inserted into DB")
	return &url, nil
}

func GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	query := `SELECT original_url FROM urls WHERE short_code = $1`

	var originalURL string
	err := db.Pool.QueryRow(ctx, query, shortCode).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	log.Println("url retrived from DB")

	return originalURL, nil
}
func GenerateUniqueShortCode(ctx context.Context) (string, error) {
	for i := 0; i < 5; i++ {
		shortCode := uuid.New().String()[:6]
		exists, err := CheckShortCodeExists(ctx, shortCode)
		if err != nil {
			return "", err
		}
		if !exists {
			return shortCode, nil
		}
		log.Println("Short code collision. Retrying...")
	}
	return "", errors.New("could not generate unique short code after multiple attempts")
}

func CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`

	var exists bool
	err := db.Pool.QueryRow(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		return false, err
	}
	log.Println("short code checked in DB")
	return exists, nil
}
