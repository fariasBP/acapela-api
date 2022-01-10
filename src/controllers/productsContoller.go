package controllers

import (
	"fmt"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	name, talla := c.FormValue("name"), c.FormValue("talla")
	larTorso, elTo := strconv.Atoi(c.FormValue("largotorso"))
	conPecho, ecPe := strconv.Atoi(c.FormValue("contornopecho"))
	conCintura, ecCi := strconv.Atoi(c.FormValue("contornocintura"))
	conCadera, ecCa := strconv.Atoi(c.FormValue("contornocadera"))
	conSisa, ecSi := strconv.Atoi(c.FormValue("contornosisa"))
	larHombro, elHo := strconv.Atoi(c.FormValue("largohombro"))
	larManga, elMa := strconv.Atoi(c.FormValue("largomanga"))

	if elTo != nil || ecPe != nil || ecCi != nil || ecCa != nil || ecSi != nil || elHo != nil || elMa != nil {
		return c.JSON(400, config.SetRes(400, "error: val is not uint8"))
	}
	fmt.Println(name, talla, larTorso, conPecho, conCintura, conCadera, conSisa, larHombro, larManga)
	// err := models.NewProductModel(name, talla, uint8(larTorso), conPecho, conCintura, conCadera, conSisa, larHombro, larManga)
	return c.JSON(200, "creating product")
}
