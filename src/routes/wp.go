package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/labstack/echo/v4"
)

func WPRoute(e *echo.Echo) {

	wp := e.Group("/wp")
	wp.POST("/registration", controllers.RegistrationWp)
	wp.POST("/code", controllers.SendCodeWpAndEmail)
	wp.POST("/linkfacebook", controllers.SendLinkFacebookWp)
	wp.POST("/default", controllers.SendDefaultMessageWp)
	wp.POST("/moreopt", controllers.SendMoreOptWp)
	wp.POST("/moreoptuno", controllers.SendMoreOptOneWp)
	wp.POST("/inactive", controllers.SendConfirmInactiveWp)
	wp.POST("/sleep", controllers.InactiveUserWp)
	wp.POST("/location", controllers.SendLocationMessageWp)
	wp.POST("/reactive", controllers.ReactiveUserWp)
	wp.POST("/deleteuser", controllers.SendConfirmDeleteWp)
	wp.POST("/confirmdeleteuser", controllers.DeleteUserWp)
	wp.GET("/newproducts", controllers.SendImgsNewProductsWp)
	// mass := wp.Group("/mass", middlewares.IsBoss)
	// mass.POST("/newproducts", controllers.SendMassMessagesFromNewProducts)
}
