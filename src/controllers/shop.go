package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

type (
	createAdminShop struct {
		IdUser string `json:"id_user"`
	}
	clientRegisterParams struct {
		Name  string `json:"name"`
		Phone int    `json:"phone"`
	}
)

// crear una tienda
func CreateShop(c echo.Context) error {
	// obteniendo variables
	body := &middlewares.CreateShopParams{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniendo al usuario que creo
	ownerId := c.Get("id").(string)
	// creando la tienda
	if err := models.CreateShop(body.Name, ownerId, body.Description); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se ha creado la tienda", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se creo correctamente la tienda"))
}

// obtener informacion de un  tienda
func GetShop(c echo.Context) error {
	// obtengo id
	id := c.QueryParam("id")
	// consultando
	shop, err := models.GetShop(id)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se ha encontrado la tienda", err.Error()))
	}
	return c.JSON(200, config.SetResJson(200, "Se encontro la tienda", shop))
}

// obtener shops
func GetShops(c echo.Context) error {
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
	types, count, err := models.GetShops(name, limitInt, pageInt)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo obtener los types.", err.Error()))
	}
	return c.JSON(200, config.SetResJsonCount(200, "se obtuvieron los types", count, types))
}

// obtener mis shops
func GetMyShops(c echo.Context) error {
	// obteniendo valores
	name := c.QueryParam("name")
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")
	// obteniendo al usuario
	ownerId := c.Get("id").(string)
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
	types, count, err := models.GetShopByOwner(ownerId, name, limitInt, pageInt)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo obtener los types.", err.Error()))
	}
	// return c.JSON(200, config.SetResJson(200, "se obtuvieron los types", types))
	return c.JSON(200, config.SetResJsonCount(200, "se obtuvieron los types", count, types))
}

func ConvertToAdminShop(c echo.Context) error {
	// obteniendo variables
	body := &createAdminShop{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniendo shop id
	idShop := c.Get("shop").(string)
	// convirtiendo a un usuario en admin
	if err := models.ConvertToAdminShop(body.IdUser, idShop); err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se pudo convertir en administrador", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se creo un administrador correctamente"))
}

// func ClientRegistrar(c echo.Context) error {
// 	// obteniendo variables
// 	body := &models.User{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// verificar que no exista un code + phone iguales
// 	existPhone := models.ExistsPhone(body.Phone)
// 	if existPhone {
// 		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
// 	}
// 	// obteniendo shop id
// 	idShop := c.Get("shop").(string)
// 	// crear usuario en BBDD
// 	err := models.ClientRegistrar(body.Name, body.Phone)
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al cliente", err.Error()))
// 	}
// 	// Capitalizando nombre
// 	name, err := capitalize.Capitalize(body.Name)
// 	if err != nil {
// 		name = body.Name
// 	}
// 	// enviar el mesaje de bienvenida
// 	err = middlewares.SendWelcomeMessage(strconv.Itoa(body.Phone), name)
// 	if err != nil {
// 		return c.JSON(200, config.SetResError(500, "Error: usario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
// 	}

// 	return c.JSON(200, config.SetRes(200, "Cliente registrado."))
// }
