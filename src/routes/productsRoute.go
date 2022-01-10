package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ProductsRoute(e *echo.Echo) {
	e.POST("/createProduct", controllers.CreateProduct, middlewares.ValidateToken, middlewares.IsAdminOrEmpl)
}
