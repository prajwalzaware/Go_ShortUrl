package utils

import (
	"context"

	"github.com/prajwalzaware/go-urlShortner/pkg/db"
)

func GetAllUrlsFromDB(ctx context.Context, tenantId *int) ([]string, error) {
	query := `SELECT short_code FROM urls where tenant_id = $1`

	rows, err := db.Pool.Query(ctx, query, &tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var shortcode string
		if err := rows.Scan(&shortcode); err != nil {
			return nil, err
		}
		urls = append(urls, shortcode)

	}
	return urls, nil
}
func CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT COUNT(*) FROM urls WHERE short_code = $1`
	var count int
	err := db.Pool.QueryRow(ctx, query, shortCode).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteUrlForUser(ctx context.Context, shortCode string, userID int) error {

	query := `DELETE FROM urls WHERE short_code = $1 AND created_by = $2`
	_, err := db.Pool.Exec(ctx, query, shortCode, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUrlForAdmin(ctx context.Context, shortCode string, userID int) error {

	query := `DELETE FROM urls WHERE short_code = $1 and tenant_id = (SELECT tenant_id FROM users WHERE id = $2)`
	_, err := db.Pool.Exec(ctx, query, shortCode, userID)
	if err != nil {
		return err
	}
	return nil
}
