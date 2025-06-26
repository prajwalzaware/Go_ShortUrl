package models

import (
	"context"
	"time"

	"github.com/prajwalzaware/go-urlShortner/pkg/db"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	TenantID  *int      `json:"tenant_id"` // pointer for null support
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (name, email, password, role,tenant_id) VALUES ($1, $2, $3, $4,$5)`
	_, err := db.Pool.Exec(ctx, query, user.Name, user.Email, user.Password, user.Role, user.TenantID)
	return err

}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password, role, tenant_id, created_at FROM users WHERE email = $1`
	row := db.Pool.QueryRow(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.TenantID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
