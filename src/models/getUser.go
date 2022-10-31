package models

import (
	"os"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// obtener usuario por email
func GetUserByEmail(email string) (*User, error) {
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)
	// establecindo filter y opciones
	opts := options.FindOne().SetProjection(bson.M{"code": 0})
	// consulta
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"email": email}, opts).Decode(user)
	return user, err
}

// obtener usuario por telefono
func GetUserByPhone(phone int) (*User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)
	// consulatando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"phone": phone}).Decode(user)
	return user, err
}

// obtener usuario por ID
func GetUserById(id primitive.ObjectID) (*User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)
	// establecindo filter y opciones
	opts := options.FindOne().SetProjection(bson.M{"code": 0})
	// consultando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"_id": id}, opts).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// obtener usuario por ID string
func GetUserByIDStr(id string) (*User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)
	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	// establecindo filter y opciones
	opts := options.FindOne().SetProjection(bson.M{"code": 0})
	// consultando
	user := &User{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}, opts).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func GetUserAndVerifyNotblockExitsAndActive(phone int) (bool, bool, bool, *User, error) {
	// variable entorno
	appName, _ := os.LookupEnv("APP_NAME")
	// Conectando a la BBDD
	ctx, client, db := config.ConnectDB()
	collApp := db.Collection("app")
	collUsers := db.Collection("users")
	defer client.Disconnect(ctx)
	// verificando que no este en modo bloqueado
	appVals := &App{}
	err := collApp.FindOne(ctx, bson.M{"name": appName}).Decode(appVals)
	if err != nil {
		return false, false, false, nil, err
	}
	if appVals.Developing {
		return false, false, false, nil, nil
	}
	// Verificando que el usuario exista
	user := &User{}
	err = collUsers.FindOne(ctx, bson.M{"phone": phone}).Decode(user)
	if err != nil {
		return true, false, false, user, err
	}
	// verificar si esta activo
	if (user.Sleep != 0 || user.SleepDate != time.Time{}) {
		return true, true, false, user, nil
	}

	return true, true, true, user, nil
}
