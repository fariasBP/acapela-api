package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	// obteniendo variables
	body := &models.Product{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()

	//verificar ID de kind
	exist := models.ExistKindIdString(body.Kind)
	if !exist {
		return c.JSON(400, config.SetResError(400, "Error: No existe un KindProduct ingresado", "Kindproduct does't exist"))
	}

	// verficar ID's de models
	for _, v := range body.Models {
		exist = models.ExistsModelIdString(v)
		if !exist {
			return c.JSON(400, config.SetResError(400, "Error: No existe un ModelProduct ingresado", "Modelproduct does't exist"))
		}
	}

	// guardar nuevo producto en BBDD
	err := models.NewProduct(body.Price, body.PriceMin, body.Photos,
		body.Kind, body.Models, body.Gender, body.Size,
		body.ModelQuality, body.MaterialQuality)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: not created product", err.Error()))
	}
	// verificar si es un nuevo producto
	// extrayendo la fecha del ultimo product
	yearP, monthP, dayP, err := models.GetLastDateFromProduct()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo extraer la fecha del ultimo product", err.Error()))
	}
	// extrayendo la fecha del ultimo created
	yearC, monthC, dayC, _ := models.GetLastDateFromCreated()
	// igualando fechas
	if yearP != yearC || monthP != monthC || dayP != dayC {
		// si es diferente se crea un created
		err := models.NewCreated()
		if err != nil {
			return c.JSON(500, config.SetResError(500, "Error: not created LastCreated", err.Error()))
		}
	}

	return c.JSON(200, config.SetRes(200, "Producto creado"))
}
func GetAllProducts(c echo.Context) error {

	products, err := models.GetAllProducts()
	if err != nil {
		c.JSON(500, config.SetResError(500, "error: not get products", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "get all products successful", products))
}

func GetProducts(c echo.Context) error {
	// oteniendo los querys
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		c.JSON(500, config.SetResError(500, "error: page no es un numero", err.Error()))
	}
	// estableciendo limite
	limit := 2
	// consultando a BBDD
	products, err := models.GetProducts(limit, (page-1)*limit)
	if err != nil {
		c.JSON(500, config.SetResError(500, "error: not get products", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "get products succesful", products))
}

func GetNewProducts(c echo.Context) error {
	// consultando a BBDD
	products, err := models.GetNewProducts()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se pudo obtener los nuevos productos de BBDD", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "Se obtuvo los nuevos productos con exito", products))
}

func SellProduct(c echo.Context) error {
	// obteniendo variables
	body := &models.Product{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe la id del producto
	exist := models.ExistProductId(body.ID)
	if !exist {
		return c.JSON(400, config.SetRes(400, "Error: no existe el product id"))
	}
	// verificar que existe el vendedor
	exist = models.ExistsSellerIDStr(body.Seller)
	if !exist {
		return c.JSON(400, config.SetRes(400, "Error: no existe el vendedor"))
	}
	// verificar si existe buyer
	exist = models.ExistsBuyerIDStr(body.Buyer)
	if exist {
		err := models.SellProductWithBuyer(body.ID, body.SellPrice, body.Seller, body.Buyer)
		if err != nil {
			return c.JSON(500, config.SetResError(500, "Error: al actualizar a vendido", err.Error()))
		}

	} else {
		err := models.SellProductWithoutBuyer(body.ID, body.SellPrice, body.Seller)
		if err != nil {
			return c.JSON(500, config.SetResError(500, "Error: al actualizar a vendido", err.Error()))
		}
	}
	// estableciendo como vendido en BBDD

	return c.JSON(200, config.SetRes(200, "sell product succeful"))
}
