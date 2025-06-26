package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/prajwalzaware/go-urlShortner/controller"
	"github.com/prajwalzaware/go-urlShortner/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", controller.Pong)

	app.Post("/url/shorten", middleware.RequireAuth(), middleware.RequireRoleAccess(), controller.ShortenURLHandler)
	app.Get("/url/redirect/:shortCode", middleware.RequireAuth(), middleware.RateLimitMiddleware, middleware.RequireRoleAccess(), controller.RedirectHandler)
	app.Get("/url/stats", middleware.RequireAuth(), middleware.RequireRoleAccess(), controller.GetStats)
	app.Get("/url/AllUrls", middleware.RequireAuth(), middleware.RequireRoleAccess(), controller.GetAllUrls)
	app.Post("/url/deleteShortCode", middleware.RequireAuth(), middleware.RequireRoleAccess(), controller.DeleteShortCode)

	app.Post("/tenant/createNewTenant", middleware.RequireAuth(), middleware.RequireRoleAccess(), controller.CreateNewTenant)
	
	app.Post("/user/signup", controller.SignUpUser)
	app.Post("/user/login", controller.LoginUser)
}
