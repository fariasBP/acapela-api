package routes

import (
	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func MessagesRoute(e *echo.Echo) {
	router := e.Group("/messages", middlewares.ValidateToken)
	router.POST("/send", controllers.SendMessageToUser, middlewares.IsBossOrAdmin)
	router.POST("/user", controllers.GetUserMessages, middlewares.IsBossOrAdmin)
	router.GET("/mailbox-not-read", controllers.GetUsersMsgNotRead, middlewares.IsBossOrAdmin)
}
