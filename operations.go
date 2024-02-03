package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Service   string             `json:"service"`
	Amount    float64            `json:"amount"`
	DueDate   time.Time          `json:"dueDate"`
	Paid      bool               `json:"paid"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
}
