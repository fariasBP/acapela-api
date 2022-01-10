package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/labstack/echo/v4"
)

func IndexRoute(e *echo.Echo) {
	e.GET("/", controllers.InfoWeb)
}
