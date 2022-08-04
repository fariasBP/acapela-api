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
	createProductValidations struct {
		Price           int      `json:"price" validate:"required,number"`
		PriceMin        int      `json:"price_min", validate:"required,number"`
		Photos          []string `json:"photos" validate:"required"`
		Kind            string   `json:"kind" validate:"required"`
		Models          []string `json:"models" validate:"required"`
		Gender          int      `json:"gender" validate:"required,number"`
		Size            []int    `json:"size" validate:"required"`
		Modelquality    int      `json:"model_quality" validate:"required,number"`
		Materialquality int      `json:"material_quality" validate:"required,number"`
	}
	sellProductValidations struct {
		ID        primitive.ObjectID `json:"id" validate:"required"`
		SellPrice int                `json:"sell_price" validate:"required,number"`
		Seller    string             `json:"seller" validate:"required"`
	}
)

func CreateProductValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.Product{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &createProductValidations{
			Price:           body.Price,
			PriceMin:        body.PriceMin,
			Photos:          body.Photos,
			Kind:            body.Kind,
			Models:          body.Models,
			Gender:          body.Gender,
			Size:            body.Size,
			Modelquality:    body.ModelQuality,
			Materialquality: body.MaterialQuality,
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
