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
	CreateModelParams struct {
		Name        string `json:"name" validate:"required,startsnotwith= ,endsnotwith= ,lowercase,min=3"`
		SubtypeId   string `json:"subtype" validate:"required"`
		Description string `json:"description"`
	}
	GetModelsParams struct {
		SubtypeId string `json:"subtype" validate:"required"`
		Name      string `json:"name"`
		Limit     string `json:"limit"`
		Page      string `json:"page"`
	}
	modelUpdateValidations struct {
		ID   primitive.ObjectID `json:"id" validate:"required"`
		Name string             `json:"name" validate:"required,lowercase,min=3"`
		Kind string             `json:"kind" validate:"required"`
	}
	modelDeleteValidations struct {
		ID primitive.ObjectID `json:"id" validate:"required"`
	}
)

func CreateModelValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &CreateModelParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &CreateModelParams{
			Name:      body.Name,
			SubtypeId: body.SubtypeId,
		}
		// realizando validacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// verificando si existe subtypeId
		exist := models.ExistsSubtypeIdString(body.SubtypeId)
		if !exist {
			return c.JSON(400, config.SetResError(400, "Error: no existe el SubtypeId.", ""))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}

// validaciones de get models
func GetModelsValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo valores
		subtypeId := c.QueryParam("subtype")
		// estableciendo los argumentos de validacion
		v := &GetModelsParams{
			SubtypeId: subtypeId,
		}
		// realizando validacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// fin del middleware
		return next(c)
	}
}

func ModelUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.ModelProduct{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &modelUpdateValidations{
			ID:   body.ID,
			Name: strings.TrimSpace(body.Name),
			Kind: body.SubtypeId,
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
func ModelDeleteValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.ModelProduct{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &modelDeleteValidations{
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
