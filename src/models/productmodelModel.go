package models

import (
	"context"
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type ProductModel struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name,omitempty"`
	Kind string             `json:"kind" bson:"kind,omitempty"`
}

func NewProductModel(name string, idKind string) error {
	newModel := &ProductModel{
		Name: name,
		Kind: idKind,
	}
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("models")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando en la BBDD
	_, err := coll.InsertOne(context.Background(), newModel)
	return err
}
func ExistsNameProductModel(name string) (b bool) {
	ctx, client, coll := config.ConnectColl("models")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	productModel := &ProductModel{}
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(productModel)
	b = true
	if err != nil {
		b = false
	}
	return
}
func GetAllModels() ([]ProductModel, error) {
	ctx, client, coll := config.ConnectColl("models")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var data []ProductModel
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, nil
}
