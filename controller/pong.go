package controllers

import "github.com/gofiber/fiber/v2"

func Pong(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong from controller",
	})
}
