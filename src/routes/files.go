package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func Files(e *echo.Group) {
	router := e.Group("/files", middlewares.ValidateToken)
	router.POST("/upload-img-product", controllers.UploadImage)
}
