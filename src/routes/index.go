package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/labstack/echo/v4"
)

func IndexRoute(e *echo.Group) {
	// e.GET("/inf", controllers.InfoWeb)
	e.GET("/dataapp", controllers.DataApp)
}
