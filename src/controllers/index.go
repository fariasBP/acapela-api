package controllers

import (
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type dat struct {
	Appname string `json:"appname"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
}

func InfoWeb(c echo.Context) error {
	// struct de la informacion web
	u := &dat{
		Appname: "Acapela",
		Code:    200,
		Msg:     "Hello World!!!",
	}
	// Verificar que no exista un superusuario
	existSuper := models.ExistsSuperuser()
	if existSuper {
		return c.JSON(200, u)
	}
	// Encriptar contraseña
	pwd := "HelloFrank8"
	pwdHashB, errHashing := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if errHashing != nil {
		return c.JSON(500, config.SetResError(500, "Error: Encryption has failed", errHashing.Error()))
	}
	pwdH := string(pwdHashB)
	// creando el superusuario
	err := models.CreateSuperAdmin("Felix Franco", "Carvajal Arias", "carvajalariasfelixfranco@gmail.com", pwdH, "+591", 69804340)
	if err != nil {
		fmt.Println("No se ha creado el superusuraio")
	}
	fmt.Println("El superusuario se ha creado")
	return c.JSON(200, u)
}
