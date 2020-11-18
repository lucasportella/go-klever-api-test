package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostServiceServer struct {
}

type Post struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	// Vote      int                `json:"vote" bson:"vote"`
}

type Posts []Post
