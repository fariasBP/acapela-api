package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateModel(c echo.Context) error {
	// obteniendo variables
	body := &middlewares.CreateModelParams{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniendo shop id
	idUser := c.Get(middlewares.ID_USER).(string)
	// creando nuevo modelo en la BBDD
	err := models.CreateModelProduct(body.Name, body.SubtypeId, idUser, body.Description)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not created new product model", err.Error()))
	}
	// incrementando used subtype
	if err := models.IncrementUsedType(body.SubtypeId); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: Se creo el Subtype pero no se incremeto el used del Type", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Modelo de producto ha sido creado"))
}

// obtener modelo de prenda
func GetModels(c echo.Context) error {
	// obteniendo params
	name := c.QueryParam("name")
	subtypeId := c.QueryParam("subtype")
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
	models, count, err := models.GetModels(name, subtypeId, limitInt, pageInt)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "error: not get all kinds", err.Error()))
	}
	return c.JSON(200, config.SetResJsonCount(200, "get all kinds successful", count, models))
}

// func GetAllModels(c echo.Context) error {
// 	models, err := models.GetAllModels()
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "error: not get all models", err.Error()))
// 	}

// 	return c.JSON(200, config.SetResJson(200, "get all models successful", models))
// }
// func UpdateModel(c echo.Context) error {
// 	// obteniendo variables
// 	body := &models.ProductModel{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// verificar si existe ID model
// 	exist := models.ExistsModelId(body.ID)
// 	if !exist {
// 		return c.JSON(400, config.SetRes(400, "Error: No existe el ID model proporcionado."))
// 	}
// 	// verificando que no existe name
// 	exist = models.ExistsNameProductModel(strings.TrimSpace(body.Name))
// 	if exist {
// 		return c.JSON(400, config.SetRes(400, "Error: ya existe el nombre."))
// 	}
// 	// verificar que existe kindID
// 	exist = models.ExistKindIdString(body.Kind)
// 	if !exist {
// 		return c.JSON(400, config.SetRes(400, "Error: ID kind no existe."))
// 	}
// 	// actualizando model
// 	err := models.UpdateModelById(body.ID, strings.TrimSpace(body.Name), body.Kind)
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "Error: no se pudo actualizar", err.Error()))
// 	}

// 	return c.JSON(200, config.SetRes(200, "actualizacion del modelo exitoso"))
// }

// func DeleteModel(c echo.Context) error {
// 	// obteniendo variables
// 	body := &models.ProductModel{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// consultamos a la BBDD
// 	err := models.DeleteModelById(body.ID)
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "no se pudo eliminar el model", err.Error()))

// 	}
// 	return c.JSON(200, config.SetRes(200, "se elimino el model exitosamente"))

// }
