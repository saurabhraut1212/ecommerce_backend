package models

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	Name         string    `bson:"name" json:"name"`
	Email        string    `bson:"email" json:"email"`
	PasswordHash string    `bson:"password" json:"-"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}
