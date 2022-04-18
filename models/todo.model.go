package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NewTodo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Item      string             `json:"item,omitempty" bson:"item,omitempty" binding:"required,min=1,max=10"`
	Completed *bool              `json:"completed,omitempty" bson:"completed,omitempty" binding:"required"`
}

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Item      string             `json:"item,omitempty" bson:"item,omitempty" binding:"max=10"`
	Completed *bool              `json:"completed,omitempty" bson:"completed,omitempty"`
}
