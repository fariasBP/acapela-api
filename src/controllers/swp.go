package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func DesactivingUserSwp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// actualizando usuario
	i, err := models.UpdNotReadedUserByPhone(body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo actualizar al usuario", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se desactivo al usuario "+strconv.Itoa(body.Phone)+" ("+strconv.Itoa(i)+")."))
}

func ActivingUserSwp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// actualizando usuario
	err := models.UpdReadedUserByPhone(body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo actualizar al usuario", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se activo al usuario "+strconv.Itoa(body.Phone)))
}
