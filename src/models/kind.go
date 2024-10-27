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
type KindProduct struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name,omitempty"`
	TypeId       string             `json:"type_id" bson:"type_id,omitempty"`
	Creator      string             `json:"creator" bson:"creator,omitempty"`
	Verification bool               `json:"verification" bson:"verification,omitempty"`
	Status       int                `json:"status" bson:"status,omitempty"`
	Suscriptions int                `json:"suscriptions" bson:"suscriptions,omitempty"`
	CreateDate   time.Time          `json:"create_date" bson:"create_date,omitempty"`
	UpdateDate   time.Time          `json:"update_date" bson:"update_date,omitempty"`
}

const KindsDB = "kinds"

func NewKindProduct(name, typeId, creator string) error {
	newModel := &KindProduct{
		Name:       name,
		TypeId:     typeId,
		Creator:    creator,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
	}

	ctx, client, coll := config.ConnectColl(KindsDB)
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(context.Background(), newModel)
	return err
}

// verifica si existe el nombre del kind (true = existe)
func ExistsNameKindProduct(name string) (b bool) {
	ctx, client, coll := config.ConnectColl("kinds")
	defer client.Disconnect(ctx)

	kindProd := &KindProduct{}

	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(kindProd)
	b = true
	if err != nil {
		b = false
	}
	return
}

// obtener kind
func GetKinds(name string, limit, page int) ([]KindProduct, int64, error) {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl(KindsDB)
	defer client.Disconnect(ctx)
	// creando parametros consulta
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(page - 1))
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
	var kinds []KindProduct
	if err = cursor.All(ctx, &kinds); err != nil {
		return nil, 0, err
	}

	return kinds, count, nil
}

// verififica si existe el kind por su Id
func ExistKindId(id primitive.ObjectID) bool {
	// conectando
	ctx, client, coll := config.ConnectColl(KindsDB)
	defer client.Disconnect(ctx)
	// consultando
	kindModel := &KindProduct{}
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
	kindModel := &KindProduct{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	return err == nil
}

// vericacion del type de un kind
func VerifyTypeIdFromKind(typeId, kindId string) bool {
	ctx, client, coll := config.ConnectColl(KindsDB)
	defer client.Disconnect(ctx)
	// convirtiendo a objectID
	ObjId, err := primitive.ObjectIDFromHex(kindId)
	if err != nil {
		return false
	}
	// consultando
	query := bson.M{"$and": []bson.M{
		bson.M{"_id": ObjId},
		bson.M{"type": typeId},
	}}
	kind := &KindProduct{}
	err = coll.FindOne(ctx, query).Decode(kind)

	return err == nil
}

// obetner el type de un kind
func GetTypeIdFromKind(kindId string) (string, error) {
	// conectado a BBDD
	ctx, client, coll := config.ConnectColl(KindsDB)
	defer client.Disconnect(ctx)
	// convirtiendo a objectID
	ObjId, err := primitive.ObjectIDFromHex(kindId)
	if err != nil {
		return "", err
	}
	// consultando
	kind := &KindProduct{}
	if err := coll.FindOne(ctx, bson.M{"_id": ObjId}).Decode(kind); err != nil {
		return "", err
	}

	return kind.TypeId, nil
}

// func UpdateNameKind(id primitive.ObjectID, name string) error {
// 	// conectado a BBDD
// 	ctx, client, coll := config.ConnectColl("kinds")
// 	defer client.Disconnect(ctx)
// 	// actualizando
// 	update := bson.M{"$set": bson.M{"name": name}}
// 	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, update)

// 	return err
// }
// func DeleteKindById(id primitive.ObjectID) error {
// 	// conectado a BBDD
// 	ctx, client, coll := config.ConnectColl("kinds")
// 	defer client.Disconnect(ctx)
// 	// eliminando de la BBDD
// 	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

//		return err
//	}
// // incrementar suscripcion
// func AddKindSuscription(idKind string) error {
// 	ctx, client, coll := config.ConnectColl("kinds")
// 	defer client.Disconnect(ctx)
// 	//convirtiendo a ObjectId
// 	id, err := primitive.ObjectIDFromHex(idKind)
// 	if err != nil {
// 		return err
// 	}
// 	// actulizando (incrementando)
// 	update := bson.M{
// 		"$inc": bson.M{
// 			"suscriptions": 1,
// 		},
// 	}
// 	_, err = coll.UpdateOne(ctx, bson.M{"_id": id}, update)

// 	return err
// }
