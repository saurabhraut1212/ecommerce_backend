package models

import "time"

type Order struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	ProductIDs []string  `bson:"product_ids" json:"product_ids"`
	Total      float64   `bson:"total" json:"total"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}
