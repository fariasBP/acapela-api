package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Comment struct {
		ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Content    string             `json:"content" bson:"content,omitempty"`
		From       string             `json:"from" bson:"from,omitempty"`
		CreateDate time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate time.Time          `json:"update_date" bson:"update_date,omitempty"`
	}
)
