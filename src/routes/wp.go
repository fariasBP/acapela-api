package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func WPRoute(e *echo.Echo) {
	// e.GET("/webhook", controllers.FirstConfirmWithApiWhatsappForReciveMessage)
	// e.POST("/webhook", controllers.RecivedMessagesWhatsapp)

	wp := e.Group("/wp", middlewares.ValidateToken)
	mass := wp.Group("/mass", middlewares.IsBoss)
	mass.POST("/newproducts", controllers.SendMassMessagesFromNewProducts)
}
