package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ShopRoute(e *echo.Group) {
	e.GET("/shop", controllers.GetShop, middlewares.GetShopValidate)
	e.GET("/shop/shops", controllers.GetShops)
	e.GET("/shop/myshops", controllers.GetMyShops, middlewares.ValidateToken)
	e.POST("/shop/create", controllers.CreateShop, middlewares.ValidateToken, middlewares.CreateShopValidate)
	// e.POST("/shop/client-register", controllers.)
	// e.PUT("/shops/to-admin", controllers.ConvertToAdminShop, middlewares.ValidateToken, middlewares.IsOwnerShop)
}
