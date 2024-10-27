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
	CreateShopParams struct {
		Name        string `json:"name" validate:"required,min=3,startsnotwith= ,endsnotwith= "`
		Description string `json:"description"`
	}
	GetShopParams struct {
		Id string `json:"id" validate:"required"`
	}
)

func GetShopValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obtengo valores
		id := c.QueryParam("id")
		// estableciendo los argumentos de validacion
		v := &GetShopParams{
			Id: id,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// fin del middleware
		return next(c)
	}
}
func CreateShopValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &models.Shop{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &CreateShopParams{
			Name: body.Name,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Valores invalidos.", err.Error()))
		}
		// verificando si existe el nombre de la tienda
		b := models.ExistsNameShop(body.Name)
		if b {
			return c.JSON(400, config.SetResError(400, "Error: Nombre de la tienda ya existe", ""))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}
