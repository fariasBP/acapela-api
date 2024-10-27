package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func TypeRoute(e *echo.Group) {
	e.GET("/type/types", controllers.GetTypes)
	e.POST("/type/create", controllers.CreateType, middlewares.ValidateToken, middlewares.IsSuperUser, middlewares.CreateTypeValidate)
}
