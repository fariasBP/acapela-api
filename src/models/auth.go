package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// ---- REGISTRADOR ----
// ---- registrador de empleados ----
func AdminEmployRegistrar(name string, phone int) error {
	// Valores del Usuario
	nUserRegister := &User{
		Name:       strings.ToLower(strings.TrimSpace(name)),
		Rol:        2,
		Phone:      phone,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando
	_, err := coll.InsertOne(ctx, nUserRegister)

	return err
}

// ---- registrador de empleados ----
func EmployRegistrar(name string, phone int) error {
	// valores de usuario
	nUserRegister := &User{
		Name:       strings.ToLower(strings.TrimSpace(name)),
		Rol:        3,
		Phone:      phone,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando
	_, err := coll.InsertOne(ctx, nUserRegister)

	return err
}

// ---- registrador de clientes ----
func ClientRegistrar(name string, phone int) error {
	// valores de usuario
	nUserRegister := &User{
		Name:       strings.ToLower(strings.TrimSpace(name)),
		Rol:        4,
		Phone:      phone,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}
	// conectando a BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando
	_, err := coll.InsertOne(ctx, nUserRegister)

	return err
}

// ---- clientes auto registrados (que se registran solos) ----
func AutoClientRegistrar(name string, phone int) error {
	// valores de usuario
	nUserRegister := &User{
		Name:       strings.ToLower(strings.TrimSpace(name)),
		Rol:        4,
		Phone:      phone,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando
	_, err := coll.InsertOne(ctx, nUserRegister)

	return err
}
func AutoClientRegistrarWithWP(phone int) error {
	// valores del usuario
	nUserRegister := &User{
		Name:               "",
		Rol:                4,
		Phone:              phone,
		WpRegistration:     true,
		CreateDate:         time.Now(),
		UpdateDate:         time.Now(),
		WpRegistrationDate: time.Now(),
	}
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando
	_, err := coll.InsertOne(ctx, nUserRegister)

	return err
}

// ---- ADMINISTRADOR ----
// ---- crear super admin ----
func CreateAdminBoss(name, lastname, email string, phone int) error {
	// valores del super usuario
	nUser := &User{
		Name:       strings.ToLower(strings.TrimSpace(name)),
		Lastname:   strings.ToLower(strings.TrimSpace(name)),
		Email:      strings.TrimSpace(email),
		Rol:        1,
		Phone:      phone,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// insertando
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
	// consultando verificar que no exista el superusuario
	b = true
	superUser := &User{}
	err := coll.FindOne(ctx, bson.M{"rol": 1}).Decode(superUser)
	if err != nil {
		b = false
	}

	return
}

// ---- crear codigo de login ----
func SetCodeByID(id primitive.ObjectID, code string) (string, error) {
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

// ---- crear codigo de login ----
func SetCodeByPhone(phone int, code string) (string, error) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// cosultando (insertando el code en usuario)
	update := bson.M{"$set": bson.M{"code": code, "code_date": time.Now()}}
	_, err := coll.UpdateOne(ctx, bson.M{"phone": phone}, update)
	if err != nil {
		return "", err
	}

	return code, nil

}
