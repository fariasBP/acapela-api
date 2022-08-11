package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
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
	// verificar si exista data app
	existDaApp := models.ExistsAppData()
	if !existDaApp {
		err := models.CreateApp()
		if err != nil {
			return c.JSON(500, config.SetResError(500, "No se puede crear la app", err.Error()))
		}
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
	valPhone, defined := os.LookupEnv("INIT_PHONE_ADMIN")
	if !defined {
		valPhone = "59169804340"
	}
	// convirtiendo valores
	valPhoneInt, err := strconv.Atoi(valPhone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo convertir a entero", err.Error()))
	}
	// creando el superusuario
	err = models.CreateAdminBoss(valName, valLastname, valEmail, valPhoneInt)
	if err != nil {
		fmt.Println("No se ha creado el superusuraio")
	}

	// enviar el primer mensaje whatsapp
	err = middlewares.SendWelcomeMessage(valPhone, valName)
	if err != nil {
		return c.JSON(200, config.SetResError(500, "Error: ususario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
	}

	fmt.Println("El superusuario se ha creado")
	return c.JSON(200, u)
}

func ChangeModeDev(c echo.Context) error {
	err := models.UpdDevelopingApp(true)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se puede cambiar a modo dev", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se cambio a modo dev"))
}

func ChangeModeProd(c echo.Context) error {
	err := models.UpdDevelopingApp(false)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se puede cambiar a modo prod", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se cambio a modo prod"))
}
