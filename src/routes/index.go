package routes

import (
	"fmt"

	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func IndexRoute(e *echo.Echo) {
	fmt.Println("aqui")
	e.GET("/", controllers.InfoWeb)
	mode := e.Group("/mode", middlewares.ValidateToken, middlewares.IsBoss)
	mode.POST("/dev", controllers.ChangeModeDev)
	mode.POST("/prod", controllers.ChangeModeProd)
}
