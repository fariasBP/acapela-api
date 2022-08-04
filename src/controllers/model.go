package controllers

import (
	"encoding/json"
	"strings"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateModel(c echo.Context) error {
	// obteniendo variables
	body := &models.ProductModel{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando si exsite en la BBDD
	exist := models.ExistsNameProductModel(strings.TrimSpace(body.Name))
	if exist {
		return c.JSON(400, config.SetResError(400, "Error: name productModel is already created.", ""))
	}
	exist = models.ExistKindIdString(body.Kind)
	if !exist {
		return c.JSON(400, config.SetResError(400, "Error: id Kind does not exist.", ""))
	}
	// creando nuevo modelo en la BBDD
	err := models.NewProductModel(strings.TrimSpace(body.Name), body.Kind)
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
func UpdateModel(c echo.Context) error {
	// obteniendo variables
	body := &models.ProductModel{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificar si existe ID model
	exist := models.ExistsModelId(body.ID)
	if !exist {
		return c.JSON(400, config.SetRes(400, "Error: No existe el ID model proporcionado."))
	}
	// verificando que no existe name
	exist = models.ExistsNameProductModel(strings.TrimSpace(body.Name))
	if exist {
		return c.JSON(400, config.SetRes(400, "Error: ya existe el nombre."))
	}
	// verificar que existe kindID
	exist = models.ExistKindIdString(body.Kind)
	if !exist {
		return c.JSON(400, config.SetRes(400, "Error: ID kind no existe."))
	}
	// actualizando model
	err := models.UpdateModelById(body.ID, strings.TrimSpace(body.Name), body.Kind)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo actualizar", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "actualizacion del modelo exitoso"))
}

func DeleteModel(c echo.Context) error {
	// obteniendo variables
	body := &models.ProductModel{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// consultamos a la BBDD
	err := models.DeleteModelById(body.ID)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "no se pudo eliminar el model", err.Error()))

	}
	return c.JSON(200, config.SetRes(200, "se elimino el model exitosamente"))

}
