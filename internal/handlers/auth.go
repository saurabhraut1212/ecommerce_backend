package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/saurabhraut1212/ecommerce_backend/internal/models"
	"github.com/saurabhraut1212/ecommerce_backend/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo  *repo.UserRepo
	JWTSecret string
}

func NewAuthHandler(userRepo *repo.UserRepo, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (h *AuthHandler) register(c *fiber.Ctx) error {
	req := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := h.UserRepo.Create(context.Background(), user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	user, err := h.UserRepo.FindByEmail(context.Background(), req.Email)

	if err != nil || user == nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, _ := token.SignedString([]byte(h.JWTSecret))
	return c.JSON(fiber.Map{"token": tokenStr})
}
