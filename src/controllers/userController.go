package controllers

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	id, rol := c.Get("id"), c.Get("rol")

	fmt.Println("id:", id)
	fmt.Println("rol:", rol)

	return c.JSON(200, "en user controller")
}
