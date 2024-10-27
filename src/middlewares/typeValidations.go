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
	CreateTypeParams struct {
		Name        string `json:"name" validate:"required,min=3,startsnotwith= ,endsnotwith= ,alpha,lowercase"`
		Description string `json:"description" validate:"required"`
	}
	GetTypeParams struct {
		Name  string `json:"name"`
		Limit string `json:"limit"`
		Page  string `json:"page"`
	}
)

// validador para typeCreate
func CreateTypeValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &CreateTypeParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &CreateTypeParams{
			Name:        body.Name,
			Description: body.Description,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// verificando que no existe type'name
		if b := models.ExistsTypeName(body.Name); b {
			return c.JSON(400, config.SetResError(400, "Error: Ya existe un TypeProduct con ese nombre.", ""))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}
