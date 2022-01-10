package models

import (
	"context"
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	Lastname string             `json:"lastname" bson:"lastname,omitempty"`
	Email    string             `json:"emails" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
	Rol      int                `json:"rol" bson:"rol,omitempty"`
}

func NewUser(name, lastname, email, pwd string, rol int) error {
	nUser := &User{
		Name:     name,
		Lastname: lastname,
		Email:    email,
		Password: pwd,
		Rol:      rol,
	}
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	_, err := coll.InsertOne(context.Background(), nUser)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	filter := bson.M{
		"email": email,
	}

	user := &User{}
	err := coll.FindOne(ctx, filter).Decode(user)
	return user, err
}

func ExistsEmail(email string) (b bool) {
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	user := &User{}
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(user)
	b = true
	if err != nil {
		b = false
	}
	return
}
