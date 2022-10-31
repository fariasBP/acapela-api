package models

import (
	"context"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type ProductModel struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
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
	defer client.Disconnect(ctx)
	// insertando en la BBDD
	_, err := coll.InsertOne(context.Background(), newModel)
	return err
}
func ExistsNameProductModel(name string) (b bool) {
	ctx, client, coll := config.ConnectColl("models")
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
func ExistsModelId(id primitive.ObjectID) bool {
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	model := &ProductModel{}
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(model)

	return err == nil
}
func ExistsModelIdString(id string) bool {
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}

	model := &ProductModel{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}).Decode(model)

	return err == nil
}
func UpdateModelById(id primitive.ObjectID, name, idKind string) error {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	update := bson.M{"$set": bson.M{"name": name, "kind": idKind}}
	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, update)

	return err
}
func DeleteModelById(id primitive.ObjectID) error {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}
