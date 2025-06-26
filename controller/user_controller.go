package controllers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/prajwalzaware/go-urlShortner/models"
	"github.com/prajwalzaware/go-urlShortner/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(c *fiber.Ctx) error {
	type userReq struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Role       string `json:"role"`
		TenantName string `json:"tenant_name,omitempty"` // cleaner naming
	}

	var req userReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate basic fields
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, password, and role are required",
		})
	}

	// Validate role
	if req.Role != "admin" && req.Role != "user" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role must be either 'admin' or 'user'",
		})
	}

	// ✅ Require tenant name if role is user
	if req.Role == "user" {
		if req.TenantName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tenant name is required for user role",
			})
		}
	}
	// 🔍 Check if user already exists
	existingUser, err := models.GetUserByEmail(context.Background(), req.Email)
	if err != nil && err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing user",
		})
	}
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// 🔐 Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	var tenantID *int // use pointer so we can store NULL in db

	if req.Role == "user" || (req.Role == "admin" && req.TenantName != "") {
		id, err := utils.GetTenantIDByName(context.Background(), req.TenantName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tenant does not exist",
			})
		}
		tenantID = &id
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		TenantID: tenantID, // NULL for admin, value for user
	}

	err = models.CreateUser(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func LoginUser(c *fiber.Ctx) error {
	type LoginUserReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	log.Println("Login request for email:", req.Email, "Password:", req.Password)
	user, err := models.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}
	token, err := utils.GenerateJWT(user.Email, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}
	return c.JSON(fiber.Map{"token": token})

}
