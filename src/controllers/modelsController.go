package controllers

import (
	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateModel(c echo.Context) error {
	// extrayendo values
	name := c.FormValue("name")
	idKind := c.FormValue("kind")
	// verificando si exsite en la BBDD
	exist := models.ExistsNameProductModel(name)
	if exist {
		return c.JSON(400, config.SetResError(400, "Error: name productModel is already created.", ""))
	}
	exist = models.VerifyKindId(idKind)
	if !exist {
		return c.JSON(400, config.SetResError(400, "Error: id Kind does not exist.", ""))
	}
	// creando nuevo modelo en la BBDD
	err := models.NewProductModel(name, idKind)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not created new product model", err.Error()))
	}
	return c.JSON(200, config.SetRes(200, "product model created"))
}
func GetAllModels(c echo.Context) error {
	models, err := models.GetAllModels()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not get all models", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "get all models successful", models))
}
