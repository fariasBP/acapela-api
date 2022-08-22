package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func Notification(e *echo.Echo) {
	router := e.Group("/notification", middlewares.ValidateToken)
	router.POST("/newproducts", controllers.NotifyNewProducts, middlewares.IsBoss)
}
