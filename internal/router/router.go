package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/saurabhraut1212/ecommerce_backend/internal/config"
	"github.com/saurabhraut1212/ecommerce_backend/internal/handlers"
	"github.com/saurabhraut1212/ecommerce_backend/internal/middleware"
	"github.com/saurabhraut1212/ecommerce_backend/internal/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

func New(cfg *config.Config, client *mongo.Client) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	//repos
	userRepo := repo.NewUserRepo(client.Database(cfg.MongoDB))
	productRepo := repo.NewProductRepo(client.Database(cfg.MongoDB))
	orderRepo := repo.NewOrderRepo(client.Database(cfg.MongoDB))

	//handlers
	authH := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	productH := handlers.NewProductHandler(productRepo)
	orderH := handlers.NewOrderHandler(productRepo, orderRepo)

	//Health
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Server running") })
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })

	api := app.Group("/api")
	//auth
	api.Post("/register", authH.Register)
	api.Post("/login", authH.Login)

	//products
	api.Get("/products", productH.List)
	api.Get("/products/:id", productH.Get)
	api.Post("/products", middleware.RequireAuth(), productH.Create)
	api.Put("/products/:id", middleware.RequireAuth(), productH.Update)
	api.Delete("/products/:id", middleware.RequireAuth(), productH.Delete)

	//orders
	api.Post("/orders", middleware.RequireAuth(), orderH.Create)
	api.Get("/orders/:id", middleware.RequireAuth(), orderH.Get)
	api.Get("/orders", middleware.RequireAuth(), orderH.ListByUser) // ?user_id=...&page=1&limit=20
	api.Patch("/orders/:id/status", middleware.RequireAuth(), orderH.UpdateStatus)
	api.Delete("/orders/:id", middleware.RequireAuth(), orderH.Delete)

	return app
}
