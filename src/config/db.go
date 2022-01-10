package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (context.Context, *mongo.Client) {
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected DB")
	return ctx, client
}
func ConnectDB() (context.Context, *mongo.Client, *mongo.Database) {
	ctx, client := Connect()
	return ctx, client, client.Database("acapela")
}
func ConnectColl(collectionName string) (context.Context, *mongo.Client, *mongo.Collection) {
	ctx, client, db := ConnectDB()
	coll := db.Collection(collectionName)
	return ctx, client, coll
}
