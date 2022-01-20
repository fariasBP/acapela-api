package controllers

import (
	"encoding/json"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateKind(c echo.Context) error {
	// extraendo name y convirtiendo en slice
	name := c.FormValue("name")
	nameArr := make([]string, 0)
	err := json.Unmarshal([]byte(name), &nameArr)
	if err != nil {
		return c.JSON(400, config.SetResError(500, "error: name field is not slice", err.Error()))
	}
	// creadndo el nuevo kind en la BBDD
	err = models.NewKindProduct(nameArr)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not created new kindProduct", err.Error()))
	}
	return c.JSON(200, config.SetRes(200, "kindProduct created"))
}
func GetAllKinds(c echo.Context) error {
	models, err := models.GetAllKinds()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not get all kinds", err.Error()))
	}
	return c.JSON(200, config.SetResJson(200, "get all kinds successful", models))
}
