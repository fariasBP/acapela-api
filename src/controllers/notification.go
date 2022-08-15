package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

type (
	bodyNotify struct {
		Type   string `json:"type"`
		Gender string `json:"gender"`
	}
)

func NotifyNewProducts(c echo.Context) error {
	// obteniendo variables
	body := &bodyNotify{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// obteniendo usuarios
	users, err := models.GetPhoneAndNameForNotificationsFromClients()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se pudo obtener a los usuarios, para enviar la notificacion de nuevos productos", err.Error()))
	}
	// Enviando el mensaje
	for _, v := range users {
		err = middlewares.SendNotificationFromNewProducts(strconv.Itoa(v.Phone), v.Name, body.Type, body.Gender)
		if err != nil {
			fmt.Println("Error: no se pudo enviar el mensaje de notificacion", err.Error())
		}
	}

	return c.JSON(200, config.SetRes(200, "in notify new products"))
}

// func SendOneNotification(c echo.Context) error {
// 	// obteniendo variables
// 	body := &middlewares.NotificationValues{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// enviando mensaje
// 	err := middlewares.SendMessageTextWP(body.To, body.Msg)
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "error don't send message wp twilio", err.Error()))
// 	}

// 	return c.JSON(200, config.SetRes(200, "user created and send welcome message."))
// }
// func SendMassFirstNotification(c echo.Context) error {
// 	// obteniendo variables
// 	body := &middlewares.NotificationValues{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// obteniendo numeros y nombres de los usuarios registrados
// 	users, err := models.GetPhoneAndNameForNotifications()
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "don't get users from DDBB", err.Error()))
// 	}
// 	// enviando mensajes masivos
// 	for _, v := range users {
// 		name, err := capitalize.Capitalize(v.Name)
// 		if err != nil {
// 			name = v.Name
// 		}
// 		to := v.Code + strconv.Itoa(v.Phone)
// 		err = middlewares.SendFirstNotification(name, to, "abrigos", "dama y varón")
// 		if err != nil {
// 			return c.JSON(500, config.SetResError(500, "error don't send message wp twilio to: "+v.Name, err.Error()))
// 		}
// 	}

// 	return c.JSON(200, config.SetRes(200, "in send mass notification"))

// }
