package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	kindValidations struct {
		Name string `json:"name" validate:"required,lowercase,min=3"`
	}
	kindUpdateNameValidations struct {
		ID   primitive.ObjectID `json:"id" validate:"required"`
		Name string             `json:"name" validate:"required,min=3"`
	}
	kindDeleteValidations struct {
		ID primitive.ObjectID `json:"id" validate:"required"`
	}
)

func KindValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.ProductKind{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &kindValidations{
			Name: strings.TrimSpace(body.Name),
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
func KindUpdateNameValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.ProductKind{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &kindUpdateNameValidations{
			ID:   body.ID,
			Name: strings.TrimSpace(body.Name),
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
func KindDeleteValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.ProductKind{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &kindDeleteValidations{
			ID: body.ID,
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
