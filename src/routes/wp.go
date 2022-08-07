package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/labstack/echo/v4"
)

func WPRoute(e *echo.Echo) {

	wp := e.Group("/wp")
	wp.POST("/registration", controllers.RegistrationWp)
	wp.POST("/code", controllers.SendCodeWp)
	wp.POST("/linkfacebook", controllers.SendLinkFacebookWp)
	wp.POST("/default", controllers.SendDefaultMessageWp)
	// mass := wp.Group("/mass", middlewares.IsBoss)
	// mass.POST("/newproducts", controllers.SendMassMessagesFromNewProducts)
}
