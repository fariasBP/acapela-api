package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

/*
	Rol:
		1: SuperAdmin (Dueño, Administrado) (Todo)
		2: AdminEmploye (Empleado Mayor) (crea, actualiza y lee)
		3: Employe (Empleado Menor) (actualiza y lee)
		4: Client (Cliente) (actualiza y lee)

	Status:
		0: Normal (recibe notificaciones)
		1: Descanzo (no recibe notificaciones)
		2: Eliminado (ha eliminado su cuenta "supestamente")
*/
type (
	User struct {
		ID                 primitive.ObjectID `bson:"_id,omitempty"`
		Name               string             `json:"name" bson:"name,omitempty"`
		Lastname           string             `json:"lastname" bson:"lastname,omitempty"`
		Email              string             `json:"email" bson:"email,omitempty"`
		Code               string             `json:"code" bson:"code,omitempty"`
		Rol                int                `json:"rol" bson:"rol,omitempty"`
		Phone              int                `json:"phone" bson:"phone"`
		Notifications      bool               `json:"notifications" bson:"notifications"`
		WpRegistration     bool               `json:"wp_registration" bson:"wp_registration"`
		CreateDate         time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate         time.Time          `json:"update_date" bson:"update_date,omitempty"`
		CodeDate           time.Time          `json:"code_date" bson:"code_date,omitempty"`
		SleepDate          time.Time          `json:"inactive_date" bson:"inactive_date,omitempty"`
		WpRegistrationDate time.Time          `json:"wp_registration_date" bson:"wp_registration_date,omitempty"`
	}
)

// ---- OBTENER USUARIO ----
// obtener usuario por email
func GetUserByEmail(email string) (*User, error) {
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
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
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// establecindo filter y opciones
	opts := options.FindOne().SetProjection(bson.M{"code": 0})
	// consulatando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"phone": phone}, opts).Decode(user)
	return user, err
}

// obtener usuario por ID
func GetUserById(id primitive.ObjectID) (*User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
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
	defer fmt.Println("Disconnected DB")
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

// ---- OBTNER USUARIOS ----
// obtener usuarios
func GetUsers() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// establecindo filter y opciones
	opts := options.Find().SetProjection(bson.M{"code": 0})
	// consultando
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// obtener numero nombre y si recibe notificaiones de usuario
func GetPhoneNameNotificationsFromUsers() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetProjection(bson.M{"phone": 1, "name": 1, "notifications": 1})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ---- obtener numero y nombre para notificaciones ----
func GetPhoneAndNameForNotifications() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetProjection(bson.M{"name": 1, "phone": 1, "code": 1})
	filter := bson.M{"notifications": true}
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ---- verificar si existe un email (true=existe) ----
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

// ---- verificar si existe un numero (true=existe) ----
func ExistsPhone(phone int) (b bool) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"phone": phone}).Decode(user)
	b = true
	if err != nil {
		b = false
	}
	return
}

// ---- verificar si existe el id string del vendedor ----
func ExistsSellerIDStr(id string) bool {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// verificando si id es correcto
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	// filtrar
	filter := bson.M{"$and": bson.M{
		"_id": ObjId,
		"rol": 4,
	}}
	// consultando
	userModel := &User{}
	err = coll.FindOne(ctx, filter).Decode(userModel)

	return err != nil
}

// ---- verificar si existe el id string del comprador ----
func ExistsBuyerIDStr(id string) bool {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// verificando si id es correcto
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	// consultando
	userModel := &User{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}).Decode(userModel)
	if userModel.Rol != 4 {
		return false
	}

	return err == nil
}

// ---- ACTUALIZAR USUARIO ----
func UpdUserNameByPhone(phone int, name string) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	update := bson.M{"$set": bson.M{"name": strings.ToLower(strings.TrimSpace(name)), "wp_registration": false}}
	_, err := coll.UpdateOne(ctx, bson.M{"phone": phone}, update)

	return err
}