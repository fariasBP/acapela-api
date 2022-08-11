package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func IndexRoute(e *echo.Echo) {
	e.GET("/", controllers.InfoWeb)
	e.Group("/mode", middlewares.ValidateToken, middlewares.IsBoss)
	e.POST("/dev", controllers.ChangeModeDev)
	e.POST("/prod", controllers.ChangeModeProd)
}
