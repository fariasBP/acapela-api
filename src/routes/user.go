package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.GET("/users/:user", controllers.GetUser, middlewares.ValidateToken)
}
