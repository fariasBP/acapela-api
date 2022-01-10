package controllers

import (
	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateModel(c echo.Context) error {
	name := c.FormValue("name")

	existE := models.ExistsNameProductModel(name)
	if existE {
		return c.JSON(400, config.SetResError(400, "Error: name product model is already registered.", ""))
	}

	err := models.NewProductModel(name)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not created new product model", err.Error()))
	}
	return c.JSON(200, config.SetRes(200, "product model created"))
}
func GetAllModels(c echo.Context) error {
	models, err := models.GetAllModels()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error get all models", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "todo ok", models))
}
