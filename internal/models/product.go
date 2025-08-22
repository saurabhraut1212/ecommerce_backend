package models

import "time"

type Product struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Price       string    `bson:"price" json:"price"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}
