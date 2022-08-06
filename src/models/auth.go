package models

import (
	"fmt"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// ---- REGISTRADOR ----
// ---- registrador de empleados ----
func AdminEmployRegistrar(name string, codePhone, phone int) error {
	nUserRegister := &User{
		Name:          name,
		Rol:           2,
		CodePhone:     codePhone,
		Phone:         phone,
		Notifications: false,
		CreateDate:    time.Now(),
		UpdateDate:    time.Now(),
	}
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	_, err := coll.InsertOne(ctx, nUserRegister)
	return err
}

// ---- registrador de empleados ----
func EmployRegistrar(name string, codePhone, phone int) error {
	nUserRegister := &User{
		Name:          name,
		Rol:           3,
		CodePhone:     codePhone,
		Phone:         phone,
		Notifications: false,
		CreateDate:    time.Now(),
		UpdateDate:    time.Now(),
	}
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	_, err := coll.InsertOne(ctx, nUserRegister)
	return err
}

// ---- registrador de clientes ----
func ClientRegistrar(name string, codePhone, phone int) error {
	nUserRegister := &User{
		Name:          name,
		Rol:           4,
		CodePhone:     codePhone,
		Phone:         phone,
		Notifications: true,
		CreateDate:    time.Now(),
		UpdateDate:    time.Now(),
	}
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	_, err := coll.InsertOne(ctx, nUserRegister)
	return err
}

// ---- clientes auto registrados (que se registran solos) ----
// func AutoClientRegistrar(name, lastname string, codePhone, phone int) error {
// 	nUserRegister := &User{
// 		Name:          name,
// 		Lastname:      lastname,
// 		Rol:           4,
// 		CodePhone:     codePhone,
// 		Phone:         phone,
// 		Notifications: true,
// 		CreateDate:    time.Now(),
// 		UpdateDate:    time.Now(),
// 	}
// 	ctx, client, coll := config.ConnectColl("users")
// 	defer fmt.Println("Disconnected DB")
// 	defer client.Disconnect(ctx)
// 	_, err := coll.InsertOne(ctx, nUserRegister)
// 	return err
// }

// ---- ADMINISTRADOR ----
// ---- crear super admin ----
func CreateAdminBoss(name, lastname, email string, codePhone, phone int) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// crear el super usuario
	nUser := &User{
		Name:          name,
		Lastname:      lastname,
		Email:         email,
		Rol:           1,
		CodePhone:     codePhone,
		Phone:         phone,
		Notifications: false,
		CreateDate:    time.Now(),
		UpdateDate:    time.Now(),
	}
	_, err := coll.InsertOne(ctx, nUser)
	if err != nil {
		return err
	}
	return nil
}

// ---- verificar si existe adminboss (true=existe) ----
func ExistsAdiminBoss() (b bool) {
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

// ---- obtener codigo de login ----
func SetCode(id primitive.ObjectID, code string) (string, error) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// cosultando (insertando el code en usuario)
	update := bson.M{"$set": bson.M{"code": code, "code_date": time.Now()}}
	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return "", err
	}

	return code, nil

}
