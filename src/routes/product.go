package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func Product(e *echo.Echo) {
	e.GET("/product/all", controllers.GetAllProducts)
	e.GET("/product/products", controllers.GetProducts)
	router := e.Group("/product", middlewares.ValidateToken)
	router.POST("/create", controllers.CreateProduct,
		middlewares.IsBossOrAdmin, middlewares.CreateProductValidate)
	router.PUT("/sell", controllers.SellProduct,
		middlewares.IsBossOrAdminOrEmpl, middlewares.SellProductValidate)
}
