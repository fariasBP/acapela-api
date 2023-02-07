package models

import (
	"context"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

/*
ID: id que da por defecto mongodb
Name: nombre de la prenda
Creator: id (string) de la tienda creadora
Verification: si la tienda se encuentra verificada es decir sí es oficial(default false)
Status: estado de la ProductKind (aun no especificado)
CreateDate: fecha de la creacion
UpdateDate: fecha de la ultima actualización
*/
type ProductKind struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name,omitempty"`
	Creator      string             `json:"creator" bson:"creator,omitempty"`
	Verification bool               `json:"verification" bson:"verification,omitempty"`
	Status       int                `json:"status" bson:"status,omitempty"`
	Suscriptions int                `json:"suscriptions" bson:"suscriptions,omitempty"`
	CreateDate   time.Time          `json:"create_date" bson:"create_date,omitempty"`
	UpdateDate   time.Time          `json:"update_date" bson:"update_date,omitempty"`
}

func NewProductKind(name, creator string) error {
	newModel := &ProductKind{
		Name:       name,
		Creator:    creator,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}

	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(context.Background(), newModel)
	return err
}

// verifica si existe el nombre del kind (true = existe)
func ExistsNameProductKind(name string) (b bool) {
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)

	kindProd := &ProductKind{}

	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(kindProd)
	b = true
	if err != nil {
		b = false
	}
	return
}
func GetAllKinds() ([]ProductKind, error) {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)
	// consultando
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var data []ProductKind
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, nil
}
func UpdateNameKind(id primitive.ObjectID, name string) error {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)
	// actualizando
	update := bson.M{"$set": bson.M{"name": name}}
	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, update)

	return err
}
func DeleteKindById(id primitive.ObjectID) error {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)
	// eliminando de la BBDD
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}
func ExistKindId(id primitive.ObjectID) bool {
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)
	// consultando
	kindModel := &ProductKind{}
	err := coll.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	return err == nil
}
func ExistKindIdString(id string) bool {
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)
	// verificando si id es correcto
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	// consultando
	kindModel := &ProductKind{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	return err == nil
}

// incrementar suscripcion
func AddKindSuscription(idKind string) error {
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)
	//convirtiendo a ObjectId
	id, err := primitive.ObjectIDFromHex(idKind)
	if err != nil {
		return err
	}
	// actulizando (incrementando)
	update := bson.M{
		"$inc": bson.M{
			"suscriptions": 1,
		},
	}
	_, err = coll.UpdateOne(ctx, bson.M{"_id": id}, update)

	return err
}
