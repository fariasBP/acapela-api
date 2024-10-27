package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

// creando tipo de prenda (Ej.(para ropa) abrigo, pantalones, blizer)
func CreateKind(c echo.Context) error {
	// obteniendo variables
	body := &models.KindProduct{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniendo shop id
	idShop := c.Get(config.Shop).(string)
	id := c.Get(config.Id).(string)
	// verificando que existe la shop y si es dueño
	if b := models.VerifyOwnerShop(id, idShop); !b {
		return c.JSON(400, config.SetResError(400, "Error: idUser y idShop incompatibles.", ""))
	}
	// verificando que existe el type
	if b := models.ExistsTypeIdString(body.TypeId); !b {
		return c.JSON(400, config.SetResError(400, "Error: typeId no existe.", ""))
	}
	// creadndo el nuevo kind en la BBDD
	err := models.NewKindProduct(body.Name, body.TypeId, idShop)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No creado un KindProduct", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "kindProduct creado"))
}

// obtener tipo de prenda
func GetKinds(c echo.Context) error {
	// obteniendo params
	name := c.QueryParam("name")
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	// convirtiendo params
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	// consultando
	models, count, err := models.GetKinds(name, limitInt, pageInt)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not get all kinds", err.Error()))
	}
	return c.JSON(200, config.SetResJsonCount(200, "get all kinds successful", count, models))
}

// func UpdateNameKind(c echo.Context) error {
// 	// obteniendo variables
// 	body := &models.KindProduct{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// verificando que existe el kindID
// 	exist := models.ExistKindId(body.ID)
// 	if !exist {
// 		return c.JSON(400, config.SetRes(400, "Error: No existe el Kind con el Id proporcionado"))
// 	}
// 	// vericando que no existe el mismo nombre
// 	exist = models.ExistsNameKindProduct(strings.TrimSpace(body.Name))
// 	if exist {
// 		return c.JSON(400, config.SetRes(400, "Error: Existe otro Kind con el mismo nombre"))
// 	}
// 	// actualizando name en BBDD
// 	err := models.UpdateNameKind(body.ID, strings.TrimSpace(body.Name))
// 	if err != nil {
// 		return c.JSON(400, config.SetResError(400, "Error: No se pudo acutalizar el nombre de ProductKind", err.Error()))
// 	}

// 	return c.JSON(200, config.SetRes(200, "Se actualizo con exito"))
// }
// func DeleteKind(c echo.Context) error {
// 	// obteniendo variables
// 	body := &models.KindProduct{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// verificando que existe el kindID
// 	exist := models.ExistKindId(body.ID)
// 	if !exist {
// 		return c.JSON(400, config.SetRes(400, "Error: No existe el Kind con el Id proporcionado"))
// 	}
// 	// eliminando kind de BBDD
// 	err := models.DeleteKindById(body.ID)
// 	if err != nil {
// 		return c.JSON(400, config.SetResError(400, "Error: No se pudo elimina el ProductKind", err.Error()))
// 	}

// 	return c.JSON(200, config.SetRes(200, "Se elimino con exito"))
// }
