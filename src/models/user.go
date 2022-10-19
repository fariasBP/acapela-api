package models

import (
	"fmt"
	"os"
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
		ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name               string             `json:"name" bson:"name,omitempty"`
		Lastname           string             `json:"lastname" bson:"lastname,omitempty"`
		Photo              string             `json:"photo" bson:"photo,omitempty"`
		Email              string             `json:"email" bson:"email,omitempty"`
		Code               string             `json:"code" bson:"code,omitempty"`
		Rol                int                `json:"rol" bson:"rol,omitempty"`
		Phone              int                `json:"phone" bson:"phone"`
		Sleep              uint8              `json:"sleep" bson:"sleep"`
		CodeDate           time.Time          `json:"code_date" bson:"code_date,omitempty"`
		SleepDate          time.Time          `json:"sleep_date" bson:"sleep_date,omitempty"`
		WpRegistrationDate time.Time          `json:"wp_registration_date" bson:"wp_registration_date,omitempty"`
		Mailbox            string             `json:"mailbox" bson:"mailbox,omitempty"`
		CountMailbox       int                `json:"count_mailbox" bson:"count_mailbox,omitempty"`
		CreateDate         time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate         time.Time          `json:"update_date" bson:"update_date,omitempty"`
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
	// consulatando
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"phone": phone}).Decode(user)
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
func GetUserAndVerifyNotblockExitsAndActive(phone int) (bool, bool, bool, *User, error) {
	// variable entorno
	appName, _ := os.LookupEnv("APP_NAME")
	// Conectando a la BBDD
	ctx, client, db := config.ConnectDB()
	collApp := db.Collection("app")
	collUsers := db.Collection("users")
	defer fmt.Println("Disconnected DB")
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

// ---- obtener numero y nombre para notificaciones ----
func GetPhoneAndNameForNotificationsFromClients() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetProjection(bson.M{"name": 1, "phone": 1})
	filter := bson.M{"$and": []bson.M{
		bson.M{"sleep": 0},
		bson.M{"rol": 4},
	}}
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

	return err == nil
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
func VerifyActiveUserByPhone(phone int) bool {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// verificar si el usuario esta activo
	user := &User{}
	err := coll.FindOne(ctx, bson.M{"phone": phone}).Decode(user)
	if err != nil {
		return false
	}
	if (user.SleepDate == time.Time{} || user.Sleep == 0) {
		return true
	}
	return false
}

// ---- ACTUALIZAR USUARIO ----
// ---- actualizar nombre de usuario ----
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

// ---- inactivar usuario ----
func UpdInactiveUserByPhone(phone int, sleep uint8) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
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
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
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

// ---- ELIMINAR USUARIO ----
func DelUserByPhone(phone int) error {
	// Conectandose a la DDBB
	ctx, client, coll := config.ConnectColl("users")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)
	// consultando
	_, err := coll.DeleteOne(ctx, bson.M{"phone": phone})

	return err
}
