package controllers

import (
	"github.com/fariasBP/acapela-api/src/config"
	"github.com/labstack/echo/v4"
)

func NotifyNewProducts(c echo.Context) error {
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
