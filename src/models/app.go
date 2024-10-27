package models

import (
	"os"
	"strconv"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// app.go SE TIENE QUE ELIMINAR
type (
	App struct {
		ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"` // a que elimnar
		Name        string             `json:"name" bson:"name,omitempty"`
		Developing  bool               `json:"developing" bson:"developing,omitempty"`
		Version     string             `json:"version" bson:"version,omitempty"`
		SetProducts time.Time          `json:"set_products" bson:"set_products,omitempty"` // a que eliminar
	}
)

func CreateApp() error {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("app")
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

func UpdDevelopingApp(dev bool) error {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("app")
	defer client.Disconnect(ctx)

	update := bson.M{"$set": bson.M{"developing": dev}}

	_, err := coll.UpdateOne(ctx, bson.M{"name": "Acapela"}, update)

	return err
}

func GetDataApp() (error, *App) {
	// obteniendo dataapp de .env
	name, _ := os.LookupEnv("NAMEAPP")
	version, _ := os.LookupEnv("VERSIONAPP")
	developing, _ := os.LookupEnv("DEVELOPINGMODE")
	// convirtiendo valores
	developingMode, err := strconv.ParseBool(developing)
	if err != nil {
		developingMode = false
	}
	// obteniendo valores
	dat := &App{
		Name:       name,
		Version:    version,
		Developing: developingMode,
	}

	return err, dat
}
