package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	// extraer datos
	kind := c.FormValue("kind")
	price, ePrice := strconv.Atoi(c.FormValue("price"))
	pricemin, ePriceM := strconv.Atoi(c.FormValue("pricemin"))
	gender, eGender := strconv.Atoi(c.FormValue("gender"))
	photo := c.FormValue("photo")
	talla := c.FormValue("talla")
	larTorso, elTo := strconv.Atoi(c.FormValue("largotorso"))
	conPecho, ecPe := strconv.Atoi(c.FormValue("contornopecho"))
	conCintura, ecCi := strconv.Atoi(c.FormValue("contornocintura"))
	conCadera, ecCa := strconv.Atoi(c.FormValue("contornocadera"))
	conSisa, ecSi := strconv.Atoi(c.FormValue("contornosisa"))
	larHombro, elHo := strconv.Atoi(c.FormValue("largohombro"))
	larManga, elMa := strconv.Atoi(c.FormValue("largomanga"))
	if ePrice != nil || ePriceM != nil || eGender != nil || elTo != nil || ecPe != nil || ecCi != nil || ecCa != nil || ecSi != nil || elHo != nil || elMa != nil {
		return c.JSON(400, config.SetRes(400, "error: values is not uint8"))
	}
	// extraer y convertir en arrays
	photos := c.FormValue("photos")
	photosArr := make([]string, 0)
	err := json.Unmarshal([]byte(photos), &photosArr)
	if err != nil {
		return c.JSON(400, config.SetRes(400, "error: photos Urls format incorrect"))
	}
	mds := c.FormValue("models")
	modelsArr := make([]string, 0)
	err = json.Unmarshal([]byte(mds), &modelsArr)
	if err != nil {
		return c.JSON(400, config.SetRes(400, "error: models ID format incorrect"))
	}
	// verficar ID's de los modelos

	// guardar nuevo producto en BBDD
	// err := models.NewProductModel(name, talla, uint8(larTorso), conPecho, conCintura, conCadera, conSisa, larHombro, larManga)
	err = models.NewProduct(kind, uint(price), uint(pricemin), uint8(gender),
		photo, photosArr, modelsArr, talla, uint8(larTorso), uint8(conPecho),
		uint8(conCintura), uint8(conCadera), uint8(conSisa), uint8(larHombro),
		uint8(larManga))
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: not created product", err.Error()))
	}
	return c.JSON(200, "creating product")
}
