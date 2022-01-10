package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func AuthRoute(e *echo.Echo) {

	e.POST("/signup", controllers.Signup, middlewares.SingupValidate)
	e.POST("/login", controllers.Login, middlewares.LoginValidate)
	e.PUT("/logout", controllers.Logout, middlewares.ValidateToken)
}
