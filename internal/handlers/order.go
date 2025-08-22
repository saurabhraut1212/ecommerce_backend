package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/saurabhraut1212/ecommerce_backend/internal/models"
	"github.com/saurabhraut1212/ecommerce_backend/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	Products *repo.ProductRepo
	Orders   *repo.OrderRepo
}

func NewOrderHandler(pr *repo.ProductRepo, or *repo.OrderRepo) *OrderHandler {
	return &OrderHandler{
		Products: pr,
		Orders:   or,
	}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
		Items  []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	userOID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	var items []models.OrderItem
	var total float64
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	for _, it := range req.Items {
		pid, err := primitive.ObjectIDFromHex(it.ProductID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid product_id"})
		}
		p, err := h.Products.GetById(ctx, pid)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if p == nil {
			return c.Status(404).JSON(fiber.Map{"error": "product not found"})
		}
		if it.Quantity < 1 {
			return c.Status(400).JSON(fiber.Map{"error": "quantity must be >=1"})
		}

		items = append(items, models.OrderItem{
			ProductID: pid,
			Quantity:  it.Quantity,
			Price:     p.Price,
		})
		total += p.Price * float64(it.Quantity)
	}

	order := &models.Order{
		UserID: userOID,
		Items:  items,
		Total:  total,
		Status: "pending",
	}

	if err := h.Orders.Create(ctx, order); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(order)
}

func (h *OrderHandler) Get(c *fiber.Ctx) error {
	idHex := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	o, err := h.Orders.GetById(ctx, oid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if o == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(o)
}

func (h *OrderHandler) ListByUser(c *fiber.Ctx) error {
	userHex := c.Query("user_id", "")
	if userHex == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id required"})
	}
	uid, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	items, err := h.Orders.ListByUser(ctx, uid, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	idHex := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return c.Status(400).JSON(fiber.Map{"error": "status required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	o, err := h.Orders.UpdateStatus(ctx, oid, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if o == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(o)
}

func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	idHex := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Orders.Delete(ctx, oid); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
