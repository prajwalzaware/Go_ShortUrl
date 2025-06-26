package utils

import (
	"context"
	"log"
	"time"

	"github.com/prajwalzaware/go-urlShortner/pkg/db"
)

type Statastics struct {
	ShortCode   string    `json:"short_code"`
	OriginalUrl string    `json:"original_url"`
	Click_count int       `json:"click_count"`
	CreatedBy   int       `json:"created_by"`
	TenantID    int       `json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetAllStats(ctx context.Context, tenantID int) ([]Statastics, error) {
	query := `SELECT short_code, original_url, click_count, created_by, tenant_id, created_at FROM urls WHERE tenant_id = $1`

	rows, err := db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		log.Println("❌ Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var stats []Statastics

	for rows.Next() {
		var stat Statastics
		err := rows.Scan(&stat.ShortCode, &stat.OriginalUrl, &stat.Click_count, &stat.CreatedBy, &stat.TenantID, &stat.CreatedAt)
		if err != nil {
			log.Println("❌ Error scanning row:", err)
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		log.Println("❌ Row iteration error:", err)
		return nil, err
	}

	log.Println("✅ Retrieved all stats for tenant:", tenantID)
	return stats, nil
}

func GetUserStats(ctx context.Context, userID int) ([]Statastics, error) {
	query := `SELECT short_code, original_url, click_count, created_by, tenant_id, created_at FROM urls WHERE created_by = $1`
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		log.Println("error while executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var stats []Statastics
	for rows.Next() {
		var stat Statastics
		if err := rows.Scan(&stat.ShortCode, &stat.OriginalUrl, &stat.Click_count, &stat.CreatedBy, &stat.TenantID, &stat.CreatedAt); err != nil {
			log.Println("error while scanning row:", err)
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		log.Println("error after iterating rows:", err)
		return nil, err
	}

	return stats, nil
}
