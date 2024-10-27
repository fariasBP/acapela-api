package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/labstack/echo/v4"
)

func PermissionRoute(e *echo.Group) {

	e.POST("/permissions/create", controllers.CreatePermission)
}
