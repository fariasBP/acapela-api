package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ProductsRoute(e *echo.Echo) {
	router := e.Group("/products", middlewares.ValidateToken)
	router.POST("/create", controllers.CreateProduct,
		middlewares.IsAdminOrEmpl)
}
