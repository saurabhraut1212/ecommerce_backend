package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/saurabhraut1212/ecommerce_backend/internal/config"
	"github.com/saurabhraut1212/ecommerce_backend/internal/handlers"
	"github.com/saurabhraut1212/ecommerce_backend/internal/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

func New(cfg *config.Config, client *mongo.Client) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	userRepo := repo.NewUserRepo(client.Database(cfg.MongoDB))
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)

	//auth
	app.Post("/api/register", authHandler.Register)
	app.Post("/api/login", authHandler.Login)

	return app
}
