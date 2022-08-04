package controllers

import (
	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

type (
	IDuser struct {
		Id string `json:"id"`
	}
	RolUser struct {
		Rol string
	}
)

// func GetUser(c echo.Context) error {
// 	fmt.Println("en getUser controller")
// 	id := fmt.Sprintf("%v", c.Get("id"))
// 	// rol := fmt.Sprintf("%v", c.Get("rol"))

// 	user, err := models.GetUserById(id)
// 	if err != nil {
// 		return c.JSON(400, config.SetResError(400, "error get user from models", err.Error()))
// 	}

// 	return c.JSON(200, config.SetResJson(200, "user recived info", user))
// }
func GetAllUsers(c echo.Context) error {
	users, err := models.GetUsers()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: al extraer todos los usuarios", err.Error()))
	}
	return c.JSON(200, config.SetResJson(200, "Se obtuvo todos los usuarios", users))
}
