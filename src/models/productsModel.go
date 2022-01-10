package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Gender    uint8              `json:"gender" bson:"gender,omitempty"`
	Models    []string           `json:"models" bson:"models,omitempty"`
	MinPrice  uint               `json:"minprice" bson:"minprice,omitempty"`
	Price     uint               `json:"price" bson:"price,omitempty"`
	PriceSold uint               `json:"pricesold" bson:"pricesold,omitempty"`
	PricaSale uint               `json:"pricesale" bson:"pricesale,omitempty"`
	SoldOut   bool               `json:"soldout" bson:"soldout,omitempty"`
	Pattern   ProductPattern     `json:"pattern" bson:"pattern,omitempty"`
	PhotoUrl  string             `json:"photourl" bson:"photourl,omitempty"`
}
