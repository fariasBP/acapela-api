package models

import (
	"context"
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type KindProduct struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name []string           `json:"name" bson:"name,omitempty"`
}

func NewKindProduct(name []string) error {
	newModel := &KindProduct{
		Name: name,
	}

	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(context.Background(), newModel)
	return err
}

// func ExistsNameKindProduct(name []string) (b bool) {
// 	ctx, client, coll := config.ConnectColl("kinds")
// 	defer fmt.Println("Disconnected DB")
// 	defer client.Disconnect(ctx)

// 	kindModel := &KindProduct{}
// 	cursor, err := coll.Find(ctx, bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	b = true
// 	if err != nil {
// 		b = false
// 	}
// 	return
// }
func GetAllKinds() ([]KindProduct, error) {
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var data []KindProduct
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	fmt.Println(data)
	return data, nil
}
func VerifyKindId(id string) bool {
	ctx, client, coll := config.ConnectColl("kinds")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}

	kindModel := &KindProduct{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	// fmt.Println(kindModel)

	return err == nil
}
