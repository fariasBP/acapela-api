package controllers

import (
	"encoding/json"
	"strings"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

type (
	createAdminShop struct {
		IdUser string `json:"id_user"`
	}
)

// crear una tienda
func CreateShop(c echo.Context) error {
	// obteniendo variables
	body := &models.Shop{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// transformando
	name := strings.TrimSpace(body.Name)
	// obteniendo al usuario que creo
	ownerId := c.Get("id").(string)

	// verificando si existe el nombre de la tienda
	b := models.ExistsNameShop(name)
	if b {
		return c.JSON(400, config.SetResError(400, "Error: Nombre de la tienda ya existe", ""))
	}

	// creando la tienda
	if err := models.CreateShop(name, ownerId, body.Description); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se ha creado la tienda", ""))
	}

	return c.JSON(200, config.SetRes(200, "Se creo correctamente la tienda"))
}

func ConvertToAdminShop(c echo.Context) error {
	// obteniendo variables
	body := &createAdminShop{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniendo shop id
	idShop := c.Get("shop").(string)
	// convirtiendo a un usuario en admin
	if err := models.ConvertToAdminShop(body.IdUser, idShop); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo convertir en administrador", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se creo un administrador correctamente"))
}
