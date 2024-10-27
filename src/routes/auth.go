package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func AuthRoute(e *echo.Group) {
	// e.POST("/signup", controllers.Signup, middlewares.SingupValidate)
	e.POST("/login", controllers.Login, middlewares.LoginValidate)
	e.POST("/codeauth", controllers.SendCodeWpAndEmail)
	e.GET("/validate", controllers.ValidateToken)
	e.GET("/issuperuser", controllers.IsSuperUser, middlewares.ValidateToken)
	// router := e.Group("/registrar", middlewares.ValidateToken, middlewares.RegisterValidator)
	// router.POST("/client", controllers.ClientRegistrar, middlewares.IsBossOrAdminOrEmpl)
	// router.POST("/employe", controllers.EmployeRegistrar, middlewares.IsBossOrAdmin)
	// router.POST("/adminemploye", controllers.AdminEmployeRegistrar, middlewares.IsBoss)
}
