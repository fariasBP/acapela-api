package models

import (
	"fmt"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

/*
	status
		- 0: no leido
		- 1: leido
	type
		- 0: mensaje de texto
		- 1: accion (boton o funcionalidad preestablecida)

*/

type (
	Message struct {
		ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Msg        string             `json:"msg" bson:"msg,omitempty"`
		From       int                `json:"from" bson:"from,omitempty"`
		To         int                `json:"to" bson:"to,omitempty"`
		Status     int                `json:"status" bson:"status,omitempty"`
		Type       int                `json:"type" bson:"type,omitempty"`
		CreateDate time.Time          `json:"create_date" bson:"create_date,omitempty"`
	}
)

func CreateMsgFromUserToApp(fromPhone int, msg string) error {
	// conectando a la BBDD

	ctx, client, db := config.ConnectDB()
	defer client.Disconnect(ctx)
	collMessages := db.Collection("messages")
	collUsers := db.Collection("users")

	// creando mensaje
	ms := &Message{
		Msg:        msg,
		From:       fromPhone,
		Status:     0,
		Type:       0,
		CreateDate: time.Now().UTC(),
	}
	_, err := collMessages.InsertOne(ctx, ms)
	if err != nil {
		return err
	}

	// estableciendo al usuario como mensaje no leido
	user := User{}
	optU := options.FindOne().SetProjection(bson.M{"count_mailbox": 1})
	err = collUsers.FindOne(ctx, bson.M{"phone": fromPhone}, optU).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	update := bson.M{"$set": bson.M{"mailbox": msg, "count_mailbox": user.CountMailbox + 1}}
	_, err = collUsers.UpdateOne(ctx, bson.M{"phone": fromPhone}, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func CreateMsgFromAppToUser(toPhone int, msg string) error {
	// conectando a la BBDD
	ctx, client, db := config.ConnectDB()
	defer client.Disconnect(ctx)
	collMessages := db.Collection("messages")
	collUsers := db.Collection("users")

	// creando mensaje
	ms := &Message{
		Msg:        msg,
		To:         toPhone,
		Type:       0,
		CreateDate: time.Now().UTC(),
	}
	_, err := collMessages.InsertOne(ctx, ms)

	// marcar como leido mailbox
	update := bson.M{"$set": bson.M{"mailbox": "", "count_mailbox": 0}}
	_, err = collUsers.UpdateOne(ctx, bson.M{"phone": toPhone}, update)

	return err
}

func GetUserMsgsByPhone(phone int) ([]Message, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("messages")
	defer client.Disconnect(ctx)

	// filtro
	filter := bson.M{"$or": []bson.M{
		bson.M{"from": phone},
		bson.M{"to": phone},
	}}

	// opciones
	opts := options.Find().SetSort(bson.M{"create_date": -1})

	// consultando
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetUsersByMailbox() ([]User, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)

	// consultando
	filter := bson.M{"$and": []bson.M{
		bson.M{"count_mailbox": bson.M{"$ne": 0}},
		bson.M{"count_mailbox": bson.M{"$ne": nil}},
	}}
	opts := options.Find().SetProjection(bson.M{"name": 1, "mailbox": 1, "count_mailbox": 1, "photo": 1, "phone": 1})
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// agregando valores
	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
