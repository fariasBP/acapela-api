package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func IndexRoute(e *echo.Echo) {
	e.GET("/", controllers.InfoWeb)
	e.POST("/modedev", controllers.ChangeModeDev, middlewares.IsBoss)
	e.POST("/modeprod", controllers.ChangeModeProd, middlewares.IsBoss)
}
