package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func User(e *echo.Echo) {
	// e.POST("/user", func(c echo.Context) error {
	// 	return c.JSON(200, "sdfasdf")
	// })
	router := e.Group("user", middlewares.ValidateToken)
	// router.POST("", controllers.GetUser)
	router.GET("/all", controllers.GetAllUsers)
	router.PUT("/name", controllers.ChangeNameUserByPhone, middlewares.IsBossOrAdmin, middlewares.NameUserValidate)
	router.GET("/profile", controllers.GetProfile)
}
