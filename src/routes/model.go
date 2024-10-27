package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ModelRoute(e *echo.Group) {
	// e.GET("models/all", controllers.GetAllModels)
	e.GET("/model/models", controllers.GetModels, middlewares.GetModelsValidate)
	e.POST("/model/create", controllers.CreateModel,
		middlewares.ValidateToken, middlewares.CreateModelValidate)
	// router.PUT("/update", controllers.UpdateModel,
	// 	middlewares.IsBossOrAdmin, middlewares.ModelUpdateValidate)
	// router.DELETE("/delete", controllers.DeleteModel,
	// 	middlewares.IsBoss, middlewares.ModelDeleteValidate)
}
