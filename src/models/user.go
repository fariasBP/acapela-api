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
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		Name          string             `json:"name" bson:"name,omitempty"`
		Lastname      string             `json:"lastname" bson:"lastname,omitempty"`
		Email         string             `json:"email" bson:"email,omitempty"`
		Code          string             `json:"code" bson:"code,omitempty"`
		Rol           int                `json:"rol" bson:"rol,omitempty"`
		CodePhone     int                `json:"code_phone" bson:"code_phone"`
		Phone         int                `json:"phone" bson:"phone"`
		Notifications bool               `json:"notifications" bson:"notifications"`
		CreateDate    time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate    time.Time          `json:"update_date" bson:"update_date,omitempty"`
		CodeDate      time.Time          `json:"code_date" bson:"code_date,omitempty"`
	}
	UserInfo struct {
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		Name          string             `json:"name" bson:"name,omitempty"`
		Lastname      string             `json:"lastname" bson:"lastname,omitempty"`
		Email         string             `json:"email" bson:"email,omitempty"`
		FirstPassword string             `json:"firstpassword" bson:"firstpassword,omitempty"`
		Rol           int                `json:"rol" bson:"rol,omitempty"`
		Code          string             `json:"code" bson:"code"`
		Phone         int                `json:"phone" bson:"phone"`
		Notifications bool               `json:"notifications" bson:"notifications"`
	}
)

// ---- OBTENER USUARIOS ----
// ---- obtener usuario por email ----
func GetUserByEmail(email string) (*User, error) {
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// filtros
	filter := bson.M{
		"email": email,
	}
	// consulta
	user := &User{}
	err := coll.FindOne(ctx, filter).Decode(user)
	return user, err
}

func GetUserByPhone(codePhone, phone int) (*User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consulatando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{
		bson.M{"code_phone": codePhone},
		bson.M{"phone": phone},
	}}).Decode(user)
	return user, err
}

func GetUserById(id primitive.ObjectID) (*UserInfo, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// // convirtiendo id en ObjectId
	// objectId, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, err
	// }
	// establecindo filter y opciones
	filter := bson.M{"_id": id}
	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	// consultando
	user := &UserInfo{}
	err := coll.FindOne(ctx, filter, opts).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func GetUsers() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetProjection(bson.M{"code_phone": 1, "phone": 1, "name": 1, "notifications": 1})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func GetPhoneNameNotificationsFromUsers() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetProjection(bson.M{"code_phone": 1, "phone": 1, "name": 1, "notifications": 1})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
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
func ExistsPhone(codePhone, phone int) (b bool) {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{
		bson.M{"code_phone": codePhone},
		bson.M{"phone": phone},
	}}).Decode(user)
	b = true
	if err != nil {
		b = false
	}
	return
}
func ExistsSellerID(id string) bool {
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
	opts := options.FindOne().SetProjection(bson.M{"password": 0, "code": 0})
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}, opts).Decode(userModel)
	if userModel.Rol == 4 {
		return false
	}

	return err == nil
}
func ExistsBuyerID(id string) bool {
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
	opts := options.FindOne().SetProjection(bson.M{"password": 0, "code": 0})
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}, opts).Decode(userModel)
	if userModel.Rol != 4 {
		return false
	}

	return err == nil
}
