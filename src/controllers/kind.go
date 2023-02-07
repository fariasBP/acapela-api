package controllers

import (
	"encoding/json"
	"strings"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

// creando tipo de prenda
func CreateKind(c echo.Context) error {
	// obteniendo variables
	body := &models.ProductKind{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()

	// obteniendo shop id
	idShop := c.Get("shop").(string)

	// // verificar que no existe un kind con el mismo nombre
	// exist := models.ExistsNameProductKind(strings.TrimSpace(body.Name))
	// if exist {
	// 	return c.JSON(400, config.SetResError(400, "Error: Ya existe un KindProduct con el mismo nombre", "same values kindproduct name"))
	// }

	// creadndo el nuevo kind en la BBDD
	err := models.NewProductKind(strings.TrimSpace(body.Name), idShop)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No creado un KindProduct", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "kindProduct creado"))
}
func GetAllKinds(c echo.Context) error {
	models, err := models.GetAllKinds()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not get all kinds", err.Error()))
	}
	return c.JSON(200, config.SetResJson(200, "get all kinds successful", models))
}
func UpdateNameKind(c echo.Context) error {
	// obteniendo variables
	body := &models.ProductKind{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe el kindID
	exist := models.ExistKindId(body.ID)
	if !exist {
		return c.JSON(400, config.SetRes(400, "Error: No existe el Kind con el Id proporcionado"))
	}
	// vericando que no existe el mismo nombre
	exist = models.ExistsNameProductKind(strings.TrimSpace(body.Name))
	if exist {
		return c.JSON(400, config.SetRes(400, "Error: Existe otro Kind con el mismo nombre"))
	}
	// actualizando name en BBDD
	err := models.UpdateNameKind(body.ID, strings.TrimSpace(body.Name))
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo acutalizar el nombre de ProductKind", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se actualizo con exito"))
}
func DeleteKind(c echo.Context) error {
	// obteniendo variables
	body := &models.ProductKind{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe el kindID
	exist := models.ExistKindId(body.ID)
	if !exist {
		return c.JSON(400, config.SetRes(400, "Error: No existe el Kind con el Id proporcionado"))
	}
	// eliminando kind de BBDD
	err := models.DeleteKindById(body.ID)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo elimina el ProductKind", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se elimino con exito"))
}
