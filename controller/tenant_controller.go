package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prajwalzaware/go-urlShortner/utils"
)

type RegisterTenantRequest struct {
	TenantName  string `json:"tenant_name"`
	AdminEmail  string `json:"admin_email"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
}

func CreateNewTenant(c *fiber.Ctx) error {
	var req RegisterTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.TenantName == "" || req.AdminEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant name and admin email are required",
		})
	}
	ctxIsvalid, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	isValidAdmin, err := utils.IsAdminEmailValid(ctxIsvalid, req.AdminEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate admin email",
		})
	}
	if !isValidAdmin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provided admin email is not valid or not an admin user",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := utils.CheckTenantExist(ctx, req.TenantName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking tenant existence",
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Tenant already exists",
		})
	}

	if err := utils.CreateTenant(ctx, req.TenantName, req.AdminEmail, req.Description, req.Industry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tenant",
		})
	}

	errTenats := utils.UpdateAdminTenantStatus(ctx, req.AdminEmail, req.AdminEmail)
	if errTenats != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update admin user with tenant",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "✅ Tenant created successfully",
	})
}
