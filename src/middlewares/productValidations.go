package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreateProductParams struct {
		//basico
		Models []string `json:"models" validate:"required"`
		//cantidad (si es unico o por lote)
		/* Quantity uint `json:"quantity" validate:"required"` */
		//precio
		Price    uint `json:"price" validate:"required,number"`
		PriceMin uint `json:"price_min" validate:"required,number"`
		//fotos y apariencias
		Photos []string `json:"photos" validate:"required"`
		/* Colors []string `json:"colors" validate:"required"` */
		//genero
		Gender models.GenderVal `json:"gender" validate:"required,number"`
		//tamaño
		Size []models.SizeVal `json:"size" validate:"required"`
	}
	sellProductValidations struct {
		ID        primitive.ObjectID `json:"id" validate:"required"`
		SellPrice int                `json:"sell_price" validate:"required,number"`
		Seller    string             `json:"seller" validate:"required"`
	}
)

type MySizes struct {
}

func CreateProductValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &CreateProductParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &CreateProductParams{
			Models:   body.Models,
			Price:    body.Price,
			PriceMin: body.PriceMin,
			Photos:   body.Photos,
			Gender:   body.Gender,
			Size:     body.Size,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// validando modelos
		// validadndo photos
		// validando gender
		// validadndo size
		sizes := []models.SizeVal{models.S, models.M, models.L, models.XL}

		for _, v := range body.Size {

		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return c.String(200, "ok se interrumpio")
		return next(c)
	}
}
func SellProductValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.Product{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &sellProductValidations{
			ID:        body.ID,
			SellPrice: body.SellPrice,
			Seller:    body.Seller,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}
