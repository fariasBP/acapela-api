package controllers

import (
	"encoding/json"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreatePermission(c echo.Context) error {
	// obteniendo variables
	body := &models.Permission{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()

	//verificando que no exista una ruta identica
	if b := models.ExistsPermisionRoute(body.Route); b {
		return c.JSON(400, config.SetResError(400, "Error: La ruta del permiso ya existe.", ""))
	}

	// creando nuevo permiso
	if err := models.CreatePermission(body.Name, body.Route); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se creo el permiso", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "se creo el permiso"))
}
