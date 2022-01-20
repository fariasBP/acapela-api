package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func KindRoute(e *echo.Echo) {
	router := e.Group("kinds", middlewares.ValidateToken)

	router.POST("/create", controllers.CreateKind,
		middlewares.IsAdmin)
	router.GET("/all", controllers.GetAllKinds,
		middlewares.IsAdminOrEmpl)
}
