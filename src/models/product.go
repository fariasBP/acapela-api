package models

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type (
	Product struct {
		ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Price           int                `json:"price" bson:"price,omitempty"`
		PriceMin        int                `json:"price_min" bson:"price_min,omitempty"`
		Photos          []string           `json:"photos" bson:"photos,omitempty"`
		Kind            string             `json:"kind" bson:"kind,omitempty"`
		Models          []string           `json:"models" bson:"models,omitempty"`
		Gender          int                `json:"gender" bson:"gender,omitempty"`
		Size            []string           `json:"size" bson:"size,omitempty"`
		ModelQuality    int                `json:"model_quality" bson:"model_quality,omitempty"`
		MaterialQuality int                `json:"material_quality" bson:"material_quality,omitempty"`
		NewPrice        int                `json:"new_price" bson:"new_price,omitempty"`
		SellPrice       int                `json:"sell_price" bson:"sell_price,omitempty"`
		Seller          string             `json:"seller" bson:"seller,omitempty"`
		Buyer           string             `json:"buyer" bson:"buyer,omitempty"`
		Day             int                `json:"day" bson:"day,omitempty"`
		Month           time.Month         `json:"month" bson:"month,omitempty"`
		Year            int                `json:"year" bson:"year,omitempty"`
		CreateDate      time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate      time.Time          `json:"update_date" bson:"update_date,omitempty"`
	}
	Created struct {
		LastCreated time.Time `json:"last_created" bson:"last_created,omitempty"`
	}
)

/*
price - precio que se muestra al cliente.
pricemin - precio minimo (cuando el cliente insiste una rebaja).
pricesold - prcio a lo que fue vendido.
pricesale - precio de rebaja que se muestra al cliente (esta se establece cuando la prenda se rebaja para que salga mas rapido)
precesalemin - precio de rebaja minimo (cuando el cliente insiste en una rebaja)
soldout - si el producto fue vendido.
seller - vendedor
gender - si el producto para hombres (1), para mujeres (2) o para ambos (3).
modelquality - calidad del modelo
	0 especial
	1 primera
	2 segunda
	3 liquidacion
	4 rezago
materialquality - calidad del material ej./ 1/10 3/10, 10/10
*/

