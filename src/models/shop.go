package models

import (
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

/*
ID: id establecido por mongodb
Name: Nombre de la tienda
Owner: id (string) del usuario que creo la tienda
Admins: lista de administradores (usuarios que tienen permisos especiales para administrar la tienda)
Description: descripcion de la tienda
Status: estado de la tienda
CreateDate: fecha de la creacion
UpdateDate: fecha de la ultima actualización
*/
type (
	Shop struct {
		ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name        string             `json:"name" bson:"name,omitempty"`
		Owner       string             `json:"owner" bson:"owner,omitempty"`
		Admins      []Admin            `json:"admins" bson:"admins,omitempty"`
		Description string             `json:"description" bson:"description,omitempty"`
		Photo       string             `json:"photo" bson:"photo,omitempty"`
		Status      int                `json:"status" bson:"status,omitempty"`
		CreateDate  time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate  time.Time          `json:"update_date" bson:"update_date,omitempty"`
	}
	Admin struct {
		IdUser      string   `json:"id_user" bson:"id_user,omitempty"`
		Permissions []string `json:"permissions" bson:"permissions"`
	}
)

const dbShops = "shops"

// crear una tienda
func CreateShop(name, ownerId, description string) error {
	// valores
	newStore := &Shop{
		Name:        name,
		Owner:       ownerId,
		Description: description,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
	}

	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(ctx, newStore)

	return err
}

// obtener tienda
func GetShop(idShop string) (*Shop, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)
	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(idShop)
	if err != nil {
		return nil, err
	}
	// consultando
	// consultando
	shop := &Shop{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(shop)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

// obtener las tiendas
func GetShops(name string, limit, page int) ([]Shop, int64, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)
	// creando parametros consulta
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(limit * (page - 1)))
	query := bson.M{}
	if name != "" {
		query = bson.M{"name": primitive.Regex{
			Pattern: `(\s` + name + `|^` + name + `|\w` + name + `\w` + `|` + name + `$` + `|` + name + `\s)`, Options: "i",
		}}
	}
	// consultando cantidad de datos
	count, err := coll.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	// consultando
	cursor, err := coll.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	// modelando datos
	var shops []Shop
	if err = cursor.All(ctx, &shops); err != nil {
		return nil, 0, err
	}

	return shops, count, nil
}

// obtener las tiendas
func GetShopByOwner(ownerId, name string, limit, page int) ([]Shop, int64, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)
	// creando parametros consulta
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(limit * (page - 1)))
	query := bson.M{"owner": ownerId}
	if name != "" {
		query = bson.M{"$and": []bson.M{
			bson.M{"name": primitive.Regex{
				Pattern: `(\s` + name + `|^` + name + `|\w` + name + `\w` + `|` + name + `$` + `|` + name + `\s)`, Options: "i",
			}},
			bson.M{"owner": ownerId},
		}}
	}
	// consultando cantidad
	count, err := coll.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	// consultando
	cursor, err := coll.Find(ctx, query, opts)
	defer cursor.Close(ctx)
	// modelando datos
	var shops []Shop
	if err = cursor.All(ctx, &shops); err != nil {
		return nil, 0, err
	}

	return shops, count, nil
}

// verificar si el nombre de la tienda ya existe (true = existe)
func ExistsNameShop(name string) bool {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)

	// consultando
	n := &Shop{}
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(n)

	return err == nil
}

// verificar si existe la tienda (true = existe)
func ExistsShopByIdString(idShop string) bool {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)

	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(idShop)
	if err != nil {
		return false
	}

	// consultando
	n := &Shop{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(n)

	return err == nil
}

// volver admins a unos usuarios
func ConvertToAdminShop(idUser, idShop string) error {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)
	// convirtiendo id en ObjectId
	objectId, err := primitive.ObjectIDFromHex(idShop)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"admins": &Admin{
				IdUser:      idUser,
				Permissions: []string{},
			},
		},
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objectId}, update)

	return err
}

// verificar si es el dueño de la tienda (true = es dueño)
func VerifyOwnerShop(idOwner string, idShop string) bool {
	// conenctado a BBDD
	ctx, client, coll := config.ConnectColl(dbShops)
	defer client.Disconnect(ctx)

	// convertir idshop
	objectId, err := primitive.ObjectIDFromHex(idShop)
	if err != nil {
		return false
	}
	// obteneniendo el idOwnder de shop
	shop := &Shop{}
	coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(shop)

	return shop.Owner == idOwner
}

// // ---- registrador de clientes ----
// func ClientRegistrar(name string, phone int, idShop string) error {
// 	// valores de usuario
// 	nUserRegister := &User{
// 		Name:       name,
// 		Rol:        4,
// 		Phone:      phone,
// 		CreateDate: time.Now(),
// 		UpdateDate: time.Now(),
// 	}
// 	nSuscription := &Suscription{
// 		IdShop: idShop,

// 	}
// 	// conectando a BBDD
// 	ctx, client, coll := config.ConnectColl("users")
// 	defer fmt.Println("Disconnected DB")
// 	defer client.Disconnect(ctx)
// 	// insertando
// 	_, err := coll.InsertOne(ctx, nUserRegister)

// 	return err
// }
