package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ModelRoute(e *echo.Echo) {
	router := e.Group("/models", middlewares.ValidateToken)

	router.POST("/create", controllers.CreateModel,
		middlewares.IsAdminOrEmpl,
		middlewares.NewProductModelValidate)
	router.GET("/all", controllers.GetAllModels,
		middlewares.IsAdminOrEmpl)
}
