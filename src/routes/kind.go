package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func KindRoute(e *echo.Group) {
	// e.GET("kind/all", controllers.GetAllKinds)
	e.GET("/kind/kinds", controllers.GetKinds)
	router := e.Group("/kind", middlewares.ValidateToken)
	router.POST("/create", controllers.CreateKind,
		middlewares.IsOwnerShop, middlewares.KindCreateValidate)
	// router.PUT("/update", controllers.UpdateNameKind,
	// 	middlewares.IsBoss, middlewares.KindUpdateNameValidate)
	// router.DELETE("/delete", controllers.DeleteKind,
	// 	middlewares.IsBoss, middlewares.KindDeleteValidate)
}
