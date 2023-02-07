package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ShopRoute(e *echo.Echo) {
	e.POST("/shop/create", controllers.CreateShop, middlewares.ValidateToken)
	e.PUT("/shops/to-admin", controllers.ConvertToAdminShop, middlewares.ValidateToken, middlewares.IsOwnerShop)
}
