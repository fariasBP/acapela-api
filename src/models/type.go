package models

import (
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type TypeProduct struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Creator     string             `json:"creator" bson:"creator,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Used        int                `json:"used" bson:"used,omitempty"`
	CreateDate  time.Time          `json:"create_date" bson:"create_date,omitempty"`
	UpdateDate  time.Time          `json:"update_date" bson:"update_date,omitempty"`
}

// crear type
func CreateType(name, creatorId, description string) error {
	// valores
	newType := &TypeProduct{
		Name:        name,
		Creator:     creatorId,
		Description: description,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
	}

	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_TYPE)
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(ctx, newType)

	return err
}

// obtener types
func GetTypes(name string, limit, page int) ([]TypeProduct, int64, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_TYPE)
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
	var types []TypeProduct
	if err = cursor.All(ctx, &types); err != nil {
		return nil, 0, err
	}

	return types, count, nil
}

// verificar si existe el mismo nombre
func ExistsTypeName(name string) bool {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_TYPE)
	defer client.Disconnect(ctx)

	// consultando
	typeProduct := &TypeProduct{}
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(typeProduct)

	return err == nil
}

// verificar si existe typeId
func ExistsTypeIdString(id string) bool {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_TYPE)
	defer client.Disconnect(ctx)
	// transformar id object
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	// consultando
	typeProduct := &TypeProduct{}
	err = coll.FindOne(ctx, bson.M{"_id": objId}).Decode(typeProduct)

	return err == nil
}

// aumentar el contado de usado
func IncrementUsedType(typeId string) error {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_TYPE)
	defer client.Disconnect(ctx)
	// transformando a idObject
	objId, err := primitive.ObjectIDFromHex(typeId)
	if err != nil {
		return err
	}
	// consultando
	_, err = coll.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$inc": bson.M{"used": 1}})

	return err
}
