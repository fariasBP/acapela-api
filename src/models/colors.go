package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Colors struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"bson,omitempty"`
	Code string             `json:"code" bson:"code,omitempty"`
}
