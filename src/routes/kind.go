package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func KindRoute(e *echo.Echo) {
	e.GET("kinds/all", controllers.GetAllKinds)
	router := e.Group("kinds", middlewares.ValidateToken)
	router.POST("/create", controllers.CreateKind,
		middlewares.IsOwnerShop, middlewares.KindValidate)
	router.PUT("/update", controllers.UpdateNameKind,
		middlewares.IsBoss, middlewares.KindUpdateNameValidate)
	router.DELETE("/delete", controllers.DeleteKind,
		middlewares.IsBoss, middlewares.KindDeleteValidate)
}
