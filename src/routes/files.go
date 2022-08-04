package routes

import (
	"fmt"

	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func Files(e *echo.Echo) {
	fmt.Println("aqui en files routes")
	router := e.Group("files", middlewares.ValidateToken)
	router.POST("/upload", controllers.UploadImage,
		middlewares.IsBossOrAdmin)
}
