package utils

import (
	"context"

	"github.com/prajwalzaware/go-urlShortner/pkg/db"
)

func IncrementClickCountIbDb(shortCode string, count int) error {
	query := `UPDATE urls SET click_count = click_count + $1 WHERE short_code = $2`
	_, err := db.Pool.Exec(context.Background(), query, count, shortCode)
	return err
}