func NewProduct(price, priceMin int, photos []string, kind string, models []string, gender int, size []string, modelQuality, materialQuality int) error {
	newProduct := &Product{
		Price:           price,
		PriceMin:        priceMin,
		Photos:          photos,
		Kind:            kind,
		Models:          models,
		Gender:          gender,
		Size:            size,
		ModelQuality:    modelQuality,
		MaterialQuality: materialQuality,
		Day:             time.Now().UTC().Day(),
		Month:           time.Now().UTC().Month(),
		Year:            time.Now().UTC().Year(),
		CreateDate:      time.Now().UTC(),
		UpdateDate:      time.Now().UTC(),
	}

	ctx, client, db := config.ConnectDB()
	collProducts := db.Collection("products")
	collApp := db.Collection("app")
	defer client.Disconnect(ctx)

	appName, _ := os.LookupEnv("APP_NAME")

	_, err := collProducts.InsertOne(ctx, newProduct)
	if err != nil {
		return err
	}

	_, err = collApp.UpdateOne(ctx, bson.M{"name": appName}, bson.M{"$set": bson.M{"set_products": time.Now().UTC()}})

	return err
}
func GetAllProducts() ([]Product, error) {
	ctx, client, coll := config.ConnectColl("products")
	defer client.Disconnect(ctx)

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func GetProducts(limit int, page int) ([]Product, error) {
	// conectando a la base de datos
	ctx, client, coll := config.ConnectColl("products")
	defer client.Disconnect(ctx)
	// obteniendo los ultimos abrigos
	opts := options.Find().SetSort(bson.M{"create_date": -1}).SetLimit(int64(limit)).SetSkip(int64(page))
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func SellProductWithoutBuyer(id primitive.ObjectID, sellprice int, seller string) error {
	ctx, client, coll := config.ConnectColl("products")
	defer client.Disconnect(ctx)
	// realizando consulta
	update := bson.M{"$set": bson.M{"sellprice": sellprice, "seller": seller}}
	_, err := coll.UpdateOne(context.Background(), bson.M{"_id": id}, update)

	return err
}
func SellProductWithBuyer(id primitive.ObjectID, sellprice int, seller, buyer string) error {
	ctx, client, coll := config.ConnectColl("products")
	defer client.Disconnect(ctx)
	// realizando consulta
	update := bson.M{"$set": bson.M{"sellprice": sellprice, "seller": seller, "buyer": buyer}}
	_, err := coll.UpdateOne(context.Background(), bson.M{"_id": id}, update)

	return err
}
func GetNewProducts() ([]Product, error) {
	// conectando a BBDD
	ctx, client, db := config.ConnectDB()
	collProducts := db.Collection("products")
	collApp := db.Collection("app")
	defer client.Disconnect(ctx)

	appName, _ := os.LookupEnv("APP_NAME")
	// obteniendo el ultimo productDate registrado
	app := &App{}
	err := collApp.FindOne(ctx, bson.M{"name": appName}).Decode(app)
	if err != nil {
		return nil, err
	}
	lastYear, lastMonth, lastDay := app.SetProducts.Date()
	// consultando
	filter := bson.M{"$and": []bson.M{
		bson.M{"day": lastDay},
		bson.M{"month": lastMonth},
		bson.M{"year": lastYear},
	}}
	cursor, err := collProducts.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	if len(products) > 3 {
		return products, nil
	}
	// enviando los 30 ultimos productos (-1 = descendente)
	opts := options.Find().SetSort(bson.M{"create_date": -1}).SetLimit(30).SetSkip(0)
	cursor, err = collProducts.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// VERIFICANDO SI EXISTE ID
func ExistProductId(id primitive.ObjectID) bool {
	fmt.Println("aqui")
	fmt.Println(id)
	ctx, client, coll := config.ConnectColl("products")
	defer client.Disconnect(ctx)
	// consultando
	productModel := &ProductModel{}
	err := coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(productModel)

	return err == nil
}

func NewCreated() error {
	newCreated := &Created{
		LastCreated: time.Now(),
	}

	ctx, client, coll := config.ConnectColl("created")
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(context.Background(), newCreated)
	return err
}

func GetLastDateFromCreated() (int, time.Month, int, error) {
	ctx, client, coll := config.ConnectColl("created")
	defer client.Disconnect(ctx)
	// consultando (obteniendo el ultimo last created)
	opts := options.Find().SetSort(bson.M{"last_created": -1}).SetLimit(1)
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, 0, 0, err
	}
	defer cursor.Close(ctx)
	var createds []Created
	if err = cursor.All(ctx, &createds); err != nil {
		return 0, 0, 0, err
	}
	if len(createds) == 0 {
		return 0, 0, 0, errors.New("length of createds is 0")
	}
	year, month, day := createds[0].LastCreated.Date()

	return year, month, day, nil
}
func GetLastDateFromProduct() (int, time.Month, int, error) {
	// conectando con BBDD
	ctx, client, coll := config.ConnectColl("products")
	defer client.Disconnect(ctx)
	// consultando
	opts := options.Find().SetSort(bson.M{"create_date": -1}).SetLimit(1)
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, 0, 0, err
	}
	defer cursor.Close(ctx)
	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return 0, 0, 0, err
	}
	year, month, day := products[0].CreateDate.Date()

	return year, month, day, nil
}

func IsEqualLastDateCreated() bool {

	ctx, client, coll := config.ConnectColl("created")
	defer client.Disconnect(ctx)
	// consultando (obteniendo el ultimo last created)
	opts := options.Find().SetSort(bson.M{"last_created": -1}).SetLimit(1)
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return false
	}
	defer cursor.Close(ctx)
	var createds []Created
	if err = cursor.All(ctx, &createds); err != nil {
		return false
	}
	if len(createds) == 0 {
		return false
	}

	ctx, client, coll = config.ConnectColl("products")
	defer client.Disconnect(ctx)

	opts = options.Find().SetSort(bson.M{"create_date": -1}).SetLimit(1)
	cursor, err = coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return false
	}
	defer cursor.Close(ctx)
	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return false
	}

	year, month, day := createds[0].LastCreated.Date()
	y, m, d := products[0].CreateDate.Date()

	fmt.Println("EL ULTIMO:")
	fmt.Println(year, month, day)
	fmt.Println("AHORA:")
	fmt.Println(y, m, d)

	if y == year && m == month && d == day {
		return true
	}

	return false
}
