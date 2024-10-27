package models

import (
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Subtype struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	TypeId      string             `json:"type" bson:"type,omitempty"`
	Creator     string             `json:"creator" bson:"creator,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	CreateDate  time.Time          `json:"create_date" bson:"create_date,omitempty"`
	UpdateDate  time.Time          `json:"update_date" bson:"update_date,omitempty"`
}

// crear una tienda
func CreateSubtype(name, typeId, creatorId, description string) error {
	// valores
	newType := &Subtype{
		Name:        name,
		TypeId:      typeId,
		Creator:     creatorId,
		Description: description,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
	}

	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_SUBTYPE)
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(ctx, newType)

	return err
}

// obtener subTypes
func GetSubtypes(name, typeId string, limit, page int) ([]Subtype, int64, error) {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_SUBTYPE)
	defer client.Disconnect(ctx)
	// creando parametros consulta
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(limit * (page - 1)))
	query := bson.M{"type": typeId}
	if name != "" {
		query = bson.M{"$and": []bson.M{
			bson.M{"name": primitive.Regex{
				Pattern: `(\s` + name + `|^` + name + `|\w` + name + `\w` + `|` + name + `$` + `|` + name + `\s)`, Options: "i",
			}},
			bson.M{"type": typeId},
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
	var subtypes []Subtype
	if err = cursor.All(ctx, &subtypes); err != nil {
		return nil, 0, err
	}

	return subtypes, count, nil
}

// verificar si existe el nombre del subtype
func ExistsSubtypeName(name string) bool {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_SUBTYPE)
	defer client.Disconnect(ctx)
	// consultando
	subtype := &Subtype{}
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(subtype)

	return err == nil
}

// verificar si existe el subtype id
func ExistsSubtypeIdString(id string) bool {
	// conectando
	ctx, client, coll := config.ConnectColl(config.DB_SUBTYPE)
	defer client.Disconnect(ctx)
	// convirtiendo id a object
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	// consultando
	kindModel := &KindProduct{}
	err = coll.FindOne(ctx, bson.M{"_id": objId}, options.FindOne().SetProjection(bson.M{"name": 0})).Decode(kindModel)

	return err == nil
}

// aumentar el contado de usado
func IncrementUsedSubtype(subtypeId string) error {
	// conectando a la BBDD
	ctx, client, coll := config.ConnectColl(config.DB_SUBTYPE)
	defer client.Disconnect(ctx)
	// transformando a idObject
	objId, err := primitive.ObjectIDFromHex(subtypeId)
	if err != nil {
		return err
	}
	// consultando
	_, err = coll.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$inc": bson.M{"used": 1}})

	return err
}
