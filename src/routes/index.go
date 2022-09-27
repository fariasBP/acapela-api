package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func IndexRoute(e *echo.Echo) {
	e.GET("/", controllers.InfoWeb)
	e.GET("/dataapp", controllers.DataApp)
	mode := e.Group("/mode", middlewares.ValidateToken, middlewares.IsBoss)
	mode.POST("/dev", controllers.ChangeModeDev)
	mode.POST("/prod", controllers.ChangeModeProd)
}
