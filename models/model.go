package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string             `bson:"name,omitempty" json:"name,omitempty"`
	Age   int                `bson:"age,omitempty" json:"age,omitempty"`
	Email string             `bson:"email,omitempty" json:"email,omitempty"`
}
