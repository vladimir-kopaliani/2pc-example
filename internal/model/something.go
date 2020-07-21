package model

import "time"

// Something ...
type Something struct {
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}
