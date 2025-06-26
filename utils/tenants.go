package utils

import (
	"context"
	"log"

	"github.com/prajwalzaware/go-urlShortner/pkg/db"
)

// Check if a tenant with the given name already exists
func CheckTenantExist(ctx context.Context, tenantName string) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM tenants WHERE name = $1`
	err := db.Pool.QueryRow(ctx, query, tenantName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create a new tenant row
func CreateTenant(ctx context.Context, name, adminEmail, description, industry string) error {
	query := `INSERT INTO tenants (name, admin_email, description, industry) VALUES ($1, $2, $3, $4)`
	_, err := db.Pool.Exec(ctx, query, name, adminEmail, description, industry)
	if err != nil {
		log.Printf("❌ Failed to insert tenant: %v", err)
		return err
	}
	log.Println("✅ Tenant created:", name)
	return nil
}

func IsAdminEmailValid(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(1) FROM users WHERE email = $1 AND role = 'admin'`
	var count int
	err := db.Pool.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateAdminTenantStatus(ctx context.Context, adminEmail, userEmail string) error {
	query := `
			UPDATE users
			SET tenant_id = (
			    SELECT id FROM tenants WHERE admin_email = $1
			)
			WHERE email = $2 AND tenant_id IS NULL
			`

	_, err := db.Pool.Exec(ctx, query, adminEmail, userEmail)
	if err != nil {
		log.Printf("❌ Failed to update admin user with tenant: %v", err)
		return err
	}
	log.Println("✅ Admin user updated with tenant:", adminEmail)
	return nil
}

func GetTenantIDByName(ctx context.Context, name string) (int, error) {
	query := `SELECT id FROM tenants WHERE name = $1`
	var id int
	err := db.Pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func GetTenantIDFromUser(ctx context.Context, userID int) (int, error) {
	query := `SELECT tenant_id FROM users WHERE id = $1`

	var tenantID int
	err := db.Pool.QueryRow(ctx, query, userID).Scan(&tenantID)
	if err != nil {
		log.Println("error retrieving tenant ID:", err)
		return 0, err // return the error so the caller can handle it
	}

	log.Println("tenant ID retrieved from user ID:", tenantID)
	return tenantID, nil
}

func GetTenantIDFromShortCode(ctx context.Context, shortCode string) (int, error) {
	query := `SELECT tenant_id FROM urls WHERE short_code = $1`

	var tenantID int
	err := db.Pool.QueryRow(ctx, query, shortCode).Scan(&tenantID)
	if err != nil {
		log.Println("error retrieving tenant ID from short code:", err)
		return 0, err // return the error so the caller can handle it
	}

	log.Println("tenant ID retrieved from short code:", tenantID)
	return tenantID, nil
}
