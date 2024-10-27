package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateType(c echo.Context) error {
	// obteniendo variables
	body := &middlewares.CreateTypeParams{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	//obteniendo id del ususario
	idUser := c.Get(config.ID_USER).(string)

	// creando el type
	if err := models.CreateType(body.Name, idUser, body.Description); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se ha creado el type", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se creo el type"))
}

func GetTypes(c echo.Context) error {
	// obteniendo valores
	name := c.QueryParam("name")
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	// convirtiendo valores
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	// consultando con la bbdd
	types, count, err := models.GetTypes(name, limitInt, pageInt)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo obtener los types.", err.Error()))
	}

	return c.JSON(200, config.SetResJsonCount(200, "se obtuvieron los types", count, types))
}
