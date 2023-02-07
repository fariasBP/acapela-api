package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func ModelRoute(e *echo.Echo) {
	e.GET("models/all", controllers.GetAllModels)
	router := e.Group("/models", middlewares.ValidateToken)
	router.POST("/create", controllers.CreateModel,
		middlewares.ModelValidate)
	router.PUT("/update", controllers.UpdateModel,
		middlewares.IsBossOrAdmin, middlewares.ModelUpdateValidate)
	router.DELETE("/delete", controllers.DeleteModel,
		middlewares.IsBoss, middlewares.ModelDeleteValidate)
}
