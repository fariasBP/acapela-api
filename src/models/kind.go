package models

import (
	"context"
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type ProductKind struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name,omitempty"`
}

func NewProductKind(name string) error {
	newModel := &ProductKind{
		Name: name,
	}

	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(context.Background(), newModel)
	return err
}
func ExistsNameProductKind(name string) (b bool) {
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	kindProd := &ProductKind{}

	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(kindProd)
	b = true
	if err != nil {
		b = false
	}
	return
}
func GetAllKinds() ([]ProductKind, error) {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var data []ProductKind
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, nil
}
func UpdateNameKind(id primitive.ObjectID, name string) error {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// actualizando
	update := bson.M{"$set": bson.M{"name": name}}
	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, update)

	return err
}
func DeleteKindById(id primitive.ObjectID) error {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// eliminando de la BBDD
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}
func ExistKindId(id primitive.ObjectID) bool {
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	kindModel := &ProductKind{}
	err := coll.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	return err == nil
}
func ExistKindIdString(id string) bool {
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// verificando si id es correcto
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	// consultando
	kindModel := &ProductKind{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	return err == nil
}
