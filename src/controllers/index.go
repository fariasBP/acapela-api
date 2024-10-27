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
type (
	App struct {
		Name       string `json:"name" bson:"name,omitempty"`
		Developing bool   `json:"developing" bson:"developing,omitempty"`
		Version    string `json:"version" bson:"version,omitempty"`
	}
)

// Se obtiene los datos actuales de la app
/* Por ejemplo:
- version actual
- si esta en desarrollo
*/
func DataApp(c echo.Context) error {
	// obteniendo dataapp de .env
	name, defined := os.LookupEnv("NAMEAPP")
	version, _ := os.LookupEnv("VERSIONAPP")
	developing, _ := os.LookupEnv("DEVELOPINGMODE")
	// convirtiendo valores
	developingMode, err := strconv.ParseBool(developing)
	if err != nil {
		developingMode = false
	}
	// valores
	dat := &App{
		Name:       name,
		Version:    version,
		Developing: developingMode,
	}
	// Verificar que no exista un superusuario
	existSuper := models.ExistsAdiminBoss()
	if existSuper {
		// return c.JSON(200, u)
		return c.JSON(200, config.SetResJson(200, "Ya se inicio la app.", dat))
	}
	// extayendo variables de entorno
	valName, defined := os.LookupEnv("INIT_NAME_ADMIN")
	if !defined {
		valName = "franco"
	}
	valLastname, defined := os.LookupEnv("INIT_LASTNAME_ADMIN")
	if !defined {
		valLastname = "carvajal"
	}
	valEmail, defined := os.LookupEnv("INIT_EMAIL_ADMIN")
	if !defined {
		valEmail = "carvajalariasfelixfranco@gmail.com"
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

	return c.JSON(200, config.SetResJson(200, "Se creo superusuario y se inicio la app.", dat))
}
