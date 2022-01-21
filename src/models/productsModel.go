package models

import (
	"context"
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Kind         string             `json:"kind" bson:"kind,omitempty"`
	Price        uint               `json:"price" bson:"price,omitempty"`
	PriceMin     uint               `json:"pricemin" bson:"pricemin,omitempty"`
	PriceSale    uint               `json:"pricesale" bson:"pricesale,omitempty"`
	PriceSaleMin uint               `json:"pricesalemin" bson:"pricesalemin,omitempty"`
	PriceSold    uint               `json:"pricesold" bson:"pricesold,omitempty"`
	SoldOut      bool               `json:"soldout" bson:"soldout,omitempty"`
	Seller       string             `json:"seller" bson:"seller,omitempty"`
	Gender       uint8              `json:"gender" bson:"gender,omitempty"`
	Photo        string             `json:"photo" bson:"photo,omitempty"`
	Photos       []string           `json:"photos" bson:"photos,omitempty"`
	Models       []string           `json:"models" bson:"models,omitempty"`
	Pattern      ProductPattern     `json:"pattern" bson:"pattern,omitempty"`
}

/*
price - precio que se muestra al cliente.
pricemin - precio minimo (cuando el cliente insiste una rebaja).
pricesold - prcio a lo que fue vendido.
pricesale - precio de rebaja que se muestra al cliente (esta se establece cuando la prenda se rebaja para que salga mas rapido)
precesalemin - precio de rebaja minimo (cuando el cliente insiste en una rebaja)
soldout - si el producto fue vendido.
seller - vendedor
gender - si el producto para hombres (1), para mujeres (2) o para ambos (3).
*/

func NewProduct(kind string, price, pricemin uint, gender uint8,
	photo string, photos, models []string, talla string, larTorso,
	conPecho, conCintura, conCadera, conSisa, larHombro, larManga uint8) error {
	newProduct := &Product{
		Kind:         kind,
		Price:        price,
		PriceMin:     pricemin,
		PriceSale:    price,
		PriceSaleMin: pricemin,
		PriceSold:    0,
		Seller:       "",
		Gender:       gender,
		Photo:        photo,
		Photos:       photos,
		Models:       models,
		Pattern: ProductPattern{
			Talla:           talla,
			LargoTorso:      larTorso,
			ContornoPecho:   conPecho,
			ContornoCintura: conCintura,
			ContornoCadera:  conCadera,
			ContornoSisa:    conSisa,
			LargoHombro:     larHombro,
			LargoManga:      larManga,
		},
	}

	ctx, client, coll := config.ConnectColl("products")
	defer fmt.Println("Disconnected DB")
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(context.Background(), newProduct)
	return err
}
func GetAllProducts() ([]Product, error) {
	ctx, client, coll := config.ConnectColl("products")
	defer fmt.Println("Disconnected DB")
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
