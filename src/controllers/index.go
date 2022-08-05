package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/fonini/go-capitalize/capitalize"
	"github.com/labstack/echo/v4"
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
	existSuper := models.ExistsAdiminBoss()
	if existSuper {
		return c.JSON(200, u)
	}
	// extayendo variables de entorno
	valName, defined := os.LookupEnv("INIT_NAME_ADMIN")
	if !defined {
		valName = "alex"
	}
	valLastname, defined := os.LookupEnv("INIT_LASTNAME_ADMIN")
	if !defined {
		valLastname = "siniatra"
	}
	valEmail, defined := os.LookupEnv("INIT_EMAIL_ADMIN")
	if !defined {
		valEmail = "francoxxxcarvajal@gmail.com"
	}
	valCodePhone, defined := os.LookupEnv("INIT_CODEPHONE_ADMIN")
	if !defined {
		valCodePhone = "591"
	}
	valPhone, defined := os.LookupEnv("INIT_PHONE_ADMIN")
	if !defined {
		valPhone = "69804340"
	}
	// convirtiendo valores
	valCodePhoneInt, err := strconv.Atoi(valCodePhone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo convertir a entero", err.Error()))
	}
	valPhoneInt, err := strconv.Atoi(valPhone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo convertir a entero", err.Error()))
	}
	// creando el superusuario
	err = models.CreateAdminBoss(valName, valLastname, valEmail, valCodePhoneInt, valPhoneInt)
	if err != nil {
		fmt.Println("No se ha creado el superusuraio")
	}

	// enviar el primer mensaje whatsapp
	name, err := capitalize.Capitalize(valName)
	if err != nil {
		name = valName
	}
	err = middlewares.SendWelcomeMessage(valCodePhoneInt, valPhoneInt, name)
	if err != nil {
		return c.JSON(200, config.SetResError(500, "Error: ususario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
	}

	fmt.Println("El superusuario se ha creado")
	return c.JSON(200, u)
}
