package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Link struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	LinkID    string             `json:"linkID" bson:"linkID"`
	LinkValue string             `json:"linkValue" bson:"linkValue"`
	ExpiresAt *time.Time         `json:"expiresAt" bson:"expiresAt"`
}
