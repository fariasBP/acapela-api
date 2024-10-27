package models

import (
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

const dbUsers = "users"

// ---- actualizar nombre de usuario ----
func UpdUserNameByPhone(phone int, name string) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// consultando
	update := bson.M{"$set": bson.M{"name": name, "wp_registration": false}}
	_, err := coll.UpdateOne(ctx, bson.M{"phone": phone}, update)

	return err
}

// ---- inactivar usuario ----
func UpdInactiveUserByPhone(phone int, sleep uint8) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// consultando
	update := bson.M{
		"$set": bson.M{
			"sleep":      sleep,
			"sleep_date": time.Now(),
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"phone": phone}, update)

	return err
}

// ---- reactivar usuario ----
func UpdReactiveUserByPhone(phone int) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// consultando
	update := bson.M{
		"$set": bson.M{
			"sleep":      0,
			"sleep_date": time.Time{},
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"phone": phone}, update)

	return err
}

// ---- actualizar a no lector ----
/* Contador que indica cuantas veces el usuario no ha leido
sus mensajes recibidos */
func UpdNotReadedUserByPhone(phone int) (int, error) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// obteniendo usuario
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"phone": phone}).Decode(user)
	if err != nil {
		return 0, err
	}
	// consultando Actualizando
	update := bson.M{
		"$set": bson.M{
			"not_readed": user.NotReaded + 1,
		},
	}
	_, err = coll.UpdateOne(ctx, bson.M{"phone": phone}, update)

	return user.NotReaded + 1, err
}

// ---- actualizar a lector ----
/* Actualiza al usuario como activo al recibir o leer sus mensajes*/
func UpdReadedUserByPhone(phone int) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// consultando
	update := bson.M{
		"$set": bson.M{
			"not_readed": 0,
		},
	}
	_, err := coll.UpdateOne(ctx, bson.M{"phone": phone}, update)

	return err
}

// ---- crea un img en espera ----
/* crea un id de una imagen subida a cloudinary, pero que no ha
sido registrado en la BBDD */
func AddWaitingCloudinaryImg(idUser, idImg string) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return err
	}
	user := &User{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(user)
	if err != nil {
		return err
	}
	// consultando
	update := bson.M{
		"$push": bson.M{
			"waiting_cloudinary": idImg,
		},
	}
	_, err = coll.UpdateOne(ctx, bson.M{"_id": objectId}, update)

	return err
}

// ---- elimina un img en espera ----
/* elimina un id de una imagen subida a cloudinary, pero que no ha
sido registrado en la BBDD */
func RemoveWaitingCloudinaryImg(idUser string) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return err
	}
	user := &User{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(user)
	if err != nil {
		return err
	}
	// consultando
	update := bson.M{
		"$unset": bson.M{
			"waiting_cloudinary": 1,
		},
	}
	_, err = coll.UpdateOne(ctx, bson.M{"_id": objectId}, update)

	return err
}

// ---- elimina imgs en espera de Cloudinary ----
/* elimina imagenes de cloudinary */
func DeleteWaitingCloudinaryImgs(idUser string) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl(dbUsers)
	defer client.Disconnect(ctx)
	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return err
	}
	user := &User{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(user)
	if err != nil {
		return err
	}

	// consultando
	update := bson.M{
		"$unset": bson.M{
			"waiting_cloudinary": 1,
		},
	}
	_, err = coll.UpdateOne(ctx, bson.M{"_id": objectId}, update)

	return err
}
