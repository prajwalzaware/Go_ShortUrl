package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prajwalzaware/go-urlShortner/background"
	"github.com/prajwalzaware/go-urlShortner/config"
	"github.com/prajwalzaware/go-urlShortner/pkg/db"
	redisClinet "github.com/prajwalzaware/go-urlShortner/pkg/redis"
	"github.com/prajwalzaware/go-urlShortner/router"
)

func main() {
	config.LoadEnv()
	redisClinet.ConnectRedis()
	db.ConnectToDB()
	db.RunSchemaMigration()

	go background.ClickFlusher()

	app := fiber.New()
	router.SetupRoutes(app)
	port := config.GetEnv("PORT", "3000")

	app.Listen(":" + port)
}
