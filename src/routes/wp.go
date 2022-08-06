package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func WPRoute(e *echo.Echo) {

	wp := e.Group("/wp", middlewares.ValidateToken)
	wp.POST("/code", controllers.SendCode)
	mass := wp.Group("/mass", middlewares.IsBoss)
	mass.POST("/newproducts", controllers.SendMassMessagesFromNewProducts)
}
