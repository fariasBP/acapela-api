package models

import (
	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type (
	Permission struct {
		ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name  string             `json:"name" bson:"name,omitempty"`
		Route string             `json:"route" bson:"route,omitempty"`
	}
)

const dbPermissions = "permissions"

func CreatePermission(name, route string) error {
	// conenctado a BBDD
	ctx, client, coll := config.ConnectColl(dbPermissions)
	defer client.Disconnect(ctx)

	// valores
	newPermission := &Permission{
		Name:  name,
		Route: route,
	}

	// creando el permiso
	_, err := coll.InsertOne(ctx, newPermission)

	return err
}

func ExistsPermisionRoute(route string) bool {
	// conenctado a BBDD
	ctx, client, coll := config.ConnectColl(dbPermissions)
	defer client.Disconnect(ctx)

	// consultando
	p := &Permission{}
	err := coll.FindOne(ctx, bson.M{"route": route}).Decode(p)

	return err != nil
}
