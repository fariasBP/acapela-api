package models

import (
	"github.com/fariasBP/acapela-api/src/config"
	"gopkg.in/mgo.v2/bson"
)

// Eliminar usuario por telefono
func DelUserByPhone(phone int) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)
	// consultando
	_, err := coll.DeleteOne(ctx, bson.M{"phone": phone})

	return err
}
