package models

import (
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

/*
	Rol:
		1: Super Admin (Dueño, Administrado) (Todo)
		2: Admin (Trabajador) (crea, actualiza y lee)
		3: Admin (Vendedor) (actualiza y lee)
		4: Client (Cliente) (actualiza y lee)
		5: Common (Posible cliente no registrado) (lee)

		1: Admin
		2: Empl
		3: Client
*/
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name,omitempty"`
	Lastname      string             `json:"lastname" bson:"lastname,omitempty"`
	Email         string             `json:"emails" bson:"email,omitempty"`
	Password      string             `json:"password" bson:"password,omitempty"`
	FirstPassword string             `json:"firstpassword" bson:"firstpassword,omitempty"`
	Rol           int                `json:"rol" bson:"rol,omitempty"`
	Code          string             `json:"code" bson:"code"`
	Phone         int                `json:"phone" bson:"phone"`
	Confirm       bool               `json:"confirm" bson:"confirm"`
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
	_, err := coll.InsertOne(ctx, nUser)
	return err
}
func NewUserRegister(name, lastname, firstpassword, code string, phone int) error {
	nUserRegister := &User{
		Name:          name,
		Lastname:      lastname,
		Email:         "",
		Password:      "",
		FirstPassword: firstpassword,
		Rol:           3,
		Code:          code,
		Phone:         phone,
		Confirm:       false,
	}
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	_, err := coll.InsertOne(ctx, nUserRegister)
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

func GetUserByPhone(code string, phone int) (*User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consulatando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{
		bson.M{"code": code},
		bson.M{"phone": phone},
	}}).Decode(user)
	return user, err
}

func ExistsEmail(email string) (b bool) {
	// Conectandose a la DDBB
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
func ExistsPhone(code string, phone int) (b bool) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	//
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{
		bson.M{"code": code},
		bson.M{"phone": phone},
	}}).Decode(user)
	b = true
	if err != nil {
		b = false
	}
	return
}

func CreateSuperAdmin(name string, lastname string, email string, pwd string, code string, phone int) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// crear el super usuario
	nUser := &User{
		Name:     name,
		Lastname: lastname,
		Email:    email,
		Password: pwd,
		Rol:      1,
		Code:     code,
		Phone:    phone,
		Confirm:  true,
	}
	_, err := coll.InsertOne(ctx, nUser)
	if err != nil {
		return err
	}
	return nil
}
func ExistsSuperuser() (b bool) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// verificar que no exista el superusuario
	b = true
	superUser := &User{}
	err := coll.FindOne(ctx, bson.M{"rol": 1}).Decode(superUser)
	if err != nil {
		b = false
	}
	return
}
