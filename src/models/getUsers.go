package models

import (
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// obtener usuarios
func GetUsers() ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
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
func GetPhoneAndNameForNotificationsFromClientsByNotReaded(notReaded int) ([]User, error) {
	// Conectando a la BBDD
	ctx, client, coll := config.ConnectColl("users")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetProjection(bson.M{"name": 1, "phone": 1})
	/*Si el valor de notReaded es 0 entonces no importa si el mensaje
	no ha sido entregado o leido, */
	var filter bson.M
	if notReaded == 0 {
		filter = bson.M{"$and": []bson.M{
			bson.M{"sleep": 0},
			bson.M{"rol": 4},
		}}
	} else {
		filter = bson.M{"$and": []bson.M{
			bson.M{"sleep": 0},
			bson.M{"rol": 4},
			bson.M{"not_readed": bson.M{"$gte": notReaded}},
		}}
	}

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
