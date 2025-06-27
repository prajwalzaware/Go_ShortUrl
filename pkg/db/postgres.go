package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prajwalzaware/go-urlShortner/config"
)

var Pool *pgxpool.Pool

func ConnectToDB() {
	dsn := config.GetEnv("DB_URL", "hhhh")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is missing")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Unable to ping database: %v\n", err)
	}
	err = Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("❌ Unable to ping database: %v\n", err)
	}
	fmt.Println("✅ Connected to PostgreSQL using pgx!")

}
func RunSchemaMigration() {
	sqlByte, err := os.ReadFile("pkg/db/schema.sql")
	if err != nil {
		log.Fatalf("❌ Unable to read schema file: %v\n", err)
	}
	_, err = Pool.Exec(context.Background(), string(sqlByte))
	if err != nil {
		log.Fatalf("❌ Unable to run schema migration: %v\n", err)
	}
	fmt.Println("✅ Schema migration completed successfully!")

}
