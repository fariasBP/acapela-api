package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	CreateSubtypeParams struct {
		Name        string `json:"name" validate:"required,min=3,startsnotwith= ,endsnotwith= ,alpha,lowercase"`
		TypeId      string `json:"type" validate:"required"`
		Description string `json:"description"`
	}
	GetSubtypesParams struct {
		TypeId string `json:"type" validate:"required"`
		Name   string `json:"name"`
		Limit  string `json:"limit"`
		Page   string `json:"page"`
	}
)

// validador para typeCreate
func CreateSubtypeValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &CreateSubtypeParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &CreateSubtypeParams{
			Name:   body.Name,
			TypeId: body.TypeId,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// verificando si existe un typeId
		if b := models.ExistsTypeIdString(body.TypeId); !b {
			return c.JSON(400, config.SetResError(400, "Error: No existe el typeId.", ""))
		}
		// verficando que no exista un nombre igual
		if b := models.ExistsSubtypeName(body.Name); b {
			return c.JSON(400, config.SetResError(400, "Error: Ya existe un SubtypeProduct con ese nombre.", ""))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}

// validador para getsubtypes
func GetSubtypesValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo valores
		typeId := c.QueryParam("type")
		// estableciendo los argumentos de validacion
		v := &GetSubtypesParams{
			TypeId: typeId,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// verificando que existe el shopid
		if b := models.ExistsTypeIdString(typeId); !b {
			return c.JSON(400, config.SetResError(400, "Error: No existe typeId.", ""))
		}
		// fin del middleware
		return next(c)
	}
}
