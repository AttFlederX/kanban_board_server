package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GoogleID string             `json:"google_id" bson:"google_id"`
	Email    string             `json:"email" bson:"email"`
	Name     string             `json:"name" bson:"name"`
	PhotoURL string             `json:"photourl" bson:"photourl"`
}
