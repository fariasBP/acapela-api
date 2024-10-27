package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateSubtype(c echo.Context) error {
	// obteniendo variables
	body := &middlewares.CreateSubtypeParams{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniedo usuario
	idUser := c.Get(config.ID_USER).(string)
	// creando el type
	if err := models.CreateSubtype(body.Name, body.TypeId, idUser, body.Description); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se ha creado el Subtype", err.Error()))
	}
	// incrementar used del type
	if err := models.IncrementUsedType(body.TypeId); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: Se creo el Subtype pero no se incremeto el used del Type", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se creo el subType"))
}

func GetSubtypes(c echo.Context) error {
	// obteniendo valores
	typeId := c.QueryParam("type")
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
	types, count, err := models.GetSubtypes(name, typeId, limitInt, pageInt)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo obtener los SubTypes.", err.Error()))
	}
	return c.JSON(200, config.SetResJsonCount(200, "se obtuvieron los SubTypes", count, types))
}
