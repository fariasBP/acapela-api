package middlewares

import (
	"github.com/fariasBP/acapela-api/src/config"
	"github.com/labstack/echo/v4"
)

func VerfiyPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo path actual
		var path string = ""
		for _, r := range c.Echo().Routes() {
			if r.Method == c.Request().Method && r.Path == c.Path() {
				path = r.Path
			}
		}
		if path == "" {
			return c.JSON(500, config.SetRes(200, "Problemas con el path y permisos"))
		}

		// verificando el permiso

		return c.NoContent(200)
	}
}
