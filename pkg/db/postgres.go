package db

import (
	"context"
	"fmt"
	"log"
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
