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
	Votes     int64              `json:"votes" bson:"votes"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Posts []Post
