package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"` // snapshot at time of order
}

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items     []OrderItem        `bson:"items" json:"items"`
	Total     float64            `bson:"total" json:"total"`
	Status    string             `bson:"status" json:"status"` // pending, paid, shipped, cancelled
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
