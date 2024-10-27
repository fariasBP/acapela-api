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
ID: identificador establecido por mongodb
Nam: nombre de modelo
*/
type (
	ModelProduct struct {
		ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name         string             `json:"name" bson:"name,omitempty"`
		Creator      string             `json:"creator" bson:"creator,omitempty"`
		SubtypeId    string             `json:"subtype" bson:"subtype,omitempty"`
		Photos       []string           `json:"photos" bson:"photos,omitempty"`
		Verification bool               `json:"verification" bson:"verification,omitempty"`
		Description  string             `json:"description" bson:"description,omitempty"`
		Suscriptions int                `json:"suscriptions" bson:"suscriptions,omitempty"`
		CreateDate   time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate   time.Time          `json:"update_date" bson:"update_date,omitempty"`
	}
	MessuresModel struct {
		Name string // aun no esta pensado bien
	}
)

func CreateModelProduct(name, subtypeId, creator, description string) error {
	newModel := &ModelProduct{
		Name:        name,
		SubtypeId:   subtypeId,
		Creator:     creator,
		Description: description,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
	}
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_MODELS)
	defer client.Disconnect(ctx)
	// insertando en la BBDD
	_, err := coll.InsertOne(context.Background(), newModel)

	return err
}

// obtener types
func GetModels(name, subtypeId string, limit, page int) ([]ModelProduct, int64, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_MODELS)
	defer client.Disconnect(ctx)
	// creando parametros consulta
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(limit * (page - 1)))
	query := bson.M{"subtype": subtypeId}
	if name != "" {
		query = bson.M{"$and": []bson.M{
			bson.M{"name": primitive.Regex{
				Pattern: `(\s` + name + `|^` + name + `|\w` + name + `\w` + `|` + name + `$` + `|` + name + `\s)`, Options: "i",
			}},
			bson.M{"subtype": subtypeId},
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
	var models []ModelProduct
	if err = cursor.All(ctx, &models); err != nil {
		return nil, 0, err
	}

	return models, count, nil
}

// verifica si ya existe el nombre del modelo (true = existe)
func ExistsNameProductModel(name string) (b bool) {
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	productModel := &ModelProduct{}
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(productModel)
	b = true
	if err != nil {
		b = false
	}
	return
}

// verificar si existe model por su id
func ExistsModelId(id primitive.ObjectID) bool {
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	model := &ModelProduct{}
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(model)

	return err == nil
}

// verificar si existe model por su id string
func ExistsModelIdString(id string) bool {
	ctx, client, coll := config.ConnectColl("models")
	defer client.Disconnect(ctx)

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}

	model := &ModelProduct{}
	err = coll.FindOne(ctx, bson.M{"_id": ObjId}).Decode(model)

	return err == nil
}

// func UpdateModelById(id primitive.ObjectID, name, idKind string) error {
// 	// conectando a la BBDD
// 	ctx, client, coll := config.ConnectColl("models")
// 	defer client.Disconnect(ctx)

// 	update := bson.M{"$set": bson.M{"name": name, "kind": idKind}}
// 	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, update)

// 	return err
// }
// func DeleteModelById(id primitive.ObjectID) error {
// 	// conectando a la BBDD
// 	ctx, client, coll := config.ConnectColl("models")
// 	defer client.Disconnect(ctx)

// 	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

// 	return err
// }
