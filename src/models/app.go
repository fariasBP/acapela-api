package models

import (
	"fmt"
	"os"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type (
	App struct {
		ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name       string             `json:"name" bson:"name,omitempty"`
		Developing bool               `json:"developing" bson:"developing,omitempty"`
		Version    string             `json:"version" bson:"version,omitempty"`
	}
)

func CreateApp() error {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("app")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// obteniendo variables de entorno
	name, _ := os.LookupEnv("APP_NAME")
	// campos
	createApp := &App{
		Name:       name,
		Developing: true,
		Version:    "1.0.0",
	}
	// consultando
	_, err := coll.InsertOne(ctx, createApp)

	return err
}

func ExistsAppData() bool {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("app")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// obteniendo variables de entorno
	name, _ := os.LookupEnv("APP_NAME")
	// consultando
	vl := &App{}
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(vl)

	return err == nil
}

func UpdDevelopingApp(dev bool) error {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("app")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	update := bson.M{"$set": bson.M{"developing": dev}}

	_, err := coll.UpdateOne(ctx, bson.M{"Name": "Acapela"}, update)

	return err
}
