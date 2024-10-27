package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func SubTypeRoute(e *echo.Group) {
	e.GET("/subtype/subtypes", controllers.GetSubtypes, middlewares.GetSubtypesValidate)
	e.POST("/subtype/create", controllers.CreateSubtype, middlewares.ValidateToken, middlewares.CreateSubtypeValidate)
}
